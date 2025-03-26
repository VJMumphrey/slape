import "./pipelineCard.css";
import {useState} from "react";
import Modal from "./Modal.tsx";
import {createPortal} from "react-dom";
import DropDownButton from "./DropDownButton.tsx";

interface pipelineProperties {
  pipeline: string;
  models: string[];
  description: string;
}

export default function PipelineCard({
  pipeline,
  models,
  description,
}: pipelineProperties) {
  const [ModalOpen, setModalOpen] = useState(false);

  let modelsString = "";
  models.forEach((element, index) => {
    if (index != models.length - 1) modelsString += `${element}, `;
    else modelsString += element;
  });

  const pipelineSettingsButtonHandler = () => {
    setModalOpen(true);
    localStorage.removeItem("models");
  };

  const modalCloseButtonHandler = () => {
    setModalOpen(false);
    if (localStorage.getItem("models") == null)
      localStorage.setItem("models", "Phi 3");
  };

  const [modelName, setmodelName] = useState("Phi 3");

  const modelNamesDropDownOptions = [
    {type: "Phi 3", name: "Phi 3"},
    {type: "Dolphin", name: "Dolphin"},
  ];

  function addModelHandler() {
    if (localStorage.getItem("models") == null)
      localStorage.setItem("models", modelName);
    else
      localStorage.setItem(
        "models",
        localStorage.getItem("models") + "," + modelName
      );
  }

  return (
    <div className="pipelineDiv">
      <h3 className="pipelineHeader">{pipeline}</h3>
      <button
        className="pipelineButton"
        onClick={pipelineSettingsButtonHandler}
      >
        Settings
      </button>
      {createPortal(
        <Modal isOpen={ModalOpen} onClose={modalCloseButtonHandler}>
          <p>
            <DropDownButton
              value={modelName}
              callBack={setmodelName}
              optionObject={modelNamesDropDownOptions}
            />
            <button onClick={addModelHandler}>Add Model</button>
          </p>
        </Modal>,
        document.body
      )}
      <p className="pipelineModels">{`Models: ${modelsString}`}</p>
      <p className="pipelineDesc">{`Description: ${description}`}</p>
    </div>
  );
}
