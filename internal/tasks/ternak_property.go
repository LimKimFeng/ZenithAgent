package tasks

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"
	"zenith-agent/internal/engine"

	"github.com/playwright-community/playwright-go"
)

type Stats struct {
	ProjectName  string    `json:"project_name"`
	LastReset    time.Time `json:"last_reset"`
	SuccessCount int       `json:"success_count"`
	FailedCount  int       `json:"failed_count"`
	FailReasons  []string  `json:"fail_reasons"`
}

var (
	statsMutex sync.Mutex
	statsFile  = "stats.json"
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

	// 1. Navigate - FIX: Changed to WaitUntilStateLoad to avoid timeout on trackers
	_, err = page.Goto("https://ternakproperty.com/as-bandung?utm_source=fb&utm_medium=mof&utm_campaign=mof_seminarbandung&utm_content=vidbeforeafterv1", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateLoad,
	})
	if err != nil {
		return err
	}

	// FIX: Wait for form to appear before interacting
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

	// 2. Fill Form with Human Delay (Type instead of Fill)
	// Delay is in milliseconds. 100ms per char is reasonable human speed.
	typeOptions := playwright.PageTypeOptions{Delay: playwright.Float(100)}
	
	page.Type("#name", selectedName, typeOptions)
	page.Type("#phone", phone, typeOptions)
	page.Type("#email", email, typeOptions)
	page.Type("#notes", generateSmartReason(), typeOptions)

	// 3. Select BCA (Base on HTML structure provided)
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
	
	return err
}

func UpdateStats(success bool, reason string) {
	statsMutex.Lock()
	defer statsMutex.Unlock()

	// FIX: Truncate long error messages
	if len(reason) > 100 {
		reason = reason[:97] + "..."
	}

	var s Stats
	data, err := os.ReadFile(statsFile)
	if err == nil {
		json.Unmarshal(data, &s)
	} else {
		s.ProjectName = "ZenithAgent-TernakProperty"
		s.LastReset = time.Now()
	}

	if success {
		s.SuccessCount++
	} else {
		s.FailedCount++
		// Deduplikasi alasan
		found := false
		for _, r := range s.FailReasons {
			if r == reason {
				found = true
				break
			}
		}
		if !found {
			s.FailReasons = append(s.FailReasons, reason)
		}
	}

	newData, _ := json.MarshalIndent(s, "", "  ")
	os.WriteFile(statsFile, newData, 0644)
}
