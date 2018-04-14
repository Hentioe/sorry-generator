package main

import (
	"path/filepath"
	"os"
	"io"
	"archive/zip"
	"strings"
	"fmt"
	"bytes"
	"regexp"
	"io/ioutil"
)

func Unzip(src string, dest string) ([]string, error) {

	var fileNames []string

	if exist, err := IsAllExist(src, dest); !exist {
		return fileNames, err
	}

	r, err := zip.OpenReader(src)
	if err != nil {
		return fileNames, err
	}
	defer r.Close()

	for _, f := range r.File {

		rc, err := f.Open()
		if err != nil {
			return fileNames, err
		}

		// Store filename/path for returning and using later on
		fPath := filepath.Join(dest, f.Name)
		if exist, _ := IsExist(fPath); exist {
			continue
		}
		// 构建模板
		if err := makeTpl(fPath, &rc); err != nil {
			return fileNames, err
		}

		fileNames = append(fileNames, fPath)

		if f.FileInfo().IsDir() {

			// Make Folder
			os.MkdirAll(fPath, os.ModePerm)

		} else {

			// Make File
			if err = os.MkdirAll(filepath.Dir(fPath), os.ModePerm); err != nil {
				return fileNames, err
			}

			outFile, err := os.OpenFile(fPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return fileNames, err
			}

			_, err = io.Copy(outFile, rc)

			// Close the file without defer to close before next iteration of loop
			outFile.Close()

			if err != nil {
				return fileNames, err
			}

		}
		rc.Close()
	}
	return fileNames, nil
}

func makeTpl(fPath string, rc *io.ReadCloser) error {

	if strings.HasSuffix(fPath, ".ass") {
		// 重写 .ass 内容
		assContentBuf := new(bytes.Buffer)
		assContentBuf.ReadFrom(*rc)
		assContent := assContentBuf.String()
		// Dialogue: 0,0:00:01.18,0:00:01.56,sorry,,0,0,0,,{{ index .sentences 0 }}
		i := 0
		var newLines []string
		var sentences []string
		var name string
		for _, line := range strings.Split(assContent, "\n") {
			var newLine string
			// 匹配字幕内容
			reg := regexp.MustCompile("^Dialogue.+sorry,,0,0,0,,(.+)$")
			if results := reg.FindStringSubmatch(line); len(results) > 0 {
				sentence := results[1]
				sentenceTpl := fmt.Sprintf("{{ index .sentences %d }}", i)
				i++
				newLine = strings.Replace(line, sentence, sentenceTpl, -1)
				sentences = append(sentences, sentence)
			} else {
				newLine = line
			}
			newLines = append(newLines, newLine)
			// 截取 title 属性
			reg = regexp.MustCompile("^Title:\\s*(.+)$")
			if results := reg.FindStringSubmatch(line); len(results) > 0 {
				name = results[1]
			}

		}
		assTplContent := strings.Join(newLines, "\n")
		*rc = ioutil.NopCloser(bytes.NewBuffer([]byte(assTplContent)))
		if pPath, err := filepath.Abs(filepath.Dir(fPath)); err != nil {
			return err
		} else {
			// 创建 sentences 文件
			if err := ioutil.WriteFile(pPath+"/sentences", []byte(strings.Join(sentences, "\n")), 0755); err != nil {
				return err
			}
			// 创建 name 文件
			if name != "" {
				if err := ioutil.WriteFile(pPath+"/name", []byte(name), 0755); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
