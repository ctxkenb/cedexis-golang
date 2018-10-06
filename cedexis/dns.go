package cedexis

import (
	"encoding/json"
	"fmt"
	"strings"
)

const dnsConfigPath = "/config/authdns.json"
const dnsRecordConfigPath = "/config/authdns.json/record"

const (
	// RecordTypeA is for IPv4 addresses
	RecordTypeA string = "A"

	// RecordTypeAAAA is for IPv6 addresses
	RecordTypeAAAA = "AAAA"

	// RecordTypeTXT is for text records
	RecordTypeTXT = "TXT"

	// RecordTypeSPF is for Sender Policy Framework records
	RecordTypeSPF = "SPF"

	// RecordTypePTR is for reverse DNS lookups
	RecordTypePTR = "PTR"

	// RecordTypeCNAME is for canonical name records
	RecordTypeCNAME = "CNAME"

	// RecordTypeNS is for domain delegation (nameserver) records
	RecordTypeNS = "NS"

	// RecordTypeOpenmix is for Openmix applications
	RecordTypeOpenmix = "OPX"

	// RecordTypeCAA is for Certification Authority Authorization records
	RecordTypeCAA = "CAA"

	// RecordTypeSRV is for defining servers for specified services
	RecordTypeSRV = "SRV"

	// RecordTypeMX is for mail server records
	RecordTypeMX = "MX"
)

// SOAResponse is used for records of type SOA
type SOAResponse struct {
	Mname   string `json:"mname"`
	Rname   string `json:"rname"`
	Refresh int    `json:"refresh"`
	Retry   int    `json:"retry"`
	Expire  int    `json:"expire"`
	Minimum int    `json:"rxttl"`
	Serial  int    `json:"serial"`
}

// AddressesResponse is used for records of type A and AAAA
type AddressesResponse struct {
	Addresses []string `json:"addresses"`
}

// TextStringsResponse is used for records of type TXT and SPF
type TextStringsResponse struct {
	TextStrings []string `json:"textStrings"`
}

// DomainNameResponse is used for records of type PTR and CNAME
type DomainNameResponse struct {
	DomainName string `json:"domainName"`
}

// DomainNamesResponse is used for records of type NS
type DomainNamesResponse struct {
	DomainNames []string `json:"domainNames"`
}

// AppResponse is used for records of type OPX
type AppResponse struct {
	AppID int `json:"appId"`
}

// CAAEntry is a single entry in a CAA record
type CAAEntry struct {
	Tag   string `json:"tag"`
	Flags int    `json:"flags"`
	Value string `json:"value"`
}

// CAAResponse is used for records of type CAA
type CAAResponse struct {
	Entries []CAAEntry `json:"entries"`
}

// SRVEntry is a single entry in an SRV record
type SRVEntry struct {
	Priority int    `json:"priority"`
	Weight   int    `json:"weight"`
	Port     int    `json:"port"`
	Target   string `json:"target"`
}

// SRVResponse is used for records of type SRV
type SRVResponse struct {
	Entries []SRVEntry `json:"textStrings"`
}

// MXHost represents an individual MX host in an MX record
type MXHost struct {
	Priority int    `json:"priority"`
	Target   string `json:"target"`
}

// MXResponse is used for records of type MX
type MXResponse struct {
	Hosts []MXHost `json:"hosts"`
}

// Record represents a DNS record in a zone.
type Record struct {
	ID            *int    `json:"id,omitempty"`
	DNSZoneID     *int    `json:"dnsZoneId,omitempty"`
	TTL           *int    `json:"ttl,omitempty"`
	SubdomainName *string `json:"subdomainName,omitempty"`
	RecordType    *string `json:"recordType,omitempty"`
	Response      *string `json:"response,omitempty"`
}

// Zone represents a DNS zone.
type Zone struct {
	ID                  *int     `json:"id,omitempty"`
	DomainName          *string  `json:"domainName,omitempty"`
	Description         *string  `json:"description,omitempty"`
	Tags                *string  `json:"tags,omitempty"`
	ImportContents      *string  `json:"importContents,omitempty"`
	IsPrimary           *bool    `json:"isPrimary,omitempty"`
	ZoneTransferEnabled *bool    `json:"zoneTranserEnabled,omitempty"`
	LastImport          *string  `json:"lastImport,omitempty"`
	Records             []Record `json:"records,omitempty"`
}

// SetResponseObject sets the Response by serializing an XXXXXResponse struct
func (r *Record) SetResponseObject(v interface{}) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	val := string(b)
	r.Response = &val
	return nil
}

// ResponseObject gets the Response field as a structure
func (r *Record) ResponseObject() (interface{}, error) {
	obj := r.responseObject()
	err := json.Unmarshal([]byte(*r.Response), obj)
	if err != nil {
		return nil, err
	}

	return obj, nil
}

func (r *Record) responseObject() interface{} {
	switch *r.RecordType {
	case "A":
		return &AddressesResponse{}
	case "AAAA":
		return &AddressesResponse{}
	case "TXT":
		return &TextStringsResponse{}
	case "SPF":
		return &TextStringsResponse{}
	case "PTR":
		return &DomainNameResponse{}
	case "CNAME":
		return &DomainNameResponse{}
	case "NS":
		return &DomainNamesResponse{}
	case "OPX":
		return &AppResponse{}
	case "CAA":
		return &CAAResponse{}
	case "SRV":
		return &SRVResponse{}
	case "MX":
		return &MXResponse{}
	default:
		return nil
	}
}

