import {useState} from "react";
import "./App.css";

export default function Prompt() {
  const [PromptInfo, setPromptInfo] = useState(""); //used to contain the current value, and to set the new value

  const handleSubmit = (event: {preventDefault: () => void}) => {
    event.preventDefault(); //makes sure the page doesn't reload when submitting the form
    alert(`the prompt you entered was: ${PromptInfo}`);
    setPromptInfo(""); //clears the prompt box after submission
  };

  function Button() {
    function handleClick() {}
    return <button onClick={handleClick}>Submit</button>;
  }

  return (
    <form onSubmit={handleSubmit}>
      <label>
        <input
          type="text"
          value={PromptInfo}
          onChange={(e) => setPromptInfo(e.target.value)} //access the current input and updates PromptInfo (e represnts the event object)
        />
        Enter Prompt
        <Button />
      </label>
    </form>
  );
}
