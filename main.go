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
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// ==========================================
// 1. LOGIKA MATEMATIKA & KALENDER JAWA
// ==========================================
// (Tidak ada perubahan di bagian ini)

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
// 2. KOMPONEN UI CUSTOM & COLORS
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

// --- TEMA: Memaksa tombol HighImportance jadi Hijau ---
type myTheme struct {
	fyne.Theme
}

func (m myTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	// Kita memaksa warna Primary (yang dipakai tombol HighImportance) jadi Hijau
	if name == theme.ColorNamePrimary {
		return ColorBadgeGreen
	}
	return m.Theme.Color(name, variant)
}

// ==========================================
// 3. LOGIKA KALENDER CUSTOM (GRID MANUAL)
// ==========================================
// Ini adalah solusi untuk masalah warna. Kita buat grid sendiri.

func createCalendarPopup(parent fyne.Window, initialDate time.Time, onSelected func(time.Time)) {
	currentMonth := initialDate // Melacak bulan yang sedang dilihat user
	selectedDate := initialDate // Melacak tanggal yang diklik user

	// Kontainer utama untuk grid tanggal
	gridContainer := container.New(layout.NewGridLayout(7))
	
	// Label Judul (Bulan Tahun)
	lblHeader := widget.NewLabel("")
	lblHeader.Alignment = fyne.TextAlignCenter
	lblHeader.TextStyle = fyne.TextStyle{Bold: true}

	var modal *dialog.CustomDialog // Deklarasi dulu agar bisa ditutup dari dalam

	// Fungsi untuk membangun ulang isi grid
	var refreshGrid func()
	refreshGrid = func() {
		gridContainer.Objects = nil // Kosongkan grid

		// 1. Header Nama Hari (Sen, Sel, Rab...)
		daysHeader := []string{"Min", "Sen", "Sel", "Rab", "Kam", "Jum", "Sab"}
		for _, dayName := range daysHeader {
			l := widget.NewLabel(dayName)
			l.Alignment = fyne.TextAlignCenter
			l.TextStyle = fyne.TextStyle{Bold: true}
			gridContainer.Add(l)
		}

		// 2. Hitung logika tanggal
		year, month, _ := currentMonth.Date()
		lblHeader.SetText(fmt.Sprintf("%s %d", BulanIndo[month], year))

		firstDayOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
		// Cari tahu hari apa tanggal 1 itu (0=Minggu, 1=Senin, dst)
		startWeekday := int(firstDayOfMonth.Weekday())
		
		// Cari tahu jumlah hari dalam bulan ini
		nextMonth := firstDayOfMonth.AddDate(0, 1, 0)
		lastDay := nextMonth.Add(-time.Hour * 24).Day()

		// 3. Isi kotak kosong sebelum tanggal 1
		for i := 0; i < startWeekday; i++ {
			gridContainer.Add(layout.NewSpacer())
		}

		// 4. Isi tombol tanggal 1 sampai akhir bulan
		for d := 1; d <= lastDay; d++ {
			dayNum := d
			dateVal := time.Date(year, month, dayNum, 0, 0, 0, 0, time.Local)
			
			btn := widget.NewButton(fmt.Sprintf("%d", dayNum), nil)
			
			// LOGIKA WARNA HIJAU:
			// Jika tanggal tombol sama dengan tanggal yang dipilih user:
			// Set tombol jadi HighImportance (karena tema kita HighImportance = Hijau)
			if dateVal.Year() == selectedDate.Year() && 
			   dateVal.Month() == selectedDate.Month() && 
			   dateVal.Day() == selectedDate.Day() {
				btn.Importance = widget.HighImportance 
			} else {
				btn.Importance = widget.MediumImportance // Atau LowImportance (transparan/abu)
			}

			// Saat tombol tanggal diklik
			btn.OnTapped = func() {
				selectedDate = dateVal
				onSelected(selectedDate) // Callback ke main app
				refreshGrid() // Refresh tampilan agar warna pindah
				// Uncomment baris bawah jika ingin dialog langsung nutup setelah pilih tanggal
				// modal.Hide() 
			}
			
			gridContainer.Add(btn)
		}
	}

	// Tombol Navigasi Bulan (Prev / Next)
	btnPrev := widget.NewButtonWithIcon("", theme.NavigateBackIcon(), func() {
		currentMonth = currentMonth.AddDate(0, -1, 0)
		refreshGrid()
	})
	btnNext := widget.NewButtonWithIcon("", theme.NavigateNextIcon(), func() {
		currentMonth = currentMonth.AddDate(0, 1, 0)
		refreshGrid()
	})

	navContainer := container.NewBorder(nil, nil, btnPrev, btnNext, lblHeader)
	
	// Layout Popup
	content := container.NewVBox(
		navContainer,
		container.NewPadded(gridContainer),
	)

	// Inisialisasi tampilan awal
	refreshGrid()

	modal = dialog.NewCustom("Pilih Tanggal", "Selesai", content, parent)
	modal.Resize(fyne.NewSize(350, 400))
	modal.Show()
}


