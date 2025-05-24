package enum

import (
	"fmt"
	"testing"
)

func TestColor(t *testing.T) {
	t.Run("Gray", func(t *testing.T) {
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
	})
}
