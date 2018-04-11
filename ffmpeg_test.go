package main

import (
	"testing"
)

func TestGeneratorResource(t *testing.T) {
	tplKey := "wangjingze"
	subs := Subs{}
	subs.Append("我王境泽就是饿死").
		Append("死外边，从这里跳下去").
		Append("不会吃你们一点东西").
		Append("真香")

	if _, err := MakeMp4(tplKey, subs); err != nil {
		t.Error(err)
	}
	if _, err := MakeGif(tplKey, subs); err != nil {
		t.Error(err)
	}
}
