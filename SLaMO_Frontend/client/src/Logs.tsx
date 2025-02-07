import MenuTabs from "./MenuTabs.tsx"
import logImageImportantDoNotRemove from "./logs.jpg";

export default function Logs() {

    function logHandler() {
        alert("Vito a Bitch");
    }

    return(
        <>
            <MenuTabs/>
            <div>
                <img onClick={logHandler} src={logImageImportantDoNotRemove} alt="Image" />
            </div>
        </>
    );
}