import { useState } from "react";
import MenuTabs from "./MenuTabs.tsx";
import PipelineCard from "./PipelineCard.tsx";
import "./pipelines.css";

export default function Pipelines() {
  if (localStorage.getItem("StyleSetting") == null)
    localStorage.setItem("StyleSetting", "Dark");

  addEventListener("changedColorTheme", () => {
    setThemeColor(localStorage.getItem("StyleSetting"));
  });

  addEventListener("currentPipelineChange", () => {
    setCurrentPipeline(null);
  });

  addEventListener("deselectedPipeline", () => {
    setCurrentPipeline("");
  })

  const [ ThemeColor, setThemeColor ] = useState(localStorage.getItem("StyleSetting"));

  const [ CurrentPipeline, setCurrentPipeline ] = useState(null);


  const simplePipeline = {
    pipeline: "simple",
    displayName: "Simple Pipeline",
    description:
      "This pipeline sends a request to one model to get a simple I/O response. This is mainly meant for simple questions.",
  };

  const cotPipeline = {
    pipeline: "cot",
    displayName: "Chain of Thought",
    description: "This pipeline is designed to pass individual thoughts between a chain of models.",
  };

  const debatePipeline = {
    pipeline: "deb",
    displayName: "Debate Pipeline",
    description: "Utilizes multiple models in order to debate on a prompt.",
  };

  const embeddingPipeline = {
    pipeline: "emb",
    displayName: "Embedding Pipeline",
    description: "Utilized to embed prompts?",
  };

  async function savePipeline() {
    if (!localStorage.getItem("currentPipeline")) {
      if ((localStorage.getItem(`${CurrentPipeline}Models`) === null) || (JSON.parse(localStorage.getItem(`${CurrentPipeline}Models`) as string).length === 0)) {
        alert("Please Select Models for this Pipeline First!");
      } else {
        localStorage.setItem("currentPipeline", JSON.stringify(CurrentPipeline));
        dispatchEvent(new Event("currentPipelineChange"));
        const modelsObjects: {name: string, fullName: string}[] = JSON.parse(localStorage.getItem(`${CurrentPipeline}Models`) as string)
        const models: string[] = [];
        modelsObjects.forEach((element) => {
          models.push(element.fullName);
        })

        await fetch(`http://localhost:8080/${CurrentPipeline}/setup`, {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({
            models: models
          }),
        });
      }
    }
  }

  function determineSubmitStatus(): string {
    if (!localStorage.getItem("currentPipeline")) {
      return `${ThemeColor}_saveButton`;
    } else {
      return `${ThemeColor}_Unselectable_saveButton`
    }
  }

  return (
    <>
      <div className={`${ThemeColor}_background`}/>
      <MenuTabs />
      {/*To do this not stupid, add an onClick Property to PipelineCard that passes the onClick to the div of PipelineCard*/}
      <div className="cardContainer">
        <div onClick={() => {setCurrentPipeline(simplePipeline.pipeline)}}>
          <PipelineCard {...simplePipeline}></PipelineCard>
        </div>
        <div onClick={() => {setCurrentPipeline(cotPipeline.pipeline)}}>
          <PipelineCard {...cotPipeline}></PipelineCard>
        </div>
        <div onClick={() => {setCurrentPipeline(debatePipeline.pipeline)}}>
          <PipelineCard {...debatePipeline}></PipelineCard>
        </div>
        <div onClick={() => {setCurrentPipeline(embeddingPipeline.pipeline)}}>
          <PipelineCard {...embeddingPipeline}></PipelineCard>
        </div>
      </div>
      <div className={`${ThemeColor}_footer`}>
        <button className={determineSubmitStatus()} onClick={savePipeline}>
          Save Pipeline
        </button>
      </div>
    </>
  );
}
