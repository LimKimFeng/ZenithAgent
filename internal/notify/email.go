package notify

import (
	"encoding/json"
	"fmt"
	"net/smtp"
	"os"
	"strings"
	"time"
)

type EmailNotifier struct {
	User      string
	Password  string
	Recipient string
	Host      string
	Port      string
}

type ProjectStats struct {
	IsRunning    bool      `json:"is_running"`
	SuccessCount int       `json:"success_count"`
	FailedCount  int       `json:"failed_count"`
	FailReasons  []string  `json:"fail_reasons"`
	LastRun      time.Time `json:"last_run"`
}

type LogEntry struct {
	Timestamp string `json:"timestamp"`
	Project   string `json:"project"`
	Message   string `json:"message"`
	Type      string `json:"type"`
}

type GlobalStats struct {
	Projects map[string]*ProjectStats `json:"projects"`
	History  []LogEntry               `json:"history"`
}

func NewEmailNotifier(user, password, recipient string) *EmailNotifier {
	return &EmailNotifier{
		User:      user,
		Password:  password,
		Recipient: recipient,
		Host:      "smtp.gmail.com", // Default
		Port:      "587",           // Default
	}
}

func (n *EmailNotifier) SendDailyReport() error {
	// Read stats.json
	data, err := os.ReadFile("stats.json")
	if err != nil {
		return n.SendEmail("Daily Status Report - Error", fmt.Sprintf("Failed to read statistics: %v", err))
	}

	var stats GlobalStats
	if err := json.Unmarshal(data, &stats); err != nil {
		return n.SendEmail("Daily Status Report - Error", fmt.Sprintf("Failed to parse statistics: %v", err))
	}

	// Generate HTML email body
	body := n.generateHTMLReport(stats)
	return n.SendHTMLEmail("Daily Status Report", body)
}

