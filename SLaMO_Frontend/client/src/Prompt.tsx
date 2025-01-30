import {useState} from "react";
import "./App.css";
import {colors} from "../../../../../../AppData/Local/deno/npm/registry.npmjs.org/debug/4.3.7/src/browser.js";

export default function Prompt() {
  const [PromptInfo, setPromptInfo] = useState(""); //used to contain the current value, and to set the new value

  const handleSubmit = (event: {preventDefault: () => void}) => {
    event.preventDefault(); //makes sure the page doesn't reload when submitting the form
    alert(`the prompt you entered was: ${PromptInfo}`);
    setPromptInfo(""); //clears the prompt box after submission
  };

  return (
    <form onSubmit={handleSubmit}>
      <label>
        {" "}
        Enter Prompt:
        <input
          type="text"
          value={PromptInfo}
          onChange={(e) => setPromptInfo(e.target.value)} //access the current input and updates PromptInfo (e represents the event object)
        />
        <button style={{backgroundColor: "gray", color: "black"}}>
          {" "}
          Submit
        </button>
      </label>
    </form>
  );
}
