import { useState, useEffect } from 'react';
import { applicationService } from '../services/authService';
import { Clock, CheckCircle, XCircle, Eye } from 'lucide-react';

const MyApplications = () => {
  const [applications, setApplications] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchApplications();
  }, []);

  const fetchApplications = async () => {
    try {
      const response = await applicationService.getMyApplications();
      setApplications(response.data || []);
    } catch (error) {
      console.error('Failed to fetch applications:', error);
    } finally {
      setLoading(false);
    }
  };

  const getStatusBadge = (status) => {
    const styles = {
      pending: 'bg-yellow-100 text-yellow-700',
      under_review: 'bg-blue-100 text-blue-700',
      approved: 'bg-green-100 text-green-700',
      rejected: 'bg-red-100 text-red-700',
    };
    return styles[status] || styles.pending;
  };

  if (loading) {
    return <div className="flex justify-center py-12">
      <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
    </div>;
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h2 className="text-2xl font-bold text-gray-800">My Applications</h2>
      </div>

      <div className="bg-white rounded-xl shadow-md p-6">
        {applications.length === 0 ? (
          <div className="text-center py-12">
            <p className="text-gray-600">No applications found</p>
          </div>
        ) : (
          <div className="space-y-4">
            {applications.map((app) => (
              <div
                key={app.id}
                className="border border-gray-200 rounded-lg p-4 hover:border-blue-300 transition"
              >
                <div className="flex items-start justify-between mb-2">
                  <div>
                    <h3 className="font-bold text-gray-800 capitalize">
                      {app.type} Certificate
                    </h3>
                    <p className="text-sm text-gray-500">
                      Application No: {app.application_number}
                    </p>
                  </div>
                  <span
                    className={`px-3 py-1 rounded-full text-sm font-medium ${getStatusBadge(
                      app.status
                    )}`}
                  >
                    {app.status.replace('_', ' ').toUpperCase()}
                  </span>
                </div>
                <p className="text-sm text-gray-600 mb-3">
                  Submitted: {new Date(app.created_at).toLocaleDateString()}
                </p>
                <button className="text-blue-600 font-medium text-sm hover:underline flex items-center gap-1">
                  <Eye size={16} /> View Details
                </button>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
};
