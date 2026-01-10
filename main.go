package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"image/color"
	"math"
	"net/http"
	"net/url"
	"os"
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
// KONFIGURASI VERSI APLIKASI
// ==========================================

// GANTI URL INI DENGAN URL ASLI ANDA
const UpdateCheckURL = "https://raw.githubusercontent.com/richstoremipad/validate/main/version.txt" 

// VERSI APLIKASI SAAT INI
const CurrentAppVersion = "1.0.0"

type UpdateData struct {
	Version     string `json:"version"`
	Title       string `json:"title"`
	Message     string `json:"message"`
	DownloadURL string `json:"download_url"` 
}

// ==========================================
// 1. EMBED RESOURCE (GAMBAR)
// ==========================================

//go:embed rich.png
var richPngData []byte

//go:embed bg.png
var bgPngData []byte

// ==========================================
// 2. DATA PENJELASAN FASE KEMATIAN
// ==========================================

var DeskripsiFase = map[string]string{
	"Geblag": `Rumus: 
Jisarji (dino ke siji pasaran ke siji).
	
Kajian:
Pada hari pertama, jasad mulai mengalami perubahan fisik yang nyata. Ruh digambarkan masih sangat dekat dengan jasadnya dan merasa kaget dengan suasana kubur yang gelap dan sempit.

Kondisi Jasad:
Bagian perut mulai membuncit karena gas mulai terbentuk di dalam usus. Warna kulit yang tadinya cerah berubah menjadi pucat kebiruan atau hijau kehitaman, terutama di area perut dan kemaluan.

Hikmah:
Inilah alasan mengapa keluarga disunnahkan memberikan sedekah pada malam pertama untuk meringankan beban "kagetnya" ruh di alam baru.`,

	"Nelung": `Rumus:
Lusarlu (Dino ke telu pasaran ke telu).
	
Kajian: 
Hari ketiga adalah fase di mana rupa manusia mulai hilang secara perlahan.

Kondisi Jasad:
Cairan mulai keluar dari lubang-lubang tubuh (hidung, mulut, dan telinga). Bau busuk mulai keluar dengan sangat menyengat karena bakteri pembusuk telah menyebar ke seluruh organ dalam.

Kondisi Organ:
Lidah mulai membengkak dan sering kali terjepit oleh gigi karena ruang di dalam mulut menyempit akibat gas. Mata mulai melunak dan tampak agak menonjol.`,

	"Mitung": `Rumus:
Tusarro (Dino ke pitu pasaran ke loro).

Kajian:
Hari ketujuh merupakan fase transisi besar dalam proses penghancuran organ dalam.

Kondisi Jasad:
Perut yang tadinya membuncit akan pecah karena tekanan gas dan aktivitas bakteri. Organ-organ vital seperti hati, paru-paru, dan lambung mulai mencair dan hancur.

Sisi Spiritual:
Berdasarkan keterangan dalam kitab Al-Hawi lil Fatawi (Imam Suyuthi) yang sering disandingkan dengan Daqa'iqul Akhbar, tujuh hari pertama adalah masa Fitnah Kubur (ujian dan pertanyaan malaikat). Oleh karena itu, sedekah makanan pada hari ke-7 sangat ditekankan.`,

	"Matang": `Rumus:
Masarma (Dino ke limo pasaran ke limo).
	
Kajian:
Pada hari ke-40, jasad sudah tidak lagi menyerupai sosok manusia yang dikenal semasa hidup.

Kondisi Jasad:
Seluruh daging mulai terlepas dari tulang belulang. Daging-daging tersebut mulai meluruh dan menyatu dengan tanah.

Kondisi Wajah:
Kulit wajah sudah hancur sepenuhnya, mata sudah hilang dari kelopaknya, dan rambut mulai rontok dari kulit kepala.

Tradisi:
Dipercaya pada hari ke-40, proses "pembersihan" sisa daging sedang terjadi secara masif, sehingga doa dikirimkan agar ruh diberikan ketenangan dalam melihat jasadnya yang hancur.`,

	"Nyatus": `Rumus:
Rosarma (Dino ke loro pasaran ke limo).

Kajian:
Memasuki hari ke-100, proses pembusukan daging sudah hampir selesai secara total.

Kondisi Jasad:
Tubuh kini didominasi oleh rangka. Hanya menyisakan sedikit jaringan otot atau kulit yang mengeras (seperti mumi) di area-area tertentu yang sulit hancur.

Bau:
Bau busuk yang menyengat sudah mulai berkurang karena sumber pembusukan (daging dan organ dalam) sudah menyatu dengan tanah.`,

	"Pendhak I": `Rumus:
Patsarpat (Dino ke papat pasaran ke papat).

Kajian: 
Istilah "Pendhak" adalah tradisi lokal Nusantara untuk menyebut Haul atau peringatan tahunan.

Kondisi Jasad:
Tulang-belulang mulai menjadi kering. Sumsum di dalam tulang sudah habis. Sendi-sendi yang menghubungkan tulang satu dengan yang lain mulai terlepas.

Kondisi Tengkorak:
Rahang bawah biasanya sudah terlepas dari tengkorak. Tubuh benar-benar sudah menjadi serpihan tulang yang terpisah-pisah.`,

	"Pendhak II": `Rumus:
Rosarpat (Dino ke siji pasaran ke telu).

Kajian:
Memasuki tahun kedua, proses dekomposisi tulang berlanjut.

Kondisi Jasad:
Tulang-belulang semakin kering dan mulai terurai oleh tanah. Sendi-sendi utama sudah lepas sepenuhnya. Struktur kerangka tubuh sudah tidak utuh lagi.

Makna:
Peringatan ini menjadi penanda bahwa hubungan fisik almarhum dengan dunia semakin pudar, dan yang tersisa hanyalah doa dari anak cucu serta amal jariyahnya.`,

	"Nyewu": `Rumus:
Nemsarma (Dino ke enem pasaran ke limo).
	
Kajian:
Ini adalah fase terakhir dalam proses dekomposisi jasad manusia secara alami.

Kondisi Jasad:
Tulang-belulang mulai melapuk dan menjadi rapuh. Dalam kitab dijelaskan bahwa pada fase ini, jasad sudah benar-benar menyatu dengan tanah (menjadi debu).

Satu Bagian yang Tersisa:
Dalam keyakinan Islam (berdasarkan Hadis Nabi), hanya satu bagian yang tidak akan hancur dimakan tanah, yaitu Ajbuz Dzamb (tulang ekor yang sangat kecil), yang darinya manusia akan dibangkitkan kembali pada hari kiamat.

Makna Doa:
Peringatan 1000 hari dimaksudkan sebagai doa pamungkas bagi keluarga untuk memohonkan ampunan total bagi almarhum/ah karena perjalanan jasadnya di bumi sudah selesai secara fisik.`,
}

