package painter

import (
	"image"
	"image/color"
	"sync"

	"github.com/Dmutre/go-visual/ui"
	"golang.org/x/exp/shiny/screen"
)

// Receiver отримує текстуру, яка була підготовлена в результаті виконання команд у циклі подій.
type Receiver interface {
	Update(t screen.Texture)
}

type State struct {
	background color.Color
	bgRect     [2]image.Point
	crosses    []*ui.Cross
}

type Loop struct {
	Receiver Receiver

	buffer screen.Texture

	state State

	mq messageQueue

	stop    chan struct{}
	stopReq bool
}

var size = image.Pt(800, 800)

func (l *Loop) Start(s screen.Screen) {
	l.buffer, _ = s.NewTexture(size)
	l.stop = make(chan struct{})
	l.state = State{
		background: color.White,
		bgRect:     [2]image.Point{{0, 0}, {0, 0}},
		crosses:    []*ui.Cross{},
	}

	go func() {
		for !(l.stopReq && l.mq.empty()) {
			op := l.mq.pull()

			update := op.Do(l.buffer, &l.state)

			if update {
				l.Receiver.Update(l.buffer)
			}
		}
		close(l.stop)
	}()
}

func (l *Loop) Post(op Operation) {
	l.mq.push(op)
}

func (l *Loop) StopAndWait() {
	l.Post(OperationFunc(func(screen.Texture, *State) {
		l.stopReq = true
	}))

	<-l.stop
}

type messageQueue struct {
	ops    []Operation
	mu     sync.Mutex
	signal chan struct{}
}

func (mq *messageQueue) push(op Operation) {
	mq.mu.Lock()
	defer mq.mu.Unlock()

	mq.ops = append(mq.ops, op)

	if mq.signal != nil {
		close(mq.signal)
	}
}

func (mq *messageQueue) pull() Operation {
	mq.mu.Lock()
	defer mq.mu.Unlock()

	if len(mq.ops) == 0 {
		mq.mu.Unlock()

		mq.signal = make(chan struct{})
		<-mq.signal
		mq.signal = nil

		mq.mu.Lock()
	}
	op := mq.ops[0]
	mq.ops[0] = nil
	mq.ops = mq.ops[1:]
	return op
}

func (mq *messageQueue) empty() bool {
	mq.mu.Lock()
	defer mq.mu.Unlock()

	return len(mq.ops) == 0
}
