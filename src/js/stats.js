const save = document.getElementById("save");
save.addEventListener("click", Save);

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
        if (response.ok) {
            console.log("OK");
        }
        OverlayOff();
    } catch (error) {
        console.error(error.message);
    }
}
