# Streaker is habit tracking cli

### Database design
Streaker uses a simple SQLite database to store user habits and their streaks. The database schema is designed to be efficient and easy to query.
![dbschema.png](./docs/images/dbschema.png)

For the good habits we store the occurances of those habits in the `streaks` table in form of a range `streak_start` to `streak_end`, gaps in which tell us if the user missed few days.
For bad habits however, we store the "good days" where user did not perform bad habit, and the gaps in this table will tell us about slip-ups.
