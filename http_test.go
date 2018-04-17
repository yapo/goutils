package goutils

import (
	"gopkg.in/stretchr/testify.v1/assert"
	"net/http/httptest"
	"strings"
	"testing"
)

type TestStruct struct {
	A int    `json:"le_a"`
	B string `json:"le_b"`
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
