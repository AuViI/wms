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
		if (ajax.responseText.length > 2) {
			callback(JSON.parse(ajax.responseText));
		}
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
	this.getEditFunc = (t) => {
		return ()=>{
			let c = document.getElementById("user_edit_container");
			let empty_c = () => {
				c.classList.remove("d-none");
				for (let i = c.childNodes.length; i>0; i--) {
					c.removeChild(c.childNodes[i-1]);
				}
			}
			empty_c();
			let header = document.createElement("h4");
			header.appendChild(document.createTextNode("Kontobearbeitung"));
			c.appendChild(header);

			let form = document.createElement("form");
			form.classList.add("mx-3", "p-3", "border", "rounded");

			let name_div = document.createElement("div");
			let colors_div = document.createElement("div");
			let icon_div = document.createElement("div");

			name_div.classList.add("form-group");
			colors_div.classList.add("form-group");
			icon_div.classList.add("form-group");

			name_div.appendChild(document.createElement("label"));
			name_div.appendChild(document.createElement("input"));
			name_div.childNodes[1].id="form-name-input";
			name_div.childNodes[0].for = name_div.childNodes[1].id;
			name_div.childNodes[0].appendChild(document.createTextNode("Anzeigename"));
			name_div.childNodes[1].placeholder = "Anzeigename";
			name_div.childNodes[1].value = t.name;
			name_div.childNodes[1].classList.add("form-control");

			let colors_label = document.createElement("label");
			let colors_row = document.createElement("div");
			colors_div.appendChild(colors_label);
			colors_div.appendChild(colors_row);
			colors_row.id="form-colors-input";
			colors_label.for = colors_row.id;
			colors_label.appendChild(document.createTextNode("Anzeigefarben"));


			let cc = (c) => {
				let hex = Number(c).toString(16);
				if (hex.length < 2) {
					hex = "0" + hex;
				}
				return hex;
			};
			let ccc = (c) => {
				return "#" + cc(c.r) + cc(c.g) + cc(c.b);
			}
			colors_row.classList.add("row")
			let colors_cols = [document.createElement("div"), document.createElement("div")];
			colors_row.appendChild(colors_cols[0]);
			colors_row.appendChild(colors_cols[1]);
			colors_cols[0].classList.add("col");
			colors_cols[1].classList.add("col");
			let cstart = document.createElement("input");
			let cend = document.createElement("input");
			cstart.type = "color";
			cend.type = "color";
			cstart.classList.add("form-control");
			cend.classList.add("form-control");
			cstart.value = ccc(t.theme.start);
			cend.value = ccc(t.theme.end);
			colors_cols[0].appendChild(cstart);
			colors_cols[1].appendChild(cend);
			let csmall = [document.createElement("small"), document.createElement("small")];
			csmall[0].classList.add("form-text", "text-muted");
			csmall[1].classList.add("form-text", "text-muted");
			colors_cols[0].appendChild(csmall[0]);
			colors_cols[1].appendChild(csmall[1]);
			csmall[0].appendChild(document.createTextNode("Startfarbe"));
			csmall[1].appendChild(document.createTextNode("Endfarbe"));


			let icon_label = document.createElement("label");
			let icon_row = document.createElement("div");
			icon_div.appendChild(icon_label);
			icon_div.appendChild(icon_row);
			icon_row.id="form-icon-input";
			icon_label.for = icon_row.id;
			icon_label.appendChild(document.createTextNode("Icon"));

			icon_row.classList.add("row")
			let icon_cols = [document.createElement("div"), document.createElement("div")];
			icon_row.appendChild(icon_cols[0]);
			icon_row.appendChild(icon_cols[1]);
			icon_cols[0].classList.add("col-9");
			icon_cols[1].classList.add("col-3");
			let icon_input = document.createElement("input");
			icon_input.value = t.theme.ilink;
			icon_input.type = "text";
			icon_input.classList.add("form-control");
			icon_input.placeholder = "Icon Link";
			icon_cols[0].appendChild(icon_input);
			let img = document.createElement("img");
			img.src = t.theme.ilink;
			img.style.maxHeight = "8rem";
			img.style.maxWidth = "8rem";
			icon_input.onkeyup = () => {img.src=icon_input.value;};
			icon_cols[1].appendChild(img);
			

			form.appendChild(name_div);
			form.appendChild(colors_div);
			form.appendChild(icon_div);

			let fbtn = document.createElement("div");
			fbtn.appendChild(document.createTextNode("Absenden"));
			fbtn.classList.add("btn", "btn-primary");
			fbtn.onclick = () => {
				let hexToRgb = (hex)=>{
					var result = /^#?([a-f\d]{2})([a-f\d]{2})([a-f\d]{2})$/i.exec(hex);
					return result ? {
						r: parseInt(result[1], 16),
						g: parseInt(result[2], 16),
						b: parseInt(result[3], 16)
					} : null;
				}

				let index = -1;
				load_users_object((u) => {
					for (let i = 0; i < u.users.length; i++) {
						if (u.users[i].id == t.id) {
							index = i;
						}
					}
					if (index < 0) {
						console.log("cannot find the user online");
						return;
					}
					u.users[index].name = name_div.childNodes[1].value;
					u.users[index].theme.ilink = icon_input.value;
					u.users[index].theme.start = hexToRgb(cstart.value);
					u.users[index].theme.end = hexToRgb(cend.value);
					set_users_object(u, console.log);
				});
				empty_c();
			};

			let sbtn = document.createElement("div");
			sbtn.appendChild(document.createTextNode("Abbrechen"));
			sbtn.classList.add("btn", "btn-warning", "mx-2");
			sbtn.onclick = empty_c;

			form.appendChild(fbtn);
			form.appendChild(sbtn);

			c.appendChild(form);
		};
	};
	this.toHTMLs = (i) => {
		let t = this.users[i];
		let li = document.createElement("div");
		li.classList.add("list-group-item");
		let row = document.createElement("div");
		row.classList.add("row");
		let c1, c2, c3, c4;
		c1 = document.createElement("div");
		c2 = document.createElement("div");
		c3 = document.createElement("div");
		c4 = document.createElement("div");
		c1.classList.add("col-1");
		c2.classList.add("col-2");
		c3.classList.add("col-7");
		c4.classList.add("col-2");
		let c1_dtage = document.createElement("a");
		c1.appendChild(c1_dtage);
		c1_dtage.appendChild(document.createTextNode(t.id));
		c1_dtage.href = "/dtage/Kühlungsborn&u=" + t.id;
		c1_dtage.target = "_blank";
		c2.appendChild(document.createTextNode(t.name));
		{ // c3 gradient
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
		{ // c4 edit button
			let btn = document.createElement("button");
			let btn_icon = document.createElement("span");
			btn_icon.classList.add("fas", "fa-edit");
			btn.appendChild(btn_icon);
			btn.classList.add("w-100", "h-100", "btn", "btn-warning");
			btn.onclick = this.getEditFunc(t);
			if (t.id==0) {
				btn.disabled = true;
			}
			c4.appendChild(btn);
		}
		row.appendChild(c1);
		row.appendChild(c2);
		row.appendChild(c3);
		row.appendChild(c4);
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
		load_default_object(
			(df)=>{
				_append_user_list(elem, df);
				load_users_object((df)=>_append_user_list(elem, df, true));
			});
	} else {
		load_users_object((df)=>_append_user_list(elem, df, true));
	}
}

function build_add_user(elem) {
	let row = document.createElement("div");
	row.classList.add("input-group", "px-3");
	let name = document.createElement("input");
	name.type = "text";
	name.placeholder = "Name";
	name.classList.add("form-control");
	let addbtn = document.createElement("button");
	addbtn.appendChild(document.createTextNode("Neuen Nutzer hinzufügen"));
	addbtn.classList.add("btn", "btn-outline-primary", "form-control");
	addbtn.onclick = () => {
		addbtn.disabled = true;
		window.setTimeout(()=>{addbtn.disabled = false;}, 9000);
		if (name.value.length == 0) {
			return;
		}
		load_default_object((d)=>{
			load_users_object((u)=>{
				let mi = 0;
				for (let i = 0; i < u.users.length; i++) {
					mi = Math.max(u.users[i].id + 1, mi);
				}
				u.add_user(mi, name.value, d.users[0].theme);
				set_users_object(u, console.log);
			});
		});
	};
	row.appendChild(name);
	row.appendChild(addbtn);
	elem.appendChild(row);
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
