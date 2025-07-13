// main.go
package main

import (
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	matrixGreenBright = lipgloss.Color("#00FF00") // Hacker green
	matrixGreenDim    = lipgloss.Color("#006400") // Dim dark green

	titleStyle = lipgloss.NewStyle().Bold(true).Foreground(matrixGreenBright).MarginBottom(1)
	descStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#444444"))
	glowStyle  = lipgloss.NewStyle().Foreground(matrixGreenBright).Bold(true)
)

type screen int

const (
	screenSplash screen = iota
	screenMenu
	screenAbout
	screenSkills
	screenExperience
	screenContact
)

type tickMsg time.Time

type model struct {
	cursor        int
	choices       []string
	screen        screen
	ticks         int
	loadingDot    int
	rainPositions []int
	width         int
	height        int
}

func initialModel() model {
	width, height := 80, 24

	rainPositions := make([]int, width)
	for i := range rainPositions {
		rainPositions[i] = rand.Intn(height)
	}

	return model{
		choices:       []string{"About Me", "Skills", "Projects", "Contact", "Exit"},
		screen:        screenSplash,
		rainPositions: rainPositions,
		width:         width,
		height:        height,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
			return tickMsg(t)
		}),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.rainPositions = make([]int, m.width)
		for i := range m.rainPositions {
			m.rainPositions[i] = rand.Intn(m.height)
		}
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			if m.screen == screenSplash {
				m.screen = screenMenu
			} else if m.screen == screenMenu {
				switch m.cursor {
				case 0:
					m.screen = screenAbout
				case 1:
					m.screen = screenSkills
				case 2:
					m.screen = screenExperience
				case 3:
					m.screen = screenContact
				case 4:
					return m, tea.Quit
				}
			} else {
				m.screen = screenMenu
			}
		case "up", "k":
			if m.cursor > 0 && m.screen == screenMenu {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 && m.screen == screenMenu {
				m.cursor++
			}
		case " ":
			if m.screen == screenSplash {
				m.screen = screenMenu
			}
		}

	case tickMsg:
		m.ticks++
		m.loadingDot = (m.loadingDot + 1) % 4
		for i := range m.rainPositions {
			m.rainPositions[i]++
			if m.rainPositions[i] > m.height+3 {
				m.rainPositions[i] = 0
			}
		}
		if m.screen == screenSplash && m.ticks > 30 {
			m.screen = screenMenu
		}
		return m, tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
			return tickMsg(t)
		})
	}

	return m, nil
}

func (m model) View() string {
	rain := renderMatrixRain(m)

	var ui string
	switch {
	case m.screen == screenSplash:
		ui = renderSplashScreen(m)
	case m.screen != screenMenu:
		ui = renderScreen(m)
	default:
		ui = renderMenu(m)
	}

	return rain + ui
}

func renderMatrixRain(m model) string {
	var b strings.Builder
	chars := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789@#$%&*")

	for r := 0; r < m.height; r++ {
		for c := 0; c < m.width; c++ {
			pos := m.rainPositions[c]
			if r == pos {
				b.WriteString(lipgloss.NewStyle().Foreground(matrixGreenBright).Render(string(chars[rand.Intn(len(chars))])))
			} else if r < pos && r > pos-3 {
				b.WriteString(lipgloss.NewStyle().Foreground(matrixGreenDim).Render(string(chars[rand.Intn(len(chars))])))
			} else {
				b.WriteString(" ")
			}
		}
		b.WriteString("\n")
	}

	return b.String()
}

func renderSplashScreen(m model) string {
	ascii := `
    ██╗   ██╗ █████╗ ██████╗ ██████╗ ██╗  ██╗ █████╗ ███╗   ███╗ █████╗ ███╗   ██╗
    ██║   ██║██╔══██╗██╔══██╗██╔══██╗██║  ██║██╔══██╗████╗ ████║██╔══██╗████╗  ██║
    ██║   ██║███████║██████╔╝██║  ██║███████║███████║██╔████╔██║███████║██╔██╗ ██║
    ╚██╗ ██╔╝██╔══██║██╔══██╗██║  ██║██╔══██║██╔══██║██║╚██╔╝██║██╔══██║██║╚██╗██║
     ╚████╔╝ ██║  ██║██║  ██║██████╔╝██║  ██║██║  ██║██║ ╚═╝ ██║██║  ██║██║ ╚████║
      ╚═══╝  ╚═╝  ╚═╝╚═╝  ╚═╝╚═════╝ ╚═╝  ╚═╝╚═╝  ╚═╝╚═╝     ╚═╝╚═╝  ╚═╝╚═╝  ╚═══╝

				██╗  ██╗ ██████╗ ████████╗██╗  ██╗ █████╗ ██████╗ ██╗
				██║ ██╔╝██╔═══██╗╚══██╔══╝██║  ██║██╔══██╗██╔══██╗██║
				█████╔╝ ██║   ██║   ██║   ███████║███████║██████╔╝██║
				██╔═██╗ ██║   ██║   ██║   ██╔══██║██╔══██║██╔══██╗██║
				██║  ██╗╚██████╔╝   ██║   ██║  ██║██║  ██║██║  ██║██║
				╚═╝  ╚═╝ ╚═════╝    ╚═╝   ╚═╝  ╚═╝╚═╝  ╚═╝╚═╝  ╚═╝╚═╝`

	coloredASCII := ""
	for _, line := range strings.Split(ascii, "\n") {
		coloredASCII += lipgloss.NewStyle().Foreground(matrixGreenBright).Render(line) + "\n"
	}

	subtitle := glowStyle.Render("Senior Software Engineer | Terminal Portfolio")
	instructions := descStyle.Render("Press ENTER or SPACE to continue...")
	if m.ticks%20 < 10 {
		instructions = glowStyle.Render("Press ENTER or SPACE to continue...")
	}

	return "\n" + coloredASCII + "\n" +
		lipgloss.NewStyle().Width(m.width).Align(lipgloss.Center).Render(subtitle) + "\n\n" +
		lipgloss.NewStyle().Width(m.width).Align(lipgloss.Center).Render(instructions)
}

