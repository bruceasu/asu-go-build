/*
  go run main.go -project <你的项目名> -output <输出目录>
 */
package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	// 解析命令行参数
	// go run main.go -project myapp -output bin
	projectName := flag.String("project", "", "项目名称")
	outputDir := flag.String("output", "bin", "输出目录")
	flag.Parse()
	
	// 如果命令行参数为空，则尝试从 go.mod 中获取 module 名称
	if *projectName == "" {
		// 如果依然没有获取到 module 名称，则使用当前目录名
		if *projectName == "" {
			wd, err := os.Getwd()
			if err != nil {
				fmt.Printf("获取当前目录失败: %v\n", err)
				os.Exit(1)
			}
			*projectName = filepath.Base(wd)
		}
	}

	// 创建输出目录
	if err := os.MkdirAll(*outputDir, 0755); err != nil {
		fmt.Printf("创建输出目录失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("编译本地执行文件...")

	// 获取本机 go env 的 GOOS / GOARCH
	localGOOS, err := runCmdAndGetStdout("go", "env", "GOOS")
	if err != nil {
		fmt.Printf("获取本机 GOOS 失败: %v\n", err)
		os.Exit(1)
	}
	localGOARCH, err := runCmdAndGetStdout("go", "env", "GOARCH")
	if err != nil {
		fmt.Printf("获取本机 GOARCH 失败: %v\n", err)
		os.Exit(1)
	}

	localOutputName := fmt.Sprintf("%s-%s-%s", *projectName, localGOOS, localGOARCH)
	if localGOOS == "windows" {
		localOutputName += ".exe"
	}
	if err := buildWithEnv(*outputDir, localOutputName, localGOOS, localGOARCH); err != nil {
		fmt.Printf("本地编译失败: %v\n", err)
		os.Exit(1)
	}

	// 指定要交叉编译的平台
	platforms := []string{
		"windows/amd64",
		"windows/386",
		"linux/amd64",
		// "darwin/amd64", // 如果需要，可以加回来
	}

	fmt.Println("开始交叉编译...")
	for _, p := range platforms {
		parts := strings.Split(p, "/")
		if len(parts) != 2 {
			fmt.Printf("无效的平台配置: %s\n", p)
			os.Exit(1)
		}
		goos := parts[0]
		goarch := parts[1]

		outputName := fmt.Sprintf("%s-%s-%s", *projectName, goos, goarch)
		if goos == "windows" {
			outputName += ".exe"
		}

		if err := buildWithEnv(*outputDir, outputName, goos, goarch); err != nil {
			fmt.Printf("交叉编译失败: %s, error: %v\n", p, err)
			os.Exit(1)
		}
	}

	fmt.Println("编译完成.")
}

// buildWithEnv 调用 go build，设置 GOOS/GOARCH 并输出到指定文件名
func buildWithEnv(outputDir, outputName, goos, goarch string) error {
	cmd := exec.Command("go", "build", "-o", filepath.Join(outputDir, outputName))
	cmd.Env = append(os.Environ(),
		"GOOS="+goos,
		"GOARCH="+goarch,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// runCmdAndGetStdout 用于获取命令执行的标准输出（去掉末尾换行）
func runCmdAndGetStdout(name string, args ...string) (string, error) {
	var out bytes.Buffer
	cmd := exec.Command(name, args...)
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}