// ==========================================
// 3. LOGIKA MATEMATIKA & KALENDER JAWA
// ==========================================

var (
	HariIndo     = []string{"Minggu", "Senin", "Selasa", "Rabu", "Kamis", "Jumat", "Sabtu"}
	Pasaran      = []string{"Legi", "Pahing", "Pon", "Wage", "Kliwon"}
	BulanIndo    = []string{"", "Januari", "Februari", "Maret", "April", "Mei", "Juni", "Juli", "Agustus", "September", "Oktober", "November", "Desember"}
	BulanJawa    = []string{"", "Suro", "Sapar", "Mulud", "Bakda Mulud", "Jumadil Awal", "Jumadil Akhir", "Rajeb", "Ruwah", "Poso", "Sawal", "Sela", "Besar"}
	NilaiHari    = []int{5, 4, 3, 7, 8, 6, 9}
	NilaiPasaran = []int{5, 9, 7, 4, 8}
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

func calculateNeptu(t time.Time) string {
	idxHari := int(t.Weekday())
	valHari := NilaiHari[idxHari]
	jd := dateToJDN(t)
	idxPasaran := jd % 5
	valPasaran := NilaiPasaran[idxPasaran]
	total := valHari + valPasaran
	return fmt.Sprintf("Jumlah Neptu: %d", total)
}

// ==========================================
// 4. KOMPONEN UI CUSTOM & COLORS
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
	ColorTextOrange = color.NRGBA{R: 255, G: 165, B: 0, A: 255}
)

type myTheme struct {
	fyne.Theme
}

func (m myTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	if name == "orange" {
		return ColorTextOrange
	}
	if name == "red" {
		return ColorBadgeRed
	}
	if name == theme.ColorNamePrimary {
		return ColorBadgeGreen
	}
	if name == theme.ColorNameError {
		return ColorBadgeRed
	}
	if name == theme.ColorNameButton {
		return color.NRGBA{R: 60, G: 63, B: 70, A: 255}
	}
	return m.Theme.Color(name, variant)
}

// --- WIDGET KLIKABLE CUSTOM ---
type clickableCard struct {
	widget.BaseWidget
	content fyne.CanvasObject
	onTap   func()
}

