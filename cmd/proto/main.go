// Package proto
// @file      : main.go
// @author    : china.gdxs@gmail.com
// @time      : 2023/11/7 17:01
// @Description: proto 生成文件 入口
package main

import (
	"github.com/china-xs/ginplus/cmd/proto/add"
	"github.com/china-xs/ginplus/cmd/proto/server"
	"github.com/spf13/cobra"
	"log"
)

const release = "v1.0.1"

func init() {
	rootCmd.AddCommand(add.CmdAdd)
	rootCmd.AddCommand(server.CmdServer)
	//rootCmd.AddCommand(project.CmdServer)

}

var rootCmd = &cobra.Command{
	Use:     "proto",
	Short:   "proto 快速生成proto文件",
	Long:    `proto 快速生成proto文件`,
	Version: release,
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
