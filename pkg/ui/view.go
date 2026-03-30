package ui

import (
	"encoding/json"
	"fmt"
	"ka-bits/pkg/game"
	"ka-bits/pkg/storage"
	"math/rand"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const DataRainChars = "01#!@$*&%"

type tickMsg time.Time
type frameMsg time.Time
type autoSaveMsg time.Time
type randomLogMsg string
type saveResultMsg struct{ err error }

type Model struct {
	Engine             *game.Engine
	Cursor             int
	Logs               []string
	ConfirmingReset    bool
	ConfirmingPrestige bool
	AnimationTick      int
	MiningEffect       int
	DataRain           []string
}

func NewModel(engine *game.Engine) Model {
	rain := make([]string, 15)
	for i := range rain {
		rain[i] = " "
	}
	return Model{
		Engine:   engine,
		Cursor:   0,
		Logs:     []string{"[INFO] System online. Ka is a wheel."},
		DataRain: rain,
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
		// Update Data Rain
		if m.AnimationTick%2 == 0 {
			copy(m.DataRain[1:], m.DataRain[:len(m.DataRain)-1])
			m.DataRain[0] = string(DataRainChars[rand.Intn(len(DataRainChars))])
			if rand.Intn(3) == 0 {
				m.DataRain[0] = " "
			}
		}
		return m, frame()

	case autoSaveMsg:
		cmd := saveCmd(m.Engine.Player)
		m.addLog("Game auto-saved.", false)
		return m, tea.Batch(cmd, autoSave())

	case randomLogMsg:
		m.addRawLog(string(msg))
		return m, randomLog()

	case saveResultMsg:
		if msg.err != nil {
			m.addLog(fmt.Sprintf("Error saving game: %v", msg.err), true)
		}
		return m, nil

	case tea.KeyMsg:
		return m.handleKey(msg)
	}

	return m, nil
}

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.ConfirmingReset {
		switch msg.String() {
		case "y", "Y":
			m.Engine.Player.Reset()
			cmd := saveCmd(m.Engine.Player)
			m.addLog("SYSTEM RESET COMPLETE. Ka begins anew.", false)
			m.ConfirmingReset = false
			return m, cmd
		case "n", "N", "esc":
			m.ConfirmingReset = false
			return m, nil
		}
		return m, nil
	}

	if m.ConfirmingPrestige {
		switch msg.String() {
		case "y", "Y":
			gain := m.Engine.CalculatePrestigeGain()
			var cmd tea.Cmd
			if gain > 0 {
				m.Engine.Player.BeamRescue(gain)
				cmd = saveCmd(m.Engine.Player)
				m.addLog(fmt.Sprintf("BEAM RESCUE COMPLETE. Gained %d Ka-Points.", gain), false)
			}
			m.ConfirmingPrestige = false
			return m, cmd
		case "n", "N", "esc":
			m.ConfirmingPrestige = false
			return m, nil
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
			return m, saveCmd(m.Engine.Player)
		}

	case "b":
		m.Engine.Player.Bits += 1
		m.Engine.Player.TotalBitsEver += 1
		m.MiningEffect = 3

	case "r":
		m.ConfirmingReset = true
	case "p":
		if m.Engine.CalculatePrestigeGain() > 0 {
			m.ConfirmingPrestige = true
		} else {
			m.addLog("Not enough total bits for Beam Rescue (Min: 500k).", true)
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
	if len(m.Logs) > 8 {
		m.Logs = m.Logs[1:]
	}
}

func (m Model) View() string {
	rainPanel := m.renderDataRain()
	leftPanel := m.renderGamePanel()
	rightPanel := m.renderTowerPanel()

	mainView := lipgloss.JoinHorizontal(lipgloss.Top, rainPanel, leftPanel, rightPanel)

	// Resolução do Conflito: Mantendo a refatoração da branch main.
	return GetContainerStyle(m.AnimationTick).Render(mainView)
}

func (m Model) renderDataRain() string {
	return renderDataRain(m.DataRain)
}

func renderDataRain(dataRain []string) string {
	var builder strings.Builder
	builder.Grow(len(dataRain) * 2)
	for _, char := range dataRain {
		builder.WriteString(char)
		builder.WriteByte('\n')
	}
	return StyleDataRain.Render(builder.String())
}

func (m Model) renderHeader(builder *strings.Builder, bps float64) {
	builder.WriteString(StyleTitle.Render(" KA-BITS: GUNSLINGER OF THE SYSTEM "))
	builder.WriteString("\n\n")

	bitsStr := fmt.Sprintf("%.2f Bits", m.Engine.Player.Bits)
	bpsStr := fmt.Sprintf("(%.2f BPS)", bps)

	moneyView := StyleMoney.Render(bitsStr)
	if m.MiningEffect > 0 {
		moneyView = StyleMiningPulse.Render(" " + bitsStr + " ")
	}

	hasBonus := false
	for _, count := range m.Engine.Player.UpgradesOwned {
		if count == 19 || count == 99 {
			hasBonus = true
			break
		}
	}

	bonusStr := ""
	if hasBonus {
		var glowStyle = StyleBonusGlow1
		switch (m.AnimationTick / 5) % 3 {
		case 1:
			glowStyle = StyleBonusGlow2
		case 2:
			glowStyle = StyleBonusGlow3
		}
		bonusStr = StyleBonusActive.Inherit(glowStyle).Render(" [Synchronicity x1.19]")
	}

	flowChars := []string{"[ - ]", "[ - ]", "[ -- ]", "[ ---]", "[----]", "[--- ]", "[--  ]", "[ -  ]"}
	flowView := ""
	if bps > 0 {
		flowView = " " + StyleFlow.Render(flowChars[m.AnimationTick%len(flowChars)])
	}

	kpView := ""
	if m.Engine.Player.KaPoints > 0 {
		kpView = " " + StyleKaPoints.Render(fmt.Sprintf("| %d Ka-Points", m.Engine.Player.KaPoints))
	}

	builder.WriteString(moneyView)
	builder.WriteString(" ")
	builder.WriteString(StyleBPS.Render(bpsStr))
	builder.WriteString(flowView)
	builder.WriteString(bonusStr)
	builder.WriteString(kpView)
	builder.WriteString("\n")

	if bps == 0 && m.Engine.Player.Bits < 100 {
		builder.WriteString(StyleHelpTip.Render("-> TIP: Press [b] to mine manually until you can buy a Terminal! <-"))
		builder.WriteString("\n")
	} else {
		builder.WriteString("\n")
	}
}

func (m Model) renderPrompts(builder *strings.Builder) {
	if m.ConfirmingReset {
		builder.WriteString(StyleResetPrompt.Render("⚠️ WARNING: RESET SYSTEM? ⚠️\nThis will wipe ALL bits and upgrades.\nPress [y] to confirm / [n] to cancel"))
		builder.WriteString("\n")
	}

	if m.ConfirmingPrestige {
		gain := m.Engine.CalculatePrestigeGain()
		builder.WriteString(StylePrestigePrompt.Render(fmt.Sprintf("✨ BEAM RESCUE ✨\nReset current run to gain %d Ka-Points?\n(+%d%% Permanent BPS Bonus)\nPress [y] to confirm / [n] to cancel", gain, gain*5)))
		builder.WriteString("\n")
	}
}

func (m Model) renderUpgrades(builder *strings.Builder) {
	builder.WriteString(StyleUpgradeHeader.Render("─── Hardware & Software Upgrades ───"))
	builder.WriteString("\n")

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
			// Pulsing cursor
			cursor = "> "
			if (m.AnimationTick/5)%2 == 0 {
				cursor = "» "
			}
			itemStyle = StyleUpgradeSelected

			itemStr := fmt.Sprintf("%s%-20s Lvl %-2d Cost: %8.2f Bits", cursor, upgrade.Name, owned, cost)
			builder.WriteString(itemStyle.Render(itemStr))
			builder.WriteString(" ")
			builder.WriteString(typeStyle.Render("[" + upgrade.Type + "]"))
			builder.WriteString("\n")
			builder.WriteString(StyleDescription.Render(upgrade.Description))
			builder.WriteString("\n")
		} else {
			itemStr := fmt.Sprintf("%s%-20s Lvl %-2d Cost: %8.2f Bits", cursor, upgrade.Name, owned, cost)
			builder.WriteString(itemStyle.Render(itemStr))
			builder.WriteString("\n")
		}
	}
}

