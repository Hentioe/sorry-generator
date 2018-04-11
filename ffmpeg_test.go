package main

import (
	"testing"
)

func TestGeneratorResource(t *testing.T) {
	tplKey := "dagong"
	subs := Subs{}
	subs.Append("没有钱啊肯定要做啊").
		Append("不做的话又没有钱用").
		Append("那你不会打工啊").
		Append("有手有脚的").
		Append("打工是不可能打工的").
		Append("这辈子不可能打工的")

	if _, err := MakeMp4(tplKey, subs); err != nil {
		t.Error(err)
	}
	if _, err := MakeGif(tplKey, subs); err != nil {
		t.Error(err)
	}
}
