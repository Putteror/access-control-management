package model

type AccessControlGroupSchedule struct {
	BaseModel
	AccessControlGroupID string  `json:"access_control_group_id"`
	DayOfWeek            int     `json:"day_of_week"`
	Date                 *string `json:"date"`
	StartTime            string  `json:"start_time"`
	EndTime              string  `json:"end_time"`
}
