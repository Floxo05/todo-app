import React, { useState } from 'react';
import './Auth.css';
import AuthHelper from "../../utils/auth/Auth";
import {Link, useNavigate} from "react-router-dom";
import ErrorMessage from "./ErrorMessage";

export interface RegisterFormValues {
    username: string;
    password: string;
    confirmPassword: string;
    email: string;
}

function RegisterForm() {
    const navigate = useNavigate();

    const [errorMessage, setErrorMessage] = useState<string | null>(null);
    const [formValues, setFormValues] = useState<RegisterFormValues>({ username: '', password: '', confirmPassword: '', email: '' });

    const changeFormValues = (key: keyof RegisterFormValues, value: string) => {
        setFormValues({ ...formValues, [key]: value });
    }

    const handleRegister = async (event: React.FormEvent) => {
        event.preventDefault();

        // check if username has at least 3 characters
        if (formValues.username.length < 3) {
            setErrorMessage('Username must be at least 3 characters long');
            return;
        }

        // Check if password and confirmPassword are the same
        if (formValues.password !== formValues.confirmPassword) {
            setErrorMessage('Passwords do not match');
            return;
        }

        if (!AuthHelper.validatePassword(formValues.password)) {
            setErrorMessage("Password must be at least 8 characters long and contain at least one number and one special character.");
            return;
        }

        try {
            await AuthHelper.register(formValues);
            navigate('/');
        } catch (error: any) {
            setErrorMessage(error.message);
        }
    };

    return (
        <>
            <h1>ToDo App - Register</h1>
            <form onSubmit={handleRegister} className={"auth-form"}>
                <label>
                    Username: <br/>
                    <input type="text" value={formValues.username}
                           onChange={e => changeFormValues('username', e.target.value)}/>
                </label>
                <label>
                    Password: <br />
                    <input type="password" value={formValues.password} onChange={e => changeFormValues('password', e.target.value)}/>
                </label>
                <label>
                    Confirm Password: <br/>
                    <input type="password" value={formValues.confirmPassword}
                           onChange={e => changeFormValues('confirmPassword', e.target.value)}/>
                </label>
                <button type="submit" value="Register"> Register</button>
                <Link to="/login" className={'link'}>Already an account? Log in here!</Link>
            </form>
            {errorMessage && <ErrorMessage message={errorMessage} />}
        </>
    );
}

export default RegisterForm;