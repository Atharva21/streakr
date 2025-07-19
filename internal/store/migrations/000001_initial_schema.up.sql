CREATE TABLE habits (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL UNIQUE CHECK (length(name) <= 20),
  description TEXT CHECK (description IS NULL OR length(description) <= 200),
  habit_type TEXT CHECK (habit_type IN ('improve', 'quit')) NOT NULL DEFAULT 'improve',
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL
);
CREATE INDEX idx_habits_name ON habits(name);

CREATE TABLE streaks (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  habit_id INTEGER NOT NULL,
  streak_start DATE NOT NULL,
  streak_end DATE NOT NULL,
  FOREIGN KEY (habit_id) REFERENCES habits(id) ON DELETE CASCADE
);
CREATE INDEX idx_streaks_habit_id ON streaks(habit_id);
