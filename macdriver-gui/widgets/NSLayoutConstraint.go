package widgets

import (
	"github.com/progrium/macdriver/core"
	"github.com/progrium/macdriver/objc"
)

type NSLayoutConstraint struct {
	objc.Object
}

func NewNSLayoutConstraintWithAttr(subView objc.Object, subAttribute NSLayoutAttribute, relation NSLayoutRelation,
	toItem objc.Object, toAttribute NSLayoutAttribute, multiplier float32, constant float32) NSLayoutConstraint {
	return NSLayoutConstraint{objc.Get("NSLayoutConstraint").Alloc().
		Send("constraintWithItem:attribute:relatedBy:toItem:attribute:multiplier:constant:",
			subView, subAttribute, relation, toItem, toAttribute, multiplier, constant)}
}

func NewNSLayoutConstraint() NSLayoutConstraint {
	return NSLayoutConstraint{Object: objc.Get("NSLayoutConstraint").Alloc().Init()}
}

func (c NSLayoutConstraint) SetConstraintWithItem(subView objc.Object, subAttribute NSLayoutAttribute, relation NSLayoutRelation,
	toItem objc.Object, toAttribute NSLayoutAttribute, multiplier float64, constant float64) {
	c.Set("constraintWithItem:attribute:relatedBy:toItem:attribute:multiplier:constant:",
		&subView, core.NSUInteger(subAttribute), core.NSUInteger(relation), &toItem, core.NSUInteger(toAttribute), multiplier, constant,
	)
}

//- (void)viewDidLoad {
//[super viewDidLoad];
//// Do any additional setup after loading the view, typically from a nib.
//self.view.backgroundColor = [UIColor yellowColor];
//
//
//UIView *subView = [[UIView alloc] init];
//subView.backgroundColor = [UIColor redColor];
//// 在设置约束前，先将子视图添加进来
//[self.view addSubview:subView];
//
//// 使用autoLayout约束，禁止将AutoresizingMask转换为约束
//[subView setTranslatesAutoresizingMaskIntoConstraints:NO];
//
//// 设置subView相对于VIEW的上左下右各40像素
//NSLayoutConstraint *constraint1 = [NSLayoutConstraint constraintWithItem:subView attribute:NSLayoutAttributeTop relatedBy:NSLayoutRelationEqual toItem:self.view attribute:NSLayoutAttributeTop multiplier:1.0 constant:40];
//NSLayoutConstraint *constraint2 = [NSLayoutConstraint constraintWithItem:subView attribute:NSLayoutAttributeLeft relatedBy:NSLayoutRelationEqual toItem:self.view attribute:NSLayoutAttributeLeft multiplier:1.0 constant:40];
//// 由于iOS坐标系的原点在左上角，所以设置右边距使用负值
//NSLayoutConstraint *constraint3 = [NSLayoutConstraint constraintWithItem:subView attribute:NSLayoutAttributeBottom relatedBy:NSLayoutRelationEqual toItem:self.view attribute:NSLayoutAttributeBottom multiplier:1.0 constant:-40];
//
//// 由于iOS坐标系的原点在左上角，所以设置下边距使用负值
//NSLayoutConstraint *constraint4 = [NSLayoutConstraint constraintWithItem:subView attribute:NSLayoutAttributeRight relatedBy:NSLayoutRelationEqual toItem:self.view attribute:NSLayoutAttributeRight multiplier:1.0 constant:-40];
//
//// 将四条约束加进数组中
//NSArray *array = [NSArray arrayWithObjects:constraint1, constraint2, constraint3, constraint4 ,nil];
//// 把约束条件设置到父视图的Contraints中
//[self.view addConstraints:array];
//}
