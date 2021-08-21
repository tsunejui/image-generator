package pkg

import "strings"

var (
	strYes = []string{"yes", "y"}
	strNo  = []string{"no", "n"}
)

func CheckYes(s string) bool {
	for _, y := range strYes {
		if strings.EqualFold(s, y) {
			return true
		}
	}
	return false
}

func CheckNo(s string) bool {
	for _, n := range strNo {
		if strings.EqualFold(s, n) {
			return true
		}
	}
	return false
}
