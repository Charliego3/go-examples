package widgets

type NSLayoutAttribute uint
type NSLayoutRelation uint

const (
	NSLayoutAttributeNotAnAttribute NSLayoutAttribute = iota
	NSLayoutAttributeLeft
	NSLayoutAttributeRight
	NSLayoutAttributeTop
	NSLayoutAttributeBottom
	NSLayoutAttributeLeading
	NSLayoutAttributeTrailing
	NSLayoutAttributeWidth
	NSLayoutAttributeHeight
	NSLayoutAttributeCenterX
	NSLayoutAttributeCenterY
	NSLayoutAttributeLastBaseline
	NSLayoutAttributeFirstBaseline
	NSLayoutAttributeLeftMargin
	NSLayoutAttributeRightMargin
	NSLayoutAttributeTopMargin
	NSLayoutAttributeBottomMargin
	NSLayoutAttributeLeadingMargin
	NSLayoutAttributeTrailingMargin
	NSLayoutAttributeCenterXWithinMargins
	NSLayoutAttributeCenterYWithinMargins
)

const (
	//NSLayoutRelationLessThanOrEqual NSLayoutRelation = iota - 1
	NSLayoutRelationEqual              NSLayoutRelation = 0
	NSLayoutRelationGreaterThanOrEqual                  = 1
)
