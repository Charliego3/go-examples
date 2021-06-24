package button

type BezelStyle uint

const (
	_                       BezelStyle = iota
	BezelStyleRounded                  // 1
	BezelStyleRegularSquare            // 2
	_
	_
	BezelStyleDisclosure        // 5
	BezelStyleShadowlessSquare  // 6
	BezelStyleCircular          // 7
	BezelStyleTexturedSquare    // 8
	BezelStyleHelpButton        // 9
	BezelStyleSmallSquare       // 10
	BezelStyleTexturedRounded   // 11
	BezelStyleRoundRect         // 12
	BezelStyleRecessed          // 13
	BezelStyleRoundedDisclosure // 14
	BezelStyleInline            // 15
)
