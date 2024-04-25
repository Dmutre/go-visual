package lang

import (
	"bufio"
	"io"
	"strings"

	"github.com/Dmutre/go-visual/painter"
)

type Parser struct{}

func complexOperation(operation painter.OperationFunc) painter.Operation {
	if operation == nil {
		return nil
	}
	return painter.OperationFunc(operation)
}

func (p *Parser) CommandParser(commandName string, args []string) painter.Operation {
	switch commandName {
	case "white":
		return painter.OperationFunc(painter.WhiteFill)
	case "green":
		return painter.OperationFunc(painter.GreenFill)
	case "bgrect":
		return complexOperation(painter.DrawRectangle(args))
	case "figure":
		return complexOperation(painter.Figure(args))
	case "move":
		return complexOperation(painter.Move(args))
	case "update":
		return painter.UpdateOp
	case "reset":
		return painter.OperationFunc(painter.Reset)
	}
	return nil
}

func (p *Parser) Parse(in io.Reader) ([]painter.Operation, error) {
	var res []painter.Operation

	scanner := bufio.NewScanner(in)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		commandLine := scanner.Text()

		sliced := strings.Split(commandLine, " ")
		args := sliced[1:]

		command := p.CommandParser(sliced[0], args)
		if command != nil {
			res = append(res, command)
		}
	}
	return res, nil
}
