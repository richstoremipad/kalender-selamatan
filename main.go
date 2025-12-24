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
// 1. LOGIKA KALENDER JAWA TERKALIBRASI
// ==========================================

var (
	HariIndo  = []string{"Minggu", "Senin", "Selasa", "Rabu", "Kamis", "Jumat", "Sabtu"}
	Pasaran   = []string{"Legi", "Pahing", "Pon", "Wage", "Kliwon"}
	BulanJawa = []string{"", "Suro", "Sapar", "Mulud", "Bakda Mulud", "Jumadil Awal", "Jumadil Akhir", "Rajeb", "Ruwah", "Poso", "Sawal", "Sela", "Besar"}
)

// dateToJDN menghitung Julian Day Number secara presisi
func dateToJDN(y, m, d int) int {
	if m <= 2 {
		y--
		m += 12
	}
	a := y / 100
	b := 2 - a + (a / 4)
	return int(math.Floor(365.25*float64(y+4716))) + int(math.Floor(30.6001*float64(m+1))) + d + b - 1524
}

func getJavaneseDate(t time.Time) string {
	jd := dateToJDN(t.Year(), int(t.Month()), t.Day())
	
	// Kalibrasi khusus untuk menyamakan dengan sistem Sultan Agungan (Mataram)
	// 1 Desember 2024 (JD 2460646) harus menghasilkan 29 Jumadil Awal
	l := jd - 1948440 + 10632 + 1 
	n := (l - 1) / 10631
	l = l - 10631*n + 354
	j := ((10985 - l) / 5316) * ((50 * l) / 17719) + (l / 5670) * ((43 * l) / 15238)
	l = l - (((30 - j) / 15) * ((17719 * j) / 50)) - ((j / 16) * ((15238 * j) / 43)) + 29
	
	hm := (24 * l) / 709
	hd := l - ((709 * hm) / 24)

	if hm < 1 || hm >= len(BulanJawa) {
		return "Eror"
	}
	return fmt.Sprintf("%d %s", hd, BulanJawa[hm])
}

func formatWeton(t time.Time) string {
	jd := dateToJDN(t.Year(), int(t.Month()), t.Day())
	
	// Hitungan Pasaran: JD % 5
	// Referensi: JD 2460646 (1 Des 2024) % 5 = 1 (Pahing)
	pIdx := jd % 5
	return fmt.Sprintf("%s %s, %s", HariIndo[t.Weekday()], Pasaran[pIdx], getJavaneseDate(t))
}

// ==========================================
// 2. UI & LAYOUT (TETAP SAMA DENGAN SEBELUMNYA)
// ==========================================

var (
	ColorBgDark     = color.NRGBA{R: 30, G: 33, B: 40, A: 255}
	ColorCardBg     = color.NRGBA{R: 45, G: 48, B: 55, A: 255}
	ColorHeaderTop  = color.NRGBA{R: 40, G: 180, B: 160, A: 255}
	ColorHeaderBot  = color.NRGBA{R: 50, G: 80, B: 160, A: 255}
	ColorTextWhite  = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	ColorTextGrey   = color.NRGBA{R: 180, G: 180, B: 180, A: 255}
	ColorBadgeGreen = color.NRGBA{R: 46, G: 125, B: 50, A: 255}
	ColorBadgeRed   = color.NRGBA{R: 198, G: 40, B: 40, A: 255}
	ColorBadgeBlue  = color.NRGBA{R: 21, G: 101, B: 192, A: 255}
)

