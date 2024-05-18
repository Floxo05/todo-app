import {Navigate} from "react-router-dom";
import { useState, useEffect } from 'react';
import AuthHelper from "../../utils/auth/Auth";

interface ProtectedRouteProps {
    children: any;
}

const ProtectedRoute = ({children}: ProtectedRouteProps) => {
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
        return null; // Oder eine Ladeanzeige, während die Token-Überprüfung durchgeführt wird
    }

    if (!isValidToken) {
        return <Navigate to="/login" replace />;
    }

    return children;
};

export default ProtectedRoute;