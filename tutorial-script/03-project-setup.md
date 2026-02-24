# Scene 3: Project Setup

> ⏱️ **Timestamp**: 8:00 — 15:00
> 🎯 **Goal**: Create project, folder structure, install dependencies, create .env and .gitignore

---

## 8:00 — Create Project

📱 **Narration**:
> "ആദ്യം നമ്മൾ project folder create ചെയ്യാം."

⌨️ **Type this**:
```bash
mkdir linkbio
cd linkbio
go mod init linkbio
```

🧠 **Explain**:
> "go mod init — ഇത് നമ്മുടെ project-ന്റെ ID card പോലെ. Go modules — dependency management system. npm init പോലെ, പക്ഷേ Go-യ്ക്ക്."

---

## 8:30 — Folder Structure

📱 **Narration**:
> "Industry-standard folder structure create ചെയ്യാം. ഇത് random folders അല്ല — Google, Uber, Stripe — large companies use ചെയ്യുന്ന pattern ആണ്."

⌨️ **Type this** (one by one, explain as you go):
```bash
mkdir -p cmd/server
mkdir -p internal/config
mkdir -p internal/server
mkdir -p internal/router
mkdir -p internal/handler
mkdir -p internal/middleware
mkdir -p internal/model
mkdir -p internal/repository
mkdir -p internal/service
mkdir -p internal/pkg/logger
mkdir -p internal/pkg/response
mkdir -p internal/pkg/templates
mkdir -p internal/testutil
mkdir -p web/templates/layouts
mkdir -p web/templates/pages
mkdir -p web/templates/partials
mkdir -p web/static/css
mkdir -p web/static/js
mkdir -p data
mkdir -p migrations
```

🎥 **Camera**: After typing, show the folder tree in VS Code sidebar.

📱 **Narration — explain each folder**:

> "`cmd/server/` — Entry point. നമ്മുടെ `main.go` ഇവിടെ live ആകും. ഒരു project-ൽ multiple binaries ഉണ്ടാകാം — API server, CLI tool, worker. ഓരോന്നും `cmd/` under ഒരു subfolder."

> "`internal/` — ഇത് Go-യുടെ special folder ആണ്. ഇതിൽ ഉള്ള code മറ്റ് projects import ചെയ്യാൻ പറ്റില്ല. Go compiler enforce ചെയ്യും!"

⚠️ **Important — Why `internal/` is special**:

📱 **Narration**:
> "ഇത് ഒരു important Go rule ആണ്. `internal/` folder-ൽ ഉള്ള code — config, handler, model — ഇവ നമ്മുടെ project-ന് മാത്രം available ആണ്. മറ്റൊരു project `import linkbio/internal/config` ചെയ്താൽ Go compiler error throw ചെയ്യും."

🎯 **Analogy**:
> "internal/ ഒരു company-ന്റെ internal documents പോലെ. Employees access ചെയ്യാം, പക്ഷേ outside-ൽ ഉള്ള ആൾക്ക് access ഇല്ല. Go compiler ആ security guard ആണ്."

📱 **Continue explaining folders**:

> "`internal/config/` — App configuration. .env file reading."

> "`internal/server/` — HTTP server setup, graceful shutdown."

> "`internal/router/` — URL routing. ഏത് URL ഏത് handler-ലേക്ക്."

> "`internal/handler/` — Request handling logic. ഓരോ feature-നും ഒരു handler file."

> "`internal/middleware/` — Auth check, logging, rate limiting — requests filter ചെയ്യുന്ന code."

> "`internal/model/` — Data structures. User, Link, Analytics — Go structs."

> "`internal/repository/` — Database queries. SQL ഇവിടെ മാത്രം."

> "`internal/service/` — Business logic. Handler-ഉം repository-യും connect ചെയ്യുന്ന layer."

> "`internal/pkg/` — Shared utilities. Logger, response helpers, template engine."

> "`internal/testutil/` — Test helpers. Test database setup, test logger."

> "`web/templates/` — HTML templates. Layouts, pages, partials — reusable pieces."

> "`web/static/` — CSS, JavaScript files. Tailwind, HTMX, Alpine.js."

> "`data/` — SQLite database file. Git ignore ചെയ്യും."

> "`migrations/` — Database migration files. Future use."

