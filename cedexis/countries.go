package cedexis

const countriesReportPath = "/reporting/countries.json"
const countriesSimpleReportPath = "/reporting/countries.json/simple"

type Location struct {
	ID      int    `json:"id"`
	ISOCode string `json:"isoCode"`
	Name    string `json:"name"`
}

type Subcontinent struct {
	Location
	Continent Location `json:"continent"`
}

type Country struct {
	Location
	OverallPercent float32      `json:"overallPercent"`
	Subcontinent   Subcontinent `json:"subcontinent"`
	Market         Location     `json:"market"`
}

// GetCountries gets all details for countries.
func (c *Client) GetCountries() ([]*Country, error) {
	var resp []*Country

	if len(c.countriesCache) == 0 {
		err := c.getJSON(baseURL+countriesReportPath, &resp)

		if err != nil {
			return nil, err
		}

		c.countriesCache = map[int]*Country{}
		for _, a := range resp {
			c.countriesCache[a.ID] = a
		}

		return resp, nil
	}

	resp = make([]*Country, 0, len(c.countriesCache))
	for _, a := range c.countriesCache {
		resp = append(resp, a)
	}

	return resp, nil
}

// GetCountryByName gets a country by name.
func (c *Client) GetCountryByName(name string) (*Country, error) {
	countries, err := c.GetCountries()
	if err != nil {
		return nil, err
	}

	for _, c := range countries {
		if c.Name == name {
			return c, nil
		}
	}

	return nil, nil
}
