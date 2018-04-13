package main

import "testing"

func TestUnzip(t *testing.T) {
	if _, err := Unzip("./res.zip", "./"); err != nil {
		t.Error(err)
	}
}
