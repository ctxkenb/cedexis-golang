package cedexis

import (
	"fmt"
	"strings"
)

const platformsConfigPath = "/config/platforms.json"
const platformsReportingPath = "/reporting/platforms.json"
const providerCategoriesPath = "/config/platforms.json/providerCategories"

// FusionCustomConfig customizes fusion for the platform
type FusionCustomConfig struct {
	Enabled         *bool   `json:"enabled,omitempty"`
	LoadURL         *string `json:"loadUrl,omitempty"`
	LoadRateSeconds *int    `json:"loadRateSeconds,omitempty"`
}

// RadarConfig is the RADAR configuration for the platform
type RadarConfig struct {
	HTTPEnabled                *bool   `json:"httpEnabled,omitempty"`
	HTTPSEnabled               *bool   `json:"httpsEnabled,omitempty"`
	UsePublicData              *bool   `json:"usePublicData,omitempty"`
	PrimeURL                   *string `json:"primeUrl,omitempty"`
	RTTURL                     *string `json:"rttUrl,omitempty"`
	XLURL                      *string `json:"xlUrl,omitempty"`
	CustomURL                  *string `json:"customUrl,omitempty"`
	PrimeSecureURL             *string `json:"primeSecureUrl,omitempty"`
	RTTSecureURL               *string `json:"rttSecureUrl,omitempty"`
	XLSecureURL                *string `json:"xlSecureUrl,omitempty"`
	CustomSecureURL            *string `json:"customSecureUrl,omitempty"`
	Weight                     *int    `json:"weight,omitempty"`
	SubProviderID              *int    `json:"subProviderId,omitempty"`
	SubProviderOwnerZoneID     *int    `json:"subProviderOwnerZoneId,omitempty"`
	SubProviderOwnerCustomerID *int    `json:"subProviderOwnerCustomerId,omitempty"`
	MajorNetworksOnly          *bool   `json:"majorNetworksOnly,omitempty"`
	WeightEnabled              *bool   `json:"weightEnabled,omitempty"`
	CacheBursting              *bool   `json:"cacheBusting,omitempty"`
	IsoWeight                  *int    `json:"isoWeight,omitempty"`
	IsoWeightList              *string `json:"isoWeightList,omitempty"`
	IsoWeightEnabled           *bool   `json:"isoWeightEnabled,omitempty"`
	MarketWeight               *int    `json:"marketWeight,omitempty"`
	MarketWeightList           *string `json:"marketWeightList,omitempty"`
	MarketWeightEnabled        *bool   `json:"marketWeightEnabled,omitempty"`
	PrimeType                  *string `json:"primeType,omitempty"`
	RTTType                    *string `json:"rttType,omitempty"`
	XLType                     *string `json:"xlType,omitempty"`
	CustomType                 *string `json:"customType,omitempty"`
	PrimeSecureType            *string `json:"primeSecureType,omitempty"`
}

// NameID represents a name-value pair
type NameID struct {
	ID   *PlatformCategory `json:"id,omitempty"`
	Name *string           `json:"name,omitempty"`
}

// SonarConfig represents the sonar configuration for a platform
type SonarConfig struct {
	Enabled             *bool   `json:"enabled,omitempty"`
	URL                 *string `json:"url,omitempty"`
	PollIntervalSeconds *int    `json:"pollIntervalSeconds,omitempty"`
	Timeout             *int    `json:"timeout,omitempty"`
	Method              *string `json:"method,omitempty"`
	IgnoreSSLErrors     *bool   `json:"ignoreSSLErrors,omitempty"`
	MaintenanceMode     *bool   `json:"maintenanceMode,omitempty"`
	Host                *string `json:"host,omitempty"`
	Market              *Market `json:"market,omitempty"`
	RequestContentType  *string `json:"requestContentType,omitempty"`
	ResponseBodyMatch   *string `json:"responseBodyMatch,omitempty"`
	ResponseMatchType   *string `json:"responseMatchType,omitempty"`
}

