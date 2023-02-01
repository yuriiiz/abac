# abac
Simple implementation for Attribute-Based Access Control (ABAC).

## Install
```bash
go get github.com/yuriiiz/abac
```

## Example
```go
policy := Policy{
    Name:        "Name",
    Description: "Description",
    Statements: []Statement{
        {
            Effect:  Allow,
            Actions: []string{"read", "write", "audit"},
            Resources: Resources{
                ResourceType:  "resource_type_1",
                ResourceNames: []string{"*"},
                CompareFunc:   StrGlob,
            },
        },
        {
            Effect:  Allow,
            Actions: []string{"read", "write", "audit"},
            Resources: Resources{
                ResourceType:  "resource_type_2",
                ResourceNames: []string{"prefix_.*"},
                CompareFunc:   StrRegex,
            },
        },
        {
            Effect:  Deny,
            Actions: []string{"write"},
            Resources: Resources{
                ResourceType:  "resource_type_1",
                ResourceNames: []string{"resource_name_1", "resource_name_2", "resource_name_3"},
                CompareFunc:   StrEqual,
            },
            Conditions: Conditions{
                Items: []Condition{
                    {
                        Operator:      "gt",
                        AttributeType: AttributeTypeEnv,
                        AttributeKey:  "current_timestamp",
                        Value:         1600000000,
                    },
                    {
                        Operator:      "lt",
                        AttributeType: AttributeTypeEnv,
                        AttributeKey:  "current_timestamp",
                        Value:         1700000000,
                    },
                    {
                        Operator:      "eq",
                        AttributeType: AttributeTypeResource,
                        AttributeKey:  "environment",
                        Value:         "production",
                    },
                },
                Formula: "$0 and $1 and $2", // TODO: to be supported
            },
        },
    },
}

var resource Resource
var attributes Attributes

resource = Resource{
    ResourceType: "resource_type_1",
    ResourceName: "resource_name",
}
fmt.Println(policy.Validate("read", resource, attributes)) // -> nil

resource = Resource{
    ResourceType: "resource_type_2",
    ResourceName: "prefix_resource_name",
}
fmt.Println(policy.Validate("read", resource, attributes)) // -> nil

resource = Resource{
    ResourceType: "resource_type_2",
    ResourceName: "resource_name",
}
fmt.Println(policy.Validate("read", resource, attributes)) // -> ErrStatementNotMatch

resource = Resource{
    ResourceType: "resource_type_1",
    ResourceName: "resource_name_1",
}
attributes = []Attribute{
    {
        Type:      AttributeTypeResource,
        Key:       "environment",
        ValueType: "string",
        Value:     "production",
    },
    {
        Type:      AttributeTypeEnv,
        Key:       "current_timestamp",
        ValueType: "number",
        Value:     1650000000,
    },
}
fmt.Println(policy.Validate("write", resource, attributes)) // -> ErrDeniedByStatement
```
