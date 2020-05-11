package server

import (
	"dogmatix/game"
	"dogmatix/network"
	"fmt"
	"sync"
	"testing"
)

//create 2 players and see if we have created 2 players
func TestPlayerCreate(t *testing.T) {
	player := game.MakePlayerCollection()
	player.SetPlayer(int(0), float64(0.0), float64(0.0), float64(0.0), float64(0.0), "name")
	player.SetPlayer(int(1), float64(0.0), float64(0.0), float64(0.0), float64(0.0), "name")

	expected := 2
	actual := len(player.Players)
	if actual == expected {
		fmt.Println("TestPlayerCreate: PASSED")
	}
	if actual != expected {
		t.Error("TestPlayerCreate: FAILED!")
	}
}

//give a player an ID and see if the ID is the same as given
func TestPlayerID(t *testing.T) {
	player := game.MakePlayerCollection()
	player.SetPlayer(int(2), float64(0.0), float64(0.0), float64(0.0), float64(0.0), "name")
	expected := 2
	actual, ok := player.GetPlayer(2)
	if ok && actual.Id == expected {
		fmt.Println("TestPlayerId: PASSED")
	}
	if !ok || actual.Id != expected {
		t.Error("TestPlayerId: FAILED!")
	}
}

//give a player a posision and see if the given posision is correct (as actual posision)
func TestPlayerPosision(t *testing.T) {
	players := game.MakePlayerCollection()
	players.SetPlayer(int(3), float64(1.0), float64(1.0), float64(0.0), float64(0.0), "name")
	expected := 1.0
	actual, ok := players.GetPlayer(3)
	if ok && actual.Xpos == 1 && actual.Ypos == 1 {
		fmt.Println("TestPlayerPosision: PASSED")
	}
	if !ok || actual.Xpos != expected || actual.Ypos != expected {
		t.Error("TestPlayerPosision: FAILED!")
	}
}

//create 2 bullets and se if we create 2 bullets correct
func TestBulletCreate(t *testing.T) {
	bullet := game.NewBulletCollection()
	bullet.SetBullet(0, 0, 0.0, 0.0, 0.0, 0.0)
	bullet.SetBullet(1, 0, 0.0, 0.0, 0.0, 0.0)

	expected := 2
	actual := len(bullet.Bullets)
	if actual == expected {
		fmt.Println("TestBulletCreate: PASSED")
	}
	if actual != expected {
		fmt.Println("TestBulletCreate: FAILED!")
	}
}

//give a bullet an ID and see if the ID is the same as given
func TestBulletID(t *testing.T) {
	bullet := game.NewBulletCollection()
	bullet.SetBullet(2, 0, 0.0, 0.0, 0.0, 0.0)
	expected := 2
	actual, ok := bullet.GetBullet(2)
	if ok && actual.Id == expected {
		fmt.Println("TestBulletId: PASSED")
	}
	if !ok || actual.Id != expected {
		t.Error("TestBulletId: FAILED!")
	}
}

//give a bullet a posision and se if the given posision is correct
func TestBulletPosision(t *testing.T) {
	bullet := game.NewBulletCollection()
	bullet.SetBullet(3, 0, 1.0, 1.0, 0.0, 0.0)
	expected := 1.0
	actual, ok := bullet.GetBullet(3)
	if ok && actual.Xpos == 1.0 && actual.Ypos == 1.0 {
		fmt.Println("TestBulletPosision: PASSED")
	}
	if !ok || actual.Xpos != expected || actual.Ypos != expected {
		t.Error("TestBulletPosision: FAILED!")
	}
}

//create a player and give the player a posision and see if the palyer has moved
func TestPlayerMove(t *testing.T) {
	createGrid(100, 100, 50, 50)
	var waitG1 sync.WaitGroup
	player1 := game.MakePlayerCollection()
	newPlayer := network.PlayerUpdate{2, 3, 3, 1, 1, false}
	player1.SetPlayer(int(2), float64(2.0), float64(2.0), float64(0.0), float64(0.0), "Player 1")
	waitG1.Add(1)
	playerMove(player1, newPlayer, &waitG1)
	waitG1.Wait()
	expected := 3.0
	actual, ok := player1.GetPlayer(2)
	if ok && actual.Xpos == expected && actual.Ypos == expected {
		fmt.Println("TestPlayerMove: PASSED")
	} else {
		t.Error("TestPlayerMove: FAILED!")
	}

}

//create a player and bulletscollision and see if we updated the bullet

