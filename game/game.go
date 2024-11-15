package game

import (
	"fmt"
	"image/color"
	"os"
	"sync"

	"code.rocketnine.space/tslocum/gohan"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"github.com/diegoxter/blockgame/asset"
	"github.com/diegoxter/blockgame/component"
	"github.com/diegoxter/blockgame/entity"
	"github.com/diegoxter/blockgame/system"
	"github.com/diegoxter/blockgame/world"
)

type game struct {
	w, h int

	op *ebiten.DrawImageOptions

	movementSystem *system.MovementSystem
	renderSystem   *system.RenderSystem

	addedSystems bool

	sync.Mutex
}

// NewGame returns a new isometric demo game.
func NewGame() (*game, error) {
	g := &game{
		op:           &ebiten.DrawImageOptions{},
	}

	err := g.loadAssets()
	if err != nil {
		return nil, err
	}

	gohan.Preallocate(30000)

	return g, nil
}

// func (g *game) tileToGameCoords(x, y int) (float64, float64) {
// 	return float64(x) * 32, float64(y) * 32
// }

func (g *game) changeMap(filePath string) {
	world.LoadMap(filePath)

	if world.World.Player == 0 {
		world.World.Player = entity.NewPlayer()
	}

	const playerStartOffset = 128
	const camStartOffset = 480

	w := float64(world.World.Map.Width * world.World.Map.TileWidth)
	h := float64(world.World.Map.Height * world.World.Map.TileHeight)

	world.World.Player.With(func(position *component.Position) {
		position.X, position.Y = w/2, h-playerStartOffset
	})

	world.World.CamX, world.World.CamY = 0, h-camStartOffset
}

// Layout is called when the game's layout changes.
func (g *game) Layout(w, h int) (int, int) {
	// return screenWidth, screenHeight
	if !world.World.NativeResolution {
		w, h = 640, 480
	}
	if w != g.w || h != g.h {
		world.World.ScreenW, world.World.ScreenH = w, h
		g.w, g.h = w, h
	}
	return g.w, g.h
}

func (g *game) Update() error {
	if ebiten.IsWindowBeingClosed() {
		g.Exit()
		return nil
	}

	if world.World.ResetGame {
		world.Reset()

    g.changeMap("image/map/map.tmx")

		if !g.addedSystems {
			g.addSystems()

			if world.World.Debug == 0 {
				//  asset.SoundTitleMusic.Play()
			}

			g.addedSystems = true // TODO
		}

		// rand.Seed(time.Now().UnixNano())

		world.World.ResetGame = false
		world.World.GameOver = false
	}

	err := gohan.Update()
	if err != nil {
		return err
	}
	return nil
}

func (g *game) Draw(screen *ebiten.Image) {
	// screen.Fill(color.NRGBA{0x00, 0x40, 0x80, 0xff})
	//screen.DrawImage(tilesetImg, nil)
	err := gohan.Draw(screen)
	if err != nil {
		panic(err)
	}
	ebitenutil.DebugPrint(screen, fmt.Sprintf("player x %2.0f player y %2.0f", world.World.PlayerX ,  world.World.PlayerY))
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%v ", world.World.GameOver), 0 , 30)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("player %v", world.World.Player), 0, 20)
}

func (g *game) addSystems() {
	// Handle input.
	g.movementSystem = system.NewMovementSystem()
	g.renderSystem = system.NewRenderSystem()

	gohan.AddSystem(system.NewPlayerMoveSystem(world.World.Player, g.movementSystem))
	gohan.AddSystem(g.movementSystem)
	gohan.AddSystem(system.NewCameraSystem())
	gohan.AddSystem(system.NewRailSystem())
	gohan.AddSystem(g.renderSystem)
	gohan.AddSystem(system.NewProfileSystem(world.World.Player))
}

func (g *game) loadAssets() error {
	asset.ImgWhiteSquare.Fill(color.White)
	return nil
}

func (g *game) Exit() {
	os.Exit(0)
}