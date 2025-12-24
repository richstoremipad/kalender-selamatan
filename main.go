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
// 1. LOGIKA MATEMATIKA & KALENDER (FIXED)
// ==========================================

var (
	HariIndo  = []string{"Minggu", "Senin", "Selasa", "Rabu", "Kamis", "Jumat", "Sabtu"}
	Pasaran   = []string{"Legi", "Pahing", "Pon", "Wage", "Kliwon"}
	BulanIndo = []string{"", "Januari", "Februari", "Maret", "April", "Mei", "Juni", "Juli", "Agustus", "September", "Oktober", "November", "Desember"}
	BulanJawa = []string{"", "Suro", "Sapar", "Mulud", "Bakda Mulud", "Jumadil Awal", "Jumadil Akhir", "Rajeb", "Ruwah", "Poso", "Sawal", "Sela", "Besar"}
)

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

	namaBulanJawa := ""
	if hm > 0 && hm < len(BulanJawa) {
		namaBulanJawa = BulanJawa[hm]
	} else {
		namaBulanJawa = "Unknown"
	}

	return fmt.Sprintf("%d %s", hd, namaBulanJawa)
}

func formatWeton(t time.Time) string {
	hari := HariIndo[t.Weekday()]
	jd := dateToJDN(t)
	pasaranIdx := jd % 5
	pasaran := Pasaran[pasaranIdx]
	jawaDate := getJavaneseDate(t)
	return fmt.Sprintf("%s %s, %s", hari, pasaran, jawaDate)
}

func formatIndoDate(t time.Time) string {
	return fmt.Sprintf("%d %s %d", t.Day(), BulanIndo[t.Month()], t.Year())
}

// ==========================================
// 2. KOMPONEN UI CUSTOM
// ==========================================

var (
	ColorBgDark     = color.NRGBA{R: 30, G: 33, B: 40, A: 200} // Transparansi sedikit agar bg terlihat
	ColorCardBg     = color.NRGBA{R: 45, G: 48, B: 55, A: 240}
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
	lblTitle.TextSize = 16
	lblTitle.TextStyle = fyne.TextStyle{Bold: true}
	lblSub := canvas.NewText(subTitle, ColorTextGrey)
	lblSub.TextSize = 12
	leftCont := container.NewVBox(lblTitle, lblSub)

	lblDate := canvas.NewText(dateStr, ColorTextWhite)
	lblDate.Alignment = fyne.TextAlignTrailing
	lblDate.TextSize = 14
	lblDate.TextStyle = fyne.TextStyle{Bold: true}
	lblWeton := canvas.NewText(wetonStr, ColorTextGrey)
	lblWeton.Alignment = fyne.TextAlignTrailing
	lblWeton.TextSize = 11
	rightCont := container.NewVBox(lblDate, lblWeton)

	topRow := container.NewBorder(nil, nil, leftCont, rightCont)
	lblBadge := canvas.NewText(badgeTextStr, ColorTextWhite)
	lblBadge.TextSize = 11
	lblBadge.TextStyle = fyne.TextStyle{Bold: true}
	badgeBg := canvas.NewRectangle(badgeColor)
	badgeBg.CornerRadius = 12
	badgeCont := container.NewStack(badgeBg, container.NewPadded(lblBadge))
	botRow := container.NewHBox(badgeCont)

	content := container.NewVBox(topRow, container.NewPadded(botRow))
	bg := canvas.NewRectangle(ColorCardBg)
	bg.CornerRadius = 10
	return container.NewStack(bg, container.NewPadded(content))
}

