package util

import (
	"os"
	"strconv"
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
	output := os.Getenv("GITHUB_OUTPUT")
	newString := key + "=" + value
	output += newString

	_ = os.Setenv("GITHUB_OUTPUT", output)
}
