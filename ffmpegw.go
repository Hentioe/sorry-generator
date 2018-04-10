package main

import (
	"fmt"
	"os/exec"
	"io/ioutil"
	"html/template"
	"os"
)

// ffmpeg CLI wrapper
// ffmpeg -i <video_file> -vf ass=<ass_file> -an <output_file>

func GeneratorToMp4(tplKey string, subs Subs) (hash string, err error) {
	tplPath := fmt.Sprintf("./resources/template/%s", tplKey)
	videoTplFile := tplPath + "/template.mp4"
	subTplFile := tplPath + "/template.ass"
	hash = subs.Hash()
	subOutputFile := fmt.Sprintf("./dist/%s", hash)
	videoOutputFile := "dist/" + hash + ".mp4"
	if _, err = os.Stat(videoOutputFile); os.IsNotExist(err) {
		tplText := ""
		if tmpBuf, err := ioutil.ReadFile(subTplFile); err != nil {
			return hash, err
		} else {
			tplText = string(tmpBuf)
		}
		t := template.New("subTitle")
		if t, err = t.Parse(tplText); err != nil {
			return
		} else {
			if f, err := os.Create(subOutputFile); err != nil {
				return hash, err
			} else {
				data := map[string][]string{
					"sentences": subs.EntrySet(),
				}
				if err = t.Execute(f, data); err != nil {
					return hash, err
				}
			}
		}
		cmd := exec.Command("ffmpeg", "-i", videoTplFile,
			"-vf", fmt.Sprintf("ass=%s", subOutputFile),
			"-an",
			videoOutputFile)
		if _, err := cmd.CombinedOutput(); err != nil {
			return hash, err
		}
	} else {
		return
	}
	return
}
