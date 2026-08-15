package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"sync"
	"strings"
)

// client data
type ConnectedClients struct{
	ID int
	Username string
	Conn net.Conn
}


var (
	activeClients = make(map[int]*ConnectedClients)
	clientsMu sync.Mutex
	idCounter = 1
)


func main(){

	lsn, err := net.Listen("tcp", "localhost:1993")
	if err!=nil{
		fmt.Println("Error while listening.... ", err)
		return
	}

	defer lsn.Close()
	
	for {
		conn, err := lsn.Accept()
		if err!=nil{
			fmt.Printf("Error %s occured while accepting connection....", err)
			return
		}
		go dealWithClient(conn)
	}

}


func dealWithClient(conn net.Conn){
	client, err := handleConnectedClients(conn)

	if err!=nil{
		fmt.Println("Error while taking username: ", err)
		return
	}
	
	defer func ()  {
		conn.Close()
		
		clientsMu.Lock()
		delete(activeClients, client.ID)
		clientsMu.Unlock()
		
		fmt.Printf("[-] %s has left chat\n", client.Username)
		
		}()
		
	readMsg := bufio.NewReader(conn)
	
	for {
		newline, err := readMsg.ReadString('\n')
		if err!=nil{
			if err!=io.EOF{
				fmt.Println("Error Reading msg!!", err)
			}
			return
		}

		fmt.Println("Message from client: ", strings.TrimSpace(newline)+"\n")
	}

}


func handleConnectedClients(conn net.Conn) (*ConnectedClients, error) {

	conn.Write([]byte("Welcome! Enter your username: \n"))
	
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err!=nil{
		conn.Close()
		return nil, err
	}

	username := strings.TrimSpace(string(buf[:n]))
	if username == ""{
		username = fmt.Sprintf("Guest%d\n: ", idCounter)
	}


	clientsMu.Lock()
	clientID := idCounter
	idCounter++

	clientData := &ConnectedClients{
		ID: clientID,
		Username: username,
		Conn: conn,
	}

	activeClients[clientID] = clientData
	clientsMu.Unlock()

	// for _, k := range activeClients{
	// 	fmt.Printf("ID:%v, Name:%v", k.ID, k.Username)
	// }

	fmt.Printf("[+] %s has entered chat\n", username)

	return clientData, nil
}



