# ZenithAgent
> **High-Performance Autonomous Web Automation Engine**

ZenithAgent is an enterprise-grade automation bot engineered in **Go (Golang)**, designed for stability, stealth, and 24/7 reliability. Unlike standard scraping scripts, ZenithAgent focuses on **human-like behavior simulation**, **security**, and **resource efficiency**, making it suitable for long-running deployments on low-cost VPS environments.

## 🚀 Key Technical Highlights

### 1. Stealth-First Architecture
- **Engine**: Built on top of `playwright-go` with custom stealth configurations.
- **Fingerprint Masking**: Automatically strips automation flags (`navigator.webdriver`) and rotates User-Agents to mimic legitimate traffic.
- **Dynamic IP Rotation (Tor)**: Integrated with Tor Network (SOCKS5 Proxy). Features a **Hashed Authenticated Rotator** that requests a new identity (`SIGNAL NEWNYM`) every 10 minutes (configurable) to prevent IP banning.

### 2. Smart Behavioral Logic
- **Human Simulation**: Features a **Non-Linear Typing Engine** that simulates human keystroke latency (50ms-200ms variance), avoiding bot detection heuristics based on typing speed.
- **Context-Aware Data**: Implements "Smart Reason Generators" that construct unique, relevant sentences for form filling (e.g., specific trading jargon for Scalping projects, property investment reasons for Real Estate).
- **Resilience Strategy**: Utilizes DOM-state awareness (`WaitUntilStateLoad`) rather than fragile network-idle triggers.

### 3. Production-Ready Reliability
- **Singleton Process Manager**: Uses a `state.json` lock-file mechanism to prevent duplicate process execution and race conditions during IP rotation.
- **Resource Management**: Strict memory policing with forced context closures after every session.
- **Bandwidth Optimization**: Integrated resource blocker (Images/CSS/Fonts) reducing bandwidth by ~90%.
- **Observability**: Built-in **24-Hour Reporting Cycle** that persists statistics locally and dispatches daily SMTP summaries.

## 🛠 Tech Stack
- **Language**: Go 1.22+
- **Automation**: Playwright
- **Network**: Tor (SOCKS5 + Control Port Auth)
- **Concurrency**: Goroutines for non-blocking reporting and background rotation.
- **Deployment**: Single-binary compilation (`ldflags -s -w`) for minimal footprint.

## 🔒 Security & Privacy
- **Authenticated Control**: Supports `HashedControlPassword` for Tor Control Port to secure the IP rotation mechanism.
- **Ephemeral Credentials**: SMTP, Tor Passwords, and configs are injected via runtime CLI memory only; never stored on disk.
- **Data Minimization**: Logs store only execution metrics, stripping all PII generated during runtime.

---
*Developed as a demonstration of high-availability software engineering and stealth automation techniques.*
