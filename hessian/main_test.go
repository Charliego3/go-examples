package main

import (
	"bytes"
	"github.com/whimthen/kits/logger"
	"golang.org/x/text/encoding/charmap"
	"testing"
)
import "reflect"
import "math"
import "encoding/json"

func TestSlice(t *testing.T) {

	arr := make([]byte, 0, 10)
	t.Logf("slice is %v\n", arr[:10])
	changeSlice(arr[:0])
	t.Logf("slice is %v\n", arr[:10])

	str := "Ljava.lang.String;"
	switch str {
	case "Ljava.lang.String;":
		t.Logf("str is Ljava.lang.String;\n")
	default:
		t.Logf("str is not Ljava.lang.String;\n")
	}

}

func changeSlice(arr []byte) {
	arr = append(arr, 0x01)
}

func TestJavaTcObject(t *testing.T) {
	// var f *os.File
	var err error

	// if f, err = os.Open("d:\\tmp\\serialize-child.data"); err != nil {
	// 	t.Fatalf("got error when open file %v\n", err)
	// }
	// defer f.Close()

	iso8859_1_buf := []byte("%C2%AC%C3%AD%00%05sr%00%17com.tenstar.api.TestObj%C2%B5%18%C2%822%C2%80%5B%C2%91%C3%A7%02%00%03I%00%03ageL%00%06amountt%00%16Ljava%2Fmath%2FBigDecimal%3BL%00%04namet%00%12Ljava%2Flang%2FString%3Bxp%00%00%00%0Csr%00%14java.math.BigDecimalT%C3%87%15W%C3%B9%C2%81%28O%03%00%02I%00%05scaleL%00%06intValt%00%16Ljava%2Fmath%2FBigInteger%3Bxr%00%10java.lang.Number%C2%86%C2%AC%C2%95%1D%0B%C2%94%C3%A0%C2%8B%02%00%00xp%00%00%00%00sr%00%14java.math.BigInteger%C2%8C%C3%BC%C2%9F%1F%C2%A9%3B%C3%BB%1D%03%00%06I%00%08bitCountI%00%09bitLengthI%00%13firstNonzeroByteNumI%00%0ClowestSetBitI%00%06signum%5B%00%09magnitudet%00%02%5BBxq%00%7E%00%06%C3%BF%C3%BF%C3%BF%C3%BF%C3%BF%C3%BF%C3%BF%C3%BF%C3%BF%C3%BF%C3%BF%C3%BE%C3%BF%C3%BF%C3%BF%C3%BE%00%00%00%01ur%00%02%5BB%C2%AC%C3%B3%17%C3%B8%06%08T%C3%A0%02%00%00xp%00%00%00%01%01xxt%00%04name")
	buf := make([]rune, len(iso8859_1_buf))
	for i, b := range iso8859_1_buf {
		buf[i] = rune(b)
	}

	bs, err := charmap.ISO8859_1.NewDecoder().Bytes([]byte(string(buf)))
	if err != nil {
		logger.Fatal(err)
	}
	reader := bytes.NewReader(bs)

	arr := make([]byte, 1<<7) // 128

	if _, err = reader.Read(arr[:4]); err != nil {
		// to be continued...
		t.Fatalf("Got error %v\n", err)
	}
	refs := make([]*JavaReferenceObject, 1<<7)
	jo := &JavaTcObject{}
	if err = jo.Deserialize(reader, refs); err != nil {
		t.Fatalf("When deserialize JavaTcObject got %v\n", err)
	}
	t.Logf("Got Tc_OBJECT %v\n", jo)
	rv := reflect.ValueOf(jo)
	if rv.Kind() == reflect.Ptr {
		t.Logf("jo type is %s\n", rv.Elem().Type().Name())
	} else {
		t.Logf("jo type is %s\n", rv.Type().Name())
	}

}
func TestJavaDeserialize(t *testing.T) {

	// var f *os.File
	var err error

	// if f, err = os.Open("d:\\tmp\\serialize-child.data"); err != nil {
	// 	t.Fatalf("got error when open file %v\n", err)
	// }
	// defer f.Close()

	bs, err := charmap.ISO8859_1.NewDecoder().Bytes([]byte("%C2%AC%C3%AD%00%05sr%00%17com.tenstar.api.TestObj%C2%B5%18%C2%822%C2%80%5B%C2%91%C3%A7%02%00%03I%00%03ageL%00%06amountt%00%16Ljava%2Fmath%2FBigDecimal%3BL%00%04namet%00%12Ljava%2Flang%2FString%3Bxp%00%00%00%0Csr%00%14java.math.BigDecimalT%C3%87%15W%C3%B9%C2%81%28O%03%00%02I%00%05scaleL%00%06intValt%00%16Ljava%2Fmath%2FBigInteger%3Bxr%00%10java.lang.Number%C2%86%C2%AC%C2%95%1D%0B%C2%94%C3%A0%C2%8B%02%00%00xp%00%00%00%00sr%00%14java.math.BigInteger%C2%8C%C3%BC%C2%9F%1F%C2%A9%3B%C3%BB%1D%03%00%06I%00%08bitCountI%00%09bitLengthI%00%13firstNonzeroByteNumI%00%0ClowestSetBitI%00%06signum%5B%00%09magnitudet%00%02%5BBxq%00%7E%00%06%C3%BF%C3%BF%C3%BF%C3%BF%C3%BF%C3%BF%C3%BF%C3%BF%C3%BF%C3%BF%C3%BF%C3%BE%C3%BF%C3%BF%C3%BF%C3%BE%00%00%00%01ur%00%02%5BB%C2%AC%C3%B3%17%C3%B8%06%08T%C3%A0%02%00%00xp%00%00%00%01%01xxt%00%04name"))
	if err != nil {
		logger.Fatal(err)
	}
	reader := bytes.NewReader(bs)

	var v JavaSerializer

	if v, err = DeserializeStream(reader); err != nil {
		t.Fatalf("When Deserialize stream, got error %v\n", err)
	} else {
		if bs, err := json.MarshalIndent(v.JsonMap(), "", "  "); err != nil {
			t.Fatalf("Json error %v\n", err)
		} else {
			t.Logf("Deserialize stream got\n %s\n", string(bs))
		}
		// t.Logf("Deserialize stream got %v\n", v.JsonMap())
	}

}

