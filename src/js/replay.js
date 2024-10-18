const output = document.getElementById("output");
const caret = document.getElementById("ghostCaret");

let interval = null;
let replay = null;

async function GetReplay() {
    try {
        const queryString = window.location.search;
        const urlParams = new URLSearchParams(queryString);
        const id = urlParams.get("id");
        const response = await fetch(`/replay?id=${id}`);
        if (!response.ok) {
            throw new Error(`Response status: ${response.status}`);
        }
        replay = await response.json();
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
    interval = setInterval(ReplayLoop, 1, startTime);
}

function ReplayLoop(startTime) {
    let lastState = null;
    const timeSinceStart = Date.now() - startTime;
    for (const state in replay.inputs) {
        if (replay.inputs[state].time >= timeSinceStart) {
            break;
        }
        lastState = state;
    }
    SetOutput(replay.inputs[lastState].value, replay.target);
    SetCaret(caret, "cursor");
    if (lastState == replay.inputs.length - 1) {
        clearInterval(interval);
    }
}

GetReplay();
