package model

import (
	"path/filepath"
	"strings"
	"time"
)

type Server struct {
	Addr  string `json:"addr" mapstructure:"addr"`
	Https Https  `json:"https" mapstructure:"https"`

	Home string `json:"home" mapstructure:"home" validate:"required"`
}

type Https struct {
	Domain string `json:"domain" mapstructure:"domain"`
	Cert   string `json:"cert" mapstructure:"cert"`
	Key    string `json:"key" mapstructure:"key"`
}

type Annymous struct {
	Enable bool       `json:"enable" mapstructure:"enable"`
	Paths  []PathRole `json:"paths" mapstructure:"paths"`
}

type User struct {
	UserName string `json:"username" mapstructure:"username" validate:"required"`
	Password string `json:"password" mapstructure:"password" validate:"required"`
}

type PathRole struct {
	Path string `json:"path" mapstructure:"path"`
	Mode string `json:"mode" mapstructure:"mode"`
}

// Valid check PathRole, and clean path
func (pr *PathRole) Valid() bool {
	if pr.Path == "" || pr.Mode == "" {
		return false
	}

	pr.Path = filepath.Clean(pr.Path)

	return (pr.Mode == "r" || pr.Mode == "w")
}

func (pr *PathRole) HasPermit(p string, mode string) bool {
	p = filepath.Clean(p)

	// 仅通过前缀判断权限，不考虑通配符
	if strings.HasPrefix(p, pr.Path) {
		if mode == "r" {
			return true
		}

		if mode == "w" && pr.Mode == "w" {
			return true
		}

		return false
	}

	return false
}

// UserAuths user 对路径的权限
type UserAuths struct {
	User      `json:"user"`
	PathRoles map[string]PathRole `json:"path_roles"`
}

type TokenAuth struct {
	Token string `json:"token"`

	// SignedUser 签发者
	SignedUser User `json:"signed_user"`

	// PathRole is a map of path and role
	PathRoles map[string]PathRole `json:"path_roles"`

	// SignAt token 签发时间
	SignAt time.Time `json:"sign_at"`

	// Duration token 有效期
	Duration time.Duration `json:"duration"`

	// UploadSizeLimit 单次上传文件大小限制 eg: 100MB
	UploadSizeLimit string `json:"size_limit"`

	// CountLimit 上传次数限制
	UploadCountLimit int `json:"count_limit"`

	// UsedCount 已使用过的次数
	UploadUsedCount int `json:"used_count"`
}

// create Token params
type CreateTokenParams struct {
	PathRoles        []PathRole `json:"path_roles" validate:"dive,required"`
	Duration         string     `json:"duration" validate:"required"`
	UploadSizeLimit  string     `json:"size_limit"`
	UploadCountLimit int        `json:"count_limit"`
}

type ServieConfig struct {
	Server
	Users []UserAuths `json:"users" mapstructure:"users"`
	Any   Annymous    `json:"any" mapstructure:"any"`
	// Admin User        `json:"admin" mapstructure:"admin"`
}
