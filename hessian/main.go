package main

// java serialize & deserialize

import (
	"fmt"
	hessian "github.com/apache/dubbo-go-hessian2"
	dec "github.com/dubbogo/gost/math/big"
)

type TestObj struct {
	Name   string
	Age    int
	Amount *dec.Decimal
}

func (o TestObj) JavaClassName() string {
	return "com.tenstar.api.TestObj"
}

func main() {
	obj := TestObj{
		Name:   "name",
		Age:    12,
		Amount: dec.NewDecFromInt(1),
	}

	encoder := hessian.NewEncoder()
	err := encoder.Encode(obj)
	if err != nil {
		panic(err)
	}

	println(string(encoder.Buffer()))

	// -----------
	serialized := []byte("��\u0000\u0005sr\u0000\aTestObj�\u0018�2�[��\u0002\u0000\u0003I\u0000\u0003ageL\u0000\u0006amountt\u0000\u0016Ljava/math/BigDecimal;L\u0000\u0004namet\u0000\u0012Ljava/lang/String;xp\u0000\u0000\u0000\fsr\u0000\u0014java.math.BigDecimalT�\u0015W��(O\u0003\u0000\u0002I\u0000\u0005scaleL\u0000\u0006intValt\u0000\u0016Ljava/math/BigInteger;xr\u0000\u0010java.lang.Number���\u001D\v���\u0002\u0000\u0000xp\u0000\u0000\u0000\u0000sr\u0000\u0014java.math.BigInteger���\u001F�;�\u001D\u0003\u0000\u0006IbitCountI\u0000\tbitLengthI\u0000\u0013firstNonzeroByteNumI\u0000\flowestSetBitI\u0000\u0006signum[\u0000\tmagnitudet\u0000\u0002[Bxq\u0000~\u0000\u0006����������������\u0000\u0000\u0000\u0001ur\u0000\u0002[B��\u0017�T�\u0002\u0000\u0000xp\u0000\u0000\u0000\u0001\u0001xxt\u0000\u0004name")
	dobj, err := hessian.NewDecoder(serialized).Decode()
	if err != nil {
		panic(err)
	}

	v, ok := dobj.(*TestObj)
	if ok {
		fmt.Printf("Name: %s, Age: %d, Amount: %s\n", v.Name, v.Age, v.Amount.String())
	}
	println(dobj)
}
