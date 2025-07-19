package tui

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Atharva21/streakr/internal/service"
	"github.com/Atharva21/streakr/internal/store"
	"github.com/Atharva21/streakr/internal/store/generated"
	"github.com/Atharva21/streakr/internal/util"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type StatsModel struct {
	Ctx                 context.Context
	Habit               generated.Habit
	TotalStreaksInMonth int
	TotalMissesInMonth  int
	HeatMap             []bool
	FirstDayOfSetMonth  time.Time // 1st of set month & year
	Today               time.Time // range after which we cannot go
	ExitError           error
	HasPreviousNbr      bool
	HasNxtNbr           bool
}

func (m StatsModel) Init() tea.Cmd {
	return func() tea.Msg {
		rangedStats, err := service.GetHabitStatsForRange(
			m.Ctx,
			m.Habit.Name,
			m.FirstDayOfSetMonth,
			m.FirstDayOfSetMonth.AddDate(0, 1, -1),
		)
		if err != nil {
			return viewErrorMsg{err: err}
		}
		return StatsModel{
			Ctx:                 m.Ctx,
			FirstDayOfSetMonth:  m.FirstDayOfSetMonth,
			Today:               m.Today,
			HeatMap:             rangedStats.Heatmap,
			Habit:               rangedStats.Habit,
			ExitError:           nil,
			HasPreviousNbr:      util.AtLeastOneMonthOlder(rangedStats.Habit.CreatedAt, m.FirstDayOfSetMonth),
			HasNxtNbr:           util.AtLeastOneMonthOlder(m.FirstDayOfSetMonth, m.Today),
			TotalStreaksInMonth: rangedStats.TotalStreakDaysInRange,
			TotalMissesInMonth:  rangedStats.TotalMissesInRange,
		}
	}
}

type neighborMonth = int

const (
	previousMonth neighborMonth = -1
	nextMonth     neighborMonth = 1
)

func getNeighbourMonthStatsCmd(m StatsModel, nbrType neighborMonth) tea.Cmd {
	firstDayOfNbrMonth := m.FirstDayOfSetMonth.AddDate(0, int(nbrType), 0)
	lastDayOfNbrMonth := firstDayOfNbrMonth.AddDate(0, 1, -1)
	habitStart := m.Habit.CreatedAt
	today := m.Today
	if firstDayOfNbrMonth.Year() < habitStart.Year() {
		return nil
	}
	if firstDayOfNbrMonth.Year() == habitStart.Year() {
		if firstDayOfNbrMonth.Month() < habitStart.Month() {
			return nil
		}
	}
	if firstDayOfNbrMonth.Year() > today.Year() {
		return nil
	}
	if firstDayOfNbrMonth.Year() == today.Year() {
		if firstDayOfNbrMonth.Month() > today.Month() {
			return nil
		}
	}
	return func() tea.Msg {
		rangedStats, err := service.GetHabitStatsForRange(m.Ctx, m.Habit.Name, firstDayOfNbrMonth, lastDayOfNbrMonth)
		if err != nil {
			return viewErrorMsg{
				err: err,
			}
		}
		if len(rangedStats.Heatmap) != lastDayOfNbrMonth.Day() {
			return viewErrorMsg{
				err: errors.New("Cannot render incomplete heatmap length mismatch"),
			}
		}
		sm := StatsModel{
			Ctx:                 m.Ctx,
			Habit:               rangedStats.Habit,
			TotalStreaksInMonth: rangedStats.TotalStreakDaysInRange,
			TotalMissesInMonth:  rangedStats.TotalMissesInRange,
			HeatMap:             rangedStats.Heatmap,
			FirstDayOfSetMonth:  firstDayOfNbrMonth,
			Today:               m.Today,
			ExitError:           nil,
			HasPreviousNbr:      util.AtLeastOneMonthOlder(m.Habit.CreatedAt, firstDayOfNbrMonth),
			HasNxtNbr:           util.AtLeastOneMonthOlder(firstDayOfNbrMonth, today),
		}
		return sm
	}
}

