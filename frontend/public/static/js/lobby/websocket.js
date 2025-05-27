import { HUB_BASE_URL } from "../config.js";
import {
  renderRoomHeader,
  renderChatMessage,
  renderActiveUsers,
} from "./ui.js";

let socket = null;

const MessageType = {
  Chat: "chat",
  UsernameUpdate: "username_update",
  UserList: "userlist",
};

// initializes connection with server hub
export function initWebSocketConn(roomID) {
  if (
    socket &&
    (socket.readyState === WebSocket.OPEN ||
      socket.readyState === WebSocket.CONNECTING)
  ) {
    socket.close(1000); // 1000 for normal close
  }

  socket = new WebSocket(`${HUB_BASE_URL}/ws?room_id=${roomID}`);
  renderRoomHeader(roomID);
  // upgrader.Upgrade() in Go server will trigger this, once updating protocol from HTTP1.1 to WebSocket
  socket.addEventListener("open", () => {
    console.log("WebSocket connected");
  });
  // triggered on clean and abnormal closes
  socket.onclose = (e) => console.warn("WebSocket closed", e);
  // connection failed to establish, transmission error, or CORS/TLS issue
  socket.onerror = (e) => console.error("WebSocket error", e);

  // -------------------------------------- WebSocket Receive ----------------------------------
  socket.addEventListener("message", (event) => {
    const data = JSON.parse(event.data);
    console.log("Received from server: ", data);
    switch (data.type) {
      case MessageType.Chat:
        renderChatMessage(data.payload);
        break;
      case MessageType.UserList:
        window.users = data.payload.users;
        renderActiveUsers(data.payload.users);
        break;
      default:
        console.warn("WebSocket message type not supported: ", data.type);
    }
  });
}

// -------------------------------------- WebSocket Send ----------------------------------
export function sendChatMessage(text) {
  let message = JSON.stringify({
    type: MessageType.Chat,
    payload: {
      text: text,
    },
  });
  sendMessage(message);
  console.log("Message sent");
}
export function sendUsernameUpdateMessage(username) {
  let message = JSON.stringify({
    type: MessageType.UsernameUpdate,
    payload: {
      username: username,
    },
  });
  sendMessage(message);
}

function sendMessage(message) {
  if (socket.readyState === WebSocket.OPEN) {
    socket.send(message);
  } else {
    console.warn("WebSocket not ready to send");
  }
}
