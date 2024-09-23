package mappath

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_EmptyOrEqualToDefault(t *testing.T) {
	tests := []struct {
		name           string
		value          PathValue
		defaultValue   PathValue
		expectedResult bool
	}{
		{
			name:           "empty value",
			value:          MustUnmarshalScalar("null"),
			defaultValue:   MustUnmarshalScalar("null"),
			expectedResult: true,
		},
		{
			name:           "boolean true equal to default",
			value:          MustUnmarshalScalar("true"),
			defaultValue:   MustUnmarshalScalar("true"),
			expectedResult: true,
		},
		{
			name:           "boolean true not equal to default",
			value:          MustUnmarshalScalar("true"),
			defaultValue:   MustUnmarshalScalar("false"),
			expectedResult: false,
		},
		{
			name:           "boolean false equal to default",
			value:          MustUnmarshalScalar("false"),
			defaultValue:   MustUnmarshalScalar("false"),
			expectedResult: true,
		},
		{
			name:           "boolean false not equal to default",
			value:          MustUnmarshalScalar("false"),
			defaultValue:   MustUnmarshalScalar("true"),
			expectedResult: false,
		},
		{
			name:           "string empty",
			value:          MustUnmarshalScalar(`""`),
			defaultValue:   MustUnmarshalScalar(`"Doesn't matter"`),
			expectedResult: true,
		},
		{
			name:           "string non-empty and equal to default",
			value:          MustUnmarshalScalar(`"something"`),
			defaultValue:   MustUnmarshalScalar(`"something"`),
			expectedResult: true,
		},
		{
			name:           "string non-empty and not equal to default",
			value:          MustUnmarshalScalar(`"something"`),
			defaultValue:   MustUnmarshalScalar(`"something else"`),
			expectedResult: false,
		},
		{
			name:           "number equal to default",
			value:          MustUnmarshalScalar("2"),
			defaultValue:   MustUnmarshalScalar("2"),
			expectedResult: true,
		},
		{
			name:           "number not equal to default",
			value:          MustUnmarshalScalar("2"),
			defaultValue:   MustUnmarshalScalar("3"),
			expectedResult: false,
		},
		{
			name:           "object empty",
			value:          MustUnmarshalMap(`{}`),
			defaultValue:   MustUnmarshalMap(`{"default_value": "meaningless" }`),
			expectedResult: true,
		},
		{
			name:           "object non-empty and equal to default",
			value:          MustUnmarshalMap(`{"value" : "something"}`),
			defaultValue:   MustUnmarshalMap(`{"value" : "something"}`),
			expectedResult: true,
		},
		{
			name:           "object non-empty and not equal to default, different keys",
			value:          MustUnmarshalMap(`{ "value" : "something" }`),
			defaultValue:   MustUnmarshalMap(`{ "value_1" : "something" }`),
			expectedResult: false,
		},
		{
			name:           "object non-empty and not equal to default, different values",
			value:          MustUnmarshalMap(`{ "value" : "something" }`),
			defaultValue:   MustUnmarshalMap(`{ "value" : "something else" }`),
			expectedResult: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := EmptyOrEqualToDefault(test.value, test.defaultValue)
			require.Equal(t, test.expectedResult, result)
		})
	}
}
