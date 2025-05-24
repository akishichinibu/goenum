# goenum

> Type-safe enum generator for Go

## 🧪 Quick Example

Please have a check at [example]("./example"). 

### 📚 Define an enum interface:

```go
// example/color.go
type _E_Color interface {
    Gray(alpha float64) (_Alpha string)
    Mix(r, g, b int) (_RGB string)
    Red()
}
```

### ⚙️ Usage & API
```go
color := EnumColor.Gray(0.5)
color.Match(
	func(v EnumColorVariantGray) {
		fmt.Println("Gray with alpha:", v.GetAlpha())
	},
	func(v EnumColorVariantMix) {
		t.Error("Expected Gray, got Mix")
	},
	func(v EnumColorVariantRed) {
		t.Error("Expected Gray, got Red")
	},
)
```

# 📝 License

MIT
