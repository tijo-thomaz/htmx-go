# Scene 12: Base Layout Template (48:00 - 50:00)

> 🎬 **Previous**: Router + server wired (Scene 11)
> 🎯 **Goal**: Base HTML layout with Tailwind, HTMX, Alpine.js, AOS, GSAP

---

**⌨️ Create `web/templates/layouts/base.html`:**
```html
{{define "base"}}
<!DOCTYPE html>
<html lang="en" x-data="{ darkMode: localStorage.getItem('darkMode') === 'true' }" :class="{ 'dark': darkMode }">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{block "title" .}}LinkBio{{end}}</title>
    
    <link href="https://fonts.googleapis.com/css2?family=Inter:wght@300;400;500;600;700;800&display=swap" rel="stylesheet">
    <script src="https://cdn.tailwindcss.com"></script>
    <script>
        tailwind.config = {
            darkMode: 'class',
            theme: { extend: { fontFamily: { sans: ['Inter', 'system-ui', 'sans-serif'] } } }
        }
    </script>
    
    <link rel="stylesheet" href="/static/css/app.css">
    <script src="https://unpkg.com/htmx.org@1.9.10"></script>
    <script defer src="https://unpkg.com/alpinejs@3.x.x/dist/cdn.min.js"></script>
    <link href="https://unpkg.com/aos@2.3.1/dist/aos.css" rel="stylesheet">
    <script src="https://unpkg.com/aos@2.3.1/dist/aos.js"></script>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/gsap/3.12.2/gsap.min.js"></script>
    
    {{block "head" .}}{{end}}
</head>
<body class="{{block "bodyClass" .}}bg-gray-50 dark:bg-gray-900{{end}} min-h-screen transition-colors duration-300">
    {{block "content" .}}{{end}}
    
    <script src="/static/js/app.js"></script>
    {{block "scripts" .}}{{end}}
</body>
</html>
{{end}}
```

> 🧠 **Explain each library:**
> 📱 "HTMX — server-sent HTML swap. JavaScript ഇല്ലാതെ dynamic pages."
> 📱 "Alpine.js — lightweight reactivity. Toggle, forms, clipboard. defer load."
> 📱 "Tailwind — utility-first CSS. Class-ൽ design."
> 📱 "AOS — Animate On Scroll. Elements scroll ചെയ്യുമ്പോൾ fade-in."
> 📱 "GSAP — smooth animations. Link cards stagger."

> 🧠 **Explain template blocks:**
> 📱 "{{block \"title\" .}} — child templates override ചെയ്യാം."
> 📱 "x-data on html tag — darkMode state global. localStorage-ൽ persist."
> 📱 ":class dark — Tailwind dark mode CSS activate."

---

## Home Page — Landing with Animations

> 📱 "ഇനി home page — landing page with gradient background, scroll animations."

**⌨️ Create `web/templates/pages/home.html`:**

```html
{{define "title"}}LinkBio - Your Link in Bio{{end}}

{{define "bodyClass"}}gradient-bg overflow-x-hidden{{end}}

{{define "content"}}
    <!-- Ambient Background -->
    <div class="fixed inset-0 overflow-hidden pointer-events-none">
        <div class="ambient-orb ambient-orb-purple absolute top-1/4 left-1/4 w-96 h-96"></div>
        <div class="ambient-orb ambient-orb-indigo absolute bottom-1/4 right-1/4 w-96 h-96"></div>
        <div class="ambient-orb ambient-orb-pink absolute top-1/2 left-1/2 w-64 h-64"></div>
    </div>
    
    <!-- Navigation -->
    <nav class="relative z-10 container mx-auto px-6 py-6">
        <div class="flex justify-between items-center">
            <a href="/" class="text-2xl font-bold text-gradient">LinkBio</a>
            <div class="flex items-center gap-4">
                <a href="/auth/login" class="text-gray-400 hover:text-white transition-colors text-sm font-medium">Sign In</a>
                <a href="/auth/register" class="btn-primary px-5 py-2.5 rounded-full text-sm font-medium">Get Started</a>
            </div>
        </div>
    </nav>
    
    <main class="relative z-10">
        <!-- Hero Section -->
        <section class="container mx-auto px-6 pt-20 pb-32">
            <div class="max-w-4xl mx-auto text-center">
                <div class="inline-flex items-center gap-2 px-4 py-2 rounded-full glass text-sm text-gray-400 mb-8" data-aos="fade-down">
                    <span class="w-2 h-2 bg-green-400 rounded-full animate-pulse"></span>
                    Built with Go, HTMX, Alpine.js
                </div>
                
                <h1 class="text-5xl md:text-7xl font-bold text-white mb-6 leading-tight" data-aos="fade-up" data-aos-delay="100">
                    One Link to <span class="text-gradient">Share Everything</span>
                </h1>
                
                <p class="text-xl text-gray-400 mb-10 max-w-2xl mx-auto" data-aos="fade-up" data-aos-delay="200">
                    Create your personalized link-in-bio page in seconds.
                </p>
                
                <div class="flex flex-col sm:flex-row gap-4 justify-center" data-aos="fade-up" data-aos-delay="300">
                    <a href="/auth/register" class="btn-primary px-8 py-4 rounded-full text-lg font-semibold">Create Your LinkBio</a>
                    <a href="#features" class="btn-secondary px-8 py-4 rounded-full text-lg font-medium">See Features</a>
                </div>
            </div>
        </section>
        
        <!-- Features Section -->
        <section id="features" class="container mx-auto px-6 py-24">
            <div class="grid md:grid-cols-3 gap-6 max-w-5xl mx-auto">
                <div class="feature-card glass rounded-2xl p-8" data-aos="fade-up" data-aos-delay="100">
                    <h3 class="text-xl font-semibold text-white mb-3">Unlimited Links</h3>
                    <p class="text-gray-400">Add as many links as you need.</p>
                </div>
                <!-- More feature cards... -->
            </div>
        </section>
        
        <!-- CTA + Footer -->
    </main>
{{end}}
```

