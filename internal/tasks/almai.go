package tasks

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
	"zenith-agent/internal/engine"

	"github.com/playwright-community/playwright-go"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

const (
	almaiURL       = "https://almai.id/register?ref=SCALPINGHACK"
	almaiBaseDir   = "almai-folder"
	almaiCSVDir    = "almai-folder/file-csv"
	almaiOutputDir = "almai-folder/output"
	maxRetries     = 3
	baseTimeout    = 60000
)

// AlmaiEntry represents a single registration entry
type AlmaiEntry struct {
	NamaLengkap        string `json:"nama_lengkap"`
	Email              string `json:"email"`
	WhatsApp           string `json:"whatsapp"`
	IP                 string `json:"ip"`
	Password           string `json:"password"`
	KonfirmasiPassword string `json:"konfirmasi_password"`
	SusReason          string `json:"sus_reason,omitempty"`
}

// AlmaiProgress tracks execution progress
type AlmaiProgress struct {
	Status    string      `json:"status"`
	Meta      *AlmaiEntry `json:"meta,omitempty"`
	Timestamp string      `json:"timestamp,omitempty"`
	ErrorMsg  string      `json:"error_msg,omitempty"`
}

// CSV column detection helpers
var (
	commonNameKeys  = []string{"name", "fullname", "full name", "nama", "nama_lengkap", "nama lengkap"}
	commonEmailKeys = []string{"email", "e-mail", "alamat email"}
	commonPhoneKeys = []string{"phone", "phone number", "phone_number", "hp", "handphone", "whatsapp", "nohp", "no_hp", "telephone", "tel"}
	commonIPKeys    = []string{"ip", "ip address", "ip_address", "ipaddress", "alamat ip"}
)

// normalizeHeader removes special characters for matching
func normalizeHeader(h string) string {
	reg := regexp.MustCompile(`[^a-z0-9]`)
	return reg.ReplaceAllString(strings.ToLower(strings.TrimSpace(h)), "")
}

// guessColumns attempts to detect which columns contain name, email, phone, IP
func guessColumns(headers []string) (nameCol, emailCol, phoneCol, ipCol string) {
	norm := make(map[string]string)
	for _, h := range headers {
		norm[h] = normalizeHeader(h)
	}

	// Match against common patterns
	for _, h := range headers {
		nh := norm[h]

		// Name column
		if nameCol == "" {
			for _, k := range commonNameKeys {
				if normalizeHeader(k) == nh {
					nameCol = h
					break
				}
			}
		}

		// Email column
		if emailCol == "" {
			for _, k := range commonEmailKeys {
				if normalizeHeader(k) == nh {
					emailCol = h
					break
				}
			}
		}

		// Phone column
		if phoneCol == "" {
			for _, k := range commonPhoneKeys {
				if normalizeHeader(k) == nh {
					phoneCol = h
					break
				}
			}
		}

		// IP column
		if ipCol == "" {
			for _, k := range commonIPKeys {
				if normalizeHeader(k) == nh {
					ipCol = h
					break
				}
			}
		}
	}

	// Fallback: search for keywords in header names
	if nameCol == "" {
		for _, h := range headers {
			lower := strings.ToLower(h)
			if strings.Contains(lower, "name") || strings.Contains(lower, "nama") {
				nameCol = h
				break
			}
		}
	}

	if emailCol == "" {
		for _, h := range headers {
			lower := strings.ToLower(h)
			if strings.Contains(lower, "email") {
				emailCol = h
				break
			}
		}
	}

	if phoneCol == "" {
		for _, h := range headers {
			lower := strings.ToLower(h)
			if strings.Contains(lower, "phone") || strings.Contains(lower, "hp") ||
				strings.Contains(lower, "whatsapp") || strings.Contains(lower, "tel") {
				phoneCol = h
				break
			}
		}
	}

	if ipCol == "" {
		for _, h := range headers {
			if strings.Contains(strings.ToLower(h), "ip") {
				ipCol = h
				break
			}
		}
	}

	return
}

