package v3

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	v3 "go.signoz.io/signoz/pkg/query-service/model/v3"
)

var testGetJSONFilterKeyData = []struct {
	Name          string
	Key           v3.AttributeKey
	IsArray       bool
	ClickhouseKey string
	Error         bool
}{
	{
		Name: "Incorrect Key",
		Key: v3.AttributeKey{
			Key:      "requestor_list[*]",
			DataType: "array(string)",
			IsJSON:   true,
		},
		IsArray: true,
		Error:   true,
	},
	{
		Name: "Using anything other than body",
		Key: v3.AttributeKey{
			Key:      "trace_id.requestor_list[*]",
			DataType: "array(string)",
			IsJSON:   true,
		},
		IsArray: true,
		Error:   true,
	},
	{
		Name: "Array String",
		Key: v3.AttributeKey{
			Key:      "body.requestor_list[*]",
			DataType: "array(string)",
			IsJSON:   true,
		},
		IsArray:       true,
		ClickhouseKey: "JSONExtract(JSON_QUERY(body, '$.requestor_list[*]'), '" + ARRAY_STRING + "')",
	},
	{
		Name: "Array String Nested",
		Key: v3.AttributeKey{
			Key:      "body.nested[*].key[*]",
			DataType: "array(string)",
			IsJSON:   true,
		},
		IsArray:       true,
		ClickhouseKey: "JSONExtract(JSON_QUERY(body, '$.nested[*].key[*]'), '" + ARRAY_STRING + "')",
	},
	{
		Name: "Array Int",
		Key: v3.AttributeKey{
			Key:      "body.int_numbers[*]",
			DataType: "array(int64)",
			IsJSON:   true,
		},
		IsArray:       true,
		ClickhouseKey: "JSONExtract(JSON_QUERY(body, '$.int_numbers[*]'), '" + ARRAY_INT64 + "')",
	},
	{
		Name: "Array Float",
		Key: v3.AttributeKey{
			Key:      "body.nested_num[*].float_nums[*]",
			DataType: "array(float64)",
			IsJSON:   true,
		},
		IsArray:       true,
		ClickhouseKey: "JSONExtract(JSON_QUERY(body, '$.nested_num[*].float_nums[*]'), '" + ARRAY_FLOAT64 + "')",
	},
	{
		Name: "Array Bool",
		Key: v3.AttributeKey{
			Key:      "body.boolarray[*]",
			DataType: "array(bool)",
			IsJSON:   true,
		},
		IsArray:       true,
		ClickhouseKey: "JSONExtract(JSON_QUERY(body, '$.boolarray[*]'), '" + ARRAY_BOOL + "')",
	},
	{
		Name: "String",
		Key: v3.AttributeKey{
			Key:      "body.message",
			DataType: "string",
			IsJSON:   true,
		},
		IsArray:       false,
		ClickhouseKey: "JSON_VALUE(body, '$.message')",
	},
	{
		Name: "Int",
		Key: v3.AttributeKey{
			Key:      "body.status",
			DataType: "int64",
			IsJSON:   true,
		},
		IsArray:       false,
		ClickhouseKey: "JSONExtract(JSON_VALUE(body, '$.status'), '" + INT64 + "')",
	},
	{
		Name: "Float",
		Key: v3.AttributeKey{
			Key:      "body.fraction",
			DataType: "float64",
			IsJSON:   true,
		},
		IsArray:       false,
		ClickhouseKey: "JSONExtract(JSON_VALUE(body, '$.fraction'), '" + FLOAT64 + "')",
	},
	{
		Name: "Bool",
		Key: v3.AttributeKey{
			Key:      "body.boolkey",
			DataType: "bool",
			IsJSON:   true,
		},
		IsArray:       false,
		ClickhouseKey: "JSONExtract(JSON_VALUE(body, '$.boolkey'), '" + BOOL + "')",
	},
}

func TestGetJSONFilterKey(t *testing.T) {
	for _, tt := range testGetJSONFilterKeyData {
		Convey("testgetKey", t, func() {
			columnName, err := getJSONFilterKey(tt.Key, tt.IsArray)
			if tt.Error {
				So(err, ShouldNotBeNil)
			} else {
				So(err, ShouldBeNil)
				So(columnName, ShouldEqual, tt.ClickhouseKey)
			}
		})
	}
}

