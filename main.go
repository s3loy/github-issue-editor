package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/s3loy/github-issue-editor/issue"
)

func printUsage() {
	fmt.Println("Usage: github-issue-editor <command> [arguments]")
	fmt.Println("\nCommands:")
	fmt.Println("  create <owner> <repo>             创建一个新的 Issue")
	fmt.Println("  read   <owner> <repo> <number>    读取指定 Issue")
	fmt.Println("  update <owner> <repo> <number>    更新指定 Issue")
	fmt.Println("  close  <owner> <repo> <number>    关闭指定 Issue")
	fmt.Println("  help                              显示帮助信息")
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "create":
		if len(os.Args) < 4 {
			log.Fatal("参数不足: create <owner> <repo>")
		}
		owner, repo := os.Args[2], os.Args[3]
		if err := issue.CreateIssue(owner, repo); err != nil {
			log.Fatalf("创建失败: %v", err)
		}
	case "read":
		if len(os.Args) < 5 {
			log.Fatal("参数不足: read <owner> <repo> <number>")
		}
		owner, repo := os.Args[2], os.Args[3]
		number, err := strconv.Atoi(os.Args[4])
		if err != nil {
			log.Fatal("Issue 编号必须是整数")
		}
		if err := issue.ReadIssue(owner, repo, number); err != nil {
			log.Fatalf("读取失败: %v", err)
		}

	case "update":
		if len(os.Args) < 5 {
			log.Fatal("参数不足: update <owner> <repo> <number>")
		}
		owner, repo := os.Args[2], os.Args[3]
		number, err := strconv.Atoi(os.Args[4])
		if err != nil {
			log.Fatal("Issue 编号必须是整数")
		}
		if err := issue.UpdateIssue(owner, repo, number); err != nil {
			log.Fatalf("更新失败: %v", err)
		}

	case "close":
		if len(os.Args) < 5 {
			log.Fatal("参数不足: close <owner> <repo> <number>")
		}
		owner, repo := os.Args[2], os.Args[3]
		number, err := strconv.Atoi(os.Args[4])
		if err != nil {
			log.Fatal("Issue 编号必须是整数")
		}
		if err := issue.CloseIssue(owner, repo, number); err != nil {
			log.Fatalf("关闭失败: %v", err)
		}

	case "help", "-h", "--help":
		printUsage()

	default:
		fmt.Printf("未知命令: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}
