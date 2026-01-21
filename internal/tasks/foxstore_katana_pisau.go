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
	foxstoreURL        = "https://foxstore.buat.io/2in1-katana-pisau-serbaguna"
	foxstoreWhatsAppID = "6285129390658"
	foxstoreTimeout    = 60000
)

// City data for autocomplete selection
var foxstoreCities = []string{
	"Balikpapan Barat, Kota Balikpapan, Kalimantan Timur",
	"Denpasar Barat, Kota Denpasar, Bali",
	"Batang Masumai, Kab. Merangin, Jambi",
	"Tamako, Kab. Kepulauan Sangihe, Sulawesi Utara",
	"Hiliserangkai, Kab. Nias, Sumatera Utara",
}

// generateFoxstoreData generates dummy order data
func generateFoxstoreData() (name, phone, address, city, product string) {
	firstNames := []string{
		"Andi", "Budi", "Cahyo", "Deni", "Eko", "Fajar", "Gilang",
		"Hadi", "Indra", "Joko", "Kevin", "Leo", "Made", "Nico",
	}

	lastNames := []string{
		"Pratama", "Santoso", "Nugroho", "Saputra", "Wibowo",
		"Ramadhan", "Kusuma", "Wijaya", "Setiawan", "Gunawan",
	}

	streets := []string{
		"Jl. Merdeka", "Jl. Sudirman", "Jl. Thamrin", "Jl. Gatot Subroto",
		"Jl. Ahmad Yani", "Jl. Diponegoro", "Jl. Imam Bonjol",
	}

	rand.Seed(time.Now().UnixNano())

	// Generate name
	firstName := firstNames[rand.Intn(len(firstNames))]
	lastName := lastNames[rand.Intn(len(lastNames))]
	name = fmt.Sprintf("%s %s", firstName, lastName)

	// Generate phone (Indonesian format)
	phoneNumber := rand.Intn(900000000) + 100000000
	phone = fmt.Sprintf("628%d", phoneNumber)

	// Generate address
	street := streets[rand.Intn(len(streets))]
	houseNumber := rand.Intn(200) + 1
	address = fmt.Sprintf("%s No. %d", street, houseNumber)

	// Select random city
	city = foxstoreCities[rand.Intn(len(foxstoreCities))]

	// Random product selection
	products := []string{
		"Beli 1 - Rp. 69.000",
		"Beli 2 - Rp. 99.000",
	}
	product = products[rand.Intn(len(products))]

	return name, phone, address, city, product
}

// fillFoxstoreProductChoice selects product option
func fillFoxstoreProductChoice(page playwright.Page, product string) error {
	fmt.Printf("[FOXSTORE] Selecting product: %s\n", product)

	// Wait for radio buttons to be available
	time.Sleep(1 * time.Second)

	// Find the radio button with matching value
	radioSelector := fmt.Sprintf("input[type='radio'][name='radio-group-xzq2wa'][value='%s']", product)

	err := page.Click(radioSelector, playwright.PageClickOptions{
		Timeout: playwright.Float(10000),
		Force:   playwright.Bool(true),
	})

	if err != nil {
		return fmt.Errorf("failed to select product: %v", err)
	}

	time.Sleep(500 * time.Millisecond)
	return nil
}

// fillFoxstoreCityAutocomplete handles the autocomplete city selection
func fillFoxstoreCityAutocomplete(page playwright.Page, cityFull string) error {
	fmt.Printf("[FOXSTORE] Selecting city: %s\n", cityFull)

	// Get first part of city name for searching
	cityParts := strings.Split(cityFull, ",")
	searchText := strings.TrimSpace(cityParts[0])

	// Focus on the district input
	districtInput := page.Locator("input#district")
	if err := districtInput.Click(); err != nil {
		return fmt.Errorf("failed to click district input: %v", err)
	}

	time.Sleep(500 * time.Millisecond)

	// Type character by character to trigger autocomplete
	fmt.Printf("[FOXSTORE] Typing city name: %s\n", searchText)
	for _, char := range searchText {
		if err := districtInput.Press(string(char)); err != nil {
			return fmt.Errorf("failed to type character %c: %v", char, err)
		}
		time.Sleep(100 * time.Millisecond) // Small delay between characters
	}

	// Wait for autocomplete dropdown to appear
	time.Sleep(1 * time.Second)

	// Try to find and click the matching li element
	liSelector := fmt.Sprintf("li:has-text('%s')", cityFull)

	fmt.Println("[FOXSTORE] Waiting for autocomplete option to appear...")

	// Wait for the li element to be visible
	_, err := page.WaitForSelector(liSelector, playwright.PageWaitForSelectorOptions{
		State:   playwright.WaitForSelectorStateVisible,
		Timeout: playwright.Float(5000),
	})

	if err != nil {
		fmt.Printf("[FOXSTORE] Warning: Autocomplete option not found, trying alternative: %v\n", err)
		// Try clicking any li with the city name
		altSelector := "li.cursor-pointer"
		if clickErr := page.Click(altSelector, playwright.PageClickOptions{
			Timeout: playwright.Float(3000),
			Force:   playwright.Bool(true),
		}); clickErr != nil {
			return fmt.Errorf("failed to select from autocomplete: %v", clickErr)
		}
	} else {
		// Click the matched option
		if err := page.Click(liSelector, playwright.PageClickOptions{
			Timeout: playwright.Float(3000),
			Force:   playwright.Bool(true),
		}); err != nil {
			return fmt.Errorf("failed to click autocomplete option: %v", err)
		}
	}

	time.Sleep(500 * time.Millisecond)
	fmt.Println("[FOXSTORE] City selected successfully")
	return nil
}

