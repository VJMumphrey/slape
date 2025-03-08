import ReactNode from 'react';
import "./modal.css"


interface modalProperties {
    isOpen: boolean,
    onClose: () => void,
    children: ReactNode
}

export default function Modal({ isOpen, onClose, children }: modalProperties) {
    if (!isOpen) {return(null);}

    return(
        <>
            <div className="overlay">
                <div className="modal-content">
                    <button className="modal-close" onClick={onClose}>
                        &times;
                    </button>
                    {children}
                </div>
            </div>
            <div className="dimGuy"></div>
        </>
  );
};