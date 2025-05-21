const socket = new WebSocket("ws://localhost:8080/ws");

socket.onopen = () => {
    console.log("Połączono z serwerem");
};

socket.onmessage = (event) => {
    const messages = document.getElementById("messages");
    const li = document.createElement("li");
    li.textContent = event.data;
    messages.appendChild(li);
};

function sendMessage() {
    const input = document.getElementById("messageInput");
    socket.send(input.value);
    input.value = "";
}
