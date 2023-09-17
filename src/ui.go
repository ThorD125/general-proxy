package main

import (
	"fmt"
	"gioui.org/widget"
	"log"
	"os"
	"strconv"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/widget/material"
)

func uii() {
	go func() {
		w := app.NewWindow(
			app.Title("Ghostly's Packet Editor"),
		)
		err := run(w)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}
func run(w *app.Window) error {
	th := material.NewTheme()
	var ops op.Ops

	var buttons []widget.Clickable
	buttonLabels := []string{
		"Button 1",
		"Button 2",
		"Button 3",
		"Button 4",
		"Button 5",
		"Button 6",
		"Button 7",
		"Button 8",
		"Button 9",
		"Button 10",
	}

	for range buttonLabels {
		buttons = append(buttons, widget.Clickable{})
	}

	for {
		e := <-w.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)

			var buttonLayouts []layout.Widget
			for i := range buttonLabels {
				i := i // Create a local variable to capture the current value of i
				for buttons[i].Clicked() {
					fmt.Println("Button", i+1, "clicked")
				}

				buttonLayouts = append(buttonLayouts, func(gtx layout.Context) layout.Dimensions {
					btn := material.Button(th, &buttons[i], strconv.Itoa(i+1))
					return btn.Layout(gtx)
				})
			}

			var flexer = layout.Flex{
				Axis: layout.Vertical,
			}

			// Convert buttonLayouts to []layout.FlexChild
			var flexChildren []layout.FlexChild
			for _, btnLayout := range buttonLayouts {
				flexChildren = append(flexChildren, layout.Rigid(btnLayout))
			}

			flexer.Layout(gtx, flexChildren...)

			e.Frame(gtx.Ops)
		}
	}
}
