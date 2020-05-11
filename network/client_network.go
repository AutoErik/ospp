package network

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

//PlayerUpdate contains an update of a player position and if the player is shooting
type PlayerUpdate struct {
	ID         int
	X          float64
	Y          float64
	DirectionX float64
	DirectionY float64
	Timestamp  int64
}

//BulletUpdate contains an update of a bullet position
type BulletUpdate struct {
	ID         int
	PlayerID   int
	X          float64
	Y          float64
	DirectionX float64
	DirectionY float64
	TimeAlive  int64
	Timestamp  int64
}

//ClientConn Represents a clients connection to the server.
type ClientConn struct {
	playerPositionUpdates []PlayerUpdate
	bulletPositionUpdates []BulletUpdate
	playerID              *intWrapper
	pingIDCounter         *intWrapper
	pingMap               map[int]int64
}

var playerMutex sync.Mutex
var bulletMutex sync.Mutex
var pingIDCounterMutex sync.Mutex
var pingMapMutex sync.Mutex

//GetPlayerPositionUpdates returns playerUpdates sent from the server.
func (c *ClientConn) GetPlayerPositionUpdates() []PlayerUpdate {

	playerMutex.Lock()

	var Updates []PlayerUpdate
	if len(c.playerPositionUpdates) > 0 {
		Updates = c.playerPositionUpdates
		c.playerPositionUpdates = make([]PlayerUpdate, 0, 10)
	}
	playerMutex.Unlock()
	return Updates
}

//GetBulletPositionUpdates returns bulletUpdates sent from the server.
func (c *ClientConn) GetBulletPositionUpdates() []BulletUpdate {
	bulletMutex.Lock()

	var Updates []BulletUpdate
	if len(c.bulletPositionUpdates) > 0 {
		Updates = c.bulletPositionUpdates
		c.bulletPositionUpdates = make([]BulletUpdate, 0, 10)
	}
	bulletMutex.Unlock()
	return Updates
}

