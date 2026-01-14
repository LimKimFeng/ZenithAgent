import React from 'react';

interface LogEntry {
    timestamp: string;
    project: string;
    message: string;
    type: string;
}

interface ActivityLogProps {
    history: LogEntry[];
}

const ActivityLog: React.FC<ActivityLogProps> = ({ history }) => {
    if (!history || history.length === 0) {
        return (
            <div className="flex flex-col items-center justify-center h-full text-slate-500 opacity-60">
                <svg
                    className="w-12 h-12 mb-3"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                >
                    <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth="1.5"
                        d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"
                    ></path>
                </svg>
                <p className="text-sm font-medium">Waiting for activity...</p>
            </div>
        );
    }

    return (
        <div className="space-y-1">
            {history.map((entry, index) => {
                const isSuccess = entry.type === 'success';
                const iconBg = isSuccess
                    ? 'bg-emerald-500/10 text-emerald-400'
                    : 'bg-rose-500/10 text-rose-400';
                const iconPath = isSuccess ? 'M5 13l4 4L19 7' : 'M6 18L18 6M6 6l12 12';
                const delay = index * 50; // Stagger animation

                return (
                    <div
                        key={index}
                        className="group mx-2 p-3 rounded-xl hover:bg-slate-800/50 transition-all border border-transparent hover:border-slate-700/50 flex items-start gap-4 animate-fade-in"
                        style={{ animationDelay: `${delay}ms` }}
                    >
                        <div
                            className={`mt-1 w-8 h-8 rounded-lg ${iconBg} flex items-center justify-center shrink-0 border border-white/5 shadow-sm`}
                        >
                            <svg
                                className="w-4 h-4"
                                fill="none"
                                viewBox="0 0 24 24"
                                stroke="currentColor"
                            >
                                <path
                                    strokeLinecap="round"
                                    strokeLinejoin="round"
                                    strokeWidth="2.5"
                                    d={iconPath}
                                ></path>
                            </svg>
                        </div>
                        <div className="flex-1 min-w-0">
                            <div className="flex justify-between items-baseline mb-0.5">
                                <h4 className="text-sm font-semibold text-slate-200 group-hover:text-white truncate pr-4">
                                    {entry.message}
                                </h4>
                                <span className="text-[10px] font-mono text-slate-500 bg-slate-800 px-1.5 py-0.5 rounded border border-slate-700/50 group-hover:border-slate-600 transition-colors whitespace-nowrap">
                                    {entry.timestamp}
                                </span>
                            </div>
                            <p className="text-xs text-slate-500 font-medium tracking-wide flex items-center gap-1.5">
                                <span className="w-1 h-1 rounded-full bg-slate-600"></span>
                                {entry.project}
                            </p>
                        </div>
                    </div>
                );
            })}
        </div>
    );
};

export default ActivityLog;
