# Scene 1: Opening Hook + App Demo

> ⏱️ **Timestamp**: 0:00 — 3:00
> 🎯 **Goal**: Grab attention in 10 seconds, show the finished app, list what they'll learn

---

## 0:00 — Intro

🔊 **Upbeat intro music** (3 seconds, fade under narration)

🎥 **Camera**: Facecam or direct screen recording. Show the finished LinkBio app in browser.

📱 **Narration**:
> "നമസ്കാരം! Kraken Coder-ലേക്ക് സ്വാഗതം."
>
> "ഈ channel-ൽ നമ്മൾ discuss ചെയ്യുന്നത്...
> frontend JS-heavy architecture...
> server-driven architecture...
> Golang industry patterns — ഇതൊക്കെയാണ്."
>
> "ഇന്ന് ഒരു interesting question നോക്കാം.
> Linktree അറിയാമല്ലോ?
> Influencers-ന് links ഒരു page-ൽ aggregate ചെയ്യാൻ ഉള്ള platform."
>
> "ഇത് build ചെയ്യാൻ React വേണോ? Angular വേണോ?
> ...അതോ hypermedia-based tools മതിയോ?
> Alpine, HTMX, Golang — ഇത് മാത്രം മതിയോ?"
>
> "Advantage എന്താ?
> React background-ൽ നിന്ന് വരുന്ന ആളുകൾക്ക്...
> ഈ aspect evaluate ചെയ്യാൻ ഒരു നല്ല chance ആണ് ഈ first video."
>
> "Web app enthusiast ആണോ?
> എങ്കിൽ deep dive ചെയ്യാം.
> Golang, HTMX, Alpine JS — ഒരു Linktree clone build ചെയ്യാം!"

---

## 0:15 — Show Finished App Demo

🎥 **Camera**: Full screen browser. Click through each feature slowly.

### Demo 1: Landing Page
🎥 Show the landing page with GSAP animations loading in.

📱 **Narration**:
> "ഇത് നമ്മുടെ landing page. GSAP animations ഉണ്ട് — smooth, professional."

### Demo 2: Registration & Login
🎥 Click "Register" → Fill form → Submit. Show HTMX form submission (no page reload).

📱 **Narration**:
> "Registration. ശ്രദ്ധിക്കൂ — page reload ഇല്ല! HTMX ആണ് magic. Server HTML return ചെയ്യുന്നു, JSON അല്ല."

### Demo 3: Dashboard with Stats
🎥 Show dashboard — link count, total views, total clicks.

📱 **Narration**:
> "Dashboard. ഇവിടെ analytics ഉണ്ട് — എത്ര views, എത്ര clicks. Real-time data."

### Demo 4: Add/Edit Links
🎥 Click "Add Link" → Fill title + URL → Submit. Show link appearing without page reload.

📱 **Narration**:
> "Links add ചെയ്യാം. HTMX use ചെയ്തിട്ട് page reload ഇല്ലാതെ list update ആകുന്നു."

### Demo 5: Public Profile Page
🎥 Navigate to `localhost:8080/@username`. Show the public profile with links.

📱 **Narration**:
> "ഇത് public profile. ആരെങ്കിലും ഈ link share ചെയ്താൽ ഇത് കാണും. Clean, beautiful."

### Demo 6: Dark Mode Toggle
🎥 Click dark mode toggle. Show Alpine.js toggling theme instantly.

📱 **Narration**:
> "Dark mode! Alpine.js ഒരു line code കൊണ്ട് toggle ചെയ്യുന്നു."

---

## 1:30 — What You'll Learn

🎥 **Camera**: Show a bullet list on screen (use a slide or text overlay in editing).

📱 **Narration**:
> "ഈ tutorial-ൽ നിങ്ങൾ പഠിക്കുന്നത്:"

📱 **Read each point, pause between them:**

> "ഒന്ന് — Go backend with industry-standard patterns. cmd, internal, repository pattern — large companies use ചെയ്യുന്ന structure."

> "രണ്ട് — HTMX for dynamic updates. JavaScript ഒരു line എഴുതാതെ interactive pages!"

> "മൂന്ന് — Alpine.js for client-side magic. Toggles, dropdowns, dark mode — minimal JavaScript."

> "നാല് — Tailwind CSS for beautiful UI. Responsive, dark mode support."

> "അഞ്ച് — SQLite database. SQL എഴുതാൻ പഠിക്കാം, ORM ഇല്ല — real SQL."

> "ആറ് — Session encryption. Password hashing, encrypted cookies — production-level security."

> "ഏഴ് — Click analytics. ഓരോ link click track ചെയ്യും, dashboard-ൽ stats കാണിക്കും."

---

## 2:30 — Transition to Architecture

🎥 **Camera**: Switch from browser to VS Code. Show an empty folder.

📱 **Narration**:
> "ഇത് എല്ലാം നമ്മൾ scratch-ൽ നിന്ന് build ചെയ്യും. ഈ empty folder-ൽ നിന്ന്. Ready ആണോ? Let's go!"

🔊 **Quick transition sound effect** (whoosh or click)

🎥 **Cut to**: Architecture diagram (Scene 2)

---

## 📝 Editing Notes

- **Intro music**: Keep it under 3 seconds. Fade under voice immediately.
- **App demo**: Record this separately with a polished version of the app. Seed test data (3-4 links, some analytics) so the dashboard looks populated.
- **Pacing**: Don't rush the demo. Viewers need to see what they're building. 15 seconds per feature minimum.
- **Text overlay**: During "What You'll Learn", show bullet points appearing one by one (add in post-production).
- **Energy**: This section sets the tone. Be enthusiastic but not over-the-top. Confident and clear.

### Test Data for Demo
Before recording the demo, seed the database:
- User: `demo` / `demo@example.com`
- Links: YouTube, GitHub, Twitter, Portfolio (4 links)
- Analytics: 150+ views, 80+ clicks (looks impressive on dashboard)
