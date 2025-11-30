import React from "react";
import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";
import { AuthProvider } from "./context/AuthContext";

import Login from "./pages/Login";
import Register from "./pages/Register";
import ForgotPassword from "./pages/ForgotPassword";

import DashboardLayout from "./layouts/DashboardLayout";

import Dashboard from "./pages/Dashboard";
import Applications from "./pages/Applications";
import Complaints from "./pages/Complaints";
import Notices from "./pages/Notices";
import Properties from "./pages/Properties";
import Payments from "./pages/Payments";
import Schemes from "./pages/Schemes";
import Citizens from "./pages/Citizens";

export default function App() {
  return (
    <AuthProvider>
      <BrowserRouter>
        <Routes>
          {/* Public Routes */}
          <Route path="/" element={<Navigate to="/login" replace />} />
          <Route path="/login" element={<Login />} />
          <Route path="/register" element={<Register />} />
          <Route path="/forgot-password" element={<ForgotPassword />} />

          {/* Dashboard Routes */}
          <Route path="/" element={<DashboardLayout />}>
            <Route path="/dashboard" element={<Dashboard />} />
            <Route path="/applications" element={<Applications />} />
            <Route path="/complaints" element={<Complaints />} />
            <Route path="/properties" element={<Properties />} />
            <Route path="/notices" element={<Notices />} />
            <Route path="/payments" element={<Payments />} />
            <Route path="/schemes" element={<Schemes />} />
            <Route path="/citizens" element={<Citizens />} />
          </Route>
        </Routes>
      </BrowserRouter>
    </AuthProvider>
  );
}
