package game

import (
	"fmt"
	"testing"
)

func TestBulletCollection(t *testing.T) {
	// Should add a new bullet of id '0'
	col := NewBulletCollection()
	col.SetBullet(3, 0, 0, 0, 0, 0)
	if len(col.Bullets) != 1 {
		fmt.Printf("no bullet")
		t.Fail()
	}
	// Should overwrite  the old bullet of id '0' /*
	col.SetBullet(0, 25, 25, 0, 0, 0)
	if len(col.Bullets) != 2 {
		fmt.Printf("SetBullet2 fail")
		t.Fail()
	}

	// Should have corrent new data
	bullet, ok := col.GetBullet(0)
	if !ok || bullet.Xpos != 25 {
		fmt.Printf("GetBullet fail")
		t.Fail()
	}
	// Should add a new bullet of id '1'
	col.SetBullet(1, 0, 0, 0, 0, 0)
	if len(col.Bullets) != 3 {
		fmt.Printf("SetBullet3 fail")
		t.Fail()
	}

	// Should remove bullet of id '0'
	col.RemoveBullet(0)
	if len(col.Bullets) != 2 {
		fmt.Printf("RemoveBullet fail")
		t.Fail()
	}

	// Should remove bullet of id '1'
	col.RemoveBullet(1)
	if len(col.Bullets) != 1 {
		fmt.Printf("RemoveBullet2 fail")
		t.Fail()
	}
	col.RemoveBullet(3)
	if len(col.Bullets) != 0 {
		fmt.Printf("RemoveBullet3 fail")
		t.Fail()
	}
}

func TestBulletUpdate(t *testing.T) {
	// Should add a new bullet of id '0'
	col := NewBulletCollection()
	col.SetBullet(3, 2, 1, 1, 2, 2)
	if len(col.Bullets) != 1 {
		fmt.Printf("no bullet")
		t.Fail()
	}

	col.BulletUpdate(1)
	col.BulletUpdate(1)
	col.BulletUpdate(1)
	bullet, ok := col.GetBullet(3)
	if !ok {
		panic("couldnt get player")
	}

	if bullet.Xpos != 25 {
		t.Errorf("Xpos = %g", bullet.Xpos)
	}

	if bullet.Ypos != 25 {
		t.Errorf("Ypos = %g", bullet.Ypos)
	}
}
