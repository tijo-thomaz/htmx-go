# Scene 14: Dashboard Template (53:00 - 56:00)

> 🎬 **Previous**: Auth pages (Scene 13)
> 🎯 **Goal**: Dashboard with stats, HTMX forms, OOB swaps, drag-drop, clipboard

---

## Analytics Stats Row

```html
<div class="grid grid-cols-2 gap-4">
    <!-- Total Views -->
    <div class="stat-card bg-white dark:bg-gray-900 rounded-2xl p-6 border border-gray-100 dark:border-gray-800">
        <p class="text-2xl font-bold text-gray-900 dark:text-white">
            {{if .Analytics}}{{.Analytics.TotalViews}}{{else}}0{{end}}
        </p>
        <p class="text-sm text-gray-500">Total Views</p>
    </div>
    <!-- Total Clicks -->
    <div class="stat-card bg-white dark:bg-gray-900 rounded-2xl p-6 border border-gray-100 dark:border-gray-800">
        <p class="text-2xl font-bold text-gray-900 dark:text-white">
            {{if .Analytics}}{{.Analytics.TotalClicks}}{{else}}0{{end}}
        </p>
        <p class="text-sm text-gray-500">Total Clicks</p>
    </div>
</div>
```

> 🧠 📱 "{{if .Analytics}} — nil check. Analytics error ആയാൽ nil, template crash ആകില്ല."
> 📱 "Click tracking work ചെയ്യുന്നുണ്ടെങ്കിൽ TotalClicks count കാണും. 0 ആയാൽ bugs check (Scene 9)."

---

## Add Link Form — Toggle + Transition + HTMX

```html
<main x-data="{ showAddForm: false }">
    <button @click="showAddForm = !showAddForm">Add Link</button>
    
    <div x-show="showAddForm" x-cloak 
         x-transition:enter="transition ease-out duration-200"
         x-transition:enter-start="opacity-0 -translate-y-2"
         x-transition:enter-end="opacity-100 translate-y-0">
        <form hx-post="/api/v1/links" 
              hx-target="#links-list" 
              hx-swap="afterbegin"
              hx-indicator="find .htmx-indicator"
              @htmx:after-request="showAddForm = false; $el.reset()">
            <input type="text" name="title" required placeholder="My Website">
            <input type="url" name="url" required placeholder="https://example.com">
            <button type="submit">
                <svg class="w-4 h-4 animate-spin htmx-indicator"><!-- spinner --></svg>
                Add Link
            </button>
        </form>
    </div>
</main>
```

> 🧠 📱 "x-transition — smooth slide-down animation. Alpine.js CSS transitions manage ചെയ്യും."
> 📱 "@htmx:after-request — response വന്നാൽ form close + reset. $el = current form element."
> 📱 "hx-swap='afterbegin' — new link list-ന്റെ top-ൽ insert."
> 📱 "htmx-indicator — HTMX request ആയാൽ spinner visible, complete ആയാൽ hidden."

> 🎯 📱 "Drawer open/close. Button click → slide down. Submit → slide up, form clear. Smooth UX!"

---

## Links List with Drag-Drop (SortableJS)

```html
<div id="links-list" x-data x-init="
    new Sortable($el, {
        animation: 200,
        ghostClass: 'sortable-ghost',
        handle: '.drag-handle',
        onEnd: function(evt) {
            const items = evt.to.querySelectorAll('[data-link-id]');
            const positions = {};
            items.forEach((item, index) => {
                positions[parseInt(item.dataset.linkId)] = index + 1;
            });
            fetch('/api/v1/links/reorder', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(positions)
            });
        }
    })
">
    {{range .Links}}
    <div class="link-card" data-link-id="{{.ID}}">
        <button class="drag-handle">⠿</button>
        <h3>{{.Title}}</h3>
        <p>{{.URL}}</p>
        <button hx-delete="/api/v1/links/{{.ID}}"
                hx-target="closest .link-card"
                hx-swap="outerHTML swap:200ms"
                hx-confirm="Delete this link?">🗑️</button>
    </div>
    {{end}}
</div>
```

> 🧠 📱 "SortableJS — drag-drop library. onEnd callback-ൽ positions collect ചെയ്ത് server-ന് send."
> 📱 "hx-delete — HTMX DELETE request. hx-confirm — browser confirm dialog."
> 📱 "hx-swap='outerHTML swap:200ms' — 200ms fade out, then remove."

---

## Copy Profile Link — Clipboard API

```html
<button x-data="{ copied: false }"
        @click="navigator.clipboard.writeText(window.location.origin + '/u/{{.User.Username}}'); 
                copied = true; setTimeout(() => copied = false, 2000)"
        :class="copied ? 'bg-green-500 text-white' : 'bg-gray-100 text-gray-600'">
    <span x-text="copied ? '✓ Copied!' : 'Copy Profile Link'"></span>
</button>
```

> 🧠 📱 "navigator.clipboard.writeText() — browser clipboard API."
> 📱 "setTimeout 2 seconds — '✓ Copied!' feedback, then reset."
> 📱 "Nested x-data — button-ന്റെ own state. Parent-ന്റെ state affect ചെയ്യില്ല."

---

> 🎥 **Transition:** "Dashboard done. ഇനി public profile page."
