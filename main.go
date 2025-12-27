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
// 1. LOGIKA MATEMATIKA & KALENDER (TETAP)
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
// 2. THEME & COLORS (UI DESIGN)
// ==========================================

// Custom Theme untuk memaksa warna seleksi kalender menjadi terang (Orange/Gold)
type myTheme struct{}

func (m myTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	if name == theme.ColorNamePrimary {
		return color.NRGBA{R: 255, G: 152, B: 0, A: 255} // Orange Gold Highlight
	}
	if name == theme.ColorNameBackground {
		return color.NRGBA{R: 30, G: 33, B: 40, A: 255} // Dark BG
	}
	return theme.DefaultTheme().Color(name, theme.VariantDark)
}
func (m myTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}
func (m myTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}
func (m myTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}

var (
	// Palet Warna Modern
	ColorBgDark     = color.NRGBA{R: 28, G: 28, B: 32, A: 255}
	ColorCardBg     = color.NRGBA{R: 44, G: 44, B: 50, A: 255}
	ColorAccent     = color.NRGBA{R: 255, G: 167, B: 38, A: 255} // Orange Gold
	ColorTextWhite  = color.NRGBA{R: 245, G: 245, B: 245, A: 255}
	ColorTextGrey   = color.NRGBA{R: 158, G: 158, B: 158, A: 255}
	ColorBadgeGreen = color.NRGBA{R: 102, G: 187, B: 106, A: 255} // Soft Green
	ColorBadgeRed   = color.NRGBA{R: 239, G: 83, B: 80, A: 255}   // Soft Red
	ColorBadgeBlue  = color.NRGBA{R: 66, G: 165, B: 245, A: 255}  // Soft Blue
	ColorInputBg    = color.NRGBA{R: 58, G: 58, B: 64, A: 255}
)

// Komponen Kartu Hasil
func createCard(title, subTitle, dateStr, wetonStr string, statusType int, diffDays int) fyne.CanvasObject {
	var badgeColor color.Color
	var badgeTextStr string

	// Logika Badge
	switch statusType {
	case 1:
		badgeColor = ColorBadgeGreen
		badgeTextStr = fmt.Sprintf("✓ Selesai (%d hari lalu)", int(math.Abs(float64(diffDays))))
	case 2:
		badgeColor = ColorBadgeRed
		badgeTextStr = "🔔 HARI INI!"
	case 3:
		badgeColor = ColorBadgeBlue
		badgeTextStr = fmt.Sprintf("⏳ %d Hari Lagi", diffDays)
	}

	// Layout Text
	lblTitle := canvas.NewText(title, ColorAccent) // Judul berwarna aksen
	lblTitle.TextSize = 18
	lblTitle.TextStyle = fyne.TextStyle{Bold: true}

	lblSub := canvas.NewText(subTitle, ColorTextGrey)
	lblSub.TextSize = 14

	lblDate := canvas.NewText(dateStr, ColorTextWhite)
	lblDate.Alignment = fyne.TextAlignTrailing
	lblDate.TextSize = 16
	lblDate.TextStyle = fyne.TextStyle{Bold: true}

	lblWeton := canvas.NewText(wetonStr, ColorTextGrey)
	lblWeton.Alignment = fyne.TextAlignTrailing
	lblWeton.TextSize = 12

	// Susunan Header Kartu
	headerRow := container.NewBorder(nil, nil,
		container.NewVBox(lblTitle, lblSub),
		container.NewVBox(lblDate, lblWeton),
	)

	// Badge Style
	lblBadge := canvas.NewText(badgeTextStr, ColorTextWhite)
	lblBadge.TextSize = 12
	lblBadge.TextStyle = fyne.TextStyle{Bold: true}
	
	badgeBg := canvas.NewRectangle(badgeColor)
	badgeBg.CornerRadius = 8
	badgeCont := container.NewStack(badgeBg, container.NewPadded(lblBadge))
	
	// Gabung Konten
	content := container.NewVBox(
		headerRow,
		layout.NewSpacer(),
		container.NewHBox(badgeCont), // Badge di kiri bawah
	)

	// Background Kartu
	bg := canvas.NewRectangle(ColorCardBg)
	bg.CornerRadius = 12

	// Shadow effect (simulated with border/color)
	return container.NewStack(bg, container.NewPadded(container.NewPadded(content)))
}

// ==========================================
// 3. MAIN APP
// ==========================================

