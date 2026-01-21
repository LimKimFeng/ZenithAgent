package tasks

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
	"zenith-agent/internal/engine"

	"github.com/playwright-community/playwright-go"
)

const (
	tokoprodukURL           = "https://tokoproduk.com/kabeldata/"
	tokoprodukWhatsAppPhone = "6289676006621"
	tokoprodukTimeout       = 60000
)

// generateTokoprodukData generates dummy order data
func generateTokoprodukData() (name, phone, address string) {
	firstNames := []string{
		"Andi", "Budi", "Cahyo", "Deni", "Eko", "Fajar", "Gilang",
		"Hadi", "Indra", "Joko", "Kevin", "Leo", "Made", "Nico",
		"Omar", "Putra", "Reza", "Sandi", "Toni", "Udin",
	}

	lastNames := []string{
		"Pratama", "Santoso", "Nugroho", "Saputra", "Wibowo",
		"Ramadhan", "Kusuma", "Wijaya", "Setiawan", "Gunawan",
		"Firmansyah", "Hermawan", "Kurniawan", "Maulana", "Hidayat",
	}

	streets := []string{
		"Jl. Merdeka", "Jl. Sudirman", "Jl. Thamrin", "Jl. Gatot Subroto",
		"Jl. Ahmad Yani", "Jl. Diponegoro", "Jl. Imam Bonjol",
		"Jl. Hayam Wuruk", "Jl. Gajah Mada", "Jl. Pemuda",
	}

	cities := []string{
		"Jakarta", "Bandung", "Surabaya", "Semarang", "Yogyakarta",
		"Malang", "Solo", "Medan", "Palembang", "Makassar",
	}

	provinces := []string{
		"DKI Jakarta", "Jawa Barat", "Jawa Timur", "Jawa Tengah",
		"DI Yogyakarta", "Sumatera Utara", "Sumatera Selatan", "Sulawesi Selatan",
	}

	rand.Seed(time.Now().UnixNano())

	// Generate name
	firstName := firstNames[rand.Intn(len(firstNames))]
	lastName := lastNames[rand.Intn(len(lastNames))]
	name = fmt.Sprintf("%s %s", firstName, lastName)

	// Generate phone (Indonesian format)
	phoneNumber := rand.Intn(900000000) + 100000000
	phone = fmt.Sprintf("628%d", phoneNumber)

	// Generate full address
	street := streets[rand.Intn(len(streets))]
	houseNum := rand.Intn(200) + 1
	rtNum := rand.Intn(20) + 1
	rwNum := rand.Intn(15) + 1
	city := cities[rand.Intn(len(cities))]
	province := provinces[rand.Intn(len(provinces))]

	address = fmt.Sprintf("%s No. %d RT %02d/RW %02d, Kel. %s, Kec. %s, %s, %s",
		street, houseNum, rtNum, rwNum,
		firstName+"an", lastName+"an", city, province)

	return name, phone, address
}

// fillTokoprodukForm fills the order form
func fillTokoprodukForm(page playwright.Page, name, phone, address string) error {
	fmt.Println("[TOKOPRODUK] Filling form...")

	// Wait for form to be ready
	time.Sleep(2 * time.Second)

	// Fill name - using class selector
	nameInput := page.Locator(".ooef-field-name")
	if err := nameInput.Fill(name); err != nil {
		return fmt.Errorf("failed to fill name: %v", err)
	}
	time.Sleep(500 * time.Millisecond)

	// Fill phone
	phoneInput := page.Locator(".ooef-field-phone")
	if err := phoneInput.Fill(phone); err != nil {
		return fmt.Errorf("failed to fill phone: %v", err)
	}
	time.Sleep(500 * time.Millisecond)

	// Fill address
	addressInput := page.Locator(".ooef-field-address")
	if err := addressInput.Fill(address); err != nil {
		return fmt.Errorf("failed to fill address: %v", err)
	}
	time.Sleep(500 * time.Millisecond)

	return nil
}

// submitTokoprodukForm submits the order
func submitTokoprodukForm(page playwright.Page) error {
	fmt.Println("[TOKOPRODUK] Submitting order...")

	// Find and click submit button with class "ooef-submit-order"
	submitButton := page.Locator("button.ooef-submit-order")

	if err := submitButton.Click(); err != nil {
		return fmt.Errorf("failed to click submit button: %v", err)
	}

	return nil
}

// waitForTokoprodukWhatsAppRedirect waits for WhatsApp redirect
func waitForTokoprodukWhatsAppRedirect(page playwright.Page) error {
	fmt.Println("[TOKOPRODUK] Waiting for WhatsApp redirect...")

	startTime := time.Now()
	timeout := 30 * time.Second

	for time.Since(startTime) < timeout {
		currentURL := page.URL()

		// Check if redirected to WhatsApp with the specific phone number
		if strings.Contains(currentURL, "api.whatsapp.com/send") &&
			strings.Contains(currentURL, tokoprodukWhatsAppPhone) {
			fmt.Println("[TOKOPRODUK] ✅ Successfully redirected to WhatsApp!")
			fmt.Printf("[TOKOPRODUK] WhatsApp URL: %s\n", currentURL)
			return nil
		}

		time.Sleep(500 * time.Millisecond)
	}

	return fmt.Errorf("timeout waiting for WhatsApp redirect")
}

// ExecuteTokoprodukKabeldata is the main execution function
func ExecuteTokoprodukKabeldata(bm *engine.BrowserManager) error {
	fmt.Println("\n=== TOKOPRODUK KABEL DATA BOT: Starting Order ===")

	// Set project status to running
	SetProjectStatus("TokoprodukKabeldata", true)
	defer SetProjectStatus("TokoprodukKabeldata", false)

	// Generate dummy data
	name, phone, address := generateTokoprodukData()

	fmt.Printf("[TOKOPRODUK] Generated Data:\n")
	fmt.Printf("  Name: %s\n", name)
	fmt.Printf("  Phone: %s\n", phone)
	fmt.Printf("  Address: %s\n", address)

	// Initialize browser
	pw, browser, context, err := bm.CreateContext()
	if err != nil {
		return fmt.Errorf("failed to create browser context: %v", err)
	}
	defer pw.Stop()
	defer browser.Close()

	page, err := context.NewPage()
	if err != nil {
		return fmt.Errorf("failed to create new page: %v", err)
	}
	defer page.Close()

	// Navigate to site
	fmt.Printf("[TOKOPRODUK] Navigating to %s\n", tokoprodukURL)
	_, err = page.Goto(tokoprodukURL, playwright.PageGotoOptions{
		Timeout:   playwright.Float(float64(tokoprodukTimeout)),
		WaitUntil: playwright.WaitUntilStateNetworkidle,
	})
	if err != nil {
		return fmt.Errorf("failed to navigate: %v", err)
	}

	// Fill form
	if err := fillTokoprodukForm(page, name, phone, address); err != nil {
		return err
	}

	// Submit form
	if err := submitTokoprodukForm(page); err != nil {
		return err
	}

	// Wait for WhatsApp redirect
	if err := waitForTokoprodukWhatsAppRedirect(page); err != nil {
		return err
	}

	fmt.Println("[TOKOPRODUK] ✅ Order completed successfully!")
	return nil
}
