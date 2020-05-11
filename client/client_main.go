package client

import (
	"dogmatix/network"
	"image"
	"math"
	"os"

	_ "image/png"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

// Huvudfunktionen för klienten, kallas från main
func ClientMain() {
	pixelgl.Run(run)
}

type Bullet struct {
	time         float64
	direction    pixel.Vec
	bulletVector pixel.Vec
}

var (
	playerId      int
	playerVector  pixel.Vec
	bulletSprites *pixel.Sprite
	bulletList    []Bullet
	direction     int
	cardinal      string
)

func remove(s []Bullet, i int) []Bullet {
	if len(s) == 0 {
		return s
	}
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}

/* Checks if the input vector is out of bounds */
func outOfBounds(vec pixel.Vec, cfg pixelgl.WindowConfig) bool {
	if vec.X >= cfg.Bounds.Max.X || vec.Y >= cfg.Bounds.Max.Y {
		return true
	}
	if vec.X < cfg.Bounds.Min.X || vec.Y < cfg.Bounds.Min.Y {
		return true
	}
	return false
}

/* Loads image with path input */
func loadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}

/* Returns an array of all the sprite objects */
/* Param:
   deltaX, the length in px between objects
   deltaY, the difference in px between objects
*/
func spriteDivider(deltaX float64, deltaY float64, image pixel.Picture) []pixel.Rect {
	var newFrames []pixel.Rect
	for x := image.Bounds().Min.X; x < image.Bounds().Max.X; x += deltaX {
		for y := image.Bounds().Min.Y; y < image.Bounds().Max.Y; y += deltaY {
			newFrames = append(newFrames, pixel.R(x, y, x+deltaX, y+deltaY))
		}
	}
	return newFrames
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Ett spel",
		Bounds: pixel.R(0, 0, 1024, 768),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	playerVector = pixel.Vec.Add(playerVector, win.Bounds().Center())

	hpbar, err := loadPicture("client/hpbar.png")
	if err != nil {
		panic(err)
	}

	hpbarFilled, err := loadPicture("client/hpbarFilled.png")
	if err != nil {
		panic(err)
	}

	// Utkommenterat då det inte används för tillfället
	// spritesheet, err := loadPicture("client/tempsprite.png")
	// if err != nil {
	// 	panic(err)
	// }

	playerSprite, err := loadPicture("client/player.png")
	if err != nil {
		panic(err)
	}

	bulletPic, err := loadPicture("client/bullet-sprite.png")
	if err != nil {
		panic(err)
	}

	// var spriteFrames = spriteDivider(64, 64, spritesheet) // Utkommenterat då det inte används för tillfället

	var playerFrame = spriteDivider(64, 64, playerSprite)
	var hpbarFilledSprite = spriteDivider(64.4, 80, hpbarFilled) /* Syntax is (deltaX, deltaY, Image) */

	hpbarSprite := pixel.NewSprite(hpbar, hpbar.Bounds())
	bulletSprite := pixel.NewSprite(bulletPic, bulletPic.Bounds())

	for !win.Closed() {
		isShooting := false

		win.Clear(colornames.Forestgreen)
		if win.JustPressed(pixelgl.MouseButtonLeft) {
			newBullet := Bullet{time: 0.0, direction: pixel.Vec.Unit(pixel.Vec.To(playerVector, win.MousePosition())), bulletVector: pixel.Vec.Add(pixel.ZV, playerVector)}
			bulletList = append(bulletList, newBullet)
			isShooting = true
		}
		if win.Pressed(pixelgl.KeyUp) {
			if !outOfBounds(playerVector.Add(pixel.V(0, 2)), cfg) {
				playerVector = playerVector.Add(pixel.V(0, 2))
			}

		}
		if win.Pressed(pixelgl.KeyRight) {
			if !outOfBounds(playerVector.Add(pixel.V(2, 0)), cfg) {
				playerVector = playerVector.Add(pixel.V(2, 0))
			}
		}

		if win.Pressed(pixelgl.KeyDown) {
			if !outOfBounds(playerVector.Add(pixel.V(0, -2)), cfg) {
				playerVector = playerVector.Add(pixel.V(0, -2))
			}
		}

		if win.Pressed(pixelgl.KeyLeft) {
			if !outOfBounds(playerVector.Add(pixel.V(-2, 0)), cfg) {
				playerVector = playerVector.Add(pixel.V(-2, 0))
			}
		}

		playerRotationMatrix := pixel.Matrix.Moved(pixel.IM, playerVector)
		angle := math.Pi + pixel.Vec.Angle(pixel.Vec.Unit(pixel.Vec.To(playerVector, win.MousePosition())))
		playerRotationMatrix = pixel.Matrix.Rotated(playerRotationMatrix, playerVector, angle)

		player := pixel.NewSprite(playerSprite, playerFrame[0])
		player.Draw(win, playerRotationMatrix)

		network.SendPlayerUpdate(network.PlayerUpdate{X: playerVector.X, Y: playerVector.Y /*Direction: angle*/, IsShooting: isShooting})

		for i := 0; i < len(bulletList); i++ {
			if bulletList[i].time < 100.0 {
				tempMatrix := pixel.Matrix.Moved(pixel.IM, bulletList[i].bulletVector)
				tempMatrix = tempMatrix.Scaled(bulletList[i].bulletVector, 0.03)

				tempMatrix = pixel.Matrix.Rotated(tempMatrix, bulletList[i].bulletVector, math.Pi+pixel.Vec.Angle(bulletList[i].direction))

				bulletSprite.Draw(win, tempMatrix)
				bulletList[i].bulletVector = bulletList[i].bulletVector.Add(pixel.Vec.Scaled(bulletList[i].direction, 8))
				bulletList[i].time++
			}

			// TODO remove
			if bulletList[i].time >= 100.0 {
				bulletList = remove(bulletList, i)
				if i <= len(bulletList) {
					i = 0
				}
			}
		}

		hpMatrix := pixel.Matrix.Moved(pixel.IM, pixel.Vec.Add(pixel.V(0, 0), pixel.V(498, 35)))
		hpFilledMatrix := pixel.Matrix.Moved(pixel.IM, pixel.Vec.Add(pixel.V(0, 0), pixel.V(500, 20)))
		hpbarSprite.Draw(win, hpMatrix)
		/* the loop should depend on hp, so 10 should be swapped out to hp instead. */
		for i := 0; i < 3; i++ {

			temp := (float64(i) * 64.4) - 290
			hpbarSpriteSelect := pixel.NewSprite(hpbarFilled, hpbarFilledSprite[i])
			hpbarSpriteSelect.Draw(win, hpFilledMatrix.Moved(pixel.V(temp, 0)))
		}
		win.Update()
	}
}

func main() {
	// TODO should be reliably set from server
	//playerId = rand.Intn(1007688)
	//pixelgl.Run(run)

}
