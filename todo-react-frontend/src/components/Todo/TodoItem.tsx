import React, {useState} from 'react';
import './TodoItem.css';
import {Category, Todo} from "./TodoList";
import {FaEdit} from 'react-icons/fa';
import EditTodoModal from "./EditTodoModal";

type TodoItemProps = {
    todo: Todo;
    updateTodo: (todo: Todo) => void;
    categories: Category[];
    selectedTodos: Todo[];
    setSelectedTodos: (todos: Todo[]) => void;
};

const TodoItem: React.FC<TodoItemProps> = ({todo, updateTodo, categories, selectedTodos, setSelectedTodos}) => {
    const [isModalOpen, setIsModalOpen] = useState(false);
    const [title, setTitle] = useState(todo.title);

    const handleTodoClick = () => {
        todo.completed = !todo.completed;
        updateTodo(todo);
    }

    const onEditConfirm = (title: string, category: string) => {
        todo.title = title;
        todo.category.title = category;
        updateTodo(todo);
        setIsModalOpen(false);
    }

    const handleSelect = () => {
        if (selectedTodos.includes(todo)) {
            setSelectedTodos(selectedTodos.filter(t => t !== todo));
        } else {
            setSelectedTodos([...selectedTodos, todo]);
        }
    }

    return (
        <div className={`todo-item ${todo.completed ? 'completed' : ''}`}>
            <div className="text-container">
                <input type="checkbox" checked={selectedTodos.includes(todo)} onChange={handleSelect}/>
                <span
                    onClick={handleTodoClick}>{todo.title}{todo.category.title !== '' ? ' - ' + todo.category.title : ''}
                </span>
            </div>
            <div className={'icon-container'}>
                {!todo.completed && (
                    <div className={'icon'}>
                        <FaEdit onClick={() => setIsModalOpen(true)} title={'edit'}/>
                    </div>
                )}
            </div>
            <EditTodoModal
                isOpen={isModalOpen}
                onRequestClose={() => setIsModalOpen(false)}
                onEditConfirm={onEditConfirm}
                initialTitle={todo.title}
                initialCategory={todo.category.title}
                categories={categories}
            />
        </div>
    );
};

export default TodoItem;