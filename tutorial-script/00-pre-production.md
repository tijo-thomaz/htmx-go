# 🎬 Pre-Production Checklist

> Complete this checklist BEFORE hitting record. Nothing kills a tutorial faster than "hold on, let me install this..."

---

## 🖥️ Software Requirements

### Must Have
- [ ] **Go 1.22+** — Run `go version` to confirm
- [ ] **VS Code** — Latest stable version
- [ ] **Browser** — Chrome or Firefox (for DevTools network tab demos)
- [ ] **Terminal** — Windows Terminal or VS Code integrated terminal
- [ ] **Git** — `git --version`

### VS Code Extensions
- [ ] **Go** (`golang.go`) — Auto-complete, formatting, debugging
- [ ] **Tailwind CSS IntelliSense** (`bradlc.vscode-tailwindcss`) — Class autocomplete in templates
- [ ] **Better Comments** — Optional, makes TODO/FIXME colorful
- [ ] **SQLite Viewer** — To show database contents on screen

### VS Code Settings for Recording
```json
{
  "editor.fontSize": 18,
  "terminal.integrated.fontSize": 16,
  "editor.lineHeight": 28,
  "editor.minimap.enabled": false,
  "breadcrumbs.enabled": false,
  "editor.renderWhitespace": "none",
  "workbench.activityBar.location": "hidden",
  "workbench.statusBar.visible": false,
  "window.zoomLevel": 1
}
```

> 💡 **Tip**: Create a separate VS Code profile called "Recording" with these settings so your normal setup isn't affected.

---

## 📹 Camera & Screen Recording

### Recording Setup
- [ ] **Screen recorder** — OBS Studio (free) or Screenflow
- [ ] **Resolution** — 1920x1080 minimum, 2560x1440 preferred
- [ ] **Frame rate** — 30fps for screen, 60fps if showing animations
- [ ] **Microphone** — External mic preferred; test audio levels
- [ ] **Quiet room** — No fan/AC noise during narration

### Screen Layout
```
┌─────────────────────────────────────────┐
│  VS Code (80% width)  │  Browser (20%) │
│                        │                │
│  ┌──────────────────┐  │  ┌──────────┐  │
│  │   Editor area    │  │  │  App     │  │
│  │   (code here)    │  │  │  preview │  │
│  │                  │  │  │          │  │
│  ├──────────────────┤  │  └──────────┘  │
│  │   Terminal       │  │                │
│  └──────────────────┘  │                │
└─────────────────────────────────────────┘
```

### Camera Notes
- [ ] If using facecam — bottom-right corner, small circle crop
- [ ] Good lighting on face (window light or ring light)
- [ ] Clean background or blur

---

## 🎨 Thumbnail Ideas

### Option A: Split Screen
- Left side: Code editor with Go/HTMX code visible
- Right side: Beautiful LinkBio app preview (dark mode)
- Text overlay: "Build LinkBio 🔗" in bold
- Malayalam subtitle: "Go + HTMX Tutorial"

### Option B: Before/After
- Left: Empty VS Code folder
- Right: Finished app with links, analytics, dark mode
- Arrow between them
- Text: "0 → Production" in big font

### Option C: Tech Stack Logos
- Go gopher + HTMX + Alpine.js + Tailwind + SQLite logos arranged
- App screenshot in center
- Text: "Full Stack Go Tutorial"

### Thumbnail Specs
- 1280x720 pixels
- Use Canva or Figma
- High contrast text (readable on mobile)
- Faces get more clicks — include facecam frame if comfortable

---

## 📱 Script on Phone

- [ ] Open `tutorial-script/01-opening-hook.md` on phone or tablet
- [ ] Use a teleprompter app (PromptSmart, BigVu) if available
- [ ] Or just keep the phone on a stand next to the monitor
- [ ] Highlight the 📱 sections — those are your narration cues
- [ ] Practice reading the Malayalam sections aloud 2-3 times

---

## 🧪 Test Run (Do This!)

### Verify Go Environment
```bash
go version
# Should show: go1.22.x or higher

go env GOPATH
# Should show a valid path
```

### Verify Empty Project Works
```bash
mkdir test-linkbio
cd test-linkbio
go mod init linkbio
go get github.com/go-chi/chi/v5
go get modernc.org/sqlite
```
> If `modernc.org/sqlite` takes too long, it's compiling pure Go SQLite. Let it finish once so it's cached for recording.

### Verify Browser DevTools
- Open Chrome DevTools → Network tab
- Filter by "Doc" to show HTML responses
- This is how we'll demo HTMX sending HTML, not JSON

### Verify Terminal
- Clear terminal history
- Set a clean prompt (short, no clutter)
- Test that copy-paste works from script to terminal

---

## ⏱️ Recording Strategy

### Record in Sections
Don't try to record 60 minutes in one take. Break into chunks:

| Section | Duration | File |
|---------|----------|------|
| Opening hook + demo | 3 min | `01-opening-hook.md` |
| Architecture diagram | 5 min | `02-architecture.md` |
| Project setup | 7 min | `03-project-setup.md` |
| Config + Logging | 5 min | `04-config-and-logging.md` |
| Database + Models | 5 min | Next section |
| ... | ... | ... |

### Between Sections
- [ ] Take a 2-minute break
- [ ] Review next script section
- [ ] Clear terminal if cluttered
- [ ] Save all files in VS Code

---

## 🗂️ File Organization

Keep these open during recording:
```
tutorial-script/
├── 00-pre-production.md    ← You are here
├── 01-opening-hook.md      ← Scene 1: Hook + Demo
├── 02-architecture.md      ← Scene 2: Diagram
├── 03-project-setup.md     ← Scene 3: Folders + deps
├── 04-config-and-logging.md ← Scene 4: Config + Logger
└── ... (more scenes)
```

---

## ✅ Final Pre-Record Checklist

- [ ] Script files loaded on phone/tablet
- [ ] VS Code open, clean, zoomed in
- [ ] Terminal ready, history cleared
- [ ] Browser open, DevTools ready
- [ ] Mic tested, audio levels good
- [ ] Screen recorder tested, resolution confirmed
- [ ] Water bottle nearby
- [ ] Phone on silent
- [ ] Notifications OFF (Do Not Disturb mode)
- [ ] Read through Scene 1 one more time

> 🎬 **You're ready. Hit record.**