🎯 **Analogy**:
> "ഈ folder structure ഒരു well-organized kitchen പോലെ. Spices ഒരിടത്ത്, utensils മറ്റൊരിടത്ത്, ingredients വേറൊരിടത്ത്. ആരെങ്കിലും kitchen-ൽ വന്നാൽ എല്ലാം find ചെയ്യാൻ easy. Code-ഉം same — organized ആയാൽ 6 months കഴിഞ്ഞ് വന്നാലും understand ചെയ്യാം."

---

## 10:00 — Install Dependencies

📱 **Narration**:
> "നമുക്ക് ആവശ്യമായ packages install ചെയ്യാം. ഓരോന്നും explain ചെയ്യാം."

⌨️ **Type and explain each**:

### 1. Chi Router
```bash
go get github.com/go-chi/chi/v5
```

🧠 **Explain**:
> "chi — lightweight HTTP router. Express.js പോലെ, പക്ഷേ Go-യ്ക്ക്. Standard library-യുടെ net/http compatible ആണ്. Middleware chaining, URL parameters — എല്ലാം support ചെയ്യും."

### 2. Gorilla Sessions
```bash
go get github.com/gorilla/sessions
```

🧠 **Explain**:
> "gorilla/sessions — login sessions manage ചെയ്യാൻ. User login ചെയ്യുമ്പോൾ ഒരു session cookie create ചെയ്യും. Next request-ൽ ആ cookie check ചെയ്ത് user-നെ identify ചെയ്യും."

### 3. SQLite Driver
```bash
go get modernc.org/sqlite
```

🧠 **Explain**:
> "modernc/sqlite — Pure Go SQLite driver. CGo വേണ്ട, C compiler install ചെയ്യണ്ട. Cross-platform — Windows, Mac, Linux — എല്ലായിടത്തും work ചെയ്യും."

⚠️ **Note**: 
> "ആദ്യ download time കൂടുതൽ ആകാം. Pure Go-ൽ SQLite compile ചെയ്യുന്നു. Patience!"

### 4. Crypto (bcrypt)
```bash
go get golang.org/x/crypto
```

🧠 **Explain**:
> "x/crypto — Go team maintain ചെയ്യുന്ന cryptography package. നമ്മൾ bcrypt use ചെയ്യും — password hashing-ന്. Plain text passwords database-ൽ store ചെയ്യരുത് — ⚠️ security disaster!"

### 5. Rate Limiter
```bash
go get golang.org/x/time
```

🧠 **Explain**:
> "x/time — rate limiting-ന്. ഒരു IP address-ൽ നിന്ന് per minute 10 requests only. Brute-force attacks prevent ചെയ്യാൻ."

### 6. Dotenv
```bash
go get github.com/joho/godotenv
```

🧠 **Explain**:
> "godotenv — `.env` file read ചെയ്ത് environment variables ആയി load ചെയ്യുന്നു. Development-ൽ convenient. Production-ൽ real environment variables use ചെയ്യും."

---

## 12:00 — Create .env File

📱 **Narration**:
> "Configuration file create ചെയ്യാം. Secrets-ഉം settings-ഉം ഇവിടെ store ചെയ്യും. ഈ file git-ൽ commit ചെയ്യരുത്!"

⌨️ **Create `.env`**:
```env
# Server Configuration
PORT=8080
ENV=development

# Logging
LOG_LEVEL=DEBUG

# Database
DATABASE_PATH=./data/linkbio.db

# Session
SESSION_SECRET=your-super-secret-key-change-in-production
SESSION_ENCRYPTION_KEY=must-be-exactly-32-bytes-long!!

# Rate Limiting
RATE_LIMIT=10
```

🧠 **Explain each variable**:

> "`PORT` — server ഏത് port-ൽ listen ചെയ്യണം. 8080 development standard."

> "`ENV` — development or production. Logging format, error details — ഇതനുസരിച്ച് change ആകും."

> "`LOG_LEVEL` — DEBUG, INFO, WARN, ERROR. Development-ൽ DEBUG — everything log ചെയ്യും."

> "`DATABASE_PATH` — SQLite file location."

> "`SESSION_SECRET` — ⚠️ ഇത് very important. Session cookies sign ചെയ്യാൻ use ചെയ്യുന്ന key."

> "`SESSION_ENCRYPTION_KEY` — ⚠️ ഇത് session data encrypt ചെയ്യാൻ. Exactly 32 bytes ആയിരിക്കണം."

> "`RATE_LIMIT` — per-IP requests per second."

---

