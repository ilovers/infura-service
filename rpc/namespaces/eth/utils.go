package eth

import (
	"bytes"

	"github.com/ethereum/go-ethereum/common"
)

// The Topic list restricts matches to particular event topics. Each event has a list
// of topics. Topics matches a prefix of that list. An empty element slice matches any
// topic. Non-empty elements represent an alternative that matches any of the
// contained topics.
//
// Examples:
// {} or nil          matches any topic list
// {{A}}              matches topic A in first position
// {{}, {B}}          matches any topic in first position AND B in second position
// {{A}, {B}}         matches topic A in first position AND B in second position
// {{A, B}, {C, D}}   matches topic (A OR B) in first position AND (C OR D) in second position
func matchTopics(topics []common.Hash, matches [][]common.Hash) bool {
	matchCount := len(matches)
	// 处理传入的参数，topic末位写入null的情况，比如传入：{{A}, nil, nil, nil}，这种情况只要第一个A满足条件，后面的nil忽略即可
	for i := len(matches) - 1; i >= 0; i-- {
		if len(matches[i]) > 0 {
			break
		}
		matchCount--
	}
	// 要求的topic数量不匹配
	if matchCount > len(topics) {
		return false
	}
	// 验证topic
	for i := 0; i < matchCount; i++ {
		if len(matches[i]) == 0 {
			continue
		}
		isMatch := false
		for _, match := range matches[i] {
			if bytes.Equal(topics[i].Bytes(), match.Bytes()) {
				isMatch = true
				break
			}
		}
		if !isMatch {
			return false
		}
	}
	return true
}
