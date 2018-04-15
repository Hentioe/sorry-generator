package main

import "testing"

func TestUnzip(t *testing.T) {
	if _, err := InstallZip("./assets/res.zip", "./resources"); err != nil {
		t.Error(err)
	}
}
