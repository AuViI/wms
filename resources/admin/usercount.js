"use strict"


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

    // TODO logic here
    let usercount = 0;
    
    for (let i = 0; i < containers.length; i++) {
        let c = containers[i];
        let s = document.createElement("span");
        let t = document.createTextNode(usercount);
        s.appendChild(t);
        s.classList.add("badge");
        s.classList.add("badge-secondary");
        c.appendChild(s);
    }
}
