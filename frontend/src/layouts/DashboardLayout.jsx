import React from "react";
import { Link, Outlet } from "react-router-dom";

export default function DashboardLayout() {
  return (
    <div className="min-h-screen grid grid-cols-12">
      {/* Sidebar */}
      <aside className="col-span-2 bg-gray-900 text-white p-6 space-y-6">
        <h2 className="text-xl font-bold">Gram Panchayat</h2>

        <nav className="space-y-3">
          <Link to="/dashboard" className="block hover:text-blue-300">Dashboard</Link>
          <Link to="/applications" className="block hover:text-blue-300">Applications</Link>
          <Link to="/complaints" className="block hover:text-blue-300">Complaints</Link>
          <Link to="/properties" className="block hover:text-blue-300">Properties</Link>
          <Link to="/notices" className="block hover:text-blue-300">Notices</Link>
          <Link to="/payments" className="block hover:text-blue-300">Payments</Link>
          <Link to="/schemes" className="block hover:text-blue-300">Schemes</Link>
          <Link to="/citizens" className="block hover:text-blue-300">Citizens</Link>
        </nav>
      </aside>

      {/* Main Content */}
      <main className="col-span-10 bg-gray-50 p-8">
        <Outlet />
      </main>
    </div>
  );
}
