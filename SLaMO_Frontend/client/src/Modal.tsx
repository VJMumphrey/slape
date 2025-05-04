import { useState, ReactNode } from "react";
import "./modal.css";

interface modalProperties {
  isOpen: boolean;
  onClose: () => void;
  children: ReactNode;
}

export default function Modal({isOpen, onClose, children}: modalProperties) {
  if (localStorage.getItem("StyleSetting") == null)
    localStorage.setItem("StyleSetting", "Dark");

  addEventListener("changedColorTheme", () => {
    setThemeColor(localStorage.getItem("StyleSetting"));
  });

  const [ ThemeColor, setThemeColor ] = useState(localStorage.getItem("StyleSetting"));

  if (!isOpen) {
    return null;
  }

  return (
    <>
      <div className={`${ThemeColor}_overlay`}>
        <div className="modal-content">
          <button className={`${ThemeColor}_modal-close`} onClick={onClose}>
            &times;
          </button>
          {children}
        </div>
      </div>
      <div className="dimGuy"></div>
    </>
  );
}
