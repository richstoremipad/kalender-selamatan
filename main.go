package main

import (
	"fmt"
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// Fungsi untuk menghitung Weton Jawa sederhana
func getWeton(t time.Time) string {
	pasaran := []string{"Legi", "Pahing", "Pon", "Wage", "Kliwon"}
	hari := []string{"Minggu", "Senin", "Selasa", "Rabu", "Kamis", "Jumat", "Sabtu"}
	
	// Referensi 1 Januari 1970 adalah Kamis Wage
	refDate := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	diff := t.Sub(refDate).Hours() / 24
	
	idxPasaran := (int(diff) + 3) % 5 // +3 karena 1 Jan 1970 adalah Wage (idx 3)
	if idxPasaran < 0 { idxPasaran += 5 }
	
	return hari[t.Weekday()] + " " + pasaran[idxPasaran]
}

func main() {
	myApp := app.NewWithID("com.richo.kalender.selamatan")
	myWindow := myApp.NewWindow("Kalkulator Selamatan Jawa")

	// Header Gradasi (Mirip Gambar 1955.jpg)
	headerText := canvas.NewText("Kalkulator Selamatan Jawa", color.White)
	headerText.TextStyle = fyne.TextStyle{Bold: true}
	headerText.TextSize = 20
	headerBg := canvas.NewLinearGradient(color.NRGBA{0, 150, 136, 255}, color.NRGBA{33, 150, 243, 255}, 0)
	header := container.NewStack(headerBg, container.NewPadded(headerText))

	inputTgl := widget.NewEntry()
	inputTgl.SetPlaceHolder("Contoh: 24-12-2024")

	listContainer := container.NewVBox()
	scroll := container.NewVScroll(listContainer)
	scroll.SetMinSize(fyne.NewSize(0, 500))

	btnHitung := widget.NewButton("HITUNG JADWAL", func() {
		t, err := time.Parse("02-01-2006", inputTgl.Text)
		if err != nil {
			dialog.ShowError(fmt.Errorf("Gunakan format: Tgl-Bln-Thn"), myWindow)
			return
		}
		
		listContainer.Objects = nil
		items := []struct {
			nama string
			hari int
		}{
			{"Geblag", 0}, {"Nelung Dino (3 Hari)", 2}, {"Mitung Dino (7 Hari)", 6},
			{"Matang Puluh (40 Hari)", 39}, {"Nyatus (100 Hari)", 99},
			{"Pendhak I (1 Tahun)", 353}, {"Pendhak II (2 Tahun)", 707}, {"Nyewu (1000 Hari)", 999},
		}

		now := time.Now().Truncate(24 * time.Hour)
		for _, item := range items {
			target := t.AddDate(0, 0, item.hari)
			weton := getWeton(target)
			
			// Penentuan Badge Status (Mirip Gambar 1955.jpg)
			statusText := ""
			statusColor := color.NRGBA{100, 100, 100, 255}
			if target.Before(now) {
				statusText = "Sudah Lewat"
				statusColor = color.NRGBA{76, 175, 80, 255} // Hijau
			} else if target.Equal(now) {
				statusText = "HARI INI!"
				statusColor = color.NRGBA{244, 67, 54, 255} // Merah
			}

			badge := canvas.NewRectangle(statusColor)
			badge.SetMinSize(fyne.NewSize(80, 20))
			badgeTxt := canvas.NewText(statusText, color.White)
			badgeTxt.TextSize = 10

			cardContent := container.NewVBox(
				container.NewHBox(
					widget.NewLabelWithStyle(item.nama, fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
					container.NewStack(badge, container.NewCenter(badgeTxt)),
				),
				widget.NewLabel(fmt.Sprintf("%s, %s", weton, target.Format("02-01-2006"))),
			)
			listContainer.Add(widget.NewCard("", "", cardContent))
		}
		listContainer.Refresh()
	})

	mainContent := container.NewBorder(
		container.NewVBox(header, widget.NewLabel("Input Tanggal Meninggal:"), inputTgl, btnHitung, widget.NewSeparator()),
		widget.NewLabelWithStyle("Matur Nuwun - Code by Richo", fyne.TextAlignCenter, fyne.TextStyle{Italic: true}),
		nil, nil,
		scroll,
	)

	myWindow.SetContent(mainContent)
	myWindow.Resize(fyne.NewSize(400, 650))
	myWindow.ShowAndRun()
}

