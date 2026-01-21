# Almai Bot - Usage Guide

## Overview
Almai Bot is an automated registration system integrated into ZenithAgent that processes CSV files and registers users to https://almai.id/referral/scalpinghack

## Features
- ✅ **Automatic CSV Detection**: Auto-detects encoding (UTF-8, Latin-1, CP1252) and delimiter (comma, semicolon, tab)
- ✅ **Smart Column Mapping**: Automatically identifies Name, Email, Phone, and IP columns
- ✅ **IP Duplicate Filtering**: Filters entries with duplicate IPs into separate `sus_data.json`
- ✅ **Session Management**: Auto-creates date-based sessions and resumes unfinished work
- ✅ **OTP Detection**: Treats OTP prompt as success indicator
- ✅ **Dashboard Integration**: Full integration with ZenithAgent monitoring system

## Folder Structure

```
almai-folder/
├── file-csv/          # Place your CSV files here
│   └── sample.csv     # Example CSV template
└── output/            # Auto-generated session folders
    └── 2026-01-21/    # Date-based session folder
        ├── data.json           # Clean entries (unique IPs)
        ├── sus_data.json       # Suspicious entries (duplicate IPs)
        ├── progress.json       # Execution progress
        └── screenshots/        # Error screenshots
```

## CSV Format

Your CSV file should contain the following columns (header names are auto-detected):

| Column      | Alternative Names                                  | Example                  |
|-------------|---------------------------------------------------|--------------------------|
| Name        | nama, full name, nama_lengkap                     | Andi Pratama             |
| Email       | e-mail, alamat email                              | andi@example.com         |
| Phone       | whatsapp, hp, phone_number, no_hp                 | 628123456789             |
| IP Address  | ip, ip_address, alamat ip                         | 192.168.1.100            |

### Sample CSV:
```csv
Nama Lengkap,Email,WhatsApp,IP Address
Andi Pratama,andi.pratama@example.com,628123456789,192.168.1.100
Budi Santoso,budi.santoso@example.com,0812-3456-890,192.168.1.101
```

**Note**: Phone numbers are automatically normalized to Indonesian international format (62xxx)

## Usage

### 1. Prepare CSV File
Place your CSV file in `almai-folder/file-csv/`

### 2. Run ZenithAgent
```bash
cd /home/linjinfeng/Documents/ZenithAgent
./zenith-agent
```

### 3. Select Almai from Menu
When prompted, select the Almai option from the project list

### 4. Monitor Progress
- The bot will automatically create a session folder based on current date
- Progress is saved to `progress.json` after each entry
- If interrupted, the bot will automatically resume from where it left off

## How It Works

1. **CSV Loading**: Bot finds the most recent CSV file in `file-csv/`
2. **IP Filtering**: Entries with duplicate IPs are moved to `sus_data.json`
3. **Processing**: Bot processes each clean entry:
   - Navigates to registration page
   - Detects and opens registration modal
   - Fills form (name, email, phone, password)
   - Submits and waits for OTP
   - OTP modal = SUCCESS ✅
4. **Auto-Resume**: On next run, bot skips already-successful entries

## Session Management

The bot uses **automatic session management** without user prompts:

- **New Session**: Creates `almai-folder/output/YYYY-MM-DD/` folder
- **Resume Session**: Automatically continues from `progress.json`
- **VPS-Ready**: No manual intervention required

## Statistics Integration

Almai integrates with ZenithAgent dashboard:
- Real-time status tracking
- Success/failure counters
- Error logging
- Multi-project support

## Troubleshooting

### No CSV Found
- Ensure CSV file is placed in `almai-folder/file-csv/`
- Check file has `.csv` extension

### Column Detection Failed
- Verify CSV has headers in first row
- Use common column names (Name, Email, Phone, IP)

### Registration Failures
- Check `screenshots/` folder for error images
- Review `progress.json` for error messages
- Verify internet connection and website availability

## Default Credentials

All registrations use default password:
- **Password**: `12345678`
- **Confirm Password**: `12345678`

*Note: Users should change passwords after registration*
