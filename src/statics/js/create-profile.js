const onUploadProfile = async () => {
    const input = document.getElementById("upload");
    const avatarId = document.getElementById("avatar-id");
    const file = input.files[0];

    if (file) {
        const allowedTypes = ["image/jpeg", "image/png", "image/gif", "image/webp"];
        
        if (!allowedTypes.includes(file.type)) {
            alert("Please upload a valid image file (JPG, PNG, GIF, or WebP).");
            input.value = "";
            return;
        }

        const avatar = document.getElementById("avatar");
        const formData = new FormData();
        formData.append("file", file);

        try{
            const mediaUploadResponse = await fetch("/media", {
                method: "POST",
                body: formData,
                headers: {},
            })

            if (mediaUploadResponse.ok) {
                const media = await mediaUploadResponse.json()
                avatar.src = media.url
                avatarId.value = media.id
            } else {
                console.error("Error:", mediaUploadResponse.statusText);
                return
            }
        }catch(e){
            console.error("Error:", e)
            return
        }

        // const reader = new FileReader();
        // reader.onload = function (e) {
        //     const avatar = document.getElementById("avatar");
        //     avatar.src = e.target.result; 
        // };
        // reader.readAsDataURL(file);
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
    const first_name = document.getElementById("first").value.trim();
    const last_name = document.getElementById("last").value.trim();
    const username = document.getElementById("username").value.trim();
    const password = document.getElementById("pass").value.trim();
    const avatarId = document.getElementById("avatar-id").value.trim();

    fetch("/users/profile", {
        method: "PUT",
        body: JSON.stringify({
            first_name,
            last_name,
            username,
            password,
            avatar_id: avatarId,
        }),
        headers: { "Content-Type": "application/json" }
    })
    .then(response => {
        if (response.ok) {
            // Redirect after successful PUT request
            window.location.href = "/auth/confirm";
        } else {
            console.error("Error:", response.statusText);
        }
    })
    .catch(error => console.error("Error:", error));

    return false;
}