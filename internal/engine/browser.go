package engine

import (
	"github.com/playwright-community/playwright-go"
)

type BrowserManager struct {
	Headless bool
}

func NewBrowserManager(headless bool) *BrowserManager {
	return &BrowserManager{Headless: headless}
}

func (bm *BrowserManager) CreateContext() (*playwright.Playwright, playwright.Browser, playwright.BrowserContext, error) {
	pw, err := playwright.Run()
	if err != nil {
		return nil, nil, nil, err
	}

	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(bm.Headless),
		Args: []string{
			"--disable-blink-features=AutomationControlled",
			"--no-sandbox",
		},
		Proxy: &playwright.Proxy{
			Server: "socks5://127.0.0.1:9050",
		},
	})
	if err != nil {
		pw.Stop()
		return nil, nil, nil, err
	}

	// Simulasi User Agent Manusia
	context, err := browser.NewContext(playwright.BrowserNewContextOptions{
		UserAgent: playwright.String("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"),
	})
	if err != nil {
		browser.Close()
		pw.Stop()
		return nil, nil, nil, err
	}

	return pw, browser, context, err
}
