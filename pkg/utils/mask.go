/*
 * 文件作用：敏感信息脱敏工具
 * 负责功能：
 *   - 对密码/Token 等敏感字段进行部分遮蔽
 */
package utils

import "strings"

// MaskSensitive 对敏感字符串进行中间遮蔽
func MaskSensitive(value string, keepStart, keepEnd int) string {
	if value == "" {
		return ""
	}
	if keepStart < 0 {
		keepStart = 0
	}
	if keepEnd < 0 {
		keepEnd = 0
	}
	if keepStart+keepEnd >= len(value) {
		return strings.Repeat("*", min(len(value), 8))
	}
	maskLen := len(value) - keepStart - keepEnd
	return value[:keepStart] + strings.Repeat("*", maskLen) + value[len(value)-keepEnd:]
}

// MaskPassword 密码脱敏（保留首尾各1位）
func MaskPassword(value string) string {
	return MaskSensitive(value, 1, 1)
}

// MaskToken Token 脱敏（保留首尾各4位）
func MaskToken(value string) string {
	return MaskSensitive(value, 4, 4)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
