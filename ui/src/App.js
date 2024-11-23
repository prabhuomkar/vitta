import React from 'react';
import {
  BrowserRouter as Router,
  Routes,
  Route,
  useNavigate
} from 'react-router-dom';
import {
  AuthProvider,
  AccountsProvider,
  TransactionsProvider,
  PayeesProvider,
  BudgetsProvider,
  GroupsProvider,
  CategoriesProvider
} from './context';
import MainLayout from './components/MainLayout';
import Auth from './components/Auth';
import { Budgets, Payees, Accounts, Account } from './pages';
import PrivateRoute from './components/common/PrivateRoute';
import './App.css';

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
                  <TransactionsProvider>
                    <PayeesProvider>
                      <BudgetsProvider>
                        <GroupsProvider>
                          <CategoriesProvider>
                            <MainLayout>
                              <Routes>
                                {/* <Route path="/" element={<Home />} />
                                <Route path="/home" element={<Home />} /> */}
                                <Route path="/budgets" element={<Budgets />} />
                                <Route path="/payees" element={<Payees />} />
                                <Route
                                  path="/accounts"
                                  element={<Accounts />}
                                />
                                <Route
                                  path="/account/:accountId"
                                  element={<Account />}
                                />
                              </Routes>
                            </MainLayout>
                          </CategoriesProvider>
                        </GroupsProvider>
                      </BudgetsProvider>
                    </PayeesProvider>
                  </TransactionsProvider>
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
