# 🎙️ OBS Studio Recording Guide — Zero Budget Setup

> Everything you need to record a professional-looking coding tutorial with just OBS Studio (free) and a mic. No fancy gear required.

---

## 📥 Step 1: Install OBS Studio

1. Go to **https://obsproject.com** → Download for Windows
2. Run the installer, accept defaults
3. On first launch, choose **"Optimize for Recording"** (NOT streaming)
4. Skip the auto-configuration wizard — we'll set everything manually

---

## ⚙️ Step 2: OBS Settings (One-Time Setup)

### Output Settings
Go to **Settings → Output → Recording**:

| Setting | Value | Why |
|---------|-------|-----|
| Recording Path | `C:\Users\tijo1\Videos\LinkBio-Tutorial` | Keep organized |
| Recording Format | **MKV** (remux to MP4 later) | MKV won't corrupt if OBS crashes |
| Encoder | **x264** (or NVENC if you have Nvidia GPU) | Best compatibility |
| Rate Control | **CRF** | Constant quality |
| CRF Value | **18-20** | Lower = better quality, bigger file |
| Preset | **veryfast** | Good balance of quality vs CPU usage |

> 💡 **Why MKV?** If OBS crashes or you forget to stop recording, MKV files are still recoverable. MP4 files get corrupted. After recording, just remux to MP4 (File → Remux Recordings).

### Video Settings
Go to **Settings → Video**:

| Setting | Value |
|---------|-------|
| Base (Canvas) Resolution | **1920x1080** |
| Output (Scaled) Resolution | **1920x1080** |
| FPS | **30** |

> 🎯 1080p30 is the sweet spot for coding tutorials. 4K is unnecessary — text readability matters more than resolution.

### Audio Settings
Go to **Settings → Audio**:

