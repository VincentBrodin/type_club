document.addEventListener("DOMContentLoaded", () => {
    const passwords = document.getElementsByClassName("password");
    for (let i = 0; i < passwords.length; i++) {
        const password = passwords[i];
        const input = password.getElementsByTagName("input")[0];
        const button = password.getElementsByTagName("button")[0];
        const icon = password.getElementsByTagName("i")[0];
        button.addEventListener("click", () => {
            if (input.type == "text") {
                icon.classList.add("fa-eye");
                icon.classList.remove("fa-eye-slash");
                input.type = "password";
            } else if (input.type == "password") {
                icon.classList.add("fa-eye-slash");
                icon.classList.remove("fa-eye");
                input.type = "text";
            }
        });
    }
});
