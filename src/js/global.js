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