## 12:30 — SESSION_SECRET vs SESSION_ENCRYPTION_KEY

⚠️ **This is the most important security concept in this section!**

🎥 **Camera**: Slow down here. Use a visual or draw on screen.

📱 **Narration**:
> "ഇവിടെ ഒരു common confusion ഉണ്ട്. SESSION_SECRET-ഉം SESSION_ENCRYPTION_KEY-ഉം — രണ്ടും different ആണ്. ഇത് understand ചെയ്യണം."

🎯 **Analogy — Seal vs Locked Box**:

📱 **Narration**:
> "ഒരു letter അയക്കുന്നു എന്ന് imagine ചെയ്യൂ."

> "SESSION_SECRET — ഇത് ഒരു **wax seal** പോലെ. Letter-ന്റെ outside-ൽ ഒരു seal — ആരെങ്കിലും open ചെയ്തോ എന്ന് detect ചെയ്യാം. ആരെങ്കിലും letter modify ചെയ്താൽ seal break ആകും. **Tampering detect ചെയ്യുന്നു, content hide ചെയ്യുന്നില്ല**. Letter read ചെയ്യാൻ ആർക്കും കഴിയും!"

> "SESSION_ENCRYPTION_KEY — ഇത് ഒരു **locked box** പോലെ. Letter box-ൽ ഇട്ടു lock ചെയ്തു. Key ഉള്ള ആൾക്ക് മാത്രം open ചെയ്യാൻ കഴിയും. **Content hide ചെയ്യുന്നു**."

> "രണ്ടും combine ചെയ്യുമ്പോൾ — letter ഒരു locked box-ൽ ഇട്ടു, box-ന് ഒരു seal ഇട്ടു. ആർക്കും read ചെയ്യാൻ പറ്റില്ല, modify ചെയ്യാനും പറ്റില്ല. **Double protection!**"

🎥 **Show diagram**:
```
SESSION_SECRET (Signing — HMAC):
┌──────────────────┐
│ session data     │ + secret key → signature
│ user_id=42       │
│ username=tijo    │   Anyone can READ the data
└──────────────────┘   But can't MODIFY without breaking signature

SESSION_ENCRYPTION_KEY (Encryption — AES):
┌──────────────────┐        ┌──────────────────┐
│ session data     │ ──────→│ a8f3k2...x9m1   │
│ user_id=42       │ encrypt│ (unreadable)     │
│ username=tijo    │        │                  │
└──────────────────┘        └──────────────────┘
                            Nobody can READ or MODIFY

Both together:
  Encrypt data → Sign the encrypted blob → Send as cookie
  Cookie arrives → Verify signature → Decrypt data → Use
```

⚠️ **Security rules**:

📱 **Narration**:
> "SESSION_SECRET — minimum 32 characters. Random generate ചെയ്യണം. `openssl rand -base64 32` terminal-ൽ run ചെയ്താൽ ഒരു random key കിട്ടും."

> "SESSION_ENCRYPTION_KEY — **exactly 32 bytes** ആയിരിക്കണം. AES-256 encryption-ന്. 16, 24, അല്ലെങ്കിൽ 32 bytes valid ആണ്, പക്ഷേ 32 strongest."

> "രണ്ടും `.env` file-ൽ. Production-ൽ real environment variables. **Git-ൽ commit ചെയ്യരുത്!**"

---

## 14:00 — Create .env.example

📱 **Narration**:
> ".env git-ൽ commit ചെയ്യില്ല. But other developers-ന് ഏത് variables വേണം എന്ന് അറിയണം. അതിന് .env.example create ചെയ്യാം."

⌨️ **Create `.env.example`**:
```env
# Server Configuration
PORT=8080
ENV=development

# Logging
LOG_LEVEL=DEBUG

# Database
DATABASE_PATH=./data/linkbio.db

# Session
SESSION_SECRET=your-super-secret-key-change-in-production
SESSION_ENCRYPTION_KEY=must-be-exactly-32-bytes-long!!

# Rate Limiting
RATE_LIMIT=10
```

🧠 **Explain**:
> "ഈ file git-ൽ commit ചെയ്യും. New developer clone ചെയ്യുമ്പോൾ `.env.example` copy ചെയ്ത് `.env` ആക്കിയാൽ മതി. Real secrets ഇതിൽ ഇല്ല."

---

## 14:15 — Create .gitignore

