import { useNavigate } from "react-router-dom";

export default function MenuTabs() {

    // You literally just have to do this. I have no idea why.
    const navigate = useNavigate();

    const promptingEventHandler = () => {
        navigate("/");
    };
    
    const settingsEventHandler = () => {
        navigate("/settings");
    };

    const modelsEventHandler = () => {
        navigate("/models")
    }

    const logsEventHandler = () => {
        navigate("/logs")
    }

    return(
        <div>
            <button onClick={modelsEventHandler}>Models</button>
            <button onClick={logsEventHandler}>Logs</button>
            <button onClick={promptingEventHandler}>Prompting</button>
            <button onClick={settingsEventHandler}>Settings</button>
        </div>
    );
}