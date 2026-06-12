package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)


func main(){
	lsn, err := net.Listen("tcp", "localhost:1993")

	if err!=nil{
		fmt.Println("Error while listening ", err)
		return
	}

	fmt.Println("Waiting for client to connect...")

	for {
		conn, err := lsn.Accept()

		if err!=nil {
			fmt.Println("Error while establishing connection... ", err)
			return
		}

		fmt.Println("connected....")

		go dealWithClient(conn)
	}
}


func dealWithClient(conn net.Conn){
	for {

		// read from client
		clientInput := bufio.NewReader(conn)
		clientMsg, err := clientInput.ReadString('\n')

		if err!=nil{
			fmt.Println("Error while reading input from client.... ", err)
			return
		}
		fmt.Println("Received Msg from Client: ", strings.TrimSpace(clientMsg))

		// respond from server
		fmt.Print("Server >: ")
		stdInput := bufio.NewReader(os.Stdin)
		writeToClient, err := stdInput.ReadString('\n')
		if err!=nil{
			fmt.Println("Error while writing to client.... ", err)
			return
		}
		_, err = conn.Write([]byte(strings.TrimSpace(writeToClient)+"\n"))
		if err != nil{
			fmt.Println("Error while writing to client... ", err)
			return
		}
	}
}
