import "./App.css";
import {useState} from "react";
interface dropDownSettings {
  className: string;
  value: string;
  callBack: (event: string) => void;
  optionObject: {name: string; type: string | number}[];
}

export default function DropDownButton({
  className,
  value,
  callBack,
  optionObject,
}: dropDownSettings) {
    if (localStorage.getItem("StyleSetting") == null)
      localStorage.setItem("StyleSetting", "Dark");

    addEventListener("changedColorTheme", () => {
      setThemeColor(localStorage.getItem("StyleSetting"));
    });
  
    const [ ThemeColor, setThemeColor ] = useState(localStorage.getItem("StyleSetting"));
  return (
    <select
      className={`${ThemeColor}_${className}`}
      value={value}
      onChange={(event) => {
        callBack(event.target.value);
      }}
    >
      {optionObject.map((type) => {
        return (
          <option value={type.type} key={type.type}>
            {type.name}
          </option>
        );
      })}
    </select>
  );
}
