package apis

// CategoryResp
// {
//	"code": 200,
//	"data": [
//		{
//			"combos": [],
//			"dishAttrCnt": 0,
//			"dishAttrs": [],
//			"dishSkuDetails": [],
//			"dishSpecCnt": 0,
//			"dishSpecs": [],
//			"dishSpus": [
//				{
//					"canWeigh": false,
//					"cateId": 2361755,
//					"cateName": "小吃",
//					"desc": "",
//					"dishAttrSource": 1,
//					"dishAttrs": [],
//					"dishSkus": [
//						{
//							"id": 23841812,
//							"memberPrice": 1000,
//							"no": "",
//							"price": 1000,
//							"saledOnWaiMai": false,
//							"spec": "份"
//						}
//					],
//					"id": 29332034,
//					"imgUrl": "http://p1.meituan.net/xianfu/15b2969962e9a5f2b9c66fa4578e70f486016.jpg",
//					"minCount": 0,
//					"modifyTime": 1589081652000,
//					"name": "凉皮",
//					"no": "",
//					"saledOnWaiMai": false,
//					"showLimit": false,
//					"sideDishSource": 2,
//					"sideDishes": [],
//					"status": 1,
//					"type": 1,
//					"unit": "份",
//					"waiMaiStatus": 2
//				}
//			],
//			"id": 2361755,
//			"modifyTime": 1632227987000,
//			"name": "小吃",
//			"sideDishCnt": 0,
//			"sideDishes": [],
//			"type": 1
//		}
//	]
// }
type CategoryResp struct {
	Code int        `json:"code"`
	Data []Category `json:"data"`
}

type Category struct {
	Combos         []interface{} `json:"combos"`
	DishAttrCnt    int           `json:"dishAttrCnt"`
	DishAttrs      []interface{} `json:"dishAttrs"`
	DishSkuDetails []interface{} `json:"dishSkuDetails"`
	DishSpecCnt    int           `json:"dishSpecCnt"`
	DishSpecs      []interface{} `json:"dishSpecs"`
	DishSpus       []struct {
		CanWeigh       bool          `json:"canWeigh"`
		CateID         int           `json:"cateId"`
		CateName       string        `json:"cateName"`
		Desc           string        `json:"desc"`
		DishAttrSource int           `json:"dishAttrSource"`
		DishAttrs      []interface{} `json:"dishAttrs"`
		DishSkus       []struct {
			ID            int    `json:"id"`
			MemberPrice   int    `json:"memberPrice"`
			No            string `json:"no"`
			Price         int    `json:"price"`
			SaledOnWaiMai bool   `json:"saledOnWaiMai"`
			Spec          string `json:"spec"`
		} `json:"dishSkus"`
		ID             int           `json:"id"`
		ImgURL         string        `json:"imgUrl"`
		MinCount       int           `json:"minCount"`
		ModifyTime     int64         `json:"modifyTime"`
		Name           string        `json:"name"`
		No             string        `json:"no"`
		SaledOnWaiMai  bool          `json:"saledOnWaiMai"`
		ShowLimit      bool          `json:"showLimit"`
		SideDishSource int           `json:"sideDishSource"`
		SideDishes     []interface{} `json:"sideDishes"`
		Status         int           `json:"status"`
		Type           int           `json:"type"`
		Unit           string        `json:"unit"`
		WaiMaiStatus   int           `json:"waiMaiStatus"`
	} `json:"dishSpus"`
	ID          int           `json:"id"`
	ModifyTime  int64         `json:"modifyTime"`
	Name        string        `json:"name"`
	SideDishCnt int           `json:"sideDishCnt"`
	SideDishes  []interface{} `json:"sideDishes"`
	Type        int           `json:"type"`
}

//  {"data":{"bizacctid":48888690,"login":"18929387993","part_type":1,"part_key":"6519763","loginSensitive":0,"nameSensitive":0,"contactSensitive":0,"access_token":"w9-eESB2lN9kb0sXnczkVFLzvz1Zsvi7jWNw0XiO5OiYxkN6pbsuE_8yUNVTnxkS4WDGNmR0umB2t7EP0KO3iA","refresh_token":"bl8sOyuPOJ2HgRlUwW0h62Dk_ePPmt8JP69u7eTi-vlBynBYjbovf4Nn3Y_ggpVzGVtiItuMz0WGRVJEd80APQ","expire_in":2592000,"refresh_in":604800}}

// AllCategory requestURL:
// https://cloud-erp.meituan.com/bossapi/api/dish/v2/spu/getAll?date=20220121&deviceId=1164252
// &operatorId=48888690&poiId=163181252&s=pos_boss&tenantId=6519763
// &v=3.20.600%3A3200600&token=1642756464&sign=JR1w%2F7NPzT9cVtsc5P9SMQzLL%2Bc%3D
func AllCategory() {

}
