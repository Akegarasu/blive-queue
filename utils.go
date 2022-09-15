package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
)

var (
	client    = &http.Client{}
	UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36 Edg/87.0.664.66"
)

func CorsMiddleWare(allowOrigin string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", allowOrigin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, Content-Length, X-CSRF-Token, Token, session, Origin, Host, Connection, Accept-Encoding, Accept-Language, X-Requested-With")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Request.Header.Del("Origin")
		c.Next()
	}
}

func GetJson(url string, headers map[string]string) (gjson.Result, error) {
	b, err := GetBytes(url, headers)
	if err != nil {
		return gjson.Result{}, err
	}
	result := gjson.ParseBytes(b)
	return result, nil
}

func GetBytes(url string, headers map[string]string) ([]byte, error) {
	reader, err := HTTPGetReadCloser(url, headers)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = reader.Close()
	}()
	return ioutil.ReadAll(reader)
}

func HTTPGetReadCloser(url string, headers map[string]string) (io.ReadCloser, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	if _, ok := headers["User-Agent"]; !ok {
		req.Header["User-Agent"] = []string{UserAgent}
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp.Body, err
}

func checkPort(port int) bool {
	s, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return false
	}
	err = s.Close()
	if err != nil {
		return false
	}
	return true
}

func checkUpdate() {
	log.Info("正在检查更新")
	updateURL := "https://api.github.com/repos/Akegarasu/blive-queue/releases/latest"
	j, err := GetJson(updateURL, nil)
	if err != nil {
		log.Error("检查更新失败~ 请手动检查更新")
		return
	}
	ver := j.Get("name").String()
	if ver != version && ver != "" {
		log.Info("---------------------------------------------")
		log.Infof("有新的版本可以更新啦: %s", ver)
		log.Infof("下载地址: %s", j.Get("html_url").String())
		log.Info("---------------------------------------------")
	}
}

func InSlice[T string | int](slice []T, elem T) bool {
	if slice == nil {
		return false
	}
	for _, v := range slice {
		if v == elem {
			return true
		}
	}
	return false
}
