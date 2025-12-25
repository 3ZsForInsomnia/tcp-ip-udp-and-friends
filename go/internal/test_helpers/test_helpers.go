package testhelpers

import "testing"

func FailTestIfErrorIsPresent(t *testing.T, err error) {
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}
