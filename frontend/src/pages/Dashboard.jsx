import { useState, useEffect, useCallback } from 'react';
import api from '../api/axios';
import { useAuth } from '../context/AuthContext';
import { Wallet, Clock, ExternalLink, Plus } from 'lucide-react';
import { toast } from 'sonner';

const Dashboard = () => {
  const { user, refreshUserData } = useAuth();
  const [rentals, setRentals] = useState([]);
  const [topups, setTopups] = useState([]);
  const [topUpAmount, setTopUpAmount] = useState('');
  const [isTopUpModalOpen, setIsTopUpModalOpen] = useState(false);

  const fetchDashboardData = useCallback(async () => {
    try {
      const [rentalsRes, topupsRes] = await Promise.all([
        api.get('/rentals/my'),
        api.get('/topup/history')
      ]);
      setRentals(rentalsRes.data.data || []);
      setTopups(topupsRes.data.data || []);
      await refreshUserData();
    } catch {
      toast.error('Failed to load dashboard data');
    }
  }, [refreshUserData]);

  useEffect(() => {
    fetchDashboardData();
  }, [fetchDashboardData]);

  const handleTopUp = async (e) => {
    e.preventDefault();
    try {
      const response = await api.post('/topup', { amount: parseFloat(topUpAmount) });
      const { payment_url } = response.data.data;
      if (payment_url) {
        window.open(payment_url, '_blank');
        setIsTopUpModalOpen(false);
        setTopUpAmount('');
        toast.success('Payment link opened! Please complete payment.');
      }
    } catch {
      toast.error('Failed to initiate top-up');
    }
  };

  const getStatusBadge = (status) => {
    const config = {
      pending: "bg-orange-500/10 text-orange-500 border-orange-500/20",
      active: "bg-emerald-500/10 text-emerald-500 border-emerald-500/20",
      completed: "bg-blue-500/10 text-blue-500 border-blue-500/20",
      cancelled: "bg-red-500/10 text-red-500 border-red-500/20",
      overdue: "bg-red-600/20 text-red-400 border-red-600/30"
    };
    
    return (
      <span className={`px-3 py-1 rounded-full text-[10px] font-black uppercase tracking-wider border ${config[status] || "bg-gray-500/10 text-gray-500"}`}>
        {status}
      </span>
    );
  };

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-10">
      <div className="mb-10">
        <h1 className="text-3xl font-extrabold text-white mb-2 tracking-tight">Dashboard</h1>
        <p className="text-gray-500 font-medium italic">Welcome back, {user?.email}</p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-12">
        <div className="bg-dark-card border border-primary/20 p-8 rounded-2xl relative overflow-hidden group glow-primary transition-all duration-300">
          <div className="absolute top-0 right-0 p-4 opacity-5 group-hover:scale-110 transition-transform">
            <Wallet size={120} />
          </div>
          <div className="relative z-10 flex flex-col h-full justify-between">
            <div>
              <span className="text-xs font-black text-primary uppercase tracking-[0.2em] mb-2 block">Available Balance</span>
              <h2 className="text-4xl font-black text-white mb-6 tracking-tight">IDR {user?.deposit_amount?.toLocaleString()}</h2>
            </div>
            <button 
              onClick={() => setIsTopUpModalOpen(true)}
              className="w-fit flex items-center gap-2 bg-primary text-white px-6 py-2.5 rounded-xl font-bold hover:bg-blue-600 transition-all shadow-lg shadow-primary/20 active:scale-95"
            >
              <Plus size={18} /> <span>Top Up Funds</span>
            </button>
          </div>
        </div>

        <div className="bg-dark-card border border-white/5 p-8 rounded-2xl flex items-center gap-6 shadow-xl">
          <div className="w-16 h-16 bg-white/5 rounded-2xl flex items-center justify-center text-primary border border-white/5">
            <Clock size={28} />
          </div>
          <div>
            <span className="text-xs font-black text-gray-500 uppercase tracking-[0.2em] mb-1 block">Active Rentals</span>
            <h2 className="text-4xl font-black text-white tracking-tight">{rentals.filter(r => r.status === 'active').length}</h2>
          </div>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
        <div className="lg:col-span-2 bg-dark-card border border-white/5 rounded-2xl overflow-hidden shadow-2xl self-start">
          <div className="px-6 py-5 border-b border-white/5 bg-white/[0.02] flex items-center gap-3">
            <div className="text-primary bg-primary/10 p-2 rounded-lg"><Clock size={18} /></div>
            <h2 className="font-bold text-white tracking-tight">Recent Rental Activity</h2>
          </div>
          <div className="overflow-x-auto">
            <table className="w-full text-left">
              <thead>
                <tr className="bg-white/[0.01]">
                  <th className="px-6 py-4 text-[10px] font-black text-gray-500 uppercase tracking-widest">Vehicle</th>
                  <th className="px-6 py-4 text-[10px] font-black text-gray-500 uppercase tracking-widest">Date</th>
                  <th className="px-6 py-4 text-[10px] font-black text-gray-500 uppercase tracking-widest">Cost</th>
                  <th className="px-6 py-4 text-[10px] font-black text-gray-500 uppercase tracking-widest">Status</th>
                  <th className="px-6 py-4 text-[10px] font-black text-gray-500 uppercase tracking-widest text-right">Action</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-white/5">
                {rentals.length > 0 ? rentals.map(rental => (
                  <tr key={rental.id} className="hover:bg-white/[0.02] transition-colors">
                    <td className="px-6 py-4"><span className="text-sm font-bold text-white tracking-tight">{rental.car_name}</span></td>
                    <td className="px-6 py-4 text-xs font-medium text-gray-500">{new Date(rental.rental_date).toLocaleDateString()}</td>
                    <td className="px-6 py-4 text-sm font-black text-primary">IDR {rental.total_cost?.toLocaleString()}</td>
                    <td className="px-6 py-4">{getStatusBadge(rental.status)}</td>
                    <td className="px-6 py-4 text-right">
                      {rental.status === 'pending' && rental.payment_url && (
                        <a href={rental.payment_url} target="_blank" rel="noreferrer" className="inline-flex items-center gap-1.5 text-xs font-black text-primary hover:underline bg-primary/10 px-3 py-1.5 rounded-lg transition-all active:scale-95 uppercase tracking-wider">
                          PAY <ExternalLink size={12} />
                        </a>
                      )}
                    </td>
                  </tr>
                )) : (
                  <tr><td colSpan="5" className="px-6 py-12 text-center text-gray-600 font-medium italic">No recent rental transactions</td></tr>
                )}
              </tbody>
            </table>
          </div>
        </div>

        <div className="bg-dark-card border border-white/5 rounded-2xl overflow-hidden shadow-2xl self-start">
          <div className="px-6 py-5 border-b border-white/5 bg-white/[0.02] flex items-center gap-3">
            <div className="text-primary bg-primary/10 p-2 rounded-lg"><Clock size={18} /></div>
            <h2 className="font-bold text-white tracking-tight">Top Up History</h2>
          </div>
          <div className="overflow-x-auto">
            <table className="w-full text-left">
              <thead>
                <tr className="bg-white/[0.01]">
                  <th className="px-6 py-4 text-[10px] font-black text-gray-500 uppercase tracking-widest">Date</th>
                  <th className="px-6 py-4 text-[10px] font-black text-gray-500 uppercase tracking-widest">Amount</th>
                  <th className="px-6 py-4 text-[10px] font-black text-gray-500 uppercase tracking-widest">Status</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-white/5">
                {topups.length > 0 ? topups.map(topup => (
                  <tr key={topup.id} className="hover:bg-white/[0.02] transition-colors">
                    <td className="px-6 py-4 text-xs font-medium text-gray-500">{new Date(topup.created_at).toLocaleDateString()}</td>
                    <td className="px-6 py-4 text-sm font-black text-white">IDR {topup.amount?.toLocaleString()}</td>
                    <td className="px-6 py-4">{getStatusBadge(topup.status)}</td>
                  </tr>
                )) : (
                  <tr><td colSpan="3" className="px-6 py-12 text-center text-gray-600 font-medium italic">No top-up records</td></tr>
                )}
              </tbody>
            </table>
          </div>
        </div>
      </div>

      {isTopUpModalOpen && (
        <div className="fixed inset-0 z-[100] flex items-center justify-center p-4">
          <div className="absolute inset-0 bg-black/80 backdrop-blur-md" onClick={() => setIsTopUpModalOpen(false)}></div>
          <div className="relative bg-dark-card border border-white/10 p-8 rounded-3xl w-full max-w-sm shadow-2xl animate-in zoom-in-95 duration-200">
            <h3 className="text-2xl font-black text-white mb-2 tracking-tight">Top Up Balance</h3>
            <p className="text-gray-500 text-sm font-medium mb-8 uppercase tracking-widest">Add funds to your account</p>
            
            <form onSubmit={handleTopUp} className="space-y-6">
              <div className="space-y-2">
                <label className="text-[10px] font-black text-gray-400 uppercase tracking-[0.2em]">Amount (IDR)</label>
                <input
                  type="number"
                  placeholder="Minimum 10,000"
                  min="10000"
                  className="w-full bg-black/40 border border-white/10 rounded-2xl py-4 px-6 text-xl font-black text-primary placeholder:text-gray-700 focus:outline-none focus:border-primary transition-all"
                  value={topUpAmount}
                  onChange={(e) => setTopUpAmount(e.target.value)}
                  required
                />
              </div>
              
              <div className="flex gap-3 pt-4">
                <button 
                  type="button" 
                  className="flex-1 border border-white/10 text-white font-bold py-3 rounded-xl hover:bg-white/5 transition-all active:scale-95" 
                  onClick={() => setIsTopUpModalOpen(false)}
                >
                  Cancel
                </button>
                <button 
                  type="submit" 
                  className="flex-1 bg-primary text-white font-black py-3 rounded-xl hover:bg-blue-600 shadow-lg shadow-primary/20 transition-all active:scale-95 uppercase tracking-widest text-xs"
                >
                  Continue
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
};

export default Dashboard;
