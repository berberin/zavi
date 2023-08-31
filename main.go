package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func RerollMiddleware(urlPrefix, spaDirectory string) gin.HandlerFunc {
	directory := static.LocalFile(spaDirectory, true)
	fileserver := http.FileServer(directory)
	if urlPrefix != "" {
		fileserver = http.StripPrefix(urlPrefix, fileserver)
	}
	return func(c *gin.Context) {
		if directory.Exists(urlPrefix, c.Request.URL.Path) {
			fileserver.ServeHTTP(c.Writer, c.Request)
			c.Abort()
		} else {
			c.Request.URL.Path = "/"
			fileserver.ServeHTTP(c.Writer, c.Request)
			c.Abort()
		}
	}
}

func main() {
	pflag.IntP("port", "p", 9000, "port")
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	r := gin.Default()
	r.Use(RerollMiddleware("/", os.Args[1]))
	r.Run(fmt.Sprintf(":%d", viper.GetInt("port")))
}
