import {StrictMode} from "react";
import {createRoot} from "react-dom/client";
import "./index.css";
import Prompt from "./Prompt.tsx";

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <div className="output"> Here's where the output will go</div>
    <div className="fixedBottom">
      <Prompt />
    </div>
  </StrictMode>
);
