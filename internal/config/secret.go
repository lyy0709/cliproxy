/*
 * 文件作用：安全密钥管理，负责初始化 JWT 密钥与数据加密密钥
 * 负责功能：
 *   - 首次启动自动生成密钥
 *   - 将密钥写回配置文件
 * 重要程度：⭐⭐⭐⭐ 重要（安全基础）
 */
package config

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

const DefaultJWTSecretPlaceholder = "cli-proxy-jwt-secret-change-in-production"

// EnsureSecrets 检查并生成必要的安全密钥
func EnsureSecrets(configPath string) error {
	if Cfg == nil {
		return errors.New("配置未加载")
	}

	// 环境变量优先时不写回文件
	jwtFromEnv := os.Getenv("JWT_SECRET") != ""
	dataKeyFromEnv := os.Getenv("DATA_ENCRYPTION_KEY") != ""

	needUpdate := false

	cfgMap, mode, err := loadConfigMap(configPath)
	if err != nil {
		return err
	}

	if !jwtFromEnv && (Cfg.JWT.Secret == "" || Cfg.JWT.Secret == DefaultJWTSecretPlaceholder) {
		secret, err := generateRandomSecret(32)
		if err != nil {
			return err
		}
		setConfigValue(cfgMap, []string{"jwt", "secret"}, secret)
		Cfg.JWT.Secret = secret
		needUpdate = true
	}

	if !dataKeyFromEnv && Cfg.Security.DataKey == "" {
		key, err := generateRandomSecret(32)
		if err != nil {
			return err
		}
		setConfigValue(cfgMap, []string{"security", "data_key"}, key)
		Cfg.Security.DataKey = key
		needUpdate = true
	}

	if needUpdate {
		if err := writeConfigMap(configPath, mode, cfgMap); err != nil {
			return fmt.Errorf("写入配置文件失败: %w", err)
		}
	}

	return nil
}

func generateRandomSecret(bytesLen int) (string, error) {
	if bytesLen <= 0 {
		return "", errors.New("密钥长度无效")
	}
	buf := make([]byte, bytesLen)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}

func loadConfigMap(path string) (map[string]interface{}, os.FileMode, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, 0, fmt.Errorf("读取配置文件失败: %w", err)
	}
	info, err := os.Stat(path)
	if err != nil {
		return nil, 0, fmt.Errorf("读取配置文件权限失败: %w", err)
	}

	cfgMap := make(map[string]interface{})
	if err := yaml.Unmarshal(data, &cfgMap); err != nil {
		return nil, 0, fmt.Errorf("解析配置文件失败: %w", err)
	}
	return cfgMap, info.Mode().Perm(), nil
}

func writeConfigMap(path string, mode os.FileMode, cfg map[string]interface{}) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("序列化配置失败: %w", err)
	}
	return os.WriteFile(path, data, mode)
}

func setConfigValue(root map[string]interface{}, keys []string, value string) {
	if len(keys) == 0 {
		return
	}
	current := root
	for i, key := range keys {
		if i == len(keys)-1 {
			current[key] = value
			return
		}
		next := toStringMap(current[key])
		if next == nil {
			next = make(map[string]interface{})
			current[key] = next
		}
		current = next
	}
}

func toStringMap(value interface{}) map[string]interface{} {
	switch m := value.(type) {
	case map[string]interface{}:
		return m
	case map[interface{}]interface{}:
		out := make(map[string]interface{})
		for k, v := range m {
			ks, ok := k.(string)
			if !ok {
				continue
			}
			out[ks] = v
		}
		return out
	default:
		return nil
	}
}
