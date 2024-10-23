const input = document.getElementById("input");
const stats = document.getElementById("data");
const output = document.getElementById("output");
const caret = document.getElementById("caret");
const ghostCaret = document.getElementById("ghostCaret");
ghostCaret.innerText = "";
let ghost = false;

const repeat = document.getElementById("repeat");
repeat.addEventListener("click", Repeat);

const words10 = document.getElementById("words10");
words10.addEventListener("click", () => {
    RemoveSelected();
    GetTarget(10);
    words10.classList.add("selected");
});
const words20 = document.getElementById("words20");
words20.addEventListener("click", () => {
    RemoveSelected();
    GetTarget(20);
    words20.classList.add("selected");
});
const words30 = document.getElementById("words30");
words30.addEventListener("click", () => {
    RemoveSelected();
    GetTarget(30);
    words30.classList.add("selected");
});
const words40 = document.getElementById("words40");
words40.addEventListener("click", () => {
    RemoveSelected();
    GetTarget(40);
    words40.classList.add("selected");
});
const words50 = document.getElementById("words50");
words50.addEventListener("click", () => {
    RemoveSelected();
    GetTarget(50);
    words50.classList.add("selected");
});

let target = "";

let started = false;
let startTime = null;

let totalInputs = 0;
let totalErrors = 0;

let inputs = [];

let replay = null;
let interval = null;

input.addEventListener("input", OnInput);
input.addEventListener("blur", OnBlur);
input.addEventListener("keydown", OnKeyDown);

document.addEventListener("DOMContentLoaded", OnLoaded);
window.addEventListener("resize", () => {
    SetCaret(caret, "cursor");
});

async function OnLoaded() {
    const queryString = window.location.search;
    const urlParams = new URLSearchParams(queryString);
    const id = urlParams.get("id");
    if (id == null) {
        GetTarget();
    } else {
        ghost = true;
        RemoveFade();
        console.log(id);
        await GetReplay(id);
        SetFade();
        SetTarget(replay.target);
    }
}
async function GetReplay(id) {
    try {
        OverlayOn();
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
    //Calculates the cursor position
    SetOutput(replay.inputs[lastState].value, replay.target);
    SetCaret(ghostCaret, "cursor");
    SetOutput(input.value, target);
    if (lastState == replay.inputs.length - 1) {
        clearInterval(interval);
    }
}

async function GetTarget(wordCount = 10) {
    try {
        OverlayOn();
        RemoveFade();
        const response = await fetch("/sentence", {
            method: "POST",
            body: JSON.stringify({
                length: wordCount,
            }),
            headers: {
                "Content-Type": "application/json; charset=UTF-8",
            },
        });
        if (response.ok) {
            SetFade();
            const text = await response.text();
            SetTarget(text);
            OverlayOff();
        }
    } catch (error) {
        console.error(error.message);
    }
}

function SetTarget(value) {
    target = value;
    input.value = "";
    inputs = [];
    started = false;
    startTime = null;
    SetOutput("", target);
    SetCaret(caret, "cursor");
}

function OnInput() {
    SetOutput(input.value, target);
    SetCaret(caret, "cursor");
    const inputWords = input.value.split(" ").length - 1;
    const targetWords = target.split(" ").length;

    stats.innerText = `${inputWords}/${targetWords}`;

    const inputLen = input.value.length;
    const targetLen = target.length;
    //Accuracy checking
    totalInputs++;
    if (input.value[inputLen - 1] !== target[inputLen - 1]) {
        totalErrors++;
    }

    //Test start
    if (inputLen !== 0 && !started) {
        started = true;
        startTime = Date.now();
        if (ghost) {
            StartReplay();
        }
    }

    //Add input to inputs
    const timeSinceStart = Date.now() - startTime;
    inputs.push({
        value: input.value,
        time: timeSinceStart,
    });

    //Test end
    if (inputLen === targetLen) {
        //Show overlay and hide input
        input.type = "hidden";
        OverlayOn();
        //WPM
        const milisecs = Date.now() - startTime;
        const secs = milisecs / 1000;
        const mins = secs / 60;
        const words = target.split(" ").length;
        //The avrage word length is 4.7 so this is more avrage
        const avrageWords = target.length / 4.7;
        const wpm = words / mins;
        const awpm = avrageWords / mins;

        //Accuracy
        const accuracy = 1 - (totalErrors / totalInputs);

        data = {
            target: target,
            html: String(output.innerHTML),
            accuracy: accuracy,
            wpm: wpm,
            awpm: awpm,
            time: secs,
            inputs: inputs,
        };
        Done(data);
    }
}

async function Done(data) {
    try {
        const response = await fetch("/done", {
            method: "POST",
            body: JSON.stringify(data),
            headers: {
                "Content-Type": "application/json; charset=UTF-8",
            },
        });

        if (response.redirected) {
            window.location.href = response.url;
        } else if (!response.ok) {
            const errorText = await response.text();
            console.error(errorText);
        }
    } catch (error) {
        console.error(error.message);
    }
}

function OnKeyDown(event) {
    if (
        event.key === "ArrowLeft" || event.key === "ArrowRight" ||
        event.key === "ArrowUp" || event.key == "ArrowDonw"
    ) {
        event.preventDefault();
    }
}

function OnBlur() {
    input.focus();
}

function RemoveFade() {
    const elementsToAnimate = document.querySelectorAll(".fade-in-target");
    elementsToAnimate.forEach((el) => {
        el.classList.remove("fade-in");
    });
}

function SetFade() {
    const elementsToAnimate = document.querySelectorAll(".fade-in-target");
    elementsToAnimate.forEach((el) => {
        el.classList.add("fade-in");
    });
}

function RemoveSelected() {
    words10.classList.remove("selected");
    words20.classList.remove("selected");
    words30.classList.remove("selected");
    words40.classList.remove("selected");
    words50.classList.remove("selected");
}

function Repeat(event) {
    event.preventDefault();
    const queryString = window.location.search;
    const urlParams = new URLSearchParams(queryString);
    let id = urlParams.get("id");
    if (id == "") {
        id = "last";
    }
    console.log(id);
    window.location.href = `/?id=${id}`;
}
