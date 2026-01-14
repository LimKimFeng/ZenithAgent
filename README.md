# ZenithAgent 2.0
> **High-Performance Autonomous Web Automation Engine with Modern Dashboard**

ZenithAgent is an enterprise-grade automation bot engineered in **Go (Golang)** with a professional **React.js dashboard**, designed for stability, stealth, and 24/7 reliability. Unlike standard scraping scripts, ZenithAgent focuses on **human-like behavior simulation**, **security**, and **resource efficiency**, making it suitable for long-running deployments on VPS environments.

## ✨ What's New in 2.0

- 🎨 **Professional React Dashboard** - Real-time monitoring with glassmorphism UI
- 🔒 **File Locking System** - Prevents concurrent execution conflicts
- 📧 **Enhanced Email Reports** - Beautiful HTML emails with complete statistics
- 🚀 **VPS Ready** - Production-ready deployment configuration

## 🚀 Quick Start

### First Time Setup (After Cloning)

```bash
# Clone the repository
git clone <your-repo-url> ZenithAgent
cd ZenithAgent

# Run setup script (installs dependencies and builds everything)
./setup.sh
```

### Running the Application

```bash
# Quick start (recommended for daily use)
./start.sh

# Access dashboard at http://localhost:8080
```

### When to Use Each Script

| Script | When to Use | What It Does |
|--------|-------------|--------------|
| `./setup.sh` | First time after cloning | Checks dependencies, installs packages, builds everything |
| `./start.sh` | Daily use / Quick start | Just runs the app (no rebuild) |
| `./run.sh` | After code changes | Rebuilds React + Go, then runs |

**Recommended workflow:**
1. First time: `./setup.sh`
2. Daily use: `./start.sh`
3. After updates: `./run.sh`

## 🎯 Key Features

### 1. Modern Dashboard
- **Real-time Monitoring**: Live stats updated every 3 seconds
- **Professional UI**: Glassmorphism design with smooth animations
- **Multi-Project Support**: Track multiple automation projects
- **Activity Log**: Color-coded history with timestamps
- **Responsive Design**: Works on desktop, tablet, and mobile

### 2. Stealth-First Architecture
- **Playwright Engine**: Custom stealth configurations
- **Fingerprint Masking**: Strips automation flags and rotates User-Agents
- **Dynamic IP Rotation**: Tor Network integration with 10-minute rotation
- **Human Simulation**: Non-linear typing with 50ms-200ms variance

### 3. Production-Ready Reliability
- **File Locking**: Prevents duplicate process execution
- **Resource Management**: Strict memory control with forced closures
- **Bandwidth Optimization**: ~90% reduction via resource blocking
- **Email Reporting**: Daily HTML reports with complete statistics

## 🛠 Tech Stack

**Backend:**
- Go 1.22+
- Playwright for browser automation
- Tor (SOCKS5 + Control Port Auth)
- SMTP for email notifications

**Frontend:**
- React 18 + TypeScript
- Vite (build tool)
- TailwindCSS 3 (styling)
- Native Fetch API

**Deployment:**
- Nginx (reverse proxy)
- Let's Encrypt (SSL/TLS)
- Systemd (process management)

## 📁 Project Structure

```
ZenithAgent/
├── cmd/agent/main.go           # Application entry point
├── internal/
│   ├── engine/                 # Browser automation
│   ├── manager/                # Process & lock management
│   ├── monitor/                # Dashboard server
│   ├── network/                # Tor integration
│   ├── notify/                 # Email notifications
│   ├── stats/                  # Statistics tracking
│   └── tasks/                  # Automation tasks
├── dashboard/                  # React frontend
│   ├── src/
│   │   ├── App.tsx            # Main component
│   │   ├── components/        # UI components
│   │   └── index.css          # Styles
│   └── dist/                  # Production build
├── setup.sh                    # First-time setup
├── start.sh                    # Quick start script
├── run.sh                      # Build and run
├── nginx.conf.example          # Nginx configuration
└── DEPLOYMENT.md               # Deployment guide
```

## 🔒 Security Features

- **Basic Authentication**: Bcrypt-hashed passwords
- **File Locking**: Exclusive process locks with flock
- **Tor Authentication**: Hashed control password support
- **SMTP Security**: TLS encryption (port 587)
- **No Credential Storage**: Runtime-only configuration

## 📊 Script Usage Guide

### setup.sh - First Time Setup
Run this **once** after cloning the repository:
```bash
./setup.sh
```
- Checks all dependencies (Node.js, Go, Tor)
- Installs npm packages
- Builds React dashboard
- Compiles Go binary
- Sets up everything automatically

### start.sh - Quick Start (Recommended)
Use this for **daily operation**:
```bash
./start.sh
```
- Starts the application immediately
- No rebuild (fast startup)
- Checks for lock files
- Best for regular use

### run.sh - Full Rebuild
Use this **after making code changes**:
```bash
./run.sh
```
- Rebuilds React dashboard
- Recompiles Go binary
- Then starts the application
- Use when you update the code

## 📧 Email Reports

Daily automated reports include:
- ✅ Project statistics (success/failure counts)
- 📊 Success rate percentages
- 📋 Last 20 activity entries
- ⚠️ Failure reasons (if any)
- 🎨 Professional HTML formatting

## 🤝 Contributing

This is a demonstration project showcasing:
- High-availability software engineering
- Stealth automation techniques
- Modern full-stack development
- Production-ready deployment practices

## 📄 License

Private project - All rights reserved

## 🚀 VPS Deployment

For production deployment instructions, see [DEPLOYMENT.md](DEPLOYMENT.md)

---

*Developed as a demonstration of enterprise-grade automation and modern web development.*
