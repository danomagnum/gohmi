var DEBUG = 0;
var DEBUG_NOREFRESH = 0;


function update_element_bulk(element, response){
	// This function is like the other update element function but
	// the response here will have the driver/ at the beginning of the tag path.
	// The goal is to have less requests to the base server by combining them.
	var tagpath = element.getAttribute("data-target");

	var tagname = tagpath;
	var newdata =  response[tagname];

	if (newdata == undefined){
		return;
	};

	if (element.hasAttribute("data-format")){
		let format = element.getAttribute("data-format");
		try{
			newdata = sprintf(format, response[tagpath]);
		} catch (error) {
  			console.error(error);
		}
	}

	if (element.hasAttribute("value")){
		if (document.activeElement == element) return;
		element.value = newdata;
	}else if(element.hasAttribute("data-value")){
		element.setAttribute("data-value", newdata);
	}else{
		element.innerHTML = newdata;
	}
	if (element.classList.contains("multistate_indicator")){
		update_indicator(element);
	}
		

	if (element.hasAttribute("onchange")){
		element.onchange();
	}

};

function refresh_all_data(response){
	var elements = document.getElementsByClassName("autoload_value");

	[].forEach.call(elements, function(element) {
		update_element_bulk(element, response);
	});
};



function get_all_datapoints(){
	var elements = document.getElementsByClassName("autoload_value");
	paths = [];
	[].forEach.call(elements, function(element){
		path = element.getAttribute("data-target");
		if (path !== null){
			paths.push(path);
		};
	});
	//for (element in elements){
		//e = elements[element];
		//paths.push(elements[element].getAttribute("data-target"));
	//}

	var xmlhttp = new XMLHttpRequest();
	xmlhttp.onreadystatechange=function() {
		if (this.readyState == 4 && this.status == 200) {
			var response = JSON.parse(this.responseText);
			refresh_all_data(response);
		}
	};

	xmlhttp.overrideMimeType("application/json");
	if (paths.length > 0){
		var path =  "/read_multi/"
		xmlhttp.open("POST", path);
		xmlhttp.setRequestHeader("Content-Type", "application/json;charset=UTF-8");
		xmlhttp.send(JSON.stringify(paths));
	}
};


function update_indicator(element){
	var sel = element.selectedIndex;
	if (sel == -1){
		//Something went wrong with the index
		return
	}
	var newclass = element.options[element.selectedIndex].className;

	//element.className = "multistate_indicator";
	element.className = element.getAttribute("data-baseclass");

	if (newclass){
		element.classList.add(newclass);
	};
}


function shiftContext(ctx, w, h, dx, dy) {
	var clamp = function(high, value) { return Math.max(0, Math.min(high, value)); };
	var imageData = ctx.getImageData(clamp(w, -dx), clamp(h, -dy), clamp(w, w-dx), clamp(h, h-dy));
	ctx.clearRect(0, 0, w, h);
	ctx.putImageData(imageData, 0, 0);
}





function setup_indicators(){
	var elements = document.getElementsByClassName("multistate_indicator");
	[].forEach.call(elements, function(element) {
		/* logic here */
		element.setAttribute("onchange", "update_indicator(this)");
		element.setAttribute("data-baseclass", element.className);
		update_indicator(element);

		if (!DEBUG){
			element.setAttribute('disabled', 1);
		};
	});

};

function pb_momentary_press(element){
	var xmlhttp = new XMLHttpRequest();
	target = element.getAttribute("data-target");

	var path =  "/write/" + target
	var params = 'value=' + element.getAttribute("data-press");
	xmlhttp.open("POST", path, true);
	xmlhttp.setRequestHeader('content-type', 'application/x-www-form-urlencoded;charset=UTF-8');
	xmlhttp.send(params);
}

function pb_momentary_release(element){
	var xmlhttp = new XMLHttpRequest();
	target = element.getAttribute("data-target");
	var path =  "/write/" + target
	var params = 'value=' + element.getAttribute("data-release");
	xmlhttp.open("POST", path, true);
	xmlhttp.setRequestHeader('content-type', 'application/x-www-form-urlencoded;charset=UTF-8');
	xmlhttp.send(params);
}

function pb_toggle(element){
	target = element.getAttribute("data-target");
	if (element.value == element.getAttribute("data-press")){
		element.value = element.getAttribute("data-release");
	}else{
		element.value = element.getAttribute("data-press");
	}

	var xmlhttp = new XMLHttpRequest();
	var path = "/write/" + target
	var params = 'value=' + element.value;
	xmlhttp.open("POST", path, true);
	xmlhttp.setRequestHeader('content-type', 'application/x-www-form-urlencoded;charset=UTF-8');
	xmlhttp.send(params);
}

function input_submit(element){
	target = element.getAttribute("data-target");
	

	var xmlhttp = new XMLHttpRequest();
	var path = "/write/" + target
	var params = 'value=' + element.value;
	xmlhttp.open("POST", path, true);
	xmlhttp.setRequestHeader('content-type', 'application/x-www-form-urlencoded;charset=UTF-8');
	xmlhttp.send(params);
};

function setup_inputs(){
	var elements = document.getElementsByTagName("input");
	[].forEach.call(elements, function(element) {
		/* logic here */
		/*element.setAttribute("onchange", "input_submit(this)");*/
		element.addEventListener("keyup", event=>{
			if (event.key !== "Enter") return;
			input_submit(element);
		});
	});
};




function setup_buttons(){
	var elements = document.getElementsByClassName("pb_momentary");
	[].forEach.call(elements, function(element) {
		/* logic here */
		if (! element.hasAttribute("data-press")){
			element.setAttribute("data-press", "1");
		}
		if (! element.hasAttribute("data-release")){
			element.setAttribute("data-release", "0");
		}

		element.setAttribute("onmousedown", "pb_momentary_press(this)");
		element.setAttribute("onmouseup", "pb_momentary_release(this)");
	});

	var elements = document.getElementsByClassName("pb_toggle");
	[].forEach.call(elements, function(element) {
		/* logic here */
		if (! element.hasAttribute("data-press")){
			element.setAttribute("data-press", "1");
		}
		if (! element.hasAttribute("data-release")){
			element.setAttribute("data-release", "0");
		}
		if (! element.hasAttribute("value")){
			element.value = 0;
		}

		element.setAttribute("onclick", "pb_toggle(this)");

	});

};


function update_all(){
	get_all_datapoints()
}

function startup(){
	setup_indicators();
	setup_buttons();
	setup_inputs();
	//refresh_data();
	if (! DEBUG_NOREFRESH){
		setInterval(update_all, 500);
	}
}


function round(item, digits){
	item.innerText = item.innerText.slice(0,digits)
}
