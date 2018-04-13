package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"errors"
	"strings"
)

var (
	NotFoundTemplateProperties = errors.New("not found any .sentences template")
)

type ResInfo struct {
	TplKey         string   `json:"tpl_key"`
	Name           string   `json:"name"`
	Sentences      []string `json:"sentences"`
	SentencesCount int      `json:"sentences_count"`
}

const baseDir = "./resources/template"

func ScanAllTemplate() (rs []ResInfo, err error) {
	rs = []ResInfo{}
	if files, err := ioutil.ReadDir(baseDir); err != nil {
		return rs, err
	} else {
		for _, f := range files {
			if r, err := ScanTemplate(f.Name()); err == nil {
				rs = append(rs, r)
			}
		}
	}
	return
}

func ScanTemplate(tplKey string) (ri ResInfo, err error) {
	ri = ResInfo{TplKey: tplKey}
	basePath := fmt.Sprintf("%s/%s", baseDir, tplKey)
	assFilePath := fmt.Sprintf("%s/%s", basePath, "template.ass")
	videoFilePath := fmt.Sprintf("%s/%s", basePath, "template.mp4")
	if exist, err := IsAllExist(basePath, assFilePath, videoFilePath); !exist {
		return ri, err
	}
	// 读取模板内容
	// 扫描句子模板数量
	if tmpBuf, err := ioutil.ReadFile(assFilePath); err != nil {
		return ri, err
	} else {
		tmpAssContent := string(tmpBuf)
		if reg, err := regexp.Compile("{{\\s*index\\s*\\.sentences\\s*[0-9]+\\s*}}"); err != nil {
			return ri, err
		} else {
			results := reg.FindAllString(tmpAssContent, -1)
			if results == nil {
				return ri, NotFoundTemplateProperties
			} else {
				ri.SentencesCount = len(results)
			}
		}
	}

	// 读取 sentences 内容
	// 扫描每一条预设句子
	sentencesFilePath := fmt.Sprintf("%s/%s", basePath, "sentences")
	if exist, _ := IsExist(sentencesFilePath); !exist {
		ri.Sentences = []string{}
	} else {
		if tmpBuf, err := ioutil.ReadFile(sentencesFilePath); err != nil {
			return ri, err
		} else {
			tmpSentencesContent := string(tmpBuf)
			results := strings.Split(tmpSentencesContent, "\n")
			if results == nil {
				results = []string{}
			}
			ri.Sentences = results[:ri.SentencesCount]
		}
	}

	// 读取 name 属性
	nameFilePath := fmt.Sprintf("%s/%s", basePath, "name")
	if exist, _ := IsExist(nameFilePath); !exist {
		ri.Name = ri.TplKey
	} else {
		if tmpBuf, err := ioutil.ReadFile(nameFilePath); err != nil {
			return ri, err
		} else {
			ri.Name = string(tmpBuf)
		}
	}

	return ri, nil
}
