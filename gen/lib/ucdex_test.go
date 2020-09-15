package lib

import (
	"bytes"
	"testing"
)

func TestGenUCDEX(t *testing.T) {
	var buf bytes.Buffer
	err := GenUCDEX(&buf)
	//If it ends normally, the test is OK
	if err != nil {
		t.Errorf("error: %s", err.Error())
	}
}

func TestGenerate(t *testing.T) {
	var buf bytes.Buffer
	err := Generate(&buf, "ucdex_test.go")
	//If it ends normally, the test is OK
	if err != nil {
		t.Errorf("error: %s", err.Error())
	}
}
