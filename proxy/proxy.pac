function matched(url, host) {
    return %s;
}

function FindProxyForURL(url, host) {
    if (matched(url, host)) {
        return "PROXY %s:8081; SOCKS %s:8081; DIRECT";
    }
    return "DIRECT";
}