// PlatformConfig represents a configured platform
type PlatformConfig struct {
	Created                     *string             `json:"created,omitempty"`
	DisplayName                 *string             `json:"displayName,omitempty"`
	FusionCustomConfig          *FusionCustomConfig `json:"fusionCustomConfig,omitempty"`
	OpenmixVisible              *bool               `json:"openmixVisible,omitempty"`
	PublicChartEnabled          *bool               `json:"publicChartEnabled,omitempty"`
	RadarConfig                 *RadarConfig        `json:"radarConfig,omitempty"`
	OwnerID                     *int                `json:"ownerId,omitempty"`
	Enabled                     *bool               `json:"enabled,omitempty"`
	Tags                        *[]string           `json:"tags,omitempty"`
	PlatformSubstitutionSources *[]int              `json:"platformSubstitutionSources,omitempty"`
	Name                        *string             `json:"name,omitempty"`
	Modified                    *string             `json:"modified,omitempty"`
	IntendedUse                 *string             `json:"intendedUse,omitempty"`
	ID                          *int                `json:"id,omitempty"`
	Category                    *NameID             `json:"category,omitempty"`
	PrivateArchetype            *bool               `json:"privateArchetype,omitempty"`
	SonarConfig                 *SonarConfig        `json:"sonarConfig,omitempty"`
	PublicProviderArchetypeID   *int                `json:"publicProviderArchetypeId,omitempty"`
	FusionArchetype             *string             `json:"fusionArchetype,omitempty"`
}

// PlatformInfo provides information about a platform
type PlatformInfo struct {
	ID                 *int    `json:"id,omitempty"`
	Name               *string `json:"name,omitempty"`
	AliasedPlatform    *NameID `json:"aliasedPlatform,omitempty"`
	IndexID            *int    `json:"indexId,omitempty"`
	Visibility         *string `json:"visibility,omitempty"`
	Category           *NameID `json:"category,omitempty"`
	IntendedUse        *string `json:"intendedUse,omitempty"`
	PublicChartEnabled *bool   `json:"publicChartEnabled,omitempty"`
	SonarConfig        *struct {
		Enabled *bool `json:"enabled,omitempty"`
	} `json:"sonarConfig,omitempty"`
	RadarConfig *struct {
		ProbeTypes []struct {
			ID int `json:"id"`
		} `json:"probeTypes"`
	} `json:"radarConfig,omitempty"`
}

// PlatformType differentiates different classes of platform
type PlatformType int

const (
	// PlatformsTypePrivate represents platforms managed by this customer
	PlatformsTypePrivate PlatformType = (1 << iota)

	// PlatformsTypeCommunity platforms are available to the cedexis community
	PlatformsTypeCommunity

	// PlatformsTypeSystem platforms
	PlatformsTypeSystem

	// PlatformsTypeAll kinds of platform
	PlatformsTypeAll PlatformType = PlatformsTypePrivate | PlatformsTypeCommunity | PlatformsTypeSystem
)

// PlatformCategory represents different categories of community platforms
type PlatformCategory int

const (
	// PlatformCategoryCloudComputing is platforms from community cloud environments
	PlatformCategoryCloudComputing PlatformCategory = 1

	// PlatformCategoryDynamicContent is platforms from community dynamic content environments
	PlatformCategoryDynamicContent PlatformCategory = 2

	// PlatformCategoryDeliveryNetwork is platforms from community CDN environments
	PlatformCategoryDeliveryNetwork PlatformCategory = 3

	// PlatformCategoryCloudStorage is platforms from community cloud storage environments
	PlatformCategoryCloudStorage PlatformCategory = 6

	// PlatformCategorySecureObjectDelivery is platforms from community secure object environments
	PlatformCategorySecureObjectDelivery PlatformCategory = 7

	// PlatformCategoryManagedDNS is platforms from community managed DNS environments
	PlatformCategoryManagedDNS PlatformCategory = 8
)

