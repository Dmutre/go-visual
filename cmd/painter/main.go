package main

import (
	"net/http"
	"os"
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