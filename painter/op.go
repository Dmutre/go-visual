package painter

import (
	"fmt"
	"image"
	"image/color"
	"strconv"

	"github.com/Dmutre/go-visual/ui"
	"golang.org/x/exp/shiny/screen"
)

func getCordsByArgs(width int, height int, floatArgs []float64) ([]int, error) {
	if len(floatArgs)%2 != 0 {
		return nil, fmt.Errorf("invalid arg count")
	}

	cords := make([]int, len(floatArgs))

	fWidth := float64(width)
	fHeight := float64(height)

	for index := range floatArgs {
		if index%2 == 0 {
			cords[index] = int(fWidth * floatArgs[index])
		} else {
			cords[index] = int(fHeight * floatArgs[index])
		}
	}

	return cords, nil
}

func convertArgs(args []string) ([]float64, error) {
	parsedValues := make([]float64, len(args))
	for i, str := range args {
		num, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return nil, err
		}
		parsedValues[i] = num
	}
	return parsedValues, nil
}

type Operation interface {
	Do(t screen.Texture, state *State) (ready bool)
}

type OperationList []Operation

func (ol OperationList) Do(t screen.Texture, state *State) (ready bool) {
	for _, o := range ol {
		ready = o.Do(t, state) || ready
	}
	return
}

var UpdateOp = updateOp{}

type updateOp struct{}

func (op updateOp) Do(t screen.Texture, state *State) bool {
	t.Fill(t.Bounds(), state.background, screen.Src)
	t.Fill(image.Rectangle{state.bgRect[0], state.bgRect[1]}, color.Black, screen.Src)
	for _, item := range state.crosses {
		item.DrawCross(t)
	}
	return true
}

type OperationFunc func(t screen.Texture, state *State)

func (f OperationFunc) Do(t screen.Texture, state *State) bool {
	f(t, state)
	return false
}

func WhiteFill(t screen.Texture, state *State) {
	state.background = color.White
}

func GreenFill(t screen.Texture, state *State) {
	state.background = color.RGBA{G: 0xff, A: 0xff}
}

func BlackFill(t screen.Texture, state *State) {
	state.background = color.Black
}

func Reset(t screen.Texture, state *State) {
	state.background = color.Black
	state.bgRect = [2]image.Point{{0, 0}, {0, 0}}
	state.crosses = []*ui.Cross{}
}

func DrawRectangle(args []string) OperationFunc {
	if len(args) != 4 {
		fmt.Println("Wrong amount of arguments to draw a rectangle")
		return nil
	}
	floatArgs, err := convertArgs(args)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return func(t screen.Texture, state *State) {
		cords, err := getCordsByArgs(t.Bounds().Dx(), t.Bounds().Dy(), floatArgs)
		if err == nil && len(cords) == 4 {
			state.bgRect[0] = image.Point{int(cords[0]), int(cords[1])}
			state.bgRect[1] = image.Point{int(cords[2]), int(cords[3])}
		}
	}
}

func Figure(args []string) OperationFunc {
	if len(args) != 2 {
		fmt.Println("Wrong amount of arguments to move figures")
		return nil
	}
	floatArgs, err := convertArgs(args)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return func(t screen.Texture, state *State) {
		cords, err := getCordsByArgs(t.Bounds().Dx(), t.Bounds().Dy(), floatArgs)
		if err == nil && len(cords) == 2 {
			cross := ui.NewCross(cords[0], cords[1])
			state.crosses = append(state.crosses, cross)
		}
	}
}

func Move(args []string) OperationFunc {
	if len(args) != 2 {
		fmt.Println("Wrong amount of arguments to move figures")
		return nil
	}
	floatArgs, err := convertArgs(args)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return func(t screen.Texture, state *State) {
		cords, err := getCordsByArgs(t.Bounds().Dx(), t.Bounds().Dy(), floatArgs)
		if err == nil && len(cords) == 2 {
			cross := ui.NewCross(cords[0], cords[1])
			state.crosses = []*ui.Cross{cross}
		}
	}
}
