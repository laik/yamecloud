package utils

import (
	"testing"
)

func Test_Sha1(t *testing.T) {
	data := "123456"
	expectValue := "7c4a8d09ca3762af61e59520943dc26494f8941b"
	if Sha1(data) != expectValue {
		t.Fatal("result not expect")
	}

}
