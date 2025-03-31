import "./pipelineCard.css";
import {useState} from "react";
import Modal from "./Modal.tsx";
import {createPortal} from "react-dom";
import DropDownButton from "./DropDownButton.tsx";
import {ReactNode} from "react";

interface pipelineProperties {
  pipeline: string;
  displayName: string;
  description: string;
}

export default function PipelineCard({
  pipeline,
  displayName,
  description,
}: pipelineProperties) {
  const [ModalOpen, setModalOpen] = useState(false);
  const [modelName, setmodelName] = useState("Phi 3");
  const [Models, setModels] = useState<object[]>([]);
  const [ AddModelsButtonStatus, setAddModelsButtonStatus ] = useState("addModelActive");

  if (localStorage.getItem("PromptSetting") == null)
    localStorage.setItem("PromptSetting", "Automatic");
  if (localStorage.getItem("StyleSetting") == null)
    localStorage.setItem("StyleSetting", "Dark");

  const themeColor: string | null = localStorage.getItem("StyleSetting");

  const pipelineSettingsButtonHandler = () => {
    setModalOpen(true);
    setModels([]);
    localStorage.removeItem(`${pipeline}Models`);
    setAddModelsButtonStatus("addModelActive");
  };

  const modalCloseButtonHandler = () => {
    setModalOpen(false);
    localStorage.setItem(`${pipeline}Models`, JSON.stringify(Models));
  };

  const modelNamesDropDownOptions = [
    {type: "Phi 3", name: "Phi 3"},
    {type: "Dolphin", name: "Dolphin"},
  ];

  // function addModelHandler() {
  //   if (localStorage.getItem(`${pipeline}Models`) == null)
  //     localStorage.setItem(`${pipeline}Models`, JSON.stringify([modelName]));
  //   else {
  //     const currentPipelineModel: string[] = JSON.parse(localStorage.getItem(`${pipeline}Models`) as string);
  //     currentPipelineModel.push(modelName);
  //     localStorage.setItem(
  //       `${pipeline}Models`,
  //       JSON.stringify(currentPipelineModel)
  //     );
  //   }
  // }

  function addModelHandler() {
    let modelObject = {name: modelName, fullName: ""};
    switch (modelName) {
      case "Phi 3": {
        modelObject = {
          ...modelObject,
          fullName: "Phi-3.5-mini-instruct.Q4_K_M.gguf",
        };
        break;
      }

      case "Dolphin": {
        modelObject = {
          ...modelObject,
          fullName: "Dolphin3.0-Llama3.2-1B-Q4_K_M.gguf",
        };
        break;
      }

      default: {
        modelObject = {...modelObject, fullName: "kys"};
      }
    }

    if (pipeline !== "simple") {
      setModels([...Models, modelObject]);
    } else if (Models.length === 0) {
      setModels([modelObject]);
      setAddModelsButtonStatus("addModelInactive");
    }
  }

  function displayCurrentModels(className: string): ReactNode {
    let modelsString = "";
    Models.forEach(
      (element: {name: string; fullName: string}, index: number) => {
        if (index !== Models.length - 1) {
          modelsString += element.name + ", ";
        } else {
          modelsString += element.name;
        }
      }
    );
    return <p className={className}>{`Current Models: ${modelsString}`}</p>;
  }

  return (
    <div className={`${themeColor}_pipelineDiv`}>
      <h3 className="pipelineHeader">{displayName}</h3>
      <button
        className={`${themeColor}_pipelineButton`}
        onClick={pipelineSettingsButtonHandler}
      >
        Settings
      </button>
      {createPortal(
        <Modal isOpen={ModalOpen} onClose={modalCloseButtonHandler}>
          <p>
            {/* {displayCurrentModels("")} */}
            <button
              className={`${themeColor}_${AddModelsButtonStatus}`}
              onClick={addModelHandler}
            >
              Add Model
            </button>
            <DropDownButton
              className="modelDropdown"
              value={modelName}
              callBack={setmodelName}
              optionObject={modelNamesDropDownOptions}
            />
          </p>
        </Modal>,
        document.body
      )}
      {displayCurrentModels("pipelineModels")}
      <p className="pipelineDesc">{`Description: ${description}`}</p>
    </div>
  );
}
