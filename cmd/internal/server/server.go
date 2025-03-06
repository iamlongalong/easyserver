package server

import (
	"context"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/iamlongalong/easyserver/assets"
	"github.com/iamlongalong/easyserver/cmd/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var (
	// validator
	dfvalidator = validator.New()

	banner = `########################################`
)

// home dir
var homeDir = "/tmp/easyserver"

func CorsMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		ctx.Writer.Header().Set("Access-Control-Allow-Methods", "GET,POST,DELETE,PUT,OPTIONS")
		ctx.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type,Content-Length,Accept-Encoding,X-CSRF-Token,Authorization")
		ctx.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(204)
			return
		}

		ctx.Next()
	}
}

// get file infos from path
func handleGetFileInfos(c *gin.Context) {
	path := c.Param("path")

	tarDir := filepath.Join(homeDir, filepath.Join("/", filepath.Clean(path)))
	// list files of tarDir

	dir, err := os.Open(tarDir)
	if err != nil {
		if os.IsNotExist(err) {
			c.JSON(200, gin.H{
				"data": []string{},
			})
			return
		}

		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}

	fileInfos, err := dir.Readdir(-1)
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}

	resFileInfos := []model.ResFileInfo{}

	// convert to file info
	for _, fi := range fileInfos {
		// 过滤所有 隐藏文件
		if strings.HasPrefix(fi.Name(), ".") {
			continue
		}

		fileType := ""

		if fi.IsDir() {
			fileType = "folder"
		} else {
			fileType = convertExtContentType(filepath.Ext(fi.Name()))
		}

		resFileInfos = append(resFileInfos, model.ResFileInfo{
			Name:         fi.Name(),
			Size:         fi.Size(),
			ModTimeStamp: fi.ModTime().UnixMilli(),
			IsDir:        fi.IsDir(),
			FileType:     fileType,
		})
	}

	c.JSON(200, gin.H{
		"message": "ok",
		"data":    resFileInfos,
	})
}

