package main

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	myApp := app.NewWithID("com.richo.kalender.selamatan")
	myWindow := myApp.NewWindow("Kalkulator Selamatan")

	// Input teks dengan panduan format
	inputTgl := widget.NewEntry()
	inputTgl.SetPlaceHolder("Contoh: 24-12-2024")

	listContainer := container.NewVBox()
	scroll := container.NewVScroll(listContainer)
	scroll.SetMinSize(fyne.NewSize(0, 400))

	btnHitung := widget.NewButton("HITUNG JADWAL", func() {
		// Validasi input tanggal
		t, err := time.Parse("02-01-2006", inputTgl.Text)
		if err != nil {
			widget.ShowError(fmt.Errorf("Format Salah! Gunakan Tgl-Bln-Thn"), myWindow)
			return
		}
		updateList(t, listContainer)
	})

	content := container.NewBorder(
		container.NewVBox(
			widget.NewLabelWithStyle("KALKULATOR SELAMATAN", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			widget.NewLabel("Masukkan Tanggal Meninggal:"),
			inputTgl,
			btnHitung,
			widget.NewSeparator(),
		),
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
		{"Geblag", 0}, {"Nelung Dino (3 Hari)", 2}, {"Mitung Dino (7 Hari)", 6},
		{"Matang Puluh (40 Hari)", 39}, {"Nyatus (100 Hari)", 99},
		{"Pendhak I (1 Tahun)", 353}, {"Pendhak II (2 Tahun)", 707}, {"Nyewu (1000 Hari)", 999},
	}

	for _, item := range items {
		target := t.AddDate(0, 0, item.hari)
		info := fmt.Sprintf("%s: %s", item.nama, target.Format("02-01-2006"))
		c.Add(widget.NewCard("", "", widget.NewLabel(info)))
	}
	c.Refresh()
}
