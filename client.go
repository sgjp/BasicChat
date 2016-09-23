package main


import (
	"net"
	"log"
	"bufio"
	"fmt"
	"os"
	"strings"
)
var userName string

func main() {
	conn, err := net.Dial("tcp","localhost:9999")

	if err!=nil {
		log.Fatalf("Conection refused %v",err)
	}
	defer conn.Close()

	setUserName(conn)

	conn.Write([]byte("U/;"+userName))
	showMenu()


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
		fmt.Printf(" %v", message)
	}
}

func InputHandler(conn net.Conn){
	reader := bufio.NewReader(os.Stdin)

	for true {
		m, _ := reader.ReadString('\n')
		option, args := parseInput(m)
		if(option==""){
			fmt.Printf("Please select an option. Remember: 0 shows the menu")
		}else{
			option = strings.Replace(option,"\n","",-1)
			switch option {

			//Show the menu
			case "0":
				showMenu()


			//Create chatroom
			case "1":
				if(args==""){
					fmt.Printf("Not Args found. Example: '1 NewChatRoom'")
				}else{
					//message := Encode(args)
					conn.Write([]byte("C/;"+args))
				}

			//List chatroom
			case "2":
				//message := Encode(args)
				conn.Write([]byte("L/;\n"))



			//Join Existing chatroom
			case "3":
				if(args==""){
					fmt.Println("Not Args found. Example: '3 ExistingChatRoom'")
				}else{
					//message := Encode(args)
					conn.Write([]byte("J/;"+args))
				}

			//Send message
			case "4":
				if(args==""){
					fmt.Println("Not Args found. Example: '4 hello everyone!'")
				}else{
					//message := Encode(args)
					conn.Write([]byte("M/;"+args))
				}

			//Leave chatroom
			case "5":
				if(args==""){
					fmt.Println("Not Args found. Example: '5 ExistingChatRoom'")
				}else{
					//message := Encode(args)
					conn.Write([]byte("Q/;"+args))
				}

			default:
				//conn.Write([]byte("M/;"+args))

			}

		}
	}
}

func parseInput(m string)(string, string){
	splitted := strings.SplitN(m," ",2)
	if(len(splitted)>1){
		return splitted[0],splitted[1]
	}
	if(len(splitted)==1){
		return splitted[0],""
	}
	return "",""
}

func showMenu(){
	fmt.Println("")
	fmt.Println("--PLEASE SELECT THE DESIRED OPTION:\n")
	fmt.Println("  1. Create a chatroom.   Args: Name")
	fmt.Println("  2. List chatrooms.")
	fmt.Println("  3. Join existing chatroom.   Args: Name")
	fmt.Println("  4. Send Message to all joined chatrooms  Args: Message")
	fmt.Println("  5. Quit chatroom.    Args: Name")
	fmt.Println("  0. Show Menu")
	fmt.Println("")
	fmt.Println("  Example:  '3 chatroom2'")
	fmt.Println("")
	fmt.Println("")
	fmt.Println("No option sends a meesage to all joined chatrooms")
	fmt.Println("")

}

func setUserName(conn net.Conn){
	fmt.Println("Please set your username:")
	reader := bufio.NewReader(os.Stdin)
	m, _ := reader.ReadString('\n')
	userName = m

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