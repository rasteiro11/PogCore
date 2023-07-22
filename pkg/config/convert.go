package config

import "strconv"

func convertToInt(value string) int {
	v, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}

	return v
}

func convertToBool(value string) bool {
	v, err := strconv.ParseBool(value)
	if err != nil {
		return false
	}

	return v
}
