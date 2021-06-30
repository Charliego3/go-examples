package widgets

// NSScrollerStyle the scroll style
// https://developer.apple.com/documentation/appkit/nsscroller/style
type NSScrollerStyle uint

const (
	NSScrollerStyleLegacy  NSScrollerStyle = iota // Specifies legacy-style scrollers as prior to macOS 10.7.
	NSScrollerStyleOverlay                        // Specifies overlay-style scrollers in macOS 10.7 and later.
)
