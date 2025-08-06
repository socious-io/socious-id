function logout(){
    window.location.href = '/auth/logout?t=' + Date.now();
    return false;
}