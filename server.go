package main


import (
	"net"
	"log"
	"bufio"
	"strings"
	"fmt"
)

var chatRooms map[string]int
var clients []*Client
var chatRoomId int = 1

func main() {

	ListenAndServe("9999")
}



func ListenAndServe(port string){
	listener, _ := net.Listen("tcp",":" + port)
	chatRooms = make(map[string]int)

	for{
		conn, err := listener.Accept()
		if err!=nil {
			log.Fatal("Error starting server %v",err)
		}
		client := Client{conn,-1,""}
		clients = append(clients, &client);

		fmt.Fprintln(client.Connection, "Welcome, please set your username,chatroom in that format:")
		channel := make(chan string)

		go handleInMessages(channel, &client)
		go handleOutMessages(channel, &client)


	}


}




func handleInMessages(out chan string, client *Client) {
	defer close(out)


	reader := bufio.NewReader(client.Connection)
	for {
		line, err := reader.ReadBytes('\n')
		log.Printf("Message %v, from %v",string(line),client)
		if err != nil {
			removeClient(client)
			return
		}
		if client.ChatRoom<=0|| client.UserName==""{

			data := strings.Split(Decode(string(line)),",")
			if len(data)>0{
				GenerateNewChatRoom(strings.Replace(data[1],"\n","",-1))
				client.UserName = data[0]
				client.ChatRoom = chatRooms[strings.Replace(data[1],"\n","",-1)]
				fmt.Fprintln(client.Connection,"Dear: " +client.UserName + ", you joined the chatroom: " + strings.Replace(data[1],"\n","",-1))
				continue
			}

		}
		out <- string(line)
	}
}

func handleOutMessages(in <-chan string, client *Client) {

	for {
		message := <- in
		if (message != "") {
			message = strings.TrimSpace(message)
			BroadcastMessage(message, client)
		}
	}
}


func BroadcastMessage(message string, client *Client){
	log.Printf("Broadcasting Message %v",message)
	for _, _client := range clients {
		if (client.ChatRoom == _client.ChatRoom) {
			fmt.Fprintln(_client.Connection, message)
		}
	}
}

func GenerateNewChatRoom(name string){
	if chatRooms[name] <= 0 {
		chatRooms[name] = chatRoomId + 1
	}

}

func removeClient(client *Client){
	var clientsTemp []*Client
	var index = -1
	client.Connection.Close()
	for i, cl := range clients{
		if(cl == client){
			index = i
			break
		}
	}

	if (index>=0){
		clientsTemp := make([]*Client, len(clients)-1)
		copy(clientsTemp,clients[:index])
		copy(clientsTemp[index:],clients[index+1:])
	}
	clients = clientsTemp

}
type Client struct {
	Connection net.Conn
	ChatRoom int
	UserName string
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