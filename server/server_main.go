package server

import (
	"dogmatix/game"
	"dogmatix/network"
	"fmt"
	"log"
	"sync"
	"time"
)

type BulletDelta struct {
	Delta     float64
	NewBullet network.BulletUpdate
}

var bullets game.BulletCollection = game.NewBulletCollection()
var players game.PlayerCollection = game.MakePlayerCollection()

var ChannelMutex sync.Mutex

// Huvudfunktionen för servern, kallas från main
func ServerMain() {

	serverConn, errChan, err := network.StartServerListener()
	if err != nil {
		log.Fatal(err, "Could not fetch ServerListener")
	}

	connClient, clientID, _, _, err := network.StartClientListener()
	if err != nil {
		log.Fatal(err)
	}

	createGrid(100, 100, 50, 50)

	var wg1 sync.WaitGroup
	var wg2 sync.WaitGroup
	var wg3 sync.WaitGroup
	var wg4 sync.WaitGroup

	bulletchan1 := make(chan BulletDelta)
	bulletchan2 := make(chan BulletDelta)
	bulletchan3 := make(chan BulletDelta)
	bulletchan4 := make(chan BulletDelta)

	channel1 := make(chan network.PlayerUpdate, 10)
	channel2 := make(chan network.PlayerUpdate, 10)
	channel3 := make(chan network.PlayerUpdate, 10)
	channel4 := make(chan network.PlayerUpdate, 10)

	go updateObjects(players, bullets, &wg1, bulletchan1, channel1)
	go updateObjects(players, bullets, &wg2, bulletchan2, channel2)
	go updateObjects(players, bullets, &wg3, bulletchan3, channel3)
	go updateObjects(players, bullets, &wg4, bulletchan4, channel4)
	update := network.PlayerUpdate{
		ID:         clientID,
		X:          22,
		Y:          74,
		DirectionX: 54,
		Timestamp:  543,
	}
	delta := 1.0

	//loop := NewLoop(1, func(delta float64) {
	for {
		time.Sleep(1000 * time.Millisecond)
		connClient.SendPlayerUpdate(update)
		update.X += 1
		fmt.Println(bullets, "bullets")
		select {
		case err := <-errChan:
			fmt.Println(err, "ERROR")
		default:
		}

		playerUpdates := serverConn.GetPlayerUpdates()
		bulletUpdates := serverConn.GetBulletUpdates()
		fmt.Println("playerupdates:", playerUpdates)
		Moved := game.PlayersThatMoved(players, playerUpdates)

		for i := 0; i < len(Moved); i++ {
			var PlayerUpdateX network.PlayerUpdate

			for k := 0; k < len(playerUpdates); k++ {

				if Moved[i] == playerUpdates[k].ID {

					PlayerUpdateX = playerUpdates[k]

				}
			}

			subgrid := findSubGrid(PlayerUpdateX.X, PlayerUpdateX.Y)
			switch subgrid {
			case (1):
				wg1.Add(1)
				fmt.Println("SubGrid1")
				channel1 <- PlayerUpdateX

			case (2):
				wg2.Add(1)
				fmt.Println("SubGrid2")
				channel2 <- PlayerUpdateX

			case (3):
				wg3.Add(1)
				fmt.Println("SubGrid3")
				channel3 <- PlayerUpdateX

			case (4):
				wg4.Add(1)
				fmt.Println("SubGrid4")
				channel4 <- PlayerUpdateX

			default:
				log.Fatal("subgrid returned", subgrid)
			}
		}

		wg1.Wait()
		wg2.Wait()
		wg3.Wait()
		wg4.Wait()
		for _, v := range bulletUpdates {
			toSend := BulletDelta{NewBullet: v, Delta: delta}
			subgridBullet := findSubGrid(v.X, v.Y)
			switch subgridBullet {
			case (1):
				wg1.Add(1)
				bulletchan1 <- toSend
				wg1.Wait()
				break
			case (2):
				wg2.Add(1)
				bulletchan2 <- toSend
				wg2.Wait()
				break
			case (3):
				wg3.Add(1)
				bulletchan3 <- toSend
				wg3.Wait()
				break
			case (4):
				wg4.Add(1)
				bulletchan4 <- toSend
				wg4.Wait()
				break
			default:
				log.Fatal("subgrid returned", subgridBullet)

			}
		}
		wg1.Wait()
		wg2.Wait()
		wg3.Wait()
		wg4.Wait()
		fmt.Println("------------------------------------------------------------------------------------------------")
		for _, v := range bullets.Bullets {
			b := network.BulletUpdate{ID: v.Id, PlayerID: v.PlayerId, X: v.Xpos, Y: v.Ypos, DirectionX: v.DirectionX, DirectionY: v.DirectionY}
			serverConn.SendWorldUpdate(b)
		}
		for _, v := range players.Players {
			p := network.PlayerUpdate{ID: v.Id, X: v.Xpos, Y: v.Ypos, DirectionX: v.DirectionX, DirectionY: v.DirectionY}
			serverConn.SendWorldUpdate(p)
		}
	}
	//	})
	/*	loop.Start()
		for { //TODO: Change to something better
		}*/
}
