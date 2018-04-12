package main

import (
	"os"
	"path/filepath"
	"log"
	"fmt"
	"github.com/emicklei/go-restful"
	"net/http"
	"flag"
)

func init() {
	if parentPath, err := filepath.Abs(filepath.Dir(os.Args[0])); err != nil {
		log.Fatal(err)
	} else {
		distDir := parentPath + "/dist"
		if _, err := os.Stat(distDir); os.IsNotExist(err) {
			fmt.Println(distDir)
			os.Mkdir(distDir, os.FileMode(0774))
		}
	}
}

var bind = flag.String("bind", ":8080", "Bind address and port")

func main() {
	flag.Parse()

	wsContainer := restful.NewContainer()
	wsContainer.Add(func() *restful.WebService {
		var ws = new(restful.WebService)
		ws.Path("/")
		ws.Route(ws.GET("/").To(func(request *restful.Request, response *restful.Response) {
			if res,err := ScanAllTemplate();err!= nil {
				response.WriteError(500, err)
			}else{
				response.WriteAsJson(map[string]interface{}{
					"res_count": len(res),
					"res":       res,
				})
			}
		}))
		ws.Route(ws.GET("/info/{tplKey}").To(func(request *restful.Request, response *restful.Response) {
			tplKey := request.PathParameter("tplKey")
			if r, err := ScanTemplate(tplKey); err != nil {
				response.WriteError(500, err)
			} else {
				response.WriteAsJson(r)
			}
		}))
		ws.Route(
			ws.POST("/generate/{tplKey}/{resType}").To(func(request *restful.Request, response *restful.Response) {
				tplKey := request.PathParameter("tplKey")
				resType := request.PathParameter("resType")
				body := map[string][]string{}
				request.ReadEntity(&body)
				subs := Subs{}
				subs.Append(body["sentences"])
				if hash, err := GenerateResource(tplKey, subs, resType); err != nil {
					response.WriteError(500, err)
				} else {
					response.WriteAsJson(map[string]string{
						"hash": hash,
					})
				}
			}))
		return ws
	}())

	server := &http.Server{Addr: *bind, Handler: wsContainer}
	log.Fatal(server.ListenAndServe())
}
