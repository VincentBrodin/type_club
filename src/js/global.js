function SetOutput(text, target) {
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
}

function SetCaret(c, t) {
    const cursor = document.getElementById(t);
    if (cursor == null) {
        c.innerText = "";
        return;
    } else {
        c.innerText = "|";
    }
    const cursorRect = cursor.getBoundingClientRect();
    c.style.left = cursorRect.left - (cursorRect.width / 2) + "px";
    c.style.top = cursorRect.top + (cursorRect.height / 8) + "px";
}

let overlay = document.getElementById("overlay");
let loading = document.getElementById("loading");
let loadingDots = 1;

function OverlayOn() {
    overlay.classList.add("d-flex");
    overlay.classList.remove("d-none");
}

function OverlayOff() {
    overlay.classList.remove("d-flex");
    overlay.classList.add("d-none");
}

document.addEventListener("DOMContentLoaded", () => {
    overlay = document.getElementById("overlay");
    loading = document.getElementById("loading");
    setInterval(() => {
        loadingDots++;
        loadingDots = loadingDots % 3;
        let text = "Loading";
        for (let i = 0; i <= loadingDots; i++) {
            text += ".";
        }
        loading.innerText = text;
    }, 500);
    OverlayOff();
});
