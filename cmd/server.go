/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"easyserver/cmd/internal/model"
	"easyserver/cmd/internal/server"
	"log"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	dfvalidator = validator.New()
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "server start",
	Long:  `server start`,
	Run: func(cmd *cobra.Command, args []string) {
		cfile := cmd.Flag("config").Value.String()

		// bind config to serverConfig
		if cfile != "" {
			viper.SetConfigFile(cfile)
			err := viper.ReadInConfig()
			if err != nil {
				log.Fatal(errors.Wrap(err, "read config error"))
			}
		}

		cmd.Flags().VisitAll(func(f *pflag.Flag) {
			viper.Set(f.Name, f.Value)
		})

		serviceConfig := model.ServieConfig{}
		err := viper.Unmarshal(&serviceConfig)
		if err != nil {
			log.Fatal(errors.Wrap(err, "unmarshal config error"))
		}

		// convert user config
		userPassStr := cmd.Flag("user").Value.String()
		if userPassStr != "" {
			userPassSlices := strings.Split(userPassStr, ",")

			for _, upStr := range userPassSlices {
				up := strings.Split(upStr, ":")
				if len(up) != 2 {
					log.Fatal(errors.Errorf("invalid user:pass format: %s", upStr))
				}

				user := model.User{
					UserName: up[0],
					Password: up[1],
				}

				userAuth := model.UserAuths{
					User: user,
					PathRoles: map[string]model.PathRole{
						user.UserName: {Path: "/", Mode: "w"},
					},
				}

				serviceConfig.Users = append(serviceConfig.Users, userAuth)
			}
		}

		err = dfvalidator.Struct(serviceConfig)
		if err != nil {
			log.Fatal(errors.Wrap(err, "validate config error"))
		}

		// fmt.Println("server config :")
		// spew.Dump(serviceConfig)

		// fmt.Println("viper all settings :")
		// spew.Dump(viper.AllSettings())

		anny := model.Annymous{
			Enable: serviceConfig.Anny.Enable,
		}
		anny.Mode = serviceConfig.Anny.Mode
		anny.Path = serviceConfig.Anny.Path

		server.Serve(serviceConfig.Server, serviceConfig.Users, anny)
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	// Define a flag for the configuration file
	serverCmd.PersistentFlags().StringP("config", "c", "", "Config file path")

	serverCmd.PersistentFlags().StringP("server.host", "H", "0.0.0.0", "Server host")
	serverCmd.PersistentFlags().IntP("server.port", "p", 8080, "Server port")

	serverCmd.PersistentFlags().Bool("server.https", false, "Server use https")

	serverCmd.PersistentFlags().String("server.cert", "", "server https cert file")
	serverCmd.PersistentFlags().String("server.key", "", "server https key file")

	serverCmd.PersistentFlags().String("server.home", "", "server home dir")

	// 匿名访问
	serverCmd.PersistentFlags().Bool("anny.enable", false, "enable annymous user")
	serverCmd.PersistentFlags().String("anny.mode", "r", "anny user mode, r: read, w: write")
	serverCmd.PersistentFlags().String("anny.path", "/", "anny user path")

	// 用户认证
	serverCmd.PersistentFlags().String("user", "admin:easyadmin123", "username:password,username:password")

}

// randomString 生成随机字符串
func randomString(n int) string {
	var letters = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	s := make([]byte, n)
	for i := range s {
		s[i] = letters[0]
	}
	return string(s)
}
