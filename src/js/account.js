//Username
const username = document.getElementById("username");
username.addEventListener("input", UsernameChanged);
let currentUsername = username.value;

const updateUsername = document.getElementById("updateUsername");
updateUsername.disabled = true;
updateUsername.addEventListener("click", UpdateUsername);

const usernameTaken = document.getElementById("usernameTaken");
usernameTaken.innerText = "";

//Email
const email = document.getElementById("email");
email.addEventListener("input", EmailChanged);
let currentEmail = email.value;

const updateEmail = document.getElementById("updateEmail");
updateEmail.disabled = true;
updateEmail.addEventListener("click", UpdateEmail);

const emailTaken = document.getElementById("emailTaken");
emailTaken.innerText = "";

//Password
const currentPassword = document.getElementById("currentPassword");
const correctPassword = document.getElementById("correctPassword");
correctPassword.innerText = "";

const newPassword = document.getElementById("newPassword");
const repeatNewPassword = document.getElementById("repeatNewPassword");
const updatePassword = document.getElementById("updatePassword");
const passwordMatch = document.getElementsByClassName("matchPassword");
for (let i = 0; i < passwordMatch.length; i++) {
    passwordMatch[i].innerText = "";
}
updatePassword.addEventListener("click", UpdatePassword);

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
        Noti(
            "type_club",
            `Updated your username from ${currentUsername} to ${username.value}.`,
        );

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
        Noti(
            "type_club",
            `Updated your email from ${currentEmail} to ${email.value}.`,
        );
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

async function UpdatePassword() {
    OverlayOn();
    valid = await CheckPassword(currentPassword.value);

    if (!valid) {
        currentPassword.classList.add("border");
        currentPassword.classList.add("border-danger");
        correctPassword.innerText = "Wrong password";
        OverlayOff();
        return;
    } else {
        currentPassword.classList.remove("border");
        currentPassword.classList.remove("border-danger");
        correctPassword.innerText = "";
    }

    if (newPassword.value != repeatNewPassword.value) {
        newPassword.classList.add("border");
        newPassword.classList.add("border-danger");
        repeatNewPassword.classList.add("border");
        repeatNewPassword.classList.add("border-danger");

        for (let i = 0; i < passwordMatch.length; i++) {
            passwordMatch[i].innerText = "Password must match";
        }

        OverlayOff();
        return;
    }
    newPassword.classList.remove("border");
    newPassword.classList.remove("border-danger");
    repeatNewPassword.classList.remove("border");
    repeatNewPassword.classList.remove("border-danger");

    for (let i = 0; i < passwordMatch.length; i++) {
        passwordMatch[i].innerText = "";
    }

    const data = {
        password: newPassword.value,
    };
    const update = await Update(data);
    if (update) {
        currentPassword.value = "";
        newPassword.value = "";
        repeatNewPassword.value = "";
    }
    OverlayOff();
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

async function CheckPassword(password) {
    try {
        const response = await fetch("/check", {
            method: "POST",
            body: JSON.stringify({
                password: password,
            }),
            headers: {
                "Content-Type": "application/json; charset=UTF-8",
            },
        });

        if (response.ok) {
            json = await response.json();
            return json.valid;
        }
        return false;
    } catch (error) {
        console.error(error.message);
        return false;
    }
}
