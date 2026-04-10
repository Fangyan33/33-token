import React from "react";
import ReactDOM from "react-dom/client";
import { BrowserRouter, Link, Route, Routes } from "react-router-dom";
import { PublicHomePage } from "./routes/public/home";
import { ConsoleHomePage } from "./routes/console/home";
import { AdminHomePage } from "./routes/admin/home";

function AppShell() {
  return (
    <div
      style={{
        minHeight: "100vh",
        background:
          "linear-gradient(160deg, #f7f0dd 0%, #efe6d2 45%, #d9e0d7 100%)",
        color: "#1f2a1f",
        fontFamily: '"Iowan Old Style", "Palatino Linotype", serif'
      }}
    >
      <header style={{ padding: "24px 32px", borderBottom: "1px solid #b9bea9" }}>
        <div style={{ display: "flex", gap: 24, alignItems: "center" }}>
          <strong style={{ letterSpacing: "0.08em" }}>MODEL API PLATFORM</strong>
          <nav style={{ display: "flex", gap: 16 }}>
            <Link to="/">官网</Link>
            <Link to="/console">用户控制台</Link>
            <Link to="/admin">管理后台</Link>
          </nav>
        </div>
      </header>
      <main style={{ padding: 32 }}>
        <Routes>
          <Route path="/" element={<PublicHomePage />} />
          <Route path="/console" element={<ConsoleHomePage />} />
          <Route path="/admin" element={<AdminHomePage />} />
        </Routes>
      </main>
    </div>
  );
}

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <BrowserRouter>
      <AppShell />
    </BrowserRouter>
  </React.StrictMode>
);
