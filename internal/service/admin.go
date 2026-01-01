/*
 * 文件作用：管理员认证服务，处理管理员登录和密码管理逻辑
 * 负责功能：
 *   - 管理员登录验证（仅数据库存储）
 *   - JWT Token 生成
 *   - 密码修改
 *   - 首次登录强制修改密码
 * 重要程度：⭐⭐⭐⭐ 重要（认证核心服务）
 * 依赖模块：config, utils, bcrypt, repository
 */
package service

import (
	"errors"
	"time"

	"cli-proxy/internal/model"
	"cli-proxy/internal/repository"
	"cli-proxy/pkg/utils"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AdminService 管理员服务
type AdminService struct {
	repo *repository.AdminConfigRepository
}

// NewAdminService 创建管理员服务实例
func NewAdminService() *AdminService {
	return &AdminService{
		repo: repository.NewAdminConfigRepository(repository.DB),
	}
}

// LoginRequest 登录请求
type AdminLoginRequest struct {
	Username    string `json:"username" binding:"required"`
	Password    string `json:"password" binding:"required"`
	CaptchaId   string `json:"captcha_id"`
	CaptchaCode string `json:"captcha_code"`
}

// LoginResponse 登录响应
type AdminLoginResponse struct {
	Token              string `json:"token"`
	Username           string `json:"username"`
	Role               string `json:"role"`
	MustChangePassword bool   `json:"must_change_password"`
}

// AdminStatusResponse 管理员状态响应
type AdminStatusResponse struct {
	Username           string     `json:"username"`
	MustChangePassword bool       `json:"must_change_password"`
	PasswordChangedAt  *time.Time `json:"password_changed_at"`
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

// ForceChangePasswordRequest 首次登录强制修改密码请求
type ForceChangePasswordRequest struct {
	NewPassword     string `json:"new_password" binding:"required"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
}

// Login 管理员登录
func (s *AdminService) Login(req *AdminLoginRequest) (*AdminLoginResponse, error) {
	adminConfig, err := s.repo.GetByUsername(req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户名或密码错误")
		}
		return nil, errors.New("管理员认证失败")
	}

	// 数据库中存在管理员配置，使用数据库验证
	if !adminConfig.VerifyPassword(req.Password) {
		return nil, errors.New("用户名或密码错误")
	}

	// 生成 JWT Token
	token, err := utils.GenerateToken(0, req.Username, "admin")
	if err != nil {
		return nil, errors.New("生成Token失败")
	}

	return &AdminLoginResponse{
		Token:              token,
		Username:           req.Username,
		Role:               "admin",
		MustChangePassword: adminConfig.MustChangePassword,
	}, nil
}

// GetStatus 获取管理员状态
func (s *AdminService) GetStatus(username string) (*AdminStatusResponse, error) {
	adminConfig, err := s.repo.GetByUsername(username)
	if err != nil {
		return nil, errors.New("管理员配置不存在")
	}

	return &AdminStatusResponse{
		Username:           adminConfig.Username,
		MustChangePassword: adminConfig.MustChangePassword,
		PasswordChangedAt:  adminConfig.PasswordChangedAt,
	}, nil
}

// ChangePassword 修改密码（需要验证旧密码）
func (s *AdminService) ChangePassword(username string, req *ChangePasswordRequest) error {
	// 验证新密码强度
	if !model.ValidatePasswordStrength(req.NewPassword) {
		return errors.New("密码强度不足：至少8位，包含字母和数字")
	}

	adminConfig, err := s.repo.GetByUsername(username)
	if err != nil {
		return errors.New("管理员配置不存在")
	}
	if adminConfig.MustChangePassword {
		return errors.New("首次登录必须先完成强制修改密码")
	}

	// 验证旧密码
	if !adminConfig.VerifyPassword(req.OldPassword) {
		return errors.New("当前密码错误")
	}

	// 设置新密码
	if err := adminConfig.SetPassword(req.NewPassword); err != nil {
		return errors.New("密码加密失败")
	}

	// 更新密码
	return s.repo.UpdatePassword(adminConfig.ID, adminConfig.PasswordHash, false)
}

// ForceChangePassword 首次登录强制修改密码
func (s *AdminService) ForceChangePassword(username string, req *ForceChangePasswordRequest) error {
	// 验证两次密码是否一致
	if req.NewPassword != req.ConfirmPassword {
		return errors.New("两次输入的密码不一致")
	}

	// 验证新密码强度
	if !model.ValidatePasswordStrength(req.NewPassword) {
		return errors.New("密码强度不足：至少8位，包含字母和数字")
	}

	adminConfig, err := s.repo.GetByUsername(username)
	if err != nil {
		return errors.New("管理员配置不存在")
	}
	if !adminConfig.MustChangePassword {
		return errors.New("无需强制修改密码，请使用修改密码功能")
	}

	// 设置新密码
	if err := adminConfig.SetPassword(req.NewPassword); err != nil {
		return errors.New("密码加密失败")
	}

	// 更新密码，标记已修改
	return s.repo.UpdatePassword(adminConfig.ID, adminConfig.PasswordHash, false)
}

// HashPassword 生成密码哈希（用于手动生成管理员密码）
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
