const save = document.getElementById("save");
if (save != null) {
    save.addEventListener("click", Save);
}

const repeat = document.getElementById("repeat");
repeat.addEventListener("click", Repeat);

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
    const id = urlParams.get("id");
    console.log(id);
    window.location.href = `/?id=${id}`;
}