> 🧠 **Explain home page structure:**
> 📱 "`{{define \"bodyClass\"}}gradient-bg{{end}}` — base template-ന്റെ default body class override ചെയ്യുന്നു. Dark gradient background set."
> 📱 "Ambient orbs — CSS blur effects. Visual depth create ചെയ്യാൻ. `filter: blur(80px)` use ചെയ്യുന്നു."
> 📱 "`data-aos` attributes — AOS scroll animations. `fade-up`, `fade-down`. `data-aos-delay` stagger effect-ന്."
> 📱 "`text-gradient` class — CSS gradient text effect. Headline-ൽ colorful text."
> 📱 "Feature cards — `glass` class with staggered AOS delays: 100, 200, 300ms. One by one appear ചെയ്യും."
> 📱 "Navigation — Sign In link + Get Started button. `/auth/login`, `/auth/register` routes-ലേക്ക്."

---

## App CSS — Custom Styles

> 📱 "ഇനി app.css — Tailwind-ന് supplement ചെയ്യുന്ന custom CSS."

**⌨️ Key parts of `web/static/css/app.css`:**

```css
:root {
    --gradient-primary: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    --gradient-bg-dark: linear-gradient(135deg, #0f0f1a 0%, #1a1a2e 50%, #16213e 100%);
    --glass-bg: rgba(255, 255, 255, 0.03);
    --glass-border: rgba(255, 255, 255, 0.05);
}

[x-cloak] { display: none !important; }

.gradient-bg { background: var(--gradient-bg-dark); }

.glass {
    background: var(--glass-bg);
    backdrop-filter: blur(20px);
    border: 1px solid var(--glass-border);
}

.text-gradient {
    background: var(--gradient-primary);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
}

.btn-primary {
    background: var(--gradient-primary);
    color: white;
    transition: all 0.3s ease;
}
.btn-primary:hover {
    transform: translateY(-2px);
    box-shadow: 0 10px 40px rgba(102, 126, 234, 0.4);
}

.link-card { animation: cardSlideIn 0.3s ease-out; }

.link-button { transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1); }
.link-button:hover { transform: translateY(-4px) scale(1.01); }

.sortable-ghost { opacity: 0.4; }

.htmx-indicator { opacity: 0; transition: opacity 200ms ease-in; }
.htmx-request .htmx-indicator { opacity: 1; }

.error-message {
    background: rgba(239, 68, 68, 0.1);
    border: 1px solid rgba(239, 68, 68, 0.3);
    color: #fca5a5;
}
```

> 🧠 **Explain key CSS concepts:**
> 📱 "CSS variables — consistent theming. Colors, gradients, shadows ഒരു place-ൽ define."
> 📱 "`[x-cloak]` — Alpine.js flash of unstyled content prevent ചെയ്യുന്നു. Alpine load ആകുന്നതിന് മുമ്പ് hide."
> 📱 "Glass morphism — `backdrop-filter: blur(20px)`. Semi-transparent background with blur. Modern UI effect."
> 📱 "Gradient text trick — `background-clip: text` + `text-fill-color: transparent`. Background gradient text-ൽ show."
> 📱 "HTMX indicator — `opacity: 0` default. `.htmx-request` class add ആകുമ്പോൾ `opacity: 1`. Loading state."
> 📱 "`.sortable-ghost` — SortableJS drag ചെയ്യുമ്പോൾ ghost element style. Opacity reduce."

---

> 🎥 **Transition:** "Base ready. ഇനി individual pages."
