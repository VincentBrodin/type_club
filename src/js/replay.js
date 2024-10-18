const output = document.getElementById("output");
const caret = document.getElementById("caret");

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
        StartReplay();
    } catch (error) {
        console.error(error.message);
    }
}

function StartReplay() {
    if (interval) {
        clearInterval(interval);
    }
    const startTime = Date.now();
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
        clearInterval(interval);
    }
}

GetReplay();
