import {useState} from "react";
import "./App.css";
import Button from "./Button.tsx";

function Prompt() {
  const [PromptInfo, setPromptInfo] = useState("");

  const handleSubmit = (event) => {
    event.preventDefault();
    alert(`the prompt you entered was: ${PromptInfo}`);
  };
  return (
    <form onSubmit={handleSubmit}>
      <label>
        <input
          type="text"
          value={PromptInfo}
          onChange={(e) => setPromptInfo(e.target.value)}
        />
        Enter Prompt
        <Button />
      </label>
    </form>
  );
}

export default Prompt;
