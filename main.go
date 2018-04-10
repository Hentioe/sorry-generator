package main

import (
	"os"
	"path/filepath"
	"log"
	"fmt"
	"github.com/emicklei/go-restful"
	"net/http"
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

func main() {
	wsContainer := restful.NewContainer()
	wsContainer.Add(func() *restful.WebService {
		var ws = new(restful.WebService)
		ws.Path("/")
		ws.Route(ws.POST("/generate/{tplKey}").To(func(request *restful.Request, response *restful.Response) {
			tplKey := request.PathParameter("tplKey")
			body := map[string][]string{}
			request.ReadEntity(&body)
			subs := Subs{}
			subs.Append(body["sentences"])
			if hash, err := GeneratorToMp4(tplKey, subs); err != nil {
				response.WriteError(500, err)
			} else {
				response.WriteAsJson(map[string]string{
					"hash": hash,
				})
			}
		}))
		return ws
	}())
	// Add container filter to enable CORS
	cors := restful.CrossOriginResourceSharing{
		ExposeHeaders:  []string{"*"},
		AllowedHeaders: []string{"*"},
		AllowedMethods: []string{"*"},
		AllowedDomains: []string{"*"},
		CookiesAllowed: false,
		Container:      wsContainer}
	wsContainer.Filter(cors.Filter)

	// Add container filter to respond to OPTIONS
	wsContainer.Filter(wsContainer.OPTIONSFilter)
	server := &http.Server{Addr: ":8080", Handler: wsContainer}
	log.Fatal(server.ListenAndServe())
}
