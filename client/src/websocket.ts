export function NewWebSocket(url: string) {
  const ws = new WebSocket(`${url}`);
  console.log(`[state ${ws.CONNECTING}] Socket has been created. The connection is not yet open.`);

  let message = "Hello from ws client.";

  // Establish connection to server
  ws.addEventListener("open", (event) => {
    console.log(`[state ${ws.OPEN}] The connection is open and ready to communicate: ${event}`);
    console.log(`[protocol] ${ws.protocol}`);
    ws.send(message)
  });

  // Handle connection errors
  ws.addEventListener("error", (event) => {
    console.log(`[error] ${event}`);
  });

  // Listen for incoming data from the server
  ws.addEventListener("message", (event) => {
    console.log(`[server message] Received data from server: ${event.data}`)
    try {
      const parsedData = JSON.parse(event.data);
      console.log(parsedData);
    } catch (error) {
      console.log(`[error] Received a non-JSON object or parsing error: ${event.data}`);
    }
  });

  // Handle connection closure
  ws.addEventListener("close", (event) => {
    if (event.wasClean) {
      console.log(`[accept] Connection successfully closed, code=${event.code}, reason=${event.reason}`);
    } else {
      console.warn(`[warn] Connection died unexpectedly`);
    }
  });
}