document.addEventListener("DOMContentLoaded", OnLoaded);

async function OnLoaded() {
    try {
        OverlayOn();
        const queryString = window.location.search;
        const urlParams = new URLSearchParams(queryString);
        const id = urlParams.get("id");
        console.log(id);

        const response = await fetch("/profile", {
            method: "POST",
            body: JSON.stringify({
                id: parseInt(id),
            }),
            headers: {
                "Content-Type": "application/json; charset=UTF-8",
            },
        });
        if (response.ok) {
            json = await response.json();
            document.getElementById("awpm").innerText = json.stats.awpm;
            document.getElementById("wpm").innerText = json.stats.wpm;
            document.getElementById("accuracy").innerText = json.stats.accuracy;

            const body = document.getElementById("body");
            for (let i = 0; i < json.runs.length; i++) {
                const run = json.runs[i];
                body.innerHTML += renderRunCard(run, json.stats);
            }
            OverlayOff();
        }
    } catch (error) {
        console.error(error.message);
    }
}

function renderRunCard(run, stats) {
    const getClass = (value, compare) => value > compare ? "green" : "red";
    const getIcon = (value, compare) =>
        value > compare
            ? '<i class="fa-solid fa-arrow-up"></i>'
            : '<i class="fa-solid fa-arrow-down"></i>';

    return `
    <div class="card mb-3">
        <div class="card-body">
            <!-- Run details -->
            <p class="card-text"><strong>Target:</strong> ${run.target}</p>
            <p class="card-text">
                <strong>Accuracy:</strong>
                <span class="${getClass(run.accuracy, stats.accuracy)}">
                    ${run.accuracy}%
                    ${getIcon(run.accuracy, stats.accuracy)}
                </span>
            </p>
            <p class="card-text">
                <strong>WPM:</strong>
                <span class="${getClass(run.awpm, stats.awpm)}">
                    ${run.awpm}
                    ${getIcon(run.awpm, stats.awpm)}
                </span>
            </p>
            <p class="card-text">
                <strong>Raw WPM:</strong>
                <span class="${getClass(run.wpm, stats.wpm)}">
                    ${run.wpm}
                    ${getIcon(run.wpm, stats.wpm)}
                </span>
            </p>
            <!-- Time -->
            <p class="card-text"><strong>Time:</strong> ${run.time} seconds</p>
            <!-- Try link -->
            <div class="d-flex flex-row">
                <a class="me-2" href="/?id=${run.id}">Try</a>
                <p class="me-2">|</p>
                <a class="me-2" href="stats/?id=${run.id}">Stats</a>
            </div>
        </div>
        </div>
    `;
}
