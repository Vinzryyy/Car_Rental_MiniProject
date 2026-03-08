import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { AuthProvider, useAuth } from './context/AuthContext';
import Navbar from './components/Navbar';
import CarList from './pages/CarList';
import Login from './pages/Login';
import Register from './pages/Register';
import Dashboard from './pages/Dashboard';
import CarDetail from './pages/CarDetail';
import AdminDashboard from './pages/AdminDashboard';
import { Toaster } from 'sonner';
import './App.css';

// Protected Route component
const ProtectedRoute = ({ children }) => {
  const { user, loading } = useAuth();
  
  if (loading) return (
    <div className="min-h-screen flex items-center justify-center bg-dark">
      <div className="w-12 h-12 border-4 border-primary border-t-transparent rounded-full animate-spin"></div>
    </div>
  );
  
  if (!user) return <Navigate to="/login" />;
  
  return children;
};

// Admin Route component
const AdminRoute = ({ children }) => {
  const { user, loading } = useAuth();
  
  if (loading) return (
    <div className="min-h-screen flex items-center justify-center bg-dark">
      <div className="w-12 h-12 border-4 border-primary border-t-transparent rounded-full animate-spin"></div>
    </div>
  );
  
  if (!user || user.role !== 'admin') return <Navigate to="/" />;
  
  return children;
};

function AppContent() {
  return (
    <div className="flex flex-col min-h-screen bg-dark transition-colors duration-500">
      <Toaster position="top-right" richColors theme="dark" closeButton />
      <Navbar />
      <main className="flex-1">
        <Routes>
          <Route path="/" element={<CarList />} />
          <Route path="/cars/:id" element={<CarDetail />} />
          <Route path="/login" element={<Login />} />
          <Route path="/register" element={<Register />} />
          <Route 
            path="/dashboard" 
            element={
              <ProtectedRoute>
                <Dashboard />
              </ProtectedRoute>
            } 
          />
          <Route 
            path="/admin" 
            element={
              <AdminRoute>
                <AdminDashboard />
              </AdminRoute>
            } 
          />
        </Routes>
      </main>
      <footer className="bg-dark-card border-t border-white/5 py-12 px-4 sm:px-6 lg:px-8 text-center mt-20">
        <div className="max-w-7xl mx-auto">
          <p className="text-gray-500 text-sm font-bold tracking-widest uppercase mb-2">&copy; 2026 Vinz Rental Car</p>
          <p className="text-gray-600 text-xs italic">Premium Car Rental Solution • Experience Excellence on Every Mile</p>
        </div>
      </footer>
    </div>
  );
}

function App() {
  return (
    <AuthProvider>
      <Router>
        <AppContent />
      </Router>
    </AuthProvider>
  );
}

export default App;
