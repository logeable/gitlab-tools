package output

import "time"

// FormatToLocalTime 将时间格式化为本地时间字符串
func FormatToLocalTime(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.In(time.Local).Format("2006-01-02 15:04:05")
}
