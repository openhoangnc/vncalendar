package vncalendar

import "fmt"

// padd is a helper function to pad a number with a leading zero if it is less than 10
func padd(digits int) string {
	if digits > 9 {
		return fmt.Sprintf("%d", digits)
	}
	return fmt.Sprintf("0%d", digits)
}
