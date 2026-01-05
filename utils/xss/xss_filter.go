package xss

import "github.com/microcosm-cc/bluemonday"

// 全局复用 Policy，避免每次调用函数都重新创建对象，提升性能
var ugcPolicy = bluemonday.UGCPolicy()

// FilterSanitize XSS 过滤工具函数
// 输入: 原始字符串
// 输出: 清洗后的字符串 (移除了 <script> 等危险标签)
func FilterSanitize(content string) string {
	return ugcPolicy.Sanitize(content)
}