func newClickableCard(content fyne.CanvasObject, onTap func()) *clickableCard {
	c := &clickableCard{content: content, onTap: onTap}
	c.ExtendBaseWidget(c)
	return c
}

func (c *clickableCard) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(c.content)
}

func (c *clickableCard) Tapped(_ *fyne.PointEvent) {
	if c.onTap != nil {
		c.onTap()
	}
}

// ==========================================
// 5. LOGIKA KALENDER CUSTOM
// ==========================================

func createCalendarPopup(parentCanvas fyne.Canvas, initialDate time.Time, onDateChanged func(time.Time), onCalculate func(time.Time)) {
	currentMonth := initialDate
	selectedDate := initialDate
	currentViewMode := 0
	hasSelected := false

	contentStack := container.NewStack()
	var popup *widget.PopUp

	toastText := canvas.NewText("Pilih tanggal dulu!", ColorTextWhite)
	toastText.TextSize = 14
	toastText.TextStyle = fyne.TextStyle{Bold: true}
	toastBg := canvas.NewRectangle(color.NRGBA{R: 0, G: 0, B: 0, A: 200})
	toastBg.CornerRadius = 8
	toastCard := container.NewStack(toastBg, container.NewPadded(toastText))
	toastWrapper := container.NewCenter(toastCard)
	toastWrapper.Hide()

	showToast := func() {
		toastWrapper.Show()
		go func() {
			time.Sleep(2 * time.Second)
			toastWrapper.Hide()
		}()
	}

	var refreshContent func()
	refreshContent = func() {
		year, month, _ := currentMonth.Date()

		if currentViewMode == 0 {
			titleText := fmt.Sprintf("%s %d", BulanIndo[month], year)
			btnHeader := widget.NewButton(titleText, func() {
				currentViewMode = 1
				refreshContent()
			})
			btnHeader.Importance = widget.LowImportance

			btnPrev := widget.NewButtonWithIcon("", theme.NavigateBackIcon(), func() {
				currentMonth = currentMonth.AddDate(0, -1, 0)
				refreshContent()
			})
			btnNext := widget.NewButtonWithIcon("", theme.NavigateNextIcon(), func() {
				currentMonth = currentMonth.AddDate(0, 1, 0)
				refreshContent()
			})

			topNav := container.NewBorder(nil, nil, btnPrev, btnNext, container.NewCenter(btnHeader))

			gridDays := container.New(layout.NewGridLayout(7))
			daysHeader := []string{"M", "S", "S", "R", "K", "J", "S"}
			for _, dayName := range daysHeader {
				l := widget.NewLabel(dayName)
				l.Alignment = fyne.TextAlignCenter
				l.TextStyle = fyne.TextStyle{Bold: true}
				gridDays.Add(l)
			}

			gridDates := container.New(layout.NewGridLayout(7))
			firstDayOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
			startWeekday := int(firstDayOfMonth.Weekday())
			nextMonth := firstDayOfMonth.AddDate(0, 1, 0)
			lastDay := nextMonth.Add(-time.Hour * 24).Day()

			for i := 0; i < startWeekday; i++ {
				gridDates.Add(layout.NewSpacer())
			}
			for d := 1; d <= lastDay; d++ {
				dayNum := d
				dateVal := time.Date(year, month, dayNum, 0, 0, 0, 0, time.Local)
				btn := widget.NewButton(fmt.Sprintf("%d", dayNum), nil)
				if hasSelected &&
					dateVal.Year() == selectedDate.Year() &&
					dateVal.Month() == selectedDate.Month() &&
					dateVal.Day() == selectedDate.Day() {
					btn.Importance = widget.HighImportance
				} else {
					btn.Importance = widget.MediumImportance
				}
				btn.OnTapped = func() {
					selectedDate = dateVal
					hasSelected = true
					refreshContent()
					if onDateChanged != nil {
						onDateChanged(selectedDate)
					}
				}
				gridDates.Add(btn)
			}
			contentStack.Objects = []fyne.CanvasObject{
				container.NewVBox(topNav, gridDays, gridDates),
			}

		} else if currentViewMode == 1 {
			btnBack := widget.NewButtonWithIcon("", theme.NavigateBackIcon(), func() {
				currentViewMode = 0
				refreshContent()
			})
			btnBack.Importance = widget.DangerImportance

			btnPrevYear := widget.NewButtonWithIcon("", theme.NavigateBackIcon(), func() {
				currentMonth = currentMonth.AddDate(-1, 0, 0)
				refreshContent()
			})

			btnNextYear := widget.NewButtonWithIcon("", theme.NavigateNextIcon(), func() {
				currentMonth = currentMonth.AddDate(1, 0, 0)
				refreshContent()
			})

			btnYearNum := widget.NewButton(fmt.Sprintf("%d", year), func() {
				currentViewMode = 2
				refreshContent()
			})
			btnYearNum.Importance = widget.LowImportance

			yearNavLayout := container.NewBorder(nil, nil, btnPrevYear, btnNextYear, container.NewCenter(btnYearNum))

			monthGrid := container.New(layout.NewGridLayout(3))
			for i := 1; i <= 12; i++ {
				mIdx := i
				mName := BulanIndo[mIdx]
				if len(mName) > 3 {
					mName = mName[:3]
				}
				btnMonth := widget.NewButton(mName, func() {
					currentMonth = time.Date(currentMonth.Year(), time.Month(mIdx), 1, 0, 0, 0, 0, time.Local)
					currentViewMode = 0
					refreshContent()
				})
				if time.Month(mIdx) == month {
					btnMonth.Importance = widget.HighImportance
				} else {
					btnMonth.Importance = widget.MediumImportance
				}
				monthGrid.Add(container.NewCenter(btnMonth))
			}

			topRow := container.NewHBox(container.NewCenter(btnBack), layout.NewSpacer())
			contentStack.Objects = []fyne.CanvasObject{
				container.NewVBox(topRow, container.NewPadded(yearNavLayout), monthGrid),
			}

		} else {
			btnBack := widget.NewButtonWithIcon("", theme.NavigateBackIcon(), func() {
				currentViewMode = 1
				refreshContent()
			})
			btnBack.Importance = widget.DangerImportance

			startYear := 1900
			endYear := 2100
			totalYears := endYear - startYear + 1

			actualIndex := year - startYear
			scrollIndex := actualIndex - 3

			if scrollIndex < 0 {
				scrollIndex = 0
			}
			if scrollIndex >= totalYears {
				scrollIndex = totalYears - 1
			}

			list := widget.NewList(
				func() int {
					return totalYears
				},
				func() fyne.CanvasObject {
					btn := widget.NewButton("Template", nil)
					btn.Alignment = widget.ButtonAlignCenter
					return btn
				},
				func(i widget.ListItemID, o fyne.CanvasObject) {
					displayYear := startYear + i
					btn := o.(*widget.Button)
					btn.SetText(fmt.Sprintf("%d", displayYear))
					btn.Importance = widget.MediumImportance

					btn.OnTapped = func() {
						currentMonth = time.Date(displayYear, month, 1, 0, 0, 0, 0, time.Local)
						currentViewMode = 1
						refreshContent()
					}
				},
			)

			listContainer := container.NewStack(list)

			topRow := container.NewBorder(nil, nil, btnBack, nil, nil)

			yearView := container.NewBorder(
				container.NewPadded(topRow),
				nil, nil, nil,
				listContainer,
			)

			contentStack.Objects = []fyne.CanvasObject{yearView}

			go func() {
				time.Sleep(100 * time.Millisecond)
				list.ScrollTo(widget.ListItemID(scrollIndex))
			}()
		}

		contentStack.Refresh()
	}

	btnHitung := widget.NewButton("Hitung", func() {
		if currentViewMode != 0 {
			showToast()
			return
		}

		if !hasSelected {
			showToast()
			return
		}
		if popup != nil {
			popup.Hide()
		}
		onCalculate(selectedDate)
	})
	btnHitung.Importance = widget.HighImportance
	btnHitung.Icon = theme.ConfirmIcon()
	bottomArea := container.NewCenter(btnHitung)

	refreshContent()

	finalLayout := container.NewBorder(
		nil,
		container.NewPadded(bottomArea),
		nil, nil,
		contentStack,
	)

	bgRect := canvas.NewRectangle(ColorCardBg)
	bgRect.CornerRadius = 12
	bgRect.SetMinSize(fyne.NewSize(280, 330))

	cardContent := container.NewStack(
		bgRect,
		container.NewPadded(finalLayout),
		toastWrapper,
	)
	centeredPopup := container.NewCenter(cardContent)

	popup = widget.NewModalPopUp(centeredPopup, parentCanvas)
	popup.Resize(fyne.NewSize(280, 330))
	popup.Show()
}