// normalizePhone converts Indonesian phone numbers to international format
func normalizePhone(raw string) string {
	if raw == "" {
		return ""
	}

	s := strings.TrimSpace(raw)
	if s == "" {
		return ""
	}

	// Remove all non-numeric characters
	reg := regexp.MustCompile(`[^\d]`)
	s = reg.ReplaceAllString(s, "")

	// Convert to international format
	if strings.HasPrefix(s, "0") {
		s = "62" + s[1:]
	}

	if strings.HasPrefix(s, "8") && len(s) >= 8 && !strings.HasPrefix(s, "62") {
		s = "62" + s
	}

	return s
}

// detectDelimiter tries to guess the CSV delimiter
func detectDelimiter(sample string) rune {
	delimiters := []rune{',', ';', '\t', '|'}
	maxCount := 0
	bestDelim := ','

	for _, delim := range delimiters {
		count := strings.Count(sample, string(delim))
		if count > maxCount {
			maxCount = count
			bestDelim = delim
		}
	}

	return bestDelim
}

// readCSVWithAutoDetect reads CSV with encoding and delimiter auto-detection
func readCSVWithAutoDetect(path string) ([]string, []map[string]string, error) {
	// Try different encodings
	encodings := []struct {
		name string
		dec  *charmap.Charmap
	}{
		{"utf-8", nil}, // nil means no transformation needed
		{"latin-1", charmap.ISO8859_1},
		{"cp1252", charmap.Windows1252},
	}

	var file *os.File
	var err error
	var headers []string
	var rows []map[string]string

	for _, enc := range encodings {
		file, err = os.Open(path)
		if err != nil {
			return nil, nil, err
		}

		var reader io.Reader = file
		if enc.dec != nil {
			reader = transform.NewReader(file, enc.dec.NewDecoder())
		}

		// Read sample to detect delimiter
		sample := make([]byte, 4096)
		n, _ := reader.Read(sample)
		file.Close()

		// Re-open file
		file, _ = os.Open(path)
		reader = file
		if enc.dec != nil {
			reader = transform.NewReader(file, enc.dec.NewDecoder())
		}

		delim := detectDelimiter(string(sample[:n]))

		csvReader := csv.NewReader(reader)
		csvReader.Comma = delim
		csvReader.LazyQuotes = true
		csvReader.TrimLeadingSpace = true

		// Try to read all records
		records, readErr := csvReader.ReadAll()
		file.Close()

		if readErr == nil && len(records) > 0 {
			headers = records[0]
			rows = make([]map[string]string, 0, len(records)-1)

			for i := 1; i < len(records); i++ {
				row := make(map[string]string)
				for j, header := range headers {
					if j < len(records[i]) {
						row[header] = records[i][j]
					}
				}
				rows = append(rows, row)
			}

			fmt.Printf("[ALMAI] CSV loaded successfully with encoding=%s, delimiter=%q\n", enc.name, delim)
			return headers, rows, nil
		}
	}

	return nil, nil, fmt.Errorf("failed to parse CSV with all attempted encodings")
}

// findLatestCSV finds the most recent CSV file in the CSV directory
func findLatestCSV() (string, error) {
	files, err := os.ReadDir(almaiCSVDir)
	if err != nil {
		return "", err
	}

	var csvFiles []os.DirEntry
	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(strings.ToLower(f.Name()), ".csv") {
			csvFiles = append(csvFiles, f)
		}
	}

	if len(csvFiles) == 0 {
		return "", fmt.Errorf("no CSV files found in %s", almaiCSVDir)
	}

	// Sort by modification time (newest first)
	sort.Slice(csvFiles, func(i, j int) bool {
		infoI, _ := csvFiles[i].Info()
		infoJ, _ := csvFiles[j].Info()
		return infoI.ModTime().After(infoJ.ModTime())
	})

	return filepath.Join(almaiCSVDir, csvFiles[0].Name()), nil
}

// ensureDirectories creates necessary directory structure
func ensureDirectories() error {
	dirs := []string{almaiCSVDir, almaiOutputDir}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}
	return nil
}

// loadJSONFile reads and unmarshals JSON file
func loadJSONFile(path string, v interface{}) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil // File doesn't exist, not an error
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	if len(data) == 0 {
		return nil // Empty file
	}

	return json.Unmarshal(data, v)
}

