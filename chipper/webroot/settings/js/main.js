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

function showPicovoiceForm() {
    var keyLabel = document.createElement("label")
    keyLabel.innerHTML = "Enter your API key here: "
    keyLabel.htmlFor = "picovoiceKey";
    var keyForm = document.createElement("input")
    keyForm.type = "text"
    keyForm.id = "picovoiceKey"
    keyForm.name = "picovoiceKey"
    sttSelect = document.getElementById("sttSettings")
    picovoiceKeyDiv = document.getElementById("picovoiceAPIInput")
    if (sttSelect.value == "picovoice_leopard") {
        picovoiceKeyDiv.innerHTML = ""
        picovoiceKeyDiv.appendChild(keyLabel)
        picovoiceKeyDiv.appendChild(keyForm)
    } else {
        picovoiceKeyDiv.innerHTML = ""
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
    var sttValue = document.getElementById("sttSettings").value
    var sttSend = ""
    var picovoiceKey = ""
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
    if (sttValue == "picovoice_leopard") {
        picovoiceKey = document.getElementById("picovoiceKey").value
        if (picovoiceKey == "") {
            alert("You must enter a Picovoice API key, or use Coqui instead.")
            return
        }
        sttSend = "leopard"
    } else {
        sttSend = "coqui"
    }
    if (hostnameValue == "escapepod_local"){
        certSend = "epod"
    }
    var sendString = "?port=" + portSend + "&certType=" + certSend + "&weatherEnable=" + weatherSend + "&weatherKey=" + weatherKey + "&houndifyEnable=" + houndSend + "&houndifyID=" + houndID + "&houndifyKey=" + houndKey + "&sttService=" + sttSend + "&picovoiceKey=" + picovoiceKey
    console.log(sendString)
    let xhr = new XMLHttpRequest();
    xhr.open("GET", "/chipper/make_config" + sendString);
    xhr.setRequestHeader("Cache-Control", "no-cache, no-store, max-age=0");
    xhr.responseType = 'json';
    xhr.send();
    xhr.onload = function() {
        console.log(xhr.response)
        restartChipperNoAlert()
        alert("Success! wire-pod has started successfully.")
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

function restartChipperNoAlert() {
    let xhr = new XMLHttpRequest();
    xhr.open("GET", "/chipper/restart_chipper");
    xhr.setRequestHeader("Cache-Control", "no-cache, no-store, max-age=0");
    xhr.send();
    xhr.onload = function() {
        console.log(xhr.response)
    }   
}

function setupBot() {
    if (document.getElementById("botIP").value == "") {
        alert("You must enter your bot's IP address. This can be found in CCIS (put Vector on the charger, double press his button, lift his lift up then down).")
        return
    }
    let sshKey = document.getElementById("sshKeyfile").files[0];  // file from input
    let req = new XMLHttpRequest();
    let formData = new FormData();
    formData.append("file", sshKey);                                
    req.open("POST", "/chipper/upload_ssh_key");
    req.send(formData);
    let xhr = new XMLHttpRequest();
    xhr.open("GET", "/chipper/setup_bot" + "?botIP=" + document.getElementById("botIP").value);
    xhr.setRequestHeader("Cache-Control", "no-cache, no-store, max-age=0");
    xhr.send();
    xhr.onload = function() {
        console.log(xhr.response)
        responseString = xhr.response
        if (responseString.includes("Unable to")) {
            alert("Unable to communicate with robot. The key may be invalid, the bot may not be unlocked, or this device and the robot are not on the same network.")
            return
        } else if (responseString.includes("Everything")) {
            alert("Success! Voice commands should now be working with your bot.")
            return
        } else {
            alert("Error: " + responseString)
            return
        }
    }
}