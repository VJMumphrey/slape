import { useNavigate } from "react-router-dom";

export default function MenuTabs() {

    // You literally just have to do this. I have no idea why.
    const navigate = useNavigate();

    const settingsEventHandler = () => {
        navigate("/settings");
    };

    const promptingEventHandler = () => {
        navigate("/");
    }

    return(
        <div>
            <button>Models</button>
            <button>Logs</button>
            <button onClick={promptingEventHandler}>Prompting</button>
            <button onClick={settingsEventHandler}>Settings</button>
        </div>
    );
}