import React, { useState } from 'react';
import '../styles/PasswordInput.css';

const PasswordInput = ({
  value,
  onChange,
  placeholder = "Contraseña",
  disabled = false,
  required = false,
  error = null,
  name,
  id
}) => {
  const [showPassword, setShowPassword] = useState(false);

  const toggleVisibility = () => {
    setShowPassword(!showPassword);
  };

  return (
    <div className="password-input-container">
      <input
        type={showPassword ? "text" : "password"}
        value={value}
        onChange={onChange}
        placeholder={placeholder}
        disabled={disabled}
        required={required}
        name={name}
        id={id}
        className={`password-input ${error ? 'error' : ''}`}
      />
      <button
        type="button"
        className="password-toggle-btn"
        onClick={toggleVisibility}
        disabled={disabled}
        tabIndex={-1}
        aria-label={showPassword ? "Ocultar contraseña" : "Mostrar contraseña"}
        title={showPassword ? "Ocultar contraseña" : "Mostrar contraseña"}
      >
        {showPassword ? "👁️" : "👁️‍🗨️"}
      </button>
      {error && <span className="error-text">{error}</span>}
    </div>
  );
};

export default PasswordInput;
