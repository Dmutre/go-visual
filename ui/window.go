package ui

import (
	"image"
	"image/color"
	"log"

	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/imageutil"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/image/draw"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/mouse"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
)

type Cross struct {
	x    int
	y    int
	size int
	w    int
}

func NewCross(nx int, ny int) *Cross {
	c := Cross{}
	c.x = nx
	c.y = ny
	c.size = 400
	c.w = 150
	return &c
}

func (c *Cross) DrawCross(t screen.Texture) {
	x1 := c.x - c.w/2
	x2 := c.x + c.w/2
	y1 := c.y - c.size/2
	y2 := c.y + c.size/2

	t.Fill(image.Rect(x1, y1, x2, y2), color.RGBA{R: 0xff, A: 0xff}, draw.Src)

	x1 = c.x - c.size/2
	x2 = c.x + c.size/2
	y1 = c.y - c.w/2
	y2 = c.y + c.w/2
	t.Fill(image.Rect(x1, y1, x2, y2), color.RGBA{R: 0xff, A: 0xff}, draw.Src)
}

func (c *Cross) drawCross(pw *Visualizer) {
	x1 := c.x + c.size
	y1 := c.y + c.size/2 + c.w/2
	y2 := c.y + c.size/2 - c.w/2

	pw.w.Fill(image.Rect(c.x, y1, x1, y2), color.RGBA{R: 0xff, A: 0xff}, draw.Src)

	x1 = c.x + c.size/2 + c.w/2
	x2 := c.x + c.size/2 - c.w/2
	y2 = c.y + c.size
	pw.w.Fill(image.Rect(x1, c.y, x2, y2), color.RGBA{R: 0xff, A: 0xff}, draw.Src)
}

func (c *Cross) move(x, y int) {
	c.x = (x - c.size/2)
	c.y = (y - c.size/2)
}

type Visualizer struct {
	Title         string
	Debug         bool
	OnScreenReady func(s screen.Screen)

	w    screen.Window
	tx   chan screen.Texture
	done chan struct{}

	sz  size.Event
	pos image.Rectangle

	initialSize int
	Crosses     []*Cross
}

func (pw *Visualizer) Main() {
	pw.tx = make(chan screen.Texture)
	pw.done = make(chan struct{})
	pw.pos.Max.Y = 200
	pw.pos.Max.X = 200
	driver.Main(pw.run)
}

func (pw *Visualizer) Update(t screen.Texture) {
	pw.tx <- t
}

func (pw *Visualizer) run(s screen.Screen) {
	w, err := s.NewWindow(&screen.NewWindowOptions{
		Title:  pw.Title,
		Width:  800,
		Height: 800,
	})
	w.Send(paint.Event{})
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
	pw.Crosses = []*Cross{{200, 200, 400, 150}}
	pw.initialSize = 400

	events := make(chan interface{})
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

func detectTerminate(e interface{}) bool {
	switch e := e.(type) {
	case lifecycle.Event:
		if e.To == lifecycle.StageDead {
			return true
		}
	case key.Event:
		if e.Code == key.CodeEscape {
			return true
		}
	}
	return false
}

func (pw *Visualizer) handleEvent(e interface{}, t screen.Texture) {
	switch e := e.(type) {

	case size.Event:
		pw.sz = e

	case error:
		log.Printf("ERROR: %s", e)

	case mouse.Event:
		if t == nil {
			switch e.Button {
			case mouse.ButtonLeft:
				if e.Direction == mouse.DirPress {
					pw.MoveAllCrosses(int(e.X), int(e.Y))
					pw.w.Send(paint.Event{})
				}
			}
		}

	case paint.Event:
		if t == nil {
			pw.drawDefaultUI()
		} else {
			pw.w.Scale(pw.sz.Bounds(), t, t.Bounds(), draw.Src, nil)
		}
		pw.w.Publish()
	}
}

func (pw *Visualizer) drawDefaultUI() {
	pw.w.Fill(pw.sz.Bounds(), color.White, draw.Src)
	for _, cross := range pw.Crosses {
		cross.drawCross(pw)
	}
	for _, br := range imageutil.Border(pw.sz.Bounds(), 10) {
		pw.w.Fill(br, color.White, draw.Src)
	}
}

func (pw *Visualizer) MoveAllCrosses(x, y int) {
	for _, cross := range pw.Crosses {
		cross.move(x, y)
	}
}

func (pw *Visualizer) AddCross(x, y int) {
	newCross := &Cross{x - pw.initialSize/2, y - pw.initialSize/2, pw.initialSize, 100}
	pw.Crosses = append(pw.Crosses, newCross)
}
