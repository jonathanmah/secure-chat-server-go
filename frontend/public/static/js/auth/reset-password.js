import { resetPassword } from "../api.js";
import { error, passwordElement, confirmPasswordElement } from "./dom.js";

document
  .getElementById("resetPasswordForm")
  .addEventListener("submit", async function (e) {
    e.preventDefault();

    const password = passwordElement.value;
    const confirm = confirmPasswordElement.value;
    const params = new URLSearchParams(window.location.search);
    const token = params.get("token");
    const passwordRegex = /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[\W_]).{8,}$/;
    if (!token) {
      error.textContent = "Missing or invalid token.";
      return;
    }
    if (password !== confirm) {
      error.textContent = "Passwords do not match.";
      return;
    }
    if (!passwordRegex.test(password)) {
      error.textContent =
        "Password must be at least 8 characters long and include uppercase, lowercase, number, and special character.";
      return;
    }
    error.textContent = "";
    try {
      await resetPassword(token, password);
    } catch (err) {
      error.textContent = err.message;
      return;
    }
    const form = document.getElementById("resetPasswordForm");
    form.style.display = "none";
    const successMsg = document.createElement("p");
    successMsg.textContent = "Password reset successful. You can now log in.";
    successMsg.style.color = "green";
    document.querySelector(".container").appendChild(successMsg);
  });
