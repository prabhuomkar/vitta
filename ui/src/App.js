// App.js
import React from 'react';
import {
  BrowserRouter as Router,
  Routes,
  Route,
  useNavigate
} from 'react-router-dom';
import { AuthProvider } from './context/AuthContext';
import MainLayout from './components/MainLayout';
import Auth from './components/Auth';
import Home from './pages/Home';
import About from './pages/About';
import Contact from './pages/Contact';
import PrivateRoute from './components/common/PrivateRoute';

const AppRoutes = () => {
  const navigate = useNavigate();

  return (
    <AuthProvider navigate={navigate}>
      <Routes>
        <Route path="/auth" element={<Auth />} />
        <Route
          path="*"
          element={
            <PrivateRoute
              element={
                <MainLayout>
                  <Routes>
                    <Route path="/" element={<Home />} />
                    <Route path="/home" element={<Home />} />
                    <Route path="/about" element={<About />} />
                    <Route path="/contact" element={<Contact />} />
                  </Routes>
                </MainLayout>
              }
            />
          }
        />
      </Routes>
    </AuthProvider>
  );
};

const App = () => {
  return (
    <Router>
      <AppRoutes />
    </Router>
  );
};

export default App;
