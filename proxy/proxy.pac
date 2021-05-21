function FindProxyForURL(url, host) {
    if (shExpMatch(host, "*weixin.qq.com") || shExpMatch(host, "*wx.qlogo.cn") || shExpMatch(host, "mmbiz.qpic.cn") || shExpMatch(host, "*taobao.com")) {
        return "DIRECT";
    }

    return "PROXY loopbackAddress:8081; SOCKS loopbackAddress:8081; DIRECT";
}