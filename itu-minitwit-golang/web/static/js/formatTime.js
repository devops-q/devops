window.addEventListener("load", (event) => formatTimeLocal)
function formatTimeLocal() {
    var times = document.getElementsByClassName("time")
    for (let element of times) {
        var datetime = new Date(element.innerHTML)
        element.innerHTML = "Written @ " + datetime.toLocaleString()
    }
}