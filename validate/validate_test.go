package validate

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type TestObj struct {
	TRange string  `verify-range:"1-20"`
	TMin   int64   `verify-range:"1-"`
	TMax   float64 `verify-range:"-20"`
	TEmpty string  `verify-nonempty:"-"`
}

func Test_verify(t *testing.T) {
	testObj := TestObj{
		"90", -1, 22, "",
	}
	_, errors := Validate(testObj)

	assert.True(t, len(errors) == 0, "testObj validate result:", errors)
}
