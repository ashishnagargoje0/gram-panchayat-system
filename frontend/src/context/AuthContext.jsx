import React, { createContext, useContext } from 'react';

const AuthContext = createContext({
  login: async (email, password) => {
    // dummy login — replace with actual API call later
    console.log('login called', email, password);
    return Promise.resolve(true);
  },
});

export const AuthProvider = ({ children }) => {
  const login = async (email, password) => {
    // implement real API call to backend later
    return Promise.resolve(true);
  };

  return (
    <AuthContext.Provider value={{ login }}>
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => useContext(AuthContext);
