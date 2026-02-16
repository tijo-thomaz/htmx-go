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

> 🎥 **Transition:** "App complete! ഇനി deploy."
