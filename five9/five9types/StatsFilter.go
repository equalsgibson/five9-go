package five9types

// Information about the statistics filter settings
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

type FilterSettingsSelectionType string

const (
	SelectAll     FilterSettingsSelectionType = "ALL"     // Show statistics for all the queues or agent groups
	SelectMy      FilterSettingsSelectionType = "MY"      // Show statistics for the user's queues or agent groups
	SelectSpecify FilterSettingsSelectionType = "SPECIFY" // Show statistics for the specified queue or agent groups.

)

// Enumeration of time periods to use as statistics filters
type StatisticsRange string

const (
	RangeCurrentDay   StatisticsRange = "CURRENT_DAY"   // Current day
	RangeCurrentMonth StatisticsRange = "CURRENT_MONTH" // Current month
	RangeCurrentShift StatisticsRange = "CURRENT_SHIFT" // Current day shift by the number of hours specified in StatsFilterSettingsInfo
	RangeCurrentWeek  StatisticsRange = "CURRENT_WEEK"  // Current week
	RangeLifetime     StatisticsRange = "LIFETIME"      // All time
	RangeRollingHour  StatisticsRange = "ROLLING_HOUR"  // Last hour divided into five minute intervals
)

// Enumeration with the time period for the list and campaign statistics
type StatisticsRollingPeriod string

const (
	FiveMinutes    StatisticsRollingPeriod = "MINUTES5"
	TenMinutes     StatisticsRollingPeriod = "MINUTES10"
	FifteenMinutes StatisticsRollingPeriod = "MINUTES15"
	ThirtyMinutes  StatisticsRollingPeriod = "MINUTES30"
	OneHour        StatisticsRollingPeriod = "HOUR1"
	TwoHours       StatisticsRollingPeriod = "HOUR2"
	ThreeHours     StatisticsRollingPeriod = "HOUR3"
	Today          StatisticsRollingPeriod = "TODAY"
)
