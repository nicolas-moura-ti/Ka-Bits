# 🎡 Ka-Bits: Gunslinger of the System

> *"Ka is a wheel; its one purpose is to turn. And your purpose is to mine."*

**Ka-Bits** is a passive (idle) terminal-based game developed in Go. You assume the role of a **System Gunslinger**, mining bits, automating infrastructure, and evolving through the "Data Tower" while facing entropy and bugs in the background.

---

## 🌌 The Narrative

- **The Wheel of Ka:** Bit generation is constant and automatic. Even when you aren't looking, the system continues to turn.
- **The Beam's Rescue:** When you reach a critical processing limit, perform a "Beam Rescue" (Prestige) to reset your hardware for permanent bonuses and code mutations.
- **SRE Meets Stephen King:** System logs mix technical SRE jargon with references to Mid-World and the Crimson King.

---

## ⚙️ Core Mechanics

- **Idle Mining (Goroutines):** The main engine uses tickers to update your bit balance in real-time without locking the interface.
- **BPS (Bits Per Second):**
  - **Hardware (Infra):** Physical items like the *Gilead Terminal* or *Mid-World Server*.
  - **Software (Automation):** Scripts like the *Lobstros Optimizer* that multiply efficiency.
- **Synchronicity Bonus:** Reach level **19** or **99** on any upgrade to activate a **x1.19** global multiplier.
- **Offline Persistence:** When you restart, the game calculates retroactive gains:
  $$Gain = OfflineTime(seconds) \times BPS \times 0.75 (Offline Penalty)$$

---

## 🛠️ Technical Stack & Architecture

- **Language:** Go (Golang)
- **TUI Framework:** [Bubble Tea](https://github.com/charmbracelet/bubbletea) (Charm.sh)
- **Styles:** [Lip Gloss](https://github.com/charmbracelet/lipgloss)
- **Architecture:**
  - `pkg/game`: Core engine, player state, and upgrade registry.
  - `pkg/storage`: JSON persistence with automatic `.bak` backups.
  - `pkg/ui`: Reactive components for Header, Upgrades, and Logs.

---

## 🚀 Installation & Usage

### Prerequisites
- [Go](https://go.dev/doc/install) 1.20+ installed.

### Build and Run
1. Clone the repository.
2. Build the binary:
   ```bash
   go build -o ka-bits main.go
   ```
3. Run the game:
   ```bash
   ./ka-bits
   ```

---

## 🎮 Controls

| Key | Action |
| :--- | :--- |
| `b` | **Mine** bits manually (useful early game). |
| `↑` / `↓` | **Navigate** through available upgrades. |
| `k` / `j` | Alternative navigation (Vim keys). |
| `Enter` / `Space` | **Buy** the selected upgrade. |
| `q` | **Quit** and save your progress. |
| `Ctrl+C` | Force quit (also triggers an emergency save). |

---

## 💾 Storage & Reliability

- **Auto-Save:** The game saves every **30 seconds** in the background.
- **Backup System:** Every save creates a `save.json.bak`. If the main file is missing, the system recovers from the backup automatically.
- **Local Progress:** All data is stored in `save.json` in the project root.

---

## 🧪 Scalability for Developers

The project uses a **Registry Pattern** for upgrades. To add a new item:
1. Open `pkg/game/registry.go`.
2. Add a new `Upgrade` struct to the `InitialUpgrades` slice.
3. The UI will automatically detect and render it on the next launch.

---

## 📜 System Log Examples

- `[INFO] System online. Ka is a wheel.`
- `[INFO] Backup performed on Lud's servers.`
- `[WARN] The Crimson King tried to intercept data packets.`
- `[ERROR] Critical Failure: The world has moved on.`

---

*“Long days and pleasant nights, Gunslinger.”*
