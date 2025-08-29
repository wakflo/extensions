package actions

import _ "embed"

//go:embed get_pin.go
var getPinDocs string

//go:embed delete_pin.go
var deletePinDocs string

//go:embed search_pins.go
var searchPinsDocs string

//go:embed update_pin.go
var updatePinDocs string

//go:embed get_pin_analytics.go
var getPinAnalyticsDocs string

//go:embed create_pin.go
var createPinDocs string