func renderMenu(m model) string {
	var s strings.Builder
	s.WriteString(titleStyle.Render("==> Welcome to Vardhaman's Terminal Portfolio\n\n"))
	s.WriteString("\n")
	for i, choice := range m.choices {
		cursor := "  "
		if m.cursor == i {
			cursor = ">>"
			s.WriteString(glowStyle.Render(fmt.Sprintf("%s %s", cursor, choice)))
			s.WriteString("\n")
		} else {
			s.WriteString(fmt.Sprintf("%s %s\n", cursor, choice))
		}
	}
	s.WriteString("\n" + descStyle.Render("[↑ ↓ arrows or j/k | Enter to select | q to quit]"))
	return s.String()
}

func renderScreen(m model) string {
	switch m.screen {
	case screenAbout:
		content := titleStyle.Render("==> ABOUT ME") + "\n\n"
		content += "Hi, I'm Vardhaman Kothari — a Senior Software Engineer with over 7 years of backend\n"
		content += "and architecture experience. My strength lies in building scalable, modular systems\n"
		content += "with Java, Spring Boot, Microservices, and Cloud-native tools.\n\n"
		content += "I've contributed to fintech platforms handling ₹1000+ Cr disbursements and love\n"
		content += "mentoring developers and solving real-world business problems with code.\n\n"
		return content + glowStyle.Render("[ Enter to go back ]")

	case screenSkills:
		content := titleStyle.Render("==> SKILLS") + "\n\n"
		content += "- Java, Spring Boot, REST APIs, Microservices\n"
		content += "- PostgreSQL, DynamoDB, MySQL\n"
		content += "- AWS (Fargate, Kinesis, Beanstalk), Docker\n"
		content += "- GitHub Actions, CI/CD, Distributed Systems\n"
		return content + "\n" + glowStyle.Render("[ Enter to go back ]")

	case screenExperience:
		content := titleStyle.Render("📂 PROFESSIONAL EXPERIENCE") + "\n\n"

		content += lipgloss.NewStyle().Bold(true).Render("Senior Software Engineer — Niro (Sep 2021 – Present)") + "\n"
		content += "---------------------------------------------------\n"
		content += "• Architected backend for Embedded Finance platform (₹1000 Cr+ disbursed)\n"
		content += "• Improved system performance by 30%\n"
		content += "• Integrated PayU, LiquiLoans, Muthoot — full lender lifecycle\n"
		content += "• Built CI/CD pipelines using AWS Fargate & GitHub Actions\n"
		content += "• Launched Loan Repayment Microservice — 95% ops automation\n"
		content += "• Mentored juniors, improved code quality and resolved tech debt\n\n"

		content += lipgloss.NewStyle().Bold(true).Render("Senior Software Engineer — Reliance Jio (Jul 2020 – Sep 2021)") + "\n"
		content += "--------------------------------------------------------------\n"
		content += "• Developed smart watchdog for multistage recovery\n"
		content += "• Created containerized microservice-based platform\n\n"

		content += lipgloss.NewStyle().Bold(true).Render("Software Engineer — DigiKredit Finance (Apr 2019 – Jul 2020)") + "\n"
		content += "-------------------------------------------------------------\n"
		content += "• Delivered MSME loan app improving financial access\n"
		content += "• Integrated APIs for DMI, Northern Arc, Paytm, PhonePe\n\n"

		content += lipgloss.NewStyle().Bold(true).Render("Software Engineer — NCR Corporation (Oct 2017 – Apr 2019)") + "\n"
		content += "----------------------------------------------------------\n"
		content += "• Developed dynamic webview suite in IMA application\n"
		content += "• Enhanced batch job management using Java, J2EE\n\n"

		content += lipgloss.NewStyle().Bold(true).Render("Intern — TickerPlant India (Apr 2017 – Oct 2017)") + "\n"
		content += "-----------------------------------------------\n"
		content += "• Learned core development practices and supported engineering teams\n\n"

		content += descStyle.Render("[ Press Enter to go back ]")
		return content

	case screenContact:
		content := titleStyle.Render("==> CONTACT") + "\n\n"
		content += "📧 vardhamank93@gmail.com\n"
		content += "🌐 linkedin.com/in/vardhaman-kothari-598843177/\n"
		content += "📞 +91 8000511720\n"
		return content + "\n" + glowStyle.Render("[ Enter to go back ]")
	default:
		return ""
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	if runtime.GOOS == "windows" {
		enableVirtualTerminalProcessing()
	}
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if err := p.Start(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func enableVirtualTerminalProcessing() {
	// Windows ANSI support placeholder
}
