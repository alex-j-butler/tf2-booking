package util

import "testing"
import "github.com/google/go-github/github"

func TestGetReleaseAsset(t *testing.T) {
	name1 := "Example asset"
	name2 := "tf2booking-amd64"
	name3 := "tf2booking-i386"

	actualResult := GetReleaseAsset([]github.ReleaseAsset{
		{
			Name: &name1,
		},
		{
			Name: &name2,
		},
		{
			Name: &name3,
		},
	})

	expectedResult := true

	if actualResult != expectedResult {
		t.Fatalf("Expected %t but got %t", expectedResult, actualResult)
	}
}
