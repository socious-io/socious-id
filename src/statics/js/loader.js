// const topBarConfig = {}

function startLoader(){
    topbar.show()
    document.getElementById("loader").style.opacity = 1;
    document.getElementById("loader").hidden = false;
}

function stopLoader(){
    topbar.hide()
    document.getElementById("loader").style.opacity = 0;
    document.getElementById("loader").hidden = true;
}

const observer = new MutationObserver(() => {
  const body = document.body;
  if (body) {
    observer.disconnect();
    startLoader()
  }
});

observer.observe(document.documentElement, { childList: true });

