import { useState } from 'react';
import { Routes, Route, useNavigate } from 'react-router-dom';
import Login from './pages/Login.jsx';
import Register from './pages/Register.jsx';
import Actividades from './pages/Actividades.jsx';
import AdminPanel from './pages/AdminPanel.jsx';
import Layout from './components/Layout.jsx';
import Home from './pages/Home.jsx';
import AlertDialog from './components/AlertDialog.jsx';

function App() {
  const [alertDialog, setAlertDialog] = useState(null);
  const navigate = useNavigate();

  const closeAlertDialog = () => {
    // Si el diálogo es por token expirado, redirigir al login
    if (alertDialog?.isTokenExpired) {
      navigate('/login');
    }
    setAlertDialog(null);
  };

  return (
    <>
      <Routes>
        <Route path="/login" element={<Login />} />
        <Route path="/register" element={<Register />} />
        <Route path="/" element={<Layout setAlertDialog={setAlertDialog} />}>
          <Route index element={<Home />} />
          <Route path="actividades" element={<Actividades />} />
          <Route path="admin" element={<AdminPanel setAlertDialog={setAlertDialog} />} />
        </Route>
      </Routes>

      {alertDialog && (
        <AlertDialog
          title={alertDialog.title}
          message={alertDialog.message}
          type={alertDialog.type}
          onClose={closeAlertDialog}
        />
      )}
    </>
  );
}

export default App;
