package models

import "time"

type SMSFilter struct {
	Status    string     `form:"status"`
	StartDate *time.Time `form:"start_date" time_format:"2006-01-02T15:04:05Z07:00"`
	EndDate   *time.Time `form:"end_date" time_format:"2006-01-02T15:04:05Z07:00"`
	Page      int        `form:"page"`
	PageSize  int        `form:"page_size"`
}