func TestMoveBullet(t *testing.T) {
	createGrid(100, 100, 50, 50)
	var waitG1 sync.WaitGroup
	player2 := game.MakePlayerCollection()
	player2.SetPlayer(int(2), float64(0.0), float64(0.0), float64(0.0), float64(0.0), "Player 1")
	bullets := game.NewBulletCollection()
	bullets.SetBullet(2, 0, 1.0, 1.0, 1.0, 1.0)
	bullet, ok := bullets.GetBullet(2)
	waitG1.Add(1)
	bulletDelta := BulletDelta{NewBullet: bullet, Delta: 1}
	bulletMove(bullets, player2, &waitG1, bulletDelta)
	waitG1.Wait()
	expected := 2.0
	actual, ok := bullets.GetBullet(2)
	fmt.Println(actual)
	if ok && actual.Xpos == expected && actual.Ypos == expected {
		fmt.Println("TestBulletMove: PASSED")
	} else {
		t.Error("TestBulletMove: FAILED!")
	}
}

func TestCollisionPlayers(t *testing.T) {
	createGrid(100, 100, 50, 50)
	expectedPlayer1 := 10.0
	expectedPlayer2 := 15.0
	var wg1 sync.WaitGroup
	playerColl := game.MakePlayerCollection()
	x := float64(10.0)
	y := float64(10.0)
	x1 := float64(20)
	y1 := float64(10.0)
	playerColl.SetPlayer(1, 10, 10, 1, 0, "")
	playerColl.SetPlayer(2, 20, 10, 1, 0, "")
	newPlayer := network.PlayerUpdate{1, x, y, 1, 0, false}
	newPlayer1 := network.PlayerUpdate{2, x1, y1, 0, 0, false}

	for i := 0; i < 20; i++ {
		wg1.Add(1)
		playerMove(playerColl, newPlayer, &wg1)
		wg1.Wait()
		wg1.Add(1)
		playerMove(playerColl, newPlayer1, &wg1)
		wg1.Wait()
		iteratePlayer, ok := playerColl.GetPlayer(1)
		iteratePlayer1, ok := playerColl.GetPlayer(2)
		if ok {
			x = iteratePlayer.Xpos + 1
			x1 = iteratePlayer1.Xpos - 1
			newPlayer = network.PlayerUpdate{1, x, y, 1, 0, false}
			newPlayer1 = network.PlayerUpdate{2, x1, y1, 0, 0, false}
		}

	}
	actual1, ok := playerColl.GetPlayer(1)
	actual2, ok := playerColl.GetPlayer(2)
	if ok && actual1.Xpos == expectedPlayer1 && actual2.Xpos == expectedPlayer2 {
		fmt.Println("TestCollisionPlayer: PASSED!")
	} else {
		t.Error("TestCollisionPlayer: FAILED!")
	}
}

func TestBulletCollision(t *testing.T) {
	createGrid(100, 100, 50, 50)
	expected := true
	actual := false
	var wg1 sync.WaitGroup
	playerColl := game.MakePlayerCollection()
	bulletsColl := game.NewBulletCollection()

	bulletsColl.SetBullet(1000, 2, 10, 10, 1, 0)
	for i := 0; i < 20; i++ {
		IterateBullet, ok := bulletsColl.GetBullet(1000)
		delta := BulletDelta{NewBullet: IterateBullet, Delta: 1}
		if ok {
			wg1.Add(1)
			bulletMove(bulletsColl, playerColl, &wg1, delta)
			wg1.Wait()
		} else {
			actual = true
			break
		}
	}
	if expected == actual {
		fmt.Println("TestBulletCollision: PASSED!")
	} else {
		t.Error("TestBulletCollision: FAILED!")
	}
}

func TestWallCollisionSet(t *testing.T) {
	createGrid(100, 100, 5, 5)
	wallColl := game.MakeWallCollection()
	wallColl.SetWall(1, 2, 1, 1, 1)
	setWallCollision(wallColl)
	//for i := 0; i < len(gridList); i++ {
	//	println("Index: %f, Wall: %f\n", gridList[i].Index, gridList[i].Wall)
	//}
	act := gridList[1].Wall
	ual := gridList[1].Occupied
	if act && ual == true {
		fmt.Println("TestPlayerCreate: PASSED")
	} else {
		t.Errorf("TestPlayerCreate: FAILED!, acutal: %v, %v", act, ual)
	}
}

func TestFindSubGrid(t *testing.T) {
	createGrid(100, 100, 10, 10)
	actual1 := findSubGrid(1, 2)
	actual2 := findSubGrid(90, 90)
	if actual1 == 1 && actual2 == 4 {
		fmt.Println("TestFindSubGrid: PASSED")
	} else {
		t.Errorf("TestFindSubGrid: FAILED!, actual: %v, %v", actual1, actual2)
	}

}
