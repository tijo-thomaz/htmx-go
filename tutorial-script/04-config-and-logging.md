# Scene 4: Config + Logging

> ⏱️ **Timestamp**: 15:00 — 20:00
> 🎯 **Goal**: Create config.go (load all env vars including SessionEncKey) and logger.go (JSON for prod, Text for dev)

---

## 15:00 — Config: Why We Need It

📱 **Narration**:
> "ആദ്യത്തെ real code! Configuration loading. .env file-ൽ ഉള്ള values read ചെയ്ത് ഒരു Go struct-ൽ store ചെയ്യും."

🎯 **Analogy**:
> "Config file ഒരു recipe card പോലെ. App start ചെയ്യുമ്പോൾ recipe card read ചെയ്ത് — 'port 8080, database ഇവിടെ, secret key ഇത്' — എല്ലാം ഒരു place-ൽ. Code-ൽ values hardcode ചെയ്യുന്നില്ല."

---

## 15:30 — Create config.go

📱 **Narration**:
> "`internal/config/config.go` create ചെയ്യാം."

🎥 **Camera**: VS Code. Create new file `internal/config/config.go`.

⌨️ **Type this** (type slowly, explain as you go):

```go
package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all application configuration
type Config struct {
	Port          string
	Env           string
	LogLevel      string
	DatabasePath  string
	SessionSecret string
	SessionEncKey string
	RateLimit     int
}
```

🧠 **Pause and explain the struct**:

📱 **Narration**:
> "Config struct — നമ്മുടെ app-ന്റെ എല്ലാ settings ഒരിടത്ത്."

> "`Port` — server listen ചെയ്യുന്ന port. String ആണ്, int അല്ല — HTTP address-ന് `:8080` format-ൽ string ആണ് വേണ്ടത്."

> "`Env` — development or production. Different behavior switch ചെയ്യാൻ."

> "`LogLevel` — DEBUG, INFO, WARN, ERROR. Development-ൽ DEBUG — everything log ചെയ്യും. Production-ൽ INFO — important things മാത്രം."

> "`DatabasePath` — SQLite file-ന്റെ path."

> "`SessionSecret` — session cookies sign ചെയ്യാൻ. Tampering detect ചെയ്യാൻ — remember wax seal analogy?"

> "`SessionEncKey` — session data encrypt ചെയ്യാൻ. Content hide ചെയ്യാൻ — locked box! Exactly 32 bytes."

> "`RateLimit` — int ആണ്. Per-IP requests per second."

---

⌨️ **Continue typing — Load function**:

```go
// Load reads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file (ignore error if not exists)
	_ = godotenv.Load()

	return &Config{
		Port:          getEnv("PORT", "8080"),
		Env:           getEnv("ENV", "development"),
		LogLevel:      getEnv("LOG_LEVEL", "INFO"),
		DatabasePath:  getEnv("DATABASE_PATH", "./data/linkbio.db"),
		SessionSecret: getEnv("SESSION_SECRET", "change-me-in-production"),
		SessionEncKey: getEnv("SESSION_ENCRYPTION_KEY", ""),
		RateLimit:     getEnvInt("RATE_LIMIT", 10),
	}, nil
}
```

🧠 **Explain**:

📱 **Narration**:
> "`Load()` function — pointer return ചെയ്യുന്നു `*Config`. Pointer ആയതിന്റെ reason — Config struct copy ആക്കാതെ reference pass ചെയ്യാൻ. Memory efficient."

> "`godotenv.Load()` — `.env` file read ചെയ്ത് environment variables set ചെയ്യുന്നു. Error ignore ചെയ്യുന്നു — production-ൽ `.env` file ഉണ്ടാകില്ല, real environment variables ആയിരിക്കും."

> "`_ = godotenv.Load()` — underscore use ചെയ്യുന്നത് 'error explicitly ignore ചെയ്യുന്നു' എന്ന് Go-യോട് പറയാൻ. `.env` file optional ആണ്."

> "`getEnv()` — environment variable read ചെയ്യും. ഇല്ലെങ്കിൽ default value return ചെയ്യും. ഇത് നമ്മൾ next write ചെയ്യും."

> "`SessionEncKey` default empty string ആണ്. Empty ആണെങ്കിൽ gorilla/sessions encryption skip ചെയ്യും — development-ൽ OK, production-ൽ set ചെയ്യണം."

> "`getEnvInt()` — string to int convert ചെയ്ത് return ചെയ്യുന്ന helper. RateLimit-ന് int value വേണം."

---

⌨️ **Continue typing — helper methods**:

```go
// IsDevelopment returns true if running in development mode
func (c *Config) IsDevelopment() bool {
	return c.Env == "development"
}

// IsProduction returns true if running in production mode
func (c *Config) IsProduction() bool {
	return c.Env == "production"
}
```

🧠 **Explain**:

📱 **Narration**:
> "Method receiver — `(c *Config)`. ഇത് Config struct-ന്റെ method ആണ്. `c` ഒരു Config pointer — ആ Config-ന്റെ `Env` field access ചെയ്യാം."