// ==========================================
// 3. MAIN APP
// ==========================================

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Kalkulator Selamatan Jawa")
	myWindow.Resize(fyne.NewSize(400, 750))

	gradient := canvas.NewHorizontalGradient(ColorHeaderTop, ColorHeaderBot)
	headerTitle := canvas.NewText("Kalkulator Selamatan Jawa", ColorTextWhite)
	headerTitle.TextStyle = fyne.TextStyle{Bold: true}
	headerTitle.TextSize = 18
	headerTitle.Alignment = fyne.TextAlignCenter
	
	headerIcon := canvas.NewImageFromResource(theme.InfoIcon())
	headerIcon.SetMinSize(fyne.NewSize(30,30))

	headerStack := container.NewStack(
		gradient,
		container.NewPadded(container.NewVBox(
			layout.NewSpacer(),
			container.NewHBox(layout.NewSpacer(), headerIcon, headerTitle, layout.NewSpacer()),
			layout.NewSpacer(),
		)),
	)
	headerContainer := container.NewVBox(headerStack)

	inputLabel := canvas.NewText("Input Tanggal / Geblag (DD/MM/YYYY)", ColorTextGrey)
	inputLabel.TextSize = 12
	inputEntry := widget.NewEntry()
	inputEntry.PlaceHolder = "Contoh: 01/12/2024"
	inputEntry.TextStyle = fyne.TextStyle{Monospace: true}
	btnCalc := widget.NewButton("Hitung", nil)
	inputRow := container.NewBorder(nil, nil, nil, btnCalc, inputEntry)
	inputCardBg := canvas.NewRectangle(ColorCardBg)
	inputCardBg.CornerRadius = 8
	inputSection := container.NewStack(inputCardBg, container.NewPadded(container.NewVBox(inputLabel, inputRow)))

	resultBox := container.NewVBox()
	scrollArea := container.NewVScroll(container.NewPadded(resultBox))

	noteText := "Notes: Perhitungan ini saya buat berdasarkan rumus jawa dari kitab yang pernah saya pelajari..."
	lblNote := widget.NewLabel(noteText)
	lblNote.Wrapping = fyne.TextWrapWord
	lblNote.TextStyle = fyne.TextStyle{Italic: true}
	lblCredit := canvas.NewText("Matur Nuwun - Code by Richo", ColorTextGrey)
	lblCredit.Alignment = fyne.TextAlignCenter
	lblCredit.TextSize = 10

	footer := container.NewVBox(lblNote, lblCredit)
	footerCardBg := canvas.NewRectangle(ColorCardBg)
	footerCardBg.CornerRadius = 8
	footerSection := container.NewStack(footerCardBg, container.NewPadded(footer))

	btnCalc.OnTapped = func() {
		dateStr := inputEntry.Text
		t, err := time.Parse("02/01/2006", dateStr)
		if err != nil {
			resultBox.Objects = []fyne.CanvasObject{widget.NewLabel("Format Salah! Gunakan DD/MM/YYYY")}
			resultBox.Refresh()
			return
		}
		resultBox.Objects = nil 
		events := []struct { Name string; Sub string; Offset int }{
			{"Geblag", "Hari H", 0}, {"Nelung", "3 Hari", 2}, {"Mitung", "7 Hari", 6},
			{"Matang", "40 Hari", 39}, {"Nyatus", "100 Hari", 99}, {"Pendhak I", "1 Tahun", 353},
			{"Pendhak II", "2 Tahun", 707}, {"Nyewu", "1000 Hari", 999},
		}
		now := time.Now()
		now = time.Date(now.Year(), now.Month(), now.Day(), 0,0,0,0, now.Location())
		for _, e := range events {
			targetDate := t.AddDate(0, 0, e.Offset)
			targetDate = time.Date(targetDate.Year(), targetDate.Month(), targetDate.Day(), 0,0,0,0, targetDate.Location())
			diff := int(targetDate.Sub(now).Hours() / 24)
			status := 3 
			if diff < 0 { status = 1 } else if diff == 0 { status = 2 }
			card := createCard(e.Name, e.Sub, formatIndoDate(targetDate), formatWeton(targetDate), status, diff)
			resultBox.Add(card)
			resultBox.Add(layout.NewSpacer()) 
		}
		resultBox.Refresh()
	}

	// --- Layout Utama (Dengan Background) ---
	bgApp := canvas.NewRectangle(ColorBgDark)
	
	// Menambahkan gambar background
	bgImage := canvas.NewImageFromFilesystem("background.png")
	bgImage.FillMode = canvas.ImageFillStretch 

	mainContent := container.NewBorder(
		container.NewVBox(headerContainer, container.NewPadded(inputSection)),
		container.NewPadded(footerSection),
		nil, nil, 
		scrollArea, 
	)

	// Menyusun tumpukan: Background paling bawah, baru UI di atasnya
	finalLayout := container.NewStack(bgImage, bgApp, mainContent)

	myWindow.SetContent(finalLayout)
	myWindow.ShowAndRun()
}

// Konversi tanggal Masehi ke Format Jawa
func getJavaneseDate(t time.Time) string {
	// 1. Hitung JDN
	jd := dateToJDN(t)

	// 2. Koreksi Tanggal Jawa (Hijriah Calendar Approximation)
	// Rumus aritmatika tabular sering meleset 1 hari tergantung kriteria (Rukyatul Hilal vs Hisab).
	// Berdasarkan feedback (1 Des 2024 = 29 Jumadil Awal), kita perlu penyesuaian +1 pada offset perhitungan.
	
	// Algoritma konversi JD ke Hijri/Jawa Tabular:
	l := jd - 1948440 + 10632 + 1 // (+1 added for calibration to Waktu.id/Mataram standard)
	n := (l - 1) / 10631
	l = l - 10631*n + 354
	j := (int)((10985 - l) / 5316) * (int)((50 * l) / 17719) + (int)(l / 5670) * (int)((43 * l) / 15238)
	l = l - (int)((30 - j) / 15) * (int)((17719 * j) / 50) - (int)(j / 16) * (int)((15238 * j) / 43) + 29
	
	hm := (int)(24 * l) / 709
	hd := l - (int)(709 * hm) / 24

	// Validasi index bulan
	namaBulanJawa := ""
	if hm > 0 && hm < len(BulanJawa) {
		namaBulanJawa = BulanJawa[hm]
	} else {
		namaBulanJawa = "Unknown"
	}

	return fmt.Sprintf("%d %s", hd, namaBulanJawa)
}

