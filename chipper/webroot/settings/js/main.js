function showWeatherForm() {
    var keyLabel = document.createElement("label")
    keyLabel.innerHTML = "Enter your API key here: "
    keyLabel.htmlFor = "weatherKey";
    var keyForm = document.createElement("input")
    keyForm.type = "text"
    keyForm.id = "weatherKey"
    keyForm.name = "weatherKey"
    weatherSelect = document.getElementById("weatherSettings")
    weatherKeyDiv = document.getElementById("weatherAPIInput")
    if (weatherSelect.value == "weather_yes") {
        weatherKeyDiv.innerHTML = ""
        weatherKeyDiv.appendChild(keyLabel)
        weatherKeyDiv.appendChild(keyForm)
    } else {
        weatherKeyDiv.innerHTML = ""
    }
}

function showHoundForm() {
    var idLabel = document.createElement("label")
    idLabel.innerHTML = "Enter your Client ID here: "
    idLabel.htmlFor = "knowledgeID";
    var idForm = document.createElement("input")
    idForm.type = "text"
    idForm.id = "knowledgeID"
    idForm.name = "knowledgeID"

    var keyLabel = document.createElement("label")
    keyLabel.innerHTML = "Enter your Client Key here: "
    keyLabel.htmlFor = "knowledgeKey";
    var keyForm = document.createElement("input")
    keyForm.type = "text"
    keyForm.id = "knowledgeKey"
    keyForm.name = "knowledgeKey"
    knowledgeSelect = document.getElementById("knowledgeSettings")
    knowledgeKeyDiv = document.getElementById("houndifyInput")
    if (knowledgeSelect.value == "houndify_yes") {
        knowledgeKeyDiv.innerHTML = ""
        knowledgeKeyDiv.appendChild(idLabel)
        knowledgeKeyDiv.appendChild(idForm)
        knowledgeKeyDiv.appendChild(keyLabel)
        knowledgeKeyDiv.appendChild(keyForm)
    } else {
        knowledgeKeyDiv.innerHTML = ""
    }
}

function setupPod() {
    var houndSend = "false"
    var weatherSend = "false"
    var certSend = "ip"
    var portSend = "443"
    var hostnameValue = document.getElementById("hostnameSettings").value
    var weatherValue = document.getElementById("weatherSettings").value
    var weatherKey = ""
    var houndValue = document.getElementById("knowledgeSettings").value
    var houndID = ""
    var houndKey = ""
    if (weatherValue == "weather_yes") {
        weatherKey = document.getElementById("weatherKey").value
        weatherSend = "true"
        if (weatherKey == "") {
            alert("You must enter a weatherapi key, or disable weather.")
            return
        }
    }
    if (houndValue == "houndify_yes") {
        houndID = document.getElementById("knowledgeID").value
        houndKey = document.getElementById("knowledgeKey").value
        houndSend = "true"
        if (houndID == "" || houndKey == "") {
            alert("You must enter a Houndify client ID and key, or disable knowledge graph.")
            return
        }
    }
    if (hostnameValue == "escapepod_local"){
        certSend = "epod"
    }
    var sendString = "?port=" + portSend + "&certType=" + certSend + "&weatherEnable=" + weatherSend + "&weatherKey=" + weatherKey + "&houndifyEnable=" + houndSend + "&houndifyID=" + houndID + "&houndifyKey=" + houndKey
    console.log(sendString)
    let xhr = new XMLHttpRequest();
    xhr.open("GET", "/chipper/make_config" + sendString);
    xhr.setRequestHeader("Cache-Control", "no-cache, no-store, max-age=0");
    xhr.responseType = 'json';
    xhr.send();
    xhr.onload = function() {
        console.log(xhr.response)
        alert("Success! You are ready to start wire-pod.")
    }   
    return
}

function startChipper() {
    let xhr = new XMLHttpRequest();
    xhr.open("GET", "/chipper/start_chipper");
    xhr.setRequestHeader("Cache-Control", "no-cache, no-store, max-age=0");
    xhr.send();
    xhr.onload = function() {
        console.log(xhr.response)
        alert(xhr.response)
    }   
}

function stopChipper() {
    let xhr = new XMLHttpRequest();
    xhr.open("GET", "/chipper/stop_chipper");
    xhr.setRequestHeader("Cache-Control", "no-cache, no-store, max-age=0");
    xhr.send();
    xhr.onload = function() {
        console.log(xhr.response)
        alert(xhr.response)
    }   
}

function restartChipper() {
    let xhr = new XMLHttpRequest();
    xhr.open("GET", "/chipper/restart_chipper");
    xhr.setRequestHeader("Cache-Control", "no-cache, no-store, max-age=0");
    xhr.send();
    xhr.onload = function() {
        console.log(xhr.response)
        alert(xhr.response)
    }   
}

function setupBot() {
    alert("This does not work at the moment!")
    return
}