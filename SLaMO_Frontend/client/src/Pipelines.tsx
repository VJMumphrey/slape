import MenuTabs from "./MenuTabs.tsx";
import PipelineCard from "./PipelineCard.tsx";

export default function Pipelines() {

    const testModel = {
        pipeline: "TestPipeline",
        models: [ "Phi-3.5", "Dolphin", "Ms. Gorman" ],
        description: "This test Pipeline tests how well I could pipeline your mom. So far the tests are successful.",
    }

    return(
        <>
            <MenuTabs/>
            <PipelineCard {...testModel}></PipelineCard>
        </>
    );
}