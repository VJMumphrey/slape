import {useNavigate} from "react-router-dom";
import "./menuTabs.css"

export default function MenuTabs() {
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
    <div className="top">
      <button onClick={pipelinesEventHandler}>Pipelines</button>
      <button onClick={logsEventHandler}>Logs</button>
      <button onClick={promptingEventHandler}>Prompting</button>
      <button onClick={settingsEventHandler}>Settings</button>
    </div>
  );
}
