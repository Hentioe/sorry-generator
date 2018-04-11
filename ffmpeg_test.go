package main

import (
	"testing"
)

func TestGeneratorToMp4(t *testing.T) {
	tplKey := "sorry"
	subs := Subs{}
	subs.Append("好啊").
		Append("就算你是一流程序员").
		Append("写出来的代码再完美").
		Append("我说这是 BUG 它就是 BUG").
		Append("毕竟我是用户").
		Append("你害我加班啊").
		Append("sorry 我就喜欢看程序猿加班").
		Append("以后天天找他 BUG").
		Append("天天找 天天找")

	if _, err := MakeMp4(tplKey, subs); err != nil {
		t.Error(err)
	}
	if _, err := MakeGif(tplKey, subs); err != nil {
		t.Error(err)
	}
}
