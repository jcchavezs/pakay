package types

type SentinelVal struct{}

type TypedConfig interface{ SentinelFn(SentinelVal) }
