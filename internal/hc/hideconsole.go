// ============================================================
// hideconsole.go - Console Hiding (Windows)
// ============================================================
package hideconsole

import (
	"log"
	"syscall"
)

// HideConsole hides the console window
func HideConsole() {
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

// IsConsoleHidden checks if console is hidden
func IsConsoleHidden() bool {
	user32 := syscall.NewLazyDLL("user32.dll")
	getConsoleWindow := user32.NewProc("GetConsoleWindow")
	isWindowVisible := user32.NewProc("IsWindowVisible")
	
	hwnd, _, _ := getConsoleWindow.Call()
	if hwnd == 0 {
		return true
	}
	
	ret, _, _ := isWindowVisible.Call(hwnd)
	return ret == 0
}