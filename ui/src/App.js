import React from 'react';
import {
  BrowserRouter as Router,
  Routes,
  Route,
  useNavigate
} from 'react-router-dom';
import { AuthProvider, AccountsProvider } from './context';
import MainLayout from './components/MainLayout';
import Auth from './components/Auth';
import { Home, About, Contact, Accounts, Account } from './pages';
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
                <AccountsProvider>
                  <MainLayout>
                    <Routes>
                      <Route path="/" element={<Home />} />
                      <Route path="/home" element={<Home />} />
                      <Route path="/about" element={<About />} />
                      <Route path="/contact" element={<Contact />} />
                      <Route path="/accounts" element={<Accounts />} />
                      <Route path="/account/:id" element={<Account />} />
                    </Routes>
                  </MainLayout>
                </AccountsProvider>
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
