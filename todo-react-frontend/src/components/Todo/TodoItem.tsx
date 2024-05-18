import React, { useState } from 'react';
import './TodoItem.css';
import {Todo} from "./TodoList";
import {FaAsymmetrik, FaEdit, FaSave} from 'react-icons/fa';

type TodoItemProps = {
    todo: Todo;
    toggleTodo: (todo: Todo) => void;
    index: number;
    updateTodo: (todo: Todo, title: string) => void;
};

const TodoItem: React.FC<TodoItemProps> = ({ todo, toggleTodo, index, updateTodo }) => {
    const [isEditing, setIsEditing] = useState(false);
    const [title, setTitle] = useState(todo.title);

    const handleTodoClick = () => {
        if (!isEditing) {
            toggleTodo(todo);
        }
    }

    const handleEditClick = () => {
        console.log(isEditing);
        setIsEditing(!isEditing);
    }

    const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        setTitle(e.target.value);
    }

    const handleInputBlur = () => {
        updateTodo(todo, title);
        setIsEditing(false);
    }

    return (
        <div className={`todo-item ${todo.completed ? 'completed' : ''}`}>
            <div className="text-container">
                <input type="checkbox" checked={todo.completed} onChange={handleTodoClick}/>
                {isEditing ? (
                    <input type="text" value={title} onChange={handleInputChange} onBlur={handleInputBlur} autoFocus />
                ) : (
                    <span onClick={handleTodoClick}>{todo.title}</span>
                )}
            </div>
            {!todo.completed && (
                <>
                    {isEditing ? (
                        <div className="icon-container">
                            <FaSave onClick={() => setIsEditing(false)}/>
                        </div>
                    ) : (
                        <div className="icon-container">
                            <FaEdit onClick={() => setIsEditing(true)}/>
                        </div>)}
                </>
            )}
        </div>
    );
};

export default TodoItem;