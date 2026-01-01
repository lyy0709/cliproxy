/*
 * 文件作用：请求/响应体读取工具，限制最大读取大小
 * 负责功能：
 *   - 防止超大请求体导致内存压力
 */
package utils

import (
	"errors"
	"io"
)

const (
	MaxRequestBodyBytes  int64 = 50 << 20 // 50MB，支持多模态请求（多张图片）
	MaxResponseBodyBytes int64 = 50 << 20 // 50MB
)

var ErrBodyTooLarge = errors.New("body too large")

// ReadAllWithLimit 在限制大小内读取数据
func ReadAllWithLimit(r io.Reader, maxBytes int64) ([]byte, error) {
	if r == nil {
		return nil, nil
	}
	if maxBytes <= 0 {
		return io.ReadAll(r)
	}

	limited := io.LimitReader(r, maxBytes+1)
	data, err := io.ReadAll(limited)
	if err != nil {
		return nil, err
	}
	if int64(len(data)) > maxBytes {
		return nil, ErrBodyTooLarge
	}
	return data, nil
}
