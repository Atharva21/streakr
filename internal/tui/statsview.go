package tui

import (
	"time"

	"github.com/Atharva21/streakr/internal/types"
)

type statsModel struct {
	habitInfo          types.HabitInfo
	heatMap            []bool
	firstDayOfSetMonth time.Time // 1st of set month & year
	streakStart        time.Time // range before which we cannot go
	today              time.Time // range after which we cannot go
}
