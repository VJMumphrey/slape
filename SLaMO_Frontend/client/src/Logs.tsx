import MenuTabs from "./MenuTabs.tsx";
import {useState} from "react";
import "./logs.css";

export default function Logs() {
  if (localStorage.getItem("StyleSetting") == null)
    localStorage.setItem("StyleSetting", "Dark");

  addEventListener("changedColorTheme", () => {
    setThemeColor(localStorage.getItem("StyleSetting"));
  });

  const [ ThemeColor, setThemeColor ] = useState(localStorage.getItem("StyleSetting"));
  const [logText, setlogText] = useState(localStorage.getItem("logsCopy") ? localStorage.getItem("logsCopy") : "");

  async function readLogs(): Promise<void> {
    let responseBody;
    const response = await fetch("http://localhost:8080/getlogs", {
      method: "GET",
    });

    if (response.ok) {
      responseBody = await response.json();
      setlogText(atob(responseBody.logs));
      localStorage.setItem("logsCopy", atob(responseBody.logs));
    } else {
      alert("Error requesting logs");
    }
  }

  setInterval(readLogs, 3000);

  return (
    <>
      <div className={`${ThemeColor}_background`} />
      <MenuTabs />
      <pre className={`${ThemeColor}_logTruck`}>{logText}</pre>
    </>
  );
}
