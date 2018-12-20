"use strict"

function get_wp_list(callback) {
    var ajax = new XMLHttpRequest();
    ajax.onreadystatechange = function() {
        if (ajax.readyState != 4) {
            return;
        }
        if (ajax.status != 200) {
            return;
        }
		callback(ajax.responseText);
    }
    ajax.open("GET", "/wp/");
    ajax.setRequestHeader("Content-Type", "application/json");
    ajax.send();
}


function build_wp_list(par) {
	window.last_wp = par;
	// 0 = same day, 1 d1 before d2, -1 d2 before d1
	var cmp = (d1, d2)=> {
		if (d1.getFullYear() < d2.getFullYear()) {
			return 3;
		} else if (d1.getFullYear() > d2.getFullYear()) {
			return -3;
		}
		if (d1.getMonth() < d2.getMonth()) {
			return 2;
		} else if (d1.getMonth() > d2.getMonth()) {
			return -2;
		}
		if (d1.getDate() < d2.getDate()) {
			return 1;
		} else if (d1.getDate() > d2.getDate()) {
			return -1;
		}
		return 0;
	}
	get_wp_list((tdata)=>{
		var data = JSON.parse(tdata);
		var dc = document.createElement("div");
		var now = new Date(Date.now());
		dc.classList.add("list-group");
		for (let i = 0; i < data.length; i++) {
			var dc1 = document.createElement("div");
			dc1.classList.add("list-group-item");

			var drow = document.createElement("div");
			drow.classList.add("row");

			var dt = new Date(data[i].y, data[i].m - 1, data[i].d, 12);
			var dtc = cmp(dt, now);
			if (dtc > 0) {
				dc1.classList.add("disabled");
			} else if (dtc == 0) {
				dc1.classList.add("active");
			} else if (dtc < 0) {
				dc1.classList.add("future");
			}

			var c1, c2, c3;
			if (dtc == 0) {
				c1 = document.createElement("a");
				c1.href = "/forecast/"+data[i].loc;
				c1.classList.add("text-light");
				c1.target = "_blank";
			} else {
				c1 = document.createElement("a");

				var year = dt.getFullYear();
				var month = dt.getMonth() + 1;
				var day = dt.getDate();

				if (month < 10) {
					month = "0" + month;
				}
				if (day < 10) {
					day = "0" + day;
				}

				c1.href = "/forecast/"+data[i].loc+"/d="year+month+day;
				c1.target = "_blank";
			}
			c2 = document.createElement("div");
			c3 = document.createElement("div");

			c1.classList.add("col-3");
			c2.classList.add("col-2");
			c3.classList.add("col-6");


			c1.appendChild(document.createTextNode(data[i].loc));
			c2.appendChild(document.createTextNode(data[i].d + "." + data[i].m + "." + data[i].y));
			c3.appendChild(document.createTextNode(data[i].c));

			drow.appendChild(c1);
			drow.appendChild(c2);
			drow.appendChild(c3);

			dc1.appendChild(drow);

			dc.appendChild(dc1);
		}
		var paren = document.getElementById(par);
		if (!paren.hasChildNodes()) {
			paren.appendChild(dc);
		} else {
			paren.replaceChild(dc, paren.lastChild);
		}
	});
}

function send_data() {
	var loc, date, content;

	loc = document.getElementById("wp_location").value;
	date = document.getElementById("wp_date").value;
	content = document.getElementById("wp_content").value;

	var url = "/wp/"+loc+"/"+date.replace(/-/g,"/");

    var ajax = new XMLHttpRequest();
    ajax.onreadystatechange = function() {
        if (ajax.readyState != 4) {
            return;
        }
        if (ajax.status != 200) {
            return;
        }
		if (window.last_wp) {
			build_wp_list(window.last_wp);
		}
    }
    ajax.open("POST", url);

	var data = new FormData();
	data.append("content", content);
    ajax.send(data);
}