//SendPlayerUpdate Send information about the client player state to the server.
func (c *ClientConn) SendPlayerUpdate(update PlayerUpdate) error {

	addr, err := net.ResolveUDPAddr("udp", ":8080")
	if err != nil {
		return err
	}
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	var OP uint8
	var buf bytes.Buffer
	OP = 2
	encoder := gob.NewEncoder(&buf)

	if err = encoder.Encode(OP); err != nil {
		return err
	}
	if err = encoder.Encode(update); err != nil {
		return err
	}

	//fmt.Println("---------")
	//fmt.Println("Client sent playerupdate to server, size:", len(buf.Bytes()))

	_, err = buf.WriteTo(conn)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

//SendBulletUpdate Shoot a new bullet.
func (c *ClientConn) SendBulletUpdate(newBullet BulletUpdate) {

	addr, err := net.ResolveUDPAddr("udp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	var OP uint8
	var buf bytes.Buffer
	OP = 3
	encoder := gob.NewEncoder(&buf)

	if err = encoder.Encode(OP); err != nil {
		log.Fatal(err)
	}
	if err = encoder.Encode(newBullet); err != nil {
		log.Fatal(err)
	}

	//fmt.Println("---------")
	//fmt.Println("Client sent playerupdate to server, size:", len(buf.Bytes()))
	for i := 0; i < 5; i++ {
		if _, err = io.Copy(conn, &buf); err != nil {
			log.Fatal(err)
		}
		time.Sleep(250 * time.Millisecond)
	}
}

func establishConnection(conn *net.UDPConn) int {

	inputBytes := make([]byte, 200)
	stop := make(chan bool)
	go sendIDRequest(stop)

	for {
		_, _, err := conn.ReadFromUDP(inputBytes)
		if err != nil {
			log.Fatal(err)
		}

		var Opcode uint8
		bufferOP := bytes.NewBuffer(inputBytes[:4])
		gob.NewDecoder(bufferOP).Decode(&Opcode)

		if Opcode == 1 {
			var newID int
			bufferOP := bytes.NewBuffer(inputBytes[4:])
			gob.NewDecoder(bufferOP).Decode(&newID)

			stop <- true
			return newID
		}
	}
}

func sendIDRequest(stop chan bool) {

	addr, err := net.ResolveUDPAddr("udp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	var OP uint8
	var buf bytes.Buffer
	OP = 1

	encoder := gob.NewEncoder(&buf)

	if err = encoder.Encode(OP); err != nil {
		log.Fatal(err)
	}

	fmt.Println("---------")
	fmt.Println("ID request sent from client to server")

	for {
		bufCopy := buf
		_, err = bufCopy.WriteTo(conn)
		if err != nil {
			log.Fatal(err)
		}
		time.Sleep(500 * time.Millisecond)

		select {
		case <-stop:
			return
		default:
			fmt.Println("another ID request sent")
		}
	}
}

//StartClientListener starts listening to server responses.
// Will automatically request and return a new player ID from the server when called.
// Returns a ping channel and an error channel, the error channel has to be emptied regularly.
// Opcode 1: Only used when requesting playerID. Will not be accepted in the main loop.
// Opcode 2: Accepts playerUpdates from the server.
// Opcode 3: Accepts bulletUpdates from the server.
// Opcode 8: Used by Ping() on the way back.
// Opcode 9: Shut down the client(or at least this function).
func StartClientListener() (*ClientConn, int, chan int, <-chan error, error) {

	client := ClientConn{
		playerPositionUpdates: make([]PlayerUpdate, 0, 10),
		bulletPositionUpdates: make([]BulletUpdate, 0, 10),
		playerID:              new(intWrapper),
		pingIDCounter:         new(intWrapper),
		pingMap:               make(map[int]int64),
	}
	errChan := make(chan error, 10) //create function to read from this channel and save it in a log!
	pingChan := make(chan int)      //create function to read from this channel and...use it somehow?

	addr, err := net.ResolveUDPAddr("udp", ":8081")
	if err != nil {
		return &client, 0, nil, nil, err
	}
	udpConn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return &client, 0, nil, nil, err
	}

	//Retrieve player ID from server before starting main loop
	yourID := establishConnection(udpConn)
	if yourID < 1 {
		log.Fatal("Failed to establish connection to server")
	}
	client.playerID.num = yourID

	ctx, cancel := context.WithCancel(context.Background())
	//Main loop
	go func() {
		defer udpConn.Close()
		defer cancel()
		for {
			inputBytes := make([]byte, 200)

			n, addr, err := udpConn.ReadFromUDP(inputBytes)
			if err != nil {
				errChan <- fmt.Errorf("Reading %d bytes from %s stream error:%w", n, addr, err)
				continue
			}

			var Opcode uint8
			bufferOP := bytes.NewBuffer(inputBytes[:4])
			err = gob.NewDecoder(bufferOP).Decode(&Opcode)
			if err != nil {
				errChan <- fmt.Errorf("Decoding opcode [%v] error:%w", inputBytes[:4], err)
				continue
			}

			switch Opcode {
			case 2:
				go handlePlayerUpdate(&client, inputBytes[4:], addr.IP.String(), errChan)
			case 3:
				go handleBulletUpdate(&client, inputBytes[4:], addr.IP.String(), errChan)
			case 8:
				go handlePingResponse(&client, inputBytes[4:], pingChan, errChan)
			case 9:
				return
			default:
				fmt.Println("Faulty Opcode")
			}

			//fmt.Println("---------")
			fmt.Println("Client with id", client.playerID.num, "has gotten a message size:", n)
			fmt.Println("Opcode:", Opcode)
			//fmt.Println("Adress of sender:", addr)
			//fmt.Println("full buffer:", inputBytes[4:])

		}
	}()

	go func() {
		for {
			pingChan <- 0
			time.Sleep(50 * time.Millisecond)

			select {
			case <-ctx.Done():
				return
			default:
				time.Sleep(100 * time.Millisecond)
				Ping(&client)
			}
		}
	}()

	return &client, yourID, pingChan, errChan, nil

}

func handlePlayerUpdate(client *ClientConn, newPlayerUpdate []byte, addr string, errChan chan<- error) {

	buffer := bytes.NewBuffer(newPlayerUpdate)
	dec := gob.NewDecoder(buffer)
	var player PlayerUpdate

	err := dec.Decode(&player)
	if err != nil {
		errChan <- fmt.Errorf("Handle update decode error:%w", err)
		return
	}
	//fmt.Println("Struct that was recieved by client:", player)
	//fmt.Println("Type of struct: Player")

	playerMutex.Lock()
	client.playerPositionUpdates = append(client.playerPositionUpdates, player)
	playerMutex.Unlock()

}

func handleBulletUpdate(client *ClientConn, newBulletUpdate []byte, addr string, errChan chan<- error) {
	buffer := bytes.NewBuffer(newBulletUpdate)
	dec := gob.NewDecoder(buffer)
	var bullet BulletUpdate

	err := dec.Decode(&bullet)
	if err != nil {
		errChan <- fmt.Errorf("Handle update decode error:%w", err)
		return
	}
	//fmt.Println("Struct that was recieved:", bullet)
	//fmt.Println("Type of struct: Bullet")

	bulletMutex.Lock()
	client.bulletPositionUpdates = append(client.bulletPositionUpdates, bullet)
	bulletMutex.Unlock()
}

//Ping sends a ping to the server from the client and measures the roundtrip time in milliseconds.
// Will send the result to the channel passed to StartClientListener().
func Ping(client *ClientConn) {

	addr, err := net.ResolveUDPAddr("udp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	udpConn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatal(err)
	}

	var OP uint8
	var pingID int
	var buf bytes.Buffer
	OP = 8

	pingIDCounterMutex.Lock()
	pingID = client.pingIDCounter.num
	client.pingIDCounter.num++
	pingIDCounterMutex.Unlock()

	t := time.Now()
	milli := t.UnixNano() / 1e6

	pingMapMutex.Lock()
	client.pingMap[pingID] = milli
	pingMapMutex.Unlock()

	encoder := gob.NewEncoder(&buf)

	if err = encoder.Encode(OP); err != nil {
		log.Fatal(err)
	}
	if err = encoder.Encode(pingID); err != nil {
		log.Fatal(err)
	}

	if _, err = buf.WriteTo(udpConn); err != nil {
		log.Fatal(err)
	}
	udpConn.Close()
}

func handlePingResponse(client *ClientConn, inputBytes []byte, pingChan chan<- int, errChan chan<- error) {

	buffer := bytes.NewBuffer(inputBytes)
	dec := gob.NewDecoder(buffer)
	var pingID int

	if err := dec.Decode(&pingID); err != nil {
		fmt.Println("decode error:", err)
		return
	}

	pingMapMutex.Lock()
	sentAt := client.pingMap[pingID]
	pingMapMutex.Unlock()

	t := time.Now()
	RecievedAt := t.UnixNano() / 1e6
	yourPing := int(RecievedAt - sentAt)
	//fmt.Println("Ping is:", yourPing)

	pingChan <- yourPing
}
