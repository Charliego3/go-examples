function directMatch(url, host) {
    return %s;
}

function FindProxyForURL(url, host) {
    if (directMatch(url, host)) {
        return "DIRECT";
    }

    return "PROXY %s:8081; SOCKS %s:8081; DIRECT";
}