package server

import (
	"easyserver/cmd/internal/model"
	"fmt"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// user auths by username and password
var userAuths = map[string]model.UserAuths{}

// user token
var userTokens = map[string]model.TokenAuth{}
var tokenLock = sync.Mutex{}

// annymous auths
var annymousAuths model.Annymous

// InitUserAuths init user auths
func InitUserAuths(uas []model.UserAuths) error {
	userAuths = map[string]model.UserAuths{}
	for _, ua := range uas {
		for _, role := range ua.PathRoles {
			if !role.Valid() {
				return errors.Errorf("invalid role [%s %s]", role.Path, role.Mode)
			}
		}
		userAuths[ua.User.UserName] = ua
	}

	return nil
}

// InitUserTokens init user tokens
func InitUserTokens() error {
	return nil
}

// InitAnnymousAuths init annymous auths
func InitAnnymousAuths(anny model.Annymous) error {
	if anny.Enable && !anny.Valid() {
		return errors.Errorf("invalid annymous auth: %s : %s", anny.Path, anny.Mode)
	}
	annymousAuths = anny
	return nil
}

var (
	ReadMode  = "r"
	WriteMode = "w"
)

// createToken auth
func createToken(u *model.User, t model.TokenAuth) error {
	// check user auth
	for tp, role := range t.PathRoles {
		if !userAuth(*u, tp, role.Mode) {
			return errors.Errorf("user [%s] not have path [%s] auth", u.UserName, t)
		}
	}

	// register token
	tokenLock.Lock()
	userTokens[t.Token] = t
	tokenLock.Unlock()

	return nil
}

func deleteToken(user model.User, token string) error {
	tokenLock.Lock()
	defer tokenLock.Unlock()

	t, ok := userTokens[token]
	if !ok {
		return errors.Errorf("token [%s] not found", token)
	}

	if t.SignedUser.UserName != user.UserName {
		return errors.Errorf("token [%s] not match user [%s]", token, user.UserName)
	}

	delete(userTokens, token)
	return nil
}

// AuthMiddleware is a middleware for static auth
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("running in auth middleware")
		// 检查接口是否需要权限
		path := c.Request.URL.Path
		method := c.Request.Method

		needRoleMode := ReadMode
		if method == "POST" || method == "PUT" || method == "DELETE" {
			needRoleMode = WriteMode
		}

		if path == "/ping" {
			c.Next()
			return
		}

		// check user auth
		// get basic auth from header
		uname, password, ok := c.Request.BasicAuth()
		if ok {
			ua, ok := userAuths[uname]
			if ok && ua.User.Password == password {
				// check user auth
				user := model.User{UserName: uname}
				c.Set("user", user)

				if userAuth(model.User{UserName: uname}, path, needRoleMode) {
					c.Next()
					return
				}
			}

		}

		// check annymous auth
		if annymousAuth(path, needRoleMode) {
			c.Next()
			return
		}

		// check token auth
		token := c.Request.URL.Query().Get("token")
		if tokenAuth(token, path, needRoleMode) {
			c.Next()
			return
		}

		c.Header("WWW-Authenticate", `Basic realm="Restricted"`)
		c.AbortWithStatusJSON(401, gin.H{
			"code":    401,
			"message": "unauthorized",
		})
	}
}

// userAuth
func userAuth(user model.User, path string, mode string) bool {
	userAuth, ok := userAuths[user.UserName]
	if !ok {
		return false
	}

	for _, pathRole := range userAuth.PathRoles {
		if pathRole.HasPermit(path, mode) {
			return true
		}
	}

	return false
}

// annymousAuth
func annymousAuth(path string, mode string) bool {
	if !annymousAuths.Enable {
		return false
	}

	if annymousAuths.HasPermit(path, mode) {
		return true
	}

	return false
}

// tokenAuth
func tokenAuth(token string, path string, mode string) bool {
	tokenLock.Lock()
	tokenAuth, ok := userTokens[token]
	tokenLock.Unlock()
	if !ok {
		return false
	}

	for _, pathRole := range tokenAuth.PathRoles {
		if pathRole.HasPermit(path, mode) {
			return true
		}
	}

	return false
}