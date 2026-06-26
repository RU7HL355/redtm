// ============================================================
// hideconsole.go - Console Hiding
// ============================================================
package hideconsole

import (
	"log"
	"runtime"
	"syscall"
)

// HideConsole hides the console window
func HideConsole() {
	if runtime.GOOS != "windows" {
		return
	}
	
	user32 := syscall.NewLazyDLL("user32.dll")
	getConsoleWindow := user32.NewProc("GetConsoleWindow")
	showWindow := user32.NewProc("ShowWindow")
	
	hwnd, _, _ := getConsoleWindow.Call()
	if hwnd == 0 {
		log.Println("No console window found")
		return
	}
	
	// SW_HIDE = 0
	showWindow.Call(hwnd, 0)
	log.Println("Console window hidden")
}

// ShowConsole shows the console window
func ShowConsole() {
	if runtime.GOOS != "windows" {
		return
	}
	
	user32 := syscall.NewLazyDLL("user32.dll")
	getConsoleWindow := user32.NewProc("GetConsoleWindow")
	showWindow := user32.NewProc("ShowWindow")
	
	hwnd, _, _ := getConsoleWindow.Call()
	if hwnd == 0 {
		return
	}
	
	// SW_SHOW = 5
	showWindow.Call(hwnd, 5)
}