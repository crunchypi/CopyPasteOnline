package config

// DB CONFIG
var (
	// How many seconds before entries in the sqlite db are flushed.
	ItemTimeoutSeconds int64 = 60 * 10
)

// MNEMONICS CONFIG
var (
	DefaultDrawN        = 3
	DefaultDrawTryLimit = 10
)

// DOS CONFIG
var (
	// FlushDeltaSeconds dictates how many seconds should pass before an accessControl
	// entry becomes stale and can be removed
	FlushDeltaSeconds int64 = 60
	// LimitDeltaSeconds dictates time in: access n/time
	LimitDeltaSeconds int64 = 60
	// AccessPerLimit dictates access in: access n/time
	AccessPerLimit int = 120
)

// APP CONFIG
var (
	IP   = "localhost"
	Port = "80"
)
