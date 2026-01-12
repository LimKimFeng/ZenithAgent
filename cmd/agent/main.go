package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
	"zenith-agent/internal/engine"
	"zenith-agent/internal/notify"
	"zenith-agent/internal/tasks"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Input SMTP User: ")
	smtpUser, _ := reader.ReadString('\n')
	smtpUser = strings.TrimSpace(smtpUser)

	fmt.Print("Input SMTP Password: ")
	smtpPass, _ := reader.ReadString('\n')
	smtpPass = strings.TrimSpace(smtpPass)

	fmt.Print("Input Email Penerima Laporan: ")
	recipient, _ := reader.ReadString('\n')
	recipient = strings.TrimSpace(recipient)

	fmt.Print("Interval eksekusi (menit): ")
	intStr, _ := reader.ReadString('\n')
	intervalMinutes, _ := strconv.Atoi(strings.TrimSpace(intStr))

	fmt.Print("Gunakan Headless mode? (y/n): ")
	hMode, _ := reader.ReadString('\n')
	headless := strings.ToLower(strings.TrimSpace(hMode)) == "y"

	// Inisialisasi Notify & Engine
	notifier := notify.NewEmailNotifier(smtpUser, smtpPass, recipient)
	browserManager := engine.NewBrowserManager(headless) // FIX: Passing argument headless

	fmt.Printf("\n--- ZenithAgent Started ---\nTarget: Ternak Property\nInterval: %d menit\n\n", intervalMinutes)

	ticker := time.NewTicker(time.Duration(intervalMinutes) * time.Minute)
	reportTicker := time.NewTicker(24 * time.Hour)

	// Jalankan pertama kali saat start
	runTask(browserManager, notifier)

	for {
		select {
		case <-ticker.C:
			runTask(browserManager, notifier)
		case <-reportTicker.C:
			notifier.SendDailyReport()
		}
	}
}

func runTask(bm *engine.BrowserManager, n *notify.EmailNotifier) {
	fmt.Printf("[%s] Executing Task...\n", time.Now().Format("15:04:05"))
	
	err := tasks.ExecuteTernakProperty(bm)
	if err != nil {
		log.Printf("Task Error: %v", err)
		tasks.UpdateStats(false, err.Error())
	} else {
		fmt.Println("Task Success!")
		tasks.UpdateStats(true, "")
	}
}