// DiffersFrom determins if any fields specified in this record differ from the other record
func (r *Record) DiffersFrom(other *Record) bool {
	if r == nil {
		return other == nil
	}

	if r.ID != nil && intsDiffer(r.ID, other.ID) {
		return true
	}

	if r.DNSZoneID != nil && intsDiffer(r.DNSZoneID, other.DNSZoneID) {
		return true
	}

	if r.TTL != nil && intsDiffer(r.TTL, other.TTL) {
		return true
	}

	if r.SubdomainName != nil && stringsDiffer(r.SubdomainName, other.SubdomainName) {
		return true
	}

	if r.RecordType != nil && stringsDiffer(r.RecordType, other.RecordType) {
		return true
	}

	if r.Response != nil && stringsDiffer(r.Response, other.Response) {
		return true
	}

	return false
}

// CreateZone creates a new DNS zone, with optional zone file
func (c *Client) CreateZone(name string, description string, tags []string, importContents *string) (*Zone, error) {
	t := true
	tagsString := strings.Join(tags, ",")

	zone := &Zone{
		DomainName:     &name,
		Description:    &description,
		Tags:           &tagsString,
		ImportContents: importContents,
		IsPrimary:      &t,
	}

	err := c.postJSON(baseURL+dnsConfigPath, zone, zone)
	if err != nil {
		return nil, err
	}

	if len(c.zoneCache) > 0 {
		c.zoneCache[*zone.ID] = zone
	}

	return zone, nil
}

// GetZones returns all configured zones.
func (c *Client) GetZones() ([]*Zone, error) {
	var resp []*Zone

	if len(c.zoneCache) == 0 {
		err := c.getJSON(baseURL+dnsConfigPath, &resp)

		if err != nil {
			return nil, err
		}

		c.zoneCache = map[int]*Zone{}
		for _, z := range resp {
			c.zoneCache[*z.ID] = z
		}

	} else {
		resp = make([]*Zone, 0, len(c.zoneCache))
		for _, z := range c.zoneCache {
			resp = append(resp, z)
		}
	}

	return resp, nil
}

// GetZone gets a zone.
func (c *Client) GetZone(id int) (*Zone, error) {
	var result *Zone

	result = c.zoneCache[id]
	if result != nil {
		return result, nil
	}

	result = &Zone{}
	err := c.getJSON(baseURL+dnsConfigPath+fmt.Sprintf("/%d", id), result)
	if err != nil {
		return nil, err
	}

	if len(c.zoneCache) > 0 {
		c.zoneCache[*result.ID] = result
	}

	return result, err
}

// DeleteZone deletes an alert.
func (c *Client) DeleteZone(id int) error {
	err := c.delete(baseURL + dnsConfigPath + fmt.Sprintf("/%d", id))
	if err != nil {
		delete(c.zoneCache, id)
	}
	return err
}

// GetZoneByName gets a zone by name, returning nil if not found
func (c *Client) GetZoneByName(name string) (*Zone, error) {
	zones, err := c.GetZones()
	if err != nil {
		return nil, err
	}

	for _, z := range zones {
		if strings.ToLower(name) == strings.ToLower(*z.DomainName) {
			return z, nil
		}
	}

	// Not found
	return nil, nil
}

// CreateRecord creates a DNS record
func (c *Client) CreateRecord(r *Record) (*Record, error) {
	out := Record{}
	err := c.postJSON(baseURL+dnsRecordConfigPath, r, &out)
	if err != nil {
		return nil, err
	}

	if c.zoneCache[*r.DNSZoneID] != nil {
		c.zoneCache[*r.DNSZoneID].Records = append(c.zoneCache[*r.DNSZoneID].Records, out)
	}

	return &out, nil
}

// GetRecordByName gets a DNS record
func (c *Client) GetRecordByName(zone string, name string, rtype string) (*Record, error) {
	z, err := c.GetZoneByName(zone)
	if err != nil {
		return nil, err
	}

	// Zone not found, so record not found
	if z == nil {
		return nil, nil
	}

	for _, r := range z.Records {
		if strings.ToLower(name) == strings.ToLower(*r.SubdomainName) && rtype == *r.RecordType {
			return &r, nil
		}
	}

	// Not found
	return nil, nil
}

// UpdateRecord updates a record
func (c *Client) UpdateRecord(r *Record) (*Record, error) {
	out := Record{}
	err := c.putJSON(baseURL+dnsRecordConfigPath+fmt.Sprintf("/%d", *r.ID), r, &out)
	if err != nil {
		return nil, err
	}

	if c.zoneCache[*r.DNSZoneID] != nil {
		for i, existing := range c.zoneCache[*r.DNSZoneID].Records {
			if existing.ID == r.ID {
				c.zoneCache[*r.DNSZoneID].Records[i] = out
			}
		}
	}

	return &out, nil
}
