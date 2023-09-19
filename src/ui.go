package main

import (
	"fmt"
	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/google/gopacket"
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

	resumeButton := widget.Clickable{}
	captureButtonLabel := "Resume Capture"

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

			if resumeButton.Clicked() {
				fmt.Println("clicked", captureButtonLabel)
				isPaused = !isPaused
				if isPaused {
					captureButtonLabel = "Resume Capture"
				} else {
					captureButtonLabel = "Pause Capture"
				}

				handleSelectDevice("\\Device\\NPF_Loopback")
			}

			buttonLayouts = append(buttonLayouts, func(gtx layout.Context) layout.Dimensions {
				btn := material.Button(th, &resumeButton, captureButtonLabel)
				return btn.Layout(gtx)
			})

			// Convert buttonLayouts to []layout.FlexChild
			var flexChildren []layout.FlexChild
			for _, btnLayout := range buttonLayouts {
				flexChildren = append(flexChildren, layout.Rigid(btnLayout))
			}

			yourDataStructureUI(gtx, th)

			flexer.Layout(gtx, flexChildren...)

			e.Frame(gtx.Ops)
		}
	}
}

func yourDataStructureUI(gtx layout.Context, th *material.Theme) layout.Dimensions {
	// Define the layout for displaying the data structure
	// Iterate through the map and display keys as titles and packets underneath
	var flexChildren []layout.FlexChild

	for key, packets := range globalPacketsMap {
		// Create a title label for each key
		title := material.H6(th, key)
		flexChildren = append(flexChildren, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return title.Layout(gtx)
		}))

		// Create a list to display packets associated with the key
		packetList := layout.List{
			Axis: layout.Vertical,
		}
		for _, packet := range packets {
			packetLabel := material.Body1(th, fmt.Sprintf("Packet: %v", packet))
			flexChildren = append(flexChildren, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return packetLabel.Layout(gtx)
			}))
		}
		flexChildren = append(flexChildren, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return packetList.Layout(gtx, len(packets), func(gtx layout.Context, index int) layout.Dimensions {
				return layout.Flex{}.Layout(gtx)
			})
		}))
	}

	// Create a scrollable list of widgets
	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx, flexChildren...)
}

var globalPacketsMap map[string][]gopacket.Packet

func updatePackageView(appsPakketList map[string][]gopacket.Packet) {
	fmt.Println("----------------------------------------")
	fmt.Println("somany apps: ", len(appsPakketList))
	for appName, Pakket := range appsPakketList {
		fmt.Println(appName, ": ", len(Pakket))
	}

	globalPacketsMap = appsPakketList
}
