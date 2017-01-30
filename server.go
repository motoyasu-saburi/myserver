package main

import (
  "strings"
  "fmt"
  "net"
  "unicode/utf8"
  "strconv"
  "bufio"
  "io"
  "time"
  "io/ioutil"
  "regexp"
)

func main() {
  println("start go http server!")

  listener, err := net.Listen("tcp", "localhost:8080")
  if err != nil {
    fmt.Printf("Listen error: %s\n", err)
    return
  }
  defer listener.Close()
  //パッケージ内変数のようにする案
  // var inputKeyFlag bool = false
  /**
   * TODO １リクエストあたりの処理を別スレッドに分けないとリクエストをさばけない
   * （複数リクエストが混じるリクエスト） ex: html内にあるimageの読み込み
   */
  for ;; {
    conn, err := listener.Accept()
    CheckError(err)
    conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
    go doServer(conn)
  }
}

func doServer(conn net.Conn) {
  // TODO 返り値が戻ってこないのでは？調査
  // inputKeyFlag = go inputKey()
  // if(inputKeyFlag) {
  //   break
  // }

  reader := bufio.NewReader(conn)
  status, err := reader.ReadString('\n')
  CheckError(err)
  //0: method, 1: パス, 2: httpのバージョン
  splitedStatus := strings.Split(status, " ")
  if(splitedStatus[0] == "POST") {
    //TODO ここで適切なContent-Lengthを取得し、ボディ部でその分だけ読み込み処理
    repContentLength := regexp.MustCompile(`^Content-Length`)
    repPartition := regexp.MustCompile(`^$`)
    messageBodyFlag := false
    // for line := ""; err == nil; line, err = reader.ReadString('\n') {
    var contentLength = 0
    for {
      line, _, err := reader.ReadLine()
      if err == io.EOF {
        break
      }
      // ここ、MatchStringしてないとcontentLengthが無いと怒られるのでどうにかせねば
      if(repContentLength.MatchString(string(line))) {
        repLength := regexp.MustCompile(`\b[0-9]+\b`)
        contentLength, err = strconv.Atoi(repLength.FindString(string(line)))
        CheckError(err)
      }
      //HTTPのボディ部の処理
      messageBodyFlag = repPartition.MatchString(string(line))
      if(messageBodyFlag) {
        //bodyの読み取り処理
        var readBodyLength = 0;
        var bodyString = ""
        var i = 0
        println("First ---------------------")
        for ; i < contentLength; i++ {
          println("@@@@@@@@@@@@@@@@@@@@@@@@@")
          line, _, err := reader.ReadLine()
          println("@@@@@@@@@@@@@@@@@@@@@@@@@")
          readBodyLength += len(string(line))
          bodyString += string(line)
          if (readBodyLength >= contentLength || err == io.EOF) {
            println("readBody is")
            println(readBodyLength)
            println("contentLength is")
            println(contentLength)
            println("break 1")
            break
          }
        }
        println(bodyString)
        println("break 2")
        break
      }
    }
  }
  path := splitedStatus[1]
  if(path == "/") {
    path = "/index.html"
  }
  //TODO Body生成部分に切り分けたい
  messageBody := readFileContent(path)
  extension := getExtension(path)
  res := GenerateHttpHeader(messageBody, extension)
  res += messageBody + "\n"
  conn.Write([]byte(res))
  defer conn.Close()
}

func readFileContent(fileName string) string {
  //TODO どでかいファイル入ると多分落ちる。
  fp, err := ioutil.ReadFile("./resources" + fileName)
  if err != nil {
    return "Not Found!"
  }
  body := fp
  return string(body)
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

func getExtension(urlPath string) string {
  extensionPosition := strings.LastIndex(urlPath, ".")
  extension := "html"
  if(extensionPosition > 0) {
    extension = urlPath[extensionPosition +1 :] // ex: "jpg", "png", "html"
  }
  return extension
}

func PostParameterParser(param string) {

}

func JsonParser() {

}

func SelectContentType(extension string) string {
  //TODO 将来的には画像に関しては先端のバイトコードを見て適切なContentTypeを送信したい。
  switch extension {
  case "html", "HTML": return "text/html; charset=utf-8;\n"
  case "png", "PNG": return "image/png;\n"
  case "jpeg", "JPEG", "jpg", "JPG": return "image/jpeg;\n"
  // case "txt", "TXT", "text", "TEXT": return "plain/text;\n"
  default: return ";\n"
  }
}

func GenerateHttpHeader(messageBody string, fileExtension string) string {
  //TODO 適切なResponseのStatsを返せるようにしたい。
  responseStatus := "HTTP/1.1 200 OK\n"
  contentType    := "Content-Type: " + SelectContentType(fileExtension)
  serverName     := "Server: goserver\n"
  keepAlive := "Connection: Keep-Alive\n"
  //TODO contentLengthがUTF8のみで行われているため、適切な長さを返せない場合がある。
  contentLength  := "Content-Length: " + strconv.Itoa(CountByteLength(messageBody) + 1) + "\n"
  return responseStatus + contentType + keepAlive + serverName + contentLength + "\n"
}
