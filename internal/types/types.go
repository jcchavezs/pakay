package types

// SentinelVal is a dummy type used for the sentinel method in TypedConfig.
type SentinelVal struct{}

// TypedConfig is an interface that all secret source configuration types must implement.
type TypedConfig interface {
	// SentinelFn is a sentinel method to ensure interface compliance can only
	// be achieved by types within this repository.
	SentinelFn(SentinelVal)
}
