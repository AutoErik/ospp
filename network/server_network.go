package network

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"net"
	"sync"
)

var playerServerMutex sync.Mutex
var bulletServerMutex sync.Mutex
var idCounterMutex sync.Mutex
var idMapMutex sync.RWMutex
var adressMapMutex sync.RWMutex
var bulletMapMutex sync.RWMutex

//ServerConn stores state for a connection in the server
type ServerConn struct {
	playerUpdates []PlayerUpdate
	bulletUpdates []BulletUpdate
	adressMap     map[int]string
	idMap         map[string]int
	bulletMap     map[int]bool
	idCounter     *intWrapper
}

type intWrapper struct {
	num int
}

// GetPlayerUpdates returns a slice of PlayerUpdates sent by the client. Can be used concurrently.
func (s *ServerConn) GetPlayerUpdates() []PlayerUpdate {
	playerServerMutex.Lock()

	var Updates []PlayerUpdate
	if len(s.playerUpdates) > 0 {
		Updates = s.playerUpdates
		s.playerUpdates = make([]PlayerUpdate, 0, 10)
	}
	playerServerMutex.Unlock()
	return Updates
}

// GetBulletUpdates returns a slice of BulletUpdates sent by the client. Can be used concurrently.
func (s *ServerConn) GetBulletUpdates() []BulletUpdate {
	bulletServerMutex.Lock()

	var Updates []BulletUpdate
	if len(s.bulletUpdates) > 0 {
		Updates = s.bulletUpdates
		s.bulletUpdates = make([]BulletUpdate, 0, 10)
	}
	bulletServerMutex.Unlock()
	return Updates
}

//SendWorldUpdate sends updates to all clients about the world state. Accepts PlayerUpdate and BulletUpdate structs.
func (s *ServerConn) SendWorldUpdate(worldUpdate interface{}) error {
	var OP uint8

	switch worldUpdate.(type) {
	case PlayerUpdate:
		OP = 2
	case BulletUpdate:
		OP = 3
	default:
		return fmt.Errorf("Error: Can only send PlayerUpdate or BulletUpdate structs so far, you sent:%v", worldUpdate)
	}

	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)

	if err := encoder.Encode(OP); err != nil {
		return fmt.Errorf("Error while encoding opcode")
	}
	if err := encoder.Encode(worldUpdate); err != nil {
		return fmt.Errorf("Error while encoding update struct")
	}

	adressMapMutex.RLock()
	for _, IPaddr := range s.adressMap { //sends to everyone on the map, at least for now

		addr, err := net.ResolveUDPAddr("udp", IPaddr+":8081")
		if err != nil {
			return fmt.Errorf("Error while resolving %v", IPaddr)
		}
		conn, err := net.DialUDP("udp", nil, addr)
		if err != nil {
			return fmt.Errorf("Error creating socket to %v", addr)
		}

		bufCopy := buf
		_, err = bufCopy.WriteTo(conn)
		if err != nil {
			return fmt.Errorf("Error while sending to %v", addr)
		}
		conn.Close()
	}
	adressMapMutex.RUnlock()

	return nil
}

// StartServerListener makes the server begin to listen for incoming requests from clients.
// Returns Error channel that has to be emptied regularly.
// Opcode 1: New client wants to be assigned a player ID.
// Opcode 2: New playerUpdates sent from clients.
// Opcode 3: New bulletUpdates sent from clients.
// Opcode 7: Test if the server is running, if so it sends a string to error.
// Opcode 8: Responds to ping requests from clients.
// Opcode 9: Shuts down the server.
func StartServerListener() (*ServerConn, <-chan error, error) {

	server := ServerConn{
		playerUpdates: make([]PlayerUpdate, 0, 10),
		bulletUpdates: make([]BulletUpdate, 0, 10),
		adressMap:     make(map[int]string),
		idMap:         make(map[string]int),
		bulletMap:     make(map[int]bool),
		idCounter:     new(intWrapper),
	}
	errChan := make(chan error, 10) //create function to read from this channel and save it in a log!

	addr, err := net.ResolveUDPAddr("udp", ":8080")
	if err != nil {
		return &server, nil, fmt.Errorf("Resolve adress error:%w", err)
	}
	udpConn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return &server, nil, fmt.Errorf("Socket creation error:%w", err)
	}

	//Main loop
	go func() {
		defer udpConn.Close()
		for {
			inputBytes := make([]byte, 200)

			n, addr, err := udpConn.ReadFromUDP(inputBytes)
			if err != nil {
				errChan <- fmt.Errorf("Reading %d bytes from %s stream error:%w", n, addr, err)
				continue
			}

			var Opcode uint8
			bufferOP := bytes.NewBuffer(inputBytes[:4])
			if err = gob.NewDecoder(bufferOP).Decode(&Opcode); err != nil {
				errChan <- fmt.Errorf("Decoding opcode [%v] error:%w", inputBytes[:4], err)
				continue
			}

			switch Opcode {
			case 1:
				go handleNewPlayer(&server, addr.IP.String(), errChan)
			case 2:
				go handleUpdate(&server, inputBytes[4:], addr.IP.String(), errChan)
			case 3:
				go handleNewBullet(&server, inputBytes[4:], addr.IP.String(), errChan)
			case 7:
				errChan <- errors.New("Server is running")
			case 8:
				go handlePingRequest(inputBytes[4:], addr.IP.String())
			case 9:
				return //Server shutdown
			default:
				errChan <- fmt.Errorf("Recieved faulty Opcode:%v", Opcode)
			}

			//fmt.Println("---------")
			//fmt.Println("Server has gotten a message size:", n)
			//fmt.Println("Opcode:", Opcode)
			//fmt.Println("Adress of sender:", addr)
			//fmt.Println("full buffer:", inputBytes[4:])
		}
	}()
	return &server, errChan, nil

}

