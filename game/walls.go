package game

type Wall struct {
	Id     int
	X      float64
	Y      float64
	Width  float64
	Height float64
}

type WallCollection struct {
	Walls map[int]Wall
}

func NewWall(Id int, X, Y, Width, Height float64) Wall {
	NewWall := Wall{Id, X, Y, Width, Height}
	return NewWall
}

func MakeWallCollection() WallCollection {
	return WallCollection{Walls: map[int]Wall{}}
}

func (col WallCollection) SetWall(Id int, X, Y, Width, Height float64) {
	col.Walls[Id] = Wall{Id: Id, X: X, Y: Y, Width: Width, Height: Height}
}

func (col WallCollection) GetWall(id int) (Wall, bool) {
	wall, ok := col.Walls[id]
	return wall, ok
}
