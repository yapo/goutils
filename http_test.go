package goutils

import (
	"gopkg.in/stretchr/testify.v1/assert"
	"math"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type TestStruct struct {
	A float64 `json:"le_a"`
	B string  `json:"le_b"`
}

func TestWriteJSONResponseString(t *testing.T) {
	response := Response{
		Code: 42,
		Body: "Don't panic",
	}
	expected := `Don't panic`
	w := httptest.NewRecorder()
	WriteJSONResponse(w, &response)
	assert.Equal(t, expected, w.Body.String())
}

func TestWriteJSONResponseStruct(t *testing.T) {
	response := Response{
		Code: 42,
		Body: TestStruct{
			A: 314159,
			B: "Pi day",
		},
	}
	expected := `{"le_a":314159,"le_b":"Pi day"}`
	w := httptest.NewRecorder()
	CreateJSON(&response)
	WriteJSONResponse(w, &response)
	assert.Equal(t, expected, w.Body.String())
}

func TestWriteJSONResponseError(t *testing.T) {
	response := Response{
		Code: 42,
		Body: TestStruct{
			A: math.Inf(1),
			B: "And beyond!",
		},
	}
	expected := Response{
		Code: http.StatusInternalServerError,
		Body: "",
	}
	w := httptest.NewRecorder()
	CreateJSON(&response)
	WriteJSONResponse(w, &response)
	assert.Equal(t, expected, response)
	//assert.Equal(t, expected, w.Body.String())
}

func TestParseJSONBodyOK(t *testing.T) {
	body := `{"le_a": 42, "le_b": "Don't panic"}`
	req := httptest.NewRequest("GET", "/someurl", strings.NewReader(body))
	input := TestStruct{}
	expected := TestStruct{
		A: 42,
		B: "Don't panic",
	}

	r := ParseJSONBody(req, &input)
	assert.Nil(t, r)
	assert.Equal(t, expected, input)
}

func TestParseJSONBodyBadRequest(t *testing.T) {
	body := `{"le_a": 42, "le_b": "Do panic"`
	req := httptest.NewRequest("GET", "/someurl", strings.NewReader(body))
	input := TestStruct{}
	expected := TestStruct{}

	r := ParseJSONBody(req, &input)
	assert.NotNil(t, r)
	assert.Equal(t, expected, input)
}

func TestParseJSONBodySomeOtherJson(t *testing.T) {
	body := `{"yapo": 42, ".cl": "yo decido"}`
	req := httptest.NewRequest("GET", "/someurl", strings.NewReader(body))
	input := TestStruct{}
	expected := TestStruct{}

	r := ParseJSONBody(req, &input)
	assert.Nil(t, r)
	assert.Equal(t, expected, input)
}
