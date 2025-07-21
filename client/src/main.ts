import { NewWebSocket } from "./websocket";

function main() {
  const port = "8080";
  const devUrl = `ws://${import.meta.env.VITE_DEV}:${port}/ws`;
  const localRemoteUrl = `ws://${import.meta.env.VITE_REMOTE}:${port}/ws`;

  NewWebSocket(devUrl);
}

export function handlerSubmit(ws: WebSocket) {
  let textbox = document.getElementById("message-input") as HTMLTextAreaElement;
  let message = textbox.value;
  appendClientMessageToWindow(message);
  textbox.value = "";

  ws.send(message);
}

export function appendClientMessageToWindow(message: string) {
  let messageWindow = document.getElementById("message-window");
  if (!messageWindow) {
    console.error("[error] element not found");
    return;
  }

  let messageDiv = document.createElement("p");
  messageDiv.classList.add("client-message");
  messageDiv.textContent = `Client: ${message}`;
  messageWindow.appendChild(messageDiv);
}

export function appendServerMessageToWindow(message: string) {
  let messageWindow = document.getElementById("message-window");
  if (!messageWindow) {
    console.error("[error] element not found");
    return;
  }

  let messageDiv = document.createElement("p");
  messageDiv.classList.add("server-message");
  messageDiv.textContent = `Server: ${message}`;
  messageWindow.appendChild(messageDiv);
}

main();