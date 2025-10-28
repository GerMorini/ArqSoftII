/**
 * Decodifica un JWT token y extrae el payload
 * @param {string} token - El token JWT
 * @returns {object|null} - El payload decodificado o null si no es válido
 */
export function getTokenPayload(token) {
  try {
    const parts = token.split('.');
    if (parts.length !== 3) {
      return null;
    }
    const decoded = atob(parts[1]);
    return JSON.parse(decoded);
  } catch {
    return null;
  }
}

/**
 * Verifica si un token está expirado
 * @param {string} token - El token JWT
 * @returns {boolean} - true si el token está expirado, false si es válido
 */
export function isTokenExpired(token) {
  const payload = getTokenPayload(token);

  if (!payload || !payload.exp) {
    return true; // Si no tiene exp claim, considerarlo expirado
  }

  // payload.exp está en segundos, Date.now() está en milisegundos
  const expirationTime = payload.exp * 1000;
  const currentTime = Date.now();

  return currentTime > expirationTime;
}

/**
 * Obtiene el tiempo restante antes de que expire el token (en segundos)
 * @param {string} token - El token JWT
 * @returns {number|null} - Segundos restantes o null si el token es inválido
 */
export function getTokenTimeRemaining(token) {
  const payload = getTokenPayload(token);

  if (!payload || !payload.exp) {
    return null;
  }

  const expirationTime = payload.exp * 1000;
  const currentTime = Date.now();
  const remaining = (expirationTime - currentTime) / 1000;

  return remaining > 0 ? remaining : 0;
}