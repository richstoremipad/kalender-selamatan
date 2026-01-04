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
// 1. LOGIKA MATEMATIKA & KALENDER JAWA
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

type myTheme struct {
	fyne.Theme
}

func (m myTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	if name == theme.ColorNamePrimary {
		return ColorBadgeGreen
	}
	return m.Theme.Color(name, variant)
}

// ==========================================
// 3. LOGIKA KALENDER CUSTOM (GRID MANUAL)
// ==========================================

// PERBAIKAN LOGIKA:
// Menambahkan parameter 'onDateChanged' agar UI utama bisa update realtime saat tombol tanggal diklik.
func createCalendarPopup(parentCanvas fyne.Canvas, initialDate time.Time, onDateChanged func(time.Time), onCalculate func(time.Time)) {
	currentMonth := initialDate
	selectedDate := initialDate

	gridContainer := container.New(layout.NewGridLayout(7))
	
	lblHeader := widget.NewLabel("")
	lblHeader.Alignment = fyne.TextAlignCenter
	lblHeader.TextStyle = fyne.TextStyle{Bold: true}

	var popup *widget.PopUp

	var refreshGrid func()
	refreshGrid = func() {
		gridContainer.Objects = nil

		daysHeader := []string{"Min", "Sen", "Sel", "Rab", "Kam", "Jum", "Sab"}
		for _, dayName := range daysHeader {
			l := widget.NewLabel(dayName)
			l.Alignment = fyne.TextAlignCenter
			l.TextStyle = fyne.TextStyle{Bold: true}
			gridContainer.Add(l)
		}

		year, month, _ := currentMonth.Date()
		lblHeader.SetText(fmt.Sprintf("%s %d", BulanIndo[month], year))

		firstDayOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
		startWeekday := int(firstDayOfMonth.Weekday())
		
		nextMonth := firstDayOfMonth.AddDate(0, 1, 0)
		lastDay := nextMonth.Add(-time.Hour * 24).Day()

		for i := 0; i < startWeekday; i++ {
			gridContainer.Add(layout.NewSpacer())
		}

		for d := 1; d <= lastDay; d++ {
			dayNum := d
			dateVal := time.Date(year, month, dayNum, 0, 0, 0, 0, time.Local)
			
			btn := widget.NewButton(fmt.Sprintf("%d", dayNum), nil)
			
			if dateVal.Year() == selectedDate.Year() && 
			   dateVal.Month() == selectedDate.Month() && 
			   dateVal.Day() == selectedDate.Day() {
				btn.Importance = widget.HighImportance 
			} else {
				btn.Importance = widget.MediumImportance
			}

			// LOGIKA REALTIME DI SINI:
			btn.OnTapped = func() {
				selectedDate = dateVal
				refreshGrid() // Refresh warna tombol di kalender
				
				// Panggil callback agar label di layar utama langsung berubah
				if onDateChanged != nil {
					onDateChanged(selectedDate)
				}
			}
			
			gridContainer.Add(btn)
		}
	}

	btnPrev := widget.NewButtonWithIcon("", theme.NavigateBackIcon(), func() {
		currentMonth = currentMonth.AddDate(0, -1, 0)
		refreshGrid()
	})
	btnNext := widget.NewButtonWithIcon("", theme.NavigateNextIcon(), func() {
		currentMonth = currentMonth.AddDate(0, 1, 0)
		refreshGrid()
	})

	navContainer := container.NewBorder(nil, nil, btnPrev, btnNext, lblHeader)
	
	btnHitung := widget.NewButton("Hitung", func() {
		if popup != nil {
			popup.Hide()
		}
		// Jalankan perhitungan akhir
		onCalculate(selectedDate)
	})
	btnHitung.Importance = widget.HighImportance
	btnHitung.Icon = theme.ConfirmIcon()

	btnBatal := widget.NewButton("Batal", func() {
		if popup != nil {
			popup.Hide()
		}
	})
	btnBatal.Importance = widget.LowImportance

	buttonRow := container.NewGridWithColumns(2, btnBatal, btnHitung)

	content := container.NewVBox(
		navContainer,
		container.NewPadded(gridContainer),
		layout.NewSpacer(),
		container.NewPadded(buttonRow),
	)

	cardWrapper := container.NewStack(
		canvas.NewRectangle(ColorCardBg), 
		container.NewPadded(content),
	)

	refreshGrid()

	popup = widget.NewModalPopUp(cardWrapper, parentCanvas)
	popup.Resize(fyne.NewSize(350, 420))
	popup.Show()
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
	myApp.Settings().SetTheme(&myTheme{Theme: theme.DefaultTheme()}) 

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

	// --- Result Container ---
	resultBox := container.NewVBox()
	scrollArea := container.NewVScroll(container.NewPadded(resultBox))

	// Variable tanggal
	calcDate := time.Now()

	// UI Komponen Tanggal di Halaman Utama
	lblDateTitle := canvas.NewText("Tanggal Wafat / Geblag:", ColorTextGrey)
	lblDateTitle.TextSize = 12
	
	lblSelectedDate := widget.NewLabel("")
	lblSelectedDate.Alignment = fyne.TextAlignCenter
	lblSelectedDate.TextStyle = fyne.TextStyle{Bold: true}

	// Fungsi Helper untuk Update Label Tanggal
	updateDateLabel := func(t time.Time) {
		lblSelectedDate.SetText(formatIndoDate(t))
	}
	// Inisialisasi label dengan tanggal hari ini
	updateDateLabel(calcDate)

	// --- Logic Calculation Function ---
	performCalculation := func(t time.Time) {
		// Pastikan label terupdate (jaga-jaga)
		updateDateLabel(t)

		// Bersihkan hasil sebelumnya
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
		
		// Normalisasi waktu ke 00:00:00
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

	// Tombol Utama
	btnOpenCalc := widget.NewButton("Pilih Tanggal & Hitung", nil)
	btnOpenCalc.Importance = widget.HighImportance
	btnOpenCalc.Icon = theme.CalendarIcon()

	btnOpenCalc.OnTapped = func() {
		createCalendarPopup(myWindow.Canvas(), calcDate, 
			// Callback 1: Realtime update saat tanggal diklik
			func(realtimeDate time.Time) {
				calcDate = realtimeDate
				updateDateLabel(calcDate) // <--- Ini yang mengubah tampilan 'Tanggal Wafat' seketika
			},
			// Callback 2: Eksekusi hitungan saat tombol Hitung ditekan
			func(finalDate time.Time) {
				calcDate = finalDate
				performCalculation(calcDate)
			},
		)
	}

	inputRow := container.NewBorder(nil, nil, nil, nil, lblSelectedDate)
	inputCardBg := canvas.NewRectangle(ColorCardBg)
	inputCardBg.CornerRadius = 8
	
	inputSection := container.NewStack(
		inputCardBg,
		container.NewPadded(container.NewVBox(
			lblDateTitle, 
			inputRow, 
			layout.NewSpacer(), 
			btnOpenCalc,
		)),
	)

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

