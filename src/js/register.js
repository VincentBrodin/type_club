const form = document.getElementById("form");

const username = document.getElementById("username");
const usernameTaken = document.getElementById("usernameTaken");
usernameTaken.innerText = "";

const email = document.getElementById("email");
const emailTaken = document.getElementById("emailTaken");
emailTaken.innerText = "";

form.addEventListener("submit", OnSubmit);

async function OnSubmit(event) {
    event.preventDefault();
    try {
        data = {
            username: username.value,
            email: email.value,
        };

        const response = await fetch("/validate", {
            method: "POST",
            body: JSON.stringify(data),
            headers: {
                "Content-Type": "application/json; charset=UTF-8",
            },
        });

        if (response.ok) {
            json = await response.json();
            if (json.username) {
                username.classList.remove("border");
                username.classList.remove("border-danger");
                usernameTaken.innerText = "";
            } else {
                username.classList.add("border");
                username.classList.add("border-danger");
                usernameTaken.innerText = "Username taken";
            }
            if (json.email) {
                email.classList.remove("border");
                email.classList.remove("border-danger");
                emailTaken.innerText = "";
            } else {
                email.classList.add("border");
                email.classList.add("border-danger");
                emailTaken.innerText = "Email taken";
            }

            if (json.username && json.email) {
                form.submit();
            }
        } else if (!response.ok) {
            const errorText = await response.text();
            console.error(errorText);
        }
    } catch (error) {
        console.error(error.message);
    }
}
