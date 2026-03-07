# Scene 16: Deploy (58:00 - 59:00)

> 🎬 **Previous**: All templates done (Scene 15)
> 🎯 **Goal**: Build binary, production env vars, security checklist

---

## Build

**⌨️ Type:**
```bash
go build -o linkbio ./cmd/server
```

> 📱 "Single binary. No runtime dependencies. Copy and run!"

---

## Production Environment Variables

```bash
export ENV=production
export SESSION_SECRET=$(openssl rand -hex 32)
export SESSION_ENCRYPTION_KEY=$(openssl rand -hex 16)
export DATABASE_PATH=/data/linkbio.db
```

> 🧠 📱 "openssl rand -hex 32 — random 64-char hex string. Signing key."
> 📱 "openssl rand -hex 16 — random 32-char hex string. Exactly 32 bytes for AES-256!"
> 📱 ".env file production-ൽ use ചെയ്യരുത്. System env vars or secrets manager."

---

## Production Checklist

| Item | Dev | Prod |
|------|-----|------|
| `ENV` | development | production |
| `Secure` cookie | false | **true** (HTTPS only) |
| `SESSION_SECRET` | any string | random, 32+ chars |
| `SESSION_ENCRYPTION_KEY` | any 32-char | random, exactly 32 chars |
| `LOG_LEVEL` | DEBUG | INFO or WARN |

> 📱 "Secure: true — production-ൽ middleware.go-ൽ change ചെയ്യണം. HTTPS ഇല്ലാതെ cookie send ചെയ്യില്ല."

---

> 🎥 **Transition:** "Deploy ready. ഇനി recap."
