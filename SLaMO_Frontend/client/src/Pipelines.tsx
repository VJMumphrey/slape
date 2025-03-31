import { useState } from "react";
import MenuTabs from "./MenuTabs.tsx";
import PipelineCard from "./PipelineCard.tsx";
import "./pipelines.css";

export default function Pipelines() {

  const [ CurrentPipeline, setCurrentPipeline ] = useState(null);

  if (localStorage.getItem("PromptSetting") == null)
    localStorage.setItem("PromptSetting", "Automatic");
  if (localStorage.getItem("StyleSetting") == null)
    localStorage.setItem("StyleSetting", "Dark");

  const themeColor: string | null = localStorage.getItem("StyleSetting");

  const simplePipeline = {
    pipeline: "simple",
    displayName: "Simple Pipeline",
    description:
      "This pipeline sends a request to one model to get a simple I/O response. This is mainly meant for simple questions.",
  };

  const testModel2 = {
    pipeline: "TestPipelines",
    displayName: "TestPipelines",
    description: "Tests to see how multiple cards look",
  };

  async function savePipeline() {
    if (localStorage.getItem(`${CurrentPipeline}Models`) === null) {
      alert("Please Select Models for this Pipeline First!");
    } else {
      localStorage.setItem("currentPipeline", JSON.stringify(CurrentPipeline));
      const modelsObjects: {name: string, fullName: string}[] = JSON.parse(localStorage.getItem(`${CurrentPipeline}Models`) as string)
      const models: string[] = [];
      modelsObjects.forEach((element) => {
        models.push(element.fullName);
      })

      const response = await fetch(`http://localhost:8080/${CurrentPipeline}/setup`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          models: models
        }),
      });

      alert(response);
    }
  }

  return (
    <>
      <div className={`${themeColor}_background`}/>
      <MenuTabs />
      {/*To do this not stupid, add an onClick Property to PipelineCard that passes the onClick to the div of PipelineCard*/}
      <div onClick={() => {setCurrentPipeline(simplePipeline.pipeline)}}>
        <PipelineCard {...simplePipeline}></PipelineCard>
      </div>
      <p></p>
      <div onClick={() => {setCurrentPipeline(testModel2.pipeline)}}>
        <PipelineCard {...testModel2}></PipelineCard>
      </div>
      <div className={`${themeColor}_footer`}>
        <button className={`${themeColor}_saveButton`} onClick={savePipeline}>
          Save Pipeline
        </button>
      </div>
    </>
  );
}
