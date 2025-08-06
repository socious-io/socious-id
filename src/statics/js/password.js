const togglePassword = (icon) => {
    const passwordInput = icon.previousElementSibling; 
    if (passwordInput.type === 'password') {
        passwordInput.type = 'text';
        icon.src = "/statics/icons/eye.svg";
    } else {
        passwordInput.type = 'password';
        icon.src = "/statics/icons/eye.svg";
    }
}

const validatePassword = () => {
    const password = document.getElementById("pass").value;
    const confirmPassword = document.getElementById("confirm-pass").value;
    const lengthHint = document.getElementById("length-hint");
    const charHint = document.getElementById("character-hint");
    const submitBtn = document.querySelector("form button[type='submit']");
    const isValidPassword = password.length >= 8 && /[\W_]/.test(password);
    const isMatchingPassword = password === confirmPassword;

    lengthHint.src = `/statics/icons/check-${password.length >= 8 ? 'green': 'grey'}.svg`;
    
    charHint.src = `/statics/icons/check-${/[\W_]/.test(password) ? 'green': 'grey'}.svg`;

    if (isValidPassword && isMatchingPassword) {
        submitBtn.removeAttribute("disabled");
    } else {
        submitBtn.setAttribute("disabled", "true");
    }
}