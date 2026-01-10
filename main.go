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
// 1. DATA & LOGIKA JAWA (LENGKAP)
// ==========================================

var (
	HariIndo  = []string{"Minggu", "Senin", "Selasa", "Rabu", "Kamis", "Jumat", "Sabtu"}
	Pasaran   = []string{"Legi", "Pahing", "Pon", "Wage", "Kliwon"}
	BulanIndo = []string{"", "Januari", "Februari", "Maret", "April", "Mei", "Juni", "Juli", "Agustus", "September", "Oktober", "November", "Desember"}
	BulanJawa = []string{"", "Suro", "Sapar", "Mulud", "Bakda Mulud", "Jumadil Awal", "Jumadil Akhir", "Rajeb", "Ruwah", "Poso", "Sawal", "Sela", "Besar"}
	
	NeptuHari    = []int{5, 4, 3, 7, 8, 6, 9} // Minggu - Sabtu
	NeptuPasaran = []int{5, 9, 7, 4, 8}       // Legi - Kliwon
	NamaWarsa    = []string{"Alip", "Ehe", "Jimawal", "Je", "Dal", "Be", "Wawu", "Jimakhir"}
)

type JavaneseDateInfo struct {
	Day       int
	MonthName string
	Year      int
	Warsa     string
}

func dateToJDN(t time.Time) int {
	a := (14 - int(t.Month())) / 12
	y := t.Year() + 4800 - a
	m := int(t.Month()) + 12*a - 3
	return t.Day() + (153*m+2)/5 + 365*y + y/4 - y/100 + y/400 - 32045
}

func getJavaneseDetail(t time.Time) JavaneseDateInfo {
	jd := dateToJDN(t)
	l := jd - 1948440 + 10632 + 1
	n := (l - 1) / 10631
	l = l - 10631*n + 354
	j := (int)((10985-l)/5316)*(int)((50*l)/17719) + (int)(l/5670)*(int)((43*l)/15238)
	l = l - (int)((30-j)/15)*(int)((17719*j)/50) - (int)(j/16)*(int)((15238*j)/43) + 29

	hm := (int)(24*l) / 709
	hd := l - (int)(709*hm)/24
	
	tahunJawa := t.Year() + 512 
	idxWarsa := (tahunJawa - 1) % 8
	if idxWarsa < 0 { idxWarsa += 8 }
	namaWarsa := NamaWarsa[idxWarsa]

	namaBulan := "Unknown"
	if hm > 0 && hm < len(BulanJawa) {
		namaBulan = BulanJawa[hm]
	}

	return JavaneseDateInfo{Day: hd, MonthName: namaBulan, Year: tahunJawa, Warsa: namaWarsa}
}

func getNeptu(t time.Time) (int, int, int) {
	wDay := t.Weekday()
	jd := dateToJDN(t)
	idxPasaran := jd % 5
	valHari := NeptuHari[wDay]
	valPasaran := NeptuPasaran[idxPasaran]
	return valHari, valPasaran, valHari + valPasaran
}

func formatWetonSimple(t time.Time) string {
	hari := HariIndo[t.Weekday()]
	jd := dateToJDN(t)
	pasaran := Pasaran[jd%5]
	return fmt.Sprintf("%s %s", hari, pasaran)
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
	ColorGold       = color.NRGBA{R: 255, G: 193, B: 7, A: 255}
)

