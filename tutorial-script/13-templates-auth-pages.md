# Scene 13: Auth Pages — Alpine.js UX (50:00 - 53:00)

> 🎬 **Previous**: Base layout (Scene 12)
> 🎯 **Goal**: Login + Register pages with Alpine.js interactivity

---

## Login Page — Password Toggle & Loading State

> 📱 "Login page-ൽ Alpine.js features: password show/hide, loading state."

**⌨️ Key parts of `web/templates/pages/login.html`:**

```html
<div class="relative z-10 w-full max-w-md" x-data="{ showPassword: false, loading: false }">

    <form hx-post="/auth/login" 
          hx-target="#error-message"
          hx-swap="innerHTML"
          @submit="loading = true"
          @htmx:after-request.window="loading = false">

        <!-- Password with toggle -->
        <div class="relative">
            <input :type="showPassword ? 'text' : 'password'" 
                   name="password" required placeholder="••••••••">
            <button type="button" @click="showPassword = !showPassword">
                <svg x-show="!showPassword"><!-- Eye icon --></svg>
                <svg x-show="showPassword" x-cloak><!-- Eye-slash --></svg>
            </button>
        </div>

        <!-- Submit with loading -->
        <button type="submit" :disabled="loading"
                x-text="loading ? 'Signing in...' : 'Sign In'">
            Sign In
        </button>
    </form>
</div>
```

> 🧠 📱 "x-data — component state. showPassword, loading രണ്ടും false start."
> 📱 ":type binding — showPassword true → type='text', false → type='password'. One line toggle!"
> 📱 "@click — showPassword flip. x-show — matching icon show."
> 📱 "x-cloak — page load flash prevent. CSS: [x-cloak] { display: none }"
> 📱 ":disabled — loading true → button disabled. Double click prevent."
> 📱 "x-text — dynamic button text. Loading state visual feedback."
> 📱 "HTMX + Alpine together! @submit → loading=true. @htmx:after-request → loading=false."

---

## Register Page — Live Preview, Password Strength, Loading

> 📱 "Register page-ൽ മൂന്ന് features: live URL preview, password strength bar, loading."

**⌨️ Key parts of `web/templates/pages/register.html`:**

```html
<div x-data="{ username: '', password: '', showPassword: false, loading: false }">

    <!-- Live URL preview -->
    <input type="text" name="username" x-model="username">
    <div x-show="username.length > 0" x-cloak>
        Your profile: linkbio.com/u/<span x-text="username"></span>
    </div>

    <!-- Password strength -->
    <input :type="showPassword ? 'text' : 'password'" x-model="password">
    <div x-show="password.length > 0" x-cloak>
        <div class="h-1.5 w-full bg-gray-700 rounded-full overflow-hidden">
            <div class="h-full rounded-full transition-all duration-300"
                 :class="password.length < 6 ? 'bg-red-500' : password.length < 10 ? 'bg-yellow-500' : 'bg-green-500'"
                 :style="'width: ' + Math.min(password.length * 10, 100) + '%'">
            </div>
        </div>
        <p :class="password.length < 6 ? 'text-red-400' : password.length < 10 ? 'text-yellow-400' : 'text-green-400'"
           x-text="password.length < 6 ? 'Too short' : password.length < 10 ? 'Fair' : 'Strong'"></p>
    </div>

    <!-- Submit -->
    <button :disabled="loading"
            x-text="loading ? 'Creating account...' : 'Create Account'">
    </button>
</div>
```

> 🧠 📱 "x-model — two-way binding. Input type → variable auto-update → preview real-time."
> 📱 ":class ternary — password length < 6 red, < 10 yellow, else green."
> 📱 ":style width — Math.min(length * 10, 100)%. ഓരോ character 10% grow, max 100%."

> 🎯 📱 "Traffic light! Red = stop (too short). Yellow = caution (fair). Green = go (strong)."

---

## 🎯 Alpine.js Directives Used So Far

| Directive | What | Example |
|-----------|------|---------|
| `x-data` | State | `{ showPassword: false }` |
| `x-model` | Two-way bind | Input ↔ variable |
| `x-show` | Show/hide (CSS) | Eye icon toggle |
| `x-cloak` | No flash | Hide until Alpine loads |
| `x-text` | Dynamic text | Button label |
| `:type` | Dynamic attr | password/text |
| `:class` | Dynamic CSS | Strength bar color |
| `:style` | Dynamic style | Bar width |
| `:disabled` | Dynamic disable | Loading button |
| `@click` | Click handler | Toggle |
| `@submit` | Form submit | Loading state |

> 📱 "11 directives. JavaScript file zero. ഇതാണ് Alpine.js!"

---

> 🎥 **Transition:** "Auth pages done. ഇനി Dashboard."
