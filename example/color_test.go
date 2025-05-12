package enum

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEnumColor(t *testing.T) {
	var red = EnumColor.Red()
	var deepGray = EnumColor.Gray(0.7)
	var lightGray = EnumColor.Gray(0.3)
	var deepGray2 = EnumColor.Gray(0.7)

	t.Log(red)
	t.Log(deepGray)
	t.Log(lightGray)

	require.False(t, red.Equal(deepGray))
	require.True(t, red.Equal(red))
	require.True(t, deepGray.Equal(deepGray2))
	require.False(t, deepGray.Equal(lightGray))

	// branch match
	var v EnumColorVariant = red
	v.Match(
		func(v EnumColorVariantGray) {
			fmt.Println("Gray", v.GetAlpha())
		},
		func(v EnumColorVariantMix) {
			fmt.Println("Mix", v.GetR(), v.GetG(), v.GetB())
		},
		func(v EnumColorVariantRed) {
			fmt.Println("Red")
		},
	)
}
