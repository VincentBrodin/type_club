const input = document.getElementById("input");
const stats = document.getElementById("data");
const output = document.getElementById("output");
const overlay = document.getElementById("overlay");
const caret = document.getElementById("caret");

overlay.style.display = "flex";

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

input.addEventListener("input", OnInput);
input.addEventListener("blur", OnBlur);
input.addEventListener("keydown", OnKeyDown);

document.addEventListener("DOMContentLoaded", OnLoaded);

async function OnLoaded() {
    const queryString = window.location.search;
    const urlParams = new URLSearchParams(queryString);
    const id = urlParams.get("id");
    if (id == null) {
        GetTarget();
    } else {
        RemoveFade();
        const response = await fetch(`/replay?id=${id}`);
        if (!response.ok) {
            throw new Error(`Response status: ${response.status}`);
        }
        const data = await response.json();
        SetFade();
        overlay.style.display = "none";
        SetTarget(data.target);
    }
}
async function GetTarget(wordCount = 10) {
    try {
        RemoveFade();
        const response = await fetch(`/random?length=${wordCount}`);
        if (!response.ok) {
            throw new Error(`Response status: ${response.status}`);
        }
        SetFade();
        overlay.style.display = "none";
        const text = await response.text();
        SetTarget(text);
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
}

function SetCaret() {
    const cursor = document.getElementById("cursor");
    if (cursor == null) {
        caret.innerText = "";
        return;
    } else {
        caret.innerText = "|";
    }
    const cursorRect = cursor.getBoundingClientRect();
    caret.style.left = cursorRect.left - (cursorRect.width / 2) + "px";
    caret.style.top = cursorRect.top + (cursorRect.height / 8) + "px";
}

function OnInput() {
    SetOutput(input.value, target);
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
        overlay.style.display = "flex";
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

function SetOutput(text, target) {
    const inputWords = input.value.split(" ").length - 1;
    const targetWords = target.split(" ").length;
    stats.innerText = `${inputWords}/${targetWords}`;

    const textLen = text.length;
    const targetLen = target.length;

    let out = "<h3>";
    for (let i = 0; i < textLen; i++) {
        if (text[i] == target[i]) {
            out += target[i];
        } else {
            out += `<span class="text-danger text-decoration-underline">${
                target[i]
            }</span>`;
        }
    }

    let first = true;
    for (let i = textLen; i < targetLen; i++) {
        if (first) {
            first = false;
            out += `<span id="cursor" class="text-secondary">${
                target[i]
            }</span>`;
        } else {
            out += `<span class="text-secondary">${target[i]}</span>`;
        }
    }

    out += "</h3>";
    output.innerHTML = out;
    SetCaret();
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
