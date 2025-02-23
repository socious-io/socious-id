const onUploadProfile = () => {
    const input = document.getElementById("upload");
    const file = input.files[0];

    if (file) {
        const allowedTypes = ["image/jpeg", "image/png", "image/gif", "image/webp"];
        
        if (!allowedTypes.includes(file.type)) {
            alert("Please upload a valid image file (JPG, PNG, GIF, or WebP).");
            input.value = "";
            return;
        }

        const reader = new FileReader();
        reader.onload = function (e) {
            const avatar = document.getElementById("avatar");
            avatar.src = e.target.result; 
        };
        reader.readAsDataURL(file);
    }

}

const validateForm = () => {
    const first = document.getElementById("first").value.trim();
    const last = document.getElementById("last").value.trim();
    const username = document.getElementById("username").value.trim();
    const password = document.getElementById("pass").value.trim();
    const lengthHint = document.getElementById("length-hint");
    const charHint = document.getElementById("character-hint");
    const isPasswordValid = password.length >= 8 && /[\W_]/g.test(password);
    const submitBtn = document.querySelector("form button[type='submit']");

    lengthHint.src = `/statics/icons/check-${password.length >= 8 ? 'green': 'grey'}.svg`;
    
    charHint.src = `/statics/icons/check-${/[\W_]/.test(password) ? 'green': 'grey'}.svg`;

    if(first && last && username && isPasswordValid) {
        console.log(first, last, username, isPasswordValid)
        submitBtn.removeAttribute("disabled");
    } else {
        submitBtn.setAttribute("disabled", "true");
    }
}

const createProfile = () => {
    const first = document.getElementById("first").value.trim();
    const last = document.getElementById("last").value.trim();
    const username = document.getElementById("username").value.trim();
    const password = document.getElementById("pass").value.trim();
    const avatar = document.getElementById("avatar").src;

    //BE logic with first, last, username, password, avatar
    window.location.href = './congrats.html';
    return false;
}