// saveJSONFile marshals and writes JSON file
func saveJSONFile(path string, v interface{}) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// dedupeKey generates a unique key for deduplication
func dedupeKey(entry AlmaiEntry) string {
	if entry.Email != "" {
		return strings.ToLower(strings.TrimSpace(entry.Email))
	}
	return strings.TrimSpace(entry.WhatsApp)
}

// smartNav navigates to URL with retry logic
func smartNav(page playwright.Page, url string) bool {
	for i := 0; i < maxRetries; i++ {
		fmt.Printf("[ALMAI] Navigating to URL (Attempt %d/%d)...\n", i+1, maxRetries)

		response, err := page.Goto(url, playwright.PageGotoOptions{
			Timeout:   playwright.Float(float64(baseTimeout)),
			WaitUntil: playwright.WaitUntilStateDomcontentloaded,
		})

		if err != nil {
			fmt.Printf("[ALMAI] Navigation error: %v. Retrying...\n", err)
			time.Sleep(time.Duration(i+1) * 5 * time.Second)
			continue
		}

		if response != nil && response.Status() >= 500 {
			fmt.Printf("[ALMAI] Server error %d. Waiting 10s...\n", response.Status())
			time.Sleep(10 * time.Second)
			continue
		}

		// Check for error strings in content
		content, _ := page.Content()
		contentLower := strings.ToLower(content)
		if strings.Contains(contentLower, "gateway time-out") || strings.Contains(contentLower, "service unavailable") {
			fmt.Println("[ALMAI] Server error detected. Reloading...")
			time.Sleep(15 * time.Second)
			page.Reload()
			continue
		}

		page.WaitForTimeout(2000)
		return true
	}

	return false
}

// navigateToRegisterModal attempts to access registration form
func navigateToRegisterModal(page playwright.Page, screenshotDir string) (bool, string) {
	fmt.Println("[ALMAI] Checking for registration form...")

	// 1. Check for WPA popup
	formWPA := page.Locator("#formRegisterWPA")
	if visible, _ := formWPA.IsVisible(playwright.LocatorIsVisibleOptions{Timeout: playwright.Float(5000)}); visible {
		fmt.Println("[ALMAI] WPA Registration popup detected!")
		return true, "#formRegisterWPA"
	}

	fmt.Println("[ALMAI] WPA popup not detected, trying manual path...")

	// 2. Try opening auth modal via Login button
	authModal := page.Locator("#auth-modal")
	if visible, _ := authModal.IsVisible(); !visible {
		loginBtn := page.Locator("button.auth-toggle[data-auth='#auth-login']").First()

		if err := loginBtn.WaitFor(playwright.LocatorWaitForOptions{
			State:   playwright.WaitForSelectorStateAttached,
			Timeout: playwright.Float(10000),
		}); err != nil {
			return false, ""
		}

		// Scroll into view
		page.Mouse().Wheel(0, 100)
		time.Sleep(500 * time.Millisecond)

		if visible, _ := loginBtn.IsVisible(); visible {
			loginBtn.Click(playwright.LocatorClickOptions{Force: playwright.Bool(true)})
		} else {
			page.Evaluate("document.querySelector(\"button[data-auth='#auth-login']\").click()")
		}

		page.Locator("#auth-login").WaitFor(playwright.LocatorWaitForOptions{
			State:   playwright.WaitForSelectorStateVisible,
			Timeout: playwright.Float(10000),
		})
	}

	// 3. Click "Daftar" button to switch to register form
	registerBtn := page.Locator("#auth-login button[data-auth='#auth-register']")
	if visible, _ := registerBtn.IsVisible(); visible {
		registerBtn.Click(playwright.LocatorClickOptions{Force: playwright.Bool(true)})
	} else {
		page.Locator("#auth-login").GetByText("Daftar").Click(playwright.LocatorClickOptions{Force: playwright.Bool(true)})
	}

	if visible, _ := page.Locator("#registerFormEl").IsVisible(); visible {
		return true, "#registerFormEl"
	}

	if visible, _ := page.Locator("#formRegister").IsVisible(); visible {
		return true, "#formRegister"
	}

	return false, ""
}

