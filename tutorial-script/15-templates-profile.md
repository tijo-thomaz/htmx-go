# Scene 15: Profile Page — Click-Tracked Links (56:00 - 58:00)

> 🎬 **Previous**: Dashboard (Scene 14)
> 🎯 **Goal**: Public profile with dark mode, share, and click-tracked links

---

## Profile Dark Mode (Independent)

```html
<div x-data="{ darkMode: false }" 
     :class="darkMode ? 'dark bg-gradient-to-b from-gray-900 to-gray-800' 
                       : 'bg-gradient-to-b from-gray-50 to-gray-100'">
```

> 🧠 📱 "Profile darkMode = independent. Dashboard darkMode localStorage save ചെയ്യും, profile-ന്റേത് save ചെയ്യില്ല. Visitor toggle, but not persist."

---

## Theme Toggle with x-if

```html
<button @click="darkMode = !darkMode">
    <template x-if="!darkMode"><span>🌙 Dark Mode</span></template>
    <template x-if="darkMode"><span>☀️ Light Mode</span></template>
</button>
```

> 🧠 📱 "x-if vs x-show: x-show = CSS display:none, element DOM-ൽ ഉണ്ട്. x-if = completely add/remove from DOM. Icon swap-ന് x-if better."

---

## ⚠️ Click-Tracked Links — The Critical Part

```html
{{range $index, $link := .Links}}
<a href="/click/{{$link.ID}}" 
   target="_blank"
   rel="noopener"
   class="link-button block w-full p-5 rounded-2xl text-center font-medium"
   data-aos="fade-up" 
   data-aos-delay="{{multiply $index 50}}">
    <span class="text-lg">{{$link.Title}}</span>
</a>
{{end}}
```

> 🧠 **Explain the click tracking flow:**
> 📱 "href='/click/{{$link.ID}}' — direct URL-ലേക്ക് point ചെയ്യുന്നില്ല!"
> 📱 "Flow: User click → GET /click/5 → Server records click (async goroutine) → Server redirects 307 to actual URL → User lands on destination."
> 📱 "User-ന് seamless experience. Redirect instant ആണ്. Analytics async."

> ⚠️ **Critical reminder:**
> 📱 "⚠️ href='/click/...' ആണ്, '/api/v1/analytics/click/...' അല്ല!"
> 📱 "API route auth-protected. Profile visitors authenticated അല്ല."
> 📱 "Wrong route use ചെയ്താൽ → login redirect → click lost → dashboard 0 clicks forever."
> 📱 "ഒരു path string mistake. Error message ഇല്ല. Silent failure. Production-ൽ ഇത് weeks കഴിഞ്ഞ് മാത്രം notice ചെയ്യും!"

---

## Share Button

```html
<button x-data="{ copied: false }"
        @click="navigator.clipboard.writeText(window.location.href); 
                copied = true; setTimeout(() => copied = false, 2000)"
        :class="copied ? 'bg-green-500 text-white' 
                : (darkMode ? 'bg-gray-800/80 text-gray-300' : 'bg-white/80 text-gray-600')">
    <template x-if="!copied"><span>📤 Share</span></template>
    <template x-if="copied"><span>✓ Copied!</span></template>
</button>
```

> 🧠 📱 "Parent darkMode + own copied state. Alpine.js scope chain — inner component parent access ചെയ്യാം!"

---

## GSAP + AOS + HTMX Integration — app.js

> 📱 "ഇപ്പോൾ animations add ചെയ്യാം. ഇത് app-നെ professional look കൊടുക്കും."

**⌨️ Create `web/static/js/app.js`:**
```js
/**
 * LinkBio - Main JavaScript
 * Handles animations, HTMX events, and interactions
 */

document.addEventListener('DOMContentLoaded', function() {
    initAOS();
    initGSAP();
    initHTMXHandlers();
});

// AOS (Animate on Scroll)
function initAOS() {
    if (typeof AOS !== 'undefined') {
        AOS.init({
            duration: 600,
            easing: 'ease-out-cubic',
            once: true,
            offset: 50,
        });
    }
}

// GSAP Animations
function initGSAP() {
    if (typeof gsap === 'undefined') return;

    // Navbar fade in from top
    gsap.from('nav', { 
        opacity: 0, y: -20, duration: 0.8, ease: 'power3.out' 
    });

    // Profile page: link buttons stagger animation
    const linkButtons = document.querySelectorAll('.link-button');
    if (linkButtons.length > 0) {
        gsap.from('.link-button', {
            opacity: 0, y: 30, stagger: 0.08,
            duration: 0.6, ease: 'power3.out', delay: 0.4
        });
    }

    // Hover effects — scale up buttons
    document.querySelectorAll('.btn-primary').forEach(function(btn) {
        btn.addEventListener('mouseenter', function() {
            gsap.to(this, { scale: 1.02, duration: 0.2 });
        });
        btn.addEventListener('mouseleave', function() {
            gsap.to(this, { scale: 1, duration: 0.2 });
        });
    });

    // Link cards — subtle slide on hover
    document.querySelectorAll('.link-card').forEach(function(card) {
        card.addEventListener('mouseenter', function() {
            gsap.to(this, { x: 4, duration: 0.2 });
        });
        card.addEventListener('mouseleave', function() {
            gsap.to(this, { x: 0, duration: 0.2 });
        });
    });
}

// HTMX event handlers
function initHTMXHandlers() {
    // Show errors from HTMX responses
    document.body.addEventListener('htmx:beforeSwap', function(evt) {
        if (evt.detail.xhr.status >= 400) {
            evt.detail.shouldSwap = true;
            evt.detail.target.innerHTML = '<div class="error-message rounded-xl px-4 py-3 text-sm">' 
                + evt.detail.xhr.responseText + '</div>';
        }
    });

    // Re-init AOS after HTMX adds new elements
    document.body.addEventListener('htmx:afterSwap', function() {
        if (typeof AOS !== 'undefined') AOS.refresh();
    });
}
```

> 🧠 **Explain key points:**
> 📱 "typeof gsap === 'undefined' — CDN load fail ആയാൽ crash ആകില്ല. Defensive coding."
> 📱 "gsap.from() — element-ന്റെ initial state define ചെയ്യുന്നു. opacity: 0, y: 30 — invisible + below. GSAP animate ചെയ്ത് normal position-ലേക്ക് bring ചെയ്യും."
> 📱 "stagger: 0.08 — ഓരോ link button 80ms gap-ൽ appear ആകും. Waterfall effect!"

> 🎯 **Analogy:**
> 📱 "stagger ഒരു dominos falling പോലെ. ഒരുമിച്ച് fall ചെയ്യുന്നതിന് പകരം one-by-one fall ചെയ്യുന്നു. Professional feel!"

> 🧠 **HTMX + AOS integration:**
> 📱 "htmx:afterSwap — HTMX new HTML insert ചെയ്യുമ്പോൾ AOS.refresh() call ചെയ്യുന്നു. ഇല്ലെങ്കിൽ new elements-ന് scroll animation ഉണ്ടാകില്ല."

---

> 🎥 **Transition:** "App complete! ഇനി deploy."
