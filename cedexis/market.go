package cedexis

// Market holds the market as two-character 'ISO codes'
type Market string

const (
	// MarketGlobal represents the global market (world-wide)
	MarketGlobal Market = "GL"

	// MarketAfrica represents the african market
	MarketAfrica Market = "AF"

	// MarketAsia represents the asian market
	MarketAsia Market = "AS"

	// MarketEurope represents the european market
	MarketEurope Market = "EU"

	// MarketNorthAmerica represents the north-american market
	MarketNorthAmerica Market = "NA"

	// MarketOceania represents the oceania market
	MarketOceania Market = "OC"

	// MarketSouthAmerica represents the south-american market
	MarketSouthAmerica Market = "SA"

	// MarketUnknown is a placeholder when the market is unknown
	MarketUnknown Market = "XX"
)
