package main
import (
	"log"
	"net"
	"fmt"
	"strings"
	"strconv"
)

func main() {
	listener, err := net.Listen("tcp", "localhost:80");
	if err != nil {
		log.Fatal(err);
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Print("Server Error:\n")
			log.Fatal(err);
		}
		fmt.Printf("访问客户端信息： con=%v 客户端ip=%v\n", conn, conn.RemoteAddr().String())
		go handlerConn(conn)
	}
}

func handlerConn(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 1024 )
	var finish = true
	var request *httpRequest
	var finishParseHeader = false
	for {
		// TODO 大消息处理
		n, err := conn.Read(buf)
		data := string(buf[:n])
		if data != "" {
			if finish {
				request = new(httpRequest)
				finishParseHeader = false
			}
			request.requestBody += data;
			if strings.Index(request.requestBody, "\r\n") == -1 {
				continue
			}
			if !finishParseHeader {
				firstLineData := parseProtocol(request.requestBody[0:strings.Index(data, "\r\n")])
				fmt.Println("协议:" + firstLineData.protocol)
				fmt.Println("方法:" + firstLineData.method)
				indexFirst := strings.Index(request.requestBody, "\r\n")
				indexSecond := strings.Index(request.requestBody, "\r\n\r\n")
				headLineData := parseHeader(request.requestBody[indexFirst+2:indexSecond])
				if headLineData != nil {
					request.protocol = *firstLineData
					request.headMap = headLineData
				}
			}
			finish = parseRequest(*request)
			if finish {
				response := handlerRequest(*request)
				conn.Write([]byte(response))
			}
		}
		if err != nil {
			fmt.Println("Connection Error")
			return
		}
	}
}

type httpRequest struct {
	protocol firstLine
	headMap map[string]string
	requestURL string
	queryMap map[string]string
	requestBody string
}

type firstLine struct {
	protocol string
	method string
	url string
}

func parseHeader(headLines string) (map[string]string){
	if strings.Index(headLines, ":") == -1 {
		return nil
	}
	mapData := make(map[string]string)
	arrays := strings.Split(headLines, "\r\n")
	for _, data := range arrays {
		header := strings.Split(data, ":")
		mapData[header[0]] = header[1]
	}
	return mapData
}


func parseProtocol(protocolLine string) (*firstLine){
	firstLine := new(firstLine)
	arrays := strings.Split(protocolLine, " ")
	firstLine.protocol = arrays[2]
	firstLine.url = arrays[1]
	firstLine.method = arrays[0]
	return firstLine
}

func parseRequest(request httpRequest) (bool){
	fmt.Println("Parse httpRequest")
	fmt.Println(request.requestBody)
	switch request.protocol.method {
		case "GET":
			return parseGet(request)
		case "POST":
			return parsePost(request)
		default:
			return true;
	}
}

func parseGet(request httpRequest) (bool) {
	url := request.protocol.url;
	index := strings.Index(url, "?")
	if index == -1 {
		return true
	}
	request.requestURL = url[0:index]
	queryParamString := url[index+1:]
	queryParamsArray := strings.Split(queryParamString, "&")
	queryMap := make(map[string]string)
	if len(queryParamsArray) > 0 {
		for _, queryPair := range queryParamsArray {
			equalIndex := strings.Index(queryPair, "=")
			if equalIndex != -1 {
				queryMap[queryPair[:equalIndex]] = queryPair[equalIndex+1:]
			}
		}
	}
	request.queryMap = queryMap
	fmt.Println("GET METHOD PARSE QUERY PARAMS")
	return true
}

func parsePost(request httpRequest) (bool){
	headMap := request.headMap
	index := strings.Index(request.requestBody, "\r\n\r\n")
	requestBody := request.requestBody[index+2:]
	length, _ := strconv.ParseInt(strings.Trim(headMap["Content-Length"], " "), 0, 32)
	fmt.Println("POST METHOD PARSE REQUEST BODY")
	/** 加上/r/n的长度*/
	if int(length) + 2 == len(requestBody) {
		return true
	}
	return false
}

func handlerRequest(request httpRequest) (string) {
	return "HTTP/1.1 200 OK\r\nConnection: Keep-Alive\r\nLast-Modified : Fri , 12 May 2020 18:53:33 GMT\r\nContent-Length:14\r\nContent-Type:application/json\r\n\r\n{\"test\":\"123\"}"
}