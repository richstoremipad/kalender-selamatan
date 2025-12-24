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
	"fyne.io/fyne/v2/dialog" // Import dialog untuk Kalender
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// ==========================================
// 1. LOGIKA MATEMATIKA & KALENDER (ORIGINAL)
// ==========================================
// Bagian ini persis seperti script awal Anda tanpa perubahan.

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
	// 1. Hitung JDN
	jd := dateToJDN(t)

	// 2. Koreksi Tanggal Jawa (Rumus Original Anda)
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

	// --- Header ---
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

	// --- Input Section (MODIFIKASI: INPUT KALENDER) ---
	// --- Input Section (MODIFIKASI: INPUT KALENDER YANG COMPATIBLE) ---
inputLabel := canvas.NewText("Pilih Tanggal / Geblag:", ColorTextGrey)
inputLabel.TextSize = 12

var selectedDate = time.Now()
lblSelectedDate := widget.NewLabel(selectedDate.Format("02 January 2006"))
lblSelectedDate.TextStyle = fyne.TextStyle{Bold: true, Monospace: true}
lblSelectedDate.Alignment = fyne.TextAlignCenter

// Perbaikan pada bagian ini:
btnPickDate := widget.NewButtonWithIcon("Buka Kalender", theme.CalendarIcon(), func() {
    // Memanggil DatePicker dengan cara yang lebih standar
    dateDialog := dialog.NewDatepicker(func(t time.Time) {
        if t.IsZero() { // Cek jika user membatalkan
            return
        }
        selectedDate = t
        bulan := BulanIndo[t.Month()]
        lblSelectedDate.SetText(fmt.Sprintf("%02d %s %d", t.Day(), bulan, t.Year()))
    }, myWindow)
    
    dateDialog.Show()
})

	btnCalc := widget.NewButton("Hitung Selamatan", nil)
	btnCalc.Importance = widget.HighImportance // Membuat tombol lebih menonjol

	// Layout Input
	inputRow := container.NewVBox(
		lblSelectedDate,
		container.NewGridWithColumns(2, btnPickDate, btnCalc),
	)
	
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
		// Tidak perlu parsing string lagi, langsung pakai variabel 'selectedDate'
		// Namun perlu kita pastikan jam-nya 00:00:00 untuk kalkulasi hari
		t := selectedDate
		
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
		// Normalize selected date to midnight
		t = time.Date(t.Year(), t.Month(), t.Day(), 0,0,0,0, t.Location())

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
}

