WSContoller = function(host,message,close){
    this.conn = undefined;
    if (window["WebSocket"]) {
        this.conn = new WebSocket(host);
        this.setOnmessage(message)
        this.setClose(close)
    }
}
WSContoller.prototype.setClose = function(funcClose){
    this.conn.onclose = funcClose;
}
WSContoller.prototype.setOnmessage = function(funcMessage){
    this.conn.onmessage = funcMessage;
}
WSContoller.prototype.send = function(msg){
    this.conn.send(msg);
}
console.log("ws_controller.js");