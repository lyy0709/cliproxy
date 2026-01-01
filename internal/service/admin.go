/*
 * 文件作用：管理员认证服务，处理管理员登录逻辑
 * 负责功能：
 *   - 管理员登录验证
 *   - JWT Token 生成
 *   - 从配置文件读取管理员凭证
 * 重要程度：⭐⭐⭐⭐ 重要（认证核心服务）
 * 依赖模块：config, utils, bcrypt
 */
package service

import (
	"errors"

	"cli-proxy/internal/config"
	"cli-proxy/pkg/utils"

	"golang.org/x/crypto/bcrypt"
)

// AdminService 管理员服务
type AdminService struct{}

// NewAdminService 创建管理员服务实例
func NewAdminService() *AdminService {
	return &AdminService{}
}

// LoginRequest 登录请求
type AdminLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应
type AdminLoginResponse struct {
	Token    string `json:"token"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

// Login 管理员登录
func (s *AdminService) Login(req *AdminLoginRequest) (*AdminLoginResponse, error) {
	// 验证用户名
	if req.Username != config.Cfg.Admin.Username {
		return nil, errors.New("用户名或密码错误")
	}

	// 验证密码
	// 如果配置的密码以 $2a$ / $2b$ / $2y$ 开头，说明是 bcrypt 加密的
	storedPassword := config.Cfg.Admin.Password
	if len(storedPassword) > 4 && (storedPassword[:4] == "$2a$" || storedPassword[:4] == "$2b$" || storedPassword[:4] == "$2y$") {
		// bcrypt 加密密码验证
		if err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(req.Password)); err != nil {
			return nil, errors.New("用户名或密码错误")
		}
	} else {
		// 明文密码验证（不推荐，仅用于开发测试）
		if req.Password != storedPassword {
			return nil, errors.New("用户名或密码错误")
		}
	}

	// 生成 JWT Token
	// 使用 ID=0 表示管理员（无数据库用户记录）
	token, err := utils.GenerateToken(0, req.Username, "admin")
	if err != nil {
		return nil, errors.New("生成Token失败")
	}

	return &AdminLoginResponse{
		Token:    token,
		Username: req.Username,
		Role:     "admin",
	}, nil
}

// HashPassword 生成密码哈希（用于生成配置文件中的密码）
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
