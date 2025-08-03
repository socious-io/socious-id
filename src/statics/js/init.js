function initOnClicks(){
	document.querySelector('button[data-event="logout"]')?.addEventListener("click",logout);
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
	document.querySelectorAll('[data-event="select-account"]').forEach((el)=>{
		el.addEventListener("click",(e)=>confirmIdentity(e.target));
	});
	document.querySelectorAll('[data-event="display-email"]').forEach((el)=>{
		el.addEventListener("load",displayEmail);
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
	document.querySelectorAll('[data-event="otp-input"]').forEach((el, idx, all)=>{
		el.addEventListener("input", (e) => moveToNext(e.target, idx+1, all));
		el.addEventListener("keyup", checkOTP);
		el.addEventListener("paste", handlePaste);
		el.addEventListener("keydown", (e) => handleBackspace(e, e.target));
	});
}

document.addEventListener("DOMContentLoaded", initOnClicks);
