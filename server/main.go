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
	clientsMu sync.RWMutex
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
	readMsg := bufio.NewReader(conn)
	client, err := handleConnectedClients(conn, readMsg)

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
		
	
	for {
		newline, err := readMsg.ReadString('\n')
		if err!=nil{
			if err!=io.EOF{
				fmt.Println("Error Reading msg!!", err)
			}
			return
		}

		// fmt.Println("Message from client: ", strings.TrimSpace(newline)+"\n")
		sentFromClients := client.Username + ": " + strings.TrimSpace(newline)
		broadCast(sentFromClients, client.ID)
	}

}


func handleConnectedClients(conn net.Conn, readMsg *bufio.Reader) (*ConnectedClients, error) {

	_, err := conn.Write([]byte("Welcome! Enter your username: \n"))

	if err != nil {
        conn.Close()
        return nil, err
    }
	
	username, err := readMsg.ReadString('\n')
	if err!=nil{
		conn.Close()
		return nil, err
	}

	username = strings.TrimSpace(username)
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

	fmt.Printf("[+] %s has entered chat\n", username)
	broadCast(fmt.Sprintf("[+] %s has entered chat\n", clientData.Username), clientData.ID)

	return clientData, nil
}


func broadCast(msg string, senderId int){
	clientsMu.RLock()
	clientSnapShot := make([]*ConnectedClients, 0, len(activeClients))
	for id, client := range activeClients{
		if id!=senderId{
			clientSnapShot = append(clientSnapShot, client)
		}
	}

	// fmt.Println("SnapShot>>> ", clientSnapShot, msg)
	clientsMu.RUnlock()

	for _, client := range clientSnapShot{
		_, err := client.Conn.Write([]byte(strings.TrimSpace(msg)+"\n"))
		if err!=nil{
			fmt.Printf("Failed to send broadcast to ID %d: %v\n", client.ID, err)
			continue
		}
	}

}