func Register(engine *gin.Engine) {

	homedirServe := func(c *gin.Context) {
		relativePath := filepath.Join("/", filepath.Clean(c.Request.URL.Path))

		// If path ends with "/", treat it as directory only
		if strings.HasSuffix(c.Request.URL.Path, "/") {
			fullPath := filepath.Join(homeDir, relativePath)
			fi, err := os.Stat(fullPath)
			if err == nil && fi.IsDir() {
				http.FileServer(http.Dir(homeDir)).ServeHTTP(c.Writer, c.Request)
				return
			}
			c.JSON(http.StatusNotFound, gin.H{"message": "directory not found"})
			return
		}

		// Try paths in sequence: xx => xx.html => xx/index.html
		pathsToTry := []string{
			filepath.Join(homeDir, relativePath),
		}

		// Only try alternate paths if the original path has no extension
		if filepath.Ext(relativePath) == "" {
			pathsToTry = append(
				pathsToTry,
				filepath.Join(homeDir, relativePath+".html"),
				filepath.Join(homeDir, relativePath),
			)
		}

		// Try each path in sequence
		for _, tryPath := range pathsToTry {
			fi, err := os.Stat(tryPath)
			if err == nil {
				// If it's a directory, only serve if we're trying the original path
				if fi.IsDir() && tryPath == pathsToTry[0] {
					continue
				}

				http.ServeFile(c.Writer, c.Request, tryPath)
				return
			}
		}

		// If we get here, none of the paths worked
		c.JSON(http.StatusNotFound, gin.H{"message": "not found"})
	}

	engine.Use(CorsMiddleware()).Use(AuthMiddleware())

	engine.
		GET("/ping", func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{
				"message": "pong",
			})
		}).
		GET("/_apilist/*path", handleGetFileInfos).
		POST("/_token", CreateToken).
		DELETE("/_token", DeleteToken)

	// f := gin.WrapH(http.FileServer(http.FS(assetsFS)))

	engine.Group("/_dash").Any("/*fp", func(ctx *gin.Context) {
		fp := ctx.Param("fp")
		if fp == "/" || fp == "/index.html" {
			fn := "/_dash/index.html"

			f, err := http.FS(assets.DashBoardFS).Open(fn)
			if err != nil {
				ctx.AbortWithError(500, err)
				return
			}

			ctx.Writer.WriteHeader(200)
			ctx.Writer.Header().Set("Content-Type", "text/html")
			_, err = io.Copy(ctx.Writer, f)
			if err != nil {
				ctx.AbortWithError(500, err)
				return
			}
		} else {
			gin.WrapH(http.FileServer(http.FS(assets.DashBoardFS)))(ctx)
		}
	})

	engine.NoRoute(func(c *gin.Context) {
		// GET 为 static file server
		if c.Request.Method == "GET" || c.Request.Method == "HEAD" {
			c.Writer.WriteHeader(200) // fix the gin write status 404

			// if path is xx, find as follows: xx => xx.html => xx/index.html

			homedirServe(c)
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

			relativePath := filepath.Join("/", filepath.Clean(c.Request.URL.Path))

			filename := filepath.Join(homeDir, relativePath)

			// 保存文件
			err = c.SaveUploadedFile(file, filename)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"message": err.Error(),
				})
				return
			} else {
				c.JSON(http.StatusOK, gin.H{
					"message": fmt.Sprintf("'%s' upload success", relativePath),
				})
				return
			}
		}

		// DELETE 为删除文件
		if c.Request.Method == "DELETE" {
			relativePath := filepath.Join("/", filepath.Clean(c.Request.URL.Path))
			filename := filepath.Join(homeDir, relativePath)
			err := os.RemoveAll(filename)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"message": err.Error(),
				})
				return
			} else {
				c.JSON(http.StatusOK, gin.H{
					"message": fmt.Sprintf("'%s' delete success", relativePath),
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

	if len(params.PathRoles) == 0 {
		c.JSON(400, gin.H{"message": "pathRoles is empty"})
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
		Token:            randString(12),
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

func Serve(s model.ServieConfig) {

	ctx, cancel := context.WithCancel(context.Background())

	_ = ctx

	// init user auths
	err := InitUserAuths(s.Users)
	if err != nil {
		log.Fatal(err)
	}

	// init annymous auths
	err = InitAnnymousAuths(s.Any)
	if err != nil {
		log.Fatal(err)
	}

	// init user tokens
	err = InitUserTokens()
	if err != nil {
		log.Fatal(err)
	}

	if filepath.IsAbs(s.Home) {
		homeDir = filepath.Clean(s.Home)
	} else {
		homeDir = filepath.Base(s.Home)
	}
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

	fi, err = os.Stat(homeDir)
	if err != nil {
		log.Fatal(err)
	}

	if !fi.IsDir() {
		// log.Fatalf("%s is not a dir", homeDir)
		// support share one file
	}

	// 检查 https
	if s.Https.Cert != "" {
		// 检查 https 证书
		if s.Https.Cert == "" || s.Https.Key == "" {
			log.Fatal("https cert file or key file is empty")
		}

		// 检查 https 证书是否存在
		_, err = os.Stat(s.Https.Cert)
		if err != nil {
			log.Fatalf("https cert file [%s] not exist", s.Https.Cert)
		}

		_, err = os.Stat(s.Https.Key)
		if err != nil {
			log.Fatalf("https key file [%s] not exist", s.Https.Key)
		}
	}

	// 根据 ENV 设置 gin 的模式
	if os.Getenv("ENV") != "DEBUG" {
		gin.SetMode(gin.ReleaseMode)
	}

	var srv *http.Server

	engine := gin.Default()

	if s.CloseConf.MaxTimes > 0 {
		currentRunTimes := 0
		engine.Use(func(c *gin.Context) {
			c.Next()

			if c.Writer.Status() < 400 { // 视为一次成功请求
				currentRunTimes += 1
			}

			if currentRunTimes > s.CloseConf.MaxTimes {
				go func() {
					log.Printf("📢 server would be close by max times: %d in 60s\n", s.CloseConf.MaxTimes)
					cancel()
					srv.Close()

					time.Sleep(time.Minute)

					log.Printf("max wait of 60s, force exit now\n")
					os.Exit(0)
				}()
			}
		})
	}

	Register(engine)

	// print users info
	if len(s.Users) > 0 {
		fmt.Println(banner)
		for _, u := range s.Users {
			proles := ""
			for _, pr := range u.PathRoles {
				proles += fmt.Sprintf("%s:%s\t", pr.Path, pr.Mode)
			}

			fmt.Printf("user: %s , paths:  %s\n", u.User.UserName, proles)
		}
	}

	if s.Any.Enable {
		fmt.Println(banner)
		anyPaths := ""
		for _, pr := range s.Any.Paths {
			anyPaths += fmt.Sprintf("%s:%s\t", pr.Path, pr.Mode)
		}

		fmt.Printf("annymous user: %s\n", anyPaths)
		fmt.Println(banner)
	}

	// print server info
	if s.Https.Domain != "" {
		fmt.Printf("\nhttps domain: https://%s\n", s.Https.Domain)
	} else {
		fmt.Printf("\nserver listen on %s\n", s.Addr)
	}

	srv = &http.Server{
		Addr:    s.Addr,
		Handler: engine,
	}

	if s.CloseConf.MaxDuration > 0 {
		go func() {
			time.Sleep(s.CloseConf.MaxDuration)

			log.Printf("📢 server closed by max duration %s\n", s.CloseConf.MaxDuration)

			cancel()
			srv.Close()

			time.Sleep(time.Second * 60) // max wait
			log.Printf("max wait of 60s, force exit now\n")
			os.Exit(0)
		}()
	}

	// 设置 https
	if s.Https.Cert != "" {
		err = srv.ListenAndServeTLS(s.Https.Cert, s.Https.Key)
	} else {
		err = srv.ListenAndServe()
	}

	if err != nil {
		log.Fatal(err)
	}
}

func DeleteToken(c *gin.Context) {
	iuser, ok := c.Get("user")
	if !ok {
		c.JSON(400, gin.H{"message": "user not found"})
		return
	}

	user, ok := iuser.(model.User)
	if !ok {
		c.JSON(401, gin.H{"message": "user not found"})
		return
	}

	err := deleteToken(user, c.Query("token"), c.Query("all") != "")
	if err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "ok"})
}

// rand string
func randString(n int) string {
	letterRunes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890$_")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))] //nolint:gosec
	}
	return string(b)
}

// convert content-type to: "image", "video", "audio", "text", "bin"
func convertContentType(contentType string) string {
	contentType = strings.ToLower(contentType)
	if strings.Contains(contentType, "image") {
		return "image"
	}

	if strings.Contains(contentType, "video") {
		return "video"
	}

	if strings.Contains(contentType, "audio") {
		return "audio"
	}

	if strings.Contains(contentType, "text") {
		return "text"
	}

	return "file"
}

// 通过文件后缀，判断文件类型
func convertExtContentType(ext string) string {
	ext = strings.ToLower(ext)
	if strings.Contains(ext, "jpg") || strings.Contains(ext, "jpeg") || strings.Contains(ext, "png") {
		return "image"
	}

	if strings.Contains(ext, "mp4") || strings.Contains(ext, "avi") || strings.Contains(ext, "mov") {
		return "video"
	}

	if strings.Contains(ext, "mp3") || strings.Contains(ext, "wav") {
		return "audio"
	}

	if strings.Contains(ext, "txt") || strings.Contains(ext, "md") || strings.Contains(ext, "html") {
		return "text"
	}

	return "file"
}
