# 🎬 LinkBio Tutorial — Scene-by-Scene Script

> **YouTube Tutorial**: Build a Production-Grade Link-in-Bio Platform
> **Duration**: ~60 minutes
> **Language**: Malayalam (Code in English)
> **Style**: Movie-like production with clear scene transitions

---

## 📁 Script Structure

Each scene is a self-contained file. Read them in order.

| # | File | Time | What Happens |
|---|------|------|-------------|
| 00 | `00-pre-production.md` | Before recording | Checklist, setup, camera notes |
| 01 | `01-opening-hook.md` | 0:00 - 3:00 | Demo finished app, hook the viewer |
| 02 | `02-architecture.md` | 3:00 - 8:00 | Draw diagram, explain layers |
| 03 | `03-project-setup.md` | 8:00 - 15:00 | Folders, deps, .env |
| 04 | `04-config-and-logging.md` | 15:00 - 20:00 | Config struct, slog logger |
| 05 | `05-database-models.md` | 20:00 - 23:00 | User, Link, Analytics models |
| 06 | `06-database-repos.md` | 23:00 - 28:00 | Repositories, migrations, queries |
| 07 | `07-middleware-security.md` | 28:00 - 33:00 | Auth, rate limit, session encryption, JWT vs Sessions |
| 08 | `08-handlers-auth.md` | 33:00 - 38:00 | Login, register, logout, bcrypt |
| 09 | `09-handlers-links.md` | 38:00 - 43:00 | CRUD, click tracking, context.Background bug fix |
| 10 | `10-handlers-dashboard.md` | 43:00 - 45:00 | Dashboard, analytics summary |
| 11 | `11-router-server.md` | 45:00 - 48:00 | Router (public vs protected), server wiring |
| 12 | `12-templates-base.md` | 48:00 - 50:00 | Base layout, Tailwind, HTMX, Alpine.js setup |
| 13 | `13-templates-auth-pages.md` | 50:00 - 53:00 | Login/Register with Alpine.js UX |
| 14 | `14-templates-dashboard.md` | 53:00 - 56:00 | Dashboard, OOB swaps, drag-drop |
| 15 | `15-templates-profile.md` | 56:00 - 58:00 | Profile page, click-tracked links |
| 16 | `16-deploy.md` | 58:00 - 59:00 | Build, env vars, production checklist |
| 17 | `17-closing.md` | 59:00 - 60:00 | Recap, bugs avoided, CTA |
| 18 | `18-bonus-testing.md` | Bonus | Tests, benchmarks, coverage |

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

1. **Signed vs Encrypted cookies** — Scene 07
2. **Click route behind auth middleware** — Scene 09, 11
3. **`r.Context()` cancelled in goroutine** — Scene 09
4. **Why JWT is wrong for user sessions** — Scene 07
