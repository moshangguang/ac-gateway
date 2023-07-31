package strutils

import (
	"hash/fnv"
	"strings"
)

func IsEmpty(str string) bool {
	return len(strings.TrimSpace(str)) == 0
}
func IsNotEmpty(str string) bool {
	return !IsEmpty(str)
}
func Hash(s string) uint32 {
	h := fnv.New32a()
	_, _ = h.Write([]byte(s))
	return h.Sum32()
}
