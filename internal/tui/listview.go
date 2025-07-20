package tui

import (
	"context"
	"errors"
	"sort"

	"github.com/Atharva21/streakr/internal/service"
	"github.com/Atharva21/streakr/internal/store"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/paginator"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type habitItem struct {
	title string
	desc  string
}

func (h habitItem) Title() string       { return h.title }
func (h habitItem) Description() string { return h.desc }
func (h habitItem) FilterValue() string { return h.title }

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type ListModel struct {
	Ctx         context.Context
	List        list.Model
	Initialized bool
}

type ListLoadedMsg struct {
	List list.Model
}

func (m ListModel) Init() tea.Cmd {
	return func() tea.Msg {
		habits, err := service.ListHabits(m.Ctx)
		if err != nil {
			return viewErrorMsg{err: err}
		}
		sort.Slice(habits, func(i, j int) bool {
			return habits[i].HabitType == store.HabitTypeImprove &&
				habits[j].HabitType != store.HabitTypeImprove
		})
		items := []list.Item{}
		for _, habit := range habits {
			description := ""
			if habit.Description.Valid {
				description = habit.Description.String
			}
			items = append(items, habitItem{
				title: habit.Name,
				desc:  description,
			})
		}

		delegate := list.NewDefaultDelegate()

		delegate.Styles.SelectedTitle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#5d8addff")). // Your desired color
			Bold(true).
			BorderLeft(true).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("#5d8addff")).
			PaddingLeft(1)

		delegate.Styles.SelectedDesc = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#cccccc")).
			PaddingLeft(2)

		listModel := list.New(items, delegate, 80, 15)
		listModel.Title = "My Habits"

		// Title style
		titleStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#5d8addff")).
			Padding(0, 1)
		listModel.Styles.Title = titleStyle

		// Fix pagination styling - use the paginator's style instead
		listModel.Paginator.Type = paginator.Dots
		listModel.Paginator.ActiveDot = "● "
		listModel.Paginator.InactiveDot = "○ "

		// Style the pagination through the paginator styles
		listModel.Styles.PaginationStyle = lipgloss.NewStyle().
			PaddingLeft(2).
			PaddingRight(2)

		listModel.SetShowPagination(true)
		listModel.SetShowStatusBar(false)

		return ListLoadedMsg{List: listModel}
	}
}

func (m ListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case viewErrorMsg:
		return m, tea.Quit
	case tea.KeyMsg:
		// Don't match any of the keys below if we're actively filtering.
		if m.List.FilterState() == list.Filtering {
			break
		}
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		if m.Initialized {
			h, v := docStyle.GetFrameSize()
			m.List.SetSize(msg.Width-h, msg.Height-v)
		}
	case ListLoadedMsg:
		m.Initialized = true
		m.List = msg.List
	}
	var cmd tea.Cmd = nil
	if m.Initialized {
		m.List, cmd = m.List.Update(msg)
	}
	return m, cmd
}

func (m ListModel) View() string {
	if m.Initialized {
		return docStyle.Render(m.List.View())
	}
	return ""
}

func RenderListView(appContext context.Context) error {
	if appContext == nil {
		return errors.New("appContext cannot be nil to render listview")
	}
	p := tea.NewProgram(ListModel{Ctx: appContext}, tea.WithAltScreen())
	go func() {
		<-appContext.Done()
		p.Send(tea.Quit())
	}()
	_, err := p.Run()
	return err
}
