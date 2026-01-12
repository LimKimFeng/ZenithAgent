# ZenithAgent
> **High-Performance Autonomous Web Automation Engine**

ZenithAgent is an enterprise-grade automation bot engineered in **Go (Golang)**, designed for stability, stealth, and 24/7 reliability. Unlike standard scraping scripts, ZenithAgent focuses on **human-like behavior simulation** and **resource efficiency**, making it suitable for long-running deployments on low-cost VPS environments.

## 🚀 Key Technical Highlights

### 1. Stealth-First Architecture
- **Engine**: Built on top of `playwright-go` with custom stealth configurations.
- **Fingerprint Masking**: Automatically strips automation flags (`navigator.webdriver`) and rotates User-Agents to mimic legitimate traffic.
- **Human Simulation**: Features a **Non-Linear Typing Engine** that simulates human keystroke latency (50ms-200ms variance), avoiding bot detection heuristics based on typing speed.

### 2. Smart Behavioral Logic
- **Dynamic Content Generation**: Implements a randomized "Smart Reason Generator" that constructs unique, context-aware sentence structures. This ensures that form inputs never look templates or repetitive.
- **Resilience Strategy**: Utilizes DOM-state awareness (`WaitUntilStateLoad`) rather than fragile network-idle triggers, ensuring success even on heavy sites loaded with tracking scripts.

### 3. production-Ready Reliability
- **Resource Management**: Strict memory policing with forced context closures after every session to prevent RAM leaks.
- **Bandwidth Optimization**: Integrated resource blocker (Images/CSS/Fonts) to reduce bandwidth usage by up to 90%.
- **Observability**: Built-in **24-Hour Reporting Cycle** that persists statistics locally and dispatches daily SMTP summaries (Success/Failure rates & granular error tracing).

## 🛠 Tech Stack
- **Language**: Go 1.22+
- **Automation**: Playwright
- **Concurrency**: Goroutines for non-blocking reporting and task execution.
- **Deployment**: Single-binary compilation with symbol stripping (`ldflags -s -w`) for minimal footprint.

## 🔒 Security & Privacy
- **Ephemeral Credentials**: SMTP and sensitive configs are injected via runtime CLI memory only; never stored on disk.
- **Data Minimization**: Logs store only execution metrics, stripping all PII (Personal Identifiable Information) generated during runtime.

---
*Developed as a demonstration of high-availability software engineering and stealth automation techniques.*
