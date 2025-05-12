package enum

// godantic:value_enum
type _G_HTTPSuccess interface {
	Ok200() int
	Created201() int
	Accepted202() int
	NoContent204() int
	PartialContent206() int
}
