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

> 🎥 **Transition:** "Base ready. ഇനി individual pages."
