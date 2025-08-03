function logout(){
    window.location.href = '/auth/logout?t=' + Date.now();
    return false;
}

function confirmIdentity(target){
    const identityIdInput = document.querySelector("#identity_id");
    identityIdInput.value = target.getAttribute("data-identity-id");
    target.form.submit();
}