/*
 * 2013-07-07
 * will@summercat.com
 */

package config

import (
	"testing"
)

// TestPopulateConfig is to test the conversion of types.
func TestPopulateConfig(t *testing.T) {
	type MyType struct {
		Str string
		Abc int64
	}
	var rawValues = map[string]string{
		"Str": "Hi there",
		"Abc": "123",
	}

	var myT MyType
	err := populateConfig(&myT, rawValues)
	if err != nil {
		t.Errorf("failed to populate: %s", err.Error())
	}
	if myT.Str != "Hi there" {
		t.Errorf("Failed to parse string")
	}
	if myT.Abc != 123 {
		t.Errorf("Failed to parse int")
	}
}
