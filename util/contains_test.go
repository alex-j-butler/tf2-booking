package util

import "testing"

func TestContains(t *testing.T) {
	actualResult := Contains([]string{
		"example",
		"test list",
		"test",
	}, "test list")

	expectedResult := true

	if actualResult != expectedResult {
		t.Fatalf("Expected %t but got %t", expectedResult, actualResult)
	}
}

func TestNotContains(t *testing.T) {
	actualResult := Contains([]string{
		"example",
		"test list",
		"test",
	}, "not in the slice")

	expectedResult := false

	if actualResult != expectedResult {
		t.Fatalf("Expected %t but got %t", expectedResult, actualResult)
	}
}

func TestContainsNumbers(t *testing.T) {
	actualResult := Contains([]string{
		"188795674523211279",
		"627789334543599593",
		"274532788613968856",
	}, "188795674523211279")

	expectedResult := true

	if actualResult != expectedResult {
		t.Fatalf("Expected %t but got %t", expectedResult, actualResult)
	}
}

func TestNotContainsNumbers(t *testing.T) {
	actualResult := Contains([]string{
		"188795674523211279",
		"627789334543599593",
		"274532788613968856",
	}, "794665738822333211")

	expectedResult := false

	if actualResult != expectedResult {
		t.Fatalf("Expected %t but got %t", expectedResult, actualResult)
	}
}
