import React from "react";
import ReactDOM from "react-dom/client";

const root = document.getElementById("root");

if (root) {
  ReactDOM.createRoot(root).render(
    <React.StrictMode>
      <main data-testid="platform-shell">
        <h1>Platform MVP</h1>
        <p>Local stack ready</p>
      </main>
    </React.StrictMode>,
  );
}
