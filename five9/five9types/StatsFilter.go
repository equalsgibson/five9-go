package five9types

// Information about the statistics filter settings.
type StatsFilterSettingsInfo struct {
	Groups               []string                    `json:"groups"`               // Array of agent groups. You set this when groupsSelectionType is set to SPECIFY
	GroupSelectionType   FilterSettingsSelectionType `json:"groupSelectionType"`   //
	Range                StatisticsRange             `json:"range"`                //
	RollingTimePeriod    StatisticsRollingPeriod     `json:"rollingTimePeriod"`    //
	ShiftHours           uint64                      `json:"shiftHours"`           // Number of hours to shift when range is set to CURRENT_SHIFT.
	Skills               []QueueID                   `json:"skills"`               // Array of queues available when skillsSelectionType is set to SPECIFY
	SkillsSelectionType  FilterSettingsSelectionType `json:"skillsSelectionType"`  //
	SubscribedHourOffset uint64                      `json:"subscribedHourOffset"` //
	TimeZone             *string                     `json:"timeZone"`             // Domain’s time zone. Default value is null
	TimeZoneID           *string                     `json:"timeZoneID"`           //
	UseAdminTimeZone     bool                        `json:"useAdminTimeZone"`     //
}
