package main

import (
	"flag"
	"fmt"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
	"os/exec"
	"time"
)

var (
	dev     bool
	input   string
	webPort = 8080
	version = "v0.3.2"
)

func init() {
	flag.BoolVar(&dev, "dev", false, "开发模式")
	log.SetFormatter(&easy.Formatter{
		TimestampFormat: "2006-01-02 15:04:05",
		LogFormat:       "[排队姬][%time%][%lvl%]: %msg% \n",
	})
	log.SetLevel(log.InfoLevel)
}

func main() {
	var err error
	flag.Parse()
	s := NewServer()
	s.Init()
	defer s.DanmakuClient.Stop()
	if !dev {
		gin.SetMode(gin.ReleaseMode)
		webPort = 18303
		go checkUpdate()
	}
	router := gin.New()
	if dev {
		router.Use(CorsMiddleWare("http://localhost:8080"))
	}
	router.GET("/eio", s.Eio.Warp)
	router.GET("/api/sync", func(c *gin.Context) {
		c.JSON(200, s.Queue.Encode())
	})

	router.Use(static.Serve("/", static.LocalFile("./frontend/dist", false)))
	router.NoRoute(func(c *gin.Context) {
		c.File("./frontend/dist/index.html")
	})
	for {
		if checkPort(webPort) {
			break
		} else {
			log.Errorf("启动失败: 端口 %d 被占用, 将尝试更换端口启动", webPort)
			webPort += 1
			time.Sleep(500 * time.Millisecond)
		}
	}
	url := fmt.Sprintf("%s%d", "http://localhost:", webPort)
	err = exec.Command("cmd", "/C", "start", url).Run()
	if err != nil {
		log.Error("打开浏览器失败了, 请手动打开网址使用排队姬 ", url)
	}
	if err = router.Run(fmt.Sprintf(":%d", webPort)); err != nil {
		log.Error("启动失败: ", err)
	}
	log.Infof("按回车键退出...")
	_, _ = fmt.Scanln(&input)
}
