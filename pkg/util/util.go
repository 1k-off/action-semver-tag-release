package util

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func GetEnvBool(key string) (bool, error) {
	s := os.Getenv(key)
	if s == "" {
		return false, nil
	}
	v, err := strconv.ParseBool(s)
	if err != nil {
		return false, err
	}
	return v, nil
}

func SetGithubOutput(key, value string) {
	fmt.Printf(`::set-output name=%s::%s`, key, value)
}

// GetEnvArray returns an array of strings from the environment variable
func GetEnvArray(key string) []string {
	s := os.Getenv(key)
	if s == "" {
		return []string{}
	}
	return deleteEmpty(strings.Split(s, "\n"))
}

func deleteEmpty(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}

func GetStringAsArray(s string) []string {
	if s == "" {
		return []string{}
	}
	return deleteEmpty(strings.Split(s, "\n"))
}
