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

function showKnowledgeForm() {
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

    var openAiKeyLabel = document.createElement("label")
    openAiKeyLabel.innerHTML = "Enter your OpenAI API key here: "
    openAiKeyLabel.htmlFor = "openaiKey";
    var openAiKeykeyForm = document.createElement("input")
    openAiKeykeyForm.type = "text"
    openAiKeykeyForm.id = "openaiKey"
    openAiKeykeyForm.name = "openaiKey"

    var graphLabel = document.createElement("label")
    graphLabel.innerHTML = "Would you like to enable the intent graph feature?: "
    graphLabel.htmlFor = "intentGraphSelect";
    var graphLabelYes = document.createElement("label")
    graphLabelYes.innerHTML = "Yes"
    graphLabelYes.htmlFor = "intentGraphYes";
    var graphLabelNo = document.createElement("label")
    graphLabelNo.innerHTML = "No"
    graphLabelNo.htmlFor = "intentGraphNo";
    var intentGraph = document.createElement("form")
    intentGraph.id = "intentGraphSelect"
    intentGraph.name = "intentGraphSelect"

    var intentYes = document.createElement("INPUT");
    intentYes.setAttribute("type", "radio");
    intentYes.value = "true"
    intentYes.name = "intentGraphSelect"
    intentYes.id = "intentGraphYes"
    intentYes.form = "intentGraphSelect"
    intentYes.checked = true

    var intentNo = document.createElement("INPUT");
    intentNo.setAttribute("type", "radio");
    intentNo.form = "intentGraphSelect"
    intentNo.name = "intentGraphSelect"
    intentNo.id = "intentGraphNo"
    intentNo.value = "false"
    intentNo.checked = false
    intentGraph.appendChild(graphLabelYes)
    intentGraph.appendChild(intentYes)
    intentGraph.appendChild(document.createElement("br"))
    intentGraph.appendChild(graphLabelNo)
    intentGraph.appendChild(intentNo)

    knowledgeSelect = document.getElementById("knowledgeSettings")
    knowledgeKeyDiv = document.getElementById("knowledgeInput")
    if (knowledgeSelect.value == "knowledge_houndify") {
        knowledgeKeyDiv.innerHTML = ""
        knowledgeKeyDiv.appendChild(idLabel)
        knowledgeKeyDiv.appendChild(idForm)
        knowledgeKeyDiv.appendChild(keyLabel)
        knowledgeKeyDiv.appendChild(keyForm)
    } else if (knowledgeSelect.value == "knowledge_openai") {
        knowledgeKeyDiv.innerHTML = ""
        knowledgeKeyDiv.appendChild(openAiKeyLabel)
        knowledgeKeyDiv.appendChild(openAiKeykeyForm)
        knowledgeKeyDiv.appendChild(graphLabel)
        knowledgeKeyDiv.appendChild(intentGraph)
    } else {
        knowledgeKeyDiv.innerHTML = ""
    }
}

function setupPod() {
    var knowledgeSend = "false"
    var weatherSend = "false"
    var certSend = "ip"
    var portSend = "443"
    var graphSetup = "false"
    var hostnameValue = document.getElementById("hostnameSettings").value
    var weatherValue = document.getElementById("weatherSettings").value
    var weatherKey = ""
    var houndValue = document.getElementById("knowledgeSettings").value
    var sttValue = document.getElementById("sttSettings").value
    var sttSend = ""
    var picovoiceKey = ""
    var knowledgeID = ""
    var knowledgeKey = ""
    var knowledgeProvider = ""
    if (weatherValue == "weather_yes") {
        weatherKey = document.getElementById("weatherKey").value
        weatherSend = "true"
        if (weatherKey == "") {
            alert("You must enter a weatherapi key, or disable weather.")
            return
        }
    }
    if (houndValue == "knowledge_houndify") {
        knowledgeID = document.getElementById("knowledgeID").value
        knowledgeKey = document.getElementById("knowledgeKey").value
        knowledgeSend = "true"
        knowledgeProvider = "houndify"
        if (knowledgeID == "" || knowledgeKey == "") {
            alert("You must enter a Houndify client ID and key, or disable knowledge graph.")
            return
        }
    } else if (houndValue == "knowledge_openai") {
        knowledgeKey = document.getElementById("openaiKey").value
        knowledgeSend = "true"
        knowledgeProvider = "openai"
        if (knowledgeKey == "") {
            alert("You must enter an OpenAI API key, or disable knowledge graph.")
            return
        }
        graphSetup = document.getElementById("intentGraphYes").value
    }
    if (sttValue == "picovoice_leopard") {
        picovoiceKey = document.getElementById("picovoiceKey").value
        if (picovoiceKey == "") {
            alert("You must enter a Picovoice API key, or use Coqui instead.")
            return
        }
        sttSend = "leopard"
    } else if (sttValue == "vosk_stt") {
        sttSend = "vosk"
    } else {
        sttSend = "coqui"
    }
    if (hostnameValue == "escapepod_local"){
        certSend = "epod"
    }
    var sendString = "?port=" + portSend + "&certType=" + certSend + "&weatherEnable=" + weatherSend + "&weatherKey=" + weatherKey + "&knowledgeEnable=" + knowledgeSend + "&knowledgeID=" + knowledgeID + "&knowledgeKey=" + knowledgeKey + "&knowledgeProvider=" + knowledgeProvider + "&knowledgeIntent=" + graphSetup + "&sttService=" + sttSend + "&picovoiceKey=" + picovoiceKey
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
    alert("Your bot is going to be set up. This will put Vector back into Onboarding mode, but will not clear user data. This may take up to a minute. Press OK to begin the process.")
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
            alert("Success! Use https://keriganc.com/vector-epod-setup on any device to finish setting up your bot.")
            return
        } else {
            alert("Error: " + responseString)
            return
        }
    }
}