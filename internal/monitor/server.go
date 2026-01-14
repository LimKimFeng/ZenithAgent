package monitor

import (
	"fmt"
	"net/http"
	"os"

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
			fmt.Printf("⚠️ [FATAL] Dashboard Failed to Start: %v\n", err)
		}
	}()
}

const dashboardHTML = `
<!DOCTYPE html>
<html lang="id">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>ZenithAgent 2.0</title>
    <script src="https://cdn.tailwindcss.com"></script>
    <link href="https://fonts.googleapis.com/css2?family=Inter:wght@300;400;500;600;700&display=swap" rel="stylesheet">
    <style>
        body { font-family: 'Inter', sans-serif; }
        .glass { background: rgba(30, 41, 59, 0.7); backdrop-filter: blur(12px); -webkit-backdrop-filter: blur(12px); border: 1px solid rgba(255, 255, 255, 0.08); }
        .glass-card { background: rgba(30, 41, 59, 0.5); backdrop-filter: blur(8px); border: 1px solid rgba(148, 163, 184, 0.1); }
        .animate-pulse-slow { animation: pulse 3s cubic-bezier(0.4, 0, 0.6, 1) infinite; }
        @keyframes pulse { 0%, 100% { opacity: 1; } 50% { opacity: .6; } }
        ::-webkit-scrollbar { width: 6px; }
        ::-webkit-scrollbar-track { background: transparent; }
        ::-webkit-scrollbar-thumb { background: #475569; border-radius: 99px; }
    </style>
</head>
<body class="bg-[#0f172a] text-slate-200 min-h-screen relative overflow-x-hidden selection:bg-cyan-500/30 selection:text-cyan-200">
    
    <!-- Background Glows -->
    <div class="fixed top-0 left-0 w-full h-full overflow-hidden pointer-events-none z-0">
        <div class="absolute top-[-20%] left-[-10%] w-[50%] h-[50%] bg-blue-600/10 rounded-full blur-[120px]"></div>
        <div class="absolute bottom-[-20%] right-[-10%] w-[50%] h-[50%] bg-emerald-600/10 rounded-full blur-[120px] mix-blend-screen"></div>
    </div>

    <div class="max-w-6xl mx-auto p-6 md:p-8 relative z-10">
        <!-- Header -->
        <header class="flex flex-col md:flex-row justify-between items-center mb-10 gap-6">
            <div class="flex items-center gap-4">
                <div class="w-12 h-12 rounded-xl bg-gradient-to-br from-cyan-500 to-blue-600 shadow-lg shadow-cyan-500/20 flex items-center justify-center">
                    <svg class="w-7 h-7 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z"></path></svg>
                </div>
                <div>
                    <h1 class="text-2xl font-bold text-white tracking-tight">ZenithAgent <span class="text-cyan-400">2.0</span></h1>
                    <p class="text-slate-400 text-xs font-medium tracking-wide">AUTONOMOUS STEALTH ENGINE</p>
                </div>
            </div>
            
            <div class="flex items-center gap-4">
                <div class="glass px-4 py-2 rounded-full flex items-center gap-3 shadow-lg shadow-black/20">
                    <div class="relative">
                        <span class="flex h-3 w-3">
                            <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75"></span>
                            <span class="relative inline-flex rounded-full h-3 w-3 bg-emerald-500"></span>
                        </span>
                    </div>
                    <span class="text-xs font-semibold text-emerald-400 tracking-wide" id="connection-status">SYSTEM ONLINE</span>
                </div>
                <div class="w-10 h-10 rounded-full bg-slate-800 border border-slate-700 flex items-center justify-center text-slate-400 hover:text-white transition-colors cursor-pointer" title="Admin">
                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"></path></svg>
                </div>
            </div>
        </header>

        <!-- Stats Overview -->
        <div id="project-cards" class="grid grid-cols-1 md:grid-cols-2 gap-6 mb-8">
            <!-- Dynamic Content -->
        </div>

        <!-- History Log -->
        <div class="glass rounded-2xl overflow-hidden shadow-2xl shadow-black/40 flex flex-col h-[500px]">
            <div class="px-6 py-5 border-b border-slate-700/50 bg-slate-900/40 flex justify-between items-center backdrop-blur-md">
                <div class="flex items-center gap-3">
                    <div class="p-2 rounded-lg bg-indigo-500/10 text-indigo-400">
                        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
                    </div>
                    <h2 class="font-semibold text-lg text-white">Live Activity Log</h2>
                </div>
                <div class="flex items-center gap-2 text-xs font-mono text-slate-500">
                    <span class="w-2 h-2 rounded-full bg-slate-600 animate-pulse"></span>
                    SYNCING
                </div>
            </div>
            
            <div id="history-list" class="flex-1 overflow-y-auto p-2 space-y-1">
                <!-- Log Items -->
            </div>
        </div>
    </div>

    <!-- Scripts -->
    <script>
        const statusEl = document.getElementById('connection-status');
        const cardsContainer = document.getElementById('project-cards');
        const historyContainer = document.getElementById('history-list');

        function formatTime(isoString) {
            return new Date(isoString).toLocaleTimeString('id-ID', { hour12: false, hour: '2-digit', minute:'2-digit', second:'2-digit' });
        }

        async function updateDashboard() {
            try {
                const response = await fetch('/api/stats');
                if (!response.ok) throw new Error('Auth Failed');
                const data = await response.json();
                
                // Status Update
                statusEl.innerText = "SYSTEM ONLINE";
                statusEl.className = "text-xs font-semibold text-emerald-400 tracking-wide";

                // Render Cards
                cardsContainer.innerHTML = '';
                if (data.projects) {
                    for (const [name, stats] of Object.entries(data.projects)) {
                        const isRunning = stats.is_running;
                        const cardClass = isRunning ? 'border-cyan-500/30 shadow-lg shadow-cyan-900/20' : 'border-rose-500/20 opacity-90 grayscale-[0.3]';
                        const statusDot = isRunning ? 'bg-cyan-400 shadow-[0_0_8px_rgba(34,211,238,0.6)]' : 'bg-rose-500';
                        const statusText = isRunning ? 'ACTIVE RUNNING' : 'PROCESS HALTED';
                        const statusTextClass = isRunning ? 'text-cyan-400' : 'text-rose-400';

                        cardsContainer.innerHTML += ` + "`" + `
                            <div class="glass p-1 rounded-2xl transition-all duration-500 hover:translate-y-[-2px] ${cardClass}">
                                <div class="bg-slate-900/60 rounded-xl p-6 h-full relative overflow-hidden group">
                                    <div class="absolute top-0 right-0 w-32 h-32 bg-white/5 rounded-full blur-2xl translate-x-10 translate-y-[-10px] group-hover:bg-white/10 transition-all duration-700"></div>
                                    
                                    <div class="flex justify-between items-start mb-6 relative">
                                        <div>
                                            <h3 class="text-lg font-bold text-white tracking-tight mb-1">${name}</h3>
                                            <div class="flex items-center gap-2">
                                                <span class="w-1.5 h-1.5 rounded-full ${statusDot} animate-pulse-slow"></span>
                                                <span class="text-[10px] font-bold tracking-widest ${statusTextClass}">${statusText}</span>
                                            </div>
                                        </div>
                                        <div class="text-right">
                                            <p class="text-[10px] text-slate-500 uppercase tracking-wider font-semibold">Last Cycle</p>
                                            <p class="text-xs font-mono text-slate-300 bg-slate-800/80 px-2 py-1 rounded border border-slate-700/50 mt-1 inline-block">
                                                ${new Date(stats.last_run).toLocaleTimeString('id-ID')}
                                            </p>
                                        </div>
                                    </div>

                                    <div class="grid grid-cols-2 gap-3 relative">
                                        <div class="bg-slate-800/40 p-3 rounded-lg border border-emerald-500/10 hover:border-emerald-500/30 transition-colors group/stat">
                                            <div class="flex items-center justify-between mb-1">
                                                <p class="text-[10px] uppercase font-bold text-slate-500 group-hover/stat:text-emerald-400 transition-colors">Success</p>
                                                <svg class="w-3 h-3 text-emerald-500 opacity-50" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M5 13l4 4L19 7"></path></svg>
                                            </div>
                                            <p class="text-2xl font-black text-white group-hover/stat:text-emerald-300 transition-colors">${stats.success_count}</p>
                                        </div>
                                        <div class="bg-slate-800/40 p-3 rounded-lg border border-rose-500/10 hover:border-rose-500/30 transition-colors group/stat">
                                            <div class="flex items-center justify-between mb-1">
                                                <p class="text-[10px] uppercase font-bold text-slate-500 group-hover/stat:text-rose-400 transition-colors">Failures</p>
                                                <svg class="w-3 h-3 text-rose-500 opacity-50" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M6 18L18 6M6 6l12 12"></path></svg>
                                            </div>
                                            <p class="text-2xl font-black text-white group-hover/stat:text-rose-300 transition-colors">${stats.failed_count}</p>
                                        </div>
                                    </div>
                                </div>
                            </div>
                        ` + "`" + `;
                    }
                }

                // Render History
                if (data.history && data.history.length > 0) {
                    historyContainer.innerHTML = data.history.map((entry, index) => {
                        const isSuccess = entry.type === 'success';
                        const iconBg = isSuccess ? 'bg-emerald-500/10 text-emerald-400' : 'bg-rose-500/10 text-rose-400';
                        const iconPath = isSuccess ? 'M5 13l4 4L19 7' : 'M6 18L18 6M6 6l12 12';
                        const delay = index * 50; // Stagger animation
                        
                        return ` + "`" + `
                        <div class="group mx-2 p-3 rounded-xl hover:bg-slate-800/50 transition-all border border-transparent hover:border-slate-700/50 flex items-start gap-4 animate-fade-in" style="animation-delay: ${delay}ms">
                            <div class="mt-1 w-8 h-8 rounded-lg ${iconBg} flex items-center justify-center shrink-0 border border-white/5 shadow-sm">
                                <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="${iconPath}"></path></svg>
                            </div>
                            <div class="flex-1 min-w-0">
                                <div class="flex justify-between items-baseline mb-0.5">
                                    <h4 class="text-sm font-semibold text-slate-200 group-hover:text-white truncate pr-4">${entry.message}</h4>
                                    <span class="text-[10px] font-mono text-slate-500 bg-slate-800 px-1.5 py-0.5 rounded border border-slate-700/50 group-hover:border-slate-600 transition-colors whitespace-nowrap">${entry.timestamp}</span>
                                </div>
                                <p class="text-xs text-slate-500 font-medium tracking-wide flex items-center gap-1.5">
                                    <span class="w-1 h-1 rounded-full bg-slate-600"></span>
                                    ${entry.project}
                                </p>
                            </div>
                        </div>
                        ` + "`" + `
                    }).join('');
                } else {
                    historyContainer.innerHTML = ` + "`" + `
                        <div class="flex flex-col items-center justify-center h-full text-slate-500 opacity-60">
                            <svg class="w-12 h-12 mb-3" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
                            <p class="text-sm font-medium">Waiting for activity...</p>
                        </div>
                    ` + "`" + `;
                }

            } catch (e) {
                console.error(e);
                statusEl.innerText = "OFFLINE";
                statusEl.parentElement.className = "bg-rose-950/30 px-4 py-2 rounded-full flex items-center gap-3 border border-rose-500/20";
                statusEl.className = "text-xs font-semibold text-rose-500 tracking-wide";
                statusEl.previousElementSibling.firstElementChild.className = "absolute inline-flex h-full w-full rounded-full bg-rose-500 opacity-20";
                statusEl.previousElementSibling.lastElementChild.className = "relative inline-flex rounded-full h-3 w-3 bg-rose-600";
            }
        }

        setInterval(updateDashboard, 3000);
        updateDashboard();
    </script>
</body>
</html>
`
