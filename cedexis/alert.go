package cedexis

import (
	"fmt"
	"strings"
)

const alertsConfigPath = "/config/alerts.json"

// Alert represents a configured alert.
type Alert struct {
	ID                *int      `json:"id,omitempty"`
	Name              *string   `json:"name,omitempty"`
	Type              *string   `json:"type,omitempty"`
	Timing            *string   `json:"timing,omitempty"`
	Threshold         *int      `json:"threshold,omitempty"`
	Sensitivity       *int      `json:"sensitivity,omitempty"`
	ProbeType         *int      `json:"probeType,omitempty"`
	Version           *int      `json:"version,omitempty"`
	Debounce          *int      `json:"debounce,omitempty"`
	CompareOperator   *string   `json:"compareOperator,omitempty"`
	Enabled           *bool     `json:"enabled,omitempty"`
	DisabledMinutes   *int      `json:"disabledMinutes,omitempty"`
	Locations         *[]string `json:"locations,omitempty"`
	Platform          *int      `json:"platform,omitempty"`
	Emails            *[]string `json:"emails,omitempty"`
	Peers             *[]int    `json:"peers,omitempty"`
	EventsLast24Hours *int      `json:"eventsLast24Hours,omitempty"`
	CountryEvents     *[]int    `json:"countryEvents,omitempty"`
	ASNEvents         *[]int    `json:"asnEvents,omitempty"`
	RefererURI        *string   `json:"refererUri,omitempty"`
	Statistic         *string   `json:"statistic,omitempty"`
	DataSource        *string   `json:"dataSource,omitempty"`
	AutoFill          *int      `json:"autoFill,omitempty"`
	NotifyChange      *string   `json:"notifyChange,omitempty"`
}

// AlertType is the type of alert.
type AlertType int

const (
	// AlertTypeSonar indicates the alert is triggered by Sonar
	AlertTypeSonar AlertType = iota

	// AlertTypeRadar indicates the alert is triggered by Radar
	AlertTypeRadar
)

// AlertChange is an enum of which events trigger the alert.
type AlertChange int

const (
	// AlertChangeAny alerts on up and down events.
	AlertChangeAny AlertChange = iota

	// AlertChangeToUp alerts on up events.
	AlertChangeToUp

	// AlertChangeToDown alerts on down events.
	AlertChangeToDown
)

// AlertTiming is an enum of when to trigger the alert.
type AlertTiming int

const (
	// AlertTimingImmediate triggers the alert immediately.
	AlertTimingImmediate AlertTiming = iota

	// AlertTimingSummary triggers the alert as a daily summary.
	AlertTimingSummary

	// AlertTimingBoth triggers the alert immediately and provides daily summary.
	AlertTimingBoth
)

func (at AlertType) String() string {
	switch at {
	case AlertTypeSonar:
		return "sonar"
	case AlertTypeRadar:
		return "radar"
	default:
		return fmt.Sprintf("<unknown %d>", int(at))
	}
}

// ParseAlertType parses an alert type 'sonar' or 'radar' to enum value
func ParseAlertType(val string) (AlertType, error) {
	switch strings.ToLower(val) {
	case "sonar":
		return AlertTypeSonar, nil
	case "radar":
		return AlertTypeRadar, nil
	default:
		return 0, fmt.Errorf("Invalid alert type '%s'", val)
	}
}

func (c AlertChange) String() string {
	switch c {
	case AlertChangeAny:
		return "ANY"
	case AlertChangeToUp:
		return "UP"
	case AlertChangeToDown:
		return "DOWN"
	default:
		return fmt.Sprintf("<unknown %d>", int(c))
	}
}

// ParseAlertChange parses an alert type 'any','up' or 'down' to enum value
func ParseAlertChange(val string) (AlertChange, error) {
	switch strings.ToLower(val) {
	case "any":
		return AlertChangeAny, nil
	case "up":
		return AlertChangeToUp, nil
	case "down":
		return AlertChangeToDown, nil
	default:
		return 0, fmt.Errorf("Invalid alert change '%s'", val)
	}
}

func (t AlertTiming) String() string {
	switch t {
	case AlertTimingImmediate:
		return "IMMEDIATE"
	case AlertTimingSummary:
		return "SUMMARY"
	case AlertTimingBoth:
		return "BOTH"
	default:
		return fmt.Sprintf("<unknown %d>", int(t))
	}
}

// ParseAlertTiming parses an alert type 'immediate','summary' or 'both' to enum value
func ParseAlertTiming(val string) (AlertTiming, error) {
	switch strings.ToLower(val) {
	case "immediate":
		return AlertTimingImmediate, nil
	case "summary":
		return AlertTimingSummary, nil
	case "both":
		return AlertTimingBoth, nil
	default:
		return 0, fmt.Errorf("Invalid alert timing '%s'", val)
	}
}

// CreateAlert creates new alerts.
func (c *Client) CreateAlert(name string, t AlertType, platform int,
	change AlertChange, timing AlertTiming, emails []string, minInterval int) error {

	atype := t.String()
	achange := change.String()
	atiming := timing.String()
	zero := 0
	operator := "GT"
	peers := []int{}
	locations := []string{}

	alert := Alert{
		Name:            &name,
		Type:            &atype,
		Platform:        &platform,
		NotifyChange:    &achange,
		Timing:          &atiming,
		Emails:          &emails,
		Debounce:        &minInterval,
		ProbeType:       &zero,
		CompareOperator: &operator,
		Peers:           &peers,
		Locations:       &locations,
	}

	return c.postJSON(baseURL+alertsConfigPath, &alert, nil)
}

// GetAlerts returns all configured alerts.
func (c *Client) GetAlerts() ([]*Alert, error) {
	var resp []*Alert
	err := c.getJSON(baseURL+alertsConfigPath, &resp)

	if err != nil {
		return nil, err
	}

	return resp, nil
}
