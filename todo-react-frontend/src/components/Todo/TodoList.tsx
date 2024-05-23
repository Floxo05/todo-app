import React, {useEffect, useState} from 'react';
import './TodoList.css';
import TodoItem from './TodoItem';
import {FaSignOutAlt} from "react-icons/fa";
import AuthHelper from "../../utils/auth/Auth";
import {useNavigate} from "react-router-dom";
import TodoService from "../../utils/todo/TodoService";
import ErrorMessage from "../Auth/ErrorMessage";
import ShareModal from "./ShareModal";
import CategoryService from "../../utils/todo/CategoryService";

export type Todo = {
    id: number;
    title: string;
    completed: boolean;
    owner_id: number;
    category: Category;
};

export type Category = {
    id: number;
    title: string;
    create_user_id: number;
}

function TodoList() {
    const navigate = useNavigate();

    const [todos, setTodos] = useState<Todo[]>([]);
    const [inputValue, setInputValue] = useState("");
    const [errorMessage, setErrorMessage] = useState("");
    const [isModalOpen, setIsModalOpen] = useState(false);
    const [username, setUsername] = useState('');
    const [categories, setCategories] = useState<Category[]>([]);
    const [selectedTodos, setSelectedTodos] = useState<Todo[]>([]);

    const addTodo = (title: string) => {
        setErrorMessage('');
        if (title === '') {
            setErrorMessage('ToDo is required');
            return;
        }

        let newTodo = {title, completed: false, id: 0, owner_id: 0, category: {id: 0, title: '', create_user_id: 0}};
        TodoService.addTodo(newTodo).then((res: any) => {
            const newTodos = [...todos, {
                id: res.id,
                title,
                completed: false,
                owner_id: res.owner_id,
                category: res.category
            }];
            setTodos(newTodos);
        }).catch((error: any) => {
            setErrorMessage('Error adding todo')
            console.error(error);
        });

        setInputValue(""); // Setzt den Wert des Eingabefelds zurück
    };

    const loadCategories = () => {
        CategoryService.getCategories().then((res: any) => {
            if (res) {
                setCategories(res);
            }
        }).catch((error: any) => {
            setErrorMessage('Error loading categories')
            console.error(error);
        });
    }

    useEffect(() => {
        // load todos
        TodoService.getTodos().then((res: any) => {
            if (res) {
                setTodos(res);
            }
        }).catch((error: any) => {
            setErrorMessage('Error loading todos')
            console.error(error);
        });

        loadCategories();
    }, [])

    const deleteCompletedTodos = () => {
        // setErrorMessage('');
        const todosToDelete: Todo[] = todos.filter(todo => todo.completed);
        let remainingTodos = todos.filter(todo => !todo.completed);

        Promise.all(todosToDelete.map(todo => TodoService.deleteTodoById(todo.id || 0)))
            .then(() => {
                setTodos(remainingTodos);
            })
            .catch((error: any) => {
                setErrorMessage(error.message)
                console.error(error);
            });
    };

    const updateTodo = (todo: Todo) => {

        const newTodos: Todo[] = [...todos];

        TodoService.updateTodo(todo)
            .then(() => {
                // replace the todo with the updated one
                const index: number = newTodos.findIndex(t => t.id === todo.id);
                newTodos[index] = todo;
                setTodos(newTodos);
            })
            .catch((error: any) => {
                setErrorMessage(error.message)
                console.error(error);
            });

        loadCategories();
    };


    const handleInputChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setInputValue(event.target.value);
    };

    const handleKeyPress = (event: React.KeyboardEvent) => {
        if (event.key === 'Enter') {
            addTodo(inputValue);
        }
    };

    const doLogout = () => {
        AuthHelper.logout();
        navigate("/login", {replace: true});
    };


    const openModal = () => {
        setIsModalOpen(true);
    };

    const closeModal = () => {
        setIsModalOpen(false);
    };

    const handleShareClick = () => {
        openModal();
    };

    const handleShareConfirm = () => {
        Promise.all(selectedTodos.map(todo => TodoService.shareTodo(username, todo)))
            .then(() => {

            })
            .catch((error: any) => {
                setErrorMessage(error.message)
                console.error(error);
            });

        closeModal();
    };

    return (
        <>
            <div className="header">
                <h1>
                    ToDo App
                    <FaSignOutAlt onClick={doLogout} className={'react-icon'} title={'logout'}/>
                </h1>
            </div>

            {todos.map((todo, index) => (
                <TodoItem key={index} todo={todo} updateTodo={updateTodo} categories={categories}
                          selectedTodos={selectedTodos} setSelectedTodos={setSelectedTodos}/>
            ))}
            <div className={"form"}>
                <label htmlFor="todo-input">New ToDo:</label>
                <input
                    type="text"
                    className="todo-input"
                    placeholder="New todo"
                    id="todo-input"
                    value={inputValue} // Bindet den Wert des Eingabefelds an den Zustand
                    onChange={handleInputChange} // Aktualisiert den Zustand, wenn sich der Wert des Eingabefelds ändert
                    onKeyDown={handleKeyPress}
                />
                <div className={"button-container"}>
                    <button
                        className="add-button"
                        onClick={() => addTodo(inputValue)} // Verwendet den aktuellen Zustand als Wert
                    >
                        Add ToDo
                    </button>
                    {
                        selectedTodos.length > 0 && (
                                <button
                                    className="share-button"
                                    onClick={handleShareClick}
                                >
                                    Share
                                </button>
                        )}
                    {
                        todos.some(todo => todo.completed) && (
                                <button
                                    className="delete-button"
                                    onClick={deleteCompletedTodos} // Ruft die Funktion deleteCompletedTodos auf, wenn geklickt wird
                                >
                                    Delete completed ToDos

                                </button>
                        )}

                </div>
                <ShareModal
                    isOpen={isModalOpen}
                    onRequestClose={closeModal}
                    onShareConfirm={handleShareConfirm}
                    username={username}
                    setUsername={setUsername}
                    todosCount={todos.length}
                />
            </div>
            {errorMessage && <ErrorMessage message={errorMessage}/>}
        </>
    );
}

export default TodoList;