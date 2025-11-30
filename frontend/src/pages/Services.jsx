import React from 'react';
import { useNavigate } from 'react-router-dom';
import { FileText, Home, DollarSign, Heart } from 'lucide-react';

const Services = () => {
  const navigate = useNavigate();

  const services = [
    { type: 'birth', name: 'Birth Certificate', icon: Heart, color: 'blue' },
    { type: 'death', name: 'Death Certificate', icon: FileText, color: 'purple' },
    { type: 'income', name: 'Income Certificate', icon: DollarSign, color: 'green' },
    { type: 'caste', name: 'Caste Certificate', icon: FileText, color: 'yellow' },
    { type: 'residence', name: 'Residence Certificate', icon: Home, color: 'red' },
    { type: 'marriage', name: 'Marriage Certificate', icon: Heart, color: 'pink' },
  ];

  return (
    <div className="space-y-6">
      <div className="bg-gradient-to-r from-blue-500 to-purple-600 rounded-xl p-6 text-white">
        <h2 className="text-2xl font-bold mb-2">Apply for Certificates Online</h2>
        <p className="opacity-90">Get your certificates quickly and hassle-free</p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {services.map((service) => (
          <div
            key={service.type}
            onClick={() => navigate(`/services/apply/${service.type}`)}
            className="bg-white rounded-xl shadow-md p-6 hover:shadow-xl transition cursor-pointer transform hover:scale-105"
          >
            <div className={`w-16 h-16 rounded-full bg-${service.color}-100 flex items-center justify-center mb-4`}>
              <service.icon className={`text-${service.color}-600`} size={28} />
            </div>
            <h3 className="text-lg font-bold text-gray-800 mb-2">{service.name}</h3>
            <p className="text-sm text-gray-600 mb-4">
              Apply online and track your application status
            </p>
            <button className="text-blue-600 font-medium hover:underline">
              Apply Now →
            </button>
          </div>
        ))}
      </div>
    </div>
  );
};
