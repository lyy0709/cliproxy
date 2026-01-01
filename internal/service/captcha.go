/*
 * 文件作用：验证码服务，生成和验证图片验证码
 * 负责功能：
 *   - 生成图片验证码
 *   - 验证码验证
 *   - 验证码存储和过期管理
 * 重要程度：⭐⭐⭐⭐ 重要（安全核心功能）
 * 依赖模块：base64Captcha
 */
package service

import (
	"image/color"
	"sync"
	"time"

	"github.com/mojocn/base64Captcha"
)

// CaptchaService 验证码服务
type CaptchaService struct {
	store  *MemoryStore
	driver *base64Captcha.DriverString
}

// MemoryStore 内存存储（带过期时间）
type MemoryStore struct {
	data sync.Map
	ttl  time.Duration
}

type captchaItem struct {
	value     string
	expiredAt time.Time
}

// NewMemoryStore 创建内存存储
func NewMemoryStore(ttl time.Duration) *MemoryStore {
	store := &MemoryStore{ttl: ttl}
	// 启动清理协程
	go store.cleanupLoop()
	return store
}

// Set 设置验证码
func (s *MemoryStore) Set(id string, value string) error {
	s.data.Store(id, captchaItem{
		value:     value,
		expiredAt: time.Now().Add(s.ttl),
	})
	return nil
}

// Get 获取验证码（不删除）
func (s *MemoryStore) Get(id string, clear bool) string {
	v, ok := s.data.Load(id)
	if !ok {
		return ""
	}
	item := v.(captchaItem)
	if time.Now().After(item.expiredAt) {
		s.data.Delete(id)
		return ""
	}
	if clear {
		s.data.Delete(id)
	}
	return item.value
}

// Verify 验证验证码
func (s *MemoryStore) Verify(id, answer string, clear bool) bool {
	v := s.Get(id, clear)
	return v != "" && v == answer
}

// cleanupLoop 定期清理过期验证码
func (s *MemoryStore) cleanupLoop() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		s.data.Range(func(key, value interface{}) bool {
			item := value.(captchaItem)
			if time.Now().After(item.expiredAt) {
				s.data.Delete(key)
			}
			return true
		})
	}
}

var (
	captchaService     *CaptchaService
	captchaServiceOnce sync.Once
)

// GetCaptchaService 获取验证码服务单例
func GetCaptchaService() *CaptchaService {
	captchaServiceOnce.Do(func() {
		store := NewMemoryStore(5 * time.Minute) // 5分钟过期

		// 配置验证码驱动
		driver := &base64Captcha.DriverString{
			Height:          60,
			Width:           200,
			NoiseCount:      0,
			ShowLineOptions: base64Captcha.OptionShowHollowLine | base64Captcha.OptionShowSlimeLine,
			Length:          4,
			Source:          "0123456789abcdefghijklmnopqrstuvwxyz",
			BgColor: &color.RGBA{
				R: 245,
				G: 245,
				B: 245,
				A: 255,
			},
		}

		captchaService = &CaptchaService{
			store:  store,
			driver: driver,
		}
	})
	return captchaService
}

// CaptchaResponse 验证码响应
type CaptchaResponse struct {
	CaptchaId string `json:"captcha_id"`
	Image     string `json:"image"`
	Enabled   bool   `json:"enabled"`
}

// Generate 生成验证码
func (s *CaptchaService) Generate() (*CaptchaResponse, error) {
	captcha := base64Captcha.NewCaptcha(s.driver, s.store)

	id, b64s, _, err := captcha.Generate()
	if err != nil {
		return nil, err
	}

	return &CaptchaResponse{
		CaptchaId: id,
		Image:     b64s,
		Enabled:   true,
	}, nil
}

// Verify 验证验证码
func (s *CaptchaService) Verify(id, code string) bool {
	if id == "" || code == "" {
		return false
	}
	return s.store.Verify(id, code, true)
}
