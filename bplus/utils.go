package bplus

import "fmt"

func keyString(k int) string {
	if k >= 'A' && k <= 'Z' || k >= 'a' && k <= 'z' {
		return fmt.Sprintf("%c", k)
	}
	return fmt.Sprintf("%d", k)
}