| Setting | Value |
|---------|-------|
| Sample Rate | **48 kHz** |
| Channels | **Mono** (for single mic) |
| Mic/Auxiliary Audio | Select your mic from dropdown |
| Desktop Audio | **Disabled** (you don't want system sounds) |

---

## 🖥️ Step 3: Create Your Scene

### Scene 1: "Coding" (Main Scene)

This is your primary recording scene. Add these sources in order (bottom to top):

#### Source 1: Display Capture (your entire screen)
1. Click **"+"** in Sources panel → **Display Capture**
2. Name it `Screen`
3. Select your main monitor
4. ✅ Check "Capture Cursor"

> ⚠️ If you have multiple monitors, only capture the one you'll be coding on.

#### Source 2: Audio Input Capture (your mic)
1. Click **"+"** → **Audio Input Capture**
2. Name it `Mic`
3. Select your microphone from the dropdown

#### (Optional) Source 3: Webcam
1. Click **"+"** → **Video Capture Device**
2. Name it `Facecam`
3. Select your webcam
4. Right-click the webcam source → **Transform → Edit Transform**:
   - Position: Bottom-right corner
   - Size: ~250x250 pixels
5. Right-click → **Filters → Add → "Chroma Key"** (if using green screen) or just leave it small

### Scene 2: "Browser" (for demo moments)
1. Create a new scene called `Browser Demo`
2. Add **Window Capture** → Select your browser window
3. Use this when showing the app full-screen

### Switching Scenes
- Go to **Settings → Hotkeys**
- Set `Switch to Scene: Coding` → **Ctrl+F1**
- Set `Switch to Scene: Browser Demo` → **Ctrl+F2**
- Now you can switch views while recording without touching OBS!

---

## 🎤 Step 4: Mic Setup & Audio Quality

### In OBS: Audio Filters (This Makes You Sound Pro!)

Right-click your Mic source → **Filters** → Add these in ORDER:

#### Filter 1: Noise Suppression
- Type: **RNNoise** (AI-based, much better than Speex)
- This removes background noise — fans, AC, traffic

#### Filter 2: Gain
- Boost by **+5 to +10 dB** if your mic is quiet
- Adjust so your voice peaks in the **yellow zone** on the audio mixer (never red)

#### Filter 3: Noise Gate
- Close Threshold: **-40 dB**
- Open Threshold: **-35 dB**
- Attack: **10 ms**
- Hold: **200 ms**
- Release: **100 ms**
- This cuts audio when you're not speaking (kills keyboard sounds between sentences)

#### Filter 4: Compressor
- Ratio: **3:1**
- Threshold: **-20 dB**
- Attack: **5 ms**
- Release: **100 ms**
- This evens out loud/quiet parts of your voice

> 🎯 **Test your filters**: Record a 30-second test clip. Talk normally, whisper, speak loudly. Listen back. Adjust Gain/Noise Gate until it sounds clean.

### Mic Tips (No Expensive Gear Needed)
- **Budget mic**: Even a ₹500 lapel mic plugged into your phone/PC sounds decent with OBS filters
- **Phone as mic**: Use **WO Mic** app (free) to use your phone as a PC microphone via USB/WiFi
- **Positioning**: Keep mic 6-8 inches from mouth, slightly off to the side (reduces plosives — "p" and "b" pops)
- **DIY pop filter**: Stretch a thin cloth/dupatta over a wire hanger loop, place between you and mic

---

## 📐 Step 5: Screen Layout for Coding Tutorials

### Prepare Your Desktop BEFORE Recording

```
┌──────────────────────────────────────────────────┐
│                                                  │
│  VS Code (Left 75%)          │ Browser (Right 25%)│
│  ┌────────────────────────┐  │ ┌──────────────┐  │
│  │                        │  │ │              │  │
│  │   Editor               │  │ │  App Preview │  │
│  │   (Font: 18px)         │  │ │  (zoomed in) │  │
│  │                        │  │ │              │  │
│  ├────────────────────────┤  │ └──────────────┘  │
│  │   Terminal             │  │                   │
│  │   (Font: 16px)         │  │    [Webcam 🔴]    │
│  └────────────────────────┘  │    (optional)     │
│                                                  │
└──────────────────────────────────────────────────┘
```

### Windows Snap Layout
1. Open VS Code → **Win + Left Arrow** (snaps to left half)
2. Drag the edge to take ~75% of screen
3. Open Browser → it fills the remaining 25%
4. Or use **PowerToys FancyZones** for custom layouts

### VS Code Recording Profile
Already in your pre-production checklist, but here's the quick setup:
1. VS Code → **Ctrl+Shift+P** → `Profiles: Create Profile`
2. Name: `Recording`
3. Apply the font/UI settings from `00-pre-production.md`

---

## 🎬 Step 6: Recording Workflow

### Before Each Section

```
1. Open OBS
2. Check mic levels (talk, watch the green/yellow bar)
3. Open your script on phone: tutorial-script/XX-section.md
4. Arrange VS Code + Browser
5. Clear terminal
6. Take a breath
7. Hit Ctrl+F9 (or your record hotkey)
```

### Set Recording Hotkeys
Go to **Settings → Hotkeys**:

| Action | Hotkey | Why |
|--------|--------|-----|
| Start Recording | **Ctrl+F9** | Start without clicking OBS |
| Stop Recording | **Ctrl+F10** | Stop without clicking OBS |
| Pause Recording | **Ctrl+F11** | Pause when you mess up |

> 🔑 **This is crucial**: You should NEVER have to click on OBS during recording. Your screen is being captured — viewers will see you clicking "Stop Recording" if you switch to OBS.

### During Recording
- **Mess up a line?** Just pause (Ctrl+F11), collect yourself, unpause and re-say it. Edit it out later.
- **Need a break?** Pause recording, don't stop it. Keeps it in one file.
- **Silence is OK**: Small pauses between sections feel natural. Don't rush.

### After Each Section
1. Stop recording (**Ctrl+F10**)
2. Go to **File → Remux Recordings** → Convert MKV to MP4
3. Rename the file: `01-opening-hook.mp4`, `02-architecture.mp4`, etc.
4. Watch it back quickly — check audio is clear, text is readable
5. Take a 2-minute break before next section

---

## 🧪 Step 7: Test Recording Checklist

Do a 2-minute test recording before your real session:

- [ ] Record yourself typing code in VS Code
- [ ] Play it back — can you read the code? (If not: increase font or zoom)
- [ ] Listen — is your voice clear? Any echo or background noise?
- [ ] Check file size — 5 min of 1080p30 should be ~200-500MB with CRF 18
- [ ] Check CPU usage — if OBS shows "Encoding overloaded", change preset to `superfast`

---

## ✂️ Step 8: Editing (Free Options)

After recording, you'll want to:
- Cut out mistakes/pauses
- Add section titles
- Maybe add zoom-ins on important code

### Free Editors
| Editor | Platform | Best For |
|--------|----------|----------|
| **DaVinci Resolve** (free) | Windows/Mac | Full editing, color, effects — industry standard |
| **Kdenlive** | Windows/Linux | Simple, fast, lightweight |
| **CapCut Desktop** | Windows/Mac | Quick edits, auto-captions, trendy effects |
| **Shotcut** | Windows | Simple cuts and transitions |

> 🎯 **Recommendation**: Use **DaVinci Resolve** (free version). It's what professionals use, and learning it once helps forever. **CapCut** is the fastest for quick YouTube tutorials with auto-subtitles.

### Minimal Editing Workflow
1. Import your `.mp4` files in order
2. Cut out long pauses and mistakes
3. Add a title card at the beginning (your thumbnail image works)
4. Add chapter markers for YouTube
5. Export at 1080p30, H.264, ~10 Mbps bitrate

---

## 📋 Quick Reference Card

```
╔════════════════════════════════════════╗
║  OBS RECORDING QUICK REFERENCE        ║
╠════════════════════════════════════════╣
║                                       ║
║  Ctrl+F9   → Start Recording          ║
║  Ctrl+F10  → Stop Recording           ║
║  Ctrl+F11  → Pause/Resume             ║
║  Ctrl+F1   → Switch to Coding Scene   ║
║  Ctrl+F2   → Switch to Browser Scene  ║
║                                       ║
║  ⚠️  BEFORE RECORDING:                ║
║  • Mic levels in yellow zone          ║
║  • Notifications OFF                  ║
║  • Script on phone                    ║
║  • Terminal cleared                   ║
║  • Water bottle ready                 ║
║                                       ║
║  🎯  RECORDING SETTINGS:              ║
║  • 1920x1080 @ 30fps                  ║
║  • MKV format (remux to MP4 after)    ║
║  • CRF 18-20, x264 veryfast          ║
║  • RNNoise + Gain + Gate + Compressor ║
║                                       ║
╚════════════════════════════════════════╝
```

---

## 💰 Total Cost

| Item | Cost |
|------|------|
| OBS Studio | **₹0** (free & open source) |
| DaVinci Resolve | **₹0** (free version) |
| Your existing mic/earphone mic | **₹0** |
| WO Mic (phone as mic) | **₹0** |

**Total: ₹0** 🎉

> You don't need expensive gear to make great tutorials. Content quality > production quality. Viewers care about **clear audio** and **readable code** — everything else is bonus.
