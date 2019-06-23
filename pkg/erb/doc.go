// Package erb provides primitives for parsing the Emlid Reach Binary protocol (ERB).
//
// Implementation is based on the ERB Protocol spec version 0.1.0:
//
//  https://files.emlid.com/ERB.pdf
package erb

// Supported protocol versions.
const (
	SupportedProtocolVersionHigh   = 0
	SupportedProtocolVersionMedium = 1
	SupportedProtocolVersionLow    = 0
)
