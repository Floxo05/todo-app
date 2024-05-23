import React from 'react';
import Modal from 'react-modal';
import './Modal.css';

type ShareModalProps = {
    isOpen: boolean;
    onRequestClose: () => void;
    onShareConfirm: () => void;
    username: string;
    setUsername: (username: string) => void;
    todosCount: number;
};

const ShareModal: React.FC<ShareModalProps> = ({isOpen, onRequestClose, onShareConfirm, username, setUsername, todosCount}) => {
    const handleInputChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setUsername(event.target.value);
    };

    return (
        <Modal
            isOpen={isOpen}
            onRequestClose={onRequestClose}
            contentLabel="Share Todos"
            className="modal-content"
            overlayClassName="modal-overlay"
        >
            <h2>Would you like to share {todosCount} {todosCount === 1 ? 'todo' : 'todos'}?</h2>
            <input
                type="text"
                placeholder="Enter username"
                value={username}
                onChange={handleInputChange}
            />
            <div className="modal-buttons">
                <button onClick={onShareConfirm}>Share</button>
                <button onClick={onRequestClose}>Abort</button>
            </div>

        </Modal>
    );
};

export default ShareModal;