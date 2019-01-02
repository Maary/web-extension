package main

import (
	"bufio"
	"crawler.center/cookie.center/web-extesion/installer/config"
	"github.com/fatih/color"
	"log"
	"os"
)

func main() {
	c := new(config.ConfigInfo)
	c.Name("server")
	ok, err := c.CreateConfig()
	if err != nil {
		log.Fatal(err)
	}
	yellow := color.New(color.FgYellow).PrintFunc()
	red := color.New(color.FgRed).PrintfFunc()
	blue := color.New(color.FgBlue).PrintfFunc()
	blue("输入回车结束...")
	reader := bufio.NewReader(os.Stdin)
	result, err := reader.ReadBytes('\n')
	terminal := result[len(result) - 1]
	if terminal == '\n' {
		if ok {
			yellow("安装完成")
		} else {
			red("安装失败")
		}
	}
}
