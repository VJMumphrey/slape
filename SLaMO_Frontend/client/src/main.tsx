import {StrictMode} from "react";
import {createRoot} from "react-dom/client";
import "./index.css";
import App from "./Prompt.tsx";
import Prompt from "./Prompt.tsx";
import Button from "./Button.tsx";

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <Prompt />
  </StrictMode>
);
