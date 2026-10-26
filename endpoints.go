package godhan

// Default API Gateway URLs
const (
	DefaultBaseURI = "https://api.dhan.co/v2"
	DefaultFeedURI = "wss://api-feed.dhan.co"
)

// API Endpoint Paths
const (
	// Authentication
	URIRenewToken = "/RenewToken"

	// IP Whitelisting
	URIGetIP    = "/ip/getIP"
	URISetIP    = "/ip/setIP"
	URIModifyIP = "/ip/modifyIP"

	// Trader Controls
	URIKillSwitch = "/killswitch"
	URIPnLExit    = "/pnlExit"

	// EDIS
	URIEDISTpin     = "/edis/tpin"
	URIEDISForm     = "/edis/form"
	URIEDISBulkForm = "/edis/bulkform"
	URIEDISInquire  = "/edis/inquire" // /edis/inquire/{isin}

	// Orders
	URIOrders         = "/orders"
	URIOrdersSlicing  = "/orders/slicing"
	URIOrdersExternal = "/orders/external" // /orders/external/{correlation-id}

	// Super Orders (Bracket Orders)
	URISuperOrders = "/super/orders"

	// Forever Orders (GTT)
	URIForeverOrders = "/forever/orders"

	// Conditional & Multi Orders (Alerts)
	URIAlertOrders      = "/alerts/orders"
	URIMultiAlertOrders = "/alerts/multi/orders"

	// Trades
	URITrades = "/trades"

	// Portfolio
	URIHoldings        = "/holdings"
	URIPositions       = "/positions"
	URIPositionsConvert = "/positions/convert"

	// Funds & Margin
	URIFundLimit            = "/fundlimit"
	URIMarginCalculator     = "/margincalculator"
	URIMarginCalculatorMulti = "/margincalculator/multi"

	// Statements
	URILedger = "/ledger"

	// Charts & Market Analytics
	URIChartsHistorical    = "/charts/historical"
	URIChartsIntraday      = "/charts/intraday"
	URIChartsRollingOption = "/charts/rollingoption"
	URIDataCompanyInfo     = "/data/companyinfo"
	URIDataMarketMovers    = "/data/marketmovers"
	URIDataNewsHeadline    = "/data/newsheadline"
	URIDataTechnical       = "/data/technical"
)
