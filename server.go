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
  for ;; {
    conn, err := listener.Accept()
    CheckError(err)
    //TimeOut処理
    conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
    go doServer(conn)
  }
}

func doServer(conn net.Conn) {
  // TODO キー入力で終了させたい。
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
      //TODO ここでPOST処理を作成。まずはHTTP Bodyをパースしていく
      messageBodyFlag = repPartition.MatchString(string(line))
      //TODO POST, GETを大きく分離したい
      if(messageBodyFlag) {
        //bodyの読み取り処理
        var readBodyLength = 0;
        var bodyString = ""
        var i = 0
        for ; i < contentLength; i++ {
          line, _, err := reader.ReadLine()
          readBodyLength += len(string(line))
          bodyString += string(line)
          if (readBodyLength >= contentLength || err == io.EOF) {
            println(readBodyLength)
            println(contentLength)
            break
          }
        }
        println(bodyString)
        break
      }
    }
  }
  //TODO 内部ディレクトリ指定 && ファイル名の指定がない場合に、index.htmlが見れない
  //TODO ファイルがない場合にフォルダ内のListを生成したい。
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

// リクエストで返すための読み込みファイルの内容を返す
func readFileContent(fileName string) string {
  //TODO どでかいファイル入ると多分落ちる。
  fp, err := ioutil.ReadFile("./resources" + fileName)
  if err != nil {
    return "Not Found!"
  }
  body := fp
  return string(body)
}

// 対象文字列をカウントをする
func CountByteLength(target string) int {
  //TODO UTF8しかできないのはまずい
  return utf8.RuneCountInString(target)
}

//汎用のエラー処理
func CheckError(err error) {
  if err != nil {
    fmt.Printf("error: %s\n", err)
    return
  }
}

// URLのパスから対象ファイルの拡張子を取得する
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
  //TODO Keep-Alive
  //TODO contentLengthがUTF8のみで行われているため、適切な長さを返せない場合がある。
  contentLength  := "Content-Length: " + strconv.Itoa(CountByteLength(messageBody) + 1) + "\n"
  return responseStatus + contentType + serverName + contentLength + "\n"
}
