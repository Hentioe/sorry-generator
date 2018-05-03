package main

import (
	"strconv"
	"errors"
	"fmt"
	"testing"

	"github.com/go-resty/resty"
)

func TestConcurrency(t *testing.T) {
	for i :=16; i > 0; i-- {
		resp, err := resty.R().
			SetHeader("Content-Type", "application/json").
			SetBody(`{"sentences":["` + strconv.Itoa(i) + `第一句","第二句","第三句","第四句","第五句","第六句","第七句","第八句","第九句"]}`).
			Post("http://localhost:8080/task/generate/sorry")
		if err != nil {
			t.Error(err)
		}
		if resp.StatusCode() != 200 {
			t.Error(errors.New(fmt.Sprintf("http err status: %d", resp.StatusCode())))
		}
	}
}