func (m StatsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case viewErrorMsg:
		return m, tea.Quit
	case contextCancelledMsg:
		return m, tea.Quit
	case StatsModel:
		return msg, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "left", "h":
			return m, getNeighbourMonthStatsCmd(m, previousMonth)
		case "right", "l":
			return m, getNeighbourMonthStatsCmd(m, nextMonth)
		}
	}
	return m, nil
}

func (m StatsModel) View() string {
	monthTitleColor := lipgloss.Color("#5d8addff")
	weekDayHeaderColor := lipgloss.Color("#5d8addff")
	todaysDateBGColor := lipgloss.Color("#5d8addff")
	weekdayStyle := lipgloss.
		NewStyle().
		Align(lipgloss.Left).
		Foreground(weekDayHeaderColor)
	streakColor := lipgloss.NewStyle().Foreground(lipgloss.Color("#25a425ff"))
	missColor := lipgloss.NewStyle().Foreground(lipgloss.Color("#c25252ff"))
	if m.Habit.HabitType == store.HabitTypeImprove {
		missColor = lipgloss.NewStyle()
	}
	futureDatesColor := lipgloss.NewStyle().Foreground(lipgloss.Color("#444444"))
	weekDaysHeader := "Mon Tue Wed Thu Fri Sat Sun"
	monthTitle := ""
	if m.HasPreviousNbr {
		monthTitle += "← "
	}
	monthTitle += fmt.Sprintf("%s %d", m.FirstDayOfSetMonth.Month(), m.FirstDayOfSetMonth.Year())
	if m.HasNxtNbr {
		monthTitle += " →"
	}
	calView := ""
	calView += lipgloss.
		NewStyle().
		Width(len(weekDaysHeader)).
		Align(lipgloss.Center).
		Foreground(monthTitleColor).
		Render(monthTitle) + "\n"
	calView += weekdayStyle.Width(len(weekDaysHeader)).Render(weekDaysHeader)
	calView += "\n"
	dayOfTheWeek := m.FirstDayOfSetMonth.Weekday()
	if dayOfTheWeek == 0 {
		dayOfTheWeek = 7
	}
	for ; dayOfTheWeek > 1; dayOfTheWeek-- {
		calView += "    "
	}
	weeksPassed := 0
	for i, val := range m.HeatMap {
		date := m.FirstDayOfSetMonth.AddDate(0, 0, i)
		if date.Day() < 10 {
			calView += "  "
		} else {
			calView += " "
		}
		style := streakColor
		if !val {
			style = missColor
		}
		if util.IsSameDate(date, time.Now()) {
			style = style.Background(todaysDateBGColor)
		}
		if util.CompareDate(date, time.Now()) == -1 {
			style = futureDatesColor
		}
		calView += style.Render(fmt.Sprintf("%d", date.Day()))
		if i != len(m.HeatMap)-1 {
			if date.Weekday() == time.Sunday {
				calView += "\n"
				weeksPassed++
			} else {
				calView += " "
			}
		}
	}
	for ; weeksPassed < 6; weeksPassed++ {
		calView += "\n"
	}

	calView += fmt.Sprintf("Completed: %d\n", m.TotalStreaksInMonth)
	calView += fmt.Sprintf("Missed: %d\n", m.TotalMissesInMonth)

	return calView
}

func RenderStatsView(appContext context.Context, year, month int, habit generated.Habit) error {
	date := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
	sm := &StatsModel{
		Ctx:                appContext,
		FirstDayOfSetMonth: date,
		Today:              time.Now(),
		ExitError:          nil,
		Habit:              habit,
	}
	if sm.Ctx == nil {
		return errors.New("Context cannot be empty for statsView")
	}
	statsViewProgram := tea.NewProgram(*sm)
	go func() {
		<-sm.Ctx.Done()
		statsViewProgram.Send(contextCancelledMsg{})
	}()
	_, err := statsViewProgram.Run()
	return err
}
