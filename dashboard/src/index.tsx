import "core-js/stable";
import "regenerator-runtime/runtime";

import { createRoot } from "react-dom/client";
import React from "react";
import App from "App";
const container = document.getElementById("output");
const root = createRoot(container!);

root.render(<App />);
