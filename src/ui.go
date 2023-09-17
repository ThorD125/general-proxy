package main

import (
	"fmt"
	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"log"
	"os"
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
	buttonLabels := selectAbleDevices()

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
				eweee := i // Create a local variable to capture the current value of eweee
				for buttons[eweee].Clicked() {
					fmt.Println("clicked", buttonLabels[eweee])
				}

				buttonLayouts = append(buttonLayouts, func(gtx layout.Context) layout.Dimensions {
					btn := material.Button(th, &buttons[eweee], buttonLabels[eweee])
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
