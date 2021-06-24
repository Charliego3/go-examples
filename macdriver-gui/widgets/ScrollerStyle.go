package widgets

// NSScrollerStyle the scroll style
// https://developer.apple.com/documentation/appkit/nsscroller/style
type NSScrollerStyle uint

const (
	Legacy  NSScrollerStyle = iota // Specifies legacy-style scrollers as prior to macOS 10.7.
	Overlay                        // Specifies overlay-style scrollers in macOS 10.7 and later.
)
