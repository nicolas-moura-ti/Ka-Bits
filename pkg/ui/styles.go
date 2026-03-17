package ui

import "github.com/charmbracelet/lipgloss"

var (
	StyleTitle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFD700")).
			Background(lipgloss.Color("#333300")).
			Padding(0, 1).
			MarginBottom(1)

	StyleMoney = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FF00")).
			Bold(true).
			Underline(true)

	StyleMiningPulse = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFFFFF")).
				Background(lipgloss.Color("#00FF00")).
				Bold(true)

	StyleBPS = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888")).
			Italic(true)

	StyleFlow = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00AA00")).
			Bold(true)

	StyleBonusActive = lipgloss.NewStyle().
				Bold(true).
				Italic(true)

	StyleBonusGlow1 = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF00FF"))
	StyleBonusGlow2 = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF66FF"))
	StyleBonusGlow3 = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF99FF"))

	StyleUpgradeHeader = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#00FFFF")).
				Bold(true).
				MarginTop(1).
				MarginBottom(1)

	StyleUpgradeSelected = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#000000")).
				Background(lipgloss.Color("#00FFFF")).
				Bold(true)

	StyleAffordable = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FF00"))

	StyleUnaffordable = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#555555"))

	StyleLogInfo = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#AAAAAA"))

	StyleLogWarn = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFA500"))

	StyleLogErr = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000"))

	StyleContainer = lipgloss.NewStyle().
			Padding(1, 2).
			Border(lipgloss.DoubleBorder()).
			BorderForeground(lipgloss.Color("#00FFFF"))
			
	StyleHelpTip = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFF00")).
			Italic(true).
			Bold(true)

	StyleResetPrompt = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#000000")).
				Background(lipgloss.Color("#FF0000")).
				Padding(1, 2).
				MarginTop(1).
				MarginBottom(1)

	StyleTypeHardware = lipgloss.NewStyle().Foreground(lipgloss.Color("#AAAAAA")).Italic(true)
	StyleTypeSoftware = lipgloss.NewStyle().Foreground(lipgloss.Color("#66CCFF")).Italic(true)
	StyleTypeCosmic   = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF00FF")).Italic(true)

	StyleDescription = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#CCCCCC")).
				Italic(true).
				Faint(true).
				PaddingLeft(4).
				MaxWidth(60)
)
