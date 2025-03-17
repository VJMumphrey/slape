import MenuTabs from "./MenuTabs.tsx";
import PipelineCard from "./PipelineCard.tsx";

export default function Pipelines() {

  if (localStorage.getItem("PromptSetting") == null)
    localStorage.setItem("PromptSetting", "Automatic");
  if (localStorage.getItem("StyleSetting") == null)
    localStorage.setItem("StyleSetting", "Dark");

  const themeColor: string | null = localStorage.getItem("StyleSetting");

  const testModel = {
    pipeline: "TestPipeline",
    models: ["Phi-3.5", "Dolphin", "Ms. Gorman"],
    description:
      "This test Pipeline tests how well I could pipeline your mom. So far the tests are successful.",
  };

  const testModel2 = {
    pipeline: "TestPipelines",
    models: ["Phi-3.5", "Dolphin"],
    description: "Tests to see how multiple cards look",
  };

  return (
    <>
      <div className={`${themeColor}_background`}/>
      <MenuTabs />
      <PipelineCard {...testModel}></PipelineCard>
      <p></p>
      <PipelineCard {...testModel2}></PipelineCard>
    </>
  );
}
