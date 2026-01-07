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
	
	// Nilai Neptu
	NeptuHari    = []int{5, 4, 3, 7, 8, 6, 9} // Minggu - Sabtu
	NeptuPasaran = []int{5, 9, 7, 4, 8}       // Legi - Kliwon

	// Nama Tahun (Warsa) dalam Siklus Windu (8 Tahun)
	NamaWarsa = []string{"Alip", "Ehe", "Jimawal", "Je", "Dal", "Be", "Wawu", "Jimakhir"}
)

// Struct untuk menampung hasil konversi lengkap
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

// Fungsi helper untuk mendapatkan detail tanggal Jawa
func getJavaneseDetail(t time.Time) JavaneseDateInfo {
	jd := dateToJDN(t)
	l := jd - 1948440 + 10632 + 1
	n := (l - 1) / 10631
	l = l - 10631*n + 354
	j := (int)((10985-l)/5316)*(int)((50*l)/17719) + (int)(l/5670)*(int)((43*l)/15238)
	l = l - (int)((30-j)/15)*(int)((17719*j)/50) - (int)(j/16)*(int)((15238*j)/43) + 29

	hm := (int)(24*l) / 709
	hd := l - (int)(709*hm)/24
	hy := n*30 + (int)((709*hm)/24) - 5 // Tahun Jawa Approximate

	// Koreksi perhitungan tahun (Simple approximation untuk konversi Masehi ke Jawa)
	// Rumus pasti sangat kompleks, ini pendekatan standar: Tahun Masehi + 512 (atau kurangi offset)
	// Kita gunakan hybrid calculation dari JDN untuk presisi
	
	// Kalkulasi Warsa (Siklus 8 tahun)
	// Tahun Jawa = Tahun Masehi + 512 (kurang lebih)
	// Kita ambil approx tahun jawa dari logic di atas atau convert manual
	tahunJawa := t.Year() + 512 
	// Jika bulan masehi < 3 biasanya masih ikut tahun jawa sebelumnya (kasarannya)
	// Untuk akurasi aplikasi sederhana, kita gunakan patokan matematika modulo:
	
	idxWarsa := (tahunJawa - 1) % 8
	if idxWarsa < 0 { idxWarsa += 8 }
	namaWarsa := NamaWarsa[idxWarsa]

	namaBulan := "Unknown"
	if hm > 0 && hm < len(BulanJawa) {
		namaBulan = BulanJawa[hm]
	}

	return JavaneseDateInfo{
		Day:       hd,
		MonthName: namaBulan,
		Year:      tahunJawa,
		Warsa:     namaWarsa,
	}
}

// Menghitung Neptu
func getNeptu(t time.Time) (int, int, int) {
	wDay := t.Weekday() // 0=Minggu, 6=Sabtu
	jd := dateToJDN(t)
	idxPasaran := jd % 5

	valHari := NeptuHari[wDay]
	valPasaran := NeptuPasaran[idxPasaran]
	total := valHari + valPasaran

	return valHari, valPasaran, total
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

// Card untuk List Selamatan (Style Original)
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

// Card Khusus untuk Hasil Weton (Style Baru tapi senada)
func createCardWetonResult(t time.Time) fyne.CanvasObject {
	// 1. Data Masehi & Weton Dasar
	hari := HariIndo[t.Weekday()]
	jd := dateToJDN(t)
	pasaran := Pasaran[jd%5]
	
	// 2. Data Neptu
	valHari, valPasaran, totalNeptu := getNeptu(t)

	// 3. Data Jawa Lengkap
	jawaInfo := getJavaneseDetail(t)

	// --- UI COMPOSITION ---
	
	// Header Weton Besar
	lblWeton := canvas.NewText(fmt.Sprintf("%s %s", hari, pasaran), ColorHeaderTop)
	lblWeton.TextSize = 24
	lblWeton.TextStyle = fyne.TextStyle{Bold: true}
	lblWeton.Alignment = fyne.TextAlignCenter

	// Divider
	line := canvas.NewRectangle(ColorTextGrey)
	line.SetMinSize(fyne.NewSize(100, 1))

	// Rincian Neptu
	txtNeptu := fmt.Sprintf("Neptu: %s (%d) + %s (%d) = %d", hari, valHari, pasaran, valPasaran, totalNeptu)
	lblNeptu := canvas.NewText(txtNeptu, ColorGold) 
	lblNeptu.TextSize = 14
	lblNeptu.Alignment = fyne.TextAlignCenter
	
	// Rincian Tanggal Jawa
	txtJawa := fmt.Sprintf("Tanggal: %d %s %d", jawaInfo.Day, jawaInfo.MonthName, jawaInfo.Year)
	lblJawaDate := canvas.NewText(txtJawa, ColorTextWhite)
	lblJawaDate.TextSize = 14
	lblJawaDate.Alignment = fyne.TextAlignCenter

	// Warsa
	txtWarsa := fmt.Sprintf("Warsa (Tahun): %s", jawaInfo.Warsa)
	lblWarsa := canvas.NewText(txtWarsa, ColorTextGrey)
	lblWarsa.TextStyle = fyne.TextStyle{Italic: true}
	lblWarsa.Alignment = fyne.TextAlignCenter

	// Container Isi
	content := container.NewVBox(
		container.NewPadded(lblWeton),
		container.NewPadded(line),
		lblNeptu,
		lblJawaDate,
		lblWarsa,
	)

	bg := canvas.NewRectangle(ColorCardBg)
	bg.CornerRadius = 12
	
	// Border tipis agar terlihat spesial
	border := canvas.NewRectangle(color.Transparent)
	border.StrokeColor = ColorHeaderBot
	border.StrokeWidth = 2
	border.CornerRadius = 12

	return container.NewStack(bg, border, container.NewPadded(content))
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
	// Pilihan Mode
	modeLabel := canvas.NewText("Pilih Fitur:", ColorTextGrey)
	modeLabel.TextSize = 12
	
	selectMode := widget.NewSelect([]string{"Hitung Selamatan", "Cek Weton Lengkap"}, nil)
	selectMode.Selected = "Hitung Selamatan" // Default

	// Tanggal Input
	dateLabel := canvas.NewText("Pilih Tanggal:", ColorTextGrey)
	dateLabel.TextSize = 12

	selectedDate := time.Now()
	btnSelectDate := widget.NewButton(selectedDate.Format("02/01/2006"), nil)
	btnSelectDate.Icon = theme.CalendarIcon()
	btnSelectDate.OnTapped = func() {
		cal := widget.NewCalendar(selectedDate, func(t time.Time) {
			selectedDate = t
			btnSelectDate.SetText(t.Format("02/01/2006"))
		})
		d := dialog.NewCustom("Pilih Tanggal", "Tutup", cal, myWindow)
		d.Resize(fyne.NewSize(300, 300))
		d.Show()
	}

	btnProcess := widget.NewButton("PROSES DATA", nil)
	btnProcess.Importance = widget.HighImportance

	// Menyusun Card Input
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
		resultBox.Objects = nil // Clear previous results

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
				// Re-normalize just in case
				targetDate = time.Date(targetDate.Year(), targetDate.Month(), targetDate.Day(), 0, 0, 0, 0, targetDate.Location())

				diff := int(targetDate.Sub(now).Hours() / 24)

				status := 3
				if diff < 0 {
					status = 1
				} else if diff == 0 {
					status = 2
				}

				// Menggunakan helper Javanese untuk detail
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
			// Buat satu kartu besar detail
			card := createCardWetonResult(selectedDate)
			
			// Tambahkan info extra di bawahnya
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

	// --- LAYOUT UTAMA ---
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

