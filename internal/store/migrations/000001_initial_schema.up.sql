CREATE TABLE habits (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL,
  description TEXT,
  habit_type TEXT CHECK (habit_type IN ('IMPROVE', 'QUIT')) NOT NULL DEFAULT 'IMPROVE',
  created_at TIMESTAMP,
  last_logged TIMESTAMP
);
CREATE INDEX idx_habits_name ON habits(name);

CREATE TABLE habit_aliases (
  habit_id INTEGER,
  alias TEXT NOT NULL UNIQUE,
  FOREIGN KEY (habit_id) REFERENCES habits(id) ON DELETE CASCADE
);
CREATE INDEX idx_habit_aliases_alias ON habit_aliases(alias);

CREATE TABLE streaks (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  habit_id INTEGER,
  streak_start DATE NOT NULL,
  streak_end DATE NOT NULL,
  FOREIGN KEY (habit_id) REFERENCES habits(id) ON DELETE CASCADE
);
CREATE INDEX idx_streaks_habit_id ON streaks(habit_id);
