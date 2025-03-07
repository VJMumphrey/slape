import "./pipelineCard.css"

interface pipelineProperties{
    pipeline: string,
    models: string[],
    description: string
}

export default function PipelineCard({ pipeline, models, description }: pipelineProperties) {

    let modelsString = ""
    models.forEach((element, index) => {
        if (index != models.length - 1)
            modelsString += `${element}, `;
        else
            modelsString += element
        }
    );

    return (
        <div className="pipelineDiv">
            <h3>{pipeline}</h3>
            <button>Settings</button>
            <p>{`Models: ${modelsString}`}</p>
            <p>{`Description: ${description}`}</p>
        </div>
    );
}