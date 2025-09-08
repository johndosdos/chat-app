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
		messageInput.value = "";
	}
});

ws.onopen = () => {
	console.log("Connected to server");
};

ws.onmessage = (event) => {
	appendMessage(messages, "server", event.data);
};

ws.onclose = (event) => {
	console.log(`Connection closed ${event.reason}`);
};

ws.onerror = (error) => {
	console.log(`[WebSocket error] ${error}`);
};

function appendMessage(
	parent: HTMLUListElement,
	source: string,
	message: string,
) {
	const li = document.createElement("li");
	const div = document.createElement("div");

	li.className = "flex flex-col w-full my-0.5";
	div.className = "rounded-2xl break-words p-2.5";

	if (source === "server") {
		const span = document.createElement("span");

		li.className += " items-start";
		div.className += " bg-gray-200 text-gray-800 ";
		span.className = "text-xs text-gray-500 pl-1";
		span.textContent = "server";

		li.appendChild(span);
		li.appendChild(div);
	} else if (source === "client") {
		li.className += " items-end";
		div.className += " bg-blue-500 text-white";

		li.appendChild(div);
	}
	div.textContent = message;
	parent.appendChild(li);
}
