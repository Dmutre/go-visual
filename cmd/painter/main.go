package main

import (
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/Dmutre/go-visual/painter"
	"github.com/Dmutre/go-visual/painter/lang"
	"github.com/Dmutre/go-visual/ui"
)

func main() {
    var (
        pv ui.Visualizer // Візуалізатор створює вікно та малює у ньому.

        // Потрібні для частини 2.
        opLoop painter.Loop // Цикл обробки команд.
        parser lang.Parser  // Парсер команд.
    )

    //pv.Debug = true
    pv.Title = "Simple painter"

    pv.OnScreenReady = opLoop.Start
    opLoop.Receiver = &pv

    go func() {
        http.Handle("/", lang.HttpHandler(&opLoop, &parser))
        _ = http.ListenAndServe("localhost:17000", nil)
    }()

    if os.Getenv("CI") == "true" {
        // If in CI, start the event loop and the tests
        go func() {
            // Wait for the event loop to start
            time.Sleep(time.Second)

            // Run the tests
            os.Args = []string{"go", "test", "./painter/loop_test.go", "./painter/lang/parser_test.go"}
            testing.Main(func(pat, str string) (bool, error) { return true, nil },
                []testing.InternalTest{},
                []testing.InternalBenchmark{},
                []testing.InternalExample{})
            
            // Stop the event loop when the tests are done
            opLoop.StopAndWait()
        }()

        // Start the event loop
        // pv.Main()
    } else {
        // If not in CI, just start the event loop
        pv.Main()
        opLoop.StopAndWait()
    }
}