func handleNewPlayer(server *ServerConn, addr string, errChan chan<- error) {

	idMapMutex.Lock()
	storedID := server.idMap[addr]
	var newID int

	if storedID == 0 {
		//idCounterMutex.Lock() behövs detta mutex?
		server.idCounter.num++
		newID = server.idCounter.num
		//idCounterMutex.Unlock()
		fmt.Println("new id was generated and sent:", newID)

		server.idMap[addr] = newID

	} else {
		newID = storedID
		fmt.Println("Previously generated id was sent out again:", newID)
	}
	idMapMutex.Unlock()

	fullAddr, err := net.ResolveUDPAddr("udp", addr+":8081")
	if err != nil {
		errChan <- fmt.Errorf("New client %s  addr resolve error:%w", addr, err)
		return
	}
	connUDP, err := net.DialUDP("udp", nil, fullAddr)
	if err != nil {
		errChan <- fmt.Errorf("New client  %s socket error:%w", addr, err)
		return
	}
	defer connUDP.Close()

	var OP uint8
	var buf bytes.Buffer
	OP = 1

	encoder := gob.NewEncoder(&buf)

	if err = encoder.Encode(OP); err != nil {
		errChan <- fmt.Errorf("New client %s opcode encode error:%w", addr, err)
		return
	}

	if err = encoder.Encode(newID); err != nil {
		errChan <- fmt.Errorf("New client %s ID encode error:%w", addr, err)
		return
	}

	if _, err = buf.WriteTo(connUDP); err != nil {
		errChan <- fmt.Errorf("New client %s WriteTo error:%w", addr, err)
		return
	}
}

func handleUpdate(server *ServerConn, newUpdate []byte, addr string, errChan chan<- error) {

	buffer := bytes.NewBuffer(newUpdate)
	dec := gob.NewDecoder(buffer)
	var playerUpdate PlayerUpdate

	err := dec.Decode(&playerUpdate)
	if err != nil {
		errChan <- fmt.Errorf("Handle player update decode error:%w", err)
		return
	}
	//fmt.Println("Struct that was recieved:", playerUpdate)

	//Checks if the client has the correct ID associated with its IP. If successful this confirms
	//that both server and client know eachothers IP and ID.
	idMapMutex.RLock()
	storedID := server.idMap[addr]
	idMapMutex.RUnlock()
	if storedID != playerUpdate.ID {
		errChan <- fmt.Errorf("stored ID and recieved ID of %v do not match:%v and %v", addr, storedID, playerUpdate.ID)
		//send error to client, they should restart
		return
	}

	//checks if the ID has an associated IP. If not, maps them together. This enables the server
	//to start sending worldUpdates to that IP, since sendWorldUpdates uses adressMap as a list of IPs to send to.
	adressMapMutex.RLock()
	storedIPAdress := server.adressMap[playerUpdate.ID]
	adressMapMutex.RUnlock()
	if storedIPAdress == "" {
		adressMapMutex.Lock()
		server.adressMap[playerUpdate.ID] = addr
		adressMapMutex.Unlock()
	}

	playerServerMutex.Lock()
	server.playerUpdates = append(server.playerUpdates, playerUpdate)
	playerServerMutex.Unlock()
}

func handleNewBullet(server *ServerConn, newUpdate []byte, addr string, errChan chan<- error) {

	buffer := bytes.NewBuffer(newUpdate)
	dec := gob.NewDecoder(buffer)
	var newBullet BulletUpdate

	err := dec.Decode(&newBullet)
	if err != nil {
		errChan <- fmt.Errorf("Handle bullet update decode error:%w", err)
		return
	}
	//fmt.Println("Struct that was recieved:", updateStruct)

	bulletMapMutex.RLock()
	alreadyFired := server.bulletMap[newBullet.ID]
	bulletMapMutex.RUnlock()

	if !alreadyFired {
		bulletMapMutex.Lock()
		server.bulletMap[newBullet.ID] = true
		bulletMapMutex.Unlock()
	} else {
		return
	}

	bulletServerMutex.Lock()
	server.bulletUpdates = append(server.bulletUpdates, newBullet)
	bulletServerMutex.Unlock()

}

func handlePingRequest(pingID []byte, IPaddr string) {

	addr, err := net.ResolveUDPAddr("udp", IPaddr+":8081")
	if err != nil {
		log.Fatal(err)
	}
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatal(err)
	}

	var buffer bytes.Buffer
	var OP uint8
	OP = 8

	encoder := gob.NewEncoder(&buffer)

	if err = encoder.Encode(OP); err != nil {
		log.Fatal(err)
	}
	buffer.Write(pingID[:20])

	_, err = buffer.WriteTo(conn)
	if err != nil {
		log.Fatal(err)
	}
	conn.Close()
}
