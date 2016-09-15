package main


import (
	"net"
	"log"
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp","localhost:9999")

	if err!=nil {
		log.Fatalf("Conection refused %v",err)
	}
	defer conn.Close()

	go ServerMessageHanlder(conn)

	for true {
		InputHandler(conn)
	}
}

func ServerMessageHanlder(conn net.Conn){
	reader := bufio.NewReader(conn)
	for true {
		message, err := reader.ReadString('\n')
		if err!=nil {
			log.Fatal("Conection lost %v",err)
		}
		message = Decode(message)
		fmt.Printf("MESSAGE RECEIVED: \n %v", message)
	}
}

func InputHandler(conn net.Conn){
	reader := bufio.NewReader(os.Stdin)

	for true {
		m, _ := reader.ReadString('\n')
		message := Encode(m)
		conn.Write([]byte(message))


	}
}

func Decode(value string) (string) {
	var ENCODING_UNENCODED_TOKENS = []string{"%", ":", "[", "]", ",", "\""}
	var ENCODING_ENCODED_TOKENS = []string{"%25", "%3A", "%5B", "%5D", "%2C", "%22"}
	return replace(ENCODING_ENCODED_TOKENS,ENCODING_UNENCODED_TOKENS, value)
}

func Encode(value string) (string) {
	var ENCODING_UNENCODED_TOKENS = []string{"%", ":", "[", "]", ",", "\""}
	var ENCODING_ENCODED_TOKENS = []string{"%25", "%3A", "%5B", "%5D", "%2C", "%22"}
	return replace(ENCODING_UNENCODED_TOKENS, ENCODING_ENCODED_TOKENS, value)
}

func replace(fromTokens []string, toTokens []string, value string) (string) {
	for i:=0; i<len(fromTokens); i++ {
		value = strings.Replace(value, fromTokens[i], toTokens[i], -1)
	}
	return value;
}