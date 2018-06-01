WSContoller = function(host,message,close){
    this.conn = undefined;
    this.host = host;
    if (window["WebSocket"]) {
        this.conn = new WebSocket(host);
        this.setOnmessage(message)
        this.setClose(close)
    }
}
WSContoller.prototype.setClose = function(funcClose){
    if(funcClose !== undefined){
        this.conn.onclose = funcClose;
        this.onclose = funcClose;
    }
    else{
        this.conn.onclose = this.onclose;
    }
}
WSContoller.prototype.setOnmessage = function(funcMessage){
    if(funcMessage !== undefined){
        this.conn.onmessage = funcMessage;
        this.onmessage = funcMessage;
    }
    else{
        this.conn.onmessage = this.onmessage;
    }
}
WSContoller.prototype.send = function(msg){
    console.log("host : "+this.host);
    this.conn.send(msg);
}
WSContoller.prototype.reConnect = function(){
    var self=this;
    console.log("re connect ....");
    setTimeout(function(){
        self.conn = new WebSocket(self.host);
        self.setOnmessage();
        self.setClose();
        console.log("re connected");
    },1000)
}
console.log("ws_controller.js");