func formatWeton(t time.Time) string {
	// Hari
	hari := HariIndo[t.Weekday()]
	
	// Pasaran (Menggunakan JDN Modulo 5)
	// Rumus: JDN % 5. 
	// 0=Legi, 1=Pahing, 2=Pon, 3=Wage, 4=Kliwon
	jd := dateToJDN(t)
	pasaranIdx := jd % 5
	pasaran := Pasaran[pasaranIdx]

	// Tanggal & Bulan Jawa
	jawaDate := getJavaneseDate(t)

	return fmt.Sprintf("%s %s, %s", hari, pasaran, jawaDate)
}

func formatIndoDate(t time.Time) string {
	return fmt.Sprintf("%d %s %d", t.Day(), BulanIndo[t.Month()], t.Year())
}

// ==========================================
// 2. KOMPONEN UI CUSTOM
// ==========================================

// Warna Palette
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

// createCard membuat tampilan kartu
func createCard(title, subTitle, dateStr, wetonStr string, statusType int, diffDays int) fyne.CanvasObject {
	var badgeColor color.Color
	var badgeTextStr string
	
	switch statusType {
	case 1: // Lewat
		badgeColor = ColorBadgeGreen
		badgeTextStr = fmt.Sprintf("✓ Sudah Lewat (%d hari)", int(math.Abs(float64(diffDays))))
	case 2: // Hari Ini
		badgeColor = ColorBadgeRed
		badgeTextStr = "🔔 HARI INI!"
	case 3: // Belum
		badgeColor = ColorBadgeBlue
		badgeTextStr = fmt.Sprintf("⏳ %d Hari Lagi", diffDays)
	}

	lblTitle := canvas.NewText(title, ColorTextWhite)
	lblTitle.TextSize = 16
	lblTitle.TextStyle = fyne.TextStyle{Bold: true}

	lblSub := canvas.NewText(subTitle, ColorTextGrey)
	lblSub.TextSize = 12

	leftCont := container.NewVBox(lblTitle, lblSub)

	lblDate := canvas.NewText(dateStr, ColorTextWhite)
	lblDate.Alignment = fyne.TextAlignTrailing
	lblDate.TextSize = 14
	lblDate.TextStyle = fyne.TextStyle{Bold: true}

	lblWeton := canvas.NewText(wetonStr, ColorTextGrey)
	lblWeton.Alignment = fyne.TextAlignTrailing
	lblWeton.TextSize = 11

	rightCont := container.NewVBox(lblDate, lblWeton)

	topRow := container.NewBorder(nil, nil, leftCont, rightCont)

	lblBadge := canvas.NewText(badgeTextStr, ColorTextWhite)
	lblBadge.TextSize = 11
	lblBadge.TextStyle = fyne.TextStyle{Bold: true}
	
	badgeBg := canvas.NewRectangle(badgeColor)
	badgeBg.CornerRadius = 12
	
	badgeCont := container.NewStack(
		badgeBg,
		container.NewPadded(lblBadge),
	)
	
	botRow := container.NewHBox(badgeCont)

	content := container.NewVBox(topRow, container.NewPadded(botRow))

	bg := canvas.NewRectangle(ColorCardBg)
	bg.CornerRadius = 10

	return container.NewStack(bg, container.NewPadded(content))
}

