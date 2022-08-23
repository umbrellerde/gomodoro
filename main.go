package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"time"
)

type ClockMode int8

const (
	Work ClockMode = iota
	ShortPause
	LongPause
)

func (cm ClockMode) String() string {
	return [...]string{"Work", "ShortPause", "LongPause"}[cm]
}

type Clock struct {
	Running   bool
	TimeLeft  time.Duration
	ModeIndex int
}

var ModeOrder = [8]ClockMode{Work, ShortPause, Work, ShortPause, Work, ShortPause, Work, LongPause}
var ModeLength = map[ClockMode]time.Duration{
	Work:       5 * time.Second,
	ShortPause: 5 * time.Second,
	LongPause:  10 * time.Second,
}

var state = Clock{
	Running:  false,
	TimeLeft: ModeLength[Work],
}

func main() {
	a := app.New()
	desk := a.(desktop.App)
	w := a.NewWindow("Hello World")
	w.SetFixedSize(true)
	w.Resize(fyne.NewSize(500, 500))

	timerLabel := widget.NewLabel("Initial Text")
	modeLabel := widget.NewLabel(ModeOrder[state.ModeIndex].String())
	startButton := widget.NewButton("Start", startPressed)
	stopButton := widget.NewButton("Stop", stopPressed)

	menu := fyne.NewMenu("Menu Label",
		fyne.NewMenuItem("Start", startPressed),
		fyne.NewMenuItem("Stop", stopPressed))
	desk.SetSystemTrayMenu(menu)

	go func() {
		for range time.Tick(time.Second) {
			if state.Running {
				print("Running....")
				state.TimeLeft = state.TimeLeft - 1*time.Second
				timerLabel.SetText(fmtDuration(state.TimeLeft))
			}
			if state.TimeLeft < 0 {
				state.ModeIndex = (state.ModeIndex + 1) % len(ModeOrder)
				newMode := ModeOrder[state.ModeIndex]
				state.TimeLeft = ModeLength[newMode]
				modeLabel.SetText(newMode.String())
				timerLabel.SetText(fmtDuration(state.TimeLeft))
				state.Running = false
				a.SendNotification(fyne.NewNotification(fmt.Sprintf("Please start %s -ing", newMode.String()), fmt.Sprintf("You have %s", fmtDuration(state.TimeLeft))))
			}
		}
	}()

	content := container.New(layout.NewVBoxLayout(), timerLabel, modeLabel, startButton, stopButton)

	w.SetContent(content)
	w.ShowAndRun()

}

func startPressed() {
	print("Starting...")
	state.Running = true
}

func stopPressed() {
	print("Stopping")
	state.Running = false
}

func fmtDuration(d time.Duration) string {
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second
	return fmt.Sprintf("%02d:%02d", m, s)
}
