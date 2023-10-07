package util

import (
	"errors"
	"net/http"
	"reflect"
	"testing"
)

func TestOutputBasic_GetResult(t *testing.T) {
	tests := []struct {
		name     string
		output   *OutputBasic
		expected map[string]interface{}
	}{
		{
			name: "Message is an error",
			output: &OutputBasic{
				Code:    http.StatusInternalServerError,
				Result:  "Internal Server Error",
				Message: errors.New("something went wrong"),
			},
			expected: map[string]interface{}{
				"code":    http.StatusInternalServerError,
				"result":  "Internal Server Error",
				"message": "something went wrong",
			},
		},
		{
			name: "Message is not an error",
			output: &OutputBasic{
				Code:    http.StatusOK,
				Result:  "OK",
				Message: "Success",
			},
			expected: map[string]interface{}{
				"code":    http.StatusOK,
				"result":  "OK",
				"message": "Success",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.output.GetResult()
			if !reflect.DeepEqual(result, test.expected) {
				t.Errorf("Got %v, expected %v", result, test.expected)
			}
		})
	}
}

func TestOutputBasic_GetCode(t *testing.T) {
	output := &OutputBasic{
		Code:    http.StatusNotFound,
		Result:  "Not Found",
		Message: errors.New("resource not found"),
	}

	code := output.GetCode()
	if code != http.StatusNotFound {
		t.Errorf("Got code %d, expected %d", code, http.StatusNotFound)
	}
}

func TestOutputBasic_GetError(t *testing.T) {
	tests := []struct {
		name     string
		output   *OutputBasic
		expected error
	}{
		{
			name: "Message is an error",
			output: &OutputBasic{
				Code:    http.StatusInternalServerError,
				Result:  "Internal Server Error",
				Message: errors.New("something went wrong"),
			},
			expected: errors.New("something went wrong"),
		},
		{
			name: "Message is not an error",
			output: &OutputBasic{
				Code:    http.StatusOK,
				Result:  "OK",
				Message: "Success",
			},
			expected: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.output.GetError()
			if !reflect.DeepEqual(err, test.expected) {
				t.Errorf("Got error %v, expected %v", err, test.expected)
			}
		})
	}
}

func TestOutputBasicObject_GetResult(t *testing.T) {
	outDataSlice := interface{}([]string{"item1", "item2"})
	tests := []struct {
		name     string
		output   *OutputBasicObject
		expected map[string]interface{}
	}{
		{
			name: "Row is a slice",
			output: &OutputBasicObject{
				Code:    http.StatusOK,
				Result:  "OK",
				OutData: &outDataSlice,
			},
			expected: map[string]interface{}{
				"code":   http.StatusOK,
				"result": "OK",
				"list":   &outDataSlice,
			},
		},
		{
			name: "Row is not a slice",
			output: &OutputBasicObject{
				Code:    http.StatusOK,
				Result:  "OK",
				OutData: ToPrt("singleItem"),
			},
			expected: map[string]interface{}{
				"code":   http.StatusOK,
				"result": "OK",
				"row":    ToPrt("singleItem"),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.output.GetResult()
			if !reflect.DeepEqual(result, test.expected) {
				t.Errorf("Got %v, expected %v", result, test.expected)
			}
		})
	}
}

func TestOutputBasicObject_GetCode(t *testing.T) {
	output := &OutputBasicObject{
		Code:    http.StatusNotFound,
		Result:  "Not Found",
		OutData: ToPrt(interface{}([]string{"item1", "item2"})),
	}

	code := output.GetCode()
	if code != http.StatusNotFound {
		t.Errorf("Got code %d, expected %d", code, http.StatusNotFound)
	}
}

func TestOutputBasicListPaging_GetResult(t *testing.T) {
	tests := []struct {
		name     string
		output   *OutputBasicListPaging
		expected map[string]interface{}
	}{
		{
			name: "List is a slice",
			output: &OutputBasicListPaging{
				Code:       http.StatusOK,
				Result:     "OK",
				List:       ToPrt(interface{}([]string{"item1", "item2"})),
				CountTotal: 2,
				Page:       1,
				Limit:      10,
			},
			expected: map[string]interface{}{
				"code":        http.StatusOK,
				"result":      "OK",
				"list":        ToPrt(interface{}([]string{"item1", "item2"})),
				"count_total": int64(2),
				"page":        int64(1),
				"limit":       int64(10),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.output.GetResult()
			if !reflect.DeepEqual(result, test.expected) {
				t.Errorf("Got %v, expected %v", result, test.expected)
			}
		})
	}
}

func TestOutputBasicListPaging_GetCode(t *testing.T) {
	output := &OutputBasicListPaging{
		Code:       http.StatusNotFound,
		Result:     "Not Found",
		List:       ToPrt(interface{}([]string{"item1", "item2"})),
		CountTotal: 2,
		Page:       1,
		Limit:      10,
	}

	code := output.GetCode()
	if code != http.StatusNotFound {
		t.Errorf("Got code %d, expected %d", code, http.StatusNotFound)
	}
}
