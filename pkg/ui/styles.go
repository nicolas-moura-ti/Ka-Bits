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
			Border(lipgloss.DoubleBorder())

	StyleThinnyFog = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#444444")).
			Italic(true)

	StyleDataRain = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#004400")).
			Faint(true)

	StylePulseBright = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FFFF")).Bold(true)
	StylePulseDim    = lipgloss.NewStyle().Foreground(lipgloss.Color("#008888"))

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

	StylePrestigePrompt = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#000000")).
				Background(lipgloss.Color("#FFD700")).
				Padding(1, 2).
				MarginTop(1).
				MarginBottom(1)

	StyleKaPoints = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFD700")).
			Bold(true)

	StyleTower = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFD700")).
			Bold(true)

	StyleTypeHardware = lipgloss.NewStyle().Foreground(lipgloss.Color("#AAAAAA")).Italic(true)
	StyleTypeSoftware = lipgloss.NewStyle().Foreground(lipgloss.Color("#66CCFF")).Italic(true)
	StyleTypeCosmic   = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF00FF")).Italic(true)

	StyleDescription = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#CCCCCC")).
				Italic(true).
				Faint(true).
				PaddingLeft(4).
				MaxWidth(60)

	glowColors  = []string{"#FFD700", "#FFC400", "#FFB100", "#FF9E00", "#FFB100", "#FFC400"}
	pulseColors = []string{"#00FFFF", "#00DDDD", "#00AAAA", "#008888", "#00AAAA", "#00DDDD"}

	GlowStyles      []lipgloss.Style
	ContainerStyles []lipgloss.Style
)

func init() {
	GlowStyles = make([]lipgloss.Style, len(glowColors))
	for i, c := range glowColors {
		GlowStyles[i] = lipgloss.NewStyle().Foreground(lipgloss.Color(c))
	}

	ContainerStyles = make([]lipgloss.Style, len(pulseColors))
	for i, c := range pulseColors {
		ContainerStyles[i] = StyleContainer.BorderForeground(lipgloss.Color(c))
	}
}

func GetPulseColor(tick int) lipgloss.Color {
	return lipgloss.Color(pulseColors[(tick/4)%len(pulseColors)])
}

func GetGlowColor(tick int) lipgloss.Color {
	return lipgloss.Color(glowColors[(tick/6)%len(glowColors)])
}

func GetGlowStyle(tick int) lipgloss.Style {
	return GlowStyles[(tick/6)%len(GlowStyles)]
}

func GetContainerStyle(tick int) lipgloss.Style {
	return ContainerStyles[(tick/4)%len(ContainerStyles)]
}
