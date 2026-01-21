package monitor

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"zenith-agent/internal/tasks"

	"golang.org/x/crypto/bcrypt"
)

// Kredensial yang di-hash (Bcrypt)
const (
	expectedUsername = "zenithagent.admin@cornelweb.com"
)

// Middleware untuk verifikasi Basic Auth
func basicAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()

		// Hardcoded Hash untuk @zenithagent04042008cornel@
		hashPass := "$2a$10$9e/z6htdS4Qrdlngl/sbgOCHXg1pnN8ZX7/kbf74na.MOTsDMScvm"

		err := bcrypt.CompareHashAndPassword([]byte(hashPass), []byte(pass))

		if !ok || user != expectedUsername || (err != nil) {
			w.Header().Set("WWW-Authenticate", `Basic realm="ZenithAgent Restricted"`)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("401 Unauthorized\n"))
			return
		}
		next.ServeHTTP(w, r)
	}
}

func StartDashboard(port string) {
	// API endpoint for stats
	http.HandleFunc("/api/stats", basicAuth(func(w http.ResponseWriter, r *http.Request) {
		data, err := os.ReadFile("stats.json")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Write(data)
	}))

	// New API endpoint for all projects (discovered + stats)
	http.HandleFunc("/api/projects", basicAuth(func(w http.ResponseWriter, r *http.Request) {
		// Initialize task registry
		if err := tasks.InitRegistry(); err != nil {
			http.Error(w, fmt.Sprintf("Failed to discover tasks: %v", err), 500)
			return
		}

		// Read stats.json
		var globalStats tasks.GlobalStats
		data, err := os.ReadFile("stats.json")
		if err == nil {
			json.Unmarshal(data, &globalStats)
		}

		// Build response with all projects
		type ProjectResponse struct {
			Key          string `json:"key"`
			Name         string `json:"name"`
			IsRunning    bool   `json:"is_running"`
			SuccessCount int    `json:"success_count"`
			FailedCount  int    `json:"failed_count"`
			LastRun      string `json:"last_run"`
			HasExecuted  bool   `json:"has_executed"`
		}

		response := struct {
			Projects []ProjectResponse `json:"projects"`
			History  []tasks.LogEntry  `json:"history"`
		}{
			Projects: []ProjectResponse{},
			History:  globalStats.History,
		}

		// Iterate through discovered tasks
		for key, task := range tasks.Registry {
			projectResp := ProjectResponse{
				Key:          key,
				Name:         task.DisplayName,
				IsRunning:    false,
				SuccessCount: 0,
				FailedCount:  0,
				LastRun:      "",
				HasExecuted:  false,
			}

			// Merge with stats if exists
			if stats, exists := globalStats.Projects[task.DisplayName]; exists {
				projectResp.IsRunning = stats.IsRunning
				projectResp.SuccessCount = stats.SuccessCount
				projectResp.FailedCount = stats.FailedCount
				if !stats.LastRun.IsZero() {
					projectResp.LastRun = stats.LastRun.Format("2006-01-02 15:04:05")
				}
				projectResp.HasExecuted = stats.SuccessCount > 0 || stats.FailedCount > 0
			}

			response.Projects = append(response.Projects, projectResp)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*") // Corrected placement
		json.NewEncoder(w).Encode(response)
	}))

	// Check if React build exists
	buildDir := "./dashboard/dist"
	if _, err := os.Stat(buildDir); err == nil {
		// Serve React production build
		fs := http.FileServer(http.Dir(buildDir))
		http.Handle("/", basicAuth(func(w http.ResponseWriter, r *http.Request) {
			// Try to serve the file
			path := buildDir + r.URL.Path
			if _, err := os.Stat(path); os.IsNotExist(err) {
				// File doesn't exist, serve index.html (SPA fallback)
				http.ServeFile(w, r, buildDir+"/index.html")
				return
			}
			fs.ServeHTTP(w, r)
		}))
		fmt.Printf("[MONITOR] Dashboard serving React build (Protected) at http://0.0.0.0:%s\n", port)
	} else {
		// Fallback to embedded HTML for development
		http.HandleFunc("/", basicAuth(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, dashboardHTML)
		}))
		fmt.Printf("[MONITOR] Dashboard HTML active (Protected) at http://0.0.0.0:%s\n", port)
	}

	go func() {
		if err := http.ListenAndServe(":"+port, nil); err != nil {
			fmt.Printf("⚠️ [FAT AL] Dashboard Failed to Start: %v\n", err)
		}
	}()
}

