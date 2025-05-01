import MenuTabs from "./MenuTabs.tsx";
import {useState} from "react";
import "./logs.css";

export default function Logs() {
  if (localStorage.getItem("PromptSetting") == null)
    localStorage.setItem("PromptSetting", "Automatic");
  if (localStorage.getItem("StyleSetting") == null)
    localStorage.setItem("StyleSetting", "Dark");

  const themeColor: string | null = localStorage.getItem("StyleSetting");
  const [logText, setlogText] = useState("");

  async function readLogs(): Promise<void> {
    let responseBody;
    const response = await fetch("http://localhost:8080/getlogs", {
      method: "GET",
    });

    if (response.ok) {
      responseBody = await response.json();
      setlogText(atob(responseBody.logs));
    } else {
      alert("Error requesting logs");
    }
  }

  readLogs();

  return (
    <>
      <div className={`${themeColor}_background`} />
      <MenuTabs />
      <pre className={`${themeColor}_logTruck`}>{logText}</pre>
    </>
  );
}
