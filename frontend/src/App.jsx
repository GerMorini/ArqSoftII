import { useState } from 'react';
import { Routes, Route, useNavigate } from 'react-router-dom';
import Login from './pages/Login.jsx';
import Register from './pages/Register.jsx';
import Actividades from './pages/Actividades.jsx';
import MisActividades from './pages/MisActividades.jsx';
import AdminPanel from './pages/AdminPanel.jsx';
import Layout from './components/Layout.jsx';
import Home from './pages/Home.jsx';
import Contact from './pages/Contact.jsx';
import NotFound from './pages/NotFound.jsx';
import AlertDialog from './components/AlertDialog.jsx';
import ProtectedRoute from './components/ProtectedRoute.jsx';

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
          <Route path="/actividades" element={<Actividades />} />
          <Route path="/mis-actividades" element={<MisActividades />} />
          <Route path="/contacto" element={<Contact />} />
          <Route path="/admin" element={<ProtectedRoute><AdminPanel setAlertDialog={setAlertDialog} /></ProtectedRoute>} />
        </Route>
        <Route path="*" element={<NotFound />} />
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
