import React from 'react';

interface ProjectStats {
    is_running: boolean;
    success_count: number;
    failed_count: number;
    fail_reasons: string[];
    last_run: string;
}

interface ProjectCardProps {
    name: string;
    stats: ProjectStats;
}

const ProjectCard: React.FC<ProjectCardProps> = ({ name, stats }) => {
    const isRunning = stats.is_running;
    const cardClass = isRunning
        ? 'border-cyan-500/30 shadow-lg shadow-cyan-900/20'
        : 'border-rose-500/20 opacity-90 grayscale-[0.3]';
    const statusDot = isRunning
        ? 'bg-cyan-400 shadow-[0_0_8px_rgba(34,211,238,0.6)]'
        : 'bg-rose-500';
    const statusText = isRunning ? 'ACTIVE RUNNING' : 'PROCESS HALTED';
    const statusTextClass = isRunning ? 'text-cyan-400' : 'text-rose-400';

    const formatTime = (isoString: string) => {
        try {
            return new Date(isoString).toLocaleTimeString('id-ID', {
                hour12: false,
                hour: '2-digit',
                minute: '2-digit',
                second: '2-digit',
            });
        } catch {
            return '--:--:--';
        }
    };

    return (
        <div
            className={`glass p-1 rounded-2xl transition-all duration-500 hover:translate-y-[-2px] ${cardClass}`}
        >
            <div className="bg-slate-900/60 rounded-xl p-6 h-full relative overflow-hidden group">
                <div className="absolute top-0 right-0 w-32 h-32 bg-white/5 rounded-full blur-2xl translate-x-10 translate-y-[-10px] group-hover:bg-white/10 transition-all duration-700"></div>

                <div className="flex justify-between items-start mb-6 relative">
                    <div>
                        <h3 className="text-lg font-bold text-white tracking-tight mb-1">
                            {name}
                        </h3>
                        <div className="flex items-center gap-2">
                            <span
                                className={`w-1.5 h-1.5 rounded-full ${statusDot} animate-pulse-slow`}
                            ></span>
                            <span
                                className={`text-[10px] font-bold tracking-widest ${statusTextClass}`}
                            >
                                {statusText}
                            </span>
                        </div>
                    </div>
                    <div className="text-right">
                        <p className="text-[10px] text-slate-500 uppercase tracking-wider font-semibold">
                            Last Cycle
                        </p>
                        <p className="text-xs font-mono text-slate-300 bg-slate-800/80 px-2 py-1 rounded border border-slate-700/50 mt-1 inline-block">
                            {formatTime(stats.last_run)}
                        </p>
                    </div>
                </div>

                <div className="grid grid-cols-2 gap-3 relative">
                    <div className="bg-slate-800/40 p-3 rounded-lg border border-emerald-500/10 hover:border-emerald-500/30 transition-colors group/stat">
                        <div className="flex items-center justify-between mb-1">
                            <p className="text-[10px] uppercase font-bold text-slate-500 group-hover/stat:text-emerald-400 transition-colors">
                                Success
                            </p>
                            <svg
                                className="w-3 h-3 text-emerald-500 opacity-50"
                                fill="none"
                                viewBox="0 0 24 24"
                                stroke="currentColor"
                            >
                                <path
                                    strokeLinecap="round"
                                    strokeLinejoin="round"
                                    strokeWidth="3"
                                    d="M5 13l4 4L19 7"
                                ></path>
                            </svg>
                        </div>
                        <p className="text-2xl font-black text-white group-hover/stat:text-emerald-300 transition-colors">
                            {stats.success_count}
                        </p>
                    </div>
                    <div className="bg-slate-800/40 p-3 rounded-lg border border-rose-500/10 hover:border-rose-500/30 transition-colors group/stat">
                        <div className="flex items-center justify-between mb-1">
                            <p className="text-[10px] uppercase font-bold text-slate-500 group-hover/stat:text-rose-400 transition-colors">
                                Failures
                            </p>
                            <svg
                                className="w-3 h-3 text-rose-500 opacity-50"
                                fill="none"
                                viewBox="0 0 24 24"
                                stroke="currentColor"
                            >
                                <path
                                    strokeLinecap="round"
                                    strokeLinejoin="round"
                                    strokeWidth="3"
                                    d="M6 18L18 6M6 6l12 12"
                                ></path>
                            </svg>
                        </div>
                        <p className="text-2xl font-black text-white group-hover/stat:text-rose-300 transition-colors">
                            {stats.failed_count}
                        </p>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default ProjectCard;
