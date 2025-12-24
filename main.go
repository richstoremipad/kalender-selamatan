package main

import (
	"fmt"
	"image/color"
	"math"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// ==========================================
// 1. DATA & WARNA
// ==========================================

var (
	HariIndo  = []string{"Minggu", "Senin", "Selasa", "Rabu", "Kamis", "Jumat", "Sabtu"}
	Pasaran   = []string{"Legi", "Pahing", "Pon", "Wage", "Kliwon"}
	BulanIndo = []string{"", "Januari", "Februari", "Maret", "April", "Mei", "Juni", "Juli", "Agustus", "September", "Oktober", "November", "Desember"}
	BulanJawa = []string{"", "Suro", "Sapar", "Mulud", "Bakda Mulud", "Jumadil Awal", "Jumadil Akhir", "Rajeb", "Ruwah", "Poso", "Sawal", "Sela", "Besar"}

	ColorBgDark     = color.NRGBA{R: 30, G: 33, B: 40, A: 150} // Transparansi agar background batik terlihat
	ColorCardBg     = color.NRGBA{R: 45, G: 48, B: 55, A: 240}
	ColorHeaderTop  = color.NRGBA{R: 40, G: 180, B: 160, A: 255}
	ColorHeaderBot  = color.NRGBA{R: 50, G: 80, B: 160, A: 255}
	ColorTextWhite  = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	ColorTextGrey   = color.NRGBA{R: 180, G: 180, B: 180, A: 255}
	ColorBadgeGreen = color.NRGBA{R: 46, G: 125, B: 50, A: 255}
	ColorBadgeRed   = color.NRGBA{R: 198, G: 40, B: 40, A: 255}
	ColorBadgeBlue  = color.NRGBA{R: 21, G: 101, B: 192, A: 255}
)

// ==========================================
// 2. LOGIKA PERHITUNGAN
// ==========================================

func dateToJDN(t time.Time) int {
	a := (14 - int(t.Month())) / 12
	y := t.Year() + 4800 - a
	m := int(t.Month()) + 12*a - 3
	return t.Day() + (153*m+2)/5 + 365*y + y/4 - y/100 + y/400 - 32045
}

func getJavaneseDate(t time.Time) string {
	jd := dateToJDN(t)
	l := jd - 1948440 + 10632 + 1 
	n := (l - 1) / 10631
	l = l - 10631*n + 354
	j := (int)((10985 - l) / 5316) * (int)((50 * l) / 17719) + (int)(l / 5670) * (int)((43 * l) / 15238)
	l = l - (int)((30 - j) / 15) * (int)((17719 * j) / 50) - (int)(j / 16) * (int)((15238 * j) / 43) + 29
	hm := (int)(24 * l) / 709
	hd := l - (int)(709 * hm) / 24
	if hm <= 0 || hm >= len(BulanJawa) { return "Unknown" }
	return fmt.Sprintf("%d %s", hd, BulanJawa[hm])
}

func formatWeton(t time.Time) string {
	return fmt.Sprintf("%s %s, %s", HariIndo[t.Weekday()], Pasaran[dateToJDN(t)%5], getJavaneseDate(t))
}

func createCard(title, dateStr, wetonStr string, statusType, diffDays int) fyne.CanvasObject {
	var badgeColor color.Color
	var badgeTxt string
	switch statusType {
	case 1: badgeColor, badgeTxt = ColorBadgeGreen, fmt.Sprintf("✓ Lewat (%d hari)", int(math.Abs(float64(diffDays))))
	case 2: badgeColor, badgeTxt = ColorBadgeRed, "🔔 HARI INI!"
	default: badgeColor, badgeTxt = ColorBadgeBlue, fmt.Sprintf("⏳ %d Hari Lagi", diffDays)
	}
	lblTitle := canvas.NewText(title, ColorTextWhite); lblTitle.TextStyle.Bold = true
	rightCont := container.NewVBox(canvas.NewText(dateStr, ColorTextWhite), canvas.NewText(wetonStr, ColorTextGrey))
	badge := container.NewStack(canvas.NewRectangle(badgeColor), container.NewPadded(canvas.NewText(badgeTxt, ColorTextWhite)))
	bg := canvas.NewRectangle(ColorCardBg); bg.CornerRadius = 10
	return container.NewStack(bg, container.NewPadded(container.NewVBox(container.NewBorder(nil, nil, lblTitle, rightCont), container.NewHBox(badge))))
}

// ==========================================
// 3. MAIN APP
// ==========================================

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Kalkulator Selamatan Jawa")
	myWindow.Resize(fyne.NewSize(400, 750))

	// Header
	gradient := canvas.NewHorizontalGradient(ColorHeaderTop, ColorHeaderBot)
	header := container.NewStack(gradient, container.NewCenter(canvas.NewText("Kalkulator Selamatan Jawa", ColorTextWhite)))

	// Input & Result
	inputEntry := widget.NewEntry(); inputEntry.PlaceHolder = "Contoh: 01/12/2024"
	resultBox := container.NewVBox()
	
	btnCalc := widget.NewButtonWithIcon("Hitung", theme.ConfirmIcon(), func() {
		t, err := time.Parse("02/01/2006", inputEntry.Text)
		if err != nil { return }
		resultBox.Objects = nil
		events := []struct { N string; O int }{
			{"Geblag", 0}, {"Nelung", 2}, {"Mitung", 6}, {"Matang", 39},
			{"Nyatus", 99}, {"Pendhak I", 353}, {"Pendhak II", 707}, {"Nyewu", 999},
		}
		now := time.Now().Truncate(24 * time.Hour)
		for _, e := range events {
			target := t.AddDate(0, 0, e.O).Truncate(24 * time.Hour)
			diff := int(target.Sub(now).Hours() / 24)
			status := 3
			if diff < 0 { status = 1 } else if diff == 0 { status = 2 }
			resultBox.Add(createCard(e.N, target.Format("02-01-2006"), formatWeton(target), status, diff))
		}
	})

	// Penempatan Background Gambar
	bgImage := canvas.NewImageFromFilesystem("background.png")
	bgImage.FillMode = canvas.ImageFillStretch

	mainContent := container.NewBorder(
		container.NewVBox(header, container.NewPadded(container.NewVBox(widget.NewLabel("Tanggal Geblag (DD/MM/YYYY):"), inputEntry, btnCalc))),
		container.NewPadded(widget.NewLabelWithStyle("Matur Nuwun - Code by Richo", fyne.TextAlignCenter, fyne.TextStyle{Italic: true})),
		nil, nil, container.NewVScroll(container.NewPadded(resultBox)),
	)

	// Menyusun tumpukan: Gambar Batik -> Layer Gelap -> Konten Aplikasi
	myWindow.SetContent(container.NewStack(bgImage, canvas.NewRectangle(ColorBgDark), mainContent))
	myWindow.ShowAndRun()
}

