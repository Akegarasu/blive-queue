package main

import (
	"flag"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
	"os/exec"
)

var (
	development bool
	webPort     = "8080"
	version     = "v0.0.2"
)

func init() {
	flag.BoolVar(&development, "dev", false, "开发模式")
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
	defer s.DanmakuClient.Stop()
	s.Init()
	if !development {
		gin.SetMode(gin.ReleaseMode)
		webPort = "18303"
		go checkUpdate()
	}
	router := gin.New()
	if development {
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

	err = exec.Command("cmd", "/C", "start", "http://localhost:"+webPort).Run()
	if err != nil {
		log.Error("打开浏览器失败了, 请手动打开网址使用排队姬 ", "http://localhost:"+webPort)
	}
	if err = router.Run(":18303"); err != nil {
		log.Fatal("启动失败: 请检查是否已经启动了~")
	}
}
