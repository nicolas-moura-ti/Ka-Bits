package ui

import (
	"fmt"
	"ka-bits/pkg/game"
	"ka-bits/pkg/storage"
	"math/rand"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type tickMsg time.Time
type frameMsg time.Time
type autoSaveMsg time.Time
type randomLogMsg string

type Model struct {
	Engine          *game.Engine
	Cursor          int
	Logs            []string
	ConfirmingReset bool
	AnimationTick   int
	MiningEffect    int
}

func NewModel(engine *game.Engine) Model {
	return Model{
		Engine: engine,
		Cursor: 0,
		Logs:   []string{"[INFO] System online. Ka is a wheel."},
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(tick(), frame(), autoSave(), randomLog())
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		m.Engine.Update(100 * time.Millisecond)
		return m, tick()

	case frameMsg:
		m.AnimationTick++
		if m.MiningEffect > 0 {
			m.MiningEffect--
		}
		return m, frame()

	case autoSaveMsg:
		storage.Save(m.Engine.Player)
		m.addLog("Game auto-saved.", false)
		return m, autoSave()

	case randomLogMsg:
		m.addRawLog(string(msg))
		return m, randomLog()

	case tea.KeyMsg:
		if m.ConfirmingReset {
			switch msg.String() {
			case "y", "Y":
				m.Engine.Player.Reset()
				storage.Save(m.Engine.Player)
				m.addLog("SYSTEM RESET COMPLETE. Ka begins anew.", false)
				m.ConfirmingReset = false
				return m, nil
			case "n", "N", "esc":
				m.ConfirmingReset = false
				return m, nil
			case "ctrl+c", "q":
				return m, tea.Quit
			}
			return m, nil
		}

		switch msg.String() {
		case "ctrl+c", "q":
			storage.Save(m.Engine.Player)
			return m, tea.Quit

		case "up", "k":
			if m.Cursor > 0 {
				m.Cursor--
			}

		case "down", "j":
			if m.Cursor < len(m.Engine.Registry.Order)-1 {
				m.Cursor++
			}

		case "enter", " ":
			upgradeID := m.Engine.Registry.Order[m.Cursor]
			success, logMsg := m.Engine.TryBuyUpgrade(upgradeID)
			m.addLog(logMsg, !success)
			if success {
				storage.Save(m.Engine.Player)
			}

		case "b":
			m.Engine.Player.Bits += 1
			m.MiningEffect = 3 // Lasts 3 frames (150ms)

		case "r":
			m.ConfirmingReset = true
		}
	}

	return m, nil
}

func (m *Model) addLog(msg string, isErr bool) {
	prefix := "[INFO]"
	if isErr {
		prefix = "[WARN]"
	}
	m.addRawLog(fmt.Sprintf("%s %s", prefix, msg))
}

func (m *Model) addRawLog(log string) {
	m.Logs = append(m.Logs, log)
	if len(m.Logs) > 10 {
		m.Logs = m.Logs[1:]
	}
}

func (m Model) View() string {
	s := StyleTitle.Render(" KA-BITS: GUNSLINGER OF THE SYSTEM ") + "\n\n"

	bitsStr := fmt.Sprintf("%.2f Bits", m.Engine.Player.Bits)
	bps := m.Engine.Player.CalculateBPS(m.Engine.Registry)
	bpsStr := fmt.Sprintf("(%.2f BPS)", bps)

	// Mining pulse effect
	moneyView := StyleMoney.Render(bitsStr)
	if m.MiningEffect > 0 {
		moneyView = StyleMiningPulse.Render(" " + bitsStr + " ")
	}

	// Check for synchronicity bonus
	hasBonus := false
	for _, count := range m.Engine.Player.UpgradesOwned {
		if count == 19 || count == 99 {
			hasBonus = true
			break
		}
	}

	bonusStr := ""
	if hasBonus {
		// Glow effect based on ticks
		glowStyles := []tea.Model{nil} // Placeholder
		_ = glowStyles
		var glowStyle = StyleBonusGlow1
		switch (m.AnimationTick / 5) % 3 {
		case 1:
			glowStyle = StyleBonusGlow2
		case 2:
			glowStyle = StyleBonusGlow3
		}
		bonusStr = StyleBonusActive.Inherit(glowStyle).Render(" [Synchronicity x1.19]")
	}

	// Flow animation
	flowChars := []string{"[ - ]", "[ - ]", "[ -- ]", "[ ---]", "[----]", "[--- ]", "[--  ]", "[ -  ]"}
	flowView := ""
	if bps > 0 {
		flowView = " " + StyleFlow.Render(flowChars[m.AnimationTick%len(flowChars)])
	}

	s += moneyView + " " + StyleBPS.Render(bpsStr) + flowView + bonusStr + "\n"
	
	if bps == 0 && m.Engine.Player.Bits < 100 {
		s += StyleHelpTip.Render("-> TIP: Press [b] to mine manually until you can buy a Terminal! <-") + "\n"
	} else {
		s += "\n"
	}

	if m.ConfirmingReset {
		s += StyleResetPrompt.Render("⚠️ WARNING: RESET SYSTEM? ⚠️\nThis will wipe ALL bits and upgrades.\nPress [y] to confirm / [n] to cancel") + "\n"
	}

	s += StyleUpgradeHeader.Render("─── Hardware & Software Upgrades ───") + "\n"
	for i, id := range m.Engine.Registry.Order {
		upgrade, _ := m.Engine.Registry.Get(id)
		owned := m.Engine.Player.UpgradesOwned[id]
		cost := game.CalculateUpgradeCost(upgrade.BaseCost, owned)
		affordable := m.Engine.Player.Bits >= cost

		cursor := "  "
		itemStyle := StyleUnaffordable
		if affordable {
			itemStyle = StyleAffordable
		}

		typeStyle := StyleTypeHardware
		switch upgrade.Type {
		case "Software":
			typeStyle = StyleTypeSoftware
		case "Cosmic":
			typeStyle = StyleTypeCosmic
		}

		if m.Cursor == i {
			cursor = "> "
			itemStyle = StyleUpgradeSelected
			
			// Render the selected item
			itemStr := fmt.Sprintf("%s%-20s Lvl %-2d Cost: %8.2f Bits", cursor, upgrade.Name, owned, cost)
			s += itemStyle.Render(itemStr) + " " + typeStyle.Render("["+upgrade.Type+"]") + "\n"
			s += StyleDescription.Render(upgrade.Description) + "\n"
		} else {
			itemStr := fmt.Sprintf("%s%-20s Lvl %-2d Cost: %8.2f Bits", cursor, upgrade.Name, owned, cost)
			s += itemStyle.Render(itemStr) + "\n"
		}
	}

	s += "\n" + StyleUpgradeHeader.Render("─── System Event Feed ───") + "\n"
	for _, log := range m.Logs {
		logStyle := StyleLogInfo
		if len(log) > 6 {
			if log[1:5] == "WARN" {
				logStyle = StyleLogWarn
			} else if log[1:5] == "ERRO" {
				logStyle = StyleLogErr
			}
		}
		s += logStyle.Render(log) + "\n"
	}

	s += "\n" + StyleBPS.Render("[q] Quit | [b] Mine | [Enter] Buy | [r] Reset | [↑/↓] Navigate")

	return StyleContainer.Render(s)
}

func tick() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func frame() tea.Cmd {
	return tea.Tick(50*time.Millisecond, func(t time.Time) tea.Msg {
		return frameMsg(t)
	})
}

func autoSave() tea.Cmd {
	return tea.Tick(30*time.Second, func(t time.Time) tea.Msg {
		return autoSaveMsg(t)
	})
}

func randomLog() tea.Cmd {
	return tea.Tick(time.Duration(10+rand.Intn(10))*time.Second, func(t time.Time) tea.Msg {
		return randomLogMsg(game.RandomLogs[rand.Intn(len(game.RandomLogs))])
	})
}