// GetProviderCategories gets the platform provider categories (CDN, Cloud, etc)
func (c *Client) GetProviderCategories() ([]*NameID, error) {
	var resp []*NameID
	err := c.getJSON(baseURL+providerCategoriesPath, &resp)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

// GetPlatforms gets information about all public and private platforms
func (c *Client) GetPlatforms(t PlatformType) ([]*PlatformInfo, error) {
	path := baseURL + platformsReportingPath
	switch t {
	case PlatformsTypeCommunity:
		path += "/community"
		break
	case PlatformsTypePrivate:
		path += "/private"
		break
	case PlatformsTypeSystem:
		path += "/system"
		break
	}

	var resp []*PlatformInfo

	// Try cache read
	if t == PlatformsTypePrivate && len(c.privatePlatformListCache) > 0 {
		resp = make([]*PlatformInfo, 0, len(c.privatePlatformListCache))
		for _, p := range c.privatePlatformListCache {
			resp = append(resp, p)
		}
		return resp, nil
	}

	// Not cached, go to service
	err := c.getJSON(path, &resp)
	if err != nil {
		return nil, err
	}

	// Cache
	c.privatePlatformListCache = map[int]*PlatformInfo{}
	for _, p := range resp {
		c.privatePlatformListCache[*p.ID] = p
	}

	return resp, nil
}

// GetEnabledPlatforms gets the avalable platforms, optionally filtered by tag
func (c *Client) GetEnabledPlatforms(tag *string) ([]*PlatformConfig, error) {
	var resp []*PlatformConfig
	err := c.getJSON(baseURL+platformsConfigPath, &resp)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

// CreatePrivatePlatform creates a new platform
func (c *Client) CreatePrivatePlatform(spec *PlatformConfig) (*PlatformConfig, error) {
	var resp = &PlatformConfig{}
	err := c.postJSON(baseURL+platformsConfigPath, spec, resp)

	if err != nil {
		return nil, err
	}

	c.privatePlatformCache[*resp.ID] = resp

	return resp, nil
}

// DeletePrivatePlatform removes a platform
func (c *Client) DeletePrivatePlatform(id int) error {
	err := c.delete(baseURL + platformsConfigPath + "/" + fmt.Sprintf("%d", id))

	if err != nil {
		delete(c.privatePlatformCache, id)
	}

	return err
}

// UpdatePrivatePlatform updates a platform
func (c *Client) UpdatePrivatePlatform(spec *PlatformConfig) error {
	var resp = &PlatformConfig{}
	err := c.putJSON(baseURL+platformsConfigPath+"/"+fmt.Sprintf("%d", *spec.ID), spec, resp)

	if err != nil {
		c.privatePlatformCache[*resp.ID] = resp
	}

	return err
}

// GetPrivatePlatform gets a platform by ID
func (c *Client) GetPrivatePlatform(id int) (*PlatformConfig, error) {
	var cfg *PlatformConfig

	cfg = c.privatePlatformCache[id]
	if cfg != nil {
		return cfg, nil
	}

	err := c.getJSON(baseURL+platformsConfigPath+"/"+fmt.Sprintf("%d", id), &cfg)

	if err != nil {
		c.privatePlatformCache[*cfg.ID] = cfg
	}

	return cfg, err
}

// GetPrivatePlatformByName gets a platform by Name
func (c *Client) GetPrivatePlatformByName(name string) (*PlatformConfig, error) {
	platforms, err := c.GetPlatforms(PlatformsTypePrivate)
	if err != nil {
		return nil, err
	}

	for _, p := range platforms {
		if strings.ToLower(name) == strings.ToLower(*p.Name) {
			return c.GetPrivatePlatform(*p.ID)
		}
	}

	// Not found
	return nil, nil
}

//
// NewPublicCloudPrivatePlatform simplifies creating a ConfiguredPlatform instance when the platform is hosted
// on a known public cloud.
//
// The parameters are defined as follows:
//     name        Unique identifier for the platform (no spaces)
//     displayName Unique friendly name for the platform (may have spaces)
//     description Description of the new platform
//     archetypeID Unique id of the underlying community public cloud platform
//     tags        Tags to apply to the new platform
//
func NewPublicCloudPrivatePlatform(
	name string,
	displayName string,
	description string,
	archetypeID int,
	tags []string) *PlatformConfig {

	platformid := PlatformCategoryCloudComputing
	radarCacheBursting := true
	radarMajorNetworksOnly := true
	radarUsePublicData := true
	sonarMethod := "GET"
	sonarPollIntervalSeconds := 60

	return &PlatformConfig{
		Category:    &NameID{ID: &platformid},
		DisplayName: &displayName,
		IntendedUse: &description,
		Name:        &name,
		PublicProviderArchetypeID: &archetypeID,
		RadarConfig: &RadarConfig{
			CacheBursting:     &radarCacheBursting,
			MajorNetworksOnly: &radarMajorNetworksOnly,
			UsePublicData:     &radarUsePublicData,
		},
		SonarConfig: &SonarConfig{
			Method:              &sonarMethod,
			PollIntervalSeconds: &sonarPollIntervalSeconds,
		},
		Tags: &tags,
	}
}

// DiffersFrom indicates if any fields in this config (that are non-nil) differ from another
// config.
func (c *PlatformConfig) DiffersFrom(other *PlatformConfig) bool {
	if c == nil {
		return other == nil
	}

	if c.Created != nil && stringsDiffer(c.Created, other.Created) {
		return true
	}

	if c.DisplayName != nil && stringsDiffer(c.DisplayName, other.DisplayName) {
		return true
	}

	if c.FusionCustomConfig != nil && c.FusionCustomConfig.DiffersFrom(other.FusionCustomConfig) {
		return true
	}

	if c.OpenmixVisible != nil && boolsDiffer(c.OpenmixVisible, other.OpenmixVisible) {
		return true
	}

	if c.PublicChartEnabled != nil && boolsDiffer(c.PublicChartEnabled, other.PublicChartEnabled) {
		return true
	}

	if c.RadarConfig != nil && c.RadarConfig.DiffersFrom(other.RadarConfig) {
		return true
	}

	if c.OwnerID != nil && intsDiffer(c.OwnerID, other.OwnerID) {
		return true
	}

	if c.Enabled != nil && boolsDiffer(c.Enabled, other.Enabled) {
		return true
	}

	if c.Tags != nil && stringArraysDiffer(c.Tags, other.Tags) {
		return true
	}

	if c.PlatformSubstitutionSources != nil && intArraysDiffer(c.PlatformSubstitutionSources, other.PlatformSubstitutionSources) {
		return true
	}

	if c.Name != nil && stringsDiffer(c.Name, other.Name) {
		return true
	}

	if c.Modified != nil && stringsDiffer(c.Modified, other.Modified) {
		return true
	}

	if c.IntendedUse != nil && stringsDiffer(c.IntendedUse, other.IntendedUse) {
		return true
	}

	if c.ID != nil && intsDiffer(c.ID, other.ID) {
		return true
	}

	if c.Category != nil && c.Category.DiffersFrom(other.Category) {
		return true
	}

	if c.PrivateArchetype != nil && boolsDiffer(c.PrivateArchetype, other.PrivateArchetype) {
		return true
	}

	if c.SonarConfig != nil && c.SonarConfig.DiffersFrom(other.SonarConfig) {
		return true
	}

	if c.PublicProviderArchetypeID != nil && intsDiffer(c.PublicProviderArchetypeID, other.PublicProviderArchetypeID) {
		return true
	}

	if c.FusionArchetype != nil && stringsDiffer(c.FusionArchetype, other.FusionArchetype) {
		return true
	}

	return false
}

// DiffersFrom indicates if any fields in this config (that are non-nil) differ from another
// config.
func (c *FusionCustomConfig) DiffersFrom(other *FusionCustomConfig) bool {
	if c == nil {
		return other == nil
	}

	if c.Enabled != nil && boolsDiffer(c.Enabled, other.Enabled) {
		return true
	}

	if c.LoadURL != nil && stringsDiffer(c.LoadURL, other.LoadURL) {
		return true
	}

	if c.LoadRateSeconds != nil && intsDiffer(c.LoadRateSeconds, other.LoadRateSeconds) {
		return true
	}

	return false
}

// DiffersFrom indicates if any fields in this config (that are non-nil) differ from another
// config.
func (c *RadarConfig) DiffersFrom(other *RadarConfig) bool {
	if c == nil {
		return other == nil
	}

	if c.HTTPEnabled != nil && boolsDiffer(c.HTTPEnabled, other.HTTPEnabled) {
		return true
	}

	if c.HTTPSEnabled != nil && boolsDiffer(c.HTTPSEnabled, other.HTTPSEnabled) {
		return true
	}

	if c.UsePublicData != nil && boolsDiffer(c.UsePublicData, other.UsePublicData) {
		return true
	}

	if c.PrimeURL != nil && stringsDiffer(c.PrimeURL, other.PrimeURL) {
		return true
	}

	if c.RTTURL != nil && stringsDiffer(c.RTTURL, other.RTTURL) {
		return true
	}

	if c.XLURL != nil && stringsDiffer(c.XLURL, other.XLURL) {
		return true
	}

	if c.CustomURL != nil && stringsDiffer(c.CustomURL, other.CustomURL) {
		return true
	}

	if c.PrimeSecureURL != nil && stringsDiffer(c.PrimeSecureURL, other.PrimeSecureURL) {
		return true
	}

	if c.RTTSecureURL != nil && stringsDiffer(c.RTTSecureURL, other.RTTSecureURL) {
		return true
	}

	if c.XLSecureURL != nil && stringsDiffer(c.XLSecureURL, other.XLSecureURL) {
		return true
	}

	if c.CustomSecureURL != nil && stringsDiffer(c.CustomSecureURL, other.CustomSecureURL) {
		return true
	}

	if c.Weight != nil && intsDiffer(c.Weight, other.Weight) {
		return true
	}

	if c.SubProviderID != nil && intsDiffer(c.SubProviderID, other.SubProviderID) {
		return true
	}

	if c.SubProviderOwnerZoneID != nil && intsDiffer(c.SubProviderOwnerZoneID, other.SubProviderOwnerZoneID) {
		return true
	}

	if c.SubProviderOwnerCustomerID != nil && intsDiffer(c.SubProviderOwnerCustomerID, other.SubProviderOwnerCustomerID) {
		return true
	}

	if c.MajorNetworksOnly != nil && boolsDiffer(c.MajorNetworksOnly, other.MajorNetworksOnly) {
		return true
	}

	if c.WeightEnabled != nil && boolsDiffer(c.WeightEnabled, other.WeightEnabled) {
		return true
	}

	if c.CacheBursting != nil && boolsDiffer(c.CacheBursting, other.CacheBursting) {
		return true
	}

	if c.IsoWeight != nil && intsDiffer(c.IsoWeight, other.IsoWeight) {
		return true
	}

	if c.IsoWeightList != nil && stringsDiffer(c.IsoWeightList, other.IsoWeightList) {
		return true
	}

	if c.IsoWeightEnabled != nil && boolsDiffer(c.IsoWeightEnabled, other.IsoWeightEnabled) {
		return true
	}

	if c.MarketWeight != nil && intsDiffer(c.MarketWeight, other.MarketWeight) {
		return true
	}

	if c.MarketWeightList != nil && stringsDiffer(c.MarketWeightList, other.MarketWeightList) {
		return true
	}

	if c.MarketWeightEnabled != nil && boolsDiffer(c.MarketWeightEnabled, other.MarketWeightEnabled) {
		return true
	}

	if c.PrimeType != nil && stringsDiffer(c.PrimeType, other.PrimeType) {
		return true
	}

	if c.RTTType != nil && stringsDiffer(c.RTTType, other.RTTType) {
		return true
	}

	if c.XLType != nil && stringsDiffer(c.XLType, other.XLType) {
		return true
	}

	if c.CustomType != nil && stringsDiffer(c.CustomType, other.CustomType) {
		return true
	}

	if c.PrimeSecureType != nil && stringsDiffer(c.PrimeSecureType, other.PrimeSecureType) {
		return true
	}

	return false
}

// DiffersFrom indicates if any fields in this config (that are non-nil) differ from another
// config.
func (c *NameID) DiffersFrom(other *NameID) bool {
	if c == nil {
		return other == nil
	}

	if c.ID != nil && c.ID.DiffersFrom(other.ID) {
		return true
	}

	if c.Name != nil && stringsDiffer(c.Name, other.Name) {
		return true
	}

	return false
}

// DiffersFrom indicates if any fields in this config (that are non-nil) differ from another
// config.
func (c *SonarConfig) DiffersFrom(other *SonarConfig) bool {
	if c == nil {
		return other == nil
	}

	if c.Enabled != nil && boolsDiffer(c.Enabled, other.Enabled) {
		return true
	}

	if c.URL != nil && stringsDiffer(c.URL, other.URL) {
		return true
	}

	if c.PollIntervalSeconds != nil && intsDiffer(c.PollIntervalSeconds, other.PollIntervalSeconds) {
		return true
	}

	if c.Timeout != nil && intsDiffer(c.Timeout, other.Timeout) {
		return true
	}

	if c.Method != nil && stringsDiffer(c.Method, other.Method) {
		return true
	}

	if c.IgnoreSSLErrors != nil && boolsDiffer(c.IgnoreSSLErrors, other.IgnoreSSLErrors) {
		return true
	}

	if c.MaintenanceMode != nil && boolsDiffer(c.MaintenanceMode, other.MaintenanceMode) {
		return true
	}

	if c.Host != nil && stringsDiffer(c.Host, other.Host) {
		return true
	}

	if c.Market != nil && c.Market.DiffersFrom(other.Market) {
		return true
	}

	if c.RequestContentType != nil && stringsDiffer(c.RequestContentType, other.RequestContentType) {
		return true
	}

	if c.ResponseBodyMatch != nil && stringsDiffer(c.ResponseBodyMatch, other.ResponseBodyMatch) {
		return true
	}

	if c.ResponseMatchType != nil && stringsDiffer(c.ResponseMatchType, other.ResponseMatchType) {
		return true
	}

	return false
}

// DiffersFrom indicates if any fields in this config (that are non-nil) differ from another
// config.
func (c *PlatformCategory) DiffersFrom(other *PlatformCategory) bool {
	if c == nil {
		return other == nil
	}

	if other == nil {
		return true
	}

	return *c != *other
}

// DiffersFrom indicates if any fields in this config (that are non-nil) differ from another
// config.
func (m *Market) DiffersFrom(other *Market) bool {
	if m == nil {
		return other == nil
	}

	if other == nil {
		return true
	}

	return *m != *other
}
