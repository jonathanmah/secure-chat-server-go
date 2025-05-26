import { emailElement, passwordElement, error } from "./dom.js";
import { login } from "../api.js";
import { SERVER_BASE_URL } from "../config.js";

document
  .getElementById("login-form")
  .addEventListener("submit", async function (e) {
    e.preventDefault();
    const email = emailElement.value.trim();
    const password = passwordElement.value.trim();
    if (!email || !password) {
      error.textContent = "Please fill in both fields.";
      return;
    }
    error.textContent = "";
    try {
      await login(email, password);
    } catch (err) {
      error.textContent = err.message;
      return;
    }
    window.location.href = "/lobby";
  });

document.getElementById("googleLogin").addEventListener("click", function () {
  // Redirect the user to go server OAuth handler
  window.location.href = `${SERVER_BASE_URL}/login/google`; //#TODO update localhost
});

// shows success status and cleans up URL
document.addEventListener("DOMContentLoaded", () => {
  const params = new URLSearchParams(window.location.search);
  const msg = params.get("msg");
  if (msg) {
    const statusDiv = document.getElementById("statusMessage");
    statusDiv.textContent = msg;
    statusDiv.style.display = "block";
    if (window.history.replaceState) {
      // msg param from the URL without reloading
      const url = new URL(window.location);
      url.searchParams.delete("msg");
      window.history.replaceState({}, "", url);
    }
  }
});
