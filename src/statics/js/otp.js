function handlePaste(event) {
    event.preventDefault(); 
    let inputs = document.querySelectorAll(".otp-input");
    let pasteData = (event.clipboardData || window.clipboardData).getData("text");
    pasteData = pasteData.replace(/\D/g, "");
    
    inputs.forEach((el, index) => {
        el.value = pasteData[index] || "";
    });
    
    const firstEmptyIndex = Array.from(inputs).findIndex(input => input.value === "");
    if(firstEmptyIndex !== -1) {
        inputs[firstEmptyIndex].focus();
    } else {
        inputs[inputs.length - 1].focus();
    }
}

function handleBackspace(event, input) {
    let inputs = document.querySelectorAll(".otp-input");
    let currentIndex = Array.from(inputs).indexOf(input);

    if (event.key === "Backspace" && input.value === "" && currentIndex > 0) {
        inputs[currentIndex - 1].value = ""; 
        inputs[currentIndex - 1].focus(); 
    }
}

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