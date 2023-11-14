package urlparams

import (
	"testing"
)

func TestMarshalling(t *testing.T) {
	type param struct {
		name     string
		given    any
		expected string
		err      error
	}

	params := []param{
		{
			name: "test with single field",
			given: struct {
				TestField string `urlparams:"test"`
			}{
				TestField: "test",
			},
			expected: "test=test",
		},
		{
			name: "test with multiple fields",
			given: struct {
				SomeField     string `urlparams:"some"`
				OptionalField string `urlparams:"optional,omitempty"`
			}{
				SomeField:     "somevalue",
				OptionalField: "optional",
			},
			expected: "optional=optional&some=somevalue",
		},
		{
			name: "test with multiple fields and empty optional",
			given: struct {
				SomeField     string `urlparams:"some2"`
				OptionalField string `urlparams:"optional,omitempty"`
			}{
				SomeField: "singlevalue",
			},
			expected: "some2=singlevalue",
		},
		{
			name: "test with multiple fields and ignored field",
			given: struct {
				F1           string `urlparams:"f1"`
				IgnoredField string
				F2           string `urlparams:"f2"`
			}{
				F1: "singlevalue1",
				F2: "singlevalue2",
			},
			expected: "f1=singlevalue1&f2=singlevalue2",
		},
		{
			name: "encoding values of fields",
			given: struct {
				TextField string `urlparams:"text"`
			}{
				TextField: "value that need to be encoded &*^#@@$8",
			},
			expected: "text=value+that+need+to+be+encoded+%26%2A%5E%23%40%40%248",
		},
		{
			name: "test different types of values",
			given: struct {
				StringField  string     `urlparams:"string"`
				IntField     int        `urlparams:"int"`
				UintField    uint       `urlparams:"uint"`
				BoolField    bool       `urlparams:"bool"`
				FloatField   float64    `urlparams:"float"`
				ComplexField complex128 `urlparams:"complex"`
			}{
				StringField:  "value",
				IntField:     -146541,
				UintField:    231,
				BoolField:    true,
				FloatField:   1.984e2,
				ComplexField: 10 + 11i,
			},
			expected: "bool=true&complex=%2810.000000%2B11.000000i%29&float=198.400000&int=-146541&string=value&uint=231",
		},
		{
			name: "test different passed types",
			given: &struct {
				StringField string `urlparams:"string"`
			}{
				StringField: "value",
			},
			expected: "string=value",
		},
	}

	for _, param := range params {
		t.Run(param.name, func(t *testing.T) {
			got, err := Marshal(param.given)
			if err != param.err {
				t.Errorf("got %s, expected %s", err, param.err)
			}
			if got != param.expected {
				t.Errorf("got %s, expected %s", got, param.expected)
			}
		})
	}

}
