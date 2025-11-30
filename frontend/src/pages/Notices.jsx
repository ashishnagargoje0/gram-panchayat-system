import React, { useState, useEffect } from 'react';
import { noticeService } from '../services/authService';
import { Bell, Calendar } from 'lucide-react';

const Notices = () => {
  const [notices, setNotices] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchNotices();
  }, []);

  const fetchNotices = async () => {
    try {
      const response = await noticeService.getNotices({});
      setNotices(response.data || []);
    } catch (error) {
      console.error('Failed to fetch notices:', error);
    } finally {
      setLoading(false);
    }
  };

  const getPriorityBadge = (priority) => {
    const styles = {
      urgent: 'bg-red-100 text-red-700',
      high: 'bg-orange-100 text-orange-700',
      normal: 'bg-blue-100 text-blue-700',
      low: 'bg-gray-100 text-gray-700',
    };
    return styles[priority] || styles.normal;
  };

  return (
    <div className="space-y-6">
      <div className="bg-gradient-to-r from-purple-500 to-pink-600 rounded-xl p-6 text-white">
        <h2 className="text-2xl font-bold mb-2">Notices & Announcements</h2>
        <p className="opacity-90">Stay updated with latest announcements</p>
      </div>

      {loading ? (
        <div className="flex justify-center py-12">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
        </div>
      ) : notices.length === 0 ? (
        <div className="bg-white rounded-xl shadow-md p-12 text-center">
          <Bell className="mx-auto mb-4 text-gray-400" size={48} />
          <h3 className="text-xl font-bold text-gray-800 mb-2">
            No Notices Available
          </h3>
          <p className="text-gray-600">Check back later for updates</p>
        </div>
      ) : (
        <div className="space-y-4">
          {notices.map((notice) => (
            <div
              key={notice.id}
              className="bg-white rounded-xl shadow-md p-6 hover:shadow-xl transition"
            >
              <div className="flex items-start justify-between mb-4">
                <div className="flex-1">
                  <div className="flex items-center gap-2 mb-2">
                    <h3 className="text-xl font-bold text-gray-800">
                      {notice.title}
                    </h3>
                    <span
                      className={`px-2 py-1 rounded-full text-xs font-medium ${getPriorityBadge(
                        notice.priority
                      )}`}
                    >
                      {notice.priority?.toUpperCase()}
                    </span>
                  </div>
                  <div className="flex items-center gap-4 text-sm text-gray-500">
                    <span className="flex items-center gap-1">
                      <Calendar size={16} />
                      {new Date(notice.published_at || notice.created_at).toLocaleDateString()}
                    </span>
                    <span className="capitalize">{notice.category}</span>
                  </div>
                </div>
              </div>

              <p className="text-gray-600 mb-4 whitespace-pre-wrap">
                {notice.content}
              </p>

              {notice.publisher && (
                <p className="text-sm text-gray-500">
                  Published by: {notice.publisher.first_name}{' '}
                  {notice.publisher.last_name}
                </p>
              )}
            </div>
          ))}
        </div>
      )}
    </div>
  );
};