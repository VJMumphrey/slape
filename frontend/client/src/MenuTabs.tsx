import { useState } from "react";
import {useNavigate} from "react-router-dom";
import "./menuTabs.css";
import { l } from "../../../../../AppData/Local/deno/npm/registry.npmjs.org/react-router/7.1.5/dist/development/fog-of-war-CCAcUMgB.d.ts";

if (localStorage.getItem("StyleSetting") == null)
  localStorage.setItem("StyleSetting", "Dark");

export default function MenuTabs() {

  addEventListener("changedColorTheme", () => {
    setThemeColor(localStorage.getItem("StyleSetting"));
  });

  const [ ThemeColor, setThemeColor ] = useState(localStorage.getItem("StyleSetting"));
  const [ CurrentTheme, setCurrentTheme ] = useState(localStorage.getItem("StyleSetting") === "Light" ? "\u2600" : "\u263D");

  async function shutdownHandler() {
    const response = await fetch("http://localhost:8080/shutdownpipes", {
      method: "GET",
    });

    if(response.ok) {
      alert("Shutting Down All Pipelines");
    } else {
      alert("Shutdown Failed!");
    }
  }
  // You literally just have to do this. I have no idea why.
  const navigate = useNavigate();

  const promptingEventHandler = () => {
    navigate("/");
  };

  const pipelinesEventHandler = () => {
    navigate("/pipelines");
  };

  const logsEventHandler = () => {
    navigate("/logs");
  };

  const handleThemeButton = () => {
    if (ThemeColor === "Light" )  {
      localStorage.setItem("StyleSetting", "Dark");
      setCurrentTheme("\u263D");
    } else {
      localStorage.setItem("StyleSetting", "Light");
      setCurrentTheme("\u2600");
    }
    dispatchEvent(new Event("changedColorTheme"));
  }

  return (
    <div className={`${ThemeColor}_top`}>
      <span className={"dropShadow"}>
        <button
          className={`${ThemeColor}_tabButton_Pipeline`}
          onClick={pipelinesEventHandler}
        >
          Pipelines
        </button>
        <button className={`${ThemeColor}_tabButton_Log`} onClick={logsEventHandler}>
          Logs
        </button>
        <button
          className={`${ThemeColor}_tabButton_Prompt`}
          onClick={promptingEventHandler}
        >
          Prompting
        </button>
      </span>
      <button className={`${ThemeColor}_themeButton`} onClick={handleThemeButton}>
        {CurrentTheme}
      </button>
      <button className={`shutdown`} onClick={shutdownHandler}>
        &times;
      </button>
    </div>
  );
}