var testGetJSONFilterData = []struct {
	Name       string
	FilterItem v3.FilterItem
	Filter     string
	Error      bool
}{
	{
		Name: "Array membership string",
		FilterItem: v3.FilterItem{
			Key: v3.AttributeKey{
				Key:      "body.requestor_list[*]",
				DataType: "array(string)",
				IsJSON:   true,
			},
			Operator: "has",
			Value:    "index_service",
		},
		Filter: "has(JSONExtract(JSON_QUERY(body, '$.requestor_list[*]'), 'Array(String)'), 'index_service')",
	},
	{
		Name: "Array membership int64",
		FilterItem: v3.FilterItem{
			Key: v3.AttributeKey{
				Key:      "body.int_numbers[*]",
				DataType: "array(int64)",
				IsJSON:   true,
			},
			Operator: "has",
			Value:    2,
		},
		Filter: "has(JSONExtract(JSON_QUERY(body, '$.int_numbers[*]'), '" + ARRAY_INT64 + "'), 2)",
	},
	{
		Name: "Array membership float64",
		FilterItem: v3.FilterItem{
			Key: v3.AttributeKey{
				Key:      "body.nested_num[*].float_nums[*]",
				DataType: "array(float64)",
				IsJSON:   true,
			},
			Operator: "nhas",
			Value:    2.2,
		},
		Filter: "NOT has(JSONExtract(JSON_QUERY(body, '$.nested_num[*].float_nums[*]'), '" + ARRAY_FLOAT64 + "'), 2.200000)",
	},
	{
		Name: "Array membership bool",
		FilterItem: v3.FilterItem{
			Key: v3.AttributeKey{
				Key:      "body.bool[*]",
				DataType: "array(bool)",
				IsJSON:   true,
			},
			Operator: "has",
			Value:    true,
		},
		Filter: "has(JSONExtract(JSON_QUERY(body, '$.bool[*]'), '" + ARRAY_BOOL + "'), true)",
	},
	{
		Name: "eq operator",
		FilterItem: v3.FilterItem{
			Key: v3.AttributeKey{
				Key:      "body.message",
				DataType: "string",
				IsJSON:   true,
			},
			Operator: "=",
			Value:    "hello",
		},
		Filter: "JSON_VALUE(body, '$.message') = 'hello'",
	},
	{
		Name: "eq operator number",
		FilterItem: v3.FilterItem{
			Key: v3.AttributeKey{
				Key:      "body.status",
				DataType: "int64",
				IsJSON:   true,
			},
			Operator: "=",
			Value:    1,
		},
		Filter: "JSONExtract(JSON_VALUE(body, '$.status'), '" + INT64 + "') = 1",
	},
	{
		Name: "neq operator number",
		FilterItem: v3.FilterItem{
			Key: v3.AttributeKey{
				Key:      "body.status",
				DataType: "float64",
				IsJSON:   true,
			},
			Operator: "=",
			Value:    1.1,
		},
		Filter: "JSONExtract(JSON_VALUE(body, '$.status'), '" + FLOAT64 + "') = 1.100000",
	},
	{
		Name: "eq operator bool",
		FilterItem: v3.FilterItem{
			Key: v3.AttributeKey{
				Key:      "body.boolkey",
				DataType: "bool",
				IsJSON:   true,
			},
			Operator: "=",
			Value:    true,
		},
		Filter: "JSONExtract(JSON_VALUE(body, '$.boolkey'), '" + BOOL + "') = true",
	},
	{
		Name: "greater than operator",
		FilterItem: v3.FilterItem{
			Key: v3.AttributeKey{
				Key:      "body.status",
				DataType: "int64",
				IsJSON:   true,
			},
			Operator: ">",
			Value:    1,
		},
		Filter: "JSONExtract(JSON_VALUE(body, '$.status'), '" + INT64 + "') > 1",
	},
	{
		Name: "regex operator",
		FilterItem: v3.FilterItem{
			Key: v3.AttributeKey{
				Key:      "body.message",
				DataType: "string",
				IsJSON:   true,
			},
			Operator: "regex",
			Value:    "a*",
		},
		Filter: "match(JSON_VALUE(body, '$.message'), 'a*')",
	},
	{
		Name: "contains operator",
		FilterItem: v3.FilterItem{
			Key: v3.AttributeKey{
				Key:      "body.message",
				DataType: "string",
				IsJSON:   true,
			},
			Operator: "contains",
			Value:    "a",
		},
		Filter: "JSON_VALUE(body, '$.message') ILIKE '%a%'",
	},
}

func TestGetJSONFilter(t *testing.T) {
	for _, tt := range testGetJSONFilterData {
		Convey("testGetJSONFilter", t, func() {
			filter, err := GetJSONFilter(tt.FilterItem)
			if tt.Error {
				So(err, ShouldNotBeNil)
			} else {
				So(err, ShouldBeNil)
				So(filter, ShouldEqual, tt.Filter)
			}
		})
	}
}
