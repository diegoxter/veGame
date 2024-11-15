package asset

import (
	"embed"
	"image"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"

)

//go:embed image
var FS embed.FS

var (
	ImgPlayer        = LoadImage("image/bat.png")
)

var ImgWhiteSquare = ebiten.NewImage(16, 16)

func LoadImage(p string) *ebiten.Image {
	f, err := FS.Open(p)
	if err != nil {
		log.Println("failed to open file: %s", err.Error())
	}
	defer f.Close()

	baseImg, _, err := image.Decode(f)
	if err != nil {
		log.Println("failed to decode image: %s", err.Error())
	}

	return ebiten.NewImageFromImage(baseImg)
}