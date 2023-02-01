package abac

import (
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

const (
	valueTypeString = "string"
	valueTypeNumber = "number"

	OperatorEQ         = "eq"
	OperatorLT         = "lt"
	OperatorLTE        = "lte"
	OperatorGT         = "gt"
	OperatorGTE        = "gte"
	OperatorNEQ        = "neq"
	OperatorGlobMatch  = "glob_match"
	OperatorRegexMatch = "regex_match"
)

func Compare(valueType, operator string, originValue, compareValue interface{}) (bool, error) {
	comparer := newComparer(valueType)
	return comparer.Compare(operator, originValue, compareValue)
}

func newComparer(valueType string) valueComparable {
	switch strings.ToLower(valueType) {
	case valueTypeString:
		return &strComparer{}
	case valueTypeNumber:
		return &numberComparer{}
	default:
		return &strComparer{}
	}
}

type valueComparable interface {
	Compare(operator string, originValue, compareValue interface{}) (bool, error)
}

type strComparer struct{}

func (c strComparer) Compare(operator string, originValue, compareValue interface{}) (bool, error) {
	var originValueStr, compareValueStr string
	var ok bool

	if originValueStr, ok = originValue.(string); !ok {
		return false, errors.New("value type mismatch")
	}

	if compareValueStr, ok = compareValue.(string); !ok {
		return false, errors.New("value type mismatch")
	}

	switch strings.ToLower(operator) {
	case OperatorEQ:
		return originValueStr == compareValueStr, nil
	case OperatorLT:
		return originValueStr < compareValueStr, nil
	case OperatorLTE:
		return originValueStr <= compareValueStr, nil
	case OperatorGT:
		return originValueStr > compareValueStr, nil
	case OperatorGTE:
		return originValueStr >= compareValueStr, nil
	case OperatorNEQ:
		return originValueStr != compareValueStr, nil
	case OperatorGlobMatch:
		return StrGlob(compareValueStr, originValueStr), nil
	case OperatorRegexMatch:
		return StrRegex(compareValueStr, originValueStr), nil
	}

	return false, errors.Errorf("unsupported operator %s", operator)
}

type numberComparer struct{}

func (c numberComparer) Compare(operator string, originValue, compareValue interface{}) (bool, error) {
	var originValueNum, compareValueNum int
	var ok bool

	if originValueNum, ok = originValue.(int); !ok {
		return false, errors.New("value type mismatch")
	}

	if compareValueNum, ok = compareValue.(int); !ok {
		return false, errors.New("value type mismatch")
	}

	switch strings.ToLower(operator) {
	case OperatorEQ:
		return originValueNum == compareValueNum, nil
	case OperatorLT:
		return originValueNum < compareValueNum, nil
	case OperatorLTE:
		return originValueNum <= compareValueNum, nil
	case OperatorGT:
		return originValueNum > compareValueNum, nil
	case OperatorGTE:
		return originValueNum >= compareValueNum, nil
	case OperatorNEQ:
		return originValueNum != compareValueNum, nil
	}

	return false, errors.Errorf("unsupported operator %s", operator)
}

type StrComparable func(pattern, subj string) bool

func StrRegex(pattern, subj string) bool {
	match, _ := regexp.MatchString(pattern, subj)
	return match
}

func StrEqual(pattern, subj string) bool {
	return pattern == subj
}

func StrGlob(pattern, subj string) bool {
	const GLOB = "*"

	// Empty pattern can only match empty subject
	if pattern == "" {
		return subj == pattern
	}

	// If the pattern _is_ a glob, it matches everything
	if pattern == GLOB {
		return true
	}

	parts := strings.Split(pattern, GLOB)

	if len(parts) == 1 {
		// No globs in pattern, so test for equality
		return subj == pattern
	}

	leadingGlob := strings.HasPrefix(pattern, GLOB)
	trailingGlob := strings.HasSuffix(pattern, GLOB)
	end := len(parts) - 1

	// Go over the leading parts and ensure they match.
	for i := 0; i < end; i++ {
		idx := strings.Index(subj, parts[i])

		switch i {
		case 0:
			// Check the first section. Requires special handling.
			if !leadingGlob && idx != 0 {
				return false
			}
		default:
			// Check that the middle parts match.
			if idx < 0 {
				return false
			}
		}

		// Trim evaluated text from subj as we loop over the pattern.
		subj = subj[idx+len(parts[i]):]
	}

	// Reached the last section. Requires special handling.
	return trailingGlob || strings.HasSuffix(subj, parts[end])
}

func MatchStrPatterns(patterns []string, subj string, compareFunc StrComparable) bool {
	if compareFunc == nil {
		compareFunc = StrGlob
	}

	for _, pattern := range patterns {
		if compareFunc(pattern, subj) {
			return true
		}
	}
	return false
}
