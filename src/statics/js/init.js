function initOnClicks(){
	document.querySelector('button[data-event="create-organization"]')?.addEventListener("click", createOrganization);
	document.querySelector('button[data-event="on-submit-profile"]')?.addEventListener("click", createProfile);
	document.querySelectorAll('input[data-event="validate-password"]').forEach((el)=>{
		el.addEventListener("input",validatePassword);
	});
	document.querySelectorAll('input[data-event="validate-email"]').forEach((el)=>{
		el.addEventListener("input",(e)=>validateEmail(e.target));
	});
	document.querySelectorAll('input[data-event="validate-organization-form"]').forEach((el)=>{
		el.addEventListener("input",validateOrgForm);
	});
	document.querySelectorAll('input[data-event="validate-form"]').forEach((el)=>{
		el.addEventListener("input",validateForm);
	});
	document.querySelectorAll('img[data-event="toggle-password"]').forEach((el)=>{
		el.addEventListener("click",(e)=>togglePassword(e.target));
	});
	document.querySelectorAll('input[data-event="on-upload-logo"]').forEach((el)=>{
		el.addEventListener("input",onUploadLogo);
	});
	document.querySelectorAll('input[data-event="on-upload-avatar"]').forEach((el)=>{
		el.addEventListener("input",onUploadProfile);
	});
	document.querySelectorAll('input[data-event="on-submit-profile"]').forEach((el)=>{
		el.addEventListener("submit",createProfile);
	});
	document.querySelectorAll('[data-event="change-language"]').forEach((el)=>{
		el.addEventListener("click",(e)=>changeLanguage(e.target, e.target.getAttribute('data-value')));
	});
	document.querySelectorAll('[data-event="toggle-lang-dropdown"]').forEach((el)=>{
		el.addEventListener("click",()=>toggleDropdown('desktop-options'));
	});
	document.querySelectorAll('[data-event="toggle-lang-dropdown-mobile"]').forEach((el)=>{
		el.addEventListener("click",()=>toggleDropdown('mobile-options'));
	});
	document.querySelectorAll('[data-event="otp-input"]').forEach((el, idx, all)=>{
		el.addEventListener("input", (e) => moveToNext(e.target, idx+1, all));
		el.addEventListener("keyup", checkOTP);
		el.addEventListener("paste", handlePaste);
		el.addEventListener("keydown", (e) => handleBackspace(e, e.target));
	});

	stopLoader();
}

document.addEventListener("DOMContentLoaded", initOnClicks);