// ==========================================
// 6. HELPER UI CARDS
// ==========================================

func createCard(title, subTitle, dateStr, wetonStr, rumusStr, descStr string, statusType int, diffDays int, parentCanvas fyne.Canvas) fyne.CanvasObject {
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

	var rightCont *fyne.Container
	if rumusStr != "" {
		lblRumus := canvas.NewText(rumusStr, ColorTextOrange)
		lblRumus.Alignment = fyne.TextAlignTrailing
		lblRumus.TextSize = 10
		lblRumus.TextStyle = fyne.TextStyle{Italic: true}
		rightCont = container.NewVBox(lblDate, lblWeton, lblRumus)
	} else {
		rightCont = container.NewVBox(lblDate, lblWeton)
	}

	topRow := container.NewBorder(nil, nil, leftCont, rightCont)

	var botRow fyne.CanvasObject
	if statusType >= 1 && statusType <= 3 {
		lblBadge := canvas.NewText(badgeTextStr, ColorTextWhite)
		lblBadge.TextSize = 11
		lblBadge.TextStyle = fyne.TextStyle{Bold: true}
		badgeBg := canvas.NewRectangle(badgeColor)
		badgeBg.CornerRadius = 12
		badgeCont := container.NewStack(badgeBg, container.NewPadded(lblBadge))
		botRow = container.NewHBox(badgeCont)
	} else {
		botRow = layout.NewSpacer()
	}

	content := container.NewVBox(topRow, container.NewPadded(botRow))
	bg := canvas.NewRectangle(ColorCardBg)
	bg.CornerRadius = 10

	visualCard := container.NewStack(bg, container.NewPadded(content))

	if descStr != "" && parentCanvas != nil {
		return newClickableCard(visualCard, func() {
			lblDesc := widget.NewLabel(descStr)
			lblDesc.Wrapping = fyne.TextWrapWord

			lblHeader := widget.NewLabel("Penjelasan Fase: " + title)
			lblHeader.Alignment = fyne.TextAlignCenter
			lblHeader.TextStyle = fyne.TextStyle{Bold: true}

			var popup *widget.PopUp
			btnClose := widget.NewButton("Tutup", func() {
				if popup != nil {
					popup.Hide()
				}
			})
			btnClose.Importance = widget.HighImportance

			scrollContainer := container.NewVScroll(container.NewPadded(lblDesc))
			scrollContainer.SetMinSize(fyne.NewSize(0, 300))

			contentBox := container.NewBorder(
				lblHeader,
				container.NewPadded(btnClose),
				nil, nil,
				scrollContainer,
			)

			bgRect := canvas.NewRectangle(ColorCardBg)
			bgRect.CornerRadius = 12
			bgRect.SetMinSize(fyne.NewSize(300, 400))

			finalPopupContent := container.NewStack(bgRect, container.NewPadded(contentBox))

			popup = widget.NewModalPopUp(container.NewCenter(finalPopupContent), parentCanvas)
			popup.Resize(fyne.NewSize(320, 450))
			popup.Show()
		})
	}
	return visualCard
}

