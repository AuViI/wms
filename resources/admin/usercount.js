"use strict"

function get_user_count(callback) {
    var ajax = new XMLHttpRequest();
    ajax.onreadystatechange = function() {
        if (ajax.readyState != 4) {
            return;
        }
        if (ajax.status != 200) {
            return;
        }
        callback(JSON.parse(ajax.responseText).length);
    }
    ajax.open("POST", "/suser/");
    ajax.setRequestHeader("Content-Type", "application/json");
    ajax.send(JSON.stringify({type:"get"}));
}

function update_usercount_badges() {
    // TODO talk to server about users
    var containers = [];
    var container_names = ["usercount-nav"];

    for (let i = 0; i < container_names.length; i++) {
        let c = document.getElementById(container_names[i]);
        if (c != undefined) {
            containers.push(c);
        }
    }

	get_user_count((usercount)=>{ 
		for (let i = 0; i < containers.length; i++) {
			let c = containers[i];
			let s = document.createElement("span");
			let t = document.createTextNode(usercount);
			s.appendChild(t);
			s.classList.add("badge");
			s.classList.add("badge-secondary");
			c.appendChild(s);
		}
	});
}
