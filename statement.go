package abac

import (
	"errors"
	"strings"
)

const (
	Allow = "allow"
	Deny  = "deny"
)

var (
	ErrDeniedByStatement = errors.New("action denied as statement declared")
	ErrStatementNotMatch = errors.New("statement does not match the action")
)

type Statement struct {
	Resources
	Effect     string
	Actions    []string
	Conditions Conditions
}

func (s Statement) Validate(action string, resource Resource, attributes Attributes) error {
	if !s.Resources.Match(resource) {
		return ErrStatementNotMatch
	}

	if !MatchStrPatterns(s.Actions, action, StrEqual) {
		return ErrStatementNotMatch
	}

	if meet, err := s.Conditions.Judge(attributes); err != nil || !meet {
		return ErrStatementNotMatch
	}

	if strings.EqualFold(s.Effect, Deny) {
		return ErrDeniedByStatement
	}

	return nil
}

type Resource struct {
	ResourceType string
	ResourceName string
}

type Resources struct {
	ResourceType  string
	ResourceNames []string
	CompareFunc   StrComparable
}

func (r Resources) Match(resource Resource) bool {
	if resource.ResourceType == r.ResourceType {
		if MatchStrPatterns(r.ResourceNames, resource.ResourceName, r.CompareFunc) {
			return true
		}
	}
	return false
}
