package main

import (
	"path/filepath"
	"os"
	"io"
	"archive/zip"
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
		filePath := filepath.Join(dest, f.Name)
		if exist, _ := IsExist(filePath); exist {
			continue
		}
		fileNames = append(fileNames, filePath)

		if f.FileInfo().IsDir() {

			// Make Folder
			os.MkdirAll(filePath, os.ModePerm)

		} else {

			// Make File
			if err = os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
				return fileNames, err
			}

			outFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
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
