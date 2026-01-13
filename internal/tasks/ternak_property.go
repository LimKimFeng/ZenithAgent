package tasks

import (
	"fmt"
	"math/rand"
	"strings" // Added missing import
	"time"
	"zenith-agent/internal/engine"

	"github.com/playwright-community/playwright-go"
)

// Smart Reason Generator
func generateSmartReason() string {
	prefix := []string{"Saya ingin ", "Mau sekali ", "Rencana saya ", "Tujuan saya "}
	action := []string{"belajar ", "mendalami ", "memahami cara ", "mencari tahu "}
	topic := []string{"passive income properti ", "strategi pensiun aman ", "jalur lelang perbankan ", "investasi aset masa depan "}
	suffix := []string{"agar keuangan stabil.", "untuk keluarga saya.", "biar ada penghasilan tambahan.", "secara legal dan benar."}

	rand.Seed(time.Now().UnixNano())
	return prefix[rand.Intn(len(prefix))] + action[rand.Intn(len(action))] + topic[rand.Intn(len(topic))] + suffix[rand.Intn(len(suffix))]
}

func ExecuteTernakProperty(bm *engine.BrowserManager) error {
    // Set status running di dashboard
    SetProjectStatus("TernakProperty", true)

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

	// Block resource berat untuk hemat VPS
	page.Route("**/*.{png,jpg,jpeg,gif,webp,css}", func(route playwright.Route) {
		route.Abort()
	})

	// 1. Navigate - Using WaitUntilStateLoad as fixed previously
	_, err = page.Goto("https://ternakproperty.com/as-bandung?utm_source=fb&utm_medium=mof&utm_campaign=mof_seminarbandung&utm_content=vidbeforeafterv1", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateLoad,
	})
	if err != nil {
		return err
	}

	// Wait for form
	_, err = page.WaitForSelector("#name", playwright.PageWaitForSelectorOptions{
		Timeout: playwright.Float(30000),
	})
	if err != nil {
		return fmt.Errorf("timeout waiting for form: %v", err)
	}

	// Data Dummy
	names := []string{"Andi Pratama", "Budi Santoso", "Cahyo Nugroho", "Deni Saputra", "Eko Wibowo", "Fajar Ramadhan", "Gilang Saputra"}
	selectedName := names[rand.Intn(len(names))]
	phone := fmt.Sprintf("628%d", rand.Intn(900000000)+100000000)
	email := fmt.Sprintf("%s%d@example.com", strings.ReplaceAll(strings.ToLower(selectedName), " ", ""), time.Now().Unix()%1000)

	// 2. Fill Form with Human Delay
	typeOptions := playwright.PageTypeOptions{Delay: playwright.Float(100)}
	
	page.Type("#name", selectedName, typeOptions)
	page.Type("#phone", phone, typeOptions)
	page.Type("#email", email, typeOptions)
	page.Type("#notes", generateSmartReason(), typeOptions)

	// 3. Select BCA
	err = page.Click("text=Bank Central Asia")
	if err != nil {
		return fmt.Errorf("failed to click BCA: %v", err)
	}

	// 4. Submit
	err = page.Click("button:has-text('Daftar Sekarang')")
	if err != nil {
		return err
	}

	// 5. Success Detection
	err = page.WaitForURL("**/success", playwright.PageWaitForURLOptions{
		Timeout: playwright.Float(30000),
	})
	
	// Set status ke false setelah selesai
	SetProjectStatus("TernakProperty", false)
	
	return err
}

// Local UpdateStats removed; using common.go version
