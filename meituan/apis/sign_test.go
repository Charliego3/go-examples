package apis

import (
	"crypto/sha1"
	"github.com/whimthen/temp/meituan/configs"
	"testing"
)

// /bizapi/loginv5?appkey=app-pos&bg_source=7&reqtime=1642756506&sign=d08dd73257830b7ed68aa9d770ed110f756260d6&utm_campaign=uisdk1.0&utm_medium=android&utm_term=3.20.600&uuid=6064b84d5dee4adfa9d89418ae1d91c0a164236619454453772
// password=yy12347890&remember_password=1&part_type=1&part_key=6519763&login=18929387993&fingerprint=fingerprint
func TestSign(t *testing.T) {
	// signMap := make(map[string]string)
	// signMap["appkey"] = "app-pos"
	// signMap["bg_source"] = "7"
	// signMap["reqtime"] = "1642756506"
	// signMap["utm_campaign"] = "uisdk1.0"
	// signMap["utm_medium"] = "android"
	// signMap["utm_term"] = "3.20.600"
	// signMap["uuid"] = "6064b84d5dee4adfa9d89418ae1d91c0a164236619454453772"

	// signMap := url.Values{
	// 	"appkey":       []string{"app-pos"},
	// 	"bg_source":    []string{"7"},
	// 	"reqtime":      []string{"1642756506"},
	// 	"utm_campaign": []string{"uisdk1.0"},
	// 	"utm_medium":   []string{"android"},
	// 	"utm_term":     []string{"3.20.600"},
	// 	"uuid":         []string{"6064b84d5dee4adfa9d89418ae1d91c0a164236619454453772"},
	// }
	//
	// // signMap["password"] = "yy12347890"
	// // signMap["remember_password"] = "1"
	// // signMap["part_type"] = "1"
	// // signMap["part_key"] = "6519763"
	// // signMap["login"] = "18929387993"
	// // signMap["fingerprint"] = "fingerprint"
	// // signMap["bizlogintoken"] = "aooWrmy32LthrreCmMmtupkU3kGe_BhJB5BKjkzZn3Tdvk0_yarcqd2H6Wml5s3k_gSW_k2-xRKGeHrHgF85Ug"
	//
	// for key, val := range signMap {
	// 	k := url.QueryEscape(key)
	// 	v := url.QueryEscape(val[0])
	// 	t.Log("Key:", key, "->", k, ", Val:", val, "->", v)
	// 	signMap[k] = v
	// }
	//
	configs.Config.AppSecret = "86951cb05a9072b4d9425bd554edc842b9512035"

	sign := Sign(1)
	t.Log("Tested Sign:", sign)
}

func TestSHA1(t *testing.T) {
	info := "This is info message"

	hash := sha1.New()
	hash.Write([]byte(info))
	infoBytes := hash.Sum(nil)
	t.Log(offset(infoBytes))
}
