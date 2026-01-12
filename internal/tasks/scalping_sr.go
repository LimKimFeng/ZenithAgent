package tasks

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
	"zenith-agent/internal/engine"

	"github.com/playwright-community/playwright-go"
)

// Smart Reason Generator Spesifik Scalping
func generateScalpingReason() string {
	prefix := []string{"Saya ingin ", "Tujuan saya ", "Mau ", "Rencana "}
	action := []string{"belajar teknik S&R ", "paham cara scalping ", "mendalami beli murah jual dekat ", "cari profit konsisten "}
	topic := []string{"di market kripto ", "buat harian ", "biar nggak boncos terus ", "sebagai side income "}
	suffix := []string{"dari rumah.", "tanpa ribet.", "biar lebih tenang.", "secara logis."}

	rand.Seed(time.Now().UnixNano())
	return prefix[rand.Intn(len(prefix))] + action[rand.Intn(len(action))] + topic[rand.Intn(len(topic))] + suffix[rand.Intn(len(suffix))]
}

func ExecuteScalpingSR(bm *engine.BrowserManager) error {
	pw, browser, context, err := bm.CreateContext()
	if err != nil { return err }
	defer pw.Stop()
	defer browser.Close()

	page, err := context.NewPage()
	if err != nil { return err }
	defer page.Close()

	// Hemat Resource
	page.Route("**/*.{png,jpg,jpeg,gif,webp,css}", func(route playwright.Route) { route.Abort() })

	// 1. Navigate
	targetURL := "https://arknamedia.com/scalping-support-resistance?utm_medium=paid&utm_source=ig&utm_id=120240391765430028&utm_content=120240391765440028&utm_term=120240391765450028&utm_campaign=120240391765430028"
	_, err = page.Goto(targetURL, playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateLoad,
		Timeout:   playwright.Float(60000),
	})
	if err != nil { return fmt.Errorf("nav_failed: %v", err) }

	// Tunggu Form
	_, err = page.WaitForSelector("#name", playwright.PageWaitForSelectorOptions{Timeout: playwright.Float(15000)})
	if err != nil { return fmt.Errorf("form_not_found") }

	// 2. Dummy Data
	names := []string{"Andi Pratama", "Budi Santoso", "Cahyo Nugroho", "Deni Saputra", "Eko Wibowo", "Fajar Ramadhan", "Gilang Saputra"}
	selectedName := names[rand.Intn(len(names))]
	phone := fmt.Sprintf("814%d%d", rand.Intn(9000)+1000, rand.Intn(9000)+1000)
	email := fmt.Sprintf("%s%d@gmail.com", strings.ReplaceAll(strings.ToLower(selectedName), " ", ""), time.Now().Unix()%1000)

	// 3. Human-like Typing
	page.Type("#name", selectedName, playwright.PageTypeOptions{Delay: playwright.Float(100)})
	time.Sleep(time.Duration(rand.Intn(2)+1) * time.Second)
	page.Type("#phone", phone, playwright.PageTypeOptions{Delay: playwright.Float(100)})
	time.Sleep(time.Duration(rand.Intn(2)+1) * time.Second)
	page.Type("#email", email, playwright.PageTypeOptions{Delay: playwright.Float(100)})

	// 4. Pilih Paket Acak (A, B, atau C)
	paketOptions := []string{
		"text=Paket A", 
		"text=Paket B", 
		"text=Paket C",
	}
	err = page.Click(paketOptions[rand.Intn(len(paketOptions))])
	if err != nil { return fmt.Errorf("click_paket_failed") }

	// 5. Submit
	// Menggunakan selector spesifik dari HTML: Ambil Sistemnya...
	err = page.Click("button:has-text('Ambil Sistemnya')")
	if err != nil { return fmt.Errorf("submit_btn_failed") }

	// 6. Success Detection
	err = page.WaitForURL("**/success", playwright.PageWaitForURLOptions{Timeout: playwright.Float(30000)})
	if err != nil { return fmt.Errorf("failed_success_redirect") }

	return nil
}