> "IsDevelopment, IsProduction — boolean return ചെയ്യുന്ന convenience methods. Code-ൽ `if cfg.IsDevelopment()` എന്ന് check ചെയ്യാം — readable."

---

⌨️ **Continue typing — helper functions**:

```go
// getEnv retrieves env variable or returns fallback
func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

// getEnvInt retrieves env variable as int or returns fallback
func getEnvInt(key string, fallback int) int {
	if val := os.Getenv(key); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return fallback
}
```

🧠 **Explain**:

📱 **Narration**:
> "`getEnv` — lowercase function. Go-ൽ lowercase means unexported — ഈ package-ന് പുറത്ത് access ചെയ്യാൻ പറ്റില്ല. Internal helper function."

> "`os.Getenv(key)` — operating system-ന്റെ environment variable read ചെയ്യുന്നു. godotenv ഇവ set ചെയ്തിട്ടുണ്ടാകും."

> "`if val := ...; val != \"\"` — Go-യുടെ short variable declaration with if. Variable declare ചെയ്ത് same line-ൽ check ചെയ്യാം. Concise, idiomatic Go."

> "`getEnvInt` — similar, but string to int convert ചെയ്യുന്നു. `strconv.Atoi` — ASCII to Integer. Error ഉണ്ടെങ്കിൽ — invalid number — fallback return ചെയ്യും."

---

## 17:00 — Full config.go File Review

🎥 **Camera**: Zoom out. Show the complete file.

📱 **Narration**:
> "Full file ഒന്ന് review ചെയ്യാം."

🎥 **Show complete file**:
```go
package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all application configuration
type Config struct {
	Port          string
	Env           string
	LogLevel      string
	DatabasePath  string
	SessionSecret string
	SessionEncKey string
	RateLimit     int
}

// Load reads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file (ignore error if not exists)
	_ = godotenv.Load()

	return &Config{
		Port:          getEnv("PORT", "8080"),
		Env:           getEnv("ENV", "development"),
		LogLevel:      getEnv("LOG_LEVEL", "INFO"),
		DatabasePath:  getEnv("DATABASE_PATH", "./data/linkbio.db"),
		SessionSecret: getEnv("SESSION_SECRET", "change-me-in-production"),
		SessionEncKey: getEnv("SESSION_ENCRYPTION_KEY", ""),
		RateLimit:     getEnvInt("RATE_LIMIT", 10),
	}, nil
}

// IsDevelopment returns true if running in development mode
func (c *Config) IsDevelopment() bool {
	return c.Env == "development"
}

// IsProduction returns true if running in production mode
func (c *Config) IsProduction() bool {
	return c.Env == "production"
}

// getEnv retrieves env variable or returns fallback
func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

// getEnvInt retrieves env variable as int or returns fallback
func getEnvInt(key string, fallback int) int {
	if val := os.Getenv(key); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return fallback
}
```

📱 **Narration**:
> "63 lines of code. Clean, simple, readable. ഓരോ function-നും ഒരു job — single responsibility. ഇത് production-ready code ആണ്."

---

## 17:30 — Logging: Why Not fmt.Println?

📱 **Narration**:
> "ഇനി logging. ഇവിടെ ഒരു question — fmt.Println use ചെയ്താൽ പോരെ? Why structured logging?"

🎯 **Analogy — Diary vs Spreadsheet**:

📱 **Narration**:
> "fmt.Println ഒരു **diary** പോലെ. Random notes — 'user logged in', 'error happened', 'server started'. ശരി, but 10,000 lines ഉണ്ടെങ്കിൽ ഒരു specific error find ചെയ്യാൻ impossible!"

> "slog ഒരു **spreadsheet** പോലെ. ഓരോ log entry-യ്ക്കും columns ഉണ്ട് — timestamp, level, message, key-value pairs. Filter ചെയ്യാം, search ചെയ്യാം, analyze ചെയ്യാം."

🎥 **Show comparison**:
```
fmt.Println:
  user logged in
  error: database connection failed
  server started on port 8080

slog (JSON format):
  {"time":"2025-01-15T10:30:00Z","level":"INFO","msg":"user logged in","user_id":42,"username":"tijo"}
  {"time":"2025-01-15T10:30:01Z","level":"ERROR","msg":"database connection failed","error":"timeout","retry":3}
  {"time":"2025-01-15T10:30:02Z","level":"INFO","msg":"server started","port":"8080","env":"production"}
```

📱 **Narration**:
> "Difference കാണുന്നുണ്ടോ? JSON format — Datadog, Splunk, Grafana — log aggregation tools ഇത് parse ചെയ്ത് dashboards create ചെയ്യും. fmt.Println-ൽ ഇത് impossible."

---

## 18:00 — Create logger.go

📱 **Narration**:
> "`internal/pkg/logger/logger.go` create ചെയ്യാം."

🎥 **Camera**: VS Code. Create new file.

⌨️ **Type this**:

```go
package logger

import (
	"log/slog"
	"os"
	"strings"
)

// New creates a new structured logger with the specified level
func New(level string) *slog.Logger {
	var logLevel slog.Level

	switch strings.ToUpper(level) {
	case "DEBUG":
		logLevel = slog.LevelDebug
	case "INFO":
		logLevel = slog.LevelInfo
	case "WARN":
		logLevel = slog.LevelWarn
	case "ERROR":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     logLevel,
		AddSource: false,
	})

	return slog.New(handler)
}

// NewDevelopment creates a logger optimized for development (text format)
func NewDevelopment() *slog.Logger {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: false,
	})

	return slog.New(handler)
}
```

---

## 18:30 — Explain Logger Code

🧠 **Explain each part**:

📱 **Narration — package and imports**:
> "`package logger` — reusable package. App-ന്റെ anywhere import ചെയ്ത് use ചെയ്യാം."

> "`log/slog` — Go 1.21-ൽ വന്ന built-in structured logging package. Third-party libraries വേണ്ട — standard library-ൽ ഉണ്ട്!"

> "`strings` — `strings.ToUpper()` use ചെയ്യാൻ. Level string uppercase ആക്കാൻ — `debug`, `Debug`, `DEBUG` — എല്ലാം work ആകും."

📱 **Narration — New function**:
> "`New(level string)` — production logger create ചെയ്യുന്നു. Level parameter — ഏത് level-ന് മുകളിൽ ഉള്ള logs മാത്രം output ചെയ്യും."

> "Switch statement — level string match ചെയ്ത് slog.Level set ചെയ്യുന്നു. Default INFO — unknown level വന്നാൽ INFO use ചെയ്യും. Safe fallback."

> "`slog.NewJSONHandler` — JSON format output. Production-ൽ ഇത്. ⚠️ Every log line ഒരു valid JSON object — machine-readable."

> "`AddSource: false` — file name, line number log-ൽ add ചെയ്യുന്നില്ല. Performance-ന് better. Debugging-ന് വേണമെങ്കിൽ true ആക്കാം."

📱 **Narration — NewDevelopment function**:
> "`NewDevelopment()` — developer-friendly logger. Text format — human-readable."

> "`slog.NewTextHandler` — key=value format. Terminal-ൽ read ചെയ്യാൻ easy."

> "Level always DEBUG — development-ൽ everything log ചെയ്യണം."

---

## 19:00 — How Logging Will Be Used

📱 **Narration**:
> "ഇത് എങ്ങനെ use ചെയ്യും? Quick preview."

🎥 **Show (don't type, just show on screen)**:
```go
// In main.go (we'll write this later):
cfg, _ := config.Load()

var log *slog.Logger
if cfg.IsDevelopment() {
    log = logger.NewDevelopment()
} else {
    log = logger.New(cfg.LogLevel)
}

// Usage anywhere in the app:
log.Info("server started", "port", cfg.Port, "env", cfg.Env)
log.Error("database failed", "error", err)
log.Debug("request received", "method", "GET", "path", "/dashboard")
```

📱 **Narration**:
> "Development-ൽ `NewDevelopment()` — text format, DEBUG level. Terminal-ൽ clear ആയി read ചെയ്യാം."

> "Production-ൽ `New(cfg.LogLevel)` — JSON format, configurable level. Log aggregators parse ചെയ്യാം."

🎥 **Show output comparison**:
```
Development (Text):
  time=2025-01-15T10:30:00.000+05:30 level=INFO msg="server started" port=8080 env=development

Production (JSON):
  {"time":"2025-01-15T05:00:00Z","level":"INFO","msg":"server started","port":"8080","env":"production"}
```

📱 **Narration**:
> "Same code, different output. Config control ചെയ്യുന്നു. ഇത് separation of concerns."

---

## 19:30 — Quick Recap

📱 **Narration**:
> "ഇവിടെ വരെ നമ്മൾ ചെയ്തത്:"

> "ഒന്ന് — Config struct. .env file read ചെയ്ത് type-safe struct-ൽ store ചെയ്യുന്നു. SessionSecret signing-ന്, SessionEncKey encryption-ന്."

> "രണ്ട് — Helper functions. getEnv default values-ഓടെ, getEnvInt integer conversion-ഓടെ."

> "മൂന്ന് — Logger. Production-ൽ JSON, development-ൽ Text. slog — Go standard library."

> "ഇത് foundation ആണ്. ഇനി database code എഴുതാം!"

🔊 **Transition sound**

🎥 **Cut to**: Database + Models (Scene 5)

---

## 📝 Editing Notes

- **Typing speed**: This is a code-heavy section. Type at a readable pace. Don't rush.
- **Split screen**: Consider showing the `.env` file on one side and `config.go` on the other when explaining which env var maps to which field.
- **Diary vs Spreadsheet**: This analogy resonates well. Consider a quick graphic — a messy diary page vs a clean spreadsheet.
- **Output comparison**: The development vs production log output is a strong visual. Show both in split terminal or overlay.
- **File save**: After finishing each file, press Ctrl+S visibly. Show the file in the VS Code explorer to confirm location.
- **Common mistake**: Viewers might forget the `strings` import for `strings.ToUpper()`. The Go extension will auto-import, but mention it.
