/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/iamlongalong/easyserver/cmd/internal/model"
	"github.com/iamlongalong/easyserver/cmd/internal/server"
	"github.com/kardianos/service"

	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	configFile *string
	addr       *string
	httpsStr   *string
	homeDir    *string
	anyStr     string
	usersStrs  *[]string
	daemon     *bool
)

func init() {
	rootCmd.AddCommand(serveCmd)

	// Define a flag for the configuration file
	configFile = serveCmd.PersistentFlags().StringP("config", "c", "", "Config file path, eg: --addr /etc/easyserver.yaml")

	addr = serveCmd.Flags().String("addr", "0.0.0.0:8080", "Server addr, eg: --addr 127.0.0.1:443")

	httpsStr = serveCmd.Flags().String("https", "", "Server use https, eg: --https file.longalong.com:/path/to/cert.pem:/path/to/key.pem")

	homeDir = serveCmd.Flags().String("home", "", "server home dir, eg: --home .")

	// 匿名访问
	serveCmd.Flags().StringVar(&anyStr, "any", "", `enable annymous user (defult diasble), eg: --any=r:/path1,w:path2 (mode: r: read, w: write), and '--any' equals '--any=r:/'`)

	// 用户认证
	usersStrs = serveCmd.Flags().StringArray("user", []string{}, `users auth, eg: --user username:password:r:/path/to/dir:w:/path2 (r: read, w: write)`)

	// 若开启匿名访问，默认为 / r 权限
	serveCmd.Flags().Lookup("any").NoOptDefVal = "r:/"

	// 后台运行
	daemon = serveCmd.PersistentFlags().BoolP("daemon", "d", false, "run as daemon, eg: -d")

}

var (
	dfvalidator = validator.New()
)

// serveCmd represents the server command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "start file server",
	Long: `server usage:

default: 
				easyserver serve .

with host and port: 
				easyserver serve . --addr 0.0.0.0:8081

with users:
				easyserver serve . --user username:password:r:/path/to/dir:w:/path2

				multi users: --user user1:password:r:/path/to/dir:w:/path2 --user user2:password:r:/path/to/dir:w:/path2

	`,
	Run: func(cmd *cobra.Command, args []string) {
		serviceConfig := model.ServieConfig{}

		var err error
		// bind config to serverConfig
		if *configFile != "" {
			viper.SetConfigFile(*configFile)
			err := viper.ReadInConfig()
			if err != nil {
				log.Fatal(errors.Wrap(err, "read config error"))
			}
		}

		// bind args 1 to homedir if homedir is nil
		if *homeDir == "" && len(args) > 0 {
			*homeDir = args[0]
		}

		// bind config file to service config
		err = viper.Unmarshal(&serviceConfig)
		if err != nil {
			log.Fatal(errors.Wrap(err, "unmarshal config error"))
		}

		// bind flag to service config

		// convert user config
		for _, usrStr := range *usersStrs {
			usr := strings.Split(usrStr, ":")
			user := model.User{}
			auths := model.UserAuths{
				PathRoles: map[string]model.PathRole{},
			}

			if len(usr) < 2 {
				log.Fatal(errors.Errorf("invalid user format: %s", usrStr))
			}

			user.UserName = usr[0]
			user.Password = usr[1]

			auths.User = user

			switch len(usr) {
			case 2:
				// 默认为 write / 权限
				auths.PathRoles["/"] = model.PathRole{Path: "/", Mode: "w"}

			default:
				// 参数数量必须为偶数
				if len(usr)%2 != 0 {
					log.Fatal(errors.Errorf("invalid user format: %s", usrStr))
				}

				for i := 2; i < len(usr); i += 2 {
					auths.PathRoles[usr[i]] = model.PathRole{Path: usr[i], Mode: usr[i+1]}
				}

			}

			serviceConfig.Users = append(serviceConfig.Users, auths)
		}

		// bind server config
		if *addr != "" {
			serviceConfig.Server.Addr = *addr
		}

		if *httpsStr != "" {
			https := strings.Split(*httpsStr, ":")
			if len(https) != 3 {
				log.Fatal(errors.Errorf("invalid https format: %s", *httpsStr))
			}
			serviceConfig.Server.Https = model.Https{
				Domain: https[0],
				Cert:   https[1],
				Key:    https[2],
			}
		}

		if *homeDir != "" {
			serviceConfig.Server.Home = filepath.Clean(*homeDir)
		}

		if anyStr != "" {
			anyPRs := strings.Split(anyStr, ",")
			for _, annStr := range anyPRs {
				anys := strings.Split(annStr, ":")

				// any 的参数需要为偶数
				if len(anys)%2 != 0 {
					log.Fatal(errors.Errorf("invalid any format: %s", annStr))
				}

				serviceConfig.Any.Enable = true

				for i := 0; i < len(anys); i += 2 {
					serviceConfig.Any.Paths = append(serviceConfig.Any.Paths,
						model.PathRole{Mode: anys[i], Path: anys[i+1]},
					)
				}
			}
		}

		// 若无 user， 则默认创建一个 admin
		if len(serviceConfig.Users) == 0 {

			serviceConfig.Users = append(serviceConfig.Users,
				model.UserAuths{
					User:      model.User{UserName: "admin", Password: "easyadmin123"},
					PathRoles: map[string]model.PathRole{"/": {Path: "/", Mode: "w"}},
				},
			)

			fmt.Printf("no user setted, using default user:  admin:easyadmin123:w:/ 🤐 \n\n")
		}

		err = dfvalidator.Struct(serviceConfig)
		if err != nil {
			log.Fatal(errors.Wrap(err, "validate config error"))
		}

		sservice, err := server.BuildServe(serviceConfig)
		if err != nil {
			log.Fatal(errors.Wrap(err, "build service error"))
		}

		if *daemon {
			// 使用后台运行
			err := runAsService(sservice)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			err := sservice.StartSync(cmd.Context())
			if err != nil {
				log.Fatal(err)
			}
		}

	},
}

func runAsService(server service.Interface) error {
	currentUser, err := user.Current()
	if err != nil {
		return err
	}

	cfg := service.Config{
		Name:         "easyserver",
		UserName:     currentUser.Username,
		Description:  "easyserver is a very easy static server, with a easy dashboard",
		Dependencies: []string{"After=network.target syslog.target"},
		Arguments:    os.Args[1:],
	}

	s, err := service.New(server, &cfg)
	if err != nil {
		return err
	}

	err = s.Install()
	if err != nil {
		if !strings.Contains(err.Error(), "Init already exists") {
			return err
		}
		// 已经注册了，就放过
	}

	err = s.Start()
	if err != nil {
		return err
	}

	return nil
}
