"use strict"

function send_request(payload, callback) {
    var ajax = new XMLHttpRequest();
    ajax.onreadystatechange = function() {
        if (ajax.readyState != 4) {
            return;
        }
        if (ajax.status != 200) {
            return;
        }
        callback(JSON.parse(ajax.responseText));
    }
    ajax.open("POST", "/suser/");
    ajax.setRequestHeader("Content-Type", "application/json");
    ajax.send(JSON.stringify(payload));
}

function _col_to_rgb(col) {
	return "rgb("+col.r+","+col.g+","+col.b+")";
}

function _theme_to_grad(theme, direction="to right") {
	return "linear-gradient(" + direction + ", "
		+ _col_to_rgb(theme.start) + " 0%, "
		+ _col_to_rgb(theme.end) + " 100%)";
}

function Users() {
	this.users = [];
	this.add_user = (id, name, theme) => {
		this.users.push({id:id, name:name, theme:theme});
	};
	this.toHTMLs = (i) => {
		let t = this.users[i];
		let li = document.createElement("div");
		li.classList.add("list-group-item");
		let row = document.createElement("div");
		row.classList.add("row");
		let c1, c2, c3;
		c1 = document.createElement("div");
		c2 = document.createElement("div");
		c3 = document.createElement("div");
		c1.classList.add("col-1");
		c2.classList.add("col-2");
		c3.classList.add("col-9");
		c1.appendChild(document.createTextNode(t.id));
		c2.appendChild(document.createTextNode(t.name));
		{
			let grad = document.createElement("div");
			grad.classList.add("grad");
			grad.style.background = _theme_to_grad(t.theme);
			grad.style.height = "4em";
			grad.classList.add("border");
			grad.classList.add("border-secondary");
			grad.classList.add("rounded");
			let img = document.createElement("img");
			img.src = t.theme.ilink;
			img.classList.add("mw-100");
			img.classList.add("mh-100");
			img.classList.add("mx-auto");
			img.classList.add("d-block");
			grad.appendChild(img);
			c3.appendChild(grad);
		}
		row.appendChild(c1);
		row.appendChild(c2);
		row.appendChild(c3);
		li.appendChild(row);
		return li;
	}
	return this;
}

function load_users(callback) {
	send_request({type:"get"}, callback);
}

function load_users_object(callback) {
	load_users((u)=>{
		if (u.error) {
			callback(null);
			return;
		}
		let users = new Users();
		for (let i = 0; i < u.length; i++) {
			let user = u[i];
			users.add_user(user.id, user.name, user.theme);
		}
		callback(users);
	});
}

function load_default(callback) {
	send_request({type:"default"}, callback);
}

function load_default_object(callback) {
	load_default((c) => {
		let users = new Users();
		users.add_user(c.id, c.name, c.theme);
		callback(users);
	});
}

function print_users() {
	load_users(console.log);
}

function set_users(users, callback) {
	send_request({type:"set", users:users}, callback);
}

function set_users_object(users, callback) {
	set_users(users.users, callback);
}

function build_user_list(elem) {
	if (!elem.hasChildNodes()) {
		elem.classList.add("list-group");
		load_default_object((df)=>_append_user_list(elem, df));
		load_users_object((df)=>_append_user_list(elem, df, true));
	} else {
		load_users_object((df)=>_append_user_list(elem, df, true));
	}
}

var last_users = [];

function _append_user_list(elem, users, real=false) {
	for (let i = 0; i < users.users.length; i++) {
		let new_u = users.toHTMLs(i);
		if (elem.childNodes.length > i+1) {
			if (last_users[i] != new_u.innerHTML) {
				console.log("Replace", i+1);
				elem.replaceChild(new_u, elem.childNodes[i+1]);
				last_users[i] = new_u.innerHTML;
			}
		} else {
			elem.appendChild(new_u);
			last_users[i]=new_u.innerHTML;
		}
	}
}
