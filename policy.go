package abac

import (
	"github.com/pkg/errors"
)

type Permission struct {
	Policy
	Identity
}

type Identity struct {
	SubjectType string
	Subject     string
}

type Policy struct {
	ID          uint64
	Name        string
	Description string
	Generated   bool
	Statements  []Statement
}

func (p Policy) Allow(action string, resource Resource, attributes Attributes) error {
	var notAllowedErr error
	var allowed bool

	for _, statement := range p.Statements {
		if err := statement.Allow(action, resource, attributes); err != nil {
			// deny immediately
			if errors.Is(err, ErrDeniedByStatement) {
				return err
			}

			// record error
			notAllowedErr = err
		} else {
			// mark as allowed
			allowed = true
		}
	}

	if allowed {
		return nil
	}

	return notAllowedErr
}

type Policies []Policy

func (ps Policies) Allow(action string, resource Resource, attributes Attributes) error {
	var notAllowedErr error
	var allowed bool

	for _, policy := range ps {
		if err := policy.Allow(action, resource, attributes); err != nil {
			// deny immediately
			if errors.Is(err, ErrDeniedByStatement) {
				return ErrDeniedByStatement
			}

			// record error
			notAllowedErr = err
		} else {
			// mark as allowed
			allowed = true
		}
	}

	if allowed {
		return nil
	}

	return notAllowedErr
}
