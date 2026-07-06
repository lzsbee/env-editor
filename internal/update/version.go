package update

import (
	"strconv"
	"strings"
)

func normalizeVersion(v string) string {
	v = strings.TrimSpace(v)
	v = strings.TrimPrefix(v, "v")
	v = strings.TrimPrefix(v, "V")
	return v
}

func parseParts(v string) [3]int {
	parts := [3]int{}
	segments := strings.Split(normalizeVersion(v), ".")
	for i := 0; i < len(segments) && i < 3; i++ {
		n, err := strconv.Atoi(strings.TrimSpace(segments[i]))
		if err != nil {
			n = 0
		}
		parts[i] = n
	}
	return parts
}

func isNewer(latest, current string) bool {
	latest = normalizeVersion(latest)
	current = normalizeVersion(current)
	if latest == "" || current == "" || latest == current {
		return false
	}

	l := parseParts(latest)
	c := parseParts(current)
	for i := 0; i < 3; i++ {
		if l[i] > c[i] {
			return true
		}
		if l[i] < c[i] {
			return false
		}
	}
	return false
}
