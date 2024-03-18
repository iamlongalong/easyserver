/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/iamlongalong/easyserver/cmd/internal/model"
	"github.com/iamlongalong/easyserver/cmd/internal/server"
	"github.com/iamlongalong/easyserver/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	shareFileTimes    *int
	shareFileDuration *time.Duration
)

// shareCmd represents the share command
// easyserver share .
// easyserver share xxx
var shareCmd = &cobra.Command{
	Use:   "share",
	Short: "share a file or dir with default config and timeout",
	Long:  `share a file or dir with default config and timeout`,
	Run: func(cmd *cobra.Command, args []string) {
		dirOrFile := "."

		if len(args) > 0 {
			dirOrFile = args[0]
		}

		// 检查 dirOrFile
		info, err := os.Stat(dirOrFile)
		if err != nil {
			log.Fatalf("share fail with dir [%s] info: %s", dirOrFile, err)
			return
		}

		absPath, err := filepath.Abs(dirOrFile)
		if err != nil {
			log.Fatalf("share fail with dir [%s] abs fail: %s", dirOrFile, err)
			return
		}

		port, err := utils.GetAvailablePort(5000, 10000)
		if err != nil {
			log.Fatalf("share fail with dir [%s], get port fail : %s", dirOrFile, err)
			return
		}

		ipStr := "127.0.0.1"
		ip, err := utils.GetPreferredOutboundIP()
		if err != nil {
			log.Printf("get prefer outbound ip fail : %s", err)
		} else {
			ipStr = ip.String()
		}

		serviceConfig := model.ServieConfig{
			Server: model.Server{
				Addr: fmt.Sprintf("0.0.0.0:%d", port),
				Home: "",
			},
			Any: model.Annymous{
				Enable: true,
				Paths:  []model.PathRole{},
			},
			CloseConf: model.CloseConf{
				MaxDuration: *shareFileDuration,
			},
		}

		if info.IsDir() { // dir 开放目录权限
			serviceConfig.Home = absPath
			serviceConfig.Any.Paths = append(serviceConfig.Any.Paths, model.PathRole{
				Path: "/",
				Mode: "r",
			})

			fmt.Printf("\n👁 share link: %s\n\n", utils.GetHttpAddrString(false, ipStr, port, "/_dash"))
		} else {
			serviceConfig.CloseConf.MaxTimes = *shareFileTimes // serve 次数

			serviceConfig.Home = filepath.Dir(absPath)

			serviceConfig.Any.Paths = append(serviceConfig.Any.Paths, model.PathRole{
				Path: filepath.Join("/", info.Name()),
				Mode: "r",
			})
			fmt.Printf("\n👁 share link: %s\n\n", utils.GetHttpAddrString(false, ipStr, port, info.Name()))
		}

		server.Serve(serviceConfig)
	},
}

func init() {
	rootCmd.AddCommand(shareCmd)

	shareFileTimes = shareCmd.Flags().Int("times", 5, "times for file share (only for share a file). 0 to disable auto close")
	shareFileDuration = shareCmd.Flags().Duration("duration", time.Minute*10, "time for share, close after the duration. 0 to disable auto close")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// shareCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// shareCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
