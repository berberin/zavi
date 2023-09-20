package main

import (
	"fmt"
	"net"
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

func GetLocalIPs() ([]net.IP, error) {
	var ips []net.IP
	addresses, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	for _, addr := range addresses {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			ips = append(ips, ipnet.IP)
		}
	}
	return ips, nil
}

func main() {
	pflag.IntP("port", "p", 9000, "port")
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	pwd, _ := os.Getwd()
	fmt.Println("Serving [ " + pwd + "/" + os.Args[1] + " ] at:")
	ips, _ := GetLocalIPs()
	for _, ip := range ips {
		if ip.To4() != nil {
			fmt.Println("http://" + ip.String() + ":" + viper.GetString("port"))
		} else {
			fmt.Println("http://[" + ip.String() + "]:" + viper.GetString("port"))
		}

	}

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(RerollMiddleware("/", os.Args[1]))
	r.Run(fmt.Sprintf(":%d", viper.GetInt("port")))
}