func (m Model) renderLogs(builder *strings.Builder) {
	builder.WriteString("\n")
	builder.WriteString(StyleUpgradeHeader.Render("─── System Event Feed ───"))
	builder.WriteString("\n")
	for _, log := range m.Logs {
		logStyle := StyleLogInfo
		if len(log) > 6 {
			if log[1:5] == "WARN" {
				logStyle = StyleLogWarn
			} else if log[1:5] == "ERRO" {
				logStyle = StyleLogErr
			}
		}
		builder.WriteString(logStyle.Render(log))
		builder.WriteString("\n")
	}

	builder.WriteString("\n")
	builder.WriteString(StyleBPS.Render("[q] Quit | [b] Mine | [Enter] Buy | [p] Prestige | [r] Reset"))
}

func (m Model) renderGamePanel() string {
	bps := m.Engine.Player.CalculateBPS(m.Engine.Registry)

	var builder strings.Builder
	builder.Grow(2048)

	m.renderHeader(&builder, bps)
	m.renderPrompts(&builder)
	m.renderUpgrades(&builder)
	m.renderLogs(&builder)

	return builder.String()
}

func (m Model) renderTowerPanel() string {
	totalUpgrades := m.Engine.Player.GetTotalUpgrades(m.Engine.Registry)

	towerLevels := []string{
		"        |        ",
		"      [ ]      ",
		"     [   ]     ",
		"    [     ]    ",
		"   [       ]   ",
		"  [         ]  ",
		" [            ] ",
		"[              ]",
		"===============",
	}

	towerArt := ""

	// Shifting fog (Thinny)
	fogChars := []string{"~  ~  ~", " ~  ~  ", "  ~  ~ ", "   ~  ~"}
	fog := StyleThinnyFog.Render(fogChars[(m.AnimationTick/8)%len(fogChars)])

	eyeChar := " "
	if totalUpgrades > 0 {
		eyeChar = "O"
	}
	if (m.AnimationTick/10)%2 == 0 && totalUpgrades > 0 {
		eyeChar = "*"
	}

	towerArt += "\n" + fog + "\n"
	towerArt += fmt.Sprintf("        %s        \n", GetGlowStyle(m.AnimationTick).Render(eyeChar))
	towerArt += "        |        \n"

	limit := (totalUpgrades / 2) + 1
	if limit > len(towerLevels) {
		limit = len(towerLevels)
	}

	for i := 0; i < limit; i++ {
		towerArt += " " + StyleTower.Render(towerLevels[i]) + "\n"
	}

	towerArt += "\n" + fog

	return lipgloss.NewStyle().
		PaddingLeft(10).
		Render(towerArt)
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

func saveCmd(p *game.Player) tea.Cmd {
	// Marshal synchronously to capture state without data races
	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return func() tea.Msg { return saveResultMsg{err} }
	}

	return func() tea.Msg {
		// Write to disk asynchronously
		writeErr := storage.WriteData(data)
		return saveResultMsg{writeErr}
	}
}

func randomLog() tea.Cmd {
	return tea.Tick(time.Duration(10+rand.Intn(10))*time.Second, func(t time.Time) tea.Msg {
		return randomLogMsg(game.RandomLogs[rand.Intn(len(game.RandomLogs))])
	})
}
