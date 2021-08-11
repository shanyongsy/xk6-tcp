function OnGetTime(tcp){
    tcp.writeStrLn('/time')
    return{code: true}
}

function OnNormalFunc(tcp, data){
    tcp.writeStrLn(data)
    return{code: true}
}

export{
    OnGetTime,
    OnNormalFunc,
}