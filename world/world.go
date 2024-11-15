package world

import (
	"image"
	"log"
	"path/filepath"

	"code.rocketnine.space/tslocum/gohan"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/lafriks/go-tiled"

	"github.com/diegoxter/blockgame/asset"
	"github.com/diegoxter/blockgame/component"
)

var World = &GameWorld{
	CamScale:     1,
	CamMoving:    true,
	PlayerWidth:  16,
	PlayerHeight: 16,
	TileImages:   make(map[uint32]*ebiten.Image),
	ResetGame:    true,
}

type GameWorld struct {
	Player gohan.Entity

	ScreenW, ScreenH int

	DisableEsc bool

	Debug int

	GameStarted      bool
	GameStartedTicks int
	GameOver         bool

	MessageVisible  bool
	MessageTicks    int
	MessageDuration int
	MessageUpdated  bool
	MessageText     string

	PlayerX, PlayerY float64

	CamX, CamY float64
	CamScale   float64
	CamMoving  bool

	PlayerWidth  float64
	PlayerHeight float64

	Map           *tiled.Map
	ObjectGroups  []*tiled.ObjectGroup
	HazardRects   []image.Rectangle
	CreepRects    []image.Rectangle
	CreepEntities []gohan.Entity
	WallEntities  []gohan.Entity
	WallRects     []image.Rectangle
	WallNames     []string

	NativeResolution bool

	BrokenPieceA, BrokenPieceB gohan.Entity

	TileImages map[uint32]*ebiten.Image

	ResetGame bool
}

func TileToGameCoords(x, y int) (float64, float64) {
	return float64(x) * 16, float64(y) * 16
}

func Reset() {
	for _, e := range gohan.AllEntities() {
		e.Remove()
	}
	World.Player = 0

	World.ObjectGroups = nil
	World.HazardRects = nil
	World.CreepRects = nil
	World.CreepEntities = nil
	World.WallEntities = nil
	World.WallRects = nil
	World.WallNames = nil

	World.MessageVisible = false
}

func LoadMap(filePath string) {
	// Parse .tmx file.
	m, err := tiled.LoadFile(filePath, tiled.WithFileSystem(asset.FS))
	if err != nil {
		log.Fatalf("error parsing world: %+v", err)
	}

	// Load tileset.

	tileset := m.Tilesets[0]

	imgPath := filepath.Join("image/map/", tileset.Image.Source)
	f, err := asset.FS.Open(filepath.FromSlash(imgPath))
	if err != nil {
		log.Fatalf("error loading tileset: %+v", err)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		panic(err)
	}
	tilesetImg := ebiten.NewImageFromImage(img)

	// Load tiles.
	for i := uint32(0); i < uint32(tileset.TileCount); i++ {
		rect := tileset.GetTileRect(i)
		World.TileImages[i+tileset.FirstGID] = tilesetImg.SubImage(rect).(*ebiten.Image)
	}

	createTileEntity := func(t *tiled.LayerTile, x int, y int) gohan.Entity {
		tileX, tileY := TileToGameCoords(x, y)

		mapTile := gohan.NewEntity()
		mapTile.AddComponent(&component.Position{
			X: tileX,
			Y: tileY,
		})

		sprite := &component.Sprite{
			Image:          World.TileImages[t.Tileset.FirstGID+t.ID],
			HorizontalFlip: t.HorizontalFlip,
			VerticalFlip:   t.VerticalFlip,
			DiagonalFlip:   t.DiagonalFlip,
		}
		mapTile.AddComponent(sprite)

		return mapTile
	}

	var t *tiled.LayerTile
	for _, layer := range m.Layers {
		for y := 0; y < m.Height; y++ {
			for x := 0; x < m.Width; x++ {
				t = layer.Tiles[y*m.Width+x]

				if t == nil || t.Nil || t.ID == 0 {
					continue // No tile at this position.
				}
				log.Println(t.ID)
				// tileObj := m.Tilesets[0].Tiles[t.ID-1]

				// if solid := tileObj.Properties.GetBool("solid"); solid {
				// 	log.Println(t.ID)
				// }

				tileImg := World.TileImages[t.Tileset.FirstGID+t.ID]
				if tileImg == nil {
					continue
				}
				createTileEntity(t, x, y)
			}
		}
	}

	// Load ObjectGroups.
	var objects []*tiled.ObjectGroup
	var loadObjects func(grp *tiled.Group)
	loadObjects = func(grp *tiled.Group) {
		for _, subGrp := range grp.Groups {
			loadObjects(subGrp)
		}
		for _, objGrp := range grp.ObjectGroups {
			objects = append(objects, objGrp)
		}
	}

	for _, grp := range m.Groups {
		loadObjects(grp)
	}

	for _, objGrp := range m.ObjectGroups {
		objects = append(objects, objGrp)
	}

	World.Map = m
	World.ObjectGroups = objects

	for _, grp := range World.ObjectGroups {
		if grp.Name == "WALLS" {
			for _, obj := range grp.Objects {
				World.WallRects = append(World.WallRects, ObjectToRect(obj))
			}
		} else if grp.Name == "HAZARDS" {
			for _, obj := range grp.Objects {
				r := ObjectToRect(obj)
				r.Min.Y += 32
				r.Max.Y += 32
				World.HazardRects = append(World.HazardRects, r)
			}
		}
	}
}

func ObjectToRect(o *tiled.Object) image.Rectangle {
	x, y, w, h := int(o.X), int(o.Y), int(o.Width), int(o.Height)
	y -= 32
	return image.Rect(x, y, x+w, y+h)
}

func LevelCoordinatesToScreen(x, y float64) (float64, float64) {
	return (x - World.CamX) * World.CamScale, (y - World.CamY) * World.CamScale
}

func StartGame() {
	if World.GameStarted {
		return
	}
	World.GameStarted = true
}
