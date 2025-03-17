import DropDownButton from "./DropDownButton.tsx";
import MenuTabs from "./MenuTabs.tsx";
import {useState} from "react";

export default function Settings() {
  if (localStorage.getItem("PromptSetting") == null)
    localStorage.setItem("PromptSetting", "Automatic");
  if (localStorage.getItem("StyleSetting") == null)
    localStorage.setItem("StyleSetting", "Dark");

  const [PromptSetting, setPromptSetting] = useState(
    localStorage.getItem("PromptSetting")
  );
  const [StyleSetting, setStyleSetting] = useState(
    localStorage.getItem("StyleSetting")
  );

  const promptDropDownOptions = [
    {type: "Automatic", name: "Automatic"},
    {type: "Manual", name: "Manual"},
    {type: "Mixed", name: "Mixed"},
  ];
  const styleDropDownOptions = [
    {type: "Pink", name: "Pink"},
    {type: "Dark", name: "Dark"},
    {type: "Light", name: "Light"},
  ];

  function settingsButtonHandler() {
    localStorage.setItem("PromptSetting", PromptSetting);
    localStorage.setItem("StyleSetting", StyleSetting);
    globalThis.location.reload();
  }

  return (
    <>
      <MenuTabs />
      <div>
        <h2>Settings</h2>
        <h4>Prompting Settings</h4>
        <DropDownButton
          value={PromptSetting}
          callBack={setPromptSetting}
          optionObject={promptDropDownOptions}
        />
        <hr />
        <h4>Style Settings</h4>
        <DropDownButton
          value={StyleSetting}
          callBack={setStyleSetting}
          optionObject={styleDropDownOptions}
        />
        <hr />
        <button className="Submit" onClick={settingsButtonHandler}>
          Display Settings
        </button>
      </div>
    </>
  );
}
