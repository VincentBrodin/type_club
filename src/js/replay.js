const output = document.getElementById("output");
let interval = null;
let data = null;

async function GetReplay() {
    try {
        const queryString = window.location.search;
        const urlParams = new URLSearchParams(queryString);
        const id = urlParams.get("id");
        const response = await fetch(`/replay?id=${id}`);
        if (!response.ok) {
            throw new Error(`Response status: ${response.status}`);
        }
        data = await response.json();
        StartReplay()
    } catch (error) {
        console.error(error.message);
    }
}

function StartReplay() {
    if (interval) {
        clearInterval(interval)
    }
    let startTime = Date.now();
    interval = setInterval(ReplayLoop, 1, data, startTime);
}

function ReplayLoop(data, startTime) {
    let lastState = null;
    const timeSinceStart = Date.now() - startTime;
    for (state in data.inputs) {
        if (data.inputs[state].time >= timeSinceStart) {
            break;
        }
        lastState = state;
    }
    SetOutput(data.inputs[lastState].value, data.target);
    if (lastState == data.inputs.length - 1) {
        clearInterval(interval)
    }
}

function SetOutput(text, target) {
    let textLen = text.length;
    let targetLen = target.length;

    let out = "<h3>";
    for (let i = 0; i < textLen; i++) {
        if (text[i] == target[i]) {
            out += target[i];
        } else {
            out += `<span class="text-danger text-decoration-underline">${target[i]}</span>`;
        }
    }

    let first = true;
    for (let i = textLen; i < targetLen; i++) {
        if (first) {
            first = false;
            out += `<span class="text-secondary text-decoration-underline">${target[i]}</span>`;
        } else {
            out += `<span class="text-secondary">${target[i]}</span>`;
        }
    }

    out += "</h3>";
    output.innerHTML = out;
}

GetReplay()