func main() {
	myApp := app.New()
	
	// Terapkan Custom Theme agar Kalender terlihat kontras
	myApp.Settings().SetTheme(&myTheme{})
	
	myWindow := myApp.NewWindow("Kalkulator Selamatan Jawa")
	myWindow.Resize(fyne.NewSize(420, 800))

	// --- 1. Header Section ---
	// Background Gradient Header
	headerBg := canvas.NewHorizontalGradient(color.NRGBA{R: 33, G: 150, B: 243, A: 255}, color.NRGBA{R: 21, G: 101, B: 192, A: 255})
	
	titleText := canvas.NewText("Selamatan Jawa", ColorTextWhite)
	titleText.TextStyle = fyne.TextStyle{Bold: true}
	titleText.TextSize = 24
	
	subTitleText := canvas.NewText("Hitung peringatan kematian & weton", color.NRGBA{R: 200, G: 200, B: 200, A: 255})
	subTitleText.TextSize = 12

	headerContent := container.NewVBox(
		titleText,
		subTitleText,
	)

	headerContainer := container.NewStack(
		headerBg,
		container.NewPadded(container.NewPadded(headerContent)),
	)

	// --- 2. Input Section Custom ---
	inputLabel := canvas.NewText("Pilih Tanggal Geblag (Wafat):", ColorTextGrey)
	inputLabel.TextSize = 14

	selectedDate := time.Now()

	// Membuat tampilan Custom untuk tombol Tanggal agar lebih cantik
	// Text tanggal akan kita ubah warnanya (ColorAccent)
	lblSelectedDate := canvas.NewText(selectedDate.Format("02 January 2006"), ColorAccent)
	lblSelectedDate.TextSize = 18
	lblSelectedDate.TextStyle = fyne.TextStyle{Bold: true}
	
	iconCal := widget.NewIcon(theme.CalendarIcon())
	
	inputBg := canvas.NewRectangle(ColorInputBg)
	inputBg.CornerRadius = 10
	inputBg.StrokeColor = ColorTextGrey
	inputBg.StrokeWidth = 1

	// Container untuk teks tanggal (Custom Widget behavior)
	dateDisplay := container.NewStack(
		inputBg,
		container.NewPadded(
			container.NewHBox(
				iconCal,
				layout.NewSpacer(),
				lblSelectedDate,
				layout.NewSpacer(),
			),
		),
	)

	// Membuat dateDisplay bisa diklik
	btnSelectDate := widget.NewButton("", nil) // Tombol transparan di atas display
	btnSelectDate.OnTapped = func() {
		cal := widget.NewCalendar(selectedDate, func(t time.Time) {
			selectedDate = t
			// Update Teks dan Warna saat dipilih
			lblSelectedDate.Text = formatIndoDate(t)
			lblSelectedDate.Color = ColorAccent // Pastikan warna Emas/Orange
			lblSelectedDate.Refresh()
		})
		
		d := dialog.NewCustom("Pilih Tanggal Wafat", "Tutup", cal, myWindow)
		d.Resize(fyne.NewSize(320, 350))
		d.Show()
	}
	
	// Stack tombol transparan diatas visual dateDisplay
	datePickerContainer := container.NewStack(dateDisplay, btnSelectDate)

	btnCalc := widget.NewButton("HITUNG SELAMATAN", nil)
	btnCalc.Importance = widget.HighImportance

	inputSection := container.NewVBox(
		inputLabel,
		datePickerContainer,
		layout.NewSpacer(),
		btnCalc,
	)

	// --- 3. Result Section ---
	resultBox := container.NewVBox()
	
	// Pesan awal kosong
	emptyImg := widget.NewIcon(theme.SearchIcon())
	emptyText := canvas.NewText("Pilih tanggal lalu klik Hitung", ColorTextGrey)
	emptyText.Alignment = fyne.TextAlignCenter
	emptyState := container.NewVBox(layout.NewSpacer(), emptyImg, emptyText, layout.NewSpacer())
	
	resultBox.Add(emptyState)

	scrollArea := container.NewVScroll(container.NewPadded(resultBox))

	// --- 4. Logic Calculation ---
	btnCalc.OnTapped = func() {
		resultBox.Objects = nil // Clear previous results

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
		
		// Normalisasi jam ke 00:00
		now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())

		for _, e := range events {
			targetDate := t.AddDate(0, 0, e.Offset)
			targetDate = time.Date(targetDate.Year(), targetDate.Month(), targetDate.Day(), 0, 0, 0, 0, targetDate.Location())

			diff := int(targetDate.Sub(now).Hours() / 24)

			status := 3 // Future
			if diff < 0 {
				status = 1 // Past
			} else if diff == 0 {
				status = 2 // Today
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
			resultBox.Add(layout.NewSpacer()) // Spacing antar kartu
		}
		resultBox.Refresh()
		scrollArea.Refresh()
	}

	// --- 5. Footer ---
	noteText := widget.NewLabel("Catatan: Perhitungan menggunakan rumus baku (3, 7, 40, 100, Pendhak, 1000). Pendhak dihitung 354 hari (Kalender Jawa/Islam).")
	noteText.Wrapping = fyne.TextWrapWord
	noteText.TextStyle = fyne.TextStyle{Italic: true}
	
	footerContainer := container.NewPadded(noteText)

	// --- 6. Final Layout Assembly ---
	bgApp := canvas.NewRectangle(ColorBgDark)
	
	topContent := container.NewVBox(
		headerContainer, 
		container.NewPadded(inputSection),
		widget.NewSeparator(),
	)

	mainContent := container.NewBorder(
		topContent,
		footerContainer,
		nil, nil,
		scrollArea,
	)

	myWindow.SetContent(container.NewStack(bgApp, mainContent))
	myWindow.ShowAndRun()
}

