package main

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"regexp"
	"testing"
)

func TestRegex(t *testing.T) {
	var re = regexp.MustCompile(`(?m)((13[0-9])|(14[5-9])|(15([0-3]|[5-9]))|(16[6-7])|(17[1-8])|(18[0-9])|(19[1|3])|(19[5|6])|(19[8|9]))\d{8}`)
	var str = `[房屋出租]
苏果超市海安恒天店一楼招商


[复制链接]

1回复
99查看










 

苏果超市恒天店








苏果超市恒天店
当前在线

UID537105
积分5
海币7 HB
好评0 点
金币0 个
活跃度17 ℃
精华0
注册时间2021-9-15
最后登录2021-9-16
在线时间0 小时









     
                
1主题2帖子5积分

新手上路



新手上路, 积分 5, 距离下一级还需 5 积分



新手上路, 积分 5, 距离下一级还需 5 积分


        积分5


  
  串个门
  发消息
加好友
打招呼
  
扫一扫，手机访问本帖

        




1楼


电梯直达







发表于 2021-9-16 20:39
|
只看该作者
|只看大图

|倒序浏览
|阅读模式


.pcb{margin-right:0}
 





马上注册，结交更多好友，享用更多功能，让你轻松玩转社区。

您需要 登录 才可以下载或查看，没有帐号？注册 



 

x

苏果超市海安恒天店，一楼招商，11-158平米，经营项目：儿童、男女服饰、箱包饰品、服务配套，餐饮，教育等
联系方式：徐女士 15996614832
地址：曙光东二路七号












4.jpg
(130.63 KB, 下载次数: 0)




下载附件



2021-9-16 20:38 上传









餐饮






















3.jpg
(103.99 KB, 下载次数: 0)




下载附件



2021-9-16 20:38 上传









服饰






















2.jpg
(183.66 KB, 下载次数: 0)




下载附件



2021-9-16 20:38 上传









教育






















1.jpg
(262.68 KB, 下载次数: 0)




下载附件



2021-9-16 20:38 上传









服务














                











分享到: 

QQ好友和群

收藏0
支持0
反对0










回复


举报







aimgcount[6290688] = ['1227118','1227117','1227116','1227115'];
attachimggroup(6290688);
var aimgfid = 0;





 
                                                

苏果超市恒天店








苏果超市恒天店
当前在线

UID537105
积分5
海币7 HB
好评0 点
金币0 个
活跃度17 ℃
精华0
注册时间2021-9-15
最后登录2021-9-16
在线时间0 小时









     
                
1主题2帖子5积分

新手上路



新手上路, 积分 5, 距离下一级还需 5 积分



新手上路, 积分 5, 距离下一级还需 5 积分


        积分5


  
  串个门
  发消息
加好友
打招呼
  


        




2楼





 楼主|
发表于 2021-9-16 20:39
|
只看该作者





沙发
                


















回复


举报








`

	for i, match := range re.FindAllString(str, -1) {
		fmt.Println(match, "found at index", i)
	}
}

func TestDb(t *testing.T) {
	db, err := sqlx.Open("postgres", "postgres://charlie:root@127.0.0.1:5432/charlie?sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	row := db.QueryRowx("SELECT 1 FROM phone_numbers WHERE phone = $1", 15996614832)
	if row.Err() != nil {
		t.Fatal(row.Err())
	}
	var phone string
	err = row.Scan(&phone)
	if err != nil {
		t.Fatal(err)
	}
	println(phone)
}

func TestReplace(t *testing.T) {
}
