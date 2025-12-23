package main

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"fyne.io/x/fyne/widget/calendar"
)

func main() {
	// Membuat aplikasi dengan ID unik untuk Package Name
	myApp := app.NewWithID("com.richo.kalender.selamatan")
	myWindow := myApp.NewWindow("Kalkulator Selamatan")

	header := widget.NewLabelWithStyle("KALKULATOR SELAMATAN", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	labelTgl := widget.NewLabel("Silakan Pilih Tanggal")
	labelTgl.Alignment = fyne.TextAlignCenter

	// Wadah hasil yang bisa di-scroll
	listContainer := container.NewVBox()
	scroll := container.NewVScroll(listContainer)
	scroll.SetMinSize(fyne.NewSize(0, 400))

	// Tombol untuk membuka kalender
	btnKalender := widget.NewButton("PILIH TANGGAL MENINGGAL", func() {
		d := calendar.NewDatePicker(time.Now(), func(t time.Time) {
			labelTgl.SetText("Tanggal dipilih: " + t.Format("02-01-2006"))
			updateList(t, listContainer)
		})
		d.Show(myWindow)
	})

	content := container.NewBorder(
		container.NewVBox(header, labelTgl, btnKalender, widget.NewSeparator()),
		widget.NewLabel("Matur Nuwun - Richo"),
		nil, nil,
		scroll,
	)

	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(400, 600))
	myWindow.ShowAndRun()
}

func updateList(t time.Time, c *fyne.Container) {
	c.Objects = nil
	items := []struct {
		nama string
		hari int
	}{
		{"Geblag", 0}, {"Nelung Dino (3)", 2}, {"Mitung Dino (7)", 6},
		{"Matang Puluh (40)", 39}, {"Nyatus (100)", 99},
		{"Pendhak I", 353}, {"Pendhak II", 707}, {"Nyewu (1000)", 999},
	}

	for _, item := range items {
		target := t.AddDate(0, 0, item.hari)
		// Menampilkan label peringatan dan tanggal masehi
		info := fmt.Sprintf("%s: %s", item.nama, target.Format("02-01-2006"))
		c.Add(widget.NewCard("", "", widget.NewLabel(info)))
	}
	c.Refresh()
}
