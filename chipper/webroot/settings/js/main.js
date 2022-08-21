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
    alert("This does not work at the moment!")
    return
}

function setupBot() {
    alert("This does not work at the moment!")
    return
}