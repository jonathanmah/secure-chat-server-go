import { updateUsername, logout } from "../api.js";
import {
  messageInput,
  sendBtn,
  logoutBtn,
  editUsernameBtn,
  usernameInput,
  usernameModal,
  cancelUsernameBtn,
  saveUsernameBtn,
  darkModeToggle,
  joinRoomBtn,
  roomInput,
  confirmJoinRoomBtn,
  cancelJoinRoomBtn,
  joinRoomModal,
} from "./dom.js";
import {
  renderCharCount,
  resizeTextarea,
  renderUsername,
  loadDarkModePref,
  clearChatMessages,
} from "./ui.js";
import {
  initWebSocketConn,
  sendChatMessage,
  sendUsernameUpdateMessage,
} from "./websocket.js";

// -------------------------------------- SEND MESSAGE ----------------------------------
// send chat messages with enter
messageInput.addEventListener("keydown", (event) => {
  if (event.key === "Enter" && !event.shiftKey) {
    event.preventDefault();
    const text = messageInput.value.trim();
    if (text) {
      sendChatMessage(text);
      messageInput.value = "";
      renderCharCount();
    }
  }
});
// send chat messages with click
sendBtn.addEventListener("click", () => {
  const text = messageInput.value.trim();
  if (text !== "") {
    sendChatMessage(text);
    messageInput.value = "";
  }
});

// handles ui resizing
messageInput.addEventListener("input", () => {
  renderCharCount();
  resizeTextarea();
});

// -------------------------------------- EDIT USERNAME MODAL ----------------------------------
// handle editing username
editUsernameBtn.addEventListener("click", () => {
  usernameInput.value = users.find((u) => u.id === id)?.username || ""; // fill in username
  usernameModal.classList.remove("hidden");
  usernameInput.focus();
});

cancelUsernameBtn.addEventListener("click", () => {
  usernameModal.classList.add("hidden");
});

document.addEventListener("keydown", function (event) {
  if (event.key === "Escape" && !usernameModal.classList.contains("hidden")) {
    cancelUsernameBtn.click();
  }
});

saveUsernameBtn.addEventListener("click", async () => {
  const newUsername = usernameInput.value.trim();
  if (!/^[a-zA-Z0-9_]{3,20}$/.test(newUsername)) {
    alert(
      "Username must be 3–20 characters, only letters, numbers, and underscores."
    );
    return;
  }
  try {
    await updateUsername(newUsername); // updates db
  } catch (err) {
    console.error(err);
  }
  window.username = newUsername;
  sendUsernameUpdateMessage(newUsername); // updates hub client and triggers a broadcast
  renderUsername(newUsername);
  usernameModal.classList.add("hidden");
});

usernameInput.addEventListener("keydown", function (event) {
  if (event.key === "Enter") {
    event.preventDefault();
    saveUsernameBtn.click();
  }
});

// -------------------------------------- JOIN ROOM MODAL ----------------------------------
joinRoomBtn.addEventListener("click", () => {
  joinRoomModal.classList.remove("hidden");
});

cancelJoinRoomBtn.addEventListener("click", () => {
  joinRoomModal.classList.add("hidden");
  roomInput.value = "";
});

confirmJoinRoomBtn.addEventListener("click", () => {
  const newRoomID = roomInput.value.trim();
  if (newRoomID) {
    initWebSocketConn(newRoomID);
    clearChatMessages();
  }
  joinRoomModal.classList.add("hidden");
  roomInput.value = "";
});
// -------------------------------------- LOGOUT ----------------------------------
logoutBtn.addEventListener("click", async () => {
  try {
    await logout();
  } catch (err) {
    console.error(err);
    alert("Failed to log out");
  }
  window.location.reload();
  window.location.href = "/";
});

// -------------------------------------- DARKMODE ----------------------------------
// toggle and save darkmode
darkModeToggle.addEventListener("change", function () {
  document.body.classList.toggle("dark-mode", this.checked);
  localStorage.setItem("darkMode", this.checked);
});

window.addEventListener("DOMContentLoaded", () => {
  loadDarkModePref();
});
