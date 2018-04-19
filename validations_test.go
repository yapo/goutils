package goutils

import (
	"gopkg.in/stretchr/testify.v1/assert"
	"testing"
)

func TestIntegerValidators(t *testing.T) {
	fns := []ValidatorFunc{
		IntegerValidator,
		StringValidator,
		EmptyStringValidator,
		PlatesValidator,
		EmailValidator,
		ValidateDate,
	}

	data := []interface{}{
		42,
		3.14159,
		"Don't Panic",
		"",
		"XY1234",
		"yo@lo.cl",
		"yo@not",
		"2018-12-13",
		"2018/12/13",
	}

	expected := [][]bool{
		{true, false, false, false, false, false},
		{true, false, false, false, false, false},
		{false, true, true, false, false, false},
		{false, true, false, false, false, false},
		{false, true, true, true, false, false},
		{false, true, true, false, true, false},
		{false, true, true, false, false, false},
		{false, true, true, false, false, true},
		{false, true, true, false, false, false},
	}
	for i, val := range data {
		for j, fn := range fns {
			r, _ := fn(val)
			assert.Equal(t, expected[i][j], r)
		}
	}

}
