package main

import (
	"testing"
)

func TestData(t *testing.T) {
	t.Parallel()

	data1 := NewData()
	data2 := NewData()

	if data1 == data2 {
		t.Error("incorrect Data", data2)
	}
}

func TestSha(t *testing.T) {
	t.Parallel()

	data := Data{}
	sha := data.Sha()

	s := sha.String()
	if s != "17b0761f87b081d5cf10757ccc89f12be355c70e2e29df288b65b30710dcbcd1" {
		t.Error("incorrect Sha", s)
	}
}

func TestBuffer(t *testing.T) {
	t.Parallel()

	buffer := NewBuffer()

	s1 := buffer.Sha.String()
	s2 := buffer.Data.Sha().String()
	if s1 != s2 {
		t.Error("incorrect Buffer", s1, s2)
	}
}
