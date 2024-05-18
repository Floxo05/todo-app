import React, {useEffect, useState} from 'react';
import './TodoList.css';
import TodoItem from './TodoItem';
import {FaSignOutAlt, FaTrash} from "react-icons/fa";
import AuthHelper from "../../utils/auth/Auth";
import {useNavigate} from "react-router-dom";
import TodoService from "../../utils/todo/TodoService";
import ErrorMessage from "../Auth/ErrorMessage";

export type Todo = {
    id?: number;
    title: string;
    completed: boolean;
};

function TodoList() {
    const navigate = useNavigate();

    const [todos, setTodos] = useState<Todo[]>([]);
    const [inputValue, setInputValue] = useState("");
    const [errorMessage, setErrorMessage] = useState("");

    const addTodo = (title: string) => {
        setErrorMessage('');
        if (title === '') {
            setErrorMessage('ToDo is required');
            return;
        }

        let newTodo = {title, completed: false};
        TodoService.addTodo(newTodo).then((res: any) => {
            const newTodos = [...todos, {id: res.todo_id, title, completed: false}];
            setTodos(newTodos);
        }).catch((error: any) => {
            setErrorMessage('Error adding todo')
            console.error(error);
        });

        setInputValue(""); // Setzt den Wert des Eingabefelds zurück
    };

    const loadTodos =  () => {
        TodoService.getTodos().then((res: any) => {
            if (res.todos) {
                setTodos(res.todos);
            }
        }).catch((error: any) => {
            setErrorMessage(prevState => 'Error loading todos')
            // console.error(error);
        });
    }

    useEffect(() => {
        // load todos
        loadTodos();
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

    const toggleTodo = (todo: Todo) => {

        const newTodos: Todo[] = [...todos];
        todo.completed = !todo.completed;

        TodoService.updateTodoStatus(todo)
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
    };

    const updateTodo = (todo: Todo, title: string) => {
        const newTodos: Todo[] = [...todos];
        const index: number = newTodos.findIndex(t => t.id === todo.id);
        newTodos[index].title = title;
        setTodos(newTodos);
    }

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

    return (
        <>
            <div className="header">
                <h1>
                    ToDo App
                    <FaSignOutAlt onClick={doLogout} className={'react-icon'} title={'abmelden'}/>
                </h1>
            </div>

            {todos.map((todo, index) => (
                <TodoItem key={index} todo={todo} toggleTodo={toggleTodo} index={index} updateTodo={updateTodo}/>
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
                        todos.some(todo => todo.completed) && (<button
                            className="delete-button"
                            onClick={deleteCompletedTodos} // Ruft die Funktion deleteCompletedTodos auf, wenn geklickt wird
                        >
                            Delete completed ToDos

                        </button>)}

                </div>
            </div>
            {errorMessage && <ErrorMessage message={errorMessage}/>}
        </>
    );
}

export default TodoList;