const dashboardHTML = `<!DOCTYPE html>
<html lang="id">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>ZenithAgent Dashboard</title>
    <link href="https://fonts.googleapis.com/css2?family=Inter:wght@300;400;500;600;700&display=swap" rel="stylesheet">
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body { font-family: 'Inter', sans-serif; background: linear-gradient(135deg, #0f172a 0%, #1e293b 100%); color: #e2e8f0; min-height: 100vh; padding: 2rem; }
        .container { max-width: 1400px; margin: 0 auto; }
        .header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 2rem; flex-wrap: wrap; gap: 1rem; }
        .logo { display: flex; align-items: center; gap: 1rem; }
        .logo-icon { width: 48px; height: 48px; background: linear-gradient(135deg, #06b6d4 0%, #3b82f6 100%); border-radius: 12px; display: flex; align-items: center; justify-content: center; box-shadow: 0 8px 16px rgba(6, 182, 212, 0.3); }
        .logo-text h1 { font-size: 1.75rem; font-weight: 700; color: white; }
        .logo-text p { font-size: 0.75rem; color: #94a3b8; font-weight: 500; letter-spacing: 1px; }
        .status-badge { display: flex; align-items: center; gap: 0.75rem; background: rgba(30, 41, 59, 0.7); backdrop-filter: blur(12px); padding: 0.5rem 1rem; border-radius: 999px; border: 1px solid rgba(255, 255, 255, 0.1); }
        .status-dot { width: 12px; height: 12px; border-radius: 50%; background: #10b981; animation: pulse 2s ease-in-out infinite; position: relative; }
        .status-dot::before { content: ''; position: absolute; width: 100%; height: 100%; border-radius: 50%; background: #10b981; opacity: 0.5; animation: ping 2s ease-in-out infinite; }
        @keyframes pulse { 0%, 100% { opacity: 1; } 50% { opacity: 0.6; } }
        @keyframes ping { 0% { transform: scale(1); opacity: 0.5; } 100% { transform: scale(2); opacity: 0; } }
        @keyframes blink { 0%, 50%, 100% { opacity: 1; } 25%, 75% { opacity: 0.3; } }
        .status-text { font-size: 0.75rem; font-weight: 600; color: #10b981; letter-spacing: 0.5px; }
        .task-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(320px, 1fr)); gap: 1.5rem; margin-bottom: 2rem; }
        .task-card { background: rgba(30, 41, 59, 0.5); backdrop-filter: blur(12px); border: 1px solid rgba(148, 163, 184, 0.2); border-radius: 16px; padding: 1.5rem; transition: all 0.3s ease; position: relative; overflow: hidden; }
        .task-card::before { content: ''; position: absolute; top: 0; right: 0; width: 100px; height: 100px; background: rgba(255, 255, 255, 0.03); border-radius: 50%; filter: blur(40px); transform: translate(30%, -30%); }
        .task-card:hover { transform: translateY(-4px); box-shadow: 0 12px 24px rgba(0, 0, 0, 0.3); border-color: rgba(148, 163, 184, 0.4); }
        .task-card.running { border-color: rgba(6, 182, 212, 0.4); box-shadow: 0 8px 16px rgba(6, 182, 212, 0.2); }
        .task-card.error { border-color: rgba(239, 68, 68, 0.4); box-shadow: 0 8px 16px rgba(239, 68, 68, 0.2); }
        .task-header { display: flex; justify-content: space-between; align-items: flex-start; margin-bottom: 1.5rem; position: relative; }
        .task-name { font-size: 1.125rem; font-weight: 600; color: white; margin-bottom: 0.5rem; }
        .task-status { display: flex; align-items: center; gap: 0.5rem; font-size: 0.65rem; font-weight: 700; letter-spacing: 1px; text-transform: uppercase; }
        .task-status-dot { width: 8px; height: 8px; border-radius: 50%; }
        .task-status.running .task-status-dot { background: #06b6d4; animation: pulse 1.5s ease-in-out infinite; }
        .task-status.running { color: #06b6d4; }
        .task-status.idle .task-status-dot { background: #3b82f6; }
        .task-status.idle { color: #3b82f6; }
        .task-status.error .task-status-dot { background: #ef4444; animation: blink 1s ease-in-out infinite; }
        .task-status.error { color: #ef4444; }
        .task-time { font-size: 0.75rem; color: #94a3b8; background: rgba(15, 23, 42, 0.5); padding: 0.375rem 0.75rem; border-radius: 6px; border: 1px solid rgba(148, 163, 184, 0.1); font-family: 'Courier New', monospace; }
        .task-stats { display: grid; grid-template-columns: 1fr 1fr; gap: 0.75rem; margin-bottom: 0.75rem; }
        .stat-box { background: rgba(15, 23, 42, 0.4); padding: 0.75rem; border-radius: 8px; border: 1px solid rgba(148, 163, 184, 0.1); transition: all 0.2s ease; }
        .stat-box:hover { border-color: rgba(148, 163, 184, 0.3); }
        .stat-label { font-size: 0.65rem; color: #94a3b8; font-weight: 600; text-transform: uppercase; letter-spacing: 0.5px; margin-bottom: 0.25rem; }
        .stat-value { font-size: 1.5rem; font-weight: 700; color: white; }
        .stat-box.success .stat-value { color: #10b981; }
        .stat-box.failed .stat-value { color: #ef4444; }
        .error-badge { background: rgba(239, 68, 68, 0.1); color: #ef4444; padding: 0.5rem 0.75rem; border-radius: 8px; font-size: 0.75rem; border: 1px solid rgba(239, 68, 68, 0.3); margin-top: 0.5rem; font-weight: 500; }
        .history-section { background: rgba(30, 41, 59, 0.5); backdrop-filter: blur(12px); border: 1px solid rgba(148, 163, 184, 0.2); border-radius: 16px; overflow: hidden; max-height: 600px; display: flex; flex-direction: column; }
        .history-header { padding: 1.25rem 1.5rem; background: rgba(15, 23, 42, 0.5); border-bottom: 1px solid rgba(148, 163, 184, 0.1); display: flex; justify-content: space-between; align-items: center; }
        .history-title { font-size: 1.125rem; font-weight: 600; color: white; }
        .history-list { flex: 1; overflow-y: auto; padding: 0.5rem; }
        .history-item { padding: 0.75rem 1rem; margin-bottom: 0.5rem; border-radius: 8px; display: flex; align-items: center; gap: 0.75rem; transition: all 0.2s ease; }
        .history-item:hover { background: rgba(15, 23, 42, 0.4); }
        .history-icon { width: 32px; height: 32px; border-radius: 8px; display: flex; align-items: center; justify-content: center; flex-shrink: 0; }
        .history-icon.success { background: rgba(16, 185, 129, 0.1); color: #10b981; }
        .history-icon.error { background: rgba(239, 68, 68, 0.1); color: #ef4444; }
        .history-content { flex: 1; min-width: 0; }
        .history-message { font-size: 0.875rem; font-weight: 500; color: #e2e8f0; margin-bottom: 0.25rem; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
        .history-meta { font-size: 0.75rem; color: #64748b; display: flex; align-items: center; gap: 0.5rem; }
        .history-time { font-family: 'Courier New', monospace; background: rgba(15, 23, 42, 0.5); padding: 0.125rem 0.5rem; border-radius: 4px; }
        .history-list::-webkit-scrollbar { width: 6px; }
        .history-list::-webkit-scrollbar-track { background: transparent; }
        .history-list::-webkit-scrollbar-thumb { background: #475569; border-radius: 999px; }
        @media (max-width: 768px) { body { padding: 1rem; } .header { flex-direction: column; align-items: flex-start; } .task-grid { grid-template-columns: 1fr; } }
        .loading { text-align: center; padding: 3rem; color: #94a3b8; }
        .no-tasks { text-align: center; padding: 3rem; color: #94a3b8; font-size: 0.875rem; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <div class="logo">
                <div class="logo-icon">
                    <svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="white" stroke-width="2">
                        <path d="M13 10V3L4 14h7v7l9-11h-7z"></path>
                    </svg>
                </div>
                <div class="logo-text">
                    <h1>ZenithAgent <span style="color: #06b6d4;">Dashboard</span></h1>
                    <p>AUTONOMOUS TASK MONITORING</p>
                </div>
            </div>
            <div class="status-badge">
                <div class="status-dot"></div>
                <span class="status-text" id="connection-status">SYSTEM ONLINE</span>
            </div>
        </div>
        <div class="task-grid" id="task-grid">
            <div class="loading">Loading tasks...</div>
        </div>
        <div class="history-section">
            <div class="history-header">
                <h2 class="history-title">📊 Live Activity Log</h2>
                <span style="font-size: 0.75rem; color: #64748b;">Auto-refresh: 3s</span>
            </div>
            <div class="history-list" id="history-list">
                <div class="no-tasks">Waiting for activity...</div>
            </div>
        </div>
    </div>
    <script>
        const taskGrid = document.getElementById('task-grid');
        const historyList = document.getElementById('history-list');
        const statusText = document.getElementById('connection-status');
        async function updateDashboard() {
            try {
                const response = await fetch('/api/projects');
                if (!response.ok) throw new Error('Failed to fetch');
                const data = await response.json();
                statusText.textContent = 'SYSTEM ONLINE';
                statusText.style.color = '#10b981';
                if (data.projects && data.projects.length > 0) {
                    taskGrid.innerHTML = data.projects.map(task => {
                        const isRunning = task.is_running;
                        const hasError = task.failed_count > 0 && !isRunning;
                        const statusClass = isRunning ? 'running' : (hasError ? 'error' : 'idle');
                        const statusTextLabel = isRunning ? 'Running' : (hasError ? 'Error' : 'Idle');
                        const cardClass = isRunning ? 'running' : (hasError ? 'error' : '');
                        return ` + "`" + `
                            <div class="task-card ${cardClass}">
                                <div class="task-header">
                                    <div>
                                        <h3 class="task-name">${task.name}</h3>
                                        <div class="task-status ${statusClass}">
                                            <div class="task-status-dot"></div>
                                            <span>${statusTextLabel}</span>
                                        </div>
                                    </div>
                                    <div class="task-time">${task.last_run || 'Never'}</div>
                                </div>
                                <div class="task-stats">
                                    <div class="stat-box success">
                                        <div class="stat-label">✓ Success</div>
                                        <div class="stat-value">${task.success_count}</div>
                                    </div>
                                    <div class="stat-box failed">
                                        <div class="stat-label">✗ Failed</div>
                                        <div class="stat-value">${task.failed_count}</div>
                                    </div>
                                </div>
                                ${hasError ? '<div class="error-badge"><strong>⚠ Last Error:</strong> Check logs for details</div>' : ''}
                            </div>
                        ` + "`" + `;
                    }).join('');
                } else {
                    taskGrid.innerHTML = '<div class="no-tasks">No tasks discovered</div>';
                }
                if (data.history && data.history.length > 0) {
                    historyList.innerHTML = data.history.map(entry => {
                        const isSuccess = entry.type === 'success';
                        const iconClass = isSuccess ? 'success' : 'error';
                        const icon = isSuccess ? '✓' : '✗';
                        return ` + "`" + `
                            <div class="history-item">
                                <div class="history-icon ${iconClass}">${icon}</div>
                                <div class="history-content">
                                    <div class="history-message">${entry.message}</div>
                                    <div class="history-meta">
                                        <span>${entry.project}</span>
                                        <span class="history-time">${entry.timestamp}</span>
                                    </div>
                                </div>
                            </div>
                        ` + "`" + `;
                    }).join('');
                } else {
                    historyList.innerHTML = '<div class="no-tasks">Waiting for activity...</div>';
                }
            } catch (error) {
                console.error('Dashboard error:', error);
                statusText.textContent = 'OFFLINE';
                statusText.style.color = '#ef4444';
                taskGrid.innerHTML = '<div class="no-tasks">Connection error</div>';
            }
        }
        updateDashboard();
        setInterval(updateDashboard, 3000);
    </script>
</body>
</html>
`
