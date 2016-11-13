//var author = {name: "shresh", books: ["Chronicles of 512", "The White Queen", "The Dark King"]};

var ws = new WebSocket("ws://" + location.host + "/ws");

ws.addEventListener("message", function(e) {console.log(e.data);});
//setTimeout(function(){ws.send(JSON.stringify(author));}, 5000);

