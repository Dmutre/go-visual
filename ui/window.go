package ui

import (
	"image"
	"image/color"
	"log"

	"golang.org/x/exp/shiny/driver"
	//"golang.org/x/exp/shiny/imageutil"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/image/draw"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/mouse"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
)

type Visualizer struct {
	Title         string
	Debug         bool
	OnScreenReady func(s screen.Screen)

	w    screen.Window
	tx   chan screen.Texture
	done chan struct{}

	sz  size.Event
	pos image.Rectangle

	point func(x, y int) image.Point
}

func (pw *Visualizer) Main() {
	pw.tx = make(chan screen.Texture)
	pw.done = make(chan struct{})
	pw.pos.Max.X = 800
	pw.pos.Max.Y = 800
	pw.point = createPoint()
	driver.Main(pw.run)
}

func (pw *Visualizer) Update(t screen.Texture) {
	pw.tx <- t
}

func (pw *Visualizer) run(s screen.Screen) {
	w, err := s.NewWindow(&screen.NewWindowOptions{
		Title: pw.Title,
		Width: pw.pos.Max.X,
		Height: pw.pos.Max.Y,
	})
	if err != nil {
		log.Fatal("Failed to initialize the app window:", err)
	}
	defer func() {
		w.Release()
		close(pw.done)
	}()

	if pw.OnScreenReady != nil {
		pw.OnScreenReady(s)
	}

	pw.w = w

	events := make(chan any)
	go func() {
		for {
			e := w.NextEvent()
			if pw.Debug {
				log.Printf("new event: %v", e)
			}
			if detectTerminate(e) {
				close(events)
				break
			}
			events <- e
		}
	}()

	var t screen.Texture

	for {
		select {
		case e, ok := <-events:
			if !ok {
				return
			}
			pw.handleEvent(e, t)

		case t = <-pw.tx:
			w.Send(paint.Event{})
		}
	}
}

func detectTerminate(e any) bool {
	switch e := e.(type) {
	case lifecycle.Event:
		if e.To == lifecycle.StageDead {
			return true // Window destroy initiated.
		}
	case key.Event:
		if e.Code == key.CodeEscape {
			return true // Esc pressed.
		}
	}
	return false
}

func createPoint() func(x, y int) image.Point {
	var pointX = 400
	var pointY = 400
	return func(x, y int) image.Point {
			if x != 0 && y != 0 {
					pointX = x
					pointY = y
			}
			return image.Point{X: pointX, Y: pointY}
	}
}

func (pw *Visualizer) handleEvent(e any, t screen.Texture) {
	switch e := e.(type) {
	case mouse.Event:
		if e.Button == mouse.ButtonLeft && e.Direction == mouse.DirPress {
			// Перетворення значень X та Y на цілі числа.
			x := int(e.X)
			y := int(e.Y)
			// Переміщення хрестика до позиції миші.
			pw.moveCrosshair(pw.point(x, y))
		}

	case paint.Event:
		// Малювання контенту вікна.
		if t == nil {
			pw.moveCrosshair(pw.point(0, 0))
		} else {
			// Використання текстури отриманої через виклик Update.
			pw.w.Scale(pw.sz.Bounds(), t, t.Bounds(), draw.Src, nil)
		}
		pw.w.Publish()
	}
}

func (pw *Visualizer) moveCrosshair(pos image.Point) {
	// Визначення нових координат центра хрестика.
	pw.pos.Max.X = pos.X
	pw.pos.Max.Y = pos.Y

	pw.sz.WidthPx = 800;
	pw.sz.HeightPx = 800;

	pw.w.Fill(pw.sz.Bounds(), color.White, draw.Src)

	// Перемалювання вікна з новими координатами хрестика.
	pw.drawCrosshair()
}

func (pw *Visualizer) drawCrosshair() {
	// Розміри хрестика.
	width := 150
	height := 400

	// Визначення позицій вертикального та горизонтального прямокутників.
	verticalRect := image.Rect(pw.pos.Max.X-width/2, pw.pos.Max.Y-height/2, pw.pos.Max.X+width/2, pw.pos.Max.Y+height/2)
	horizontalRect := image.Rect(pw.pos.Max.X-height/2, pw.pos.Max.Y-width/2, pw.pos.Max.X+height/2, pw.pos.Max.Y+width/2)

	// Малювання прямокутників.
	pw.w.Fill(verticalRect, color.RGBA{R: 0xff, A: 0xff}, draw.Src)
	pw.w.Fill(horizontalRect, color.RGBA{R: 0xff, A: 0xff}, draw.Src)
}