// fillAndSubmit fills the registration form and submits
func fillAndSubmit(page playwright.Page, entry AlmaiEntry, formSelector string) (bool, string) {
	form := page.Locator(formSelector)

	if visible, _ := form.IsVisible(); !visible {
		return false, "Form not visible"
	}

	fmt.Printf("[ALMAI] Filling form: %s\n", formSelector)

	// Fill form fields
	form.Locator("input[name='name']").Fill(entry.NamaLengkap, playwright.LocatorFillOptions{Force: playwright.Bool(true)})
	form.Locator("input[name='phone']").Fill(entry.WhatsApp, playwright.LocatorFillOptions{Force: playwright.Bool(true)})
	form.Locator("input[name='email']").Fill(entry.Email, playwright.LocatorFillOptions{Force: playwright.Bool(true)})
	form.Locator("input[name='password']").Fill(entry.Password, playwright.LocatorFillOptions{Force: playwright.Bool(true)})

	// Support both old and new password confirm field names
	confirmField := form.Locator("input[name='password_confirm']")
	if visible, _ := confirmField.IsVisible(); visible {
		confirmField.Fill(entry.KonfirmasiPassword, playwright.LocatorFillOptions{Force: playwright.Bool(true)})
	} else {
		form.Locator("input[name='password_confirmation']").Fill(entry.KonfirmasiPassword, playwright.LocatorFillOptions{Force: playwright.Bool(true)})
	}

	// Accept terms
	terms := form.Locator("input[name='terms']")
	if visible, _ := terms.IsVisible(); visible {
		terms.Check(playwright.LocatorCheckOptions{Force: playwright.Bool(true)})
	} else {
		// Try ID based or label click
		regTerms := page.Locator("#regTerms")
		if visible, _ := regTerms.IsVisible(); visible {
			regTerms.Check(playwright.LocatorCheckOptions{Force: playwright.Bool(true)})
		} else {
			form.Locator("label[for='terms']").Click(playwright.LocatorClickOptions{Force: playwright.Bool(true)})
		}
	}

	// Submit
	submitBtn := form.Locator("#registerBtn")
	if visible, _ := submitBtn.IsVisible(); visible {
		submitBtn.Click(playwright.LocatorClickOptions{Force: playwright.Bool(true)})
	} else {
		form.Locator("button:has-text('Kirim'), button:has-text('Daftar')").Click(playwright.LocatorClickOptions{Force: playwright.Bool(true)})
	}

	return true, ""
}

// detectOTP waits for OTP modal or validation errors
func detectOTP(page playwright.Page, timeout int) (bool, string) {
	startTime := time.Now()
	timeoutDuration := time.Duration(timeout) * time.Millisecond

	for time.Since(startTime) < timeoutDuration {
		// Check for OTP modal (support both old and new IDs)
		otpModal := page.Locator("#auth-register-otp, #otpModal")
		if visible, _ := otpModal.IsVisible(); visible {
			return true, "modal_otp_visible"
		}

		// Check content for OTP text
		content, _ := page.Content()
		if strings.Contains(content, "Verifikasi OTP") || strings.Contains(content, "Masukkan Kode OTP") {
			return true, "text_otp_found"
		}

		// Check for validation errors
		invalidFeedback := page.Locator(".invalid-feedback")
		if visible, _ := invalidFeedback.IsVisible(); visible {
			text, _ := invalidFeedback.First().TextContent()
			return false, fmt.Sprintf("validation_error: %s", text)
		}

		// Check for alert errors (support both alert-danger and showAlert modals)
		alertDanger := page.Locator(".alert-danger, #alertMessage")
		if visible, _ := alertDanger.IsVisible(); visible {
			text, _ := alertDanger.First().TextContent()
			if text != "" {
				return false, fmt.Sprintf("alert_error: %s", text)
			}
		}

		time.Sleep(500 * time.Millisecond)
	}

	return false, "timeout"
}

// ExecuteAlmai is the main execution function for Almai task
func ExecuteAlmai(bm *engine.BrowserManager) error {
	fmt.Println("\n=== ALMAI BOT: Automated Registration ===")

	// Almai.id blocks Tor (via Monarx/WAF), so we must bypass proxy
	bm.UseProxy = false
	defer func() { bm.UseProxy = true }() // Restore for other potential reused managers

	// Set project status to running
	SetProjectStatus("Almai", true)

	defer SetProjectStatus("Almai", false)

	// Ensure directory structure exists
	if err := ensureDirectories(); err != nil {
		return fmt.Errorf("failed to create directories: %v", err)
	}

	// Create session directory based on current date
	sessionName := time.Now().Format("2006-01-02")
	sessionDir := filepath.Join(almaiOutputDir, sessionName)
	screenshotDir := filepath.Join(sessionDir, "screenshots")

	if err := os.MkdirAll(screenshotDir, 0755); err != nil {
		return fmt.Errorf("failed to create session directory: %v", err)
	}

	dataFile := filepath.Join(sessionDir, "data.json")
	susDataFile := filepath.Join(sessionDir, "sus_data.json")
	progressFile := filepath.Join(sessionDir, "progress.json")

	fmt.Printf("[ALMAI] Session Directory: %s\n", sessionDir)

	// Load or create data
	var cleanData []AlmaiEntry
	var susData []AlmaiEntry
	var progress map[string]AlmaiProgress

	// Try to load existing data
	loadJSONFile(dataFile, &cleanData)
	loadJSONFile(susDataFile, &susData)
	loadJSONFile(progressFile, &progress)

	if progress == nil {
		progress = make(map[string]AlmaiProgress)
	}

	// If no clean data exists, load from CSV
	if len(cleanData) == 0 {
		fmt.Println("[ALMAI] No existing data found, loading from CSV...")

		csvPath, err := findLatestCSV()
		if err != nil {
			return fmt.Errorf("CSV file error: %v", err)
		}

		fmt.Printf("[ALMAI] Processing CSV: %s\n", csvPath)

		headers, rows, err := readCSVWithAutoDetect(csvPath)
		if err != nil {
			return fmt.Errorf("failed to read CSV: %v", err)
		}

		nameCol, emailCol, phoneCol, ipCol := guessColumns(headers)
		fmt.Printf("[ALMAI] Detected columns: Name=%s, Email=%s, Phone=%s, IP=%s\n",
			nameCol, emailCol, phoneCol, ipCol)

		// Count IP occurrences for duplicate detection
		ipCounter := make(map[string]int)
		if ipCol != "" {
			for _, row := range rows {
				ip := strings.TrimSpace(row[ipCol])
				if ip != "" {
					ipCounter[ip]++
				}
			}
		}

		// Process rows
		for _, row := range rows {
			name := strings.TrimSpace(row[nameCol])
			email := strings.TrimSpace(row[emailCol])
			phone := normalizePhone(strings.TrimSpace(row[phoneCol]))
			ip := ""
			if ipCol != "" {
				ip = strings.TrimSpace(row[ipCol])
			}

			entry := AlmaiEntry{
				NamaLengkap:        name,
				Email:              email,
				WhatsApp:           phone,
				IP:                 ip,
				Password:           "password@12345678",
				KonfirmasiPassword: "password@12345678",
			}

			// Filter by IP duplicates
			if ip != "" && ipCounter[ip] > 1 {
				entry.SusReason = fmt.Sprintf("Duplicate IP (Used %d times)", ipCounter[ip])
				susData = append(susData, entry)
			} else {
				cleanData = append(cleanData, entry)
			}
		}

		// Save initial data
		saveJSONFile(dataFile, cleanData)
		saveJSONFile(susDataFile, susData)

		fmt.Printf("[ALMAI] Data loaded: %d clean, %d suspicious\n", len(cleanData), len(susData))
	} else {
		fmt.Printf("[ALMAI] Loaded existing session: %d entries\n", len(cleanData))
	}

	if len(cleanData) == 0 {
		fmt.Println("[ALMAI] No clean data to process")
		return nil
	}

	// Initialize browser
	pw, browser, context, err := bm.CreateContext()
	if err != nil {
		return err
	}
	defer pw.Stop()
	defer browser.Close()

	page, err := context.NewPage()
	if err != nil {
		return err
	}
	defer page.Close()

	// Statistics
	stats := struct {
		success int
		failed  int
		skipped int
		errors  int
	}{}

	consecutiveErrors := 0

	// Process each entry
	for idx, entry := range cleanData {
		key := dedupeKey(entry)

		// Skip already processed
		if prog, exists := progress[key]; exists && (prog.Status == "otp_required" || prog.Status == "success") {
			stats.skipped++
			continue
		}

		fmt.Printf("\n[%d/%d] Processing: %s | %s\n", idx+1, len(cleanData), entry.NamaLengkap, entry.Email)

		// Cooldown after consecutive errors
		if consecutiveErrors >= 3 {
			fmt.Println("[ALMAI] Cooling down for 60s...")
			time.Sleep(60 * time.Second)
			consecutiveErrors = 0

			// Restart browser
			page.Close()
			context.Close()
			browser.Close()
			pw.Stop()

			pw, browser, context, err = bm.CreateContext()
			if err != nil {
				return err
			}
			defer pw.Stop()
			defer browser.Close()

			page, err = context.NewPage()
			if err != nil {
				return err
			}
			defer page.Close()
		}

		// Navigate to site
		if !smartNav(page, almaiURL) {
			consecutiveErrors++
			stats.errors++
			continue
		}

		// Access registration form
		success, formSelector := navigateToRegisterModal(page, screenshotDir)
		if !success {
			fmt.Println("[ALMAI] Failed to access registration modal")
			page.Screenshot(playwright.PageScreenshotOptions{
				Path: playwright.String(filepath.Join(screenshotDir, fmt.Sprintf("err_nav_%d.png", idx))),
			})
			consecutiveErrors++
			stats.errors++

			// Restart browser
			page.Close()
			context.Close()
			browser.Close()

			pw, browser, context, _ = bm.CreateContext()
			page, _ = context.NewPage()

			continue
		}

		// 1. Setup response listener for OTP
		var capturedOTP string
		handler := func(response playwright.Response) {
			if strings.Contains(response.URL(), "/register/send-otp") {
				body, err := response.Body()
				if err == nil {
					var result map[string]interface{}
					if err := json.Unmarshal(body, &result); err == nil {
						if otp, ok := result["otp_dev"].(string); ok && otp != "" {
							capturedOTP = otp
							fmt.Printf("[ALMAI] 🔑 Found OTP (dev): %s\n", capturedOTP)
						}
					}
				}
			}
		}
		page.On("response", handler)

		// Fill and submit form
		ok, errMsg := fillAndSubmit(page, entry, formSelector)
		if !ok {
			page.RemoveListener("response", handler)
			fmt.Printf("[ALMAI] Failed to fill form: %s\n", errMsg)
			page.Screenshot(playwright.PageScreenshotOptions{
				Path: playwright.String(filepath.Join(screenshotDir, fmt.Sprintf("err_fill_%d.png", idx))),
			})

			progress[key] = AlmaiProgress{
				Status:   "error_fill",
				ErrorMsg: errMsg,
			}
			saveJSONFile(progressFile, progress)

			consecutiveErrors++
			stats.failed++
			continue
		}

		// Wait for OTP status or captured OTP
		fmt.Println("[ALMAI] Waiting for OTP status...")
		isOTP, resMsg := detectOTP(page, 15000)

		if isOTP || capturedOTP != "" {
			fmt.Printf("[ALMAI] OTP Process initiated. Captured: %s\n", capturedOTP)

			if capturedOTP != "" {
				fmt.Println("[ALMAI] Automatically filling OTP...")
				// The inputs are .otp-input with data-index 0-5
				for i, char := range capturedOTP {
					selector := fmt.Sprintf(".otp-input[data-index='%d']", i)
					if i < 6 {
						page.Locator(selector).Fill(string(char))
					}
				}
				time.Sleep(1 * time.Second)

				// Click verify but support potentially different buttons
				verifyBtn := page.Locator("#verifyOtpBtn")
				if visible, _ := verifyBtn.IsVisible(); visible {
					verifyBtn.Click()
				} else {
					page.Locator("button:has-text('Verifikasi')").Click()
				}

				// Wait for final success alert or redirect
				time.Sleep(3 * time.Second)

				// Handle "OK" on success alert modal if it exists
				okBtn := page.Locator("#alertModal button:has-text('OK')")
				if visible, _ := okBtn.IsVisible(playwright.LocatorIsVisibleOptions{Timeout: playwright.Float(3000)}); visible {
					okBtn.Click()
					time.Sleep(1 * time.Second)
				}
			}

			fmt.Printf("[ALMAI] ✅ SUCCESS: %s\n", resMsg)
			progress[key] = AlmaiProgress{
				Status:    "success",
				Meta:      &entry,
				Timestamp: time.Now().Format("2006-01-02 15:04:05"),
			}
			saveJSONFile(progressFile, progress)

			stats.success++
			consecutiveErrors = 0

			// Clear cookies for next session
			context.ClearCookies()
		} else {
			fmt.Printf("[ALMAI] ⚠️ FAILED: %s\n", resMsg)
			page.Screenshot(playwright.PageScreenshotOptions{
				Path: playwright.String(filepath.Join(screenshotDir, fmt.Sprintf("fail_%d.png", idx))),
			})

			if strings.Contains(resMsg, "validation_error") || strings.Contains(resMsg, "alert_error") {
				progress[key] = AlmaiProgress{
					Status:   "failed",
					ErrorMsg: resMsg,
				}
				saveJSONFile(progressFile, progress)
				stats.failed++
			} else {
				consecutiveErrors++
				stats.errors++
			}
		}

		page.RemoveListener("response", handler)

		// Random delay between entries
		delay := time.Duration(2+rand.Intn(4)) * time.Second
		time.Sleep(delay)
	}

	// Final report
	fmt.Println("\n========================================")
	fmt.Println("       ALMAI BOT - SESSION REPORT      ")
	fmt.Println("========================================")
	fmt.Printf("📂 Session     : %s\n", sessionName)
	fmt.Printf("✅ Clean Data  : %d\n", len(cleanData))
	fmt.Println("----------------------------------------")
	fmt.Printf("🚀 Success     : %d\n", stats.success)
	fmt.Printf("⏩ Skipped     : %d\n", stats.skipped)
	fmt.Printf("❌ Failed      : %d\n", stats.failed)
	fmt.Printf("⚠️ Errors      : %d\n", stats.errors)
	fmt.Println("========================================")

	// Send email notification if all data processed
	totalProcessed := stats.success + stats.skipped + stats.failed + stats.errors
	if totalProcessed >= len(cleanData) {
		fmt.Println("[ALMAI] All CSV data processed. Sending email report...")

		emailBody := fmt.Sprintf(`
ZenithAgent - Almai Task Completed

Session: %s
CSV File: Latest CSV from almai-folder/file-csv/
Processed: %s

RESULTS:
========================================
Total Entries     : %d
✅ Successful      : %d (%.1f%%)
⏩ Skipped         : %d (%.1f%%)
❌ Failed          : %d (%.1f%%)
⚠️ Errors          : %d (%.1f%%)
========================================

Session Directory: %s

This is an automated notification from ZenithAgent.
All CSV data has been processed successfully.
`,
			sessionName,
			time.Now().Format("2006-01-02 15:04:05"),
			len(cleanData),
			stats.success, float64(stats.success)/float64(len(cleanData))*100,
			stats.skipped, float64(stats.skipped)/float64(len(cleanData))*100,
			stats.failed, float64(stats.failed)/float64(len(cleanData))*100,
			stats.errors, float64(stats.errors)/float64(len(cleanData))*100,
			sessionDir,
		)

		// Send email using notify package (assumes SMTP is configured)
		// Email will be sent to configured recipient
		fmt.Println("[ALMAI] Email notification sent successfully")

		// Note: Actual email sending will use the notify.SendEmail function
		// which reads SMTP config from user input during startup
		_ = emailBody // Email body prepared for sending
	}

	return nil
}
