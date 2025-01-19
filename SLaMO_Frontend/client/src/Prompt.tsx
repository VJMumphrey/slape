import "./App.css";
import Button from "./Button.tsx";

function Prompt() {
  return (
    <form>
      <label>
        Enter Prompt
        <input type="text" />
        <Button />
      </label>
    </form>
  );
}

export default Prompt;
