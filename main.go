package main

import (
	"runtime"
	"os"
	"path/filepath"
	"log"
	"fmt"
	"flag"
	"github.com/gin-gonic/gin"
)

func init() {
	if parentPath, err := filepath.Abs(filepath.Dir(os.Args[0])); err != nil {
		log.Fatal(err)
	} else {
		distDir := parentPath + "/dist"
		resourcesDir := parentPath + "/resources"
		tmpDir := parentPath + "/tmp"
		if err := IfNotExistMkAllMir(0774, distDir, resourcesDir, tmpDir); err != nil {
			log.Fatal(err)
		}
	}
}

var bind = flag.String("bind", ":8080", "bind address and port")
var installRes = flag.String("i", "", "install resources for a zip file")
var mode = flag.String("mode", "test", "running mode, e.g. debug/test/release")
var cl = flag.Int("cl", runtime.NumCPU(), "concurrency limits")

func main() {
	flag.Parse()
	if *installRes != "" {
		if _, err := InstallZip(*installRes, "./resources"); err != nil {
			fmt.Printf("install template resources failed, %s\n", err)
			os.Exit(1)
		}
		fmt.Println("install template resources succcess.")
		os.Exit(0)
	}
	gin.SetMode(*mode)
	server := Server{router: gin.Default(), bind: *bind}
	go asyncMakeAction()
	log.Fatal(server.Run())
}
