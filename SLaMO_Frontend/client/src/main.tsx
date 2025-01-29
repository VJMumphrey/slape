import {StrictMode} from "react";
import {createRoot} from "react-dom/client";
import "./index.css";
import Prompt from "./Prompt.tsx";
import PromptInfo from "./Prompt.tsx";

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <Prompt />
  </StrictMode>
);
