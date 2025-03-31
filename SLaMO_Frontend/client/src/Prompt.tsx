import {useState} from "react";
import "./App.css";
import "./prompt.css";
import MenuTabs from "./MenuTabs.tsx";
import Markdown from "react-markdown";
import DropDownButton from "./DropDownButton.tsx";

export default function Prompt() {
  const [PromptInfo, setPromptInfo] = useState(""); //used to contain the current value, and to set the new value
  const [ResponseAnswer, setResponseAnswer] = useState("... Awaiting Response");
  const [PromptMode, setPromptMode] = useState("simple");
  const [loadingAnimation, setloadingAnimation] = useState("");

  if (localStorage.getItem("PromptSetting") == null)
    localStorage.setItem("PromptSetting", "Automatic");
  if (localStorage.getItem("StyleSetting") == null)
    localStorage.setItem("StyleSetting", "Dark");

  const themeColor: string | null = localStorage.getItem("StyleSetting");

  const promptTypes = [
    {name: "Simple", type: "simple"},
    {name: "Chain of Thought", type: "cot"},
    {name: "Tree of Thought", type: "tot"},
    {name: "Graph of Thought", type: "got"},
    {name: "Thinking Hats", type: "thinkinghats"},
    {name: "Mixture of Experts", type: "moe"},
  ];

  async function handleSubmit(event: {preventDefault: () => void}) {
    setPromptInfo(""); //clears the prompt box after submission
    setloadingAnimation(<div className={`${themeColor}_spinner`} />);
    setResponseAnswer(
      <>
        <p className={`${themeColor}_left`}>{`Prompt: ${PromptInfo}`}</p>{" "}
        <p
          className={`${themeColor}_left`}
        >{`Response: Generating Response...`}</p>
      </>
    );
    event.preventDefault(); //makes sure the page doesn't reload when submitting the form
    const response = await fetch("http://localhost:8080/simple", {
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
    setloadingAnimation("");
    setResponseAnswer(
      <>
        <p className={`${themeColor}_left`}>{`Prompt: ${PromptInfo}`}</p>{" "}
        <p className={`${themeColor}_left`}>
          Response:
          <Markdown>
            {`${JSON.parse(JSON.stringify(await response.json())).answer}`}
          </Markdown>
        </p>
      </>
    );
  }
  return (
    <>
      <div className={`${themeColor}_background`} />
      <MenuTabs />
      <div className={`${themeColor}_output`}>{ResponseAnswer}</div>
      <div className={`${themeColor}_fixedBottom`}>
        <form onSubmit={handleSubmit}>
          <label>
            {" "}
            <input
              className={`${themeColor}_prompt`}
              type="text"
              value={PromptInfo}
              placeholder="Enter Prompt"
              onChange={(e) => setPromptInfo(e.target.value)} //access the current input and updates PromptInfo (e represents the event object)
            />
            <DropDownButton
              className="inference"
              value={PromptMode}
              callBack={setPromptMode}
              optionObject={promptTypes}
            />
            <button className={`${themeColor}_promptSubmit`}> Submit</button>
          </label>
        </form>
      </div>
      <p className="loading">{loadingAnimation}</p>
    </>
  );
}
