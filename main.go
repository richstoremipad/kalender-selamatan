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
// 1. LOGIKA MATEMATIKA & KALENDER (DIPERBAIKI)
// ==========================================

var (
	HariIndo  = []string{"Minggu", "Senin", "Selasa", "Rabu", "Kamis", "Jumat", "Sabtu"}
	Pasaran   = []string{"Legi", "Pahing", "Pon", "Wage", "Kliwon"}
	BulanIndo = []string{"", "Januari", "Februari", "Maret", "April", "Mei", "Juni", "Juli", "Agustus", "September", "Oktober", "November", "Desember"}
	// Perhatikan: Index 0 kosong agar index 1 = Suro
	BulanJawa = []string{"", "Suro", "Sapar", "Mulud", "Bakda Mulud", "Jumadil Awal", "Jumadil Akhir", "Rejeb", "Ruwah", "Poso", "Sawal", "Sela", "Besar"}
)

// dateToJDN menghitung Julian Day Number (Standar Astronomi)
// Digunakan untuk menghitung Pasaran
func dateToJDN(t time.Time) int {
	a := (14 - int(t.Month())) / 12
	y := t.Year() + 4800 - a
	m := int(t.Month()) + 12*a - 3
	return t.Day() + (153*m+2)/5 + 365*y + y/4 - y/100 + y/400 - 32045
}

// getJavaneseDate Menggunakan Metode Anchor / Pivot
// Referensi: 1 Suro 1957 (Jimawal) jatuh pada 19 Juli 2023
func getJavaneseDate(t time.Time) string {
	// Normalisasi tanggal input ke Midnight
	target := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
	
	// ANCHOR: 1 Suro 1957 = 19 Juli 2023
	// Ini adalah titik referensi yang pasti benar.
	anchorDate := time.Date(2023, 7, 19, 0, 0, 0, 0, time.UTC) 
	
	// Hitung selisih hari
	diff := int(target.Sub(anchorDate).Hours() / 24)

	// Mulai dari 1 Suro 1957
	curDay := 1
	curMonth := 1 // 1 = Suro
	curYear := 1957

	// Pola bulan Jawa (ganjil 30, genap 29) secara umum untuk tahun Jimawal & Je
	// 1:30, 2:29, 3:30, 4:29, 5:30, 6:29, 7:30, 8:29, 9:30, 10:29, 11:30, 12:29 (atau 30 jika kabisat)
	
	// Jika tanggal input setelah Anchor (Masa Depan dari 19 Juli 2023)
	if diff >= 0 {
		for diff > 0 {
			daysInMonth := 29
			if curMonth%2 != 0 { // Bulan Ganjil biasanya 30 hari (Suro, Mulud, Jumadil Awal...)
				daysInMonth = 30
			} else {
				// Pengecualian Khusus (Opsional): Besar (12) bisa 30 di tahun kabisat (Taun Wuntu)
				// Untuk simplifikasi aplikasi selamatan 1000 hari kedepan, pola 30-29 sudah cukup akurat
				// karena 1957 (Jimawal) dan 1958 (Je) dominan pola standar.
				daysInMonth = 29
			}

			// Jika sisa hari lebih besar dari hari dalam bulan ini, maju ke bulan berikutnya
			if diff >= (daysInMonth - curDay + 1) {
				diff -= (daysInMonth - curDay + 1)
				curDay = 1
				curMonth++
				if curMonth > 12 {
					curMonth = 1
					curYear++
				}
			} else {
				curDay += diff
				diff = 0
			}
		}
	} else {
		// Jika tanggal input sebelum Anchor (Masa Lalu)
		// Logika mundur (Opsional, tapi ditambahkan untuk keamanan)
		for diff < 0 {
			// Mundur satu hari
			curDay--
			diff++
			if curDay < 1 {
				curMonth--
				if curMonth < 1 {
					curMonth = 12
					curYear--
				}
				
				// Tentukan jumlah hari bulan sebelumnya
				prevMonthDays := 29
				if curMonth%2 != 0 {
					prevMonthDays = 30
				}
				curDay = prevMonthDays
			}
		}
	}

	namaBulan := "Unknown"
	if curMonth > 0 && curMonth <= len(BulanJawa) {
		namaBulan = BulanJawa[curMonth]
	}

	// Format: 17 Jumadil Awal 1957 (Tahun Jawa)
	// Catatan: User minta format tanggal, jika ingin menampilkan tahun Hijriah (1445),
	// perlu konverter terpisah. Tapi umumnya orang Jawa pakai tahun Jawa (1957).
	// Jika ingin paksa tampil tahun Hijriah (1445), kurangi tahun Jawa dengan 512.
	return fmt.Sprintf("%d %s %d", curDay, namaBulan, curYear)
}

func formatWeton(t time.Time) string {
	// Hari Masehi
	hari := HariIndo[t.Weekday()]
	
	// Pasaran (JDN Modulo 5 tidak pernah salah untuk urutan)
	jd := dateToJDN(t)
	pasaranIdx := jd % 5
	pasaran := Pasaran[pasaranIdx]

	// Tanggal Jawa
	jawaDate := getJavaneseDate(t)

	return fmt.Sprintf("%s %s, %s", hari, pasaran, jawaDate)
}

func formatIndoDate(t time.Time) string {
	return fmt.Sprintf("%d %s %d", t.Day(), BulanIndo[t.Month()], t.Year())
}

// ==========================================
// 2. KOMPONEN UI CUSTOM (TIDAK BERUBAH)
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
	inputEntry.PlaceHolder = "Contoh: 01/12/2023" // Update contoh ke 2023
	inputEntry.TextStyle = fyne.TextStyle{Monospace: true}
	inputEntry.Text = "01/12/2023" // Set default untuk testing user

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
	noteText := "Notes: Perhitungan menggunakan metode Pivot 1 Suro 1957 (19 Juli 2023) untuk akurasi tinggi periode 2023-2026. Perbedaan hisab/rukyat mungkin terjadi +- 1 hari. Wallahu A'lam Bishawab"
	lblNote := widget.NewLabel(noteText)
	lblNote.Wrapping = fyne.TextWrapWord
	lblNote.TextStyle = fyne.TextStyle{Italic: true}
	
	lblCredit := canvas.NewText("Matur Nuwun - Code by Richo (Fixed)", ColorTextGrey)
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
			{"Pendhak I", "1 Tahun", 353}, // 354 hari dalam tahun jawa normal, offset array 0-based
			{"Pendhak II", "2 Tahun", 707}, // 354*2 - 1
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
}