// ==========================================
// 7. MAIN APP
// ==========================================

// --- UPDATE CHECKER LOGIC (MAXIMUM ANTI-CACHE) ---
func checkForUpdates(myCanvas fyne.Canvas, myApp fyne.App) {
	go func() {
		// Logika: Tunggu aplikasi tampil (3 detik), lalu cek.
		time.Sleep(3 * time.Second)

		fmt.Println("--- Memulai Pengecekan Update ---")

		var updateInfo UpdateData

		// 1. URL Unik (Cache Busting)
		// Menambahkan parameter ?t=123456789 agar server menganggap ini request baru
		uniqueURL := fmt.Sprintf("%s?t=%d", UpdateCheckURL, time.Now().UnixNano())
		fmt.Println("Target URL:", uniqueURL)

		// 2. HTTP Request dengan Header Agresif
		req, err := http.NewRequest("GET", uniqueURL, nil)
		if err != nil {
			fmt.Println("Gagal membuat request:", err)
			return
		}
		// Header standar untuk memaksa server memberikan data terbaru
		req.Header.Set("Cache-Control", "no-cache, no-store, must-revalidate")
		req.Header.Set("Pragma", "no-cache")
		req.Header.Set("Expires", "0")

		// 3. Transport Layer (Mematikan Cache Internal Go)
		tr := &http.Transport{
			DisableKeepAlives: true, // Jangan pakai koneksi lama
		}
		client := &http.Client{
			Transport: tr,
			Timeout:   10 * time.Second,
		}

		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Gagal koneksi (Offline?):", err)
			return 
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Println("Server merespon error:", resp.StatusCode)
			return
		}

		if err := json.NewDecoder(resp.Body).Decode(&updateInfo); err != nil {
			fmt.Println("Gagal decode JSON:", err)
			return
		}

		fmt.Printf("Data Server: %+v\n", updateInfo)
		fmt.Printf("Versi App: %s, Versi Server: %s\n", CurrentAppVersion, updateInfo.Version)

		// 4. Bandingkan Versi
		if updateInfo.Version == CurrentAppVersion {
			fmt.Println("Versi sudah paling baru. Tidak menampilkan popup.")
			return 
		}

		fmt.Println("Update ditemukan! Menampilkan Popup...")

		// 5. TAMPILKAN POPUP
		
		lblTitle := canvas.NewText(updateInfo.Title, ColorTextWhite)
		lblTitle.TextStyle = fyne.TextStyle{Bold: true}
		lblTitle.TextSize = 16
		lblTitle.Alignment = fyne.TextAlignCenter

		lblVer := canvas.NewText("Versi Baru: "+updateInfo.Version, ColorTextWhite)
		lblVer.TextSize = 12
		badgeBg := canvas.NewRectangle(ColorBadgeGreen)
		badgeBg.CornerRadius = 8
		badgeVer := container.NewStack(badgeBg, container.NewPadded(lblVer))

		msgText := widget.NewRichTextFromMarkdown(updateInfo.Message)
		msgText.Wrapping = fyne.TextWrapWord

		// Tombol Keluar (KIRI) - AMAN DENGAN OS.EXIT
		btnExit := widget.NewButton("Keluar", func() {
			os.Exit(0) 
		})
		btnExit.Importance = widget.DangerImportance 

		// Tombol Update (KANAN)
		btnUpdate := widget.NewButton("Update", func() {
			u, err := url.Parse(updateInfo.DownloadURL)
			if err == nil {
				myApp.OpenURL(u)
			}
		})
		btnUpdate.Importance = widget.HighImportance 

		// Layout Tombol: [Keluar] <--- SPASI ---> [Update]
		buttonRow := container.NewHBox(
			btnExit,
			layout.NewSpacer(), 
			btnUpdate,
		)

		mainContent := container.NewVBox(
			lblTitle,
			container.NewCenter(badgeVer),
			widget.NewSeparator(),
			msgText,
		)

		finalLayout := container.NewBorder(
			nil,                         
			container.NewPadded(buttonRow), 
			nil,                          
			nil,                         
			container.NewPadded(mainContent),
		)

		bgRect := canvas.NewRectangle(ColorCardBg)
		bgRect.CornerRadius = 12
		bgRect.SetMinSize(fyne.NewSize(280, 250))

		popupContent := container.NewStack(
			bgRect,
			container.NewPadded(finalLayout),
		)

		popup := widget.NewModalPopUp(container.NewCenter(popupContent), myCanvas)
		popup.Resize(fyne.NewSize(320, 300))
		popup.Show()

	}()
}

