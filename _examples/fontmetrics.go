package main

import (
	"fmt"

	. "modernc.org/tk9.0"
	_ "modernc.org/tk9.0/themes/azure"
)

func main() {
	for _, family := range []string{
		DefaultFont, TextFont, FixedFont, MenuFont, HeadingFont,
		CaptionFont, SmallCaptionFont, IconFont, TooltipFont,
		CourierFont(),
	} {
		font := NewFont(Family(family))
		fmt.Printf("%20s : ascent=%d, descent=%d, linespace=%d, fixed=%t\n",
			family,
			font.MetricsAscent(App),
			font.MetricsDescent(App),
			font.MetricsLinespace(App),
			font.MetricsFixed(App),
		)
	}
}
