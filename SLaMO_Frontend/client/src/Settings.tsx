import * as fs from "node:fs";
import MenuTabs from "./MenuTabs.tsx";
import { useState } from "react";

export default function Settings() {

    const [ PromptSetting, setPromptSetting ] = useState("Automatic")
    const [ StyleSetting, setStyleSetting ] = useState("Pink")
    const [ SettingsObject, setSettingsObject ] = useState({promptType: "", styleType: ""})

    const promptRadioOptions = ["Automatic", "Manual", "Mixed"];
    const styleRadioOptions = ["Pink", "Dark", "Light"];

    const promptRadioElements = promptRadioOptions.map((option) => {
        return (
            <div>
                <label htmlFor={option}>{`${option}:`}</label>
                <input type="radio" id={option} value={option} checked={ PromptSetting === option} onChange={() => {setPromptSetting(option); setSettingsObject({...SettingsObject, promptType: option});}}/>
            </div>
        )
    })
    
    const styleRadioElements = styleRadioOptions.map((option) => {
        return (
            <div>
                <label htmlFor={option}>{`${option}:`}</label>
                <input type="radio" id={option} value={option} checked={ StyleSetting === option} onChange={() => {setStyleSetting(option); setSettingsObject({...SettingsObject, styleType: option});}}/>
            </div>
        )
    })

    async function settingsButtonHandler() {
        console.log(SettingsObject);



    }

    return(
        <>
            <MenuTabs/>
            <div>
                <h3>Settings</h3>
                <p>Prompting Settings</p>
                {promptRadioElements}
                <hr/>
                {styleRadioElements}
                <hr/>
                <button onClick={settingsButtonHandler}>Display Settings</button>
            </div>
        </>
    );
}