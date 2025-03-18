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
  const [modelName, setmodelName] = useState("Phi 3");

  if (localStorage.getItem("PromptSetting") == null)
    localStorage.setItem("PromptSetting", "Automatic");
  if (localStorage.getItem("StyleSetting") == null)
    localStorage.setItem("StyleSetting", "Dark");

  const colorTheme: string | null = localStorage.getItem("StyleSetting");

  let modelsString = "";
  models.forEach((element, index) => {
    if (index != models.length - 1) modelsString += `${element}, `;
    else modelsString += element;
  });

  const pipelineSettingsButtonHandler = () => {
    setModalOpen(true);
    localStorage.removeItem(`${pipeline}Models`);
  };

  const modalCloseButtonHandler = () => {
    setModalOpen(false);
  };


  const modelNamesDropDownOptions = [
    {type: "Phi 3", name: "Phi 3"},
    {type: "Dolphin", name: "Dolphin"},
  ];

  function addModelHandler() {
    if (localStorage.getItem(`${pipeline}Models`) == null)
      localStorage.setItem(`${pipeline}Models`, JSON.stringify([modelName]));
    else {
      const currentPipelineModel: string[] = JSON.parse(localStorage.getItem(`${pipeline}Models`) as string);
      currentPipelineModel.push(modelName);
      localStorage.setItem(
        `${pipeline}Models`,
        JSON.stringify(currentPipelineModel)
      );
    }
  }

  return (
    <div className={`${colorTheme}_pipelineDiv`}>
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
