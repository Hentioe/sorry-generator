package main

import (
	"fmt"
	"os/exec"
	"io/ioutil"
	"html/template"
	"os"
	"errors"
)

func MakeGif(tplKey string, subs Subs) (string, error) {
	return GenerateResource(tplKey, subs, "gif")
}
func MakeMp4(tplKey string, subs Subs) (string, error) {
	return GenerateResource(tplKey, subs, "mp4")
}

// GenerateResource Generate resources(gif/mp4)
// ffmpeg CLI wrapper
func GenerateResource(tplKey string, subs Subs, resType string) (hash string, err error) {
	tplPath := fmt.Sprintf("./resources/template/%s", tplKey)
	videoTplFile := tplPath + "/template.mp4"
	subTplFile := tplPath + "/template.ass"
	hash = subs.Hash(tplKey)
	subOutputFile := fmt.Sprintf("./dist/%s", hash)
	outputResource := "dist/" + hash + "." + resType
	if _, err = os.Stat(outputResource); os.IsNotExist(err) {
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
		var cmd = &exec.Cmd{}
		switch resType {
		case "gif":
			cmd = exec.Command("ffmpeg", "-i", videoTplFile,
				"-vf", fmt.Sprintf("ass=%s,scale=300:-1", subOutputFile),
				"-r", "8",
				"-y", outputResource)
		case "mp4":
			cmd = exec.Command("ffmpeg", "-i", videoTplFile,
				"-vf", fmt.Sprintf("ass=%s", subOutputFile),
				"-an",
				"-y", outputResource)
		default:
			return "", errors.New("Unknown resType: " + resType)
		}
		if _, err := cmd.CombinedOutput(); err != nil {
			return hash, err
		}
	} else {
		return
	}
	return
}
