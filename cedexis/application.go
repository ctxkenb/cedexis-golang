package cedexis

import (
	"fmt"
	"sort"
)

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
	ID                    *int                   `json:"id,omitempty"`
	Name                  *string                `json:"name,omitempty"`
	Description           *string                `json:"description,omitempty"`
	Enabled               *bool                  `json:"enabled,omitempty"`
	RelativePercent       *float64               `json:"relativePercent,omitempty"`
	Version               *int                   `json:"version,omitempty"`
	AppData               *string                `json:"appData,omitempty"`
	Created               *string                `json:"created,omitempty"`
	Modified              *string                `json:"modified,omitempty"`
	Status                *string                `json:"status,omitempty"`
	Tags                  *[]string              `json:"tags,omitempty"`
	Protocol              *string                `json:"protocol,omitempty"`
	Type                  *string                `json:"type,omitempty"`
	ModifiedBy            *string                `json:"modifiedBy,omitempty"`
	AvailabilityThreshold *int                   `json:"availabilityThreshold,omitempty"`
	UseRadarAvailability  *bool                  `json:"useRadarAvailability,omitempty"`
	Platforms             *[]ApplicationPlatform `json:"platforms,omitempty"`
	FallbackCname         *string                `json:"fallbackCname,omitempty"`
	Cname                 *string                `json:"cname,omitempty"`
	CorsHeader            *string                `json:"corsHeader,omitempty"`
	EDNSEnabled           *bool                  `json:"eDNSEnabled,omitempty"`
	TTL                   *int                   `json:"ttl,omitempty"`
	HOPXSecureApp         *bool                  `json:"hopxSecureApp,omitempty"`
}

// NewApplication creates a new Application object, use CreateApplication to actually create the application.
func NewApplication(name string, description string, appType string, fallbackName string, availabilityThreshold int, targets []ApplicationPlatform) *Application {
	return &Application{
		Name:                  &name,
		Description:           &description,
		Type:                  &appType,
		FallbackCname:         &fallbackName,
		AvailabilityThreshold: &availabilityThreshold,
		Platforms:             &targets,
	}
}

// DiffersFrom determines if fields in this application differ from the other application
func (h *Handicap) DiffersFrom(other *Handicap) bool {
	if h == nil {
		return other == nil
	}

	if h.Handicap != nil && intsDiffer(h.Handicap, other.Handicap) {
		return true
	}

	if h.Enabled != nil && boolsDiffer(h.Enabled, other.Enabled) {
		return true
	}

	return false
}

func handicapMapDiffers(a map[string]Handicap, b map[string]Handicap) bool {
	if len(a) != len(b) {
		return true
	}

	for k, v := range a {
		h := &v
		h2 := b[k]
		if h.DiffersFrom(&h2) {
			return true
		}
	}

	return false
}

func platformArraysDiffer(a *[]ApplicationPlatform, b *[]ApplicationPlatform) bool {
	if a == nil {
		return b == nil
	}

	if b == nil {
		return true
	}

	if len(*a) != len(*b) {
		return true
	}

	sort.Slice(*a, func(i, j int) bool {
		return *(*a)[i].Cname < *(*a)[j].Cname
	})

	sort.Slice(*b, func(i, j int) bool {
		return *(*b)[i].Cname < *(*b)[j].Cname
	})

	for i := range *a {
		bPlatform := (*b)[i]
		if (*a)[i].DiffersFrom(&bPlatform) {
			return true
		}
	}

	return false
}

// DiffersFrom determines if fields in this application differ from the other application
func (g *ApplicationGeo) DiffersFrom(other *ApplicationGeo) bool {
	if handicapMapDiffers(g.Country, other.Country) {
		return true
	}

	if handicapMapDiffers(g.Market, other.Market) {
		return true
	}

	if g.Global != nil && g.Global.DiffersFrom(other.Global) {
		return true
	}

	return false
}

// DiffersFrom determines if fields in this application differ from the other application
func (p *ApplicationPlatform) DiffersFrom(other *ApplicationPlatform) bool {
	if p == nil {
		return other == nil
	}

	if p.ID != nil && intsDiffer(p.ID, other.ID) {
		return true
	}

	if p.Cname != nil && stringsDiffer(p.Cname, other.Cname) {
		return true
	}

	if p.Handicap != nil && intsDiffer(p.Handicap, other.Handicap) {
		return true
	}

	if p.Weight != nil && intsDiffer(p.Weight, other.Weight) {
		return true
	}

	if p.Enabled != nil && boolsDiffer(p.Enabled, other.Enabled) {
		return true
	}

	if p.Description != nil && stringsDiffer(p.Description, other.Description) {
		return true
	}

	if p.SonarEnabled != nil && boolsDiffer(p.SonarEnabled, other.SonarEnabled) {
		return true
	}

	if p.Geo != nil && p.Geo.DiffersFrom(other.Geo) {
		return true
	}

	return false
}

