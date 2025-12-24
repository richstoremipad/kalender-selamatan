import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/dialog" // Tambahkan ini
)

// ... di dalam btnHitung ...
if err != nil {
    dialog.ShowError(fmt.Errorf("Format Salah! Gunakan Tgl-Bln-Thn"), myWindow) // Ganti jadi dialog
    return
}
