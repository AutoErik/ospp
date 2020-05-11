package game

import (
	"dogmatix/network"
	"fmt"
	"sync"
)

type Player struct {
	Id         int
	Xpos       float64
	Ypos       float64
	Name       string
	Width      float64
	Height     float64
	DirectionX float64
	DirectionY float64
}

type PlayerCollection struct {
	Players map[int]Player
}

var PlayerCollectionMutex sync.Mutex

func MakePlayerCollection() PlayerCollection {
	return PlayerCollection{Players: map[int]Player{}}
}

func (col PlayerCollection) SetPlayer(id int, X float64, Y float64, DirectionX float64, DirectionY float64, Name string) {
	PlayerCollectionMutex.Lock()
	col.Players[id] = Player{Id: id, Xpos: X, Ypos: Y, DirectionX: DirectionX, DirectionY: DirectionY, Name: Name, Height: 1, Width: 1}
	PlayerCollectionMutex.Unlock()
}

func (col PlayerCollection) GetPlayer(id int) (Player, bool) {
	PlayerCollectionMutex.Lock()
	player, ok := col.Players[id]
	PlayerCollectionMutex.Unlock()
	return player, ok

}

func (col PlayerCollection) Remove(id int) {
	PlayerCollectionMutex.Lock()
	delete(col.Players, id)
	PlayerCollectionMutex.Unlock()
}

//func newPlayer(Id int, Xpos float64, Ypos float64, Name string, Width float64, Height float64, DirectionX float64, DirectionY float64) Player {
//	var newPlayer = Player{Id, Xpos, Ypos, Name, Width, Height, DirectionX, DirectionY}
//	return newPlayer
//}

//func GetAndSetPlayerPos(LocalPlayers []*Player) {
//	positions := network.GetPlayerUpdates()
//
//	for i := 0; i <= 2; i++ { // Kanske kan ändra så att man alltid uppdaterar? Ändra till fler players?
//		if LocalPlayers[i].Id == positions[i].PlayerId { // Kanske redundant pga loopen ovan
//			if LocalPlayers[i].Xpos != positions[i].X && LocalPlayers[i].Ypos != positions[i].Y {
//				LocalPlayers[i].Xpos = positions[i].X //Ändra till de inkommande positionerna
//				LocalPlayers[i].Ypos = positions[i].Y
//
//				//network.SendPlayerUpdate(positions[i]) Skicka tillbaka något? Isåfall localplayers?
//			}
//		}
//	}
//
//}

func DetectCollisionPlayers(player1 Player, player2 Player) {
	if player1.Xpos < player2.Xpos+player2.Width &&
		player1.Xpos+player1.Width > player2.Xpos &&
		player1.Ypos < player2.Ypos+player2.Height &&
		player1.Ypos+player1.Height > player2.Ypos {
		fmt.Println("collision detected") //tbi vad som sker vid collision
	}
}

//func DetectCollisionWallPlayer(wall Wall, player Player) {
//	if wall.Xpos < player.Xpos+player.Width &&
//		wall.Xpos+wall.Width > player.Xpos &&
//		wall.Ypos < player.Ypos+player.Height &&
//		wall.Ypos+wall.Height > player.Ypos {
//		fmt.Println("collision detected") //tbi vad som sker vid collision
//	}
//}

func PlayersThatMoved(col PlayerCollection, PlayerUpdates []network.PlayerUpdate) []int {

	var MovedPlayersIDs []int

	for i := 0; i < len(PlayerUpdates); i++ {
		Id := PlayerUpdates[i].ID
		currPlayer, ok := col.GetPlayer(Id)

		if !ok {
			col.SetPlayer(PlayerUpdates[i].ID, PlayerUpdates[i].X, PlayerUpdates[i].Y, PlayerUpdates[i].DirectionX, PlayerUpdates[i].DirectionY, "")
			MovedPlayersIDs = append(MovedPlayersIDs, PlayerUpdates[i].ID)
		} else if PlayerUpdates[i].X != currPlayer.Xpos || PlayerUpdates[i].Y != currPlayer.Ypos { // that player has moved.
			MovedPlayersIDs = append(MovedPlayersIDs, PlayerUpdates[i].ID)
		}

	}
	return MovedPlayersIDs
}

// vi vet vilka som flyttat sig. nu måste vi hitta subgrid och skicka in allt i go rutins.
