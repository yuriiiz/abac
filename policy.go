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

func (p Policy) Validate(action string, resource Resource, attributes Attributes) error {
	validators := make([]permissionValidator, 0, len(p.Statements))

	for _, statement := range p.Statements {
		validators = append(validators, statement.Validate)
	}

	return mergePermissionValidators(validators...)(action, resource, attributes)
}

type Policies []Policy

func (ps Policies) Validate(action string, resource Resource, attributes Attributes) error {
	validators := make([]permissionValidator, 0, len(ps))

	for _, policy := range ps {
		validators = append(validators, policy.Validate)
	}

	return mergePermissionValidators(validators...)(action, resource, attributes)
}

type permissionValidator func(action string, resource Resource, attributes Attributes) error

func mergePermissionValidators(validators ...permissionValidator) permissionValidator {
	return func(action string, resource Resource, attributes Attributes) error {
		var notAllowedErr error
		var allowed bool

		for _, validator := range validators {
			if err := validator(action, resource, attributes); err != nil {
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
}
