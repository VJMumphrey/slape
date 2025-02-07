import MenuTabs from "./MenuTabs.tsx";
import { useState } from "react";

export default function Settings() {

    const [ promptSetting, setPromptSetting ] = useState("Automatic")

    const radioOptions = ["Automatic", "Manual", "Mixed"];

    const radioElements = radioOptions.map((option) => {
        return (
            <div>
                <label htmlFor={option}>{`${option}:`}</label>
                <input type="radio" id={option} value={option} checked={ promptSetting === option} onChange={() => {setPromptSetting(option)}}/>
            </div>
        )
    })

    return(
        <>
            <MenuTabs/>
            <div>
                <h3>Settings</h3>
                <p>Prompting Settings</p>
                {radioElements}
            </div>
        </>
    );
}