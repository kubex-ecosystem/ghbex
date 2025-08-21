// Package utils provides utility functions for the application.
package utils

import "time"

func Cutoff(days int) time.Time {
	if days <= 0 {
		return time.Time{}
	}
	return time.Now().AddDate(0, 0, -days)
}
