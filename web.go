package main
import (
	"log"
	"net"
	"fmt"
	"strings"
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
	for {
		// TODO 大消息处理
		n, err := conn.Read(buf)
		data := string(buf[:n])
		if data != "" {
			firstLineData := parseProtocol(data[0:strings.Index(data, "\r\n")])
			fmt.Println("协议:" + firstLineData.protocol)
			fmt.Println("方法:" + firstLineData.method)
			indexFirst := strings.Index(data, "\r\n")
			indexSecond := strings.Index(data, "\r\n\r\n")
			headLineData := parseHeader(data[indexFirst+2:indexSecond])
			if headLineData != nil {
				request := new(httpRequest)
				(*request).protocol = *firstLineData
				request.headMap = headLineData
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

func handlerRequest(request httpRequest) (string){
	fmt.Println("Handler httpRequest")
	switch request.protocol.method {
		case "GET":
			parseURL(request)
			break;
		case "POST":
			parseBody()
			break;
		default:
			break;
	}
	return "HTTP/1.1 200 OK\r\nConnection: Keep-Alive\r\nLast-Modified : Fri , 12 May 2020 18:53:33 GMT\n\rContent-Length:14\n\rContent-Type:application/json\r\n\r\n{\"test\":\"123\"}"
}

func parseURL(request httpRequest) {
	url := request.protocol.url;
	index := strings.Index(url, "?")
	if index == -1 {
		return
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
}

func parseBody() {
	fmt.Println("POST METHOD PARSE REQUEST BODY")
}