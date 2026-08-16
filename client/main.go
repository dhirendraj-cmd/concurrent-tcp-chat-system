package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)


func main(){
	conn, err := net.Dial("tcp", "localhost:1993")

	fmt.Println("Waiting to establish connection checking error ... ")

	if err!=nil{
		fmt.Println("Error while connection to server.... ", err)
		return
	}

	defer conn.Close()
	fmt.Println("Conneted on: ", conn.RemoteAddr())

	go func ()  {
		serverInput := bufio.NewReader(conn)
		for {
			// read from server
			fmt.Print("> ")
			serverMsg, err := serverInput.ReadString('\n')
			if err != nil{
				fmt.Println("Error while reading msg from server: ", err)
				return
			}
		
			fmt.Println(strings.TrimSpace(serverMsg))
		
		}
	}()

	// start writing to server
	stdInput := bufio.NewReader(os.Stdin)
	for {
		
		// try to read one line at a time
		writeToServer, err := stdInput.ReadString('\n')
		if err!=nil{
			fmt.Println("Error while reading lines from terminal....", err)
			return
		}

		// send data to server by writing
		_, err = conn.Write([]byte(strings.TrimSpace(writeToServer)+"\n"))
		if err != nil{
			fmt.Println("Error while writing to server... ", err)
			return
		}
	}



}

