import { sendPasswordResetEmail } from "../api.js";
import { emailElement } from "./dom.js";

document
  .getElementById("forgot-password-form")
  .addEventListener("submit", async function (e) {
    e.preventDefault();
    const email = emailElement.value;

    error.textContent = "";
    try {
      await sendPasswordResetEmail(email);
    } catch (err) {
      error.textContent = err.message;
      return;
    }
    const form = document.getElementById("forgot-password-form");
    form.style.display = "none";
    const successMsg = document.createElement("p");
    successMsg.textContent =
      "If this email exists, a reset link has been sent.";
    successMsg.style.color = "green";
    document.querySelector(".container").appendChild(successMsg);
  });
