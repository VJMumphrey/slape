import {useState, useRef, useEffect } from "react";
import "./App.css";
import "./prompt.css";
import MenuTabs from "./MenuTabs.tsx";
import Markdown from "react-markdown";
import DropDownButton from "./DropDownButton.tsx";
import { allGeneratedPositionsFor } from "../../../../../AppData/Local/deno/npm/registry.npmjs.org/@jridgewell/trace-mapping/0.3.25/dist/types/trace-mapping.d.ts";

export default function Prompt() {
  const [PromptInfo, setPromptInfo] = useState(""); //used to contain the current value, and to set the new value
  const [ResponseAnswer, setResponseAnswer] = useState("... Awaiting Response");
  const [PromptMode, setPromptMode] = useState("simple");
  const [ ThinkingMode, setThinkingMode ] = useState(false);
  const [ InternetSearchMode, setInternetSearchMode ] = useState(false);
  const [loadingAnimation, setloadingAnimation] = useState("");
  const [ShiftPressed, setShiftPressed] = useState(false);

  if (localStorage.getItem("PromptSetting") == null)
    localStorage.setItem("PromptSetting", "Automatic");
  if (localStorage.getItem("StyleSetting") == null)
    localStorage.setItem("StyleSetting", "Dark");

  const themeColor: string | null = localStorage.getItem("StyleSetting");

  const textAreaRef = useRef<HTMLTextAreaElement>(null);

  const promptTypes = [
    {name: "Simple", type: "simple"},
    {name: "Chain of Thought", type: "cot"},
    {name: "Tree of Thought", type: "tot"},
    {name: "Graph of Thought", type: "got"},
    {name: "Thinking Hats", type: "thinkinghats"},
    {name: "Mixture of Experts", type: "moe"},
  ];

  function handleKeyDown(e: React.KeyboardEvent<HTMLTextAreaElement>) {
    if (e.key === "Shift") {
      setShiftPressed(true);
    } else if (ShiftPressed && e.key === "Enter") {
      e.preventDefault();
      setPromptInfo(PromptInfo + "\n");
    } else if (!ShiftPressed && e.key === "Enter") {
      e.preventDefault();
      handleSubmit(e);
    }
  }

  function keyUpMapHandler(e: React.KeyboardEvent<HTMLTextAreaElement>) {
    if (e.key === "Shift") {
      setShiftPressed(false);
    }
  }

  const useAutosizeTextArea = (
    textAreaRef: HTMLTextAreaElement | null,
    value: string
  ) => {
    useEffect(() => {
      if (textAreaRef) {
        // We need to reset the height momentarily to get the correct scrollHeight for the textarea
        textAreaRef.style.height = "0px";
        const scrollHeight = textAreaRef.scrollHeight;
  
        // We then set the height directly, outside of the render loop
        // Trying to set this with state or a ref will product an incorrect value.
        textAreaRef.style.height = scrollHeight + "px";
      }
    }, [textAreaRef, value]);
  };

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
    if (localStorage.getItem("currentPipeline") !== null) {
      const response = await fetch(`http://localhost:8080/${JSON.parse(localStorage.getItem("currentPipeline") as string)}/generate`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          prompt: PromptInfo,
          thinking: String(ThinkingMode),
          search: String(InternetSearchMode),
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
    } else {
      alert("Please Select a Pipeline from the Pipelines Tab");
    }
  }

  useAutosizeTextArea(textAreaRef.current, PromptInfo);
  return (
    <>
      <div className={`${themeColor}_background`} />
      <MenuTabs />
      <div className={`${themeColor}_output`}>{ResponseAnswer}</div>
      <div className={`${themeColor}_fixedBottom`}>
          <textarea
            className={`${themeColor}_prompt`}
            value={PromptInfo}
            placeholder="Enter Prompt"
            onChange={(e) => setPromptInfo(e.target.value)} //access the current input and updates PromptInfo (e represents the event object)
            rows={1}
            ref={textAreaRef}
            onKeyDown={(e) => {handleKeyDown(e);}}
            onKeyUp={(e) => {keyUpMapHandler(e);}}
          />
          <DropDownButton
            className="inference"
            value={PromptMode}
            callBack={setPromptMode}
            optionObject={promptTypes}
          />
          <label>
            <input
              className={`${themeColor}_thinkingBox`}
              type="checkbox"
              checked={ThinkingMode}
              onChange={() => {setThinkingMode(!ThinkingMode)}}
            />
            <span>Thinking</span>
          </label>
          <label>
            <input
              className={`${themeColor}_internetSearchBox`}
              type="checkbox"
              checked={InternetSearchMode}
              onChange={() => {setInternetSearchMode(!InternetSearchMode)}}
            />
            <span>Internet Search</span>
          </label>
          <button className={`${themeColor}_promptSubmit`}onClick={handleSubmit}> Submit</button>
      </div>
      <p className="loading">{loadingAnimation}</p>
    </>
  );
}
