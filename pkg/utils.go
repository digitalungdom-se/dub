package pkg

import (
	"github.com/bwmarrin/snowflake"
)

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func IsSnowflake(input string) bool {
	_, err := snowflake.ParseString(input)

	if err != nil {
		return false
	}

	return true
}
