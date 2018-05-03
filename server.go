package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Server struct {
	bind   string
	router *gin.Engine
}

// Run 启动 web 服务
func (s *Server) Run() error {
	router := s.router

	router.GET("/", func(c *gin.Context) {
		if res, err := ScanAllTemplate(); err != nil {
			c.JSON(http.StatusInternalServerError, err)
		} else {
			c.JSON(http.StatusOK, map[string]interface{}{
				"res_count": len(res),
				"res":       res,
			})
		}
	})
	router.GET("/info/:tpl_key", func(c *gin.Context) {
		tplKey := c.Param("tpl_key")
		if r, err := ScanTemplate(tplKey); err != nil {
			c.JSON(http.StatusInternalServerError, err)
		} else {
			c.JSON(http.StatusOK, r)
		}
	})
	router.POST("/generate/:tpl_key/:res_type", func(c *gin.Context) {
		tplKey := c.Param("tpl_key")
		resType := c.Param("res_type")
		body := map[string][]string{}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusInternalServerError, err)
		} else {
			subs := Subs{}
			subs.Append(body["sentences"])
			if hash, err := GenerateResource(tplKey, subs, resType); err != nil {
				c.JSON(http.StatusInternalServerError, err)
			} else {
				c.JSON(http.StatusOK, map[string]string{
					"hash": hash,
				})
			}
		}
	})
	router.POST("/task/generate/:tpl_key", func(c *gin.Context) {
		tplKey := c.Param("tpl_key")
		body := map[string][]string{}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusInternalServerError, err)
		} else {
			subs := Subs{}
			subs.Append(body["sentences"])
			hash := addMakeTask(Task{RunnableList: []makeFunc{MakeMp4, MakeGif}, TplKey: tplKey, Subs: subs})
			c.JSON(http.StatusOK, map[string]interface{}{
				"hash":  hash,
				"state": taskState[hash],
			})
		}
	})
	router.GET("/task/generate/:hash", func(c *gin.Context) {
		hash := c.Param("hash")
		c.JSON(http.StatusOK, map[string]interface{}{
			"hash":  hash,
			"state": loadTaskState(hash),
		})
	})
	router.POST("/upload/res", func(c *gin.Context) {
		if file, err := c.FormFile("file"); err != nil {
			c.JSON(http.StatusInternalServerError, err)
		} else {
			if err := c.SaveUploadedFile(file, "./tmp"+"/"+file.Filename); err != nil {
				c.JSON(http.StatusInternalServerError, err)
			} else {
				if files, err := InstallZip("./tmp/"+file.Filename, "./resources"); err != nil {
					c.JSON(http.StatusInternalServerError, err)
				} else {
					c.JSON(http.StatusOK, map[string]interface{}{
						"make_files": files,
					})
				}
			}
		}
	})

	return router.Run(s.bind)
}
