const testButton = document.getElementById("test-button")

const time = new EventSource('/time');
time.addEventListener('time', (e) => {
    document.getElementById("actual-time").innerHTML = e.data;
}, false);

testButton.addEventListener("click", function() {
    alert("Button Was Clicked")
})

