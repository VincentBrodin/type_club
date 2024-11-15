const output = document.getElementById("output");
const caret = document.getElementById("ghostCaret");

let interval = null;
let replay = null;

async function GetReplay() {
    try {
        OverlayOn();
        const queryString = window.location.search;
        const urlParams = new URLSearchParams(queryString);
        const id = urlParams.get("id");

        const response = await fetch("/stats", {
            method: "POST",
            body: JSON.stringify({
                id: id,
            }),
            headers: {
                "Content-Type": "application/json; charset=UTF-8",
            },
        });
        if (response.ok) {
            replay = await response.json();
            StartReplay();
            OverlayOff();
        }
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

document.addEventListener("DOMContentLoaded", GetReplay);
