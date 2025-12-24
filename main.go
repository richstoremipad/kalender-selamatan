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

/* ===============================
   1. DATA & KALENDER
================================ */

var (
	HariIndo  = []string{"Minggu", "Senin", "Selasa", "Rabu", "Kamis", "Jumat", "Sabtu"}
	Pasaran   = []string{"Legi", "Pahing", "Pon", "Wage", "Kliwon"}
	BulanIndo = []string{"", "Januari", "Februari", "Maret", "April", "Mei", "Juni", "Juli", "Agustus", "September", "Oktober", "November", "Desember"}
	BulanJawa = []string{"", "Suro", "Sapar", "Mulud", "Bakda Mulud", "Jumadil Awal", "Jumadil Akhir", "Rajeb", "Ruwah", "Poso", "Sawal", "Sela", "Besar"}
)

// Julian Day Number
func dateToJDN(t time.Time) int {
	a := (14 - int(t.Month())) / 12
	y := t.Year() + 4800 - a
	m := int(t.Month()) + 12*a - 3
	return t.Day() + (153*m+2)/5 + 365*y + y/4 - y/100 + y/400 - 32045
}

// Tanggal Jawa
func getJavaneseDate(t time.Time) string {
	jd := dateToJDN(t)
	l := jd - 1948440 + 10632 + 1
	n := (l - 1) / 10631
	l = l - 10631*n + 354
	j := ((10985 - l) / 5316) * ((50 * l) / 17719) + (l / 5670) * ((43 * l) / 15238)
	l = l - ((30 - j) / 15) * ((17719 * j) / 50) - (j / 16) * ((15238 * j) / 43) + 29

	hm := (24 * l) / 709
	hd := l - (709 * hm) / 24

	if hm < 1 || hm >= len(BulanJawa) {
		return "Unknown"
	}
	return fmt.Sprintf("%d %s", hd, BulanJawa[hm])
}

func formatWeton(t time.Time) string {
	hari := HariIndo[t.Weekday()]
	jd := dateToJDN(t)
	pasaran := Pasaran[jd%5]
	return fmt.Sprintf("%s %s, %s", hari, pasaran, getJavaneseDate(t))
}

func formatIndoDate(t time.Time) string {
	return fmt.Sprintf("%d %s %d", t.Day(), BulanIndo[t.Month()], t.Year())
}

/* ===============================
   2. UI STYLE
================================ */

var (
	ColorBgDark     = color.NRGBA{30, 33, 40, 255}
	ColorCardBg     = color.NRGBA{45, 48, 55, 255}
	ColorHeaderTop  = color.NRGBA{40, 180, 160, 255}
	ColorHeaderBot  = color.NRGBA{50, 80, 160, 255}
	ColorTextWhite  = color.NRGBA{255, 255, 255, 255}
	ColorTextGrey   = color.NRGBA{180, 180, 180, 255}
	ColorBadgeGreen = color.NRGBA{46, 125, 50, 255}
	ColorBadgeRed   = color.NRGBA{198, 40, 40, 255}
	ColorBadgeBlue  = color.NRGBA{21, 101, 192, 255}
)

func createCard(title, sub, dateStr, weton string, status, diff int) fyne.CanvasObject {
	var badgeColor color.Color
	var badgeText string

	switch status {
	case 1:
		badgeColor = ColorBadgeGreen
		badgeText = fmt.Sprintf("✓ Sudah Lewat (%d hari)", int(math.Abs(float64(diff))))
	case 2:
		badgeColor = ColorBadgeRed
		badgeText = "🔔 HARI INI!"
	default:
		badgeColor = ColorBadgeBlue
		badgeText = fmt.Sprintf("⏳ %d Hari Lagi", diff)
	}

	titleLbl := canvas.NewText(title, ColorTextWhite)
	titleLbl.TextStyle = fyne.TextStyle{Bold: true}
	subLbl := canvas.NewText(sub, ColorTextGrey)

	dateLbl := canvas.NewText(dateStr, ColorTextWhite)
	dateLbl.Alignment = fyne.TextAlignTrailing
	dateLbl.TextStyle = fyne.TextStyle{Bold: true}
	wetonLbl := canvas.NewText(weton, ColorTextGrey)
	wetonLbl.Alignment = fyne.TextAlignTrailing

	top := container.NewBorder(nil, nil,
		container.NewVBox(titleLbl, subLbl),
		container.NewVBox(dateLbl, wetonLbl),
	)

	badge := canvas.NewRectangle(badgeColor)
	badge.CornerRadius = 12
	badgeTextLbl := canvas.NewText(badgeText, ColorTextWhite)
	badgeTextLbl.TextStyle = fyne.TextStyle{Bold: true}

	content := container.NewVBox(
		top,
		container.NewStack(badge, container.NewPadded(badgeTextLbl)),
	)

	bg := canvas.NewRectangle(ColorCardBg)
	bg.CornerRadius = 10

	return container.NewStack(bg, container.NewPadded(content))
}

/* ===============================
   3. MAIN APP
================================ */

func main() {
	app := app.New()
	win := app.NewWindow("Kalkulator Selamatan Jawa")
	win.Resize(fyne.NewSize(420, 760))

	// Header
	grad := canvas.NewHorizontalGradient(ColorHeaderTop, ColorHeaderBot)
	title := canvas.NewText("Kalkulator Selamatan Jawa", ColorTextWhite)
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.Alignment = fyne.TextAlignCenter

	header := container.NewStack(
		grad,
		container.NewCenter(container.NewHBox(theme.InfoIcon(), title)),
	)

	// Date Picker
	label := canvas.NewText("Pilih Tanggal / Geblag", ColorTextGrey)
	datePicker := widget.NewDatePicker()
	datePicker.SetDate(time.Now())

	btn := widget.NewButton("Hitung", nil)
	input := container.NewBorder(nil, nil, nil, btn, datePicker)

	inputCard := container.NewStack(
		canvas.NewRectangle(ColorCardBg),
		container.NewPadded(container.NewVBox(label, input)),
	)

	// Result
	result := container.NewVBox()
	scroll := container.NewVScroll(container.NewPadded(result))

	btn.OnTapped = func() {
		result.Objects = nil

		t := datePicker.Date
		t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local)

		now := time.Now()
		now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)

		events := []struct {
			Name, Sub string
			Offset    int
		}{
			{"Geblag", "Hari H", 0},
			{"Nelung", "3 Hari", 2},
			{"Mitung", "7 Hari", 6},
			{"Matang", "40 Hari", 39},
			{"Nyatus", "100 Hari", 99},
			{"Pendhak I", "1 Tahun", 353},
			{"Pendhak II", "2 Tahun", 707},
			{"Nyewu", "1000 Hari", 999},
		}

		for _, e := range events {
			d := t.AddDate(0, 0, e.Offset)
			diff := int(d.Sub(now).Hours() / 24)

			status := 3
			if diff < 0 {
				status = 1
			} else if diff == 0 {
				status = 2
			}

			result.Add(createCard(
				e.Name, e.Sub,
				formatIndoDate(d),
				formatWeton(d),
				status, diff,
			))
			result.Add(layout.NewSpacer())
		}
		result.Refresh()
	}

	bg := canvas.NewRectangle(ColorBgDark)
	content := container.NewBorder(
		container.NewVBox(header, container.NewPadded(inputCard)),
		nil, nil, nil,
		scroll,
	)

	win.SetContent(container.NewStack(bg, content))
	win.ShowAndRun()
}
