package lib

import (
	"bytes"
	"testing"
)

func TestHeavyGenUCD(t *testing.T) {
	var buf bytes.Buffer
	err := GenUCD(&buf)
	//If it ends normally, the test is OK
	if err != nil {
		t.Errorf("error: %s", err.Error())
	}
}
