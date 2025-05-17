import {useState, useRef, useEffect} from "react";
import "./App.css";
import "./prompt.css";
import MenuTabs from "./MenuTabs.tsx";
import Markdown from "react-markdown";
import DropDownButton from "./DropDownButton.tsx";

export default function Prompt() {
  if (localStorage.getItem("StyleSetting") == null)
    localStorage.setItem("StyleSetting", "Dark");

  addEventListener("changedColorTheme", () => {
    setThemeColor(localStorage.getItem("StyleSetting"));
  });

  const [ ThemeColor, setThemeColor ] = useState(localStorage.getItem("StyleSetting"));
  const [PromptInfo, setPromptInfo] = useState(""); //used to contain the current value, and to set the new value
  const [PromptMode, setPromptMode] = useState("simple");
  const [ThinkingMode, setThinkingMode] = useState(false);
  const [InternetSearchMode, setInternetSearchMode] = useState(false);
  const [loadingAnimation, setloadingAnimation] = useState("");
  const [ShiftPressed, setShiftPressed] = useState(false);
  const [PipelineSelected, setPipelineSelected] = useState(Boolean(localStorage.getItem("currentPipeline")));
  const [ResponseAnswer, setResponseAnswer] = useState(PipelineSelected ? "...Awaiting Prompt" : "Please Select a Pipeline to begin Prompting!");
  

  const textAreaRef = useRef<HTMLTextAreaElement>(null);
  const footerDivRef = useRef<HTMLDivElement>(null);
  const responseDivRef = useRef<HTMLDivElement>(null);

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
    footerDivRef: HTMLDivElement | null,
    value: string
  ) => {
    useEffect(() => {
      if (textAreaRef && footerDivRef) {
        const defaultDivHeight = 100;
        const defaultTextHeight = 26;
        if (defaultDivHeight + textAreaRef.scrollHeight <= 310) {
          // We need to reset the height momentarily to get the correct scrollHeight for the textarea
          textAreaRef.style.height = "0px";
          const scrollHeight = textAreaRef.scrollHeight;

          // We then set the height directly, outside of the render loop
          // Trying to set this with state or a ref will product an incorrect value.
          textAreaRef.style.height = scrollHeight + "px";
          if (textAreaRef.scrollHeight >= defaultTextHeight)
            footerDivRef.style.height = String((defaultDivHeight - defaultTextHeight) + Number(scrollHeight)) + "px";
          else
            footerDivRef.style.height = defaultDivHeight + "px";
          // if (value === "") {
          //   textAreaRef.style.height = defaultTextHeight + "px";
          //   footerDivRef.style.height = defaultDivHeight + "px"
          // }
        }
      }
    }, [textAreaRef, footerDivRef, value]);
  };

  async function handleSubmit(event: {preventDefault: () => void}) {
      setPromptInfo(""); //clears the prompt box after submission
      setloadingAnimation(<div className={`${ThemeColor}_spinner`} />);
      setResponseAnswer(
        <>
          <p className={`${ThemeColor}_leftTitle`}>Prompt:</p>
          <p className={`${ThemeColor}_left`}>{`${PromptInfo}`}</p>{" "}
          <p className={`${ThemeColor}_leftTitle`}>Response:</p>
          <p className={`${ThemeColor}_left`}>{`Generating Response...`}</p>
        </>
      );
      event.preventDefault(); //makes sure the page doesn't reload when submitting the form
      if (localStorage.getItem("currentPipeline") !== null) {
        const response = await fetch(
          `http://localhost:8080/${JSON.parse(
            localStorage.getItem("currentPipeline") as string
          )}/generate`,
          {
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
          }
        );
        setloadingAnimation("");
        setResponseAnswer(
          <>
            <p className={`${ThemeColor}_leftTitle`}>Prompt:</p>
            <p className={`${ThemeColor}_left`}>{`${PromptInfo}`}</p>{" "}
            <p className={`${ThemeColor}_leftTitle`}>Response:</p>
            <p className={`${ThemeColor}_left`}>
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

  function determinePipelineStatus(): string {
    if (PipelineSelected && ResponseAnswer === "") {
      return "...Awaiting Prompt";
    } else if (PipelineSelected) {
      return ResponseAnswer;
    } else {
      return "Please Select a Pipeline to begin Prompting!";
    }
  }

  useAutosizeTextArea(textAreaRef.current, footerDivRef.current, PromptInfo);
  return (
    <>
      <div className={`${ThemeColor}_background`} />
      <MenuTabs />
      <div className={`${ThemeColor}_output`} ref={responseDivRef}>{determinePipelineStatus()}</div>
      <div className={`${ThemeColor}_fixedBottom`} ref={footerDivRef}>
        <textarea
          className={`${ThemeColor}_prompt`}
          value={PromptInfo}
          placeholder="Enter Prompt"
          onChange={(e) => setPromptInfo(e.target.value)} //access the current input and updates PromptInfo (e represents the event object)
          rows={1}
          ref={textAreaRef}
          onKeyDown={(e) => {
            handleKeyDown(e);
          }}
          onKeyUp={(e) => {
            keyUpMapHandler(e);
          }}
        />
        <DropDownButton
          className="inference"
          value={PromptMode}
          callBack={setPromptMode}
          optionObject={promptTypes}
        />
        <label>
          <input
            className={`${ThemeColor}_thinkingBox`}
            type="checkbox"
            checked={ThinkingMode}
            onChange={() => {
              setThinkingMode(!ThinkingMode);
            }}
          />
          <span>Thinking</span>
        </label>
        <label>
          <input
            className={`${ThemeColor}_internetSearchBox`}
            type="checkbox"
            checked={InternetSearchMode}
            onChange={() => {
              setInternetSearchMode(!InternetSearchMode);
            }}
          />
          <span>Internet Search</span>
        </label>
        <button className={`${ThemeColor}_promptSubmit`} onClick={(e) => handleSubmit(e)}>
          {" "}
          Submit
        </button>
      </div>
      <div className="loading">{loadingAnimation}</div>
    </>
  );
}
