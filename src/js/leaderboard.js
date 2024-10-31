document.addEventListener("DOMContentLoaded", OnLoaded);

async function OnLoaded() {
    try {
        OverlayOn();
        const response = await fetch("/leaderboard", {
            method: "POST",
        });
        if (response.ok) {
            const json = await response.json();
            console.log(json);
            render(json);
            OverlayOff();
        }
    } catch (error) {
        console.error(error.message);
    }
}

function render(json) {
    const body = document.getElementById("body");
    for (let i = 0; i < json.length; i++) {
        const data = json[i];
        body.innerHTML += renderCard(data, i + 1);
    }
}

function renderCard(data, rank) {
    return `
        <div class="card mb-3">
            <div class="card-body">
                <p class="card-text d-flex justify-content-between">
                    <span>
                        <strong>${rank}</strong>
                        <span>|</span>
                        <a href="/profile?id=${data.user.id}">${data.user.username}</a>
                    </span>
                    <a href="/stats?id=${data.run.id}">${data.run.awpm} WPM</a>
                </p>
            </div>
        </div>

`;
}
