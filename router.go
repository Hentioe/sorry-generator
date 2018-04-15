package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Server struct {
	bind   string
	router *gin.Engine
}

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

	return router.Run(s.bind)
}
