package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
)

var host = flag.String("host", "0.0.0.0", "host")
var port = flag.String("port", "8080", "port")
var headers = map[string]string{
	"Accept":          "能接受的返回内容类型",
	"Accept-Charset":  "能接受的字符集",
	"Accept-Encoding": "能接受的编码方式列表",
	"Accept-Language": "能够接受的回应内容的自然语言列表",
	"Accept-Datetime": "能够接受的按照时间来表示的版本",
	"Authorization":   "认证信息",
	"Cache-Control":   "缓存机制",
	"Connection":      "客户端想要优先使用的连接类型",
	"Cookie":          "客户端传来的 Cookie",
	"Content-Length":  "请求体长度",
	"Content-MD5":     "请求体内容 MD5",
	"Content-Type":    "请求体类型",
	"Date":            "发送该消息的日期和时间",
	"Expect":          "表明客户端要求服务器做出特定的行为",
	"From":            "发起此请求的用户的邮件地址",
	"Host":            "服务器的域名和端口号",
	"Origin":          "跨域来源",
	"Referer":         "Referer 页面地址",
	"TE":              "浏览器预期接受的传输编码方式",
	"User-Agent":      "浏览器的浏览器身份标识字符串",
	"Upgrade":         "要求服务器升级到另一个协议",
	"Via":             "请求所使用的代理",
	"DNT": "请求不要跟踪",
	"Upgrade-Insecure-Requests": "请求使用加密链接",
}

func main() {
	flag.Parse()
	l, err := net.Listen("tcp", *host+":"+*port)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()
	fmt.Println("Listening on " + *host + ":" + *port)
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		// Handle connections in a new goroutine.
		go handleRequest(conn)
	}
}

func printHeaders(reader *bufio.Reader) {
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading:", err.Error())
			return
		}
		if len(line) == 2 { // \r\n
			return
		}
		fmt.Println(line)
		pair := strings.SplitAfterN(line, ":", 2)
		pair[0] = strings.Trim(pair[0], ":")
		name, ok := headers[pair[0]]
		if !ok {
			fmt.Println("Unknown header:", pair[0])
		} else {
			fmt.Println(name, ": ", pair[1], "\n")
		}
	}
}

func serveFile(writer *bufio.Writer, path string) {

}

func serveFolder(writer *bufio.Writer, path string) {

}

func handleRequest(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	firstLine, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	firstLine = strings.Trim(firstLine, "\r\n")
	entries := strings.Split(firstLine, " ")
	if len(entries) != 3 {
		fmt.Println("Invalid request.")
		return
	}
	fmt.Println(firstLine)
	method := entries[0]
	path := entries[1]
	fmt.Println("HTTP 请求类型：", method, " 路径：", path, "版本：", entries[2], "\n")
	go printHeaders(reader)

	writer := bufio.NewWriter(conn)
	if strings.HasSuffix(path, "/") {
		go serveFolder(writer, path)
	} else {
		go serveFile(writer, path)
	}
	conn.Write([]byte("Message received."))
}
