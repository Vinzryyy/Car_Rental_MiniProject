import React, { useState, useEffect } from 'react';
import api from '../api/axios';
import { useAuth } from '../context/AuthContext';
import { Wallet, CreditCard, Clock, CheckCircle, AlertCircle, ExternalLink, Plus } from 'lucide-react';
import './Dashboard.css';

const Dashboard = () => {
  const { user, refreshUserData } = useAuth();
  const [rentals, setRentals] = useState([]);
  const [topups, setTopups] = useState([]);
  const [loading, setLoading] = useState(true);
  const [topUpAmount, setTopUpAmount] = useState('');
  const [isTopUpModalOpen, setIsTopUpModalOpen] = useState(false);

  useEffect(() => {
    fetchDashboardData();
  }, []);

  const fetchDashboardData = async () => {
    try {
      setLoading(true);
      const [rentalsRes, topupsRes] = await Promise.all([
        api.get('/rentals/my'),
        api.get('/topup/history')
      ]);
      setRentals(rentalsRes.data.data || []);
      setTopups(topupsRes.data.data || []);
      await refreshUserData();
    } catch (err) {
      console.error('Failed to fetch dashboard data', err);
    } finally {
      setLoading(false);
    }
  };

  const handleTopUp = async (e) => {
    e.preventDefault();
    try {
      const response = await api.post('/topup', { amount: parseFloat(topUpAmount) });
      const { payment_url } = response.data.data;
      if (payment_url) {
        window.open(payment_url, '_blank');
        setIsTopUpModalOpen(false);
        setTopUpAmount('');
        // Alert user to refresh after payment
        alert('Payment link opened in new tab. Please refresh dashboard after completing payment.');
      }
    } catch (err) {
      alert('Failed to initiate top-up');
    }
  };

  const getStatusBadge = (status) => {
    const badges = {
      pending: <span className="badge pending"><Clock size={12} /> Pending</span>,
      active: <span className="badge active"><CheckCircle size={12} /> Active</span>,
      completed: <span className="badge completed"><CheckCircle size={12} /> Completed</span>,
      cancelled: <span className="badge cancelled"><AlertCircle size={12} /> Cancelled</span>,
      overdue: <span className="badge overdue"><AlertCircle size={12} /> Overdue</span>
    };
    return badges[status] || <span>{status}</span>;
  };

  return (
    <div className="dashboard-container">
      <div className="dashboard-header">
        <h1>Welcome, {user?.email}</h1>
        <p>Manage your rentals and balance here.</p>
      </div>

      <div className="stats-grid">
        <div className="stat-card balance">
          <div className="stat-icon"><Wallet size={24} /></div>
          <div className="stat-content">
            <span className="stat-label">Available Balance</span>
            <h2 className="stat-value">IDR {user?.deposit_amount?.toLocaleString()}</h2>
          </div>
          <button className="add-funds-btn" onClick={() => setIsTopUpModalOpen(true)}>
            <Plus size={18} /> Top Up
          </button>
        </div>

        <div className="stat-card">
          <div className="stat-icon"><Clock size={24} /></div>
          <div className="stat-content">
            <span className="stat-label">Active Rentals</span>
            <h2 className="stat-value">{rentals.filter(r => r.status === 'active').length}</h2>
          </div>
        </div>
      </div>

      <div className="dashboard-content">
        <div className="section">
          <div className="section-header">
            <h2><Clock size={20} /> Recent Rentals</h2>
          </div>
          <div className="table-container">
            <table className="data-table">
              <thead>
                <tr>
                  <th>Car</th>
                  <th>Date</th>
                  <th>Cost</th>
                  <th>Status</th>
                  <th>Action</th>
                </tr>
              </thead>
              <tbody>
                {rentals.length > 0 ? rentals.map(rental => (
                  <tr key={rental.id}>
                    <td><strong>{rental.car_name}</strong></td>
                    <td>{new Date(rental.rental_date).toLocaleDateString()}</td>
                    <td>IDR {rental.total_cost?.toLocaleString()}</td>
                    <td>{getStatusBadge(rental.status)}</td>
                    <td>
                      {rental.status === 'pending' && rental.payment_url && (
                        <a href={rental.payment_url} target="_blank" rel="noreferrer" className="pay-link">
                          Pay Now <ExternalLink size={14} />
                        </a>
                      )}
                    </td>
                  </tr>
                )) : (
                  <tr><td colSpan="5" className="empty-row">No rentals found</td></tr>
                )}
              </tbody>
            </table>
          </div>
        </div>

        <div className="section">
          <div className="section-header">
            <h2><History size={20} /> Top Up History</h2>
          </div>
          <div className="table-container">
            <table className="data-table">
              <thead>
                <tr>
                  <th>Date</th>
                  <th>Amount</th>
                  <th>Status</th>
                </tr>
              </thead>
              <tbody>
                {topups.length > 0 ? topups.map(topup => (
                  <tr key={topup.id}>
                    <td>{new Date(topup.created_at).toLocaleDateString()}</td>
                    <td>IDR {topup.amount?.toLocaleString()}</td>
                    <td>{getStatusBadge(topup.status)}</td>
                  </tr>
                )) : (
                  <tr><td colSpan="3" className="empty-row">No top-up history</td></tr>
                )}
              </tbody>
            </table>
          </div>
        </div>
      </div>

      {isTopUpModalOpen && (
        <div className="modal-overlay">
          <div className="modal">
            <h3>Top Up Balance</h3>
            <form onSubmit={handleTopUp}>
              <div className="form-group">
                <label>Amount (IDR)</label>
                <input
                  type="number"
                  placeholder="Min. 10,000"
                  min="10000"
                  value={topUpAmount}
                  onChange={(e) => setTopUpAmount(e.target.value)}
                  required
                />
              </div>
              <div className="modal-actions">
                <button type="button" className="cancel-btn" onClick={() => setIsTopUpModalOpen(false)}>Cancel</button>
                <button type="submit" className="confirm-btn">Proceed to Payment</button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
};

// Simple import icon for Top Up History header
const History = ({ size }) => <Clock size={size} />;

export default Dashboard;
