package main

import (
	"testing"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"fmt"
)

func TestTplConvContent(t *testing.T) {
	assBuf, _ := ioutil.ReadFile("./resources/template/sorry/template.ass")
	outputAss, _ := os.Create("./dist/output.ass")
	makeAss(string(assBuf), outputAss)
}

func makeAss(tplContentText string, fWriter io.Writer) {
	sentences := []string{
		"好啊",
		"就算你是一流程序员",
		"写出来的代码再完美",
		"我说这是 BUG 它就是 BUG",
		"毕竟我是用户",
		"你害我加班啊",
		"sorry 我就喜欢看程序猿加班",
		"以后天天找他 BUG",
		"天天找 天天找",
	}
	data := map[string][]string{
		"sentences": sentences,
	}
	tpl := template.New("subTitle")
	tpl, _ = tpl.Parse(tplContentText)
	tpl.Execute(fWriter, data)
}

func TestMakeVideo(t *testing.T)  {
	makeVideo()
}

func makeVideo() {
	cmd := exec.Command("ffmpeg", "-i", "./resources/template/sorry/template.mp4",
		"-vf", fmt.Sprintf("ass=%s", "./dist/output.ass"),
		"-an",
		"-y", "./dist/output.mp4")
	cmd.Start()
}