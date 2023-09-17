package main

import (
	"gioui.org/widget"
	"log"
	"os"

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

	var resumeCapture widget.Clickable

	resumeCapture.Clicked()

	for {
		e := <-w.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)

			layout.Flex{
				Axis: layout.Vertical,
			}.Layout(gtx,
				layout.Rigid(
					func(gtx layout.Context) layout.Dimensions {
						btn := material.Button(th, &resumeCapture, "Start")
						return btn.Layout(gtx)
					},
				),
			)

			e.Frame(gtx.Ops)
		}
	}
}
