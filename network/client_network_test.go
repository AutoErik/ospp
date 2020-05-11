package network

import (
	"bytes"
	"encoding/gob"
	"log"
	"net"
	"testing"
	"time"
)

func TestClientIDRequest(t *testing.T) {

	_, _, err := StartServerListener()
	if err != nil {
		log.Fatal(err)
	}
	_, clientID, _, _, err := StartClientListener()
	if err != nil {
		log.Fatal(err)
	}
	if clientID != 1 {
		t.Error("Client recieved wrong ID, should start from 1:", clientID)
	}

	helpFunCloseServer()
	helpFunCloseClient()
	time.Sleep(50 * time.Millisecond)
}

func TestPing(t *testing.T) {

	_, _, err := StartServerListener()
	if err != nil {
		log.Fatal(err)
	}
	_, _, pingChan, _, err := StartClientListener()
	if err != nil {
		log.Fatal(err)
	}

	timesPinged := 0
	//go helpFunPinger(connClient) //runs Ping() 10 times
	for {
		ping := <-pingChan
		//fmt.Println("Ping:", ping)
		if ping > 4000 {
			t.Error("Pings are taking too long")
		}
		timesPinged++
		if timesPinged > 4 {
			break
		}
	}
	helpFunCloseServer()
	helpFunCloseClient()
	time.Sleep(50 * time.Millisecond)
}

func TestClientSendUpdate(t *testing.T) {

	connServer, _, err := StartServerListener()
	if err != nil {
		log.Fatal(err)
	}

	connClient, clientID, _, errChan, err := StartClientListener()
	if err != nil {
		log.Fatal(err)
	}

	update := PlayerUpdate{
		ID:         clientID,
		X:          22,
		Y:          74,
		DirectionX: 54,
		Timestamp:  543,
	}
	connClient.SendPlayerUpdate(update)
	time.Sleep(50 * time.Millisecond)

	select {
	case err := <-errChan:
		log.Fatal("errChan recieved error:", err)
	default:
	}

	var i int
	for i = 0; i < 5; i++ {
		if i == 4 {
			t.Error("Timeout while waiting for getPlayerUpdates to recieve data")
		}
		time.Sleep(100 * time.Millisecond)
		playerUpdates := connServer.GetPlayerUpdates()
		if len(playerUpdates) != 0 {
			if playerUpdates[0].ID != 1 ||
				playerUpdates[0].X != 22 ||
				playerUpdates[0].Y != 74 ||
				playerUpdates[0].DirectionX != 54 ||
				playerUpdates[0].Timestamp != 543 {
				t.Error("Incorrect playerupdate values recieved from client")
			}
			break
		}
	}

	select {
	case err := <-errChan:
		log.Fatal("errChan recieved error:", err)
	default:
	}

	helpFunCloseServer()
	helpFunCloseClient()
	time.Sleep(50 * time.Millisecond)
}

func TestClientSendBullet(t *testing.T) {

	connServer, _, err := StartServerListener()
	if err != nil {
		log.Fatal(err)
	}

	connClient, _, _, errChan, err := StartClientListener()
	if err != nil {
		log.Fatal(err)
	}

	update := BulletUpdate{
		ID:         49,
		X:          22,
		Y:          74,
		DirectionX: 54,
		Timestamp:  543,
	}
	connClient.SendBulletUpdate(update)
	time.Sleep(100 * time.Millisecond)

	select {
	case err := <-errChan:
		log.Fatal("errChan recieved error:", err)
	default:
	}

	var i int
	for i = 0; i < 5; i++ {
		if i == 4 {
			t.Error("Timeout while waiting for getPlayerUpdates to recieve data")
		}
		time.Sleep(100 * time.Millisecond)
		bulletUpdates := connServer.GetBulletUpdates()
		if len(bulletUpdates) != 0 {
			if bulletUpdates[0].ID != 49 ||
				bulletUpdates[0].X != 22 ||
				bulletUpdates[0].Y != 74 ||
				bulletUpdates[0].DirectionX != 54 ||
				bulletUpdates[0].Timestamp != 543 {
				t.Error("Incorrect Bulletupdate values recieved from client")
			}
			break
		}
	}

	select {
	case err := <-errChan:
		log.Fatal("errChan recieved error:", err)
	default:
	}

	helpFunCloseServer()
	helpFunCloseClient()
	time.Sleep(50 * time.Millisecond)
}

func helpFunCloseServer() {

	addr, err := net.ResolveUDPAddr("udp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatal(err)
	}

	var OP uint8
	var buf bytes.Buffer
	OP = 9
	encoder := gob.NewEncoder(&buf)

	if err = encoder.Encode(OP); err != nil {
		log.Fatal(err)
	}

	_, err = buf.WriteTo(conn)
	if err != nil {
		log.Fatal(err)
	}
	conn.Close()
	time.Sleep(20 * time.Millisecond)
}

func helpFunCloseClient() {

	addr, err := net.ResolveUDPAddr("udp", ":8081")
	if err != nil {
		log.Fatal(err)
	}
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatal(err)
	}

	var OP uint8
	var buf bytes.Buffer
	OP = 9
	encoder := gob.NewEncoder(&buf)

	if err = encoder.Encode(OP); err != nil {
		log.Fatal(err)
	}

	_, err = buf.WriteTo(conn)
	if err != nil {
		log.Fatal(err)
	}
	conn.Close()
}

func helpFunPinger(conn *ClientConn) {

	//	parent, cancel := context.WithCancel(context.Background())

	//	child := context.Context(parent)

	for i := 0; i < 10; i++ {
		time.Sleep(50 * time.Millisecond)
		go Ping(conn)
	}
	//	cancel()
}
