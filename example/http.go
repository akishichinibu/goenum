package enum

// godantic:enum
type _E_HTTPSuccess interface {
	Ok200() (_200 int)
	Created201() (_201 int)
	Accepted202() (_202 int)
	NoContent204() (_204 int)
	PartialContent206() (_206 int)
	BadRequest400() (_400 int)
}
