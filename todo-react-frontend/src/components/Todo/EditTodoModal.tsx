import React from 'react';
import Modal from 'react-modal';
import './Modal.css';
import {Category} from "./TodoList";

type EditTodoModalProps = {
    isOpen: boolean;
    onRequestClose: () => void;
    onEditConfirm: (title: string, category: string) => void;
    initialTitle: string;
    initialCategory: string;
    categories: Category[];
};

const EditTodoModal: React.FC<EditTodoModalProps> = ({
                                                         isOpen,
                                                         onRequestClose,
                                                         onEditConfirm,
                                                         initialTitle,
                                                         initialCategory,
                                                         categories
                                                     }) => {
    const [title, setTitle] = React.useState(initialTitle);
    const [category, setCategory] = React.useState(initialCategory);

    const handleTitleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setTitle(event.target.value);
    };

    const handleCategoryChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setCategory(event.target.value);
    };

    const handleEditConfirm = () => {
        onEditConfirm(title, category);
    };

    return (
        <Modal
            isOpen={isOpen}
            onRequestClose={onRequestClose}
            contentLabel="Edit Todo"
            className="modal-content"
            overlayClassName="modal-overlay"
        >
            <h2>Edit Todo</h2>
            <label>
                Title:
                <input
                    type="text"
                    placeholder="Enter title"
                    value={title}
                    onChange={handleTitleChange}
                />
            </label>
            <label>
                Category:
                <input
                    type="text"
                    placeholder="Enter category"
                    value={category}
                    onChange={handleCategoryChange}
                    list={"category-options"}
                />
                <datalist id="category-options">
                    {categories.map((category, index) => (
                        <option key={index} value={category.title}/>
                    ))}
                </datalist>
            </label>
            <div className="modal-buttons">
                <button onClick={handleEditConfirm}>Save</button>
                <button onClick={onRequestClose}>Cancel</button>
            </div>
        </Modal>
    );
};

export default EditTodoModal;