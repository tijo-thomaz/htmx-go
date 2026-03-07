# Scene 2: Architecture Diagram

> ⏱️ **Timestamp**: 3:00 — 8:00
> 🎯 **Goal**: Explain the 3-layer architecture, why HTMX sends HTML not JSON, and the full request flow

---

## 3:00 — Show the Diagram

🎥 **Camera**: Full screen. Either draw this live on a whiteboard app (Excalidraw) or show a pre-made diagram. Highlight each layer as you explain.

🔊 **Soft background music** (very low, thinking/focus music)

📱 **Narration**:
> "Code എഴുതുന്നതിന് മുമ്പ്, നമ്മുടെ app-ന്റെ architecture understand ചെയ്യാം. ഈ diagram നോക്കൂ."

### The Diagram

```
┌─────────────────────────────────────────────────────────────┐
│                         BROWSER                              │
│                                                              │
│  ┌──────────────────────┐  ┌───────────────────────────┐    │
│  │  Alpine.js            │  │  GSAP                     │    │
│  │  • Dark mode toggle  │  │  • Page load animations   │    │
│  │  • Dropdowns         │  │  • Fade-in effects        │    │
│  │  • Form validation   │  │  • Stagger animations     │    │
│  └──────────────────────┘  └───────────────────────────┘    │
│                                                              │
│  ┌──────────────────────────────────────────────────────┐    │
│  │  HTMX                                                │    │
│  │  • hx-post="/links"     → sends form data            │    │
│  │  • hx-delete="/links/5" → sends DELETE request       │    │
│  │  • hx-swap="innerHTML"  → replaces HTML on page      │    │
│  │  • Server returns HTML, NOT JSON!                    │    │
│  └──────────────────────────────────────────────────────┘    │
└─────────────────────────┬───────────────────────────────────┘
                          │
                          │  HTTP Request (HTML response!)
                          ▼
┌─────────────────────────────────────────────────────────────┐
│                       GO SERVER                              │
│                                                              │
│  ┌──────────┐  ┌────────────┐  ┌──────────┐  ┌──────────┐  │
│  │  Router  │→ │ Middleware │→ │ Handler  │→ │ Template │  │
│  │          │  │            │  │          │  │          │  │
│  │ chi      │  │ • Auth     │  │ • Logic  │  │ • HTML   │  │
│  │ matches  │  │ • Logging  │  │ • DB     │  │ • Render │  │
│  │ URL path │  │            │  │   calls  │  │ • Data   │  │
│  │          │  │            │  │          │  │          │  │
│  └──────────┘  └────────────┘  └──────────┘  └──────────┘  │
└─────────────────────────┬───────────────────────────────────┘
                          │
                          │  SQL Queries
                          ▼
┌─────────────────────────────────────────────────────────────┐
│                       SQLite Database                        │
│                                                              │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────────┐  │
│  │    users      │  │    links      │  │    analytics      │  │
│  │              │  │              │  │                  │  │
│  │ id           │  │ id           │  │ id               │  │
│  │ username     │  │ user_id  FK  │  │ user_id      FK  │  │
│  │ email        │  │ title        │  │ link_id      FK  │  │
│  │ password_hash│  │ url          │  │ event_type       │  │
│  │ display_name │  │ icon         │  │ referrer         │  │
│  │ bio          │  │ position     │  │ user_agent       │  │
│  │ avatar_url   │  │ is_active    │  │ created_at       │  │
│  │ theme        │  │ created_at   │  │                  │  │
│  │ created_at   │  │              │  │                  │  │
│  └──────────────┘  └──────────────┘  └──────────────────┘  │
│                                                              │
│  Foreign Keys: links.user_id → users.id                      │
│                analytics.user_id → users.id                  │
│                analytics.link_id → links.id                  │
└─────────────────────────────────────────────────────────────┘
```

---

## 3:30 — Layer 1: Browser

🎥 **Camera**: Highlight the Browser section of the diagram.

📱 **Narration**:
> "ആദ്യ layer — Browser. User-ന്റെ browser-ൽ മൂന്ന് tools work ചെയ്യുന്നു."

📱 **Alpine.js**:
> "Alpine.js — ചെറിയ interactive things handle ചെയ്യാൻ. Dark mode toggle, dropdown menus. React-ന്റെ baby brother എന്ന് വിചാരിക്കാം — simple, lightweight."

