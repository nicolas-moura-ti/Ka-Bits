# 🛠️ Development & Scalability Guide

This guide explains the internal architecture of **Ka-Bits** and how to extend it with new features.

---

## 🏗️ Project Structure

```text
.
├── main.go                # Application entry point & TUI runner
├── pkg/
│   ├── game/              # Core logic & Engine
│   │   ├── player.go      # Player state (Resources, Upgrades)
│   │   ├── engine.go      # Passive income & Purchase logic
│   │   ├── registry.go    # Registry of all available upgrades
│   │   └── upgrades.go    # Data structures for items
│   ├── storage/           # Persistence
│   │   └── save.go        # JSON & Backup handling
│   └── ui/                # Terminal Interface
│       ├── styles.go      # Lip Gloss styling
│       └── view.go        # Bubble Tea Model & View logic
└── save.json              # Local save file (Auto-generated)
```

---

## ➕ Adding New Upgrades

The game uses a **Registry** to decouple data from logic. To add a new upgrade:

1. Open `pkg/game/registry.go`.
2. Locate the `InitialUpgrades` slice.
3. Add a new `Upgrade` object:
   ```go
   {
       ID:          "quantum_rose",
       Name:        "Quantum Rose",
       Description: "The focal point of existence, blooming in binary.",
       BaseCost:    2000,
       BaseBPS:     100,
       Type:        "cosmic",
   },
   ```
4. Rebuild the project (`go build`). The item will appear in the UI automatically.

---

## 📈 New Resources (Prestige / Ka-Points)

To implement a prestige system or new currencies:

1. **Player Struct:** Add a new field in `pkg/game/player.go` (e.g., `PrestigePoints float64`).
2. **Logic:** In `pkg/game/engine.go`, create a `Reset()` function that clears bits/upgrades and calculates the new currency.
3. **UI:** Update the `View()` in `pkg/ui/view.go` to display the new resource.

---

## ⏳ Persistence & Backups

The `pkg/storage/save.go` module handles local saving.
- **Save Interval:** Managed in `pkg/ui/view.go` via the `autoSave()` command (default: 30s).
- **Versioning:** The `Player.Version` field allows you to write migration logic if you change the JSON structure in the future.
- **Corruption Protection:** If `save.json` is corrupted, the system attempts to load `save.json.bak`.

---

## 🧪 Testing

Run tests to ensure core logic remains intact after refactoring:
```bash
go test ./...
```

The tests cover:
- Bit accumulation (BPS).
- Purchase logic (validating costs and balances).
- JSON save/load cycles.

---

*“May you find your tower.”*
