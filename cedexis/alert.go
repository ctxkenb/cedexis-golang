package cedexis

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
	ASNEvents         *[]int    `json:"asnEvents"`
	RefererURI        *string   `json:"refererUri"`
	Statistic         *string   `json:"statistic"`
	DataSource        *string   `json:"dataSource"`
	AutoFill          *int      `json:"autoFill"`
	NotifyChange      *string   `json:"notifyChange"`
}

// AlertType is the type of alert.
type AlertType int

const (
	// AlertTypeSonar indicates the alert is triggered by Sonar
	AlertTypeSonar = iota
)

// AlertChange is an enum of which events trigger the alert.
type AlertChange int

const (
	AlertChangeAny = iota
	AlertChangeToUp
	AlertChangeToDown
)

// AlertTiming is an enum of when to trigger the alert.
type AlertTiming int

const (
	// AlertTimingImmediate triggers the alert immediately.
	AlertTimingImmediate = iota
)

// GetAlerts returns all configured alerts.
func (c *Client) GetAlerts() ([]*Alert, error) {
	var resp []*Alert
	err := c.getJSON(baseURL+alertsConfigPath, &resp)

	if err != nil {
		return nil, err
	}

	return resp, nil
}
