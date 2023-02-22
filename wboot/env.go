package wboot

import (
	"os"
	"strconv"
	"strings"
)

func EnvStr(key string) string {
	return strings.TrimSpace(os.Getenv(key))
}

func EnvStrOr(key string, d string) string {
	if v := EnvStr(key); v != "" {
		return v
	}
	return d
}

func EnvBool(key string) bool {
	v, _ := strconv.ParseBool(EnvStr(key))
	return v
}

func EnvBoolOr(key string, d bool) bool {
	if v, err := strconv.ParseBool(EnvStr(key)); err != nil {
		return d
	} else {
		return v
	}
}
