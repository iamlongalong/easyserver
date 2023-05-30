package server

import (
	"easyserver/cmd/internal/model"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/gofrs/uuid"
)

var (
	// validator
	dfvalidator = validator.New()

	banner = `########################################`
)

// home dir
var homeDir = "/tmp/easyserver"

func Register(engine *gin.Engine) {

	fileServer := http.FileServer(http.Dir(homeDir))
	fileServer = http.StripPrefix("/", fileServer) // 去掉 URL 前面的斜杠

	// 开启文件列表功能

	engine.Use(AuthMiddleware()).
		GET("/ping", func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{
				"message": "pong",
			})
		}).
		POST("/_token", CreateToken).
		DELETE("/_token", DeleteToken)

	engine.NoRoute(func(c *gin.Context) {
		// GET 为 static file server
		if c.Request.Method == "GET" {
			gin.WrapF(func(w http.ResponseWriter, r *http.Request) {
				fileServer.ServeHTTP(w, r)
			})(c)
			return
		}

		// POST 为上传文件
		if c.Request.Method == "POST" {
			// 接收上传的文件
			file, err := c.FormFile("file")
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"message": err.Error(),
				})
				return
			}

			filename := filepath.Join(homeDir, c.Request.URL.Path)

			// 保存文件
			err = c.SaveUploadedFile(file, filename)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"message": err.Error(),
				})
				return
			} else {
				c.JSON(http.StatusOK, gin.H{
					"message": fmt.Sprintf("'%s' upload success", filename),
				})
				return
			}
		}

		// DELETE 为删除文件
		if c.Request.Method == "DELETE" {
			filename := filepath.Join(homeDir, c.Request.URL.Path)
			err := os.RemoveAll(filename)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"message": err.Error(),
				})
				return
			} else {
				c.JSON(http.StatusOK, gin.H{
					"message": fmt.Sprintf("'%s' delete success", filename),
				})
				return
			}
		}

		c.JSON(http.StatusNotFound, gin.H{"message": "not found"})
		return

	})
}

func CreateToken(c *gin.Context) {
	params := model.CreateTokenParams{}

	err := c.ShouldBindJSON(&params)
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}

	err = dfvalidator.Struct(params)
	if err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	iuser, ok := c.Get("user")
	if !ok {
		c.JSON(400, gin.H{"message": "user not found"})
		return
	}

	user, ok := iuser.(model.User)
	if !ok {
		c.JSON(400, gin.H{"message": "user not found"})
		return
	}

	rmaps := make(map[string]model.PathRole)
	for _, pr := range params.PathRoles {
		rmaps[pr.Path] = pr
	}

	d, err := time.ParseDuration(params.Duration)
	if err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	token := model.TokenAuth{
		Token:            genUUID(),
		SignedUser:       user,
		PathRoles:        rmaps,
		SignAt:           time.Now(),
		Duration:         d,
		UploadSizeLimit:  params.UploadSizeLimit,
		UploadCountLimit: params.UploadCountLimit,
	}

	err = createToken(&user, token)
	if err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"token": token.Token,
	})
}

func Serve(server model.Server, uas []model.UserAuths, anny model.Annymous) {
	// init user auths
	err := InitUserAuths(uas)
	if err != nil {
		log.Fatal(err)
	}

	// init annymous auths
	err = InitAnnymousAuths(anny)
	if err != nil {
		log.Fatal(err)
	}

	// init user tokens
	err = InitUserTokens()
	if err != nil {
		log.Fatal(err)
	}

	homeDir = filepath.Base(server.Home)
	// 检查 homeDir 是否存在，不存在则创建
	fi, err := os.Stat(homeDir)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(homeDir, 0755)
		}

		if err != nil {
			log.Fatal(err)
		}
	}

	if !fi.IsDir() {
		log.Fatalf("%s is not a dir", homeDir)
	}

	// 检查 https
	if server.Https {
		// 检查 https 证书
		if server.Cert == "" || server.Key == "" {
			log.Fatal("https cert file or key file is empty")
		}

		// 检查 https 证书是否存在
		_, err = os.Stat(server.Cert)
		if err != nil {
			log.Fatalf("https cert file [%s] not exist", server.Cert)
		}

		_, err = os.Stat(server.Key)
		if err != nil {
			log.Fatalf("https key file [%s] not exist", server.Key)
		}
	}

	// 根据 ENV 设置 gin 的模式
	if os.Getenv("ENV") != "DEBUG" {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.Default()
	Register(engine)

	addr := fmt.Sprintf("%s:%d", server.Host, server.Port)

	// print users info
	fmt.Println(banner)
	for _, u := range uas {
		proles := ""
		for _, r := range u.PathRoles {
			proles += fmt.Sprintf("%s:%s\t", r.Path, r.Mode)
		}

		fmt.Printf("user: %s , paths:  %s\n", u.User.UserName, proles)
	}
	fmt.Println(banner)
	if anny.Enable {
		fmt.Printf("annymous user: %s:%s\n", anny.Path, anny.Mode)
		fmt.Println(banner)
	}

	// print server info
	log.Printf("server listen on %s", addr)

	// 设置 https
	if server.Https {
		err = engine.RunTLS(addr, server.Cert, server.Key)
	} else {
		err = engine.Run(addr)
	}

	if err != nil {
		log.Fatal(err)
	}
}

func DeleteToken(c *gin.Context) {
	user, ok := c.MustGet("user").(model.User)
	if !ok {
		c.JSON(400, gin.H{"message": "user not found"})
		return
	}

	deleteToken(user, c.Query("token"))

	c.JSON(200, gin.H{"message": "ok"})
}

func genUUID() string {
	u, err := uuid.NewV4()
	if err != nil {
		return genUUID()
	}

	return u.String()
}
