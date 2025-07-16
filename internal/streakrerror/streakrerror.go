package streakrerror

type StreakrError struct {
	Err         error
	TerminalMsg string
	ShowUsage   bool
}

func (e *StreakrError) Error() string {
	if e.TerminalMsg != "" {
		return e.TerminalMsg
	}
	if e.Err == nil {
		return "Streakr error"
	}
	return e.Err.Error()
}

func (e *StreakrError) Unwrap() error {
	return e.Err
}
