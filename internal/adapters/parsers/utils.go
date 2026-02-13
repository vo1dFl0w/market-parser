package parsers

import "strconv"

func ParseStringToInteger(s string) (int, error) {
	return strconv.Atoi(s)
}
