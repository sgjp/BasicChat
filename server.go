package main


import (
	"net"
	"log"
	"bufio"
	"strings"
	"fmt"
)

var chatRooms map[string]Chatroom
//var clients []*Client
//var chatRoomId int = 1

func main() {

	ListenAndServe("9999")
}



func ListenAndServe(port string){
	listener, _ := net.Listen("tcp",":" + port)
	chatRooms = make(map[string]Chatroom)

	for{
		conn, err := listener.Accept()
		if err!=nil {
			log.Fatal("Error starting server %v",err)
		}
		client := Client{conn,""}
//		clients = append(clients, &client);

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
			//removeClient(client)
			return
		}


		command, message := parseMessage(Decode(string(line)))
		log.Printf("command: %v, message: %v",command,message)

		switch command {

		//Create chatroom
		case "C":
			created := CreateNewChatRoom(message)
			if created{
				fmt.Fprintln(client.Connection,"The chatroom was created")
			}else{
				fmt.Fprintln(client.Connection,"The chatroom already exists")
			}

		//List chatrooms
		case "L":
			if len(chatRooms)==0{
				fmt.Fprintln(client.Connection,"There are no chatrooms created")
			}else{
				for k := range chatRooms {
					fmt.Fprint(client.Connection,"* "+k)
				}
			}


		//Join existing chatroom
		case "J":
			joined := joinChatroom(client,message)
			if joined{
				fmt.Fprintln(client.Connection,"You joined the chatroom")
			}else{
				fmt.Fprintln(client.Connection,"The chatroom doesn't exists")
			}

		//Message to chatrooms
		case "M":
			out <- string(client.UserName+": "+message)

		//Quit chatroom
		case "Q":
			left := leaveChatRoom(client,message)
			if left{
				fmt.Fprintln(client.Connection,"You left the chatroom")
			}else{
				fmt.Fprintln(client.Connection,"You were not in the chatroom or it doesn't exist")
			}
		//Set username
		case "U":
			client.UserName=message
			client.UserName = strings.Replace(client.UserName,"\n","",-1)
			//fmt.Fprintf(client.Connection,"Welcome " + client.UserName + "!!" )

		//Send message
		default:
			out <- string(client.UserName+": "+message)


		}
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
	//Go through all the chatrooms
	for k := range chatRooms {
		//Go through the clients
		for i:= range chatRooms[k].clients{
			//Is the client in this chatroom?
			if (chatRooms[k].clients[i] == *client) {
				addMessageToChatroom(k,message)
				//Go through all the clients in this chatroom
				for j:= range chatRooms[k].clients{
					//Avoid sending back the message to the same user
					if (chatRooms[k].clients[j] != *client) {
						fmt.Fprintln(chatRooms[k].clients[j].Connection, message)
					}
				}

			}
		}


	}
}

func CreateNewChatRoom(name string) bool{

	if _, ok := chatRooms[name]; ok {
		return false
	}
	log.Printf("CHs %v",chatRooms)
	var clients []Client
	var messages []string
	chatRoom := Chatroom{name,clients,messages}

	chatRooms[name] = chatRoom
	return true

}
func joinChatroom(client *Client, name string)  bool{

	for k := range chatRooms {
		if(k==name){
			chatRoom := chatRooms[k]
			chatRoom.clients = append(chatRoom.clients, *client)
			chatRooms[k] = chatRoom
			return true
		}
	}
	return false
}

func leaveChatRoom(client *Client, name string)  bool{

	//Go through all the chatRooms
	for k := range chatRooms {
		//Find the specified chatroom
		if(k==name){
			//Go througn all the clients for the chatroom
			for i:= range chatRooms[k].clients{
				//Is the user in this chatroom?
				if (chatRooms[k].clients[i] == *client) {
					chatRoom := chatRooms[k]
					clients := chatRoom.clients
					clients = append(clients[:i],clients[i+1:]...)
					chatRoom.clients = clients
					chatRooms[k] = chatRoom
					return true
				}
			}
		}
	}
	return false
}

func addMessageToChatroom(name, message string){
	chatRoom := chatRooms[name]
	chatRoom.meesages = append(chatRoom.meesages,message)
}

func parseMessage(data string) (string, string){
	result := strings.SplitN(data, "/;",2)
	if len(result)<2{
		return "",result[0]
	}
	return result[0],result[1]

}

type Client struct {
	Connection net.Conn
	UserName string
}

type Chatroom struct{
	name string
	clients []Client
	meesages []string
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