package server

import (
	"dogmatix/game"
	"dogmatix/network"
	"fmt"
	"math"
	"sync"
)

type Tile struct {
	Index    int
	Occupied bool
	Player   bool
	PlayerId int
	Bullet   bool
	Wall     bool
	Floor    bool
}

//Funderar på att lägga till en Rectangle struct. Avvaktar

var (
	BulletHeight float64
	BulletWidth  float64
	BulletSpeed  float64
	mapHeight    float64
	mapWidth     float64
	gridAmount   int
	YOffset      int
	gridList     []Tile
	RowSize      float64
	ColumnSize   float64

)

//Skapar en grid baserat på skärmstorlek och gridmängd
func createGrid(width float64, height float64, gridAmountX int, gridAmountY int) {
	YOffset = gridAmountY
	BulletHeight = 2 //Magic number
	BulletWidth = 2  //Magic number
	BulletSpeed = 1  //Magico

	gridAmount = gridAmountY * gridAmountX
	mapHeight = height
	mapWidth = width
	ColumnSize = width / float64(gridAmountX)
	RowSize = height / float64(gridAmountY)
	q := 0

	for i := 0; i < gridAmountY; i++ {

		for k := 0; k < gridAmountX; k++ { //En tile med flags som t.ex Occupied.
			q += 1
			var tileToAppend = Tile{Index: q, Occupied: false, Player: false, PlayerId: -1, Bullet: false, Wall: false, Floor: false}

			gridList = append(gridList, tileToAppend)

		}
	}
}

// Tar reda på vilka Tiles som blir intersectade i både x och y
// Returvärde kommer bli en lista med Tiles
func findIntersectTiles(minX float64, minY float64, maxX float64, maxY float64) []Tile {
	xTiles := 1
	yTiles := 1
	TileminX := int(math.Floor(minX / ColumnSize))
	TilemaxX := int(math.Floor(maxX / ColumnSize))
	TileminY := int(math.Floor(minY / RowSize))
	TilemaxY := int(math.Floor(maxY / RowSize))

	var YTiles []int
	var XTiles []int

	YTiles = append(YTiles, TileminY, TilemaxY)
	XTiles = append(XTiles, TileminX, TilemaxX)

	if TileminX != TilemaxX {
		xTiles = 2
	}

	if TileminY != TilemaxY {
		yTiles = 2
	}

	var IntersectTiles []Tile
	for i := 0; i < yTiles; i++ {

		for k := 0; k < xTiles; k++ {

			IntersectTiles = append(IntersectTiles, gridList[((YTiles[i]*(YOffset-1))+XTiles[k])]) //Hämtar fram alla relevanta tiles. Logiken är rader * yoffset + xoffset

		}
	}

	return IntersectTiles
}

// Här borde en lista med koordinater vara input. Om man vill, också height/width
func setWallCollision(walls game.WallCollection) {

	for _, coords := range walls.Walls { // skapa två listor för alla väggars x o y tiles
		Intersecting := findIntersectTiles((coords.X - (coords.Width / 2)), (coords.Y - (coords.Height / 2)), (coords.X + (coords.Width / 2)), (coords.Y + (coords.Height / 2)))
		for _, tile := range Intersecting {
			gridList[tile.Index].Wall = true
			gridList[tile.Index].Occupied = true
		}

	}
	/*
		AllTileminX := int(math.Floor(minX / ColumnSize))
		AllTilemaxX := int(math.Floor(maxX / ColumnSize))
		AllTileminY := int(math.Floor(minY / RowSize))
		AllTilemaxY := int(math.Floor(maxY / RowSize))

		var AllYTiles []int
		var AllXTiles []int
		//AllYTiles := list.New()
		//AllXTiles := list.New()

		AllYTiles = append(AllYTiles, AllTileminY, AllTilemaxY) //Lista av alla X tiles
		AllXTiles = append(AllXTiles, AllTileminX, AllTilemaxX) //Lista av alla Y tiles
	*/

	/*
		MinWallX := int(math.Floor(FirstX / ColumnSize))
		MaxWallX := int(math.Floor((FirstX + coords.Width) / ColumnSize))
		MinWallY := int(math.Floor(FirstY / RowSize))
		MaxWallY := int(math.Floor((FirstY + coords.Height) / RowSize))

		AllWallTilesX = append(AllWallTilesX, MinWallX, MaxWallX)
		AllWallTilesY = append(AllWallTilesY, MinWallY, MaxWallY)
		for i, v := range AllYTiles {
			if v == AllWallTilesY[i] {
				fmt.Println("Ändra till True på wall")
				gridList[i].Wall = true
			}
		}

		for i, v := range AllXTiles {
			if v == AllWallTilesX[i] {
				fmt.Println("Ändra till True på wall")
				gridList[i].Wall = true
			}
		}
	*/

}

