import React, { useState, useEffect } from 'react';
import api from '../api/axios';
import { Car, Search, Filter, Fuel, Users, MapPin, Loader } from 'lucide-react';
import './CarList.css';

const CarList = () => {
  const [cars, setCars] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [search, setSearch] = useState('');
  const [category, setCategory] = useState('');

  useEffect(() => {
    fetchCars();
  }, [category]);

  const fetchCars = async () => {
    try {
      setLoading(true);
      const params = {
        available: true,
        category: category || undefined,
        search: search || undefined
      };
      const response = await api.get('/cars', { params });
      setCars(response.data.data.cars);
    } catch (err) {
      setError('Failed to fetch cars');
    } finally {
      setLoading(false);
    }
  };

  const handleSearch = (e) => {
    e.preventDefault();
    fetchCars();
  };

  if (loading && cars.length === 0) {
    return (
      <div className="loading-container">
        <Loader className="spinner" size={48} />
        <p>Loading premium cars...</p>
      </div>
    );
  }

  return (
    <div className="car-list-container">
      <section className="hero">
        <h1>Premium Car Rental</h1>
        <p>Drive your dream car today with our easy and affordable rental service.</p>
      </section>

      <div className="filter-bar">
        <form onSubmit={handleSearch} className="search-box">
          <Search size={20} />
          <input
            type="text"
            placeholder="Search cars by name or spec..."
            value={search}
            onChange={(e) => setSearch(e.target.value)}
          />
          <button type="submit">Search</button>
        </form>

        <div className="category-filters">
          <button 
            className={category === '' ? 'active' : ''} 
            onClick={() => setCategory('')}
          >
            All
          </button>
          <button 
            className={category === 'SUV' ? 'active' : ''} 
            onClick={() => setCategory('SUV')}
          >
            SUV
          </button>
          <button 
            className={category === 'Sedan' ? 'active' : ''} 
            onClick={() => setCategory('Sedan')}
          >
            Sedan
          </button>
          <button 
            className={category === 'Luxury' ? 'active' : ''} 
            onClick={() => setCategory('Luxury')}
          >
            Luxury
          </button>
        </div>
      </div>

      {error && <div className="error-message">{error}</div>}

      <div className="car-grid">
        {cars.map((car) => (
          <div key={car.id} className="car-card">
            <div className="car-image">
              <img src={car.image_url || 'https://images.unsplash.com/photo-1503376780353-7e6692767b70?auto=format&fit=crop&q=80&w=800'} alt={car.name} />
              <div className="car-category">{car.category}</div>
            </div>
            <div className="car-info">
              <h3>{car.name}</h3>
              <div className="car-specs">
                <span><MapPin size={14} /> Jakarta</span>
                <span><Fuel size={14} /> Hybrid</span>
              </div>
              <p className="car-description">{car.description}</p>
              <div className="car-footer">
                <div className="car-price">
                  <span className="amount">IDR {car.rental_costs?.toLocaleString()}</span>
                  <span className="unit">/day</span>
                </div>
                <button className="rent-btn" onClick={() => window.location.href = `/cars/${car.id}`}>
                  View Details
                </button>
              </div>
            </div>
          </div>
        ))}
      </div>

      {cars.length === 0 && !loading && (
        <div className="no-results">
          <Car size={64} />
          <h3>No cars found</h3>
          <p>Try adjusting your search or filters.</p>
        </div>
      )}
    </div>
  );
};

export default CarList;
