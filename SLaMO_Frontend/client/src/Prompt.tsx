import {useState} from "react";
import "./App.css";
import MenuTabs from "./MenuTabs.tsx";
import Markdown from "react-markdown";

export default function Prompt() {
  const [PromptInfo, setPromptInfo] = useState(""); //used to contain the current value, and to set the new value
  const [ResponseAnswer, setResponseAnswer] = useState("... Awaiting Response");
  const [PromptMode, setPromptMode] = useState("simple");

  const promptTypes = [
    {name: "Simple", type: "simple"},
    {name: "Chain of Thought", type: "cot"},
    {name: "Tree of Thought", type: "tot"},
    {name: "Graph of Thought", type: "got"},
    {name: "Thinking Hats", type: "thinkinghats"},
    {name: "Mixture of Experts", type: "moe"},
  ];

  function dropDownButton() {
    return (
      <select
        className="inference"
        value={PromptMode}
        onChange={(event) => {
          setPromptMode(event.target.value);
        }}
      >
        {promptTypes.map((type) => {
          return (
            <option value={type.type} key={type.type}>
              {type.name}
            </option>
          );
        })}
      </select>
    );
  }

  async function handleSubmit(event: {preventDefault: () => void}) {
    // event.preventDefault();
    // const model = await fetch("http://localhost:3069/up");

    // if (model.status !== 200) {
    //   alert("model not ready");
    // } else {
    setPromptInfo(""); //clears the prompt box after submission
    setResponseAnswer(
      <>
        <p className="left">{`Prompt: ${PromptInfo}`}</p>{" "}
        <p className="left">{`Response: Generating Response...`}</p>
      </>
    );
    event.preventDefault(); //makes sure the page doesn't reload when submitting the form
    const response = await fetch("http://localhost:3069/simple", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        prompt: PromptInfo,
        model: "gay",
        mode: PromptMode,
      }),
    });
    setResponseAnswer(
      <>
        <p className="left">{`Prompt: ${PromptInfo}`}</p>{" "}
        <p className="left">
          Response:
          <Markdown>
            {`${JSON.parse(JSON.stringify(await response.json())).answer}`}
          </Markdown>
        </p>
      </>
    );
    // }
  }
  return (
    <>
      <MenuTabs />
      <div className="output">{ResponseAnswer}</div>
      <div className="fixedBottom">
        <form onSubmit={handleSubmit}>
          <label>
            {" "}
            Enter Prompt:
            <input
              className="prompt"
              type="text"
              value={PromptInfo}
              onChange={(e) => setPromptInfo(e.target.value)} //access the current input and updates PromptInfo (e represents the event object)
            />
            <button style={{backgroundColor: "gray", color: "black"}}>
              {" "}
              Submit
            </button>
            {dropDownButton()}
          </label>
        </form>
      </div>
    </>
  );
}