// ==========================================
// 4. HELPER UI CARDS
// ==========================================

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
// 5. MAIN APP
// ==========================================

func main() {
	myApp := app.New()
	myApp.Settings().SetTheme(&myTheme{Theme: theme.DefaultTheme()}) // Terapkan tema hijau

	myWindow := myApp.NewWindow("Kalkulator Selamatan Jawa")
	myWindow.Resize(fyne.NewSize(400, 750))

	// Header
	gradient := canvas.NewHorizontalGradient(ColorHeaderTop, ColorHeaderBot)
	headerTitle := canvas.NewText("Kalkulator Selamatan Jawa", ColorTextWhite)
	headerTitle.TextStyle = fyne.TextStyle{Bold: true}
	headerTitle.TextSize = 18
	headerIcon := canvas.NewImageFromResource(theme.InfoIcon())
	headerIcon.SetMinSize(fyne.NewSize(30, 30))
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
	inputLabel := canvas.NewText("Pilih Tanggal Geblag / Wafat:", ColorTextGrey)
	inputLabel.TextSize = 12

	selectedDate := time.Now()

	btnSelectDate := widget.NewButton(selectedDate.Format("02/01/2006"), nil)
	btnSelectDate.Icon = theme.CalendarIcon()
	btnSelectDate.Importance = widget.LowImportance

	// LOGIKA UTAMA: PANGGIL KALENDER CUSTOM
	btnSelectDate.OnTapped = func() {
		createCalendarPopup(myWindow, selectedDate, func(newDate time.Time) {
			selectedDate = newDate
			btnSelectDate.SetText(newDate.Format("02/01/2006"))
		})
	}

	btnCalc := widget.NewButton("Hitung Selamatan", nil)
	btnCalc.Importance = widget.HighImportance

	inputRow := container.NewBorder(nil, nil, nil, nil, btnSelectDate)
	inputCardBg := canvas.NewRectangle(ColorCardBg)
	inputCardBg.CornerRadius = 8
	inputSection := container.NewStack(
		inputCardBg,
		container.NewPadded(container.NewVBox(inputLabel, inputRow, layout.NewSpacer(), btnCalc)),
	)

	// --- Result Container ---
	resultBox := container.NewVBox()
	scrollArea := container.NewVScroll(container.NewPadded(resultBox))

	// --- Footer ---
	noteText := "Notes: Perhitungan ini menggunakan rumus (3, 7, 40, 100, Pendhak 1 & 2, 1000). Jika ada selisih 1 hari, itu wajar karena perbedaan penentuan awal bulan Hijriah/Jawa."
	lblNote := widget.NewLabel(noteText)
	lblNote.Wrapping = fyne.TextWrapWord
	lblNote.TextStyle = fyne.TextStyle{Italic: true}
	lblCredit := canvas.NewText("Code by Richo", ColorTextGrey)
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
		resultBox.Objects = nil
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
		t := selectedDate
		now := time.Now()
		now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())

		for _, e := range events {
			targetDate := t.AddDate(0, 0, e.Offset)
			targetDate = time.Date(targetDate.Year(), targetDate.Month(), targetDate.Day(), 0, 0, 0, 0, targetDate.Location())
			diff := int(targetDate.Sub(now).Hours() / 24)
			status := 3
			if diff < 0 {
				status = 1
			} else if diff == 0 {
				status = 2
			}
			card := createCard(e.Name, e.Sub, formatIndoDate(targetDate), formatWeton(targetDate), status, diff)
			resultBox.Add(card)
			resultBox.Add(layout.NewSpacer())
		}
		resultBox.Refresh()
	}

	bgApp := canvas.NewRectangle(ColorBgDark)
	mainContent := container.NewBorder(
		container.NewVBox(headerContainer, container.NewPadded(inputSection)),
		container.NewPadded(footerSection),
		nil, nil,
		scrollArea,
	)
	myWindow.SetContent(container.NewStack(bgApp, mainContent))
	myWindow.ShowAndRun()
}

