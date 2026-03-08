import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import api from '../api/axios';
import { useAuth } from '../context/AuthContext';
import { Calendar, Shield, Zap, Info, Loader, ArrowLeft, CheckCircle } from 'lucide-react';
import './CarDetail.css';

const CarDetail = () => {
  const { id } = useParams();
  const navigate = useNavigate();
  const { user } = useAuth();
  const [car, setCar] = useState(null);
  const [days, setDays] = useState(1);
  const [loading, setLoading] = useState(true);
  const [renting, setRenting] = useState(false);
  const [error, setError] = useState('');

  useEffect(() => {
    fetchCarDetails();
  }, [id]);

  const fetchCarDetails = async () => {
    try {
      const response = await api.get(`/cars/${id}`);
      setCar(response.data.data);
    } catch (err) {
      setError('Car not found');
    } finally {
      setLoading(false);
    }
  };

  const handleRent = async () => {
    if (!user) {
      navigate('/login');
      return;
    }

    try {
      setRenting(true);
      const response = await api.post('/rentals', {
        car_id: id,
        rental_days: parseInt(days)
      });
      
      const { payment_url } = response.data.data;
      if (payment_url) {
        window.location.href = payment_url;
      } else {
        // If no payment URL, maybe it was paid by deposit
        navigate('/dashboard');
      }
    } catch (err) {
      setError(err.response?.data?.message || 'Failed to rent car');
      setRenting(false);
    }
  };

  if (loading) return <div className="loading-container"><Loader className="spinner" /></div>;
  if (!car) return <div className="error-container"><h2>Car not found</h2><button onClick={() => navigate('/')}>Go Back</button></div>;

  return (
    <div className="car-detail-container">
      <button className="back-btn" onClick={() => navigate('/')}>
        <ArrowLeft size={20} /> Back to Gallery
      </button>

      <div className="detail-grid">
        <div className="detail-image">
          <img src={car.image_url || 'https://images.unsplash.com/photo-1503376780353-7e6692767b70?auto=format&fit=crop&q=80&w=800'} alt={car.name} />
        </div>

        <div className="detail-content">
          <div className="detail-header">
            <span className="category-pill">{car.category}</span>
            <h1>{car.name}</h1>
            <p className="price-tag">IDR {car.rental_costs?.toLocaleString()} <span>/ day</span></p>
          </div>

          <div className="features-grid">
            <div className="feature"><Zap size={18} /> Instant Booking</div>
            <div className="feature"><Shield size={18} /> Insurance Included</div>
            <div className="feature"><CheckCircle size={18} /> Clean & Sanitized</div>
            <div className="feature"><Info size={18} /> 24/7 Support</div>
          </div>

          <div className="description">
            <h3>Description</h3>
            <p>{car.description || 'No description available for this premium vehicle.'}</p>
          </div>

          <div className="rental-options">
            <div className="days-picker">
              <label>Rental Duration (Days)</label>
              <div className="number-input">
                <button onClick={() => setDays(Math.max(1, days - 1))}>-</button>
                <input type="number" value={days} readOnly />
                <button onClick={() => setDays(days + 1)}>+</button>
              </div>
            </div>

            <div className="total-summary">
              <div className="summary-row">
                <span>Daily Rate</span>
                <span>IDR {car.rental_costs?.toLocaleString()}</span>
              </div>
              <div className="summary-row">
                <span>Duration</span>
                <span>{days} days</span>
              </div>
              <div className="summary-row total">
                <span>Total Amount</span>
                <span>IDR {(car.rental_costs * days).toLocaleString()}</span>
              </div>
            </div>

            <button 
              className="confirm-rent-btn" 
              onClick={handleRent}
              disabled={renting || !car.availability}
            >
              {renting ? 'Processing...' : !car.availability ? 'Not Available' : 'Rent Now'}
            </button>
            
            {!user && <p className="login-hint">Please login to complete your booking.</p>}
          </div>
        </div>
      </div>
    </div>
  );
};

export default CarDetail;
