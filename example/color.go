package enum

// godantic:enum
type _E_Color interface {
	Gray(alpha float64)
	Mix(r, g, b int)
	Red()
}
