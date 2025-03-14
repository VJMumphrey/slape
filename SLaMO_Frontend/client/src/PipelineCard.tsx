import "./pipelineCard.css";
import {useState} from 'react';
import Modal from "./Modal.tsx";
import {createPortal} from "react-dom";

interface pipelineProperties{
    pipeline: string,
    models: string[],
    description: string
}

export default function PipelineCard({ pipeline, models, description }: pipelineProperties) {

    const [ ModalOpen, setModalOpen ] = useState(false);

    let modelsString = ""
    models.forEach((element, index) => {
        if (index != models.length - 1)
            modelsString += `${element}, `;
        else
            modelsString += element
        }
    );

    const pipelineSettingsButtonHandler = () => {
        setModalOpen(true);
    };

    const modalCloseButtonHandler = () => {
        setModalOpen(false);
    };

    return (
        <div className="pipelineDiv">
            <h3 className="pipelineHeader">{pipeline}</h3>
            <button className="pipelineButton" onClick={pipelineSettingsButtonHandler}>Settings</button>
            {createPortal(
                <Modal isOpen={ModalOpen} onClose={modalCloseButtonHandler}>
                    <h1>Stone A Bitch?</h1>
                </Modal>,
                document.body
            )}
            <p className="pipelineModels">{`Models: ${modelsString}`}</p>
            <p className="pipelineDesc">{`Description: ${description}`}</p>
        </div>
    );
}