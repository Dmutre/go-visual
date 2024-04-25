package lang

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/Dmutre/go-visual/painter"
)

func TestParser_CommandParser(t *testing.T) {
	p := &Parser{}

	tests := []struct {
		name         string
		commandName  string
		args         []string
		expectedFunc painter.OperationFunc
	}{
		{
			name:         "WhiteFill",
			commandName:  "white",
			expectedFunc: painter.OperationFunc(painter.WhiteFill),
		},
		{
			name:         "GreenFill",
			commandName:  "green",
			expectedFunc: painter.OperationFunc(painter.GreenFill),
		},
		{
			name:         "DrawRectangle",
			commandName:  "bgrect",
			args:         []string{"0.25", "0.25", "30", "40"},
			expectedFunc: painter.DrawRectangle([]string{"0.25", "0.25", "30", "40"}),
		},
		{
			name:         "Figure",
			commandName:  "figure",
			args:         []string{"0.25", "0.25"},
			expectedFunc: painter.Figure([]string{"0.25", "0.25"}),
		},
		{
			name:         "Move",
			commandName:  "move",
			args:         []string{"0.5", "0.5"},
			expectedFunc: painter.Move([]string{"0.5", "0.5"}),
		},
		{
			name:         "Reset",
			commandName:  "reset",
			expectedFunc: painter.OperationFunc(painter.Reset),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fmt.Println(test.commandName)
			result := p.CommandParser(test.commandName, test.args)
			if result == nil {
				t.Errorf("Expected non-nil result for command %s", test.commandName)
			}
			if func1, ok := result.(painter.OperationFunc); ok {
				if !painterOperationFuncEquals(func1, test.expectedFunc) {
					t.Errorf("Expected function %v, got %v", test.expectedFunc, func1)
				}
			}
		})
	}
}

func TestParser_Parse(t *testing.T) {
	p := &Parser{}

	tests := []struct {
		name          string
		input         string
		expectedCount int
	}{
		{
			name:          "SingleCommand",
			input:         "white\n",
			expectedCount: 1,
		},
		{
			name:          "MultipleCommands",
			input:         "white\ngreen\n",
			expectedCount: 2,
		},
		{
			name:          "CommandWithArgs",
			input:         "bgrect 0.5 0.4 30 40\n",
			expectedCount: 1,
		},
		{
			name:          "EmptyInput",
			input:         "",
			expectedCount: 0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := strings.NewReader(test.input)
			ops, err := p.Parse(r)
			if err != nil {
				t.Fatalf("Error parsing input: %v", err)
			}
			if len(ops) != test.expectedCount {
				t.Errorf("Expected %d operations, got %d", test.expectedCount, len(ops))
			}
		})
	}
}

// painterOperationFuncEquals compares two OperationFunc for equality by comparing their underlying functions.
func painterOperationFuncEquals(f1, f2 painter.OperationFunc) bool {
	return reflect.ValueOf(f1).Pointer() == reflect.ValueOf(f2).Pointer()
}
