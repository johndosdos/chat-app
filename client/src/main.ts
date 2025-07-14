import { NewWebSocket } from "./websocket";

function main() {
  const port = "8600";
  let url = `ws://localhost:${port}/ws`;

  NewWebSocket(url);
}

main();