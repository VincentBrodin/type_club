console.log("Hello world")

document.addEventListener("DOMContentLoaded", () => {
    const elementsToAnimate = document.querySelectorAll('.fade-in-target');
    elementsToAnimate.forEach(el => {
        el.classList.add('fade-in');
    });
});
