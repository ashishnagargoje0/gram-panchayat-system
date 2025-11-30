import React, { useState } from 'react';
import { Outlet, Link, useLocation } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import {
  Menu,
  X,
  Home,
  FileText,
  MessageSquare,
  DollarSign,
  Bell,
  User,
  Settings,
  LogOut,
  Users,
  ClipboardList,
} from 'lucide-react';

const MainLayout = () => {
  const [sidebarOpen, setSidebarOpen] = useState(true);
  const { user, logout } = useAuth();
  const location = useLocation();

  const isAdmin = user?.role === 'admin';

  const citizenMenuItems = [
    { path: '/dashboard', label: 'Dashboard', icon: Home },
    { path: '/services', label: 'Services', icon: FileText },
    { path: '/my-applications', label: 'My Applications', icon: ClipboardList },
    { path: '/complaints', label: 'Complaints', icon: MessageSquare },
    { path: '/property-tax', label: 'Property Tax', icon: DollarSign },
    { path: '/notices', label: 'Notices', icon: Bell },
  ];

  const adminMenuItems = [
    { path: '/dashboard', label: 'Dashboard', icon: Home },
    { path: '/admin/users', label: 'Users', icon: Users },
    { path: '/admin/applications', label: 'Applications', icon: FileText },
    { path: '/admin/complaints', label: 'Complaints', icon: MessageSquare },
    { path: '/property-tax', label: 'Tax Management', icon: DollarSign },
    { path: '/notices', label: 'Notices', icon: Bell },
  ];

  const menuItems = isAdmin ? adminMenuItems : citizenMenuItems;

  const isActive = (path) => location.pathname === path;

  return (
    <div className="flex h-screen bg-gray-50">
      {/* Sidebar */}
      <div
        className={`${
          sidebarOpen ? 'w-64' : 'w-20'
        } bg-gradient-to-b from-blue-600 to-purple-600 text-white transition-all duration-300 flex flex-col`}
      >
        <div className="p-4 flex items-center justify-between">
          {sidebarOpen && (
            <h2 className="text-xl font-bold">Gram Panchayat</h2>
          )}
          <button
            onClick={() => setSidebarOpen(!sidebarOpen)}
            className="p-2 hover:bg-white/10 rounded-lg"
          >
            {sidebarOpen ? <X size={20} /> : <Menu size={20} />}
          </button>
        </div>

        <nav className="flex-1 px-2 py-4 space-y-2">
          {menuItems.map((item) => (
            <Link
              key={item.path}
              to={item.path}
              className={`w-full flex items-center gap-3 px-4 py-3 rounded-lg transition ${
                isActive(item.path)
                  ? 'bg-white/20'
                  : 'hover:bg-white/10'
              }`}
            >
              <item.icon size={20} />
              {sidebarOpen && <span>{item.label}</span>}
            </Link>
          ))}
        </nav>

        <div className="p-4 space-y-2">
          <Link
            to="/profile"
            className={`w-full flex items-center gap-3 px-4 py-3 rounded-lg transition ${
              isActive('/profile') ? 'bg-white/20' : 'hover:bg-white/10'
            }`}
          >
            <Settings size={20} />
            {sidebarOpen && <span>Settings</span>}
          </Link>
          <button
            onClick={logout}
            className="w-full flex items-center gap-3 px-4 py-3 rounded-lg hover:bg-white/10 transition"
          >
            <LogOut size={20} />
            {sidebarOpen && <span>Logout</span>}
          </button>
        </div>
      </div>

      {/* Main Content */}
      <div className="flex-1 flex flex-col overflow-hidden">
        {/* Header */}
        <header className="bg-white shadow-sm px-8 py-4">
          <div className="flex items-center justify-between">
            <h1 className="text-2xl font-bold text-gray-800">
              {menuItems.find((item) => item.path === location.pathname)
                ?.label || 'Gram Panchayat'}
            </h1>
            <div className="flex items-center gap-4">
              <button className="p-2 hover:bg-gray-100 rounded-lg transition relative">
                <Bell size={20} />
                <span className="absolute top-1 right-1 w-2 h-2 bg-red-500 rounded-full"></span>
              </button>
              <div className="flex items-center gap-3">
                <div className="w-10 h-10 rounded-full bg-gradient-to-r from-blue-500 to-purple-500 flex items-center justify-center text-white font-bold">
                  {user?.first_name?.[0] || user?.email?.[0]?.toUpperCase()}
                </div>
                <div className="text-sm">
                  <p className="font-medium text-gray-800">
                    {user?.first_name} {user?.last_name}
                  </p>
                  <p className="text-gray-500 text-xs capitalize">
                    {user?.role}
                  </p>
                </div>
              </div>
            </div>
          </div>
        </header>

        {/* Page Content */}
        <main className="flex-1 overflow-auto p-8">
          <Outlet />
        </main>
      </div>
    </div>
  );
};

export default MainLayout;
