const save = document.getElementById("save");
if (save != null) {
    save.addEventListener("click", Save);
}

const repeat = document.getElementById("repeat");
repeat.addEventListener("click", Repeat);

const share = document.getElementById("share");
if (share != null) {
    share.addEventListener("click", Share);
}

async function Save(event) {
    event.preventDefault();
    try {
        OverlayOn();
        const response = await fetch("/save", {
            method: "POST",
            headers: {
                "Content-Type": "application/json; charset=UTF-8",
            },
        });
        if (response.redirected) {
            window.location.href = response.url;
        }
        OverlayOff();
    } catch (error) {
        console.error(error.message);
    }
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

function Share(event) {
    event.preventDefault();
    const path = window.location.toString();
    navigator.clipboard.writeText(path);
    Noti("type_club", "Saved to clipboard");
}
