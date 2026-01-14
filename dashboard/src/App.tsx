import { useState, useEffect } from 'react';
import ProjectCard from './components/ProjectCard';
import ActivityLog from './components/ActivityLog';
import './index.css';

interface ProjectStats {
  is_running: boolean;
  success_count: number;
  failed_count: number;
  fail_reasons: string[];
  last_run: string;
}

interface LogEntry {
  timestamp: string;
  project: string;
  message: string;
  type: string;
}

interface StatsData {
  projects: Record<string, ProjectStats>;
  history: LogEntry[];
}

const AUTH_CREDENTIALS = btoa('zenithagent.admin@cornelweb.com:@zenithagent04042008cornel@');

function App() {
  const [stats, setStats] = useState<StatsData | null>(null);
  const [connectionStatus, setConnectionStatus] = useState<'online' | 'offline'>('offline');
  const [error, setError] = useState<string | null>(null);

  const updateDashboard = async () => {
    try {
      const response = await fetch('/api/stats', {
        headers: {
          'Authorization': `Basic ${AUTH_CREDENTIALS}`,
        },
      });

      if (!response.ok) {
        throw new Error('Authentication failed');
      }

      const data = await response.json();
      setStats(data);
      setConnectionStatus('online');
      setError(null);
    } catch (err) {
      console.error('Failed to fetch stats:', err);
      setConnectionStatus('offline');
      setError('Failed to connect to backend');
    }
  };

  useEffect(() => {
    updateDashboard();
    const interval = setInterval(updateDashboard, 3000);
    return () => clearInterval(interval);
  }, []);

  return (
    <div className="min-h-screen relative overflow-x-hidden">
      {/* Background Glows */}
      <div className="fixed top-0 left-0 w-full h-full overflow-hidden pointer-events-none z-0">
        <div className="absolute top-[-20%] left-[-10%] w-[50%] h-[50%] bg-blue-600/10 rounded-full blur-[120px]"></div>
        <div className="absolute bottom-[-20%] right-[-10%] w-[50%] h-[50%] bg-emerald-600/10 rounded-full blur-[120px] mix-blend-screen"></div>
      </div>

      <div className="max-w-6xl mx-auto p-6 md:p-8 relative z-10">
        {/* Header */}
        <header className="flex flex-col md:flex-row justify-between items-center mb-10 gap-6">
          <div className="flex items-center gap-4">
            <div className="w-12 h-12 rounded-xl bg-gradient-to-br from-cyan-500 to-blue-600 shadow-lg shadow-cyan-500/20 flex items-center justify-center">
              <svg
                className="w-7 h-7 text-white"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth="2"
                  d="M13 10V3L4 14h7v7l9-11h-7z"
                ></path>
              </svg>
            </div>
            <div>
              <h1 className="text-2xl font-bold text-white tracking-tight">
                ZenithAgent <span className="text-cyan-400">2.0</span>
              </h1>
              <p className="text-slate-400 text-xs font-medium tracking-wide">
                AUTONOMOUS STEALTH ENGINE
              </p>
            </div>
          </div>

          <div className="flex items-center gap-4">
            <div
              className={`glass px-4 py-2 rounded-full flex items-center gap-3 shadow-lg shadow-black/20 ${connectionStatus === 'offline'
                  ? 'bg-rose-950/30 border-rose-500/20'
                  : ''
                }`}
            >
              <div className="relative">
                <span className="flex h-3 w-3">
                  <span
                    className={`animate-ping absolute inline-flex h-full w-full rounded-full ${connectionStatus === 'online'
                        ? 'bg-emerald-400 opacity-75'
                        : 'bg-rose-500 opacity-20'
                      }`}
                  ></span>
                  <span
                    className={`relative inline-flex rounded-full h-3 w-3 ${connectionStatus === 'online'
                        ? 'bg-emerald-500'
                        : 'bg-rose-600'
                      }`}
                  ></span>
                </span>
              </div>
              <span
                className={`text-xs font-semibold tracking-wide ${connectionStatus === 'online'
                    ? 'text-emerald-400'
                    : 'text-rose-500'
                  }`}
              >
                {connectionStatus === 'online' ? 'SYSTEM ONLINE' : 'OFFLINE'}
              </span>
            </div>
            <div
              className="w-10 h-10 rounded-full bg-slate-800 border border-slate-700 flex items-center justify-center text-slate-400 hover:text-white transition-colors cursor-pointer"
              title="Admin"
            >
              <svg
                className="w-5 h-5"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth="2"
                  d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"
                ></path>
              </svg>
            </div>
          </div>
        </header>

        {/* Error Message */}
        {error && (
          <div className="mb-6 glass p-4 rounded-lg border-rose-500/30 bg-rose-900/20">
            <p className="text-rose-400 text-sm">{error}</p>
          </div>
        )}

        {/* Stats Overview */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-8">
          {stats?.projects &&
            Object.entries(stats.projects).map(([name, projectStats]) => (
              <ProjectCard key={name} name={name} stats={projectStats} />
            ))}
          {(!stats?.projects || Object.keys(stats.projects).length === 0) && (
            <div className="col-span-2 glass rounded-2xl p-8 text-center">
              <p className="text-slate-400">No projects running yet...</p>
            </div>
          )}
        </div>

        {/* History Log */}
        <div className="glass rounded-2xl overflow-hidden shadow-2xl shadow-black/40 flex flex-col h-[500px]">
          <div className="px-6 py-5 border-b border-slate-700/50 bg-slate-900/40 flex justify-between items-center backdrop-blur-md">
            <div className="flex items-center gap-3">
              <div className="p-2 rounded-lg bg-indigo-500/10 text-indigo-400">
                <svg
                  className="w-5 h-5"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth="2"
                    d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"
                  ></path>
                </svg>
              </div>
              <h2 className="font-semibold text-lg text-white">
                Live Activity Log
              </h2>
            </div>
            <div className="flex items-center gap-2 text-xs font-mono text-slate-500">
              <span className="w-2 h-2 rounded-full bg-slate-600 animate-pulse"></span>
              SYNCING
            </div>
          </div>

          <div className="flex-1 overflow-y-auto p-2">
            <ActivityLog history={stats?.history || []} />
          </div>
        </div>
      </div>
    </div>
  );
}

export default App;
