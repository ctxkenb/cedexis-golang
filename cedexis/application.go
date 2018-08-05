package cedexis

import "fmt"

const appsConfigPath = "/config/applications/dns.json"

const (
	ApplicationTypeJavascriptV1  = "V1_JS"
	ApplicationTypeOptimalRTT    = "RT_HTTP_PERFORMANCE"
	ApplicationTypeFailover      = "STATIC_FAILOVER"
	ApplicationTypeRoundRobin    = "RR_PURE_WEIGHTED"
	ApplicationTypeStaticRouting = "STATIC_ROUTING"
	ApplicationTypeThroughput    = "KBPS_HTTP_PERFORMANCE"
)

const (
	ApplicationProtocolDns = "dns"
)

// Handicap represents a restriction against a given location
type Handicap struct {
	Handicap *int  `json:"handicap,omitempty"`
	Enabled  *bool `json:"enabled,omitempty"`
}

// ApplicationGeo represents geo restrictions against a platform for an app
type ApplicationGeo struct {
	Market  map[string]Handicap `json:"market,omitempty"`
	Country map[string]Handicap `json:"country,omitempty"`
	Global  *Handicap           `json:"global,omitempty"`
}

// ApplicationPlatform represents a platform eligable for an app
type ApplicationPlatform struct {
	ID           *int            `json:"id,omitempty"`
	Cname        *string         `json:"cname,omitempty"`
	Handicap     *int            `json:"handicap,omitempty"`
	Weight       *int            `json:"weight,omitempty"`
	Enabled      *bool           `json:"enabled,omitempty"`
	Description  *string         `json:"description,omitempty"`
	SonarEnabled *bool           `json:"sonarEnabled,omitempty"`
	Geo          *ApplicationGeo `json:"geo,omitempty"`
}

// Application represents an Openmix application
type Application struct {
	ID                    *int                  `json:"id,omitempty"`
	Name                  *string               `json:"name,omitempty"`
	Description           *string               `json:"description,omitempty"`
	Enabled               *bool                 `json:"enabled,omitempty"`
	RelativePercent       *float64              `json:"relativePercent,omitempty"`
	Version               *int                  `json:"version,omitempty"`
	AppData               *string               `json:"appData,omitempty"`
	Created               *string               `json:"created,omitempty"`
	Modified              *string               `json:"modified,omitempty"`
	Status                *string               `json:"status,omitempty"`
	Tags                  []string              `json:"tags,omitempty"`
	Protocol              *string               `json:"protocol,omitempty"`
	Type                  *string               `json:"type,omitempty"`
	ModifiedBy            *string               `json:"modifiedBy,omitempty"`
	AvailabilityThreshold *int                  `json:"availabilityThreshold,omitempty"`
	UseRadarAvailability  *bool                 `json:"useRadarAvailability,omitempty"`
	Platforms             []ApplicationPlatform `json:"platforms,omitempty"`
	FallbackCname         *string               `json:"fallbackCname,omitempty"`
	Cname                 *string               `json:"cname,omitempty"`
	CorsHeader            *string               `json:"corsHeader,omitempty"`
	EDNSEnabled           *bool                 `json:"eDNSEnabled,omitempty"`
	TTL                   *int                  `json:"ttl,omitempty"`
	HOPXSecureApp         *bool                 `json:"hopxSecureApp,omitempty"`
}

// GetApplications gets all applications.
func (c *Client) GetApplications() ([]*Application, error) {
	var resp []*Application
	err := c.getJSON(baseURL+appsConfigPath, &resp)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

// GetApplication gets an alert.
func (c *Client) GetApplication(id int) (*Alert, error) {
	result := Alert{}
	err := c.getJSON(baseURL+appsConfigPath+fmt.Sprintf("/%d", id), &result)
	if err != nil {
		return nil, err
	}
	return &result, err
}

// CreateApplication creates an application
func (c *Client) CreateApplication(name string, description string, t string, fallbackName string, availabilityThreshold int, targets []ApplicationPlatform) error {
	app := Application{
		Name:                  &name,
		Description:           &description,
		Type:                  &t,
		FallbackCname:         &fallbackName,
		AvailabilityThreshold: &availabilityThreshold,
		Platforms:             targets,
	}
	return c.postJSON(baseURL+appsConfigPath, &app, nil)
}

// DeleteApplication deletes an app.
func (c *Client) DeleteApplication(id int) error {
	return c.delete(baseURL + appsConfigPath + fmt.Sprintf("/%d", id))
}
