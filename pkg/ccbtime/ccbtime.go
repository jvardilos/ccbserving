package ccbtime

import (
	"fmt"
	"time"
)

func FormatDate(now time.Time) string {
	return fmt.Sprintf("%d-%d-%d", now.Year(), now.Month(), now.Day())
}

func FormatLastMonthDate(now time.Time) string {
	fmt.Println("Input date.   :", now)
	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	fmt.Println("Input date fmt:", now)
	lastMonth := now.AddDate(0, -1, 0)
	fmt.Println("Last month fmt:", lastMonth)
	return fmt.Sprintf("%d-%d-%d", lastMonth.Year(), lastMonth.Month(), lastMonth.Day())
}