📱 **Narration**:
> "Git-ൽ commit ചെയ്യരുത്ത files specify ചെയ്യാം."

⌨️ **Create `.gitignore`**:
```gitignore
# Binaries
bin/
*.exe
*.exe~
*.dll
*.so
*.dylib

# Test binary
*.test

# Output of go coverage
*.out

# Dependency directories
vendor/

# Environment files
.env
.env.local
.env.*.local

# Database
data/*.db
*.db

# IDE
.idea/
.vscode/
*.swp
*.swo

# OS
.DS_Store
Thumbs.db

# Air (hot reload)
tmp/
```

🧠 **Explain key entries**:
> "`.env` — secrets git-ൽ commit ചെയ്യരുത്. ⚠️ ഇത് forget ചെയ്താൽ secrets public ആകും!"

> "`data/*.db` — database file. Each developer-ന് own data ഉണ്ടാകും."

> "`*.exe` — compiled binaries. Source code commit ചെയ്താൽ മതി."

---

## 14:30 — Create Makefile

📱 **Narration**:
> "Development workflow easy ആക്കാൻ ഒരു Makefile create ചെയ്യാം. ഇത് required അല്ല — just a convenience tool. Long commands remember ചെയ്യണ്ട."

⌨️ **Create `Makefile`**:
```makefile
.PHONY: run build test clean dev tidy

# Run the server
run:
	go run ./cmd/server

# Build binary
build:
	go build -o bin/linkbio.exe ./cmd/server

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf bin/ data/*.db

# Download dependencies
tidy:
	go mod tidy

# Development with hot reload (requires air)
dev:
	air

# Create .env from example
env:
	cp .env.example .env
```

🧠 **Explain each command**:

> "`make run` — server quick ആയി start ചെയ്യാൻ. `go run ./cmd/server` full command type ചെയ്യണ്ട."

> "`make build` — binary compile ചെയ്യാൻ. `bin/linkbio.exe` ആയി output ആകും."

> "`make test` — എല്ലാ tests run ചെയ്യാൻ. `-v` verbose — ഓരോ test-ന്റെ result കാണിക്കും."

> "`make clean` — build artifacts-ഉം database file-ഉം delete ചെയ്യാൻ. Fresh start വേണമെങ്കിൽ."

> "`.PHONY` — Make-നോട് പറയുന്നു ഇവ files അല്ല, commands ആണ്. `run` എന്ന folder ഉണ്ടെങ്കിലും `make run` command execute ചെയ്യും."

📱 **Narration**:
> "Makefile ഒരു convenience tool ആണ്. ഇത് ഇല്ലാതെയും project run ചെയ്യാം — direct go commands use ചെയ്താൽ മതി. But team projects-ൽ Makefile ഉണ്ടെങ്കിൽ everyone same commands use ചെയ്യും."

---

## 14:45 — Verify Setup

📱 **Narration**:
> "ഇനി verify ചെയ്യാം — everything correct ആണോ."

⌨️ **Type this**:
```bash
go mod tidy
```

🧠 **Explain**:
> "go mod tidy — unnecessary dependencies remove ചെയ്യും, missing dependencies add ചെയ്യും. go.sum file update ചെയ്യും."

🎥 **Camera**: Show `go.mod` file in VS Code. Point out the module name and dependencies.

📱 **Narration**:
> "go.mod നോക്കൂ. Module name `linkbio`, Go version, എല്ലാ dependencies list ചെയ്തിട്ടുണ്ട്. Package.json പോലെ."

---

## 14:55 — Transition

📱 **Narration**:
> "Project setup complete. Folders ready, dependencies installed, secrets configured. ഇനി actual code എഴുതാൻ തുടങ്ങാം — config loading-ൽ നിന്ന്!"

🔊 **Transition sound**

🎥 **Cut to**: Config file creation (Scene 4)

---

## 📝 Editing Notes

- **Terminal commands**: Type them live, don't paste. Viewers learn better watching you type.
- **Folder structure**: After creating all folders, show the VS Code sidebar with the full tree expanded. Hold for 3-4 seconds so viewers can screenshot.
- **Security section**: The seal vs locked box analogy is crucial. Consider adding a simple graphic in post-production.
- **Pacing**: Dependencies section can feel tedious. Keep energy up. One-sentence explanation per package, move on.
- **On Windows**: If using PowerShell, use `New-Item -ItemType Directory -Force -Path` instead of `mkdir -p`. Or show both.
