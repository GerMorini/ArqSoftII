import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import '../styles/index.css'
import Login from '../pages/Login.jsx'
import Register from '../pages/Register.jsx'
import Actividades from '../pages/Actividades.jsx'
import AdminPanel from '../pages/AdminPanel.jsx'
import Layout from '../components/Layout.jsx'
import Home from '../pages/Home.jsx'
import { BrowserRouter, Routes, Route } from "react-router-dom"

createRoot(document.getElementById('root')).render(
  <StrictMode>
    <BrowserRouter>
      <Routes>
        <Route path="/login" element={<Login />} />
        <Route path="/register" element={<Register />} />
        <Route path="/" element={<Layout />}>
          <Route index element={<Home />} />
          <Route path="actividades" element={<Actividades />} />
          <Route path="admin" element={<AdminPanel />} />
        </Route>
      </Routes>
    </BrowserRouter>
  </StrictMode>,
)