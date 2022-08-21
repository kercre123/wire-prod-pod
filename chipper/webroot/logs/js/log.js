logDiv = document.getElementById("wirePodLogs")
logP = document.createElement("p")


setInterval(function() {
    let xhr = new XMLHttpRequest();
    xhr.open("GET", "/chipper/get_logs");
    //xhr.setRequestHeader("Cache-Control", "no-cache, no-store, max-age=0");
    xhr.send();
    xhr.onload = function() {
        logDiv.innerHTML = ""
        logP.innerHTML = xhr.response
        logDiv.appendChild(logP)
    }
}, 200)