func createCard(title, subTitle, dateStr, wetonStr string, statusType int, diffDays int) fyne.CanvasObject {
	var badgeColor color.Color
	var badgeTextStr string
	
	switch statusType {
	case 1: 
		badgeColor = ColorBadgeGreen
		badgeTextStr = fmt.Sprintf("✓ Sudah Lewat (%d hari)", int(math.Abs(float64(diffDays))))
	case 2: 
		badgeColor = ColorBadgeRed
		badgeTextStr = "🔔 HARI INI!"
	case 3: 
		badgeColor = ColorBadgeBlue
		badgeTextStr = fmt.Sprintf("⏳ %d Hari Lagi", diffDays)
	}

	lblTitle := canvas.NewText(title, ColorTextWhite)
	lblTitle.TextStyle = fyne.TextStyle{Bold: true}
	lblSub := canvas.NewText(subTitle, ColorTextGrey)
	lblSub.TextSize = 12

	lblDate := canvas.NewText(dateStr, ColorTextWhite)
	lblDate.Alignment = fyne.TextAlignTrailing
	lblDate.TextStyle = fyne.TextStyle{Bold: true}
	lblWeton := canvas.NewText(wetonStr, ColorTextGrey)
	lblWeton.Alignment = fyne.TextAlignTrailing
	lblWeton.TextSize = 11

	badgeBg := canvas.NewRectangle(badgeColor)
	badgeBg.CornerRadius = 12
	lblBadge := canvas.NewText(badgeTextStr, ColorTextWhite)
	lblBadge.TextSize = 11
	badgeCont := container.NewStack(badgeBg, container.NewPadded(lblBadge))
	
	topRow := container.NewBorder(nil, nil, container.NewVBox(lblTitle, lblSub), container.NewVBox(lblDate, lblWeton))
	content := container.NewVBox(topRow, container.NewHBox(badgeCont))

	bg := canvas.NewRectangle(ColorCardBg)
	bg.CornerRadius = 10
	return container.NewStack(bg, container.NewPadded(content))
}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Kalkulator Selamatan Jawa")
	myWindow.Resize(fyne.NewSize(400, 700))

	resultBox := container.NewVBox()
	inputEntry := widget.NewEntry()
	inputEntry.SetPlaceHolder("01/12/2024")

	btnCalc := widget.NewButton("Hitung", func() {
		t, err := time.Parse("02/01/2006", inputEntry.Text)
		if err != nil {
			resultBox.Objects = []fyne.CanvasObject{widget.NewLabel("Format Salah (DD/MM/YYYY)")}
			resultBox.Refresh()
			return
		}
		resultBox.Objects = nil
		now := time.Now()
		now = time.Date(now.Year(), now.Month(), now.Day(), 0,0,0,0, time.Local)

		events := []struct{n, s string; o int}{
			{"Geblag", "Hari H", 0}, {"Nelung", "3 Hari", 2}, {"Mitung", "7 Hari", 6},
			{"Matang", "40 Hari", 39}, {"Nyatus", "100 Hari", 99},
			{"Pendhak I", "1 Tahun", 353}, {"Pendhak II", "2 Tahun", 707}, {"Nyewu", "1000 Hari", 999},
		}

		for _, e := range events {
			target := t.AddDate(0, 0, e.o)
			diff := int(target.Sub(now).Hours() / 24)
			status := 3
			if diff < 0 { status = 1 } else if diff == 0 { status = 2 }

			resultBox.Add(createCard(e.n, e.s, target.Format("02 January 2006"), formatWeton(target), status, diff))
		}
		resultBox.Refresh()
	})

	header := container.NewStack(canvas.NewHorizontalGradient(ColorHeaderTop, ColorHeaderBot), container.NewPadded(widget.NewLabelWithStyle("Kalkulator Selamatan Jawa", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})))
	
	content := container.NewBorder(container.NewVBox(header, container.NewPadded(container.NewVBox(widget.NewLabel("Input Tanggal (DD/MM/YYYY)"), container.NewBorder(nil, nil, nil, btnCalc, inputEntry)))), nil, nil, nil, container.NewVScroll(resultBox))
	
	myWindow.SetContent(container.NewStack(canvas.NewRectangle(ColorBgDark), content))
	myWindow.ShowAndRun()
}

