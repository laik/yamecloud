package utils

import "testing"

func Test_ContainStringItem(t *testing.T) {
	list := []string{"A", "B"}
	if ContainStringItem(list, "C") {
		t.Fatal("expected not equal")
	}
	if !ContainStringItem(list, "A") || !ContainStringItem(list, "B") {
		t.Fatal("expected not equal")
	}
}
