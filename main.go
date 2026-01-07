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
// 1. LOGIKA MATEMATIKA & KALENDER
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
	j := (int)((10985-l)/5316)*(int)((50*l)/17719) + (int)(l/5670)*(int)((43*l)/15238)
	l = l - (int)((30-j)/15)*(int)((17719*j)/50) - (int)(j/16)*(int)((15238*j)/43) + 29

	hm := (int)(24*l) / 709
	hd := l - (int)(709*hm)/24

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

// Helper khusus untuk mengambil nama pasaran saja
func getPasaranOnly(t time.Time) string {
	jd := dateToJDN(t)
	pasaranIdx := jd % 5
	return Pasaran[pasaranIdx]
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

	badgeCont := container.NewStack(badgeBg, container.NewPadded(lblBadge))
	botRow := container.NewHBox(badgeCont)
	content := container.NewVBox(topRow, container.NewPadded(botRow))

	bg := canvas.NewRectangle(ColorCardBg)
	bg.CornerRadius = 10
	return container.NewStack(bg, container.NewPadded(content))
}

// ==========================================
// 3. SCREEN BUILDERS
// ==========================================

// --- Screen 1: Hitung Selamatan (Fitur Lama) ---
func makeSelamatanScreen(myWindow fyne.Window) fyne.CanvasObject {
	// Header
	gradient := canvas.NewHorizontalGradient(ColorHeaderTop, ColorHeaderBot)
	headerTitle := canvas.NewText("Hitung Selamatan", ColorTextWhite)
	headerTitle.TextStyle = fyne.TextStyle{Bold: true}
	headerTitle.TextSize = 18

	headerIcon := canvas.NewImageFromResource(theme.HistoryIcon()) // Ganti icon
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

	inputLabel := canvas.NewText("Pilih Tanggal Geblag / Wafat:", ColorTextGrey)
	inputLabel.TextSize = 12

	selectedDate := time.Now()

	btnSelectDate := widget.NewButton(selectedDate.Format("02/01/2006"), nil)
	btnSelectDate.Icon = theme.CalendarIcon()
	btnSelectDate.Importance = widget.LowImportance

	btnSelectDate.OnTapped = func() {
		cal := widget.NewCalendar(selectedDate, func(t time.Time) {
			selectedDate = t
			btnSelectDate.SetText(t.Format("02/01/2006"))
		})
		d := dialog.NewCustom("Pilih Tanggal Wafat", "Tutup", cal, myWindow)
		d.Resize(fyne.NewSize(300, 300))
		d.Show()
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

	resultBox := container.NewVBox()
	scrollArea := container.NewVScroll(container.NewPadded(resultBox))

	// Footer Logic
	noteText := "Rumus: 3, 7, 40, 100, Pendhak 1&2, 1000 hari."
	lblNote := widget.NewLabel(noteText)
	lblNote.Wrapping = fyne.TextWrapWord
	lblNote.TextStyle = fyne.TextStyle{Italic: true}
	footerCardBg := canvas.NewRectangle(ColorCardBg)
	footerCardBg.CornerRadius = 8
	footerSection := container.NewStack(
		footerCardBg,
		container.NewPadded(container.NewVBox(lblNote)),
	)

	// Calculation Logic
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
		// Normalize time
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

	bgApp := canvas.NewRectangle(ColorBgDark)
	content := container.NewBorder(
		container.NewVBox(headerContainer, container.NewPadded(inputSection)),
		container.NewPadded(footerSection),
		nil, nil,
		scrollArea,
	)

	return container.NewStack(bgApp, content)
}

// --- Screen 2: Cek Weton (Fitur Baru) ---
func makeWetonScreen(myWindow fyne.Window) fyne.CanvasObject {
	// Header
	gradient := canvas.NewHorizontalGradient(ColorHeaderTop, ColorHeaderBot)
	headerTitle := canvas.NewText("Cek Weton Lahir", ColorTextWhite)
	headerTitle.TextStyle = fyne.TextStyle{Bold: true}
	headerTitle.TextSize = 18
	headerIcon := canvas.NewImageFromResource(theme.SearchIcon())
	headerIcon.SetMinSize(fyne.NewSize(30, 30))

	headerStack := container.NewStack(
		gradient,
		container.NewPadded(container.NewVBox(
			layout.NewSpacer(),
			container.NewHBox(layout.NewSpacer(), headerIcon, headerTitle, layout.NewSpacer()),
			layout.NewSpacer(),
		)),
	)

	// Input Section
	inputLabel := canvas.NewText("Pilih Tanggal Lahir:", ColorTextGrey)
	inputLabel.TextSize = 12

	selectedDate := time.Now()

	btnSelectDate := widget.NewButton(selectedDate.Format("02/01/2006"), nil)
	btnSelectDate.Icon = theme.CalendarIcon()

	btnSelectDate.OnTapped = func() {
		cal := widget.NewCalendar(selectedDate, func(t time.Time) {
			selectedDate = t
			btnSelectDate.SetText(t.Format("02/01/2006"))
		})
		d := dialog.NewCustom("Pilih Tanggal Lahir", "Tutup", cal, myWindow)
		d.Resize(fyne.NewSize(300, 300))
		d.Show()
	}

	// Result Display Components
	lblResultBig := canvas.NewText("-", ColorTextWhite)
	lblResultBig.TextSize = 24
	lblResultBig.TextStyle = fyne.TextStyle{Bold: true}
	lblResultBig.Alignment = fyne.TextAlignCenter

	lblResultJawa := canvas.NewText("-", ColorHeaderTop)
	lblResultJawa.TextSize = 16
	lblResultJawa.Alignment = fyne.TextAlignCenter

	lblResultMasehi := canvas.NewText("-", ColorTextGrey)
	lblResultMasehi.TextSize = 14
	lblResultMasehi.Alignment = fyne.TextAlignCenter

	resultContainer := container.NewVBox(
		layout.NewSpacer(),
		lblResultBig,
		lblResultJawa,
		lblResultMasehi,
		layout.NewSpacer(),
	)
	
	resultCardBg := canvas.NewRectangle(ColorCardBg)
	resultCardBg.CornerRadius = 15
	resultStack := container.NewStack(resultCardBg, container.NewPadded(resultContainer))

	// Button Action
	btnCheck := widget.NewButton("Cek Weton", func() {
		hari := HariIndo[selectedDate.Weekday()]
		pasaran := getPasaranOnly(selectedDate)
		
		// Set Text
		lblResultBig.Text = fmt.Sprintf("%s %s", hari, pasaran)
		lblResultJawa.Text = fmt.Sprintf("Kalender Jawa: %s", getJavaneseDate(selectedDate))
		lblResultMasehi.Text = formatIndoDate(selectedDate)
		
		// Refresh
		resultStack.Refresh()
	})
	btnCheck.Importance = widget.HighImportance

	// Layout Assembly
	inputCardBg := canvas.NewRectangle(ColorCardBg)
	inputCardBg.CornerRadius = 8
	inputSection := container.NewStack(
		inputCardBg,
		container.NewPadded(container.NewVBox(inputLabel, btnSelectDate, layout.NewSpacer(), btnCheck)),
	)

	// Main Layout for this Tab
	bgApp := canvas.NewRectangle(ColorBgDark)
	mainContent := container.NewVBox(
		headerStack,
		container.NewPadded(inputSection),
		layout.NewSpacer(),
		container.NewPadded(resultStack),
		layout.NewSpacer(),
	)

	return container.NewStack(bgApp, mainContent)
}

// ==========================================
// 4. MAIN APP
// ==========================================

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Javanese Calc") // Judul Window
	myWindow.Resize(fyne.NewSize(400, 750))

	// Buat Screen untuk Tab
	screenSelamatan := makeSelamatanScreen(myWindow)
	screenWeton := makeWetonScreen(myWindow)

	// Buat Tab Container
	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon("Selamatan", theme.ListIcon(), screenSelamatan),
		container.NewTabItemWithIcon("Cek Weton", theme.SearchIcon(), screenWeton),
	)

	// Set posisi tab di bawah atau atas (Default atas)
	tabs.SetTabLocation(container.TabLocationTop)

	myWindow.SetContent(tabs)
	myWindow.ShowAndRun()
}

