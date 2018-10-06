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

// NewAlert creates a new alert object, use CreateAlert to create in Cedexis.
func (c *Client) NewAlert(name string, t AlertType, platform int,
	change AlertChange, timing AlertTiming, emails []string, minInterval int) *Alert {
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

	return &alert
}

// CreateAlert creates new alerts.
func (c *Client) CreateAlert(alert *Alert) (*Alert, error) {
	out := Alert{}
	err := c.postJSON(baseURL+alertsConfigPath, &alert, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateAlert creates new alerts.
func (c *Client) UpdateAlert(alert *Alert) (*Alert, error) {
	out := Alert{}
	err := c.putJSON(baseURL+alertsConfigPath+fmt.Sprintf("/%d", *alert.ID), &alert, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// GetAlert gets an alert.
func (c *Client) GetAlert(id int) (*Alert, error) {
	result := Alert{}
	err := c.getJSON(baseURL+alertsConfigPath+fmt.Sprintf("/%d", id), &result)
	if err != nil {
		return nil, err
	}
	return &result, err
}

// DeleteAlert deletes an alert.
func (c *Client) DeleteAlert(id int) error {
	return c.delete(baseURL + alertsConfigPath + fmt.Sprintf("/%d", id))
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

// GetAlertByName returns an alert by name.
func (c *Client) GetAlertByName(name string) (*Alert, error) {
	alerts, err := c.GetAlerts()
	if err != nil {
		return nil, err
	}

	for _, a := range alerts {
		if *a.Name == name {
			return a, nil
		}
	}

	return nil, nil
}

// DiffersFrom indicates if any fields in this config (that are non-nil) differ from another
// config.
func (a *Alert) DiffersFrom(other *Alert) bool {
	if a == nil {
		return other == nil
	}

	if a.ID != nil && intsDiffer(a.ID, other.ID) {
		return true
	}

	if a.Name != nil && stringsDiffer(a.Name, other.Name) {
		return true
	}

	if a.Type != nil && stringsDiffer(a.Type, other.Type) {
		return true
	}

	if a.Timing != nil && stringsDiffer(a.Timing, other.Timing) {
		return true
	}

	if a.Threshold != nil && intsDiffer(a.Threshold, other.Threshold) {
		return true
	}

	if a.Sensitivity != nil && intsDiffer(a.Sensitivity, other.Sensitivity) {
		return true
	}

	if a.ProbeType != nil && intsDiffer(a.ProbeType, other.ProbeType) {
		return true
	}

	if a.Version != nil && intsDiffer(a.Version, other.Version) {
		return true
	}

	if a.Debounce != nil && intsDiffer(a.Debounce, other.Debounce) {
		return true
	}

	if a.CompareOperator != nil && stringsDiffer(a.CompareOperator, other.CompareOperator) {
		return true
	}

	if a.Enabled != nil && boolsDiffer(a.Enabled, other.Enabled) {
		return true
	}

	if a.DisabledMinutes != nil && intsDiffer(a.DisabledMinutes, other.DisabledMinutes) {
		return true
	}

	if a.Locations != nil && stringArraysDiffer(a.Locations, other.Locations) {
		return true
	}

	if a.Platform != nil && intsDiffer(a.Platform, other.Platform) {
		return true
	}

	if a.Emails != nil && stringArraysDiffer(a.Emails, other.Emails) {
		return true
	}

	if a.Peers != nil && intArraysDiffer(a.Peers, other.Peers) {
		return true
	}

	if a.EventsLast24Hours != nil && intsDiffer(a.EventsLast24Hours, other.EventsLast24Hours) {
		return true
	}

	if a.CountryEvents != nil && intArraysDiffer(a.CountryEvents, other.CountryEvents) {
		return true
	}

	if a.ASNEvents != nil && intArraysDiffer(a.ASNEvents, other.ASNEvents) {
		return true
	}

	if a.RefererURI != nil && stringsDiffer(a.RefererURI, other.RefererURI) {
		return true
	}

	if a.Statistic != nil && stringsDiffer(a.Statistic, other.Statistic) {
		return true
	}

	if a.DataSource != nil && stringsDiffer(a.DataSource, other.DataSource) {
		return true
	}

	if a.AutoFill != nil && intsDiffer(a.AutoFill, other.AutoFill) {
		return true
	}

	if a.NotifyChange != nil && stringsDiffer(a.NotifyChange, other.NotifyChange) {
		return true
	}

	return false
}
