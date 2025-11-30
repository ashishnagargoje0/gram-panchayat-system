import React from "react";
import { createRoot } from "react-dom/client";
import App from "./App";
import "./index.css"; // ensure Tailwind / styles are loaded

const root = createRoot(document.getElementById("root"));
root.render(<App />);
