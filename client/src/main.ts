import "./style.css";

const messages = document.getElementById("messages") as HTMLUListElement;

const urlScheme = window.location.protocol === "https:" ? "wss" : "ws";
const host =
	window.location.hostname === "localhost"
		? `${window.location.hostname}:8080`
		: `${window.location.hostname}`;
const ws = new WebSocket(`${urlScheme}://${host}/ws`);
const messageInput = document.getElementById(
	"messageInput",
) as HTMLInputElement;

const message = {
	user_id: "",
	content: "",
};

messageInput.addEventListener("keydown", (event) => {
	if (event.key === "Enter") {
		if (!messageInput.value) {
			return;
		}

		// NEED TO DISPLAY WHO SENT THE MESSAGE/S...

		message.user_id = crypto.randomUUID();
		message.content = messageInput.value;

		const msgJSON = JSON.stringify(message);
		ws.send(msgJSON);

		appendMessage(messages, "client", messageInput.value);
		console.log(`[Client] ${messageInput.value}`);
		messageInput.value = "";
	}
});

function appendMessage(parent: HTMLUListElement, source: string, message: any) {
	const li = document.createElement("li");
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
	appendMessage(messages, "server", event.data);
};

ws.onclose = (event) => {
	console.log(`Connection closed ${event.reason}`);
};

ws.onerror = (error) => {
	console.log(`[WebSocket error] ${error}`);
};
