// const topBarConfig = {}

function startLoader(){
  console.log("calling start loader")
    topbar.show()
    document.getElementById("loader").style.opacity = 1;
    document.getElementById("loader").hidden = false;
}

function stopLoader(){
    topbar.hide()
    document.getElementById("loader").style.opacity = 0;
    document.getElementById("loader").hidden = true;
}

(function check() {
  if (document.body) {
    startLoader();
  } else {
    requestAnimationFrame(check); // wait for next frame
  }
})();

