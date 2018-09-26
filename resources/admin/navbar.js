"use strict"

function navbar_draw(page) {
    var ajax = new XMLHttpRequest();
    ajax.onreadystatechange = function() {
        if (ajax.readyState != 4) {
            return;
        }
        if (ajax.status != 200) {
            return;
        }
        navbar_update(ajax.responseText);
        update_usercount_badges();
    }
    ajax.open("POST", "/resources/admin/nav");
    ajax.setRequestHeader("Content-Type", "application/json");
    ajax.send(JSON.stringify({origin: page}));
}

function navbar_update(innerHTML) {
    let wrapper = document.getElementById("navbar-wrapper");
    wrapper.innerHTML = innerHTML;
}
