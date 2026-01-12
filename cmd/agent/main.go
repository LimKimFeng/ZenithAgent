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
	"zenith-agent/internal/manager"
	"zenith-agent/internal/monitor"
	"zenith-agent/internal/network"
	"zenith-agent/internal/notify"
	"zenith-agent/internal/tasks"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	// User Input Section
	fmt.Print("Input SMTP User: ")
	smtpUser, _ := reader.ReadString('\n')
	smtpUser = strings.TrimSpace(smtpUser)

	fmt.Print("Input SMTP Password: ")
	smtpPass, _ := reader.ReadString('\n')
	smtpPass = strings.TrimSpace(smtpPass)

	// Tor Authentication Prompt
	fmt.Print("Input Tor Control Password (leave empty if none): ")
	torPass, _ := reader.ReadString('\n')
	torPass = strings.TrimSpace(torPass)

	fmt.Print("Input Email Penerima Laporan: ")
	recipient, _ := reader.ReadString('\n')
	recipient = strings.TrimSpace(recipient)

	// Interval
	fmt.Print("Interval eksekusi (menit): ")
	intStr, _ := reader.ReadString('\n')
	intervalMinutes, _ := strconv.Atoi(strings.TrimSpace(intStr))

	// Project Selection
	fmt.Println("Pilih Project:")
	fmt.Println("1. Ternak Property")
	fmt.Println("2. Scalping Support & Resistance")
	fmt.Print("Pilihan (1/2): ")
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)
	
	projectName := "Ternak Property"
	if choice == "2" {
		projectName = "Scalping SR"
	}

	fmt.Print("Gunakan Headless mode? (y/n): ")
	hMode, _ := reader.ReadString('\n')
	headless := strings.ToLower(strings.TrimSpace(hMode)) == "y"

	// Singleton Check & Tor Start
	state := manager.GetState()
	if state.IsRunning {
		fmt.Printf("⚠️  Warning: Agent is already running (PID: %d). This instance will run tasks but will NOT handle IP Rotation.\n", state.LastPid)
	} else {
		// Master Process
		manager.UpdateState(true, true)
		defer manager.UpdateState(false, false)
		
		// Start Tor Rotator (Background) with Password
		go network.StartRotator(10, torPass)

		// Start Dashboard Server (Background)
		monitor.StartDashboard("8080")
	}

	// Inisialisasi Notify & Engine
	notifier := notify.NewEmailNotifier(smtpUser, smtpPass, recipient)
	browserManager := engine.NewBrowserManager(headless)

	fmt.Printf("\n--- ZenithAgent Started ---\nTarget: %s\nInterval: %d menit\n\n", projectName, intervalMinutes)

	ticker := time.NewTicker(time.Duration(intervalMinutes) * time.Minute)
	reportTicker := time.NewTicker(24 * time.Hour)

	// Jalankan pertama kali saat start
	runTask(browserManager, notifier, choice)

	for {
		select {
		case <-ticker.C:
			runTask(browserManager, notifier, choice)
		case <-reportTicker.C:
			notifier.SendDailyReport()
		}
	}
}

func runTask(bm *engine.BrowserManager, n *notify.EmailNotifier, choice string) {
	fmt.Printf("[%s] Executing Task...\n", time.Now().Format("15:04:05"))
	
	var err error
	var pName string
	
	if choice == "2" {
		pName = "ZenithAgent-ScalpingSR"
		err = tasks.ExecuteScalpingSR(bm)
	} else {
		pName = "ZenithAgent-TernakProperty"
		err = tasks.ExecuteTernakProperty(bm)
	}

	if err != nil {
		log.Printf("Task Error: %v", err)
		tasks.UpdateStats(pName, false, err.Error())
	} else {
		fmt.Println("Task Success!")
		tasks.UpdateStats(pName, true, "")
	}
}
