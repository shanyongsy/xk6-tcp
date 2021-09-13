function C2GVerifyMessage(tcp) {
    tcp.c2GVerifyMessage()
    return { code: true }
}

function C2GTestMessage(tcp) {
    tcp.c2GTestMessage()
    return { code: true }
}

function C2GHeartbeatMessage(tcp) {
    tcp.c2GHeartbeatMessage()
    return { code: true }
}

export {
    C2GVerifyMessage,
    C2GTestMessage,
    C2GHeartbeatMessage,
}
