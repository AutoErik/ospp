package network

import (
	"bytes"
	"encoding/gob"
	"log"
	"net"
	"testing"
	"time"
)

func TestGetPlayerUpdates(t *testing.T) {

	conn, _, err := StartServerListener()
	if err != nil {
		log.Fatal(err)
	}

	updateStruct := PlayerUpdate{
		ID:         1,
		X:          22,
		Y:          74,
		DirectionX: 54,
		Timestamp:  543,
	}

	playerUpdates := conn.GetPlayerUpdates()
	if len(playerUpdates) > 0 {
		t.Error("Playerupdate should be empty")
	}

	conn.playerUpdates = append(conn.playerUpdates, updateStruct)
	playerUpdates = conn.GetPlayerUpdates()

	playerUpdatesTest := conn.GetPlayerUpdates()
	if len(playerUpdatesTest) > 0 {
		t.Error("Playerupdate should be empty")
	}

	if playerUpdates[0].ID != 1 ||
		playerUpdates[0].X != 22 ||
		playerUpdates[0].Y != 74 ||
		playerUpdates[0].DirectionX != 54 ||
		playerUpdates[0].Timestamp != 543 {
		t.Error("Incorrect playerupdate values recieved from client")
	}

	helpFunCloseServer()
	time.Sleep(50 * time.Millisecond)
}

func TestServerStart(t *testing.T) {

	_, errChan, err := StartServerListener()
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(250 * time.Millisecond)

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
	OP = 7
	encoder := gob.NewEncoder(&buf)

	if err = encoder.Encode(OP); err != nil {
		log.Fatal(err)
	}

	_, err = buf.WriteTo(conn)
	if err != nil {
		log.Fatal(err)
	}

	select {
	case err = <-errChan:
		if err.Error() != "Server is running" {
			t.Error("Server responded with wrong string:", err)
		}
	case <-time.After(5 * time.Second):
		t.Error("Timeout from server")
	}

	helpFunCloseServer()
	time.Sleep(50 * time.Millisecond)
}

func TestServerIPHandling(t *testing.T) {
	connServer, errChan, err := StartServerListener()
	if err != nil {
		log.Fatal(err)
	}

	select {
	case err := <-errChan:
		log.Fatal("errChan recieved error:", err)
	default:
	}

	connClient, clientID, _, _, err := StartClientListener()
	if err != nil {
		log.Fatal(err)
	}

	playerTestUpdate := PlayerUpdate{
		ID:         clientID,
		X:          22,
		Y:          74,
		DirectionX: 54,
		Timestamp:  543,
	}
	connClient.SendPlayerUpdate(playerTestUpdate)
	time.Sleep(100 * time.Millisecond)

	storedIP := connServer.adressMap[clientID]
	if storedIP != "127.0.0.1" {
		t.Error("Stored IP adress does not map to IP of client")
	}

	if connServer.idMap[storedIP] != clientID {
		t.Error("Stored ID does not map to ID given to client")
	}

	helpFunCloseServer()
	helpFunCloseClient()
	time.Sleep(50 * time.Millisecond)
}

func TestServerSendUpdates(t *testing.T) {

	connServer, _, err := StartServerListener()
	if err != nil {
		log.Fatal(err)
	}

	connClient, clientID, _, errChan, err := StartClientListener()
	if err != nil {
		log.Fatal(err)
	}

	playerTestUpdate := PlayerUpdate{
		ID:         clientID,
		X:          22,
		Y:          74,
		DirectionX: 54,
		Timestamp:  543,
	}
	connClient.SendPlayerUpdate(playerTestUpdate)

	bulletUpdate := BulletUpdate{
		ID:         3,
		X:          7542,
		Y:          85334,
		DirectionX: 4.4,
		Timestamp:  999999,
	}

	for i := 0; i < 5; i++ {
		if i == 4 {
			t.Error("Timeout while waiting for getPlayerPositionUpdates to recieve data")

		}
		time.Sleep(100 * time.Millisecond)
		playerUpdates := connServer.GetPlayerUpdates()

		if len(playerUpdates) != 0 {
			err = connServer.SendWorldUpdate(playerUpdates[0])
			if err != nil {
				log.Fatal(err)
			}
			err = connServer.SendWorldUpdate(bulletUpdate)
			if err != nil {
				log.Fatal(err)
			}
			break
		}
	}

	for i := 0; i < 5; i++ {
		if i == 4 {
			t.Error("Timeout while waiting for getPlayerPositionUpdates to recieve data")
		}
		time.Sleep(100 * time.Millisecond)
		playerUpdates := connClient.GetPlayerPositionUpdates()

		if len(playerUpdates) != 0 {
			if playerUpdates[0].ID != 1 ||
				playerUpdates[0].X != 22 ||
				playerUpdates[0].Y != 74 ||
				playerUpdates[0].DirectionX != 54 ||
				playerUpdates[0].Timestamp != 543 {
				t.Error("Incorrect playerUpdate values recieved from the server")
			}
			break
		}
	}

	for i := 0; i < 5; i++ {
		if i == 4 {
			t.Error("Timeout while waiting for getBulletPositionUpdates to recieve data")
		}
		time.Sleep(100 * time.Millisecond)
		bulletUpdates := connClient.GetBulletPositionUpdates()
		if len(bulletUpdates) != 0 {
			if bulletUpdates[0].ID != 3 ||
				bulletUpdates[0].X != 7542 ||
				bulletUpdates[0].Y != 85334 ||
				bulletUpdates[0].DirectionX != 4.4 ||
				bulletUpdates[0].Timestamp != 999999 {
				t.Error("Incorrect bulletUpdate values recieved from the server")
			}
			break
		}
	}

	select {
	case err := <-errChan:
		log.Fatal("errChan recieved error:", err)
	default:
	}
	//fmt.Println("idmap:", connServer.idMap)
	//fmt.Println("adressmap:", connServer.adressMap)
	//fmt.Println("counter:", connServer.idCounter.num)

	helpFunCloseServer()
	helpFunCloseClient()
	time.Sleep(50 * time.Millisecond)
}