🎯 **Analogy**:
> "Alpine.js ഒരു Swiss Army knife പോലെ. ചെറുതാണ്, പക്ഷേ ആവശ്യമുള്ളത് എല്ലാം ഉണ്ട്. React ഒരു full toolbox ആണ് — powerful, but overkill for our needs."

📱 **GSAP**:
> "GSAP — GreenSock Animation Platform. CSS animations-നേക്കാൾ powerful."
> "Page load ചെയ്യുമ്പോൾ navbar slide-in, link cards stagger ചെയ്ത് appear — smooth, professional feel."
> "Button hover-ൽ scale effect, link cards-ൽ subtle slide — ഇതൊക്കെ GSAP."
> "CSS animations-ൽ stagger, timeline, easing control ചെയ്യാൻ complicated ആണ്. GSAP-ൽ one line."

📱 **AOS**:
> "AOS — Animate On Scroll. Scroll ചെയ്യുമ്പോൾ elements fade-up, fade-in ചെയ്യും. Landing page features section-ൽ use ചെയ്യും."

📱 **HTMX**:
> "HTMX — ഇത് നമ്മുടെ star player. Server-ലേക്ക് requests അയക്കുന്നത് ഇതാണ്."

---

## 4:30 — Why HTMX Sends HTML, Not JSON

🎥 **Camera**: Show a side-by-side comparison (draw or show a slide).

📱 **Narration**:
> "ഇവിടെ ഒരു important concept ഉണ്ട്. React, Vue — ഇവ JSON APIs use ചെയ്യും. HTMX different ആണ്."

🎥 **Show comparison**:

```
Traditional SPA (React):
  Browser → GET /api/links → Server
  Server → {"links": [{"title": "YouTube", "url": "..."}]} → Browser
  Browser: JavaScript parses JSON, creates DOM elements, renders

HTMX approach:
  Browser → GET /links → Server
  Server → <div class="link-card"><h3>YouTube</h3>...</div> → Browser
  Browser: HTMX swaps HTML directly into the page. Done!
```

📱 **Narration**:
> "React approach-ൽ server JSON അയക്കുന്നു. Browser-ൽ JavaScript ആ JSON parse ചെയ്ത് HTML create ചെയ്യണം. Double work!"

> "HTMX approach-ൽ server direct HTML അയക്കുന്നു. Browser-ൽ HTMX ആ HTML page-ൽ insert ചെയ്യും. Simple, fast, no JavaScript needed!"

🎯 **Analogy**:
> "Restaurant-ൽ ഒരു comparison. React approach: Kitchen ingredients അയക്കുന്നു, customer cook ചെയ്യണം. HTMX approach: Kitchen ready-made food അയക്കുന്നു, customer eat ചെയ്താൽ മതി. ആരാണ് smart?"

---

## 5:30 — Layer 2: Go Server

🎥 **Camera**: Highlight the Go Server section. Point to each box left to right.

📱 **Narration**:
> "രണ്ടാം layer — Go server. ഇത് നമ്മുടെ app-ന്റെ brain. ഇവിടെ നാല് stages ഉണ്ട്."

📱 **Router**:
> "Router — chi library use ചെയ്യുന്നു. URL match ചെയ്യുന്ന പണി. `/login` വന്നാൽ login handler-ലേക്ക്, `/dashboard` വന്നാൽ dashboard handler-ലേക്ക്."

🎯 **Analogy**:
> "Router ഒരു building-ലെ receptionist പോലെ. നിങ്ങൾ വരുമ്പോൾ ചോദിക്കും — 'എവിടെ പോകണം?' — ശരിയായ room-ലേക്ക് direct ചെയ്യും."

📱 **Middleware**:
> "Middleware — request handler-ൽ reach ആകുന്നതിന് മുമ്പ് check ചെയ്യുന്ന code. Authentication, logging."

🎯 **Analogy**:
> "Middleware ഒരു security guard പോലെ. Building-ൽ enter ചെയ്യുന്നതിന് മുമ്പ് ID check ചെയ്യും. Valid അല്ലെങ്കിൽ — 'Sorry, you can't enter.' Valid ആണെങ്കിൽ — 'Go ahead.'"