func createCardSelamatan(title, subTitle, dateStr, wetonStr string, statusType int, diffDays int) fyne.CanvasObject {
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

func createCardWetonResult(t time.Time) fyne.CanvasObject {
	hari := HariIndo[t.Weekday()]
	jd := dateToJDN(t)
	pasaran := Pasaran[jd%5]
	valHari, valPasaran, totalNeptu := getNeptu(t)
	jawaInfo := getJavaneseDetail(t)

	lblWeton := canvas.NewText(fmt.Sprintf("%s %s", hari, pasaran), ColorHeaderTop)
	lblWeton.TextSize = 24
	lblWeton.TextStyle = fyne.TextStyle{Bold: true}
	lblWeton.Alignment = fyne.TextAlignCenter

	line := canvas.NewRectangle(ColorTextGrey)
	line.SetMinSize(fyne.NewSize(100, 1))

	txtNeptu := fmt.Sprintf("Neptu: %s (%d) + %s (%d) = %d", hari, valHari, pasaran, valPasaran, totalNeptu)
	lblNeptu := canvas.NewText(txtNeptu, ColorGold) 
	lblNeptu.TextSize = 14
	lblNeptu.Alignment = fyne.TextAlignCenter
	
	txtJawa := fmt.Sprintf("Tanggal: %d %s %d", jawaInfo.Day, jawaInfo.MonthName, jawaInfo.Year)
	lblJawaDate := canvas.NewText(txtJawa, ColorTextWhite)
	lblJawaDate.TextSize = 14
	lblJawaDate.Alignment = fyne.TextAlignCenter

	txtWarsa := fmt.Sprintf("Warsa (Tahun): %s", jawaInfo.Warsa)
	lblWarsa := canvas.NewText(txtWarsa, ColorTextGrey)
	lblWarsa.TextStyle = fyne.TextStyle{Italic: true}
	lblWarsa.Alignment = fyne.TextAlignCenter

	content := container.NewVBox(
		container.NewPadded(lblWeton),
		container.NewPadded(line),
		lblNeptu,
		lblJawaDate,
		lblWarsa,
	)

	bg := canvas.NewRectangle(ColorCardBg)
	bg.CornerRadius = 12
	
	border := canvas.NewRectangle(color.Transparent)
	border.StrokeColor = ColorHeaderBot
	border.StrokeWidth = 2
	border.CornerRadius = 12

	return container.NewStack(bg, border, container.NewPadded(content))
}

// === UTILS: TOAST NOTIFICATION ===
func showToast(w fyne.Window, message string) {
	lbl := widget.NewLabel(message)
	lbl.Alignment = fyne.TextAlignCenter
	
	bg := canvas.NewRectangle(color.NRGBA{R: 50, G: 50, B: 50, A: 220})
	bg.CornerRadius = 8
	
	content := container.NewStack(bg, container.NewPadded(lbl))
	
	// Create PopUp centered
	popup := widget.NewModalPopUp(content, w.Canvas())
	popup.Show()

	// Auto hide after 1.5 seconds
	time.AfterFunc(1500*time.Millisecond, func() {
		popup.Hide()
	})
}

// ==========================================
// 3. MAIN APP
// ==========================================

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Kalkulator Jawa Pro")
	myWindow.Resize(fyne.NewSize(400, 750))

	// --- HEADER ---
	gradient := canvas.NewHorizontalGradient(ColorHeaderTop, ColorHeaderBot)
	headerTitle := canvas.NewText("Kalkulator Jawa", ColorTextWhite)
	headerTitle.TextStyle = fyne.TextStyle{Bold: true}
	headerTitle.TextSize = 18

	headerIcon := canvas.NewImageFromResource(theme.GridIcon())
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

	// --- INPUT SECTION ---
	modeLabel := canvas.NewText("Pilih Fitur:", ColorTextGrey)
	modeLabel.TextSize = 12
	
	selectMode := widget.NewSelect([]string{"Hitung Selamatan", "Cek Weton Lengkap"}, nil)
	selectMode.Selected = "Hitung Selamatan" 

	dateLabel := canvas.NewText("Pilih Tanggal:", ColorTextGrey)
	dateLabel.TextSize = 12

	// STATE MANAGEMENT
	selectedDate := time.Now()
	isDatePicked := false // Flag apakah sudah pilih tanggal

	// BUTTON SELECT DATE
	btnSelectDate := widget.NewButton("Pilih Tanggal...", nil)
	btnSelectDate.Icon = theme.CalendarIcon()
	btnSelectDate.Importance = widget.LowImportance // Default: Abu-abu / Outline

	btnSelectDate.OnTapped = func() {
		// Logika Calendar
		cal := widget.NewCalendar(selectedDate, func(t time.Time) {
			selectedDate = t
			isDatePicked = true
			
			// Update Tampilan Tombol jadi HIJAU dan Teks Berubah
			btnSelectDate.SetText(t.Format("02/01/2006"))
			btnSelectDate.Importance = widget.SuccessImportance // Jadi Hijau
			btnSelectDate.Refresh()
		})
		
		d := dialog.NewCustom("Pilih Tanggal", "Batal", cal, myWindow)
		d.Resize(fyne.NewSize(300, 300))
		d.Show()
	}

	btnProcess := widget.NewButton("PROSES DATA", nil)
	btnProcess.Importance = widget.HighImportance

	inputForm := container.NewVBox(
		modeLabel,
		selectMode,
		layout.NewSpacer(),
		dateLabel,
		container.NewBorder(nil,nil,nil,nil, btnSelectDate),
		layout.NewSpacer(),
		btnProcess,
	)

	inputCardBg := canvas.NewRectangle(ColorCardBg)
	inputCardBg.CornerRadius = 8
	inputSection := container.NewStack(
		inputCardBg,
		container.NewPadded(inputForm),
	)

	// --- RESULT CONTAINER ---
	resultBox := container.NewVBox()
	scrollArea := container.NewVScroll(container.NewPadded(resultBox))

	// --- LOGIC HANDLING ---
	btnProcess.OnTapped = func() {
		// 1. VALIDASI: Cek apakah tanggal sudah dipilih
		if !isDatePicked {
			showToast(myWindow, "⚠️ Silakan Pilih Tanggal Terlebih Dahulu!")
			return
		}

		resultBox.Objects = nil 

		// Mode 1: Hitung Selamatan
		if selectMode.Selected == "Hitung Selamatan" {
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

				jawaInfo := getJavaneseDetail(targetDate)
				wetonFull := formatWetonSimple(targetDate) + fmt.Sprintf(", %d %s", jawaInfo.Day, jawaInfo.MonthName)

				card := createCardSelamatan(
					e.Name,
					e.Sub,
					formatIndoDate(targetDate),
					wetonFull,
					status,
					diff,
				)
				resultBox.Add(card)
				resultBox.Add(layout.NewSpacer())
			}
		
		// Mode 2: Cek Weton Lengkap
		} else {
			card := createCardWetonResult(selectedDate)
			
			lblInfo := widget.NewLabel("Perhitungan Warsa & Tanggal Jawa menggunakan pendekatan aritmatika Masehi-Jawa.")
			lblInfo.Wrapping = fyne.TextWrapWord
			lblInfo.TextStyle = fyne.TextStyle{Italic: true}
			
			resultBox.Add(card)
			resultBox.Add(layout.NewSpacer())
			resultBox.Add(container.NewPadded(lblInfo))
		}
		
		resultBox.Refresh()
	}

	// --- FOOTER ---
	footerText := canvas.NewText("Code by Richo", ColorTextGrey)
	footerText.TextSize = 10
	footerText.Alignment = fyne.TextAlignCenter
	footerContainer := container.NewPadded(footerText)

	bgApp := canvas.NewRectangle(ColorBgDark)
	mainContent := container.NewBorder(
		container.NewVBox(headerContainer, container.NewPadded(inputSection)),
		footerContainer,
		nil, nil,
		scrollArea,
	)

	myWindow.SetContent(container.NewStack(bgApp, mainContent))
	myWindow.ShowAndRun()
}
