const validateEmail = () => {
    const input = document.getElementById("email");
    const email = input.value.trim();
    const emailPattern = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    const submitBtn = document.querySelector("form button[type='submit']");

    if(email && emailPattern.test(email)) {
        submitBtn.removeAttribute("disabled");
    } else {
        submitBtn.setAttribute("disabled", "true");
    }
}