// DiffersFrom determines if fields in this application differ from the other application
func (a *Application) DiffersFrom(other *Application) bool {
	if a == nil {
		return other == nil
	}

	if a.ID != nil && intsDiffer(a.ID, other.ID) {
		return true
	}

	if a.Name != nil && stringsDiffer(a.Name, other.Name) {
		return true
	}

	if a.Description != nil && stringsDiffer(a.Description, other.Description) {
		return true
	}

	if a.Enabled != nil && boolsDiffer(a.Enabled, other.Enabled) {
		return true
	}

	if a.RelativePercent != nil && float64sDiffer(a.RelativePercent, other.RelativePercent) {
		return true
	}

	if a.Version != nil && intsDiffer(a.Version, other.Version) {
		return true
	}

	if a.AppData != nil && stringsDiffer(a.AppData, other.AppData) {
		return true
	}

	if a.Created != nil && stringsDiffer(a.Created, other.Created) {
		return true
	}

	if a.Modified != nil && stringsDiffer(a.Modified, other.Modified) {
		return true
	}

	if a.Status != nil && stringsDiffer(a.Status, other.Status) {
		return true
	}

	if a.Tags != nil && stringArraysDiffer(a.Tags, other.Tags) {
		return true
	}

	if a.Protocol != nil && stringsDiffer(a.Protocol, other.Protocol) {
		return true
	}

	if a.Type != nil && stringsDiffer(a.Type, other.Type) {
		return true
	}

	if a.ModifiedBy != nil && stringsDiffer(a.ModifiedBy, other.ModifiedBy) {
		return true
	}

	if a.AvailabilityThreshold != nil && intsDiffer(a.AvailabilityThreshold, other.AvailabilityThreshold) {
		return true
	}

	if a.UseRadarAvailability != nil && boolsDiffer(a.UseRadarAvailability, other.UseRadarAvailability) {
		return true
	}

	if a.Platforms != nil && platformArraysDiffer(a.Platforms, other.Platforms) {
		return true
	}

	if a.FallbackCname != nil && stringsDiffer(a.FallbackCname, other.FallbackCname) {
		return true
	}

	if a.Cname != nil && stringsDiffer(a.Cname, other.Cname) {
		return true
	}

	if a.CorsHeader != nil && stringsDiffer(a.CorsHeader, other.CorsHeader) {
		return true
	}

	if a.EDNSEnabled != nil && boolsDiffer(a.EDNSEnabled, other.EDNSEnabled) {
		return true
	}

	if a.TTL != nil && intsDiffer(a.TTL, other.TTL) {
		return true
	}

	if a.HOPXSecureApp != nil && boolsDiffer(a.HOPXSecureApp, other.HOPXSecureApp) {
		return true
	}

	return false
}

// GetApplications gets all applications.
func (c *Client) GetApplications() ([]*Application, error) {
	var resp []*Application

	if len(c.appCache) == 0 {
		err := c.getJSON(baseURL+appsConfigPath, &resp)

		if err != nil {
			return nil, err
		}

		c.appCache = map[int]*Application{}
		for _, a := range resp {
			c.appCache[*a.ID] = a
		}

	} else {
		resp = make([]*Application, 0, len(c.appCache))
		for _, a := range c.appCache {
			resp = append(resp, a)
		}
	}

	return resp, nil
}

// GetApplication gets an alert.
func (c *Client) GetApplication(id int) (*Application, error) {
	var result *Application

	result = c.appCache[id]
	if result != nil {
		return result, nil
	}

	result = &Application{}
	err := c.getJSON(baseURL+appsConfigPath+fmt.Sprintf("/%d", id), result)
	if err != nil {
		return nil, err
	}

	if len(c.appCache) > 0 {
		c.appCache[*result.ID] = result
	}

	return result, err
}

// GetApplicationByName gets an application by name.
func (c *Client) GetApplicationByName(name string) (*Application, error) {
	apps, err := c.GetApplications()
	if err != nil {
		return nil, err
	}

	for _, a := range apps {
		if *a.Name == name {
			return a, nil
		}
	}

	return nil, nil
}

// CreateApplication creates an application
func (c *Client) CreateApplication(app *Application) (*Application, error) {
	out := &Application{}
	err := c.postJSON(baseURL+appsConfigPath, app, out)
	if err != nil {
		return nil, err
	}

	if len(c.appCache) > 0 {
		c.appCache[*out.ID] = out
	}

	return out, nil
}

// UpdateApplication updates an app.
func (c *Client) UpdateApplication(app *Application) (*Application, error) {
	err := c.putJSON(baseURL+appsConfigPath+fmt.Sprintf("/%d", *app.ID), app, nil)
	if err != nil {
		return nil, err
	}

	out := &Application{}
	err = c.getJSON(baseURL+appsConfigPath+fmt.Sprintf("/%d", *app.ID), out)
	if err != nil && len(c.appCache) > 0 {
		c.appCache[*out.ID] = out
	}

	return out, nil
}

// DeleteApplication deletes an app.
func (c *Client) DeleteApplication(id int) error {
	err := c.delete(baseURL + appsConfigPath + fmt.Sprintf("/%d", id))
	if err != nil {
		delete(c.appCache, id)
	}
	return err
}
