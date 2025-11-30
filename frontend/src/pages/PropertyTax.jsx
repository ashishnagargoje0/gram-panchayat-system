import React, { useState, useEffect } from 'react';
import { propertyService } from '../services/authService';
import { Home, DollarSign, CreditCard } from 'lucide-react';

const PropertyTax = () => {
  const [properties, setProperties] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchProperties();
  }, []);

  const fetchProperties = async () => {
    try {
      const response = await propertyService.getProperties({});
      setProperties(response.data || []);
    } catch (error) {
      console.error('Failed to fetch properties:', error);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="space-y-6">
      <div className="bg-gradient-to-r from-green-500 to-teal-600 rounded-xl p-6 text-white">
        <h2 className="text-2xl font-bold mb-2">Property Tax Management</h2>
        <p className="opacity-90">View and pay your property taxes online</p>
      </div>

      {loading ? (
        <div className="flex justify-center py-12">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
        </div>
      ) : properties.length === 0 ? (
        <div className="bg-white rounded-xl shadow-md p-12 text-center">
          <Home className="mx-auto mb-4 text-gray-400" size={48} />
          <h3 className="text-xl font-bold text-gray-800 mb-2">
            No Properties Found
          </h3>
          <p className="text-gray-600 mb-4">
            You don't have any registered properties yet
          </p>
          <button className="bg-blue-600 text-white px-6 py-3 rounded-lg font-medium hover:bg-blue-700 transition">
            Register Property
          </button>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          {properties.map((property) => (
            <div
              key={property.id}
              className="bg-white rounded-xl shadow-md p-6 hover:shadow-xl transition"
            >
              <div className="flex items-start justify-between mb-4">
                <div>
                  <h3 className="text-lg font-bold text-gray-800 mb-1">
                    {property.address}
                  </h3>
                  <p className="text-sm text-gray-600">
                    Property No: {property.property_number}
                  </p>
                </div>
                <span className="px-3 py-1 bg-green-100 text-green-700 text-sm rounded-full font-medium">
                  Active
                </span>
              </div>

              <div className="space-y-2 mb-4">
                <div className="flex justify-between text-sm">
                  <span className="text-gray-600">Type:</span>
                  <span className="font-medium text-gray-800 capitalize">
                    {property.property_type}
                  </span>
                </div>
                <div className="flex justify-between text-sm">
                  <span className="text-gray-600">Area:</span>
                  <span className="font-medium text-gray-800">
                    {property.area} sq ft
                  </span>
                </div>
                <div className="flex justify-between text-sm">
                  <span className="text-gray-600">Annual Tax:</span>
                  <span className="font-medium text-gray-800">
                    ₹{property.annual_tax_amount?.toLocaleString()}
                  </span>
                </div>
              </div>

              <button className="w-full bg-gradient-to-r from-green-600 to-teal-600 text-white py-3 rounded-lg font-medium hover:shadow-lg transition flex items-center justify-center gap-2">
                <CreditCard size={20} /> Pay Tax
              </button>
            </div>
          ))}
        </div>
      )}
    </div>
  );
};