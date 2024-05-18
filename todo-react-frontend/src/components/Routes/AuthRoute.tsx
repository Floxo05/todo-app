import {Navigate} from "react-router-dom";
import { useState, useEffect } from 'react';
import AuthHelper from "../../utils/auth/Auth";

interface AuthRouteProps {
    children: any;
}

const AuthRoute = ({children}: AuthRouteProps) => {
    const [isValidToken, setIsValidToken] = useState<boolean | null>(null);

    useEffect(() => {
        AuthHelper.checkToken().then(isValid => {
            setIsValidToken(isValid);
            if (!isValid) {
                localStorage.removeItem('token');
            }
        });
    }, []);

    if (isValidToken === null) {
        return null;
    }

    if (isValidToken) {
        return <Navigate to="/" replace />;
    }

    return children;
};

export default AuthRoute;