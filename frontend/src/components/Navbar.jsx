import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { Car, User, LogOut, LayoutDashboard } from 'lucide-react';

const Navbar = () => {
  const { user, logout } = useAuth();
  const navigate = useNavigate();

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  return (
    <nav className="sticky top-0 z-50 bg-dark-card/80 backdrop-blur-xl border-b border-white/5 shadow-xl">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between items-center h-16">
          <Link to="/" className="flex items-center gap-2 text-primary font-bold text-xl hover:scale-105 transition-transform">
            <Car size={32} />
            <span>Vinz Rental</span>
          </Link>

          <div className="flex items-center gap-6">
            <Link to="/" className="text-gray-300 hover:text-primary hover:bg-primary/5 px-3 py-2 rounded-lg font-medium transition-all">
              Cars
            </Link>
            
            {user ? (
              <div className="flex items-center gap-4">
                <Link to="/dashboard" className="flex items-center gap-2 text-gray-300 hover:text-primary hover:bg-primary/5 px-3 py-2 rounded-lg font-medium transition-all">
                  <LayoutDashboard size={18} />
                  <span className="hidden sm:inline">Dashboard</span>
                </Link>
                
                <div className="group relative flex items-center gap-2 bg-white/5 border border-white/5 px-4 py-1.5 rounded-full cursor-pointer hover:bg-white/10 transition-all">
                  <User size={18} className="text-primary" />
                  <span className="text-sm font-medium text-gray-300">{user.email}</span>
                  
                  <div className="absolute right-0 top-full mt-2 w-48 bg-dark-card border border-white/5 rounded-xl shadow-2xl opacity-0 translate-y-2 pointer-events-none group-hover:opacity-100 group-hover:translate-y-0 transition-all duration-200 p-1">
                    <button 
                      onClick={handleLogout} 
                      className="w-full flex items-center gap-2 px-4 py-2.5 text-sm text-red-400 hover:bg-red-500/10 rounded-lg transition-colors"
                    >
                      <LogOut size={16} />
                      <span>Logout</span>
                    </button>
                  </div>
                </div>
              </div>
            ) : (
              <div className="flex items-center gap-3">
                <Link to="/login" className="text-primary border border-primary/30 hover:border-primary hover:bg-primary/10 px-5 py-2 rounded-lg font-semibold transition-all text-sm">
                  Login
                </Link>
                <Link to="/register" className="bg-primary text-white shadow-lg shadow-primary/30 hover:bg-blue-600 hover:-translate-y-0.5 px-5 py-2 rounded-lg font-semibold transition-all text-sm">
                  Register
                </Link>
              </div>
            )}
          </div>
        </div>
      </div>
    </nav>
  );
};

export default Navbar;
