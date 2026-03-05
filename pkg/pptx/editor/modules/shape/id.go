package shape

import (
	"regexp"
	"strconv"
)

func MaxObjectID(content []byte, pattern *regexp.Regexp, submatchSize int) int {
	matches := pattern.FindAllSubmatch(content, -1)
	maxID := 0
	for _, match := range matches {
		if len(match) < submatchSize {
			continue
		}
		id, err := strconv.Atoi(string(match[1]))
		if err != nil {
			continue
		}
		if id > maxID {
			maxID = id
		}
	}
	return maxID
}
