var SYMBOLSET = {};

SYMBOLSET.baseURL = window.location.origin;

// From http://stackoverflow.com/a/8809472
SYMBOLSET.generateUUID = function() {
    var d = new Date().getTime();
    if(window.performance && typeof window.performance.now === "function"){
        d += performance.now(); //use high-precision timer if available
    }
    var uuid = 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function(c) {
        var r = (d + Math.random()*16)%16 | 0;
        d = Math.floor(d/16);
        return (c=='x' ? r : (r&0x3|0x8)).toString(16);
    });
    return uuid;
};


SYMBOLSET.UploadFileModel = function () {
    var self = this; 
    
    
    self.uuid = SYMBOLSET.generateUUID();

    self.serverMessage = ko.observable("_");
    
    self.maxMessages = 10;
    self.serverMessages = ko.observableArray();
    
    self.connectWebSock = function() {
	var zock = new WebSocket(SYMBOLSET.baseURL.replace("http://", "ws://") + "/websockreg" );
	zock.onopen = function() {
	    console.log("SYMBOLSET.connectWebSock: sending uuid over zock: "+ self.uuid);
	    zock.send("CLIENT_ID: "+ self.uuid);
	};
	zock.onmessage = function(e) {
	    // Just drop the keepalive message
	    if(e.data === "WS_KEEPALIVE") {
		// var d = new Date();
		// var h = d.getHours();
		// var m = d.getMinutes();
		// var s = d.getSeconds();
		// var msg = "Websocket keepalive recieved "+ h +":"+ m +":"+ s;
		// self.serverMessage(msg);
	    }
	    else {
		//console.log("Websocket got: "+ e.data)
		self.serverMessage(e.data);
	    };
	};
	zock.onerror = function(e){
	    console.log("websocket error: " + e.data);
	};
	zock.onclose = function (e) {
	    console.log("websocket got close event: "+ e.code)
	};
    };
    
    self.message = ko.observable("");
    
    self.selectedFile = ko.observable(null);
    self.hasSelectedFile = ko.observable(false);   
    
    self.setSelectedFile = function(symbolsetFile) {
	self.selectedFile(symbolsetFile);
	console.log("selected file: ", self.selectedFile())
	self.hasSelectedFile(true);
    }

    self.uploadFile = function() {
	console.log("uploading file: ", self.selectedFile())
	var url = SYMBOLSET.baseURL + "/symbolset/upload"
	var xhr = new XMLHttpRequest();
	var fd = new FormData();
	xhr.open("POST", url, true);
	xhr.onreadystatechange = function() {
            if (xhr.readyState === 4 && xhr.status === 200) {
		// Everything ok, file uploaded
		console.log("uploadFile returned response text ", xhr.responseText); // handle response.
		self.message("Upload completed without errors");
	    } else {
		self.message("Upload failed: " + xhr.responseText);
	    };
	};
	fd.append("client_uuid", self.uuid);
	fd.append("upload_file", self.selectedFile());
	xhr.send(fd);
    };
    
};

var upload = new SYMBOLSET.UploadFileModel();
ko.applyBindings(upload);
upload.connectWebSock();

console.log("UUID: "+ upload.uuid);
