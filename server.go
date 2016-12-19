package main

import (
  "strings"
  "fmt"
  "net"
  "unicode/utf8"
  "strconv"
  "bufio"
  "os"
  // "runtime"
)

func main() {
  listener, err := net.Listen("tcp", "localhost:8080")
  if err != nil {
    fmt.Printf("Listen error: %s\n", err)
    return
  }
  defer listener.Close()
  //パッケージ内変数のようにする案
  // var inputKeyFlag bool = false
  for ;; {
    //TODO 返り値が戻ってこないのでは？調査
    // inputKeyFlag = go inputKey()
    // if(inputKeyFlag) {
    //   break
    // }

    conn, err := listener.Accept()
    CheckError(err)
    defer conn.Close()

    status, err := bufio.NewReader(conn).ReadString('\n')
    CheckError(err)

    //1: method, 2: パス, 3: httpのバージョン
    splitedStatus := strings.Split(status, " ")
    path := splitedStatus[1]
    if(path == "/") {
      path = "/index.html"
    }
    messageBody := readFileContent(path)

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
        fmt.Printf("error: %s\n", err)
        return
      }
  		fmt.Print(string(buf[:n]))
  	}
  }
}

func readFileContent(fileName string) string {
  fp, err := os.Open("./resources" + fileName)
  if err != nil {
    return "Not Found!"
  }
  scanner := bufio.NewScanner(fp)
  body := ""
  for scanner.Scan() {
    body += scanner.Text()
  }
  return body
}

// func inputKey() bool {
//   runtime.Gosched()
//   var key string
//   fmt.Scan(&key)
//   return (key == "q")
// }

func CountByteLength(target string) int {
  return utf8.RuneCountInString(target)
}

func CheckError(err error) {
  if err != nil {
    fmt.Printf("error: %s\n", err)
    return
  }
}

func GenerateHttpHeader(messageBody string) string {
  responseStatus := "HTTP/1.1 200 OK\n"
  contentType    := "Content-Type: text/html; charset=utf-8;"
  serverName     := "Server: goserver\n"
  contentLength  := "Content-Length: " + strconv.Itoa(CountByteLength(messageBody) + 1) + "\n"

  return responseStatus + contentType + serverName + contentLength + "\n"
}
