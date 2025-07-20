package tui

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sort"
	"time"

	"github.com/Atharva21/streakr/internal/service"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type statsLoadedMsg struct {
	table table.Model
}

type OverallStats struct {
	Ctx   context.Context
	table table.Model
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
			Foreground(lipgloss.Color("15")).
			Bold(false)
		t.SetStyles(style)
		return statsLoadedMsg{
			table: t,
		}
	}
}

func (m OverallStats) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case viewErrorMsg:
		return m, tea.Quit
	case statsLoadedMsg:
		m.table = msg.table
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		case "enter":
			habitName := m.table.SelectedRow()[0]
			habit, err := service.GetHabitByName(m.Ctx, habitName)
			if err != nil {
				slog.Error("error in getting habit by name in statsview", "err", err.Error())
				return m, tea.Quit
			}
			now := time.Now()
			sm := &StatsModel{
				Ctx:                m.Ctx,
				FirstDayOfSetMonth: time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local),
				Today:              now,
				Habit:              habit,
				ParentTable:        &m.table,
			}
			return sm, sm.Init()
		}
	}
	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m OverallStats) View() string {
	helpStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#666666"))
	helpMsg := helpStyle.Render("↑↓ navigate • enter select • q/esc quit")
	return lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).Render(m.table.View()) + "\n" + helpMsg
}

func RenderOverallStats(appContext context.Context) error {
	if appContext == nil {
		return errors.New("Context cannot be nil")
	}
	overallStats := OverallStats{
		Ctx: appContext,
	}
	p := tea.NewProgram(overallStats, tea.WithAltScreen())
	go func() {
		<-appContext.Done()
		slog.Error("app context closed, closing overall stats view")
		p.Send(tea.Quit())
	}()
	_, err := p.Run()
	if err != nil {
		slog.Error("error in overall stats view", "err", err.Error())
		return err
	}
	return nil
}