func main() {
	myApp := app.New()
	myApp.Settings().SetTheme(&myTheme{Theme: theme.DefaultTheme()})

	myWindow := myApp.NewWindow("Kalkulator Selamatan Jawa & Weton")
	myWindow.Resize(fyne.NewSize(400, 750))

	resBg := fyne.NewStaticResource("bg.png", bgPngData)
	imgBg := canvas.NewImageFromResource(resBg)
	imgBg.FillMode = canvas.ImageFillCover

	gradient := canvas.NewHorizontalGradient(ColorHeaderTop, ColorHeaderBot)
	headerTitle := canvas.NewText("Kalkulator Selamatan & Weton", ColorTextWhite)
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

	// =======================================================
	// BAGIAN 1: TAB SELAMATAN
	// =======================================================

	resultBox := container.NewVBox()
	scrollArea := container.NewVScroll(container.NewPadded(resultBox))

	calcDate := time.Now()
	lblDateTitle := canvas.NewText("Tanggal Wafat / Geblag:", ColorTextGrey)
	lblDateTitle.TextSize = 12

	lblSelectedDate := widget.NewLabel("Belum dipilih")
	lblSelectedDate.Alignment = fyne.TextAlignCenter
	lblSelectedDate.TextStyle = fyne.TextStyle{Bold: true}

	// Helper update text
	updateDateLabel := func(t time.Time) {
		lblSelectedDate.SetText(formatIndoDate(t))
	}
	updateDateLabel(calcDate)

	performCalculation := func(t time.Time) {
		updateDateLabel(t)
		resultBox.Objects = nil

		type Event struct {
			Name   string
			Sub    string
			Offset int
			Rumus  string
		}
		events := []Event{
			{"Geblag", "Hari H", 0, "Jisarji"},
			{"Nelung", "3 Hari", 2, "Lusarlu"},
			{"Mitung", "7 Hari", 6, "Tusarro"},
			{"Matang", "40 Hari", 39, "Masarma"},
			{"Nyatus", "100 Hari", 99, "Rosarma"},
			{"Pendhak I", "1 Tahun", 353, "Patsarpat"},
			{"Pendhak II", "2 Tahun", 707, "Rosarpat"},
			{"Nyewu", "1000 Hari", 999, "Nemsarmo"},
		}

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
			desc := DeskripsiFase[e.Name]
			card := createCard(e.Name, e.Sub, formatIndoDate(targetDate), formatWeton(targetDate), e.Rumus, desc, status, diff, myWindow.Canvas())
			resultBox.Add(card)
			resultBox.Add(layout.NewSpacer())
		}
		resultBox.Refresh()
	}

	btnOpenCalc := widget.NewButton("Hitung Selamatan", nil)
	btnOpenCalc.Importance = widget.HighImportance
	btnOpenCalc.Icon = theme.CalendarIcon()

	// CALL CREATE CALENDAR POPUP DENGAN REALTIME CALLBACK
	btnOpenCalc.OnTapped = func() {
		createCalendarPopup(myWindow.Canvas(), calcDate,
			// Callback 1: Realtime Update
			func(realtimeDate time.Time) {
				updateDateLabel(realtimeDate)
			},
			// Callback 2: Final Selection
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
			container.NewCenter(btnOpenCalc),
		)),
	)

	tabContentSelamatan := container.NewBorder(
		container.NewPadded(inputSection),
		nil, nil, nil,
		scrollArea,
	)

	// =======================================================
	// BAGIAN 2: TAB CEK WETON
	// =======================================================

	wetonResultBox := container.NewVBox()
	wetonScrollArea := container.NewVScroll(container.NewPadded(wetonResultBox))

	wetonDate := time.Now()
	lblWetonTitle := canvas.NewText("Tanggal Lahir:", ColorTextGrey)
	lblWetonTitle.TextSize = 12

	lblSelectedWetonDate := widget.NewLabel("Belum dipilih")
	lblSelectedWetonDate.Alignment = fyne.TextAlignCenter
	lblSelectedWetonDate.TextStyle = fyne.TextStyle{Bold: true}

	updateWetonDateLabel := func(t time.Time) {
		lblSelectedWetonDate.SetText(formatIndoDate(t))
	}
	updateWetonDateLabel(wetonDate)

	performWetonCheck := func(t time.Time) {
		updateWetonDateLabel(t)
		wetonResultBox.Objects = nil
		t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
		neptuStr := calculateNeptu(t)
		card := createCard("Hasil Weton", neptuStr, formatIndoDate(t), formatWeton(t), "", "", 4, 0, nil)
		wetonResultBox.Add(card)
		wetonResultBox.Refresh()
	}

	btnOpenWeton := widget.NewButton("Pilih Tanggal Lahir", nil)
	btnOpenWeton.Importance = widget.HighImportance
	btnOpenWeton.Icon = theme.AccountIcon()

	// CALL CREATE CALENDAR POPUP DENGAN REALTIME CALLBACK (WETON)
	btnOpenWeton.OnTapped = func() {
		createCalendarPopup(myWindow.Canvas(), wetonDate,
			// Callback 1: Realtime Update
			func(realtimeDate time.Time) {
				updateWetonDateLabel(realtimeDate)
			},
			// Callback 2: Final Selection
			func(finalDate time.Time) {
				wetonDate = finalDate
				performWetonCheck(wetonDate)
			},
		)
	}

	inputRowWeton := container.NewBorder(nil, nil, nil, nil, lblSelectedWetonDate)
	inputCardBgWeton := canvas.NewRectangle(ColorCardBg)
	inputCardBgWeton.CornerRadius = 8

	inputSectionWeton := container.NewStack(
		inputCardBgWeton,
		container.NewPadded(container.NewVBox(
			lblWetonTitle,
			inputRowWeton,
			layout.NewSpacer(),
			container.NewCenter(btnOpenWeton),
		)),
	)

	tabContentWeton := container.NewBorder(
		container.NewPadded(inputSectionWeton),
		nil, nil, nil,
		wetonScrollArea,
	)

	// =======================================================
	// FOOTER SETUP
	// =======================================================

	richNoteSelamatan := widget.NewRichText(
		&widget.TextSegment{
			Text: "Notes: ",
			Style: widget.RichTextStyle{
				ColorName: "orange",
				Inline:    true,
				TextStyle: fyne.TextStyle{Italic: true, Bold: true},
			},
		},
		&widget.TextSegment{
			Text: "Perhitungan ini menggunakan rumus ",
			Style: widget.RichTextStyle{
				Inline:    true,
				TextStyle: fyne.TextStyle{Italic: true},
			},
		},
		&widget.TextSegment{
			Text: "lusarlu ",
			Style: widget.RichTextStyle{
				ColorName: "red",
				Inline:    true,
				TextStyle: fyne.TextStyle{Italic: true, Bold: true},
			},
		},
		&widget.TextSegment{
			Text: "hingga ",
			Style: widget.RichTextStyle{
				Inline:    true,
				TextStyle: fyne.TextStyle{Italic: true},
			},
		},
		&widget.TextSegment{
			Text: "nemsarmo ",
			Style: widget.RichTextStyle{
				ColorName: "red",
				Inline:    true,
				TextStyle: fyne.TextStyle{Italic: true, Bold: true},
			},
		},
		&widget.TextSegment{
			Text: ". Silahkan klik pada hasil hari/pasaran untuk melihat rumus dan filosofinya.",
			Style: widget.RichTextStyle{
				Inline:    true,
				TextStyle: fyne.TextStyle{Italic: true},
			},
		},
	)
	richNoteSelamatan.Wrapping = fyne.TextWrapWord

	richNoteWeton := widget.NewRichText(
		&widget.TextSegment{
			Text: "Notes: ",
			Style: widget.RichTextStyle{ColorName: "orange", Inline: true, TextStyle: fyne.TextStyle{Italic: true, Bold: true}},
		},
		&widget.TextSegment{
			Text: "Perhitungan Weton ini menjumlahkan neptu ",
			Style: widget.RichTextStyle{Inline: true, TextStyle: fyne.TextStyle{Italic: true}},
		},
		&widget.TextSegment{
			Text: "Hari dan Pasaran ",
			Style: widget.RichTextStyle{ColorName: "primary", Inline: true, TextStyle: fyne.TextStyle{Italic: true, Bold: true}},
		},
		&widget.TextSegment{
			Text: "sesuai pakem Primbon Jawa.",
			Style: widget.RichTextStyle{Inline: true, TextStyle: fyne.TextStyle{Italic: true}},
		},
	)
	richNoteWeton.Wrapping = fyne.TextWrapWord

	noteContainer := container.NewStack()
	noteContainer.Add(richNoteSelamatan)

	resRich := fyne.NewStaticResource("rich.png", richPngData)
	imgCredit := canvas.NewImageFromResource(resRich)
	imgCredit.FillMode = canvas.ImageFillContain
	imgCredit.SetMinSize(fyne.NewSize(150, 50))

	footerContent := container.NewVBox(noteContainer, container.NewCenter(imgCredit))

	footerCardBg := canvas.NewRectangle(ColorCardBg)
	footerCardBg.CornerRadius = 8
	footerSection := container.NewStack(
		footerCardBg,
		container.NewPadded(footerContent),
	)

	// =======================================================
	// TAB CONTROL
	// =======================================================

	tabs := container.NewAppTabs(
		container.NewTabItem("Hitung Selamatan", tabContentSelamatan),
		container.NewTabItem("Cek Weton Lahir", tabContentWeton),
	)
	tabs.SetTabLocation(container.TabLocationTop)

	tabs.OnSelected = func(i *container.TabItem) {
		noteContainer.Objects = nil
		if i.Text == "Hitung Selamatan" {
			noteContainer.Add(richNoteSelamatan)
		} else {
			noteContainer.Add(richNoteWeton)
		}
		noteContainer.Refresh()
	}

	mainContent := container.NewBorder(
		headerContainer,
		container.NewPadded(footerSection),
		nil, nil,
		container.NewPadded(tabs),
	)

	myWindow.SetContent(container.NewStack(imgBg, mainContent))

	// ========================================
	// PANGGIL FUNGSI CEK UPDATE DI SINI
	// ========================================
	checkForUpdates(myWindow.Canvas(), myApp)

	myWindow.ShowAndRun()
}
