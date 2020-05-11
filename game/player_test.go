package game

import (
	"fmt"
	"testing"
)


func TestPlayerCollection(t *testing.T) {
	// Should add a new player of id '0'
	players := MakePlayerCollection()
	players.SetPlayer(int(0), float64(0.0), float64(0.0), float64(0.0), float64(0.0), "name")
	if len(players.Players) != 1 {
		fmt.Println(len(players.Players))
		t.Fail()
	}

	// Should overwrite player of id '0'

	players.SetPlayer(int(0), float64(0.0), float64(0.0), float64(0.0), float64(0.0), "name")
	if len(players.Players) != 1 {
		t.Fail()
	}
	// Data should be correct
	player, ok := players.GetPlayer(0)
	if !ok || player.Xpos != 0 {
		t.Fail()
	}

	// Should add a new player of id '2'

	players.SetPlayer(int(1), float64(2.0), float64(2.0), float64(2.0), float64(2.0), "2")
	if len(players.Players) != 2 {
		t.Fail()
	}

	players.Remove(0)
	if len(players.Players) != 1 {
		fmt.Println("Remove with id 2 fails")
		t.Fail()
	}

}

//func TestPlayersThatMoved(t *testing.T) {
//	//TBI, network tester?
//}
