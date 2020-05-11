package game

import (
	"fmt"
	"sync"
)

type Bullet struct {
	Id         int
	PlayerId   int
	DirectionX float64
	DirectionY float64
	Xpos       float64
	Ypos       float64
}

var BulletChannelCollection = make(map[int]chan float64)

var BulletCollectionMutex sync.Mutex

type BulletCollection struct {
	Bullets map[int]Bullet
}

func NewBullet(Id int, PlayerId int, Xpos, Ypos float64, DirectionX float64, DirectionY float64) Bullet {
	NewBullet := Bullet{Id, PlayerId, DirectionX, DirectionY, Xpos, Ypos}
	return NewBullet
}

func (col BulletCollection) GetBullet(Id int) (Bullet, bool) {
	BulletCollectionMutex.Lock()
	bullet, ok := col.Bullets[Id]
	BulletCollectionMutex.Unlock()
	return bullet, ok
}

func (col BulletCollection) SetBullet(Id int, PlayerId int, Xpos, Ypos float64, DirectionX float64, DirectionY float64) {
	BulletCollectionMutex.Lock()
	col.Bullets[Id] = Bullet{Id: Id, Xpos: Xpos, Ypos: Ypos, DirectionX: DirectionX, DirectionY: DirectionY, PlayerId: PlayerId}
	BulletCollectionMutex.Unlock()
}

// a working implementation of bulletupdate, however, i dont know to support arbitrary size of array.
func (col BulletCollection) BulletUpdate(delta float64) {
	var BulletSpeed float64 = 4
	for _, v := range col.Bullets {
		v.Xpos += v.DirectionX * delta * BulletSpeed
		v.Ypos += v.DirectionY * delta * BulletSpeed
		fmt.Printf("xpos = %g, Ypos = %g \n", v.Xpos, v.Ypos)
		col.SetBullet(v.Id, v.PlayerId, v.Xpos, v.Ypos, v.DirectionX, v.DirectionY)
	}

}

func (col BulletCollection) RemoveBullet(Id int) {
	BulletCollectionMutex.Lock()
	delete(col.Bullets, Id)
	BulletCollectionMutex.Unlock()
}

func NewBulletCollection() BulletCollection {
	return BulletCollection{Bullets: map[int]Bullet{}}
}

//Direction kommer kunna vara 0<= till <= 364
// vi måste ta ett ut så att det blir en najs beräkning på våra x och y poskoordinater.
// vi gör detta genom sin och cos. om vi kollar på enhetscirkeln så förstår man.
