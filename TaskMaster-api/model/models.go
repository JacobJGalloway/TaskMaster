package model

  type User struct {
      ID   int64 `json:"id"`
      Name string `json:"name"`
  }

  type StudyTask struct {
      ID               int64 `json:"id"`
      UserID           int64 `json:"user_id"`
      Title          string `json:"title"`
      DurationMinutes  int `json:"duration_minutes"`
      PeriodStart      string `json:"period_start"`
      PeriodEnd        string `json:"period_end"`
      Completed        bool `json:"completed"`
      CompletedAt      *string `json:"completed_at,omitempty"` // pointer since it's nullable
  }