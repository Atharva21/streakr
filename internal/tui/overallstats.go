package tui

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sort"

	"github.com/Atharva21/streakr/internal/service"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type statsLoadedMsg struct {
	t table.Model
}

type OverallStats struct {
	Ctx context.Context
	t   table.Model
}

func (m OverallStats) Init() tea.Cmd {
	return func() tea.Msg {
		s, err := service.GetOverallStats(m.Ctx)
		if err != nil {
			slog.Error("error in getting overall stats from service", "err", err.Error())
			return viewErrorMsg{err: err}
		}
		if s == nil {
			return viewErrorMsg{err: errors.New("could not fetch overall stats")}
		}
		sort.Slice(s.HabitInfos, func(i, j int) bool {
			return s.HabitInfos[i].CurrentStreak > s.HabitInfos[j].CurrentStreak
		})
		slog.Info("got overall stats", "len", len(s.HabitInfos))
		cols := []table.Column{
			{Title: "Habit", Width: 22},
			{Title: "Current", Width: 9},
			{Title: "Max", Width: 9},
			{Title: "Total", Width: 6},
			{Title: "Missed", Width: 6},
		}
		rows := []table.Row{}
		for _, habitInfo := range s.HabitInfos {
			currentStreakStr := fmt.Sprintf("%d", habitInfo.CurrentStreak)
			maxStreakStr := fmt.Sprintf("%d", habitInfo.MaxStreak)
			if habitInfo.CurrentStreak == habitInfo.MaxStreak && habitInfo.CurrentStreak != 0 {
				currentStreakStr += "⚡"
				maxStreakStr += "⚡"
			}
			rows = append(rows, table.Row{
				habitInfo.Habit.Name,
				currentStreakStr,
				maxStreakStr,
				fmt.Sprintf("%d", habitInfo.TotalPerformedDays),
				fmt.Sprintf("%d", habitInfo.TotalMissedDays),
			})
		}
		t := table.New(
			table.WithColumns(cols),
			table.WithRows(rows),
			table.WithFocused(true),
			table.WithHeight(7),
		)
		style := table.DefaultStyles()
		style.Header = style.Header.
			BorderStyle(lipgloss.NormalBorder()).
			BorderBottom(true).
			Bold(false)
		style.Selected = style.Selected.
			Background(lipgloss.Color("#5d8addff")).
			Bold(false)
		t.SetStyles(style)
		return statsLoadedMsg{
			t: t,
		}
	}
}

func (m OverallStats) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case viewErrorMsg:
		return m, tea.Quit
	case statsLoadedMsg:
		m.t = msg.t
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	var cmd tea.Cmd
	m.t, cmd = m.t.Update(msg)
	return m, cmd
}

func (m OverallStats) View() string {
	return lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).Render(m.t.View()) + "\n"
}

func RenderOverallStats(appContext context.Context) error {
	if appContext == nil {
		return errors.New("Context cannot be nil")
	}
	overallStats := OverallStats{
		Ctx: appContext,
	}
	p := tea.NewProgram(overallStats)
	go func() {
		<-appContext.Done()
		slog.Error("app context closed, closing overall stats view")
		p.Send(tea.Quit())
	}()
	slog.Info("running tea program for overall stats")
	_, err := p.Run()
	if err != nil {
		slog.Error("error in overall stats view", "err", err.Error())
		return err
	}
	return nil
}