func (n *EmailNotifier) generateHTMLReport(stats GlobalStats) string {
	var html strings.Builder

	html.WriteString(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <style>
        body { font-family: 'Segoe UI', Arial, sans-serif; background-color: #f4f4f4; margin: 0; padding: 20px; }
        .container { max-width: 800px; margin: 0 auto; background: white; border-radius: 8px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
        .header { background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); color: white; padding: 30px; border-radius: 8px 8px 0 0; }
        .header h1 { margin: 0; font-size: 28px; }
        .header p { margin: 10px 0 0 0; opacity: 0.9; }
        .content { padding: 30px; }
        .project-card { background: #f9fafb; border-left: 4px solid #667eea; padding: 20px; margin-bottom: 20px; border-radius: 4px; }
        .project-name { font-size: 20px; font-weight: bold; color: #1f2937; margin-bottom: 10px; }
        .stats { display: flex; gap: 20px; margin-top: 15px; }
        .stat { flex: 1; }
        .stat-label { font-size: 12px; color: #6b7280; text-transform: uppercase; font-weight: 600; }
        .stat-value { font-size: 32px; font-weight: bold; margin-top: 5px; }
        .success { color: #10b981; }
        .failure { color: #ef4444; }
        .status { display: inline-block; padding: 4px 12px; border-radius: 12px; font-size: 11px; font-weight: 600; text-transform: uppercase; }
        .status-running { background: #d1fae5; color: #065f46; }
        .status-stopped { background: #fee2e2; color: #991b1b; }
        .history { margin-top: 30px; }
        .history h2 { color: #1f2937; margin-bottom: 15px; }
        .log-entry { padding: 12px; margin-bottom: 8px; border-radius: 4px; display: flex; align-items: center; gap: 12px; }
        .log-success { background: #f0fdf4; border-left: 3px solid #10b981; }
        .log-error { background: #fef2f2; border-left: 3px solid #ef4444; }
        .log-icon { width: 20px; height: 20px; border-radius: 50%; display: flex; align-items: center; justify-content: center; font-size: 12px; }
        .log-icon-success { background: #10b981; color: white; }
        .log-icon-error { background: #ef4444; color: white; }
        .log-details { flex: 1; }
        .log-message { font-weight: 500; color: #1f2937; }
        .log-meta { font-size: 12px; color: #6b7280; margin-top: 4px; }
        .footer { text-align: center; padding: 20px; color: #6b7280; font-size: 12px; }
        table { width: 100%; border-collapse: collapse; margin-top: 10px; }
        th, td { text-align: left; padding: 8px; }
        th { background: #f3f4f6; color: #374151; font-weight: 600; }
        tr:nth-child(even) { background: #f9fafb; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>⚡ ZenithAgent Daily Report</h1>
            <p>Autonomous Stealth Engine - Activity Summary</p>
            <p>` + time.Now().Format("Monday, 02 January 2006 15:04:05 WIB") + `</p>
        </div>
        <div class="content">`)

	// Project Statistics
	if len(stats.Projects) > 0 {
		html.WriteString("<h2 style='color: #1f2937; margin-bottom: 20px;'>📊 Project Statistics</h2>")
		for name, project := range stats.Projects {
			statusClass := "status-stopped"
			statusText := "Stopped"
			if project.IsRunning {
				statusClass = "status-running"
				statusText = "Running"
			}

			html.WriteString(fmt.Sprintf(`
            <div class="project-card">
                <div class="project-name">%s <span class="status %s">%s</span></div>
                <p style="color: #6b7280; font-size: 14px; margin: 5px 0;">Last Run: %s</p>
                <div class="stats">
                    <div class="stat">
                        <div class="stat-label">Success</div>
                        <div class="stat-value success">%d</div>
                    </div>
                    <div class="stat">
                        <div class="stat-label">Failures</div>
                        <div class="stat-value failure">%d</div>
                    </div>
                    <div class="stat">
                        <div class="stat-label">Success Rate</div>
                        <div class="stat-value" style="color: #667eea;">%.1f%%</div>
                    </div>
                </div>`,
				name,
				statusClass,
				statusText,
				project.LastRun.Format("02 Jan 2006 15:04:05"),
				project.SuccessCount,
				project.FailedCount,
				calculateSuccessRate(project.SuccessCount, project.FailedCount)))

			if len(project.FailReasons) > 0 {
				html.WriteString("<div style='margin-top: 15px;'><strong style='color: #ef4444;'>Failure Reasons:</strong><ul style='margin: 5px 0; padding-left: 20px;'>")
				for _, reason := range project.FailReasons {
					html.WriteString(fmt.Sprintf("<li style='color: #6b7280;'>%s</li>", reason))
				}
				html.WriteString("</ul></div>")
			}

			html.WriteString("</div>")
		}
	} else {
		html.WriteString("<p style='color: #6b7280; font-style: italic;'>No projects have been executed yet.</p>")
	}

	// Recent Activity
	if len(stats.History) > 0 {
		html.WriteString(`<div class="history"><h2>📋 Recent Activity (Last 20)</h2>`)
		limit := 20
		if len(stats.History) < limit {
			limit = len(stats.History)
		}
		for i := 0; i < limit; i++ {
			entry := stats.History[i]
			logClass := "log-error"
			iconClass := "log-icon-error"
			icon := "✕"
			if entry.Type == "success" {
				logClass = "log-success"
				iconClass = "log-icon-success"
				icon = "✓"
			}

			html.WriteString(fmt.Sprintf(`
            <div class="log-entry %s">
                <div class="log-icon %s">%s</div>
                <div class="log-details">
                    <div class="log-message">%s</div>
                    <div class="log-meta">%s • %s</div>
                </div>
            </div>`,
				logClass,
				iconClass,
				icon,
				entry.Message,
				entry.Project,
				entry.Timestamp))
		}
		html.WriteString("</div>")
	}

	html.WriteString(`
        </div>
        <div class="footer">
            <p>This is an automated report from ZenithAgent</p>
            <p>Autonomous Stealth Engine © 2026</p>
        </div>
    </div>
</body>
</html>`)

	return html.String()
}

func calculateSuccessRate(success, failed int) float64 {
	total := success + failed
	if total == 0 {
		return 0.0
	}
	return (float64(success) / float64(total)) * 100
}

func (n *EmailNotifier) SendEmail(subject, body string) error {
	if n.User == "" || n.Password == "" || n.Recipient == "" {
		return fmt.Errorf("email configuration missing")
	}

	addr := fmt.Sprintf("%s:%s", n.Host, n.Port)
	auth := smtp.PlainAuth("", n.User, n.Password, n.Host)

	msg := []byte(fmt.Sprintf("To: %s\r\n"+
		"Subject: ZenithAgent: %s\r\n"+
		"\r\n"+
		"%s\r\n", n.Recipient, subject, body))

	return smtp.SendMail(addr, auth, n.User, []string{n.Recipient}, msg)
}

func (n *EmailNotifier) SendHTMLEmail(subject, htmlBody string) error {
	if n.User == "" || n.Password == "" || n.Recipient == "" {
		return fmt.Errorf("email configuration missing")
	}

	addr := fmt.Sprintf("%s:%s", n.Host, n.Port)
	auth := smtp.PlainAuth("", n.User, n.Password, n.Host)

	// MIME headers for HTML email
	headers := make(map[string]string)
	headers["From"] = n.User
	headers["To"] = n.Recipient
	headers["Subject"] = "ZenithAgent: " + subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=UTF-8"

	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + htmlBody

	return smtp.SendMail(addr, auth, n.User, []string{n.Recipient}, []byte(message))
}