// ==========================================
// 3. MAIN APP
// ==========================================

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Kalkulator Selamatan Jawa")
	myWindow.Resize(fyne.NewSize(400, 750))

	// --- Header Gradient ---
	gradient := canvas.NewHorizontalGradient(ColorHeaderTop, ColorHeaderBot)
	headerTitle := canvas.NewText("Kalkulator Selamatan Jawa", ColorTextWhite)
	headerTitle.TextStyle = fyne.TextStyle{Bold: true}
	headerTitle.TextSize = 18
	headerTitle.Alignment = fyne.TextAlignCenter
	
	headerIcon := canvas.NewImageFromResource(theme.InfoIcon())
	headerIcon.SetMinSize(fyne.NewSize(30,30))

	headerStack := container.NewStack(
		gradient,
		container.NewPadded(container.NewVBox(
			layout.NewSpacer(),
			container.NewHBox(layout.NewSpacer(), headerIcon, headerTitle, layout.NewSpacer()),
			layout.NewSpacer(),
		)),
	)
	headerContainer := container.NewVBox(headerStack)


	// --- Input Section ---
	inputLabel := canvas.NewText("Input Tanggal / Geblag (DD/MM/YYYY)", ColorTextGrey)
	inputLabel.TextSize = 12

	inputEntry := widget.NewEntry()
	inputEntry.PlaceHolder = "Contoh: 01/12/2024"
	inputEntry.TextStyle = fyne.TextStyle{Monospace: true}

	btnCalc := widget.NewButton("Hitung", nil)

	inputRow := container.NewBorder(nil, nil, nil, btnCalc, inputEntry)
	
	inputCardBg := canvas.NewRectangle(ColorCardBg)
	inputCardBg.CornerRadius = 8
	inputSection := container.NewStack(
		inputCardBg,
		container.NewPadded(container.NewVBox(inputLabel, inputRow)),
	)


	// --- Result Container ---
	resultBox := container.NewVBox()
	scrollArea := container.NewVScroll(container.NewPadded(resultBox))

	// --- Footer ---
	noteText := "Notes: Perhitungan ini saya buat berdasarkan rumus jawa dari kitab yang pernah saya pelajari yaitu lusarlu (3), tusarmo (7), masarmo (40), rosarmo (100), patsarpat (Pendhak 1), rosarji (Pendhak 2), nemsarmo (1000). Adapun perbedaan dari hitungan anda mungkin hanya 1/2 hari saja yang berarti tidak masalah. Wallahu A'lam Bishawab"
	lblNote := widget.NewLabel(noteText)
	lblNote.Wrapping = fyne.TextWrapWord
	lblNote.TextStyle = fyne.TextStyle{Italic: true}
	
	lblCredit := canvas.NewText("Matur Nuwun - Code by Richo", ColorTextGrey)
	lblCredit.Alignment = fyne.TextAlignCenter
	lblCredit.TextSize = 10

	footer := container.NewVBox(lblNote, lblCredit)
	footerCardBg := canvas.NewRectangle(ColorCardBg)
	footerCardBg.CornerRadius = 8
	footerSection := container.NewStack(
		footerCardBg,
		container.NewPadded(footer),
	)


	// --- Logic Calculation ---
	btnCalc.OnTapped = func() {
		dateStr := inputEntry.Text
		layoutFormat := "02/01/2006"
		t, err := time.Parse(layoutFormat, dateStr)
		if err != nil {
			resultBox.Objects = []fyne.CanvasObject{
				widget.NewLabel("Format Salah! Gunakan DD/MM/YYYY"),
			}
			resultBox.Refresh()
			return
		}

		resultBox.Objects = nil // Clear previous

		type Event struct {
			Name   string
			Sub    string
			Offset int
		}

		events := []Event{
			{"Geblag", "Hari H", 0},
			{"Nelung", "3 Hari", 2},
			{"Mitung", "7 Hari", 6},
			{"Matang", "40 Hari", 39},
			{"Nyatus", "100 Hari", 99},
			{"Pendhak I", "1 Tahun", 353},
			{"Pendhak II", "2 Tahun", 707},
			{"Nyewu", "1000 Hari", 999},
		}

		now := time.Now()
		// Normalize now to midnight
		now = time.Date(now.Year(), now.Month(), now.Day(), 0,0,0,0, now.Location())

		for _, e := range events {
			targetDate := t.AddDate(0, 0, e.Offset)
			targetDate = time.Date(targetDate.Year(), targetDate.Month(), targetDate.Day(), 0,0,0,0, targetDate.Location())

			diff := int(targetDate.Sub(now).Hours() / 24)

			status := 3 // Future
			if diff < 0 {
				status = 1 // Lewat
			} else if diff == 0 {
				status = 2 // Hari ini
			}

			card := createCard(
				e.Name,
				e.Sub,
				formatIndoDate(targetDate),
				formatWeton(targetDate),
				status,
				diff,
			)
			
			resultBox.Add(card)
			resultBox.Add(layout.NewSpacer()) 
		}
		resultBox.Refresh()
	}

	// --- Layout Utama ---
	bgApp := canvas.NewRectangle(ColorBgDark)

	mainContent := container.NewBorder(
		container.NewVBox(headerContainer, container.NewPadded(inputSection)),
		container.NewPadded(footerSection),
		nil, nil, 
		scrollArea, 
	)

	finalLayout := container.NewStack(bgApp, mainContent)

	myWindow.SetContent(finalLayout)
	myWindow.ShowAndRun()

	// Load gambar background
bgImage := canvas.NewImageFromFilesystem("background.png")
bgImage.FillMode = canvas.ImageFillStretch // Agar gambar menutupi seluruh layar

// Gunakan container.NewStack agar konten berada di atas background
contentWithBg := container.NewStack(bgImage, mainContent)

myWindow.SetContent(contentWithBg)

}
