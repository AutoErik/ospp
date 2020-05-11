package server


import (
//	"fmt"
)


type Bullet struct {
	Xpos      float64
	Ypos      float64

	Direction float32



	//För att lyckas lägga till kollision behöver vi en bestämd höjd/bredd på objektet

}

func NewBullet(Xpos , Ypos float64, Direction float32) Bullet {
	NewBullet := Bullet{Xpos, Ypos, Direction}
	return NewBullet
}

// a working implementation of bulletupdate, however, i dont know to support arbitrary size of array.
func BulletUpdate(bullets  [10]*Bullet, delta float64) {

	// a multiplier that increases the bullets speed. Feel free to change this value if its too high or low.
	var BULLETSPEED float64 = 4 

	// delta sets it so that the bullets travel a set  amount per second independent of server lag.
	// i think the BULLETSPEED might screw things up a little bit. we will have to try and see. im not sure. its hard to think.
	for i:=0; bullets[i] != nil; i++ {
		bullets[i].Xpos += DegreesToX( float64(bullets[i].Direction) ) * delta * BULLETSPEED
		bullets[i].Ypos += DegreesToY( float64(bullets[i].Direction) ) * delta * BULLETSPEED
	} 

//	fmt.Printf("delta: %g \n", delta)
	
}


func DetectBulletCollision (testingBullet *Bullet, bullets  [10]*Bullet)  bool  {

	collision := false

	
	for i:= 0 ; bullets[i] != nil; i++ {
		if(testingBullet != bullets[i]) {
			if(testingBullet.Xpos == bullets[i].Xpos && testingBullet.Ypos == bullets[i].Ypos) {
				collision = true
			}
		}

	}
	return collision
}


//Direction kommer kunna vara 0<= till <= 364
// vi måste ta ett ut så att det blir en najs beräkning på våra x och y poskoordinater.
// vi gör detta genom sin och cos. om vi kollar på enhetscirkeln så förstår man.

