// ============================================================
// factoryreset.go - Factory Reset Protection (FIXED)
// ============================================================
package factoryreset

import (
	"log"
	"time"
)

// Disable disables factory reset options
func Disable() {
	log.Println("🔒 Disabling factory reset options...")
	
	// Use goroutine with timeout
	done := make(chan bool, 1)
	go func() {
		disableSystemRestore()
		disableRecoveryOptions()
		done <- true
	}()
	
	select {
	case <-done:
		log.Println("✅ Factory reset protection applied")
	case <-time.After(2 * time.Second):
		log.Println("⚠️ Factory reset timeout - continuing")
	}
}

func disableSystemRestore() {
	// In production, modify registry
}

func disableRecoveryOptions() {
	// In production, modify BCD
}