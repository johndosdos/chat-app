import { appendServerMessageToWindow, handlerSubmit } from "./main";

export function NewWebSocket(url: string) {
  const ws = new WebSocket(`${url}`);
  console.log(`${new Date().toLocaleString()} [state] Socket has been created. The connection is not yet open.`);

  // Establish connection to server
  let submitBtn = document.getElementById("submit-button");

  ws.addEventListener("open", () => {
    console.log(`${new Date().toLocaleString()} [state] The connection is open and ready to communicate`);

    submitBtn?.addEventListener("click", (event) => {
      event.preventDefault();

      handlerSubmit(ws);
    })
  });

  // Handle connection errors
  ws.addEventListener("error", (event) => {
    console.log(`[error] ${event}`);
  });

  // Listen for incoming data from the server
  ws.addEventListener("message", (event) => {
    const parsedData = event.data;
    appendServerMessageToWindow(parsedData);
    console.log(`${new Date().toLocaleDateString()} data: ${parsedData}`);
  });

  // Handle connection closure
  ws.addEventListener("close", (event) => {
    if (event.wasClean) {
      console.log(`${new Date().toLocaleString()} [accept] Connection successfully closed, code=${event.code}, reason=${event.reason}`);
    } else {
      console.warn(`[warn] Connection died unexpectedly`);
    }
  });
}