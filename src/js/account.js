const username = document.getElementById("username");
username.addEventListener("input", UsernameChanged);
let currentUsername = username.value;

const updateUsername = document.getElementById("updateUsername");
updateUsername.disabled = true;
updateUsername.addEventListener("click", UpdateUsername);

const usernameTaken = document.getElementById("usernameTaken");
usernameTaken.innerText = "";

const email = document.getElementById("email");
email.addEventListener("input", EmailChanged);
let currentEmail = email.value;

const updateEmail = document.getElementById("updateEmail");
updateEmail.disabled = true;
updateEmail.addEventListener("click", UpdateEmail);

const emailTaken = document.getElementById("emailTaken");
emailTaken.innerText = "";

const overlay = document.getElementById("overlay");
const loading = document.getElementById("loading");
let loadingDots = 1;
OverlayOn();

async function UpdateUsername() {
    OverlayOn();

    const data = {
        username: username.value,
    };
    const validate = await Validate(data);
    if (validate.username) {
        username.classList.remove("border");
        username.classList.remove("border-danger");
        usernameTaken.innerText = "";
    } else {
        username.classList.add("border");
        username.classList.add("border-danger");
        usernameTaken.innerText = "Username taken";
        OverlayOff();
        return;
    }
    const update = await Update(data);
    if (update) {
        currentUsername = username.value;
    } else {
        username.value = currentUsername;
    }

    UsernameChanged();
    OverlayOff();
}

function UsernameChanged() {
    updateUsername.disabled = username.value == currentUsername;
}

async function UpdateEmail() {
    OverlayOn();

    const data = {
        email: email.value,
    };

    const validate = await Validate(data);
    if (validate.email) {
        email.classList.remove("border");
        email.classList.remove("border-danger");
        emailTaken.innerText = "";
    } else {
        email.classList.add("border");
        email.classList.add("border-danger");
        emailTaken.innerText = "Email taken";
        OverlayOff();
        return;
    }

    const update = await Update(data);
    if (update) {
        currentEmail = email.value;
    } else {
        email.value = currentEmail;
    }

    EmailChanged();
    OverlayOff();
}

function EmailChanged() {
    updateEmail.disabled = email.value == currentEmail;
}

async function Update(data) {
    try {
        const response = await fetch("/update", {
            method: "POST",
            body: JSON.stringify(data),
            headers: {
                "Content-Type": "application/json; charset=UTF-8",
            },
        });
        return response.ok;
    } catch (error) {
        console.error(error.message);
        return false;
    }
}

async function Validate(data) {
    try {
        const response = await fetch("/validate", {
            method: "POST",
            body: JSON.stringify(data),
            headers: {
                "Content-Type": "application/json; charset=UTF-8",
            },
        });

        if (response.ok) {
            json = await response.json();
            return json;
        }
        return {
            username: false,
            email: false,
        };
    } catch (error) {
        console.error(error.message);
        return {
            username: false,
            email: false,
        };
    }
}

function OverlayOn() {
    overlay.classList.add("d-flex");
    overlay.classList.remove("d-none");
}

function OverlayOff() {
    overlay.classList.remove("d-flex");
    overlay.classList.add("d-none");
}

document.addEventListener("DOMContentLoaded", () => {
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
