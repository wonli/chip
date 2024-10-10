package chip

import (
	"encoding/json"
	"fmt"
	"html"
	"html/template"
	"math"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"
)

type Format struct {
	chip *Chip
}

func (f Format) Str2slice(str string) []string {
	var data []string
	sl := strings.Split(str, ",")
	for _, s := range sl {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}

		data = append(data, s)
	}

	return data
}

func (f Format) Substr(s string, pos int) string {
	cleanString := f.Strips(s)
	if utf8.RuneCountInString(cleanString) <= pos {
		return html.UnescapeString(cleanString)
	}

	truncated := ""
	for _, runeValue := range cleanString {
		if pos == 0 {
			break
		}

		truncated += string(runeValue)
		pos--
	}

	return html.UnescapeString(truncated)
}

func (f Format) Strips(s string) string {
	re := regexp.MustCompile(`<[^>]*>`)
	result := re.ReplaceAllString(s, "")

	// Remove control characters like newlines.
	result = strings.ReplaceAll(result, "\n", "")
	result = strings.ReplaceAll(result, "\r", "")
	result = strings.ReplaceAll(result, "\t", "")
	return result
}

func (f Format) Res(assets string) string {
	rr := f.chip.GetEventRoute()
	directoryPath := rr.Route
	if strings.HasPrefix(assets, "/") || directoryPath == "/" {
		return assets
	}

	// 去除assets路径前的"./"
	if strings.HasPrefix(assets, "./") {
		assets = assets[2:]
	}

	if lastSlash := strings.LastIndex(directoryPath, "/"); lastSlash != -1 {
		directoryPath = directoryPath[:lastSlash]
	}

	// 计算路径中包含的"/"数量，每一个"/"代表一级目录
	upLevels := strings.Count(directoryPath, "/")
	relativePath := strings.Repeat("../", upLevels)
	if strings.HasPrefix(relativePath, "../") && strings.HasPrefix(assets, "./") {
		assets = assets[2:]
	}

	resultPath := "./" + relativePath + assets
	return resultPath
}

func (f Format) StripScheme(url string) string {
	if strings.HasPrefix(url, "http:") {
		return strings.TrimPrefix(url, "http:")
	}

	if strings.HasPrefix(url, "https:") {
		return strings.TrimPrefix(url, "https:")
	}

	return url
}

func (f Format) ToHTTPS(url string) string {
	if strings.HasPrefix(url, "http:") {
		return "https:" + strings.TrimPrefix(url, "http:")
	}

	return url
}

func (f Format) FriTime(pastTime time.Time) string {
	currentTime := time.Now()
	diff := currentTime.Sub(pastTime).Seconds()

	if diff <= 31536000 {
		timeFrames := []struct {
			Seconds int64
			Label   string
		}{
			{31536000, "年前"},
			{2592000, "个月前"},
			{604800, "星期前"},
			{86400, "天前"},
			{3600, "小时前"},
			{60, "分钟前"},
			{1, "秒前"},
		}

		for _, frame := range timeFrames {
			c := math.Floor(diff / float64(frame.Seconds))
			remainder := diff - c*float64(frame.Seconds)

			if remainder > 0.80*float64(frame.Seconds) {
				c++
				return fmt.Sprintf("约%d%s", int(c), frame.Label)
			} else if c != 0 {
				return fmt.Sprintf("%d%s", int(c), frame.Label)
			}
		}
	}

	return pastTime.Format(time.DateTime)
}

func (f Format) FriNumber(value float64) string {
	var unit string
	var ffNumber string
	if value >= 1000000 {
		unit = "M"
		ffNumber = fmt.Sprintf("%.2f", value/1000000)
	} else if value >= 1000 {
		unit = "K"
		ffNumber = fmt.Sprintf("%.2f", value/1000)
	} else {
		ffNumber = fmt.Sprintf("%.2f", value)
	}

	if strings.Contains(ffNumber, ".") {
		ffNumber = strings.TrimRight(ffNumber, "0")
		ffNumber = strings.TrimRight(ffNumber, ".")
	}

	return ffNumber + unit
}

func (f Format) Url(routeName string, params ...any) string {
	var router *Route
	for _, r := range f.chip.config.Routes {
		if r.Name == routeName {
			router = r
			break
		}
	}

	if router == nil {
		return routeName
	}

	re := regexp.MustCompile(`\{[^}]+\}`)
	url := router.Route
	url = re.ReplaceAllString(url, "%s")

	count := strings.Count(url, "%s")
	stringVars := make([]any, count)
	if count > 0 && len(params) >= count {
		for i := 0; i < count; i++ {
			stringVars[i] = fmt.Sprintf("%v", params[i])
		}
	}

	return f.chip.config.BaseLinkPath + fmt.Sprintf(url, stringVars...)
}

func (f Format) UnixDate(unix int, format string) string {
	t := time.Unix(int64(unix), 0)
	return t.Format(f.toDateTimeFormat(format))
}

func (f Format) DateTime(t time.Time, format string) string {
	return t.Format(f.toDateTimeFormat(format))
}

func (f Format) DateFormat(t time.Time, layout string) string {
	return t.Format(layout)
}

func (f Format) FloatFormat(number float64, precision int) string {
	return fmt.Sprintf("%.*f", precision, number)
}

func (f Format) Capitalize(text string) string {
	if len(text) == 0 {
		return ""
	}
	return strings.ToUpper(text[:1]) + text[1:]
}

func (f Format) EscapeHTML(input string) string {
	return template.HTMLEscapeString(input)
}

func (f Format) UnescapeHTML(input string) string {
	return html.UnescapeString(input)
}

func (f Format) ToJsonSlice(data string) []string {
	var result []string
	err := json.Unmarshal([]byte(data), &result)
	if err != nil {
		return nil
	}

	return result
}

func (f Format) Kb(kb int64) string {
	const (
		KB = 1024
		MB = 1024 * KB
		GB = 1024 * MB
	)

	var result string

	kbf := kb * 1024
	switch {
	case kbf >= GB:
		result = fmt.Sprintf("%.2f GB", float64(kbf)/GB)
	case kbf >= MB:
		result = fmt.Sprintf("%.2f MB", float64(kbf)/MB)
	case kbf >= KB:
		result = fmt.Sprintf("%.2f KB", float64(kbf)/KB)
	default:
		result = fmt.Sprintf("%d Bytes", kbf)
	}

	return result
}

func (f Format) toDateTimeFormat(format string) string {
	format = strings.ToUpper(format)
	replacements := map[string]string{
		"Y": "2006",
		"M": "01",
		"D": "02",
		"H": "15",
		"m": "04",
		"S": "05",
	}

	// 替换所有自定义格式
	for k, v := range replacements {
		format = strings.ReplaceAll(format, k, v)
	}

	return format
}
