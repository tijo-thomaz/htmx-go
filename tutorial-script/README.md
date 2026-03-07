# 🎬 LinkBio Tutorial — 3-Part YouTube Series

> **YouTube Series**: Build a Production-Grade Link-in-Bio Platform
> **Total Duration**: ~2 hours (3 videos)
> **Language**: Malayalam (Code in English)
> **Style**: Movie-like production with clear scene transitions

---

## 📺 Series Structure

### Part 1 — Foundation & Setup (~40 min)
*"From empty folder to running server"*

| # | File | What Happens |
|---|------|-------------|
| 00 | `00-pre-production.md` | Checklist, setup, camera notes |
| 00a | `00a-obs-recording-guide.md` | OBS recording setup |
| 01 | `01-opening-hook.md` | Demo finished app, hook the viewer |
| 02 | `02-architecture.md` | Draw diagram, explain layers |
| 03 | `03-project-setup.md` | Folders, deps, .env |
| 04 | `04-config-and-logging.md` | Config struct, slog logger |

**Ends with**: Server starts, `/health` returns `{"status":"ok"}`

---

### Part 2 — Backend & Logic (~45 min)
*"From models to working auth"*

| # | File | What Happens |
|---|------|-------------|
| 05 | `05-database-models.md` | User, Link, Analytics models |
| 06 | `06-database-repos.md` | Repositories, migrations, queries |
| 07 | `07-middleware-security.md` | Auth, session encryption, JWT vs Sessions |
| 08 | `08-handlers-auth.md` | Login, register, logout, bcrypt |
| 09 | `09-handlers-links.md` | CRUD, click tracking, context.Background bug |
| 10 | `10-handlers-dashboard.md` | Dashboard, analytics, stats polling, profile |
| 11 | `11-router-server.md` | Router (public vs protected), server wiring, main.go |

**Ends with**: Register, login, add links, click tracking — all working via curl/Postman

---

### Part 3 — Frontend & Polish (~35 min)
*"From raw HTML to beautiful app"*

| # | File | What Happens |
|---|------|-------------|
| 12 | `12-templates-base.md` | Base layout, Tailwind, HTMX, Alpine.js setup |
| 13 | `13-templates-auth-pages.md` | Login/Register with Alpine.js UX |
| 14 | `14-templates-dashboard.md` | Dashboard, HTMX polling, drag-drop |
| 15 | `15-templates-profile.md` | Profile page, click-tracked links |
| 16 | `16-deploy.md` | Build, env vars, production checklist |
| 17 | `17-closing.md` | Recap, bugs avoided, CTA |

**Ends with**: Full working app with beautiful UI, dark mode, analytics

---

### Bonus (Separate Video)
| # | File | What Happens |
|---|------|-------------|
| 18 | `18-bonus-testing.md` | Tests, benchmarks, coverage |

---

## 🎬 Script Conventions

| Icon | Meaning |
|------|---------|
| 📱 | Read this narration aloud (on phone/teleprompter) |
| ⌨️ | Type this code on screen |
| 🧠 | Explain to viewer after typing |
| 🎯 | Analogy or teaching moment |
| ⚠️ | Bug or security issue to highlight |
| 🎥 | Camera/editing instruction |
| 🔊 | Sound effect or music cue |

---

## 🐛 Key Bugs We Fix (Teaching Highlights)

These are the "aha moments" viewers will remember:

1. **Signed vs Encrypted cookies** — Part 2, Scene 07
2. **Click route behind auth middleware** — Part 2, Scene 09, 11
3. **`r.Context()` cancelled in goroutine** — Part 2, Scene 09
4. **Why JWT is wrong for user sessions** — Part 2, Scene 07

---

## 🎬 Part Transitions

### End of Part 1
> 📱 "Foundation ready — config, logging, server running. Part 2-ൽ database, authentication, handlers — actual business logic build ചെയ്യും. Subscribe ചെയ്യൂ, miss ആകരുത്!"

### Start of Part 2
> 📱 "Part 1-ൽ നമ്മൾ project setup, config, logging — foundation build ചെയ്തു. Server running ആണ്. ഇന്ന് real code — database, auth, handlers!"

### End of Part 2
> 📱 "Backend complete — database, auth, CRUD, analytics. Part 3-ൽ frontend — HTMX, Alpine.js, beautiful UI. ഇത് miss ആകരുത്!"

### Start of Part 3
> 📱 "Part 2-ൽ backend complete ചെയ്തു — working API. ഇന്ന് frontend — templates, HTMX magic, dark mode, drag-drop. App beautiful ആക്കാം!"
