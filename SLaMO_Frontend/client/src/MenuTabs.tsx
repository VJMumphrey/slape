import { useState } from "react";
import {useNavigate} from "react-router-dom";
import "./menuTabs.css";
import { l } from "../../../../../AppData/Local/deno/npm/registry.npmjs.org/react-router/7.1.5/dist/development/fog-of-war-CCAcUMgB.d.ts";

if (localStorage.getItem("StyleSetting") == null)
  localStorage.setItem("StyleSetting", "Dark");

let themeColor: string | null = localStorage.getItem("StyleSetting");

export default function MenuTabs() {

  const [ CurrentTheme, setCurrentTheme ] = useState(themeColor==="Light" ? "\u2600" : "\u263D");

  function shutdownHandler() {
    fetch("http://localhost:8080/shutdown", {
      method: "GET",
    });
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
    if (themeColor === "Light" )  {
      localStorage.setItem("StyleSetting", "Dark");
      setCurrentTheme("\u263D");
    } else {
      localStorage.setItem("StyleSetting", "Light");
      setCurrentTheme("\u2600");
    }
    globalThis.location.reload();
  }

  return (
    <div className={`${themeColor}_top`}>
      <button
        className={`${themeColor}_tabButton`}
        onClick={pipelinesEventHandler}
      >
        Pipelines
      </button>
      <button className={`${themeColor}_tabButton`} onClick={logsEventHandler}>
        Logs
      </button>
      <button
        className={`${themeColor}_tabButton`}
        onClick={promptingEventHandler}
      >
        Prompting
      </button>
      <button className={`${themeColor}_themeButton`} onClick={handleThemeButton}>
        {CurrentTheme}
      </button>
      <button className={`shutdown`} onClick={shutdownHandler}>
        &times;
      </button>
    </div>
  );
}
