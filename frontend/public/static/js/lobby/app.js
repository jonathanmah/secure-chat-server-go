import { getUserInfo } from "../api.js";
import { renderUsername } from "./ui.js";
import { initWebSocketConn } from "./websocket.js";
import "./events.js";

async function init() {
  try {
    const data = await getUserInfo();
    window.id = data.id;
    window.username = data.username;
    window.users = [];
    renderUsername(data.username);
    initWebSocketConn();
  } catch (err) {
    console.error(err);
  }
}

init();
