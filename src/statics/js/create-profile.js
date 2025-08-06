const isValidUsername = (username) => {
  const allowedPattern = /^[a-z0-9._-]+$/;
  const startsWithInvalidChar = /^[._]/;
  const hasConsecutiveSpecials = /[._]{2,}/;

  return (
    username &&
    username.length >= 6 &&
    username.length <= 24 &&
    allowedPattern.test(username) &&
    !startsWithInvalidChar.test(username) &&
    !hasConsecutiveSpecials.test(username)
  );
};

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

		try {
			const mediaUploadResponse = await fetch("/media", {
				method: "POST",
				body: formData,
				headers: {},
			});

			if (mediaUploadResponse.ok) {
				const media = await mediaUploadResponse.json();
				avatar.src = media.url;
				avatarId.value = media.id;
			} else {
				console.error("Error:", mediaUploadResponse.statusText);
				return;
			}
		} catch (e) {
			console.error("Error:", e);
			return;
		}

		// const reader = new FileReader();
		// reader.onload = function (e) {
		//     const avatar = document.getElementById("avatar");
		//     avatar.src = e.target.result;
		// };
		// reader.readAsDataURL(file);
	}
};

const validateForm = () => {
    const first = document.getElementById("first").value.trim();
    const last = document.getElementById("last").value.trim();
    const username = document.getElementById("username").value.trim();
    const errorSpan = document.getElementById("username-error");
    const password = document.getElementById("pass").value.trim();
    const lengthHint = document.getElementById("length-hint");
    const charHint = document.getElementById("character-hint");
	const submitBtn = document.querySelector('button[data-event="on-submit-profile"]');
	
    
    //Password Validation
    const isPasswordValid = password.length >= 8 && /[\W_]/g.test(password);
    lengthHint.src = `/statics/icons/check-${password.length >= 8 ? 'green': 'grey'}.svg`;
    charHint.src = `/statics/icons/check-${/[\W_]/.test(password) ? 'green': 'grey'}.svg`;

    //Username Validation
    const isUsernameValid = isValidUsername(username);
    if (!isUsernameValid) {
        errorSpan.style.display = "block";
    } else {
        errorSpan.style.display = "none";
    }

    if(first && last && isUsernameValid && isPasswordValid) {
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
			avatar_id: avatarId.length > 0 ? avatarId : null,
		}),
		headers: { "Content-Type": "application/json" },
	})
	.then((response) => {
		console.log(response)
		if (response.ok) {
			if (response.redirected) {  // Detect if a redirect happened
				window.location.href = response.url;  // Redirect user
			} else {
				return response.text();  // Handle normal response
			}
		} else {
			response.text().then(html=>{
				const parser = new DOMParser();
				const doc = parser.parseFromString(html, "text/html");
				const error = window.document.getElementById("error");
				error.innerHTML = doc.querySelector("#error").innerHTML;
			});
		}
	})
	.catch((error) => console.error("Error:", error));

	return false;
};

const onUploadLogo = async () => {
	const input = document.getElementById("upload");
	const logoId = document.getElementById("logo-id");
	const file = input.files[0];

	if (file) {
		const allowedTypes = ["image/jpeg", "image/png", "image/gif", "image/webp"];

		if (!allowedTypes.includes(file.type)) {
			alert("Please upload a valid image file (JPG, PNG, GIF, or WebP).");
			input.value = "";
			return;
		}

		const logo = document.getElementById("logo");
		const formData = new FormData();
		formData.append("file", file);

		try {
			const mediaUploadResponse = await fetch("/media", {
				method: "POST",
				body: formData,
				headers: {},
			});

			if (mediaUploadResponse.ok) {
				const media = await mediaUploadResponse.json();
				logo.src = media.url;
				logoId.value = media.id;
			} else {
				console.error("Error:", mediaUploadResponse.statusText);
				return;
			}
		} catch (e) {
			console.error("Error:", e);
			return;
		}

		// const reader = new FileReader();
		// reader.onload = function (e) {
		//     const avatar = document.getElementById("avatar");
		//     avatar.src = e.target.result;
		// };
		// reader.readAsDataURL(file);
	}
};

const validateOrgForm = () => {
	const name = document.getElementById("name").value.trim();
    const shortname = document.getElementById("shortname").value.trim();
	const errorSpan = document.getElementById("shortname-error");
	const emailInput = document.getElementById("email");
	const email = emailInput.value.trim();
	const submitBtn = document.querySelector('button[data-event="create-organization"]');

    //Shortname Validation
    const isShortnameValid = isValidUsername(shortname);
    if (!isShortnameValid) {
        errorSpan.style.display = "block";
    } else {
        errorSpan.style.display = "none";
    }

	//Email Validation
	const isEmailValid = email && emailInput.checkValidity();

	if (name && isShortnameValid && isEmailValid) {
		submitBtn.removeAttribute("disabled");
	} else {
		submitBtn.setAttribute("disabled", "true");
	}
};

const createOrganization = () => {
	const name = document.getElementById("name").value.trim();
	const shortname = document.getElementById("shortname").value.trim();
	const email = document.getElementById("email").value.trim();
	const logoId = document.getElementById("logo-id").value.trim();
	const submitBtn = document.querySelector('button[data-event="create-organization"]');

	submitBtn.setAttribute("disabled", "true");
	
    fetch("/organizations/register", {
		method: "POST",
		body: JSON.stringify({
			name,
			shortname,
			email,
			logo_id: logoId.length > 0 ? logoId : null,
		}),
		headers: { "Content-Type": "application/json" },
	})
	.then((response) => {
		if (response.ok) {
			if (response.redirected) {  // Detect if a redirect happened
				window.location.href = response.url;  // Redirect user
			} else {
				submitBtn.removeAttribute("disabled");
				return response.text();  // Handle normal response
			}
		} else {
			response.text().then(html=>{
				const parser = new DOMParser();
				const doc = parser.parseFromString(html, "text/html");
				const error = window.document.getElementById("error");
				error.innerHTML = doc.querySelector("#error").innerHTML;
				submitBtn.removeAttribute("disabled");
			});
		}
	})
	.catch((error) => {
		console.error("Error:", error)
	});

	return false;
}   
