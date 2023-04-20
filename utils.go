package ginche

import "time"

// String Converts a string to a string pointer
func String(str string) *string {
	return &str
}

// Int Converts int to int pointer
func Int(i int) *int {
	return &i
}

// Int64 Converts int64 to int64 pointer
func Int64(i int64) *int64 {
	return &i
}

// Float64 Converts float64 to float64 pointer
func Float64(f float64) *float64 {
	return &f
}

// Bool Converts bool to bool pointer
func Bool(b bool) *bool {
	return &b
}

// Duration Converts time.Duration to time.Duration pointer
func Duration(t time.Duration) *time.Duration {
	return &t
}
