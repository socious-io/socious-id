function checkOTP() {
    const inputs = document.querySelectorAll(".otp-input");
    const otp = Array.from(inputs).map(input => input.value).join("");
    const submitBtn = document.querySelector("form button[type='submit']");
    const codeInput = document.querySelector("input[name='code']");
    
    if(otp.length === 6) {
        submitBtn.removeAttribute("disabled");
        codeInput.value = otp;
    } else {
        submitBtn.setAttribute("disabled", "true");
    }
}

const moveToNext = (input, nextIndex) => {
    const inputs = document.querySelectorAll(".otp-input");
    if (input.value.length === 1 && nextIndex < inputs.length) {
        inputs[nextIndex].focus();
    }
}