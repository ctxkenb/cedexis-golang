package cedexis

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
	ID   *int    `json:"id,omitempty"`
	Name *string `json:"name,omitempty"`
}

// SonarConfig represents the sonar configuration for a platform
type SonarConfig struct {
	Enabled              *bool   `json:"enabled,omitempty"`
	URL                  *string `json:"url,omitempty"`
	PollInterviewSeconds *int    `json:"pollIntervalSeconds,omitempty"`
	Timeout              *int    `json:"timeout,omitempty"`
	Method               *string `json:"method,omitempty"`
	IgnoreSSLErrors      *bool   `json:"ignoreSSLErrors,omitempty"`
	MaintenanceMode      *bool   `json:"maintenanceMode,omitempty"`
	Host                 *string `json:"host,omitempty"`
	Market               *string `json:"market,omitempty"`
	RequestContentType   *string `json:"requestContentType,omitempty"`
	ResponseBodyMatch    *string `json:"responseBodyMatch,omitempty"`
	ResponseMatchType    *string `json:"responseMatchType,omitempty"`
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
	err := c.getJSON(path, &resp)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

// GetPrivatePlatforms gets the configured platforms, optionally filtered by tag
func (c *Client) GetPrivatePlatforms(tag *string) ([]*PlatformConfig, error) {
	var resp []*PlatformConfig
	err := c.getJSON(baseURL+platformsConfigPath, &resp)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

// CreatePrivatePlatform creates a new platform
func (c *Client) CreatePrivatePlatform(spec *PlatformConfig) (*PlatformConfig, error) {
	var resp PlatformConfig
	err := c.postJSON(baseURL+platformsConfigPath, spec, &resp)

	if err != nil {
		return nil, err
	}

	return &resp, nil
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

	platformid := 1
	fusionArchetype := "CUSTOM"
	openmixVisible := true
	privateArchetype := false
	publicChartEnabled := false
	radarCacheBursting := true
	radarMajorNetworksOnly := true
	radarUsePublicData := true
	sonarMethod := "GET"
	sonarPollIntervalSeconds := 60

	return &PlatformConfig{
		Category:                  &NameID{ID: &platformid},
		DisplayName:               &displayName,
		FusionArchetype:           &fusionArchetype,
		IntendedUse:               &description,
		Name:                      &name,
		OpenmixVisible:            &openmixVisible,
		PrivateArchetype:          &privateArchetype,
		PublicChartEnabled:        &publicChartEnabled,
		PublicProviderArchetypeID: &archetypeID,
		RadarConfig: &RadarConfig{
			CacheBursting:     &radarCacheBursting,
			MajorNetworksOnly: &radarMajorNetworksOnly,
			UsePublicData:     &radarUsePublicData,
		},
		SonarConfig: &SonarConfig{
			Method:               &sonarMethod,
			PollInterviewSeconds: &sonarPollIntervalSeconds,
		},
		Tags: &tags,
	}

}
