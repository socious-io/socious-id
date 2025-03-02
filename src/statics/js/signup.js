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

const handleSignUp = ()  => {
    const emailInput = document.getElementById("email").value.trim();
    if (!emailInput) return;

    window.location.href = `./otp?email=${emailInput}`;
    return false;
}

const displayEmail = () => {
    const params = new URLSearchParams(window.location.search);
    const email = params.get("email");
    if (email) {
        document.getElementById("user-email").textContent = decodeURIComponent(email);
    }
}

const verifyOTP = () => {
    const inputs = document.querySelectorAll(".otp-input");
    const otp = Array.from(inputs).map(input => input.value).join("");
    const params = new URLSearchParams(window.location.search);
    const email = params.get("email");

    if (otp.length !== 6 || isNaN(otp)) {
        return false; 
    }

    //BE logic with email
    window.location.href = './create-profile.html';
    return false;
}