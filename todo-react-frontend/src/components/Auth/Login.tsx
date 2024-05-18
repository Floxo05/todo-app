import React, { useState } from 'react';
import './Auth.css';
import AuthHelper from "../../utils/auth/Auth";
import {Link, useNavigate} from "react-router-dom";
import ErrorMessage from './ErrorMessage';

export interface LoginFormValues {
    username: string;
    password: string;
}

function LoginForm() {
    const navigate = useNavigate();

    const [formValues, setFormValues] = useState<LoginFormValues>({ username: '', password: '' });
    const [errorMessage, setErrorMessage] = useState("");

    const changeFormValues = (key: keyof LoginFormValues, value: string) => {
        setFormValues({ ...formValues, [key]: value });
    }

    const handleLogin = async (event: React.FormEvent) => {
        event.preventDefault();

        try {
            await AuthHelper.login(formValues);
            navigate('/');
        } catch (error: any) {
            setErrorMessage(error.message);
            console.error(error);
        }
    };

    return (
        <>
            <h1>ToDo App - Login</h1>
            <form onSubmit={handleLogin} className={"auth-form"}>
                <label>
                    Username: <br/>
                    <input type="text" value={formValues.username}
                           onChange={e => changeFormValues('username', e.target.value)}/>
                </label>
                <label>
                    Password: <br />
                    <input type="password" value={formValues.password} onChange={e => changeFormValues('password', e.target.value)}/>
                </label>
                <button type="submit" value="Login"> Login</button>
                <Link to="/register" className={'link'}>No Account? Register here!</Link>
                {errorMessage && <ErrorMessage message={errorMessage} />}
            </form>
        </>
    );
}

export default LoginForm;