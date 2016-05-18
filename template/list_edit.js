print = window.console.log
var cdiv = document.getElementById("container")
var lon = 0
function seti(num) {
    return document.getElementById("s"+num)
}
function dels(num) {
    cdiv.removeChild(seti(num))
    update()
}
function news() {
    lon += 1
    newDiv = document.createElement("div")
    newDiv.className = "select"
    newDiv.id = "s"+lon
    newDiv.innerHTML = '<select name="service'+lon+'">\
            <option value="txt">Text</option>\
            <option value="view">Anzeige</option>\
            <option value="csv">CSV</option>\
        </select>\
        <input type="text" placeholder="Ort (Parameter)" name="'+lon+"_para"+1+'">\
        <input type="button" onclick="news()" name="okBtn'+lon+'" value="OK">\
        <input type="button" onclick="dels('+lon+')" name="delBtn'+lon+'" value="Entfernen">'
    cdiv.appendChild(newDiv)
    update()
}
function lastp() {
    res = 0
    while (true) {
        n = cdiv.innerHTML.indexOf("<p>",res+1)
        if (n==-1) {
            return res
        } else {
            res = n
        }
    }
}
function save() {
    console.log("saving")
}
function update() {
    var dels = document.getElementsByClassName("_del_on_update")
    for (var ob in dels) {
        if (dels[ob].parentNode) {
            dels[ob].parentNode.removeChild(dels[ob])
        }
    }
    for( var i = 0; i < cdiv.childElementCount; i++) {
        if ( i < cdiv.children.length - 1 ) {
            for ( var c = 0; c < cdiv.children[i].childElementCount - 1; c++) {
                cdiv.children[i].children[c].disabled = true;
            }
        } else {
            for ( var c in cdiv.children[i].children ) {
                cdiv.children[i].children[c].disabled = false;
            }
        }
    }
}
function output() {
    update()
    res = []
    dbg_string = "Output:\n"
    for (var i = 0; i < cdiv.childElementCount; i++) {
        res[i] = {}
        res[i]["data"] = {}
        res[i]["type"] = cdiv.children[i].children[0].selectedOptions[0].value
        dbg_string += i + " => " + res[i]["type"] + " {\n"
        for (var c = 1; c < cdiv.children[i].childElementCount - 2; c++) {
            r = cdiv.children[i].children[c]
            var c = r.value
            if (c == "") {
                c = "nil"
            } else {
                c = "\"" + c + "\""
            }
            res[i]["data"][r.name] = c
            dbg_string += "\t"+ r.name + " => " + c + "\n"
        }
        var l = 0
        var v = ""
        for (var c in res[i]["data"]) {
            v = res[i].data[c]
            l++
        }
        if (l == 1 && v != "nil"){
            var link = "/"+res[i]["type"]+"/"+v.substring(1,v.length-1)+"/"
            dbg_string += "\t__link__ => " + link
            var anchor = document.createElement("a")
            anchor.href=link
            anchor.target ="_blank"
            anchor.className="_del_on_update"
            anchor.innerHTML="<button>Ansehen: " + link + "</button>"
            cdiv.children[i].appendChild(anchor)
        }
        res[i].data.length=l
        dbg_string +="}\n"
    }
    console.log(dbg_string)
    return res
}

news()
