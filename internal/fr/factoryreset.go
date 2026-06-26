// ============================================================
// factoryreset.go - Factory Reset Protection
// ============================================================
package factoryreset

import (
	"log"
	"runtime"
)

// Disable disables factory reset options
func Disable() {
	if runtime.GOOS != "windows" {
		return
	}
	
	log.Println("🔒 Disabling factory reset options...")
	
	// Disable System Restore
	disableSystemRestore()
	
	// Disable recovery options
	disableRecoveryOptions()
}

// disableSystemRestore disables System Restore
func disableSystemRestore() {
	// In production, modify registry:
	// HKLM\SOFTWARE\Microsoft\Windows NT\CurrentVersion\SystemRestore
	// DisableSR = 1
}

// disableRecoveryOptions disables recovery options
func disableRecoveryOptions() {
	// In production, modify BCD:
	// bcdedit /set {default} recoveryenabled no
}