//Dåligt namn på denna, den returnerar en lista med allt newlist har som inte oldlisthar
func findNewTiles(newlist []Tile, oldlist []Tile) []Tile {
	var newTiles []Tile

	for tile1 := range newlist {

		for tile2 := range oldlist {

			if newlist[tile1].Index != oldlist[tile2].Index {

				newTiles = append(newTiles, newlist[tile1])
			}
		}
	}
	return newTiles
}

// Tittar om en kollision har skett
func collisionChecker(tiles []Tile) Tile {
	for i := 0; i < len(tiles); i++ {
		if tiles[i].Occupied {

			return tiles[i]
		}
	}
	return Tile{Index: -1, Occupied: false, Bullet: false, Wall: false, Floor: false}
}

// Tar reda på vilken/vilka tiles spelaren står på
// Returnerar en lista med Tiles
func playerFindCurrentTiles(player game.Player) []Tile {

	minX := player.Xpos - player.Width/2
	minY := player.Ypos - player.Height/2
	maxX := player.Xpos + player.Width/2
	maxY := player.Ypos + player.Height/2
	return findIntersectTiles(minX, minY, maxX, maxY)

}

// BulletWidth och BulletHeight är inte definierade ännu
// Flyttar alla Bullets i BulletList om de inte kolliderar.


func bulletMove(BulletCollection game.BulletCollection, PlayerCollection game.PlayerCollection, wg *sync.WaitGroup, ch BulletDelta) {

	delta, Bullet := ch.Delta, ch.NewBullet

	remove := false
	BulletTiles := findIntersectTiles(Bullet.X-BulletWidth/2, Bullet.Y-BulletHeight/2, Bullet.X+BulletWidth/2, Bullet.Y+BulletWidth/2)
	CollideTile := collisionChecker(BulletTiles)
	if CollideTile.Occupied {
		remove = true

		for k := 0; k < len(BulletTiles); k++ {
			if BulletTiles[k].Player {

				if BulletTiles[k].PlayerId == Bullet.PlayerID {
					remove = false
				}
				//else {
				//tbi: dealDamageToPlayer(BulletTiles[k].PlayerId, PlayerCollection) }


			}

		}
	}
	if remove {
		fmt.Println("BULLET REMOVED")
		BulletCollection.RemoveBullet(Bullet.PlayerID) //Den har kolliderat, alltså vill vi ta bort kulan

	}
	remove = true
	if !CollideTile.Occupied {
		Bullet.X += Bullet.DirectionX * delta * BulletSpeed
		Bullet.Y += Bullet.DirectionY * delta * BulletSpeed
		BulletCollection.SetBullet(Bullet.ID, Bullet.PlayerID, Bullet.X, Bullet.Y, Bullet.DirectionX, Bullet.DirectionY)

	}
	wg.Done()

}

func findSubGrid(X float64, Y float64) int {
	X = math.Floor(X / ColumnSize)
	Y = math.Floor(Y / RowSize)
	Index := X + (Y * float64(YOffset-1))
	subGridSize := float64(gridAmount / 4)
	subGrid := int(math.Floor(Index/subGridSize) + 1)
	if subGrid > 4 {
		return 4
	}
	if subGrid < 1 {
		return 1
	}
	return subGrid
}

func outOfBoundsCheck(Xpos float64, Ypos float64, Width float64, Height float64) (newX, newY float64) {

	if (Xpos - Width/2) < 0 {
		Xpos = Width / 2
	}
	if Xpos+Width/2 > mapWidth {
		Xpos = mapWidth - Width/2
	}
	if (Ypos - Height/2) < 0 {
		Ypos = Height / 2
	}
	if (Ypos + Height/2) > mapHeight {
		Ypos = mapHeight - Height
	}
	return Xpos, Ypos


}

func updateObjects(players game.PlayerCollection, bullets game.BulletCollection, wg *sync.WaitGroup, bulletChan chan BulletDelta, playerChan chan network.PlayerUpdate) {
	for {
		select {
		case playerUpdate := <-playerChan:
			playerMove(players, playerUpdate, wg)
		case bulletUpdate := <-bulletChan:
			bulletMove(bullets, players, wg, bulletUpdate)
		}
	}
}

func playerMove(col game.PlayerCollection, ch network.PlayerUpdate, wg *sync.WaitGroup) {

	playerUpdate := ch

	player, ok := col.GetPlayer(playerUpdate.ID)
	if ok == false {
		panic(ok)
	}

	newX, newY := outOfBoundsCheck(playerUpdate.X, playerUpdate.Y, player.Width, player.Height)
	newLeft := newX - player.Width/2
	newBottom := newY - player.Height/2
	newRight := newX + player.Width/2
	newTop := newY + player.Height/2

	oldLeft := player.Xpos - player.Width/2
	oldRight := player.Xpos + player.Width/2
	oldBottom := player.Ypos - player.Height/2
	oldTop := player.Ypos + player.Height/2

	//Ser till att vi inte hamnar utanför den förbestämda mappen
	playerTiles := playerFindCurrentTiles(player)
	intersectTiles := findIntersectTiles(newLeft, oldBottom, newRight, oldTop) //Tittar vilka tiles vi står på efter en ändring i x-led
	HasCollidedTile := collisionChecker(intersectTiles)

	if HasCollidedTile.Occupied { //X kolliderade så vi nollställer testar med de nya Y koordinaterna
		if HasCollidedTile.Wall || HasCollidedTile.Player {
			newLeft = oldLeft
			newRight = oldRight
			newX = player.Xpos
		}

	}
	intersectTiles = findIntersectTiles(newLeft, newBottom, newRight, newTop)
	HasCollidedTile = collisionChecker(intersectTiles)
	if HasCollidedTile.Occupied {
		if HasCollidedTile.Wall || HasCollidedTile.Player {
			newTop = oldTop
			newBottom = oldBottom
			newY = player.Ypos
		}
	}
	intersectTiles = findIntersectTiles(newLeft, newBottom, newRight, newTop)
	newTiles := findNewTiles(intersectTiles, playerTiles)
	oldTiles := findNewTiles(playerTiles, newTiles)
	playerUpdate.X = newX
	playerUpdate.Y = newY

	// Sätter korrekta flaggor på spelarens nya tiles.
	for i := 0; i < len(newTiles); i++ {
		gridList[newTiles[i].Index].Occupied = true
		gridList[newTiles[i].Index].Player = true
		gridList[newTiles[i].Index].PlayerId = player.Id
	}

	for i := 0; i < len(oldTiles); i++ {
		if gridList[oldTiles[i].Index].PlayerId != player.Id {
			gridList[oldTiles[i].Index].Occupied = false
			gridList[oldTiles[i].Index].Player = false
			gridList[oldTiles[i].Index].PlayerId = -1
		}

	}
	col.SetPlayer(player.Id, newX, newY, playerUpdate.DirectionX, playerUpdate.DirectionY, "")
	wg.Done()
}