// TestLong it's ok
// we must declare it to hold the number before using it
// davidwang2006@aliyun.com 2018-02-01 11:17:24
func TestLong(t *testing.T) {
	var it int64 = -3665804199014368530
	var uit uint64 = uint64(it)
	var it2 int64 = int64(uit)
	t.Logf("%d %x %d", it, uit, it2)
	var i32 int32 = 0x3f400000
	var f32 float32 = math.Float32frombits(uint32(i32)) // float32(i32)
	t.Logf("%d %f", i32, f32)
	buff := make([]byte, 0, 4)
	t.Logf("buff is %v\n", buff)
	t.Logf("buff is %v\n", buff[:4]) // 4 is okay
}

// TestObjectSerialize test object serialize object to stream
func TestObjectSerialize(t *testing.T) {

	jo := NewJavaTcObject(1)
	clz := NewJavaTcClassDesc("com.tenstar.api.TestObj", 5397420999990079001, SC_SERIALIZABLE)
	jfa := NewJavaField(TC_STRING, "name", "name")
	jfa.FieldObjectClassName = "java.lang.String"
	jfb := NewJavaField(TC_PRIM_INTEGER, "age", 12)
	// jfc := NewJavaField(TC_OBJ_OBJECT, "amount", decimal.NewFromInt(1))
	// jfc.FieldObjectClassName = "java.math.BigDecimal"
	clz.AddField(jfa)
	clz.AddField(jfb)
	// clz.AddField(jfc)
	clz.SortFields()

	jo.AddClassDesc(clz)

	// var f *os.File
	var err error

	// if f, err = os.OpenFile("serialize-go.data", os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0755); err != nil {
	// 	t.Fatalf("got error when open file %v\n", err)
	// }
	// defer f.Close()

	buffer := bytes.NewBuffer(nil)
	if err = SerializeJavaEntity(buffer, jo); err != nil {
		t.Fatalf("SerializeJavaEntity got %v\n", err)
	} else {
		t.Logf("%s", buffer.String())
	}
}
