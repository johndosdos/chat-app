const messages = document.getElementById("messages");

const ws = new WebSocket("ws://localhost:8080/ws");
const messageInput = document.getElementById("messageInput")

messageInput.addEventListener("keydown", event => {
  if (event.key === "Enter") {
    if (!messageInput.value) {
      return;
    }

    ws.send(messageInput.value);
    appendMessage(messages, "client", messageInput.value);
    console.log(`[Client] ${messageInput.value}`);
    messageInput.value = "";
  }
})

function appendMessage(parent, source, message) {
  let li = document.createElement("li");
  li.textContent = message;

  if (source === "server") {
    li.className = "received-message";
  } else if (source === "client") {
    li.className = "sent-message";
  }
  parent.appendChild(li);
}

ws.onopen = () => {
  console.log("Connected to server");
};

ws.onmessage = (event) => {
  console.log(`[Server] ${event.data}`);
  appendMessage(messages, "server", event.data)
};

ws.onclose = (event) => {
  console.log(`[Connection closed] ${event.reason}`);
};

ws.onerror = (error) => {
  console.log(`[WebSocket error] ${error}`)
}
