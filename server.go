package main

import (
  "fmt"
  "net"
  "unicode/utf8"
  "strconv"
)

func main() {
  listener, err := net.Listen("tcp", "localhost:8080")
  if err != nil {
    fmt.Printf("Listen error: %s\n", err)
    return
  }
  defer listener.Close()

  conn, err := listener.Accept()
  if err != nil {
    fmt.Printf("Accept error: %s\n", err)
    return
  }
  defer conn.Close()

  var messageBody = "<h1>hoge!!!!!!!!!</h1>"
	res := GenerateHttpHeader(messageBody)
  res += messageBody + "\n"

	conn.Write([]byte(res))

	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if n == 0 {
			break
		}
		if err != nil {
			fmt.Printf("Read error: %s\n", err)
		}
		fmt.Print(string(buf[:n]))
	}
}

func CountByteLength(target string)(int) {
  return utf8.RuneCountInString(target)
}

func GenerateHttpHeader(messageBody string)(string) {
  var responseStatus = "HTTP/1.1 200 OK\n"
  var contentType    = "Content-Type: text/html\n"
  var charset        = "charset=utf-8\n";
  var serverName     = "Server: goserver\n"
  var contentLength  = "Content-Length: " + strconv.Itoa(CountByteLength(messageBody) + 1) + "\n"

  return responseStatus + contentType + charset + serverName + contentLength + "\n"
}
