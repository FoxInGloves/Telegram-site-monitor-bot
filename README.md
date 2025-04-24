# 🚨 Telegram Site Monitor Bot

**Telegram Site Monitor Bot** is a lightweight and convenient solution on **Go** that allows you to **monitor website
availability** and **promptly receive notifications in Telegram** in case of failures.

---

## 🔧 What does this bot do?

- 📡 **Checks website availability** via HTTP/HTTPS.
- 📬 **Sends notifications to Telegram** if the site becomes unavailable or works again.
- 📈 **Logs all checks**: both successful and unsuccessful to the file [year-month-date].log.
- 📋 Provides a summary of the `/status` command in telegram with the current state of sites.

---

## 📦 Key features

- ✅ Support for multiple sites.
- ⏱ Configurable intervals and request timeouts.
- 🔁 Automatic startup via systemd.
- ⚙️ Configuration file can be changed without recompilation
- 🔐 Security: configuration is stored locally, without third-party services.

---

## 📂 Installation and configuration

1. **Install Go** (v1.24.2 or higher).
2. **Configure the server** (Linux, Ubuntu/Debian recommended).
3. **Create a `config.toml` configuration file:**

```toml
[telegram]
bot_token = "YOUR_TELEGRAM_BOT_TOKEN"
chat_id = 123456789

[sites]
urls = [
  "https://example.com",
  "https://google.com",
]

[settings]
check_interval = 300 # Check interval (in seconds)
timeout = 10 # Request timeout (in seconds)
```

4. **Build and run the bot:**

- When building, the config must be located in the root directory of the project
```bash
go build -o site-monitor-bot
./site-monitor-bot
```

- When running via the `-path` flag or `-p` for short, pass the path to config.
  If you do not use the flag, the program will use the config in the root directory
```bash
go run main.go -path /path/to/your/config.toml
```

---

## 📫 Feedback

If you have any ideas, suggestions or bugs, feel free to open an issue or make a pull request!