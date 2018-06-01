BoxAction = function(canvas){
    this.canvas = canvas;
    this.canvas.setAttribute('tabindex', 0);
    this.context = canvas.getContext('2d');
    this.init();
}
BoxAction.prototype.init = function(){
    this.pl = {
        cx:100,
        cy:100,
        cw:100,
        ch:100
    }
}
BoxAction.prototype.setEvent = function(key,callback){
    canvas.addEventListener(key, callback, false);
}
BoxAction.prototype.setSampleEvent = function(){
    this.setEvent('down',function (e) { console.log("down"); })
    this.setEvent('up',function (e) { console.log("up"); })
    this.setEvent('click',function (e) { console.log("click"); })
    this.setEvent('mouseover',function (e) { console.log("mouseover"); })
    this.setEvent('mouseout',function (e) { console.log("mouseout"); })
    this.setEvent('down',function (e) { console.log("down"); })
    this.setEvent('down',function (e) { console.log("down"); })
    this.setEvent('keydown',function (e) { console.log("key down"); })
    this.setEvent('keyup',function (e) { console.log("key down"); })
}

BoxAction.prototype.drawRect = function(x, y, width, height) {
    this.context.clearRect(0, 0, this.canvas.width, this.canvas.height);
    this.context.fillRect(pl.x, pl.y, pl.width, pl.height);
    this.boxs.forEach(function(b){
        this.context.fillRect(b.x, b.y, b.width, b.height); 
    })
}
BoxAction.prototype.moveRect = function(x, y){
    this.pl.cx = x*5;
    this.pl.cy = y*5;
}


Box = function(){
    this.cx = 100;
    this.cy = 100;
    this.width = 100;
    this.height = 100;
}
Box.prototype.move = function(x,y){
    this.box.cx = x*5;
    this.box.cy = y*5;
}