// fillFoxstoreForm fills all form fields
func fillFoxstoreForm(page playwright.Page, name, phone, address, city, product string) error {
	fmt.Println("[FOXSTORE] Filling form...")

	// Wait for form to be fully loaded
	time.Sleep(2 * time.Second)

	// Select product
	if err := fillFoxstoreProductChoice(page, product); err != nil {
		return err
	}

	// Fill name
	nameInput := page.Locator("input#name")
	if err := nameInput.Fill(name); err != nil {
		return fmt.Errorf("failed to fill name: %v", err)
	}
	time.Sleep(300 * time.Millisecond)

	// Fill phone
	phoneInput := page.Locator("input#phone")
	if err := phoneInput.Fill(phone); err != nil {
		return fmt.Errorf("failed to fill phone: %v", err)
	}
	time.Sleep(300 * time.Millisecond)

	// Fill address
	addressTextarea := page.Locator("textarea#address")
	if err := addressTextarea.Fill(address); err != nil {
		return fmt.Errorf("failed to fill address: %v", err)
	}
	time.Sleep(300 * time.Millisecond)

	// Fill city with autocomplete
	if err := fillFoxstoreCityAutocomplete(page, city); err != nil {
		return err
	}

	// Select COD payment method
	codRadio := page.Locator("input[type='radio'][value='cod']")
	if err := codRadio.Click(playwright.LocatorClickOptions{
		Force: playwright.Bool(true),
	}); err != nil {
		return fmt.Errorf("failed to select COD: %v", err)
	}
	time.Sleep(300 * time.Millisecond)

	return nil
}

// submitFoxstoreForm submits the order form
func submitFoxstoreForm(page playwright.Page) error {
	fmt.Println("[FOXSTORE] Submitting form...")

	// Find and click submit button
	submitButton := page.Locator("button[type='submit']")

	if err := submitButton.Click(); err != nil {
		return fmt.Errorf("failed to click submit button: %v", err)
	}

	return nil
}

// waitForFoxstoreWhatsAppRedirect waits for redirect to WhatsApp
func waitForFoxstoreWhatsAppRedirect(page playwright.Page) error {
	fmt.Println("[FOXSTORE] Waiting for WhatsApp redirect...")

	startTime := time.Now()
	timeout := 30 * time.Second

	for time.Since(startTime) < timeout {
		currentURL := page.URL()

		// Check if redirected to WhatsApp
		if strings.Contains(currentURL, "api.whatsapp.com/send") &&
			strings.Contains(currentURL, foxstoreWhatsAppID) {
			fmt.Println("[FOXSTORE] ✅ Successfully redirected to WhatsApp!")
			return nil
		}

		time.Sleep(500 * time.Millisecond)
	}

	return fmt.Errorf("timeout waiting for WhatsApp redirect")
}

// ExecuteFoxstoreKatanaPisau is the main execution function
func ExecuteFoxstoreKatanaPisau(bm *engine.BrowserManager) error {
	fmt.Println("\n=== FOXSTORE KATANA PISAU BOT: Starting Order ===")

	// Set project status to running
	SetProjectStatus("FoxstoreKatanaPisau", true)
	defer SetProjectStatus("FoxstoreKatanaPisau", false)

	// Generate dummy data
	name, phone, address, city, product := generateFoxstoreData()

	fmt.Printf("[FOXSTORE] Generated Data:\n")
	fmt.Printf("  Name: %s\n", name)
	fmt.Printf("  Phone: %s\n", phone)
	fmt.Printf("  Address: %s\n", address)
	fmt.Printf("  City: %s\n", city)
	fmt.Printf("  Product: %s\n", product)

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
	fmt.Printf("[FOXSTORE] Navigating to %s\n", foxstoreURL)
	_, err = page.Goto(foxstoreURL, playwright.PageGotoOptions{
		Timeout:   playwright.Float(float64(foxstoreTimeout)),
		WaitUntil: playwright.WaitUntilStateNetworkidle,
	})
	if err != nil {
		return fmt.Errorf("failed to navigate: %v", err)
	}

	// Fill form
	if err := fillFoxstoreForm(page, name, phone, address, city, product); err != nil {
		return err
	}

	// Submit form
	if err := submitFoxstoreForm(page); err != nil {
		return err
	}

	// Wait for WhatsApp redirect
	if err := waitForFoxstoreWhatsAppRedirect(page); err != nil {
		return err
	}

	fmt.Println("[FOXSTORE] ✅ Order completed successfully!")
	return nil
}
