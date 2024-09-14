package chip

import (
	"fmt"
	"math"
	"time"
)

func TimeSince(start, end time.Time) string {
	if start.IsZero() {
		return ""
	}

	duration := end.Sub(start)
	seconds := int(duration.Seconds())
	minutes := int(duration.Minutes())
	hours := int(duration.Hours())
	days := hours / 24

	if seconds < 1 {
		return fmt.Sprintf("%.2f秒", duration.Seconds())
	} else if seconds < 60 {
		return fmt.Sprintf("%d秒", seconds)
	} else if minutes < 60 {
		return fmt.Sprintf("%d分%d秒", minutes, seconds%60)
	} else if hours < 24 {
		return fmt.Sprintf("%d小时%d分钟", hours, minutes%60)
	} else {
		return fmt.Sprintf("%d天%d小时%d分钟", days, hours%24, minutes%60)
	}
}

func FormatBites(size float64) string {
	if size <= 0 {
		return ""
	}

	unit := []string{"B", "KB", "MB", "GB", "TB", "PB"}
	s := math.Floor(math.Log(size) / math.Log(1024))
	i := int(s)

	if i < len(unit) {
		return fmt.Sprintf("%.2f %s", size/math.Pow(1024, s), unit[i])
	}

	return fmt.Sprintf("%f %s", size, unit[0])
}
