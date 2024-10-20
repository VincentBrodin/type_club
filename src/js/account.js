const username = document.getElementById("username");
username.addEventListener("input", UsernameChanged);
const currentUsername = username.value;

const updateUsername = document.getElementById("updateUsername");
updateUsername.disabled = true;
updateUsername.addEventListener("click", UpdateUsername);

const email = document.getElementById("email");
email.addEventListener("input", EmailChanged);
const currentEmail = email.value;

const updateEmail = document.getElementById("updateEmail");
updateEmail.disabled = true;
updateEmail.addEventListener("click", UpdateEmail);

function UpdateUsername() {
    console.log(username.value);
}

function UsernameChanged() {
    updateUsername.disabled = username.value == currentUsername;
}

function UpdateEmail() {
    console.log(email.value);
}

function EmailChanged() {
    updateEmail.disabled = email.value == currentEmail;
}
