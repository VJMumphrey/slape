import {useNavigate} from "react-router-dom";
import "./menuTabs.css";

if (localStorage.getItem("PromptSetting") == null)
  localStorage.setItem("PromptSetting", "Automatic");
if (localStorage.getItem("StyleSetting") == null)
  localStorage.setItem("StyleSetting", "Dark");

const themeColor: string | null = localStorage.getItem("StyleSetting");

export default function MenuTabs() {
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

  const settingsEventHandler = () => {
    navigate("/settings");
  };

  const pipelinesEventHandler = () => {
    navigate("/pipelines");
  };

  const logsEventHandler = () => {
    navigate("/logs");
  };

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
      <button
        className={`${themeColor}_tabButton`}
        onClick={settingsEventHandler}
      >
        Settings
      </button>
      <button className={`shutdown`} onClick={shutdownHandler}>
        X
      </button>
    </div>
  );
}
