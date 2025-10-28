import { getTokenPayload } from '../utils/tokenUtils';

/**
 * Hook para acceder al estado del usuario actual decodificando el JWT
 * Previene manipulación de datos de sesión (solo el JWT firmado importa)
 */
const useCurrentUser = () => {
    const token = localStorage.getItem("access_token");

    if (!token) {
        return {
            userId: null,
            isLoggedIn: false,
            isAdmin: false,
            username: "Usuario"
        };
    }

    try {
        const payload = getTokenPayload(token);

        return {
            userId: payload.id_usuario,
            isLoggedIn: true,
            isAdmin: payload.is_admin,
            username: payload.username || "Usuario"
        };
    } catch (error) {
        return {
            userId: null,
            isLoggedIn: false,
            isAdmin: false,
            username: "Usuario"
        };
    }
};

export default useCurrentUser;
