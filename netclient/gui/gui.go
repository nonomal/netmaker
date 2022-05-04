package gui

import (
	"embed"
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/gravitl/netmaker/netclient/functions"
	"github.com/gravitl/netmaker/netclient/gui/components"
	"github.com/gravitl/netmaker/netclient/gui/components/views"
)

//go:embed nm-logo-sm.png
var logoContent embed.FS

func run(networks []string) error {
	a := app.New()
	window := a.NewWindow("Netclient")

	img, err := logoContent.ReadFile("nm-logo.png")
	if err != nil {
		return err
	}

	window.SetIcon(&fyne.StaticResource{StaticName: "Netmaker logo", StaticContent: img})
	window.Resize(fyne.NewSize(600, 400))

	networkView := container.NewVScroll(views.GetNetworksView(networks))
	networkView.SetMinSize(fyne.NewSize(400, 300))
	views.SetView(views.Networks, networkView)

	netDetailsViews := container.NewVScroll(views.GetSingleNetworkView(""))
	netDetailsViews.SetMinSize(fyne.NewSize(400, 300))
	views.SetView(views.NetDetails, netDetailsViews)
	window.SetFixedSize(true)

	toolbar := container.NewCenter(widget.NewToolbar(
		components.NewToolbarLabelButton("Networks", theme.HomeIcon(), func() {
			views.ShowView(views.Networks)
			views.ClearNotification()
		}, components.Blue_color),
		components.NewToolbarLabelButton("Join new", theme.ContentAddIcon(), func() {
			views.ShowView(views.Join)
		}, components.Gravitl_color),
		components.NewToolbarLabelButton("Uninstall", theme.ErrorIcon(), func() {
			confirmView := views.GetConfirmation("Confirm Netclient uninstall?", func() {
				views.ShowView(views.Networks)
			}, func() {
				views.LoadingNotify()
				err := functions.Uninstall()
				time.Sleep(time.Second >> 1)
				if err != nil {
					views.ErrorNotify("Failed to uninstall: \n" + err.Error())
				} else {
					views.SuccessNotify("Uninstalled Netclient!")
				}
				time.Sleep(time.Second >> 1)
				views.ShowView(views.Networks)
			})
			views.RefreshComponent(views.Confirm, confirmView)
			views.ShowView(views.Confirm)
			// TODO:
			// - call uninstall
			// - Refresh networks view when finished
		}, components.Red_color),
	))

	joinView := views.GetJoinView()
	views.SetView(views.Join, joinView)

	confirmView := views.GetConfirmation("", func() {}, func() {})
	views.SetView(views.Confirm, confirmView)

	views.ShowView(views.Networks)

	initialNotification := views.GenerateNotification("", color.Transparent)
	views.SetView(views.Notify, initialNotification)

	views.CurrentContent = container.NewVBox()

	views.CurrentContent.Add(container.NewGridWithRows(
		1,
		toolbar,
	))
	views.CurrentContent.Add(views.GetView(views.Networks))
	views.CurrentContent.Add(views.GetView(views.NetDetails))
	views.CurrentContent.Add(views.GetView(views.Notify))
	views.CurrentContent.Add(views.GetView(views.Join))

	window.SetContent(views.CurrentContent)
	window.ShowAndRun()
	return nil
}
