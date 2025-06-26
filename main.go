package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

func dotToAscii(dot string, fancy bool) (string, error) {
	dotURL := "https://dot-to-ascii.ggerganov.com/dot-to-ascii.php"
	boxart := 0
	if fancy {
		boxart = 1
	}

	params := url.Values{
		"boxart": {fmt.Sprintf("%d", boxart)},
		"src":    {dot},
	}

	// 拼接 GET 请求参数
	fullURL := dotURL + "?" + params.Encode()
	resp, err := http.Get(fullURL)
	if err != nil {
		return "", fmt.Errorf("发送 GET 请求时出错: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应内容时出错: %v", err)
	}

	if len(respBody) == 0 {
		return "", fmt.Errorf("DOT 字符串格式不正确")
	}

	return string(respBody), nil
}

func main() {
	var (
		fancy = flag.Bool("fancy", false, "是否使用 boxart 字符")
		show  = flag.Bool("show", false, "是否打印读取的文件内容")
	)
	flag.Usage = func() {
		fmt.Println("用法: dot <dot文件路径> [--fancy] [--show]")
		fmt.Println("  <dot文件路径> 必填，dot文件路径")
		fmt.Println("  --fancy       可选，是否使用 boxart 字符（默认 false）")
		fmt.Println("  --show        可选，是否打印读取的文件内容（默认 false）")
	}
	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
		return
	}

	filePath := flag.Arg(0)

	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("读取文件时出错: %v\n", err)
		return
	}

	contentString := string(fileContent)
	if *show {
		fmt.Println("读取的文件内容:")
		fmt.Println(contentString)
	}

	asciiContent, err := dotToAscii(contentString, *fancy)
	if err != nil {
		fmt.Printf("转换时出错: %v\n", err)
		return
	}

	if *show {
		fmt.Println("转换后的 ASCII 内容:")
	}
	fmt.Println(asciiContent)
}