📱 **Handler**:
> "Handler — actual business logic. Form data read ചെയ്യുന്നത്, database-ൽ save ചെയ്യുന്നത്, validation — എല്ലാം ഇവിടെ."

📱 **Template**:
> "Template — Go-യുടെ built-in html/template use ചെയ്ത് HTML generate ചെയ്യുന്നു. Data + Template = Final HTML page."

---

## 6:30 — Layer 3: Database

🎥 **Camera**: Highlight the SQLite section. Point to each table.

📱 **Narration**:
> "മൂന്നാം layer — SQLite database. ഒരു single file — `linkbio.db`. Install ഒന്നും വേണ്ട, server ഒന്നും run ചെയ്യണ്ട."

📱 **Tables**:
> "മൂന്ന് tables ഉണ്ട്. `users` — accounts store ചെയ്യാൻ. `links` — user-ന്റെ links. `analytics` — click tracking."

📱 **Foreign Keys**:
> "Foreign keys — tables-നെ connect ചെയ്യുന്നു. ഒരു user delete ചെയ്താൽ, ആ user-ന്റെ links-ഉം analytics-ഉം automatic delete ആകും. ON DELETE CASCADE — ഇത് database handle ചെയ്യും."

🎯 **Analogy**:
> "Foreign keys ഒരു family tree പോലെ. Parent delete ചെയ്താൽ children-ഉം delete ആകും. Orphan data ഉണ്ടാകില്ല."

---

## 7:00 — Full Request Flow

🎥 **Camera**: Draw or animate the full request path. Arrow from browser through each layer and back.

📱 **Narration**:
> "ഇപ്പോൾ full picture. ഒരു user 'Add Link' button click ചെയ്യുമ്പോൾ എന്ത് സംഭവിക്കുന്നു?"

📱 **Step by step**:
> "Step 1: HTMX browser-ൽ നിന്ന് POST request അയക്കുന്നു — `/links` endpoint-ലേക്ക്, form data-യോട് കൂടി."

> "Step 2: Go server-ൽ Router match ചെയ്യുന്നു — 'POST /links? അത് LinkHandler.Create-ലേക്ക്.'"

> "Step 3: Middleware check ചെയ്യുന്നു — 'ഈ user logged in ആണോ? Session valid ആണോ?'"

> "Step 4: Handler form data read ചെയ്യുന്നു, validate ചെയ്യുന്നു, database-ൽ save ചെയ്യുന്നു."

> "Step 5: Template render ചെയ്യുന്നു — new link-ന്റെ HTML fragment."

> "Step 6: HTML browser-ലേക്ക് return. HTMX ആ HTML page-ൽ insert ചെയ്യുന്നു. Done! Page reload ഇല്ല!"

🎥 **Show flow as arrows**:
```
User clicks "Add Link"
    ↓
HTMX: POST /links {title: "YouTube", url: "https://youtube.com"}
    ↓
Router: POST /links → LinkHandler.Create
    ↓
Middleware: Auth check ✅ → Logging ✅
    ↓
Handler: Validate → Save to SQLite → Render template
    ↓
Response: <div class="link-card">YouTube</div>
    ↓
HTMX: Swaps HTML into page → User sees new link instantly
```

---

## 7:45 — Transition

📱 **Narration**:
> "Architecture clear ആയോ? ഓരോ layer-ഉം separate ആണ്, maintain ചെയ്യാൻ easy ആണ്. ഇനി code എഴുതാൻ തുടങ്ങാം!"

🔊 **Transition sound** (whoosh)

🎥 **Cut to**: VS Code with empty folder (Scene 3)

---

## 📝 Editing Notes

- **Diagram**: Use Excalidraw (free) or draw.io. Export as clean PNG. Or draw live for engagement.
- **Highlighting**: When explaining each layer, dim the other layers and highlight the current one (use rectangles/borders in post-production).
- **Request flow**: Animate with arrows appearing one at a time. Each arrow appears as you narrate that step.
- **Pacing**: This is a dense section. Speak slowly. Pause between layers. Viewers need time to process.
- **Side-by-side**: The JSON vs HTML comparison is the key "aha moment" of this scene. Make it visually clear.
