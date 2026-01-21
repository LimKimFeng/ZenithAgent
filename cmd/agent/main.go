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

	fmt.Print("Gunakan Headless mode? (y/n): ")
	hMode, _ := reader.ReadString('\n')
	headless := strings.ToLower(strings.TrimSpace(hMode)) == "y"

	// Initialize task registry
	fmt.Println("\n🔍 Discovering available tasks...")
	if err := tasks.InitRegistry(); err != nil {
		fmt.Printf("❌ Failed to discover tasks: %v\n", err)
		os.Exit(1)
	}

	// Display available projects
	fmt.Println("\n🎯 Available Projects:")
	taskList := tasks.GetTaskList()

	if len(taskList) == 0 {
		fmt.Println("❌ No tasks found!")
		os.Exit(1)
	}

	// Display menu
	for i, key := range taskList {
		task, _ := tasks.GetTask(key)
		fmt.Printf("%d. %s\n", i+1, task.DisplayName)
	}

	fmt.Printf("Masukkan pilihan (1-%d): ", len(taskList))
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	// Parse choice
	choiceNum, err := strconv.Atoi(choice)
	if err != nil || choiceNum < 1 || choiceNum > len(taskList) {
		fmt.Println("❌ Pilihan tidak valid!")
		os.Exit(1)
	}

	selectedProject := taskList[choiceNum-1]
	selectedTask, _ := tasks.GetTask(selectedProject)

	fmt.Printf("\n✓ Selected: %s\n", selectedTask.DisplayName)

	// Acquire lock to prevent multiple instances
	// if err := manager.AcquireLock(); err != nil {
	// 	fmt.Printf("\n%v\n\n", err)
	// 	fmt.Println("If you believe this is an error, remove the .zenith.lock file manually.")
	// 	os.Exit(1)
	// }
	// defer manager.ReleaseLock()

	// --- PERBAIKAN: DASHBOARD DINYALAKAN DI LUAR KONDISI STATE ---
	// Ini memastikan port 8080 selalu aktif begitu binary dijalankan
	monitor.StartDashboard("8080")

	// Singleton Check & Tor Start
	state := manager.GetState()
	if state.IsRunning {
		fmt.Printf("⚠️  Cleaning up previous state (PID: %d)...\n", state.LastPid)
		manager.UpdateState(true, true)
	} else {
		// Master Process
		manager.UpdateState(true, true)
		defer manager.UpdateState(false, false)
	}

	// Start Tor Rotator (Background) with Password
	go network.StartRotator(10, torPass)

	// Inisialisasi Notify & Engine
	notifier := notify.NewEmailNotifier(smtpUser, smtpPass, recipient)
	browserManager := engine.NewBrowserManager(headless)

	fmt.Printf("\n--- ZenithAgent Started ---\nTarget: %s\nInterval: %d menit\n\n", selectedProject, intervalMinutes)

	ticker := time.NewTicker(time.Duration(intervalMinutes) * time.Minute)
	reportTicker := time.NewTicker(24 * time.Hour)

	// Jalankan pertama kali saat start
	runTask(browserManager, notifier, selectedProject)

	for {
		select {
		case <-ticker.C:
			runTask(browserManager, notifier, selectedProject)
		case <-reportTicker.C:
			notifier.SendDailyReport()
		}
	}
}

func runTask(bm *engine.BrowserManager, n *notify.EmailNotifier, projectKey string) {
	fmt.Printf("[%s] Executing Task...\n", time.Now().Format("15:04:05"))

	// Get task from registry
	task, exists := tasks.GetTask(projectKey)
	if !exists {
		log.Printf("Task not found: %s", projectKey)
		return
	}

	// Execute task
	err := task.Function(bm)

	// Update stats
	tasks.GlobalUpdateStats(task.DisplayName, err == nil, "")

	if err != nil {
		log.Printf("Task Error: %v", err)
		tasks.GlobalUpdateStats(task.DisplayName, false, err.Error())
	} else {
		fmt.Println("Task Success!")
	}
}
