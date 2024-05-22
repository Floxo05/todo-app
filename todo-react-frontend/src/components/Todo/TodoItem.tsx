import React, {useState} from 'react';
import './TodoItem.css';
import {Todo} from "./TodoList";
import {FaEdit, FaSave} from 'react-icons/fa';

type TodoItemProps = {
    todo: Todo;
    updateTodo: (todo: Todo) => void;
};

const TodoItem: React.FC<TodoItemProps> = ({todo, updateTodo}) => {
    const [isEditing, setIsEditing] = useState(false);
    const [title, setTitle] = useState(todo.title);

    const handleTodoClick = () => {
        if (!isEditing) {
            todo.completed = !todo.completed;
            updateTodo(todo);
        }
    }

    const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        setTitle(e.target.value);
    }

    const handleInputBlur = () => {
        todo.title = title;
        updateTodo(todo);
        setIsEditing(false);
    }

    return (
        <div className={`todo-item ${todo.completed ? 'completed' : ''}`}>
            <div className="text-container">
                <input type="checkbox" checked={todo.completed} onChange={handleTodoClick}/>
                {isEditing ? (
                    <input type="text" value={title} onChange={handleInputChange} onBlur={handleInputBlur} autoFocus/>
                ) : (
                    <span onClick={handleTodoClick}>{todo.title}</span>
                )}
            </div>
            <div className={'icon-container'}>
                {!todo.completed && (
                    <div className={'icon'}>
                        {isEditing ? (
                            <FaSave onClick={() => setIsEditing(false)} title={'edit'}/>
                        ) : (
                            <FaEdit onClick={() => setIsEditing(true)} title={'edit'}/>
                        )}
                    </div>
                )}
            </div>
        </div>
    );
};

export default TodoItem;