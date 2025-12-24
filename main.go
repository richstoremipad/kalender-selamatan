package main

import (
	"fmt"
	"image/color"
	"math"
	"strconv"
	"strings"
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
// 1. LOGIKA MATEMATIKA & KALENDER
// ==========================================

var (
	HariIndo  = []string{"Minggu", "Senin", "Selasa", "Rabu", "Kamis", "Jumat", "Sabtu"}
	Pasaran   = []string{"Legi", "Pahing", "Pon", "Wage", "Kliwon"}
	BulanIndo = []string{"", "Januari", "Februari", "Maret", "April", "Mei", "Juni", "Juli", "Agustus", "September", "Oktober", "November", "Desember"}
	BulanJawa = []string{"", "Suro", "Sapar", "Mulud", "Bakda Mulud", "Jumadil Awal", "Jumadil Akhir", "Rajeb", "Ruwah", "Poso", "Sawal", "Sela", "Besar"}
)

// Konversi tanggal Masehi ke Format Jawa (Logic Porting dari Bash)
func getJavaneseDate(t time.Time) string {
	d := t.Day()
	m := int(t.Month())
	y := t.Year()

	// Logic Weton (Pasaran)
	// Unix timestamp / 86400 + 4 % 5 (Mirip logic bash)
	// Kita gunakan referensi simple: 20 Januari 2027 adalah Rabu Wage (Untuk validasi)
	// Tapi cara paling aman pakai Modulo Julian Day atau referensi tanggal pas.
	// Menggunakan logic bash: ((target_ts / 86400) + 4) % 5
	unixDays := t.Unix() / 86400
	pasaranIdx := (unixDays + 4) % 5
	if pasaranIdx < 0 {
		pasaranIdx += 5
	}

	// Logic Tahun/Bulan Jawa (Porting Matematika Julian Day dari Bash)
	if m <= 2 {
		y -= 1
		m += 12
	}
	
	a := y / 100
	b := 2 - a + (a / 4)
	
	// Float calculation untuk presisi Julian Day
	jd := math.Floor(365.25*float64(y+4716)) + math.Floor(30.6001*float64(m+1)) + float64(d) + float64(b) - 1524.5
	l := int(jd) - 1948440 + 10632
	n := (l - 1) / 10631
	l = l - 10631*n + 354
	j := ((10985 - l) / 5316) * ((50 * l) / 17719) + (l / 5670) * ((43 * l) / 15238)
	l = l - (((30 - j) / 15) * ((17719 * j) / 50)) - ((j / 16) * ((15238 * j) / 43)) + 29
	hm := (24 * l) / 709
	hd := l - ((709 * hm) / 24)

	// Pastikan index bulan valid (1-12)
	namaBulanJawa := ""
	if hm > 0 && hm < len(BulanJawa) {
		namaBulanJawa = BulanJawa[hm]
	} else {
		// Fallback jika kalkulasi math meleset sedikit (sangat jarang)
		namaBulanJawa = "Unknown" 
	}

	return fmt.Sprintf("%d %s", hd, namaBulanJawa)
}

func formatWeton(t time.Time) string {
	hari := HariIndo[t.Weekday()]
	unixDays := t.Unix() / 86400
	pasaranIdx := (unixDays + 4) % 5
	if pasaranIdx < 0 {
		pasaranIdx += 5
	}
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

// Warna Palette (Mirip Gambar)
var (
	ColorBgDark     = color.NRGBA{R: 30, G: 33, B: 40, A: 255}    // Background Utama
	ColorCardBg     = color.NRGBA{R: 45, G: 48, B: 55, A: 255}    // Abu Gelap Card
	ColorHeaderTop  = color.NRGBA{R: 40, G: 180, B: 160, A: 255}  // Teal/Tosca
	ColorHeaderBot  = color.NRGBA{R: 50, G: 80, B: 160, A: 255}   // Biru
	ColorBtnOrange  = color.NRGBA{R: 230, G: 150, B: 50, A: 255}  // Orange Tombol
	ColorTextWhite  = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	ColorTextGrey   = color.NRGBA{R: 180, G: 180, B: 180, A: 255}
	ColorBadgeGreen = color.NRGBA{R: 46, G: 125, B: 50, A: 255}
	ColorBadgeRed   = color.NRGBA{R: 198, G: 40, B: 40, A: 255}
	ColorBadgeBlue  = color.NRGBA{R: 21, G: 101, B: 192, A: 255}
)

// createCard membuat tampilan kartu per baris data
func createCard(title, subTitle, dateStr, wetonStr string, statusType int, diffDays int) fyne.CanvasObject {
	// 1. Status Badge Logic
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

	// 2. Komponen UI
	
	// Title & Subtitle (Kiri)
	lblTitle := canvas.NewText(title, ColorTextWhite)
	lblTitle.TextSize = 16
	lblTitle.TextStyle = fyne.TextStyle{Bold: true}

	lblSub := canvas.NewText(subTitle, ColorTextGrey)
	lblSub.TextSize = 12

	leftCont := container.NewVBox(lblTitle, lblSub)

	// Date & Weton (Kanan)
	lblDate := canvas.NewText(dateStr, ColorTextWhite)
	lblDate.Alignment = fyne.TextAlignTrailing
	lblDate.TextSize = 14
	lblDate.TextStyle = fyne.TextStyle{Bold: true}

	lblWeton := canvas.NewText(wetonStr, ColorTextGrey)
	lblWeton.Alignment = fyne.TextAlignTrailing
	lblWeton.TextSize = 11

	rightCont := container.NewVBox(lblDate, lblWeton)

	topRow := container.NewBorder(nil, nil, leftCont, rightCont)

	// Badge (Bawah)
	lblBadge := canvas.NewText(badgeTextStr, ColorTextWhite)
	lblBadge.TextSize = 11
	lblBadge.TextStyle = fyne.TextStyle{Bold: true}
	
	badgeBg := canvas.NewRectangle(badgeColor)
	badgeBg.CornerRadius = 12
	
	// Padding untuk badge text
	badgeCont := container.NewStack(
		badgeBg,
		container.NewPadded(lblBadge),
	)
	
	// Layout Bawah (Badge di kiri atau kanan? Gambar referensi di kiri/kanan. Kita taruh kiri & kanan sesuai context)
	// Kita buat baris bawah untuk badge
	botRow := container.NewHBox(badgeCont)

	// Gabungkan Top dan Bot
	content := container.NewVBox(topRow, container.NewPadded(botRow))

	// 3. Background Card
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
	
	headerIcon := canvas.NewImageFromResource(theme.InfoIcon()) // Placeholder icon
	headerIcon.SetMinSize(fyne.NewSize(30,30))

	headerStack := container.NewStack(
		gradient,
		container.NewPadded(container.NewVBox(
			layout.NewSpacer(),
			container.NewHBox(layout.NewSpacer(), headerIcon, headerTitle, layout.NewSpacer()),
			layout.NewSpacer(),
		)),
	)
	// Header height fix
	headerContainer := container.NewVBox(headerStack)


	// --- Input Section ---
	inputLabel := canvas.NewText("Input Dagall / Geblag (DD/MM/YYYY)", ColorTextGrey)
	inputLabel.TextSize = 12

	inputEntry := widget.NewEntry()
	inputEntry.PlaceHolder = "Contoh: 20/01/2027"
	inputEntry.TextStyle = fyne.TextStyle{Monospace: true}

	btnCalc := widget.NewButton("Hitung", nil) // Logic nanti
	// Styling button agak tricky di Fyne standard, kita pakai default dulu tapi logic-nya kuat.
	// Untuk warna orange, kita bisa bungkus logic high importance theme, atau biarkan default primary.
	// Agar mirip gambar (tombol di kanan), kita pakai Border Layout.

	inputRow := container.NewBorder(nil, nil, nil, btnCalc, inputEntry)
	
	inputCardBg := canvas.NewRectangle(ColorCardBg)
	inputCardBg.CornerRadius = 8
	inputSection := container.NewStack(
		inputCardBg,
		container.NewPadded(container.NewVBox(inputLabel, inputRow)),
	)


	// --- Result Container ---
	resultBox := container.NewVBox()
	
	// Scroll Container
	scrollArea := container.NewVScroll(container.NewPadded(resultBox))

	// --- Footer ---
	noteText := "Notes: Perhitungan menggunakan rumus lusarlu (3), tusarmo (7), masarmo (40), rosarmo (100), patsarpat (Pendhak 1), rosarji (Pendhak 2), nemsarmo (1000)."
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
		// Parsing DD/MM/YYYY
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

		// Definisi Selamatan (Sesuai Script Bash)
		// Format: Label, Sublabel (hari), Offset Hari
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
		// Normalize now to midnight for fair day comparison
		now = time.Date(now.Year(), now.Month(), now.Day(), 0,0,0,0, now.Location())

		for _, e := range events {
			targetDate := t.AddDate(0, 0, e.Offset)
			// Normalize target
			targetDate = time.Date(targetDate.Year(), targetDate.Month(), targetDate.Day(), 0,0,0,0, targetDate.Location())

			diff := int(targetDate.Sub(now).Hours() / 24)

			status := 3 // Future
			if diff < 0 {
				status = 1 // Lewat
			} else if diff == 0 {
				status = 2 // Hari ini
			}

			// Render Card
			card := createCard(
				e.Name,
				e.Sub,
				formatIndoDate(targetDate),
				formatWeton(targetDate),
				status,
				diff,
			)
			
			// Add spacer and card
			resultBox.Add(card)
			resultBox.Add(layout.NewSpacer()) // Sedikit jarak
		}
		resultBox.Refresh()
	}

	// --- Layout Utama ---
	// Background Utama App
	bgApp := canvas.NewRectangle(ColorBgDark)

	// Konten disusun
	mainContent := container.NewBorder(
		container.NewVBox(headerContainer, container.NewPadded(inputSection)), // Top
		container.NewPadded(footerSection), // Bottom
		nil, nil, // Left, Right
		scrollArea, // Center (Isian)
	)

	finalLayout := container.NewStack(bgApp, mainContent)

	myWindow.SetContent(finalLayout)
	myWindow.ShowAndRun()
}

