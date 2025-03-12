import DropDownButton from "./DropDownButton.tsx";
import MenuTabs from "./MenuTabs.tsx";
import {useState} from "react";

export default function Settings() {
  const [PromptSetting, setPromptSetting] = useState("Automatic");
  const [StyleSetting, setStyleSetting] = useState("Pink");

  const promptDropDownOptions = [
    { type: "Automatic", name: "Automatic" },
    { type: "Manual", name: "Manual" },
    { type: "Mixed", name: "Mixed" }
  ];
  const styleDropDownOptions = [
    { type: "Pink", name: "Pink" },
    { type: "Dark", name: "Dark" },
    { type: "Light", name: "Light" }
  ];

  function settingsButtonHandler() {
    localStorage.setItem("PromptSetting", PromptSetting);
    localStorage.setItem("StyleSetting", StyleSetting);
  }

  return (
    <>
      <MenuTabs />
      <div>
        <h2>Settings</h2>
        <h4>Prompting Settings</h4>
        <DropDownButton value={PromptSetting} callBack={setPromptSetting} optionObject={promptDropDownOptions}/>
        <hr />
        <h4>Style Settings</h4>
        <DropDownButton value={StyleSetting} callBack={setStyleSetting} optionObject={styleDropDownOptions}/>
        <hr />
        <button onClick={settingsButtonHandler}>Display Settings</button>
      </div>
    </>
  );
}
