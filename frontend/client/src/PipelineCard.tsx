import "./pipelineCard.css";
import {useState, ReactNode} from "react";
import Modal from "./Modal.tsx";
import {createPortal} from "react-dom";
import DropDownButton from "./DropDownButton.tsx";

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

  if (localStorage.getItem("StyleSetting") == null)
    localStorage.setItem("StyleSetting", "Dark");

  addEventListener("currentPipelineChange", () => {
    if (pipeline === JSON.parse(localStorage.getItem("currentPipeline") as string)) {
      setCurrentlySelected(true);
      setPipelineStatus(": Running");
    } else {
      setPipelineStatus("  ");
    }
  });

  addEventListener("deselectedPipeline", () => {
    setPipelineStatus(" ");
  });

  addEventListener("changedColorTheme", () => {
    setThemeColor(localStorage.getItem("StyleSetting"));
  });

  const [ ThemeColor, setThemeColor ] = useState(localStorage.getItem("StyleSetting"));

  const [ModalOpen, setModalOpen] = useState(false);
  const [modelName, setmodelName] = useState("Phi 3");
  const [ModelList, setModelList] = useState([]);
  const [Models, setModels] = useState<object[]>(localStorage.getItem(`${pipeline}Models`) ? JSON.parse(localStorage.getItem(`${pipeline}Models`) as string) : []);
  const [AddModelsButtonStatus, setAddModelsButtonStatus] =
    useState("addModelActive");
  const [ CurrentlySelected, setCurrentlySelected ] = useState(pipeline === JSON.parse(localStorage.getItem("currentPipeline") as string));
  const [ PipelineStatus, setPipelineStatus ] = useState(JSON.parse(localStorage.getItem("currentPipeline") as string) === pipeline ? ": Running" : "");


  const pipelineSettingsButtonHandler = () => {
    setModalOpen(true);
    getModelsFromVito();
    setModels([]);
    localStorage.removeItem(`${pipeline}Models`);
    setAddModelsButtonStatus("addModelActive");
  };

  const modalCloseButtonHandler = () => {
    setModalOpen(false);
    localStorage.setItem(`${pipeline}Models`, JSON.stringify(Models));
  };

  async function getModelsFromVito() {
    const dropDownOptions: {type: string; name: string}[] = [];
    const response = await fetch(`http://localhost:8080/getmodels`, {
      method: "GET",
    });

    let responseBody: {models: string[]};
    let models: string[] = [];

    if (response.ok) {
      responseBody = await response.json();
      models = responseBody.models;
    } else {
      alert(
        "Failed to pull models. Please restart front end and ensure backend is running."
      );
    }

    models.forEach((element) => {
      dropDownOptions.push({type: element, name: element.slice(0, -5)});
    });

    setModelList(dropDownOptions);
    setmodelName(dropDownOptions[0].type);
  }
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
    console.log(modelName);
    const modelObject = {name: modelName.slice(0, -5), fullName: modelName};

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

  async function shutdownCurrentPipeline() {
    setPipelineStatus(": Attempting Shutdown, Please Wait.");
    const response = await fetch(`http://localhost:8080/${pipeline}/shutdown`, {method: "GET",});
    if (response.ok) {
      localStorage.removeItem("currentPipeline");
      setCurrentlySelected(false);
      setPipelineStatus("");
      globalThis.dispatchEvent(new Event("deselectedPipeline"));
    } else {
      alert("Pipeline Failed to Shutdown. Please Try Again.");
      setPipelineStatus(": Running");
    }
  }

  function determinePipelineClass(): string {
    if (CurrentlySelected) {
      return `${ThemeColor}_Selected_pipelineDiv`;
    } else if (localStorage.getItem("currentPipeline")) {
      return `${ThemeColor}_notSelected_pipelineDiv`;
    } else {
      return `${ThemeColor}_pipelineDiv`;
    }
  }

  function determineCardButton(): ReactNode {
    if (CurrentlySelected) {
      return <button onClick={shutdownCurrentPipeline} className={`${ThemeColor}_shutdownPipelineButton`}>&times;</button>
    } else if (localStorage.getItem("currentPipeline")) {
      return <button className={`${ThemeColor}_notSelected_button`}>Select Models</button>
    } else {
      return <button className={`${ThemeColor}_pipelineButton`} onClick={pipelineSettingsButtonHandler}>Select Models</button>
    }
  }

  return (
    <div className={determinePipelineClass()} tabIndex={0}>
      <h3 className="pipelineHeader">{displayName + PipelineStatus}</h3>
      {determineCardButton()}
      {createPortal(
        <Modal isOpen={ModalOpen} onClose={modalCloseButtonHandler}>
          <p>
            {/* {displayCurrentModels("")} */}
            <button
              className={`${ThemeColor}_${AddModelsButtonStatus}`}
              onClick={addModelHandler}
            >
              Add Model
            </button>
            <DropDownButton
              className="modelDropdown"
              value={modelName}
              callBack={setmodelName}
              optionObject={ModelList}
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
