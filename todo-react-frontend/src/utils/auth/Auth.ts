import {LoginFormValues} from "../../components/Auth/Login";
import {RegisterFormValues} from "../../components/Auth/Register";


class AuthHelper {
    static async login(formValues: LoginFormValues): Promise<void> {
        const response = await fetch(process.env.REACT_APP_API + "/login", {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(formValues)
        });

        const data = await response.json();

        if (!response.ok) {
            throw new Error(data.error);
        }

        // Speichern Sie den Token in localStorage
        localStorage.setItem('token', data.token);
    }

    static async register(formValues: RegisterFormValues): Promise<void>  {
        const response = await fetch(process.env.REACT_APP_API + "/register", {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(formValues)
        });

        const data = await response.json();

        if (!response.ok) {
            throw new Error(data.error);
        }

        // Speichern Sie den Token in localStorage
        localStorage.setItem('token', data.token);
    }

    static logout(): void {
        localStorage.removeItem('token');
    }

    static getToken(): string | null {
        return localStorage.getItem('token');
    }

    static async checkToken(): Promise<boolean> {
        const response = await fetch(process.env.REACT_APP_API + '/auth/check-token', {
            headers: {
                "Authorization": "Bearer " + AuthHelper.getToken() || ""
            }
        });

        return response.status === 200;
    }
}

export default AuthHelper;