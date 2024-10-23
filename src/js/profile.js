document.addEventListener("DOMContentLoaded", OnLoaded);
const sort = document.getElementById("sort");
const type = document.getElementById("type");
const order = document.getElementById("order");
let json;
async function OnLoaded() {
    document.getElementById("sort").addEventListener("click", OnClick);
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
            render(json.runs, json.stats);

            OverlayOff();
        }
    } catch (error) {
        console.error(error.message);
    }
}

function OnClick() {
    const sortType = type.value;
    const sortOrder = order.value;

    const sortedRuns = json.runs.sort((a, b) => {
        if (sortOrder === "asc") {
            return a[sortType] - b[sortType];
        } else {
            return b[sortType] - a[sortType];
        }
    });

    const body = document.getElementById("body");
    body.innerHTML = "";
    render(sortedRuns, json.stats);
}

function render(runs, stats) {
    const body = document.getElementById("body");
    for (let i = 0; i < runs.length; i++) {
        const run = runs[i];
        body.innerHTML += renderRunCard(run, stats);
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
            <p class="card-text">
                <strong>Accuracy:</strong>
                <span class="${getClass(run.accuracy, stats.accuracy)}">
                    ${run.accuracy}%
                    ${getIcon(run.accuracy, stats.accuracy)}
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
