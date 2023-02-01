package abac

import (
	"fmt"

	"github.com/pkg/errors"
)

const (
	AttributeTypeSubject  = "AttributeTypeSubject"
	AttributeTypeResource = "AttributeTypeResource"
	AttributeTypeEnv      = "AttributeTypeEnv"
)

type Attribute struct {
	Type      string // subject/resource/action/env
	Key       string
	ValueType string
	Value     interface{}
}

type Attributes []Attribute

type Conditions struct {
	Items   []Condition
	Formula string
}

type Condition struct {
	Operator      string // equal/gt/lt
	AttributeType string
	AttributeKey  string
	Value         interface{}
}

func (cs Conditions) Judge(attributes Attributes) (bool, error) {
	attrPathGen := func(attrType, attrKey string) string {
		return fmt.Sprintf("%s:%s", attrType, attrKey)
	}

	attributesMap := make(map[string]Attribute)
	for _, attribute := range attributes {
		attributesMap[attrPathGen(attribute.Type, attribute.Key)] = attribute
	}

	meet := true

	for _, condition := range cs.Items {
		attrPath := attrPathGen(condition.AttributeType, condition.AttributeKey)
		attribute, exist := attributesMap[attrPath]
		if !exist {
			return false, errors.Errorf("attribute %s not given", attrPath)
		}

		conditionMeet, err := Compare(attribute.ValueType, condition.Operator, attribute.Value, condition.Value)
		if err != nil {
			return false, err
		}

		meet = meet && conditionMeet
	}

	// TODO: support condition formula
	return meet, nil
}
