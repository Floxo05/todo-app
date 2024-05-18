import React from 'react';
import './App.css';
import TodoList from "./components/Todo/TodoList";
import Login from "./components/Auth/Login";
import Register from "./components/Auth/Register";
import {BrowserRouter as Router, Route, Routes} from 'react-router-dom';
import ProtectedRoute from "./components/Routes/ProtectedRoute";
import AuthRoute from "./components/Routes/AuthRoute";

function App() {

    const loginAuthRoute = <AuthRoute><Login/></AuthRoute>
    const registerAuthRoute = <AuthRoute><Register/></AuthRoute>
    const protectedTodoListRoute = <ProtectedRoute><TodoList/></ProtectedRoute>

    return (
        <Router>
            <div className="App">
                <main className="container">
                    <section className={"section"}>
                        <Routes>
                            <Route path="/login" element={loginAuthRoute} />
                            <Route path="/register" element={registerAuthRoute} />
                            <Route path="/" element={protectedTodoListRoute} />
                        </Routes>
                    </section>
                </main>
            </div>
        </Router>
    );
}

export default App;