package trackingmore

import sdkcore "github.com/wakflo/go-sdk/core"

var courierCodes = []*sdkcore.AutoFormSchema{
	{Const: "fedex", Title: "FedEx"},
	{Const: "ups", Title: "UPS"},
	{Const: "usps", Title: "USPS"},
	{Const: "dhl", Title: "DHL"},
	{Const: "speedaf", Title: "Speedaf"},
	{Const: "ontrac", Title: "OnTrac"},
	{Const: "lasership", Title: "LaserShip"},
	{Const: "tforce", Title: "TForce"},
	{Const: "rrdonnelley", Title: "RR Donnelley"},
	{Const: "estes", Title: "Estes"},
	{Const: "old_dominion", Title: "Old Dominion Freight Line"},
	{Const: "saia", Title: "Saia"},
	{Const: "xpo", Title: "XPO Logistics"},
	{Const: "forward_air", Title: "Forward Air"},
	{Const: "speedee", Title: "SpeeDee Delivery"},
	{Const: "gso", Title: "GSO (Golden State Overnight)"},
	{Const: "tnt", Title: "TNT"},
	{Const: "pitney_bowes", Title: "Pitney Bowes"},
	{Const: "purolator", Title: "Purolator"},
	{Const: "newgistics", Title: "Newgistics"},
	{Const: "yrc_freight", Title: "YRC Freight"},
	{Const: "land_air", Title: "Land Air Express"},
	{Const: "u_ship", Title: "uShip"},
	{Const: "xpress_global", Title: "Xpress Global Systems"},
	{Const: "daylight_transport", Title: "Daylight Transport"},
	{Const: "r_l_carriers", Title: "R+L Carriers"},
	{Const: "central_transport", Title: "Central Transport"},
	{Const: "southeastern_freight", Title: "Southeastern Freight Lines"},
	{Const: "ward_trucking", Title: "Ward Trucking"},
	{Const: "cross_country", Title: "Cross Country Courier"},
	{Const: "gls", Title: "GLS"},
	{Const: "roadrunner", Title: "Roadrunner Transportation"},
	{Const: "aaa_cooper", Title: "AAA Cooper Transportation"},
	{Const: "dohrn_transfer", Title: "Dohrn Transfer"},
	{Const: "rwc", Title: "RWC (Regional West)"},
	{Const: "brown_integrated", Title: "Brown Integrated Logistics"},
	{Const: "forwardair", Title: "Forward Air"},
	{Const: "r_l_freight", Title: "R+L Global Logistics"},
	{Const: "sudden_valley", Title: "Sudden Valley"},
	{Const: "best_overland", Title: "Best Overland Freight"},
	{Const: "western_freight", Title: "Western Freight"},
	{Const: "americold", Title: "Americold"},
	{Const: "allied_express", Title: "Allied Express"},
	{Const: "yellow_corp", Title: "Yellow Corporation"},
	{Const: "safelite", Title: "Safelite"},
	{Const: "mainfreight", Title: "Mainfreight"},
	{Const: "dynamex", Title: "Dynamex"},
	{Const: "hermes", Title: "Hermes"},
	{Const: "ups_freight", Title: "UPS Freight"},
	{Const: "fedex_ground", Title: "FedEx Ground"},
	{Const: "usps_priority", Title: "USPS Priority Mail"},
	{Const: "china-post", Title: "China Post"},
	{Const: "china-ems", Title: "China EMS"},
	{Const: "postnord", Title: "PostNord"},
	{Const: "fastway_uk", Title: "Fastway UK"},
	{Const: "xpo_logistics", Title: "XPO Logistics"},
	{Const: "dpd_canada", Title: "DPD Canada"},
	{Const: "intime_express", Title: "InTime Express"},
	{Const: "p4d", Title: "P4D (UK)"},
	{Const: "frakt24", Title: "Frakt24"},
	{Const: "bring_express", Title: "Bring Express"},
	{Const: "mbe_canada", Title: "MBE Canada"},
	{Const: "chit_chat_express", Title: "Chit Chat Express"},
	{Const: "gophr", Title: "Gophr"},
	{Const: "city_sprint", Title: "City Sprint"},
	{Const: "purolator_courier", Title: "Purolator Courier"},
	{Const: "logistics_xpress", Title: "Logistics Xpress"},
	{Const: "ceva_logistics", Title: "CEVA Logistics"},
	{Const: "cool_express", Title: "Cool Express"},
	{Const: "hong-kong-post", Title: "Hong Kong Post"},
	{Const: "singapore-post", Title: "Singapore Post"},
	{Const: "swiss-post", Title: "Swiss Post"},
	{Const: "royal-mail", Title: "Royal Mail"},
	{Const: "postnl-parcels", Title: "PostNL International"},
	{Const: "canada-post", Title: "Canada Post"},
	{Const: "australia-post", Title: "Australia Post"},
	{Const: "new-zealand-post", Title: "New Zealand Post"},
	{Const: "parcel-force", Title: "Parcelforce"},
	{Const: "belgium-post", Title: "Bpost"},
	{Const: "brazil-correios", Title: "Brazil Correios"},
	{Const: "russian-post", Title: "Russian Post"},
	{Const: "malaysia-post", Title: "Malaysia Post"},
	{Const: "maldives-post", Title: "Maldives Post"},
	{Const: "malta-post", Title: "Malta Post"},
	{Const: "mauritius-post", Title: "Mauritius Post"},
	{Const: "correos-mexico", Title: "Mexico Post"},
	{Const: "moldova-post", Title: "Moldova Post"},
	{Const: "la-poste-monaco", Title: "Monaco Post"},
	{Const: "monaco-ems", Title: "Monaco EMS"},
	{Const: "mongol-post", Title: "Mongol Post"},
	{Const: "posta-crne-gore", Title: "Montenegro Post"},
}
