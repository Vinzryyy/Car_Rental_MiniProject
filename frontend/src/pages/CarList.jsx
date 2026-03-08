import React, { useState, useEffect } from 'react';
import api from '../api/axios';
import { Car, Search, Filter, Fuel, Users, MapPin } from 'lucide-react';
import { motion } from 'framer-motion';
import CarSkeleton from '../components/CarSkeleton';
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

  const container = {
    hidden: { opacity: 0 },
    show: {
      opacity: 1,
      transition: {
        staggerChildren: 0.1
      }
    }
  };

  const item = {
    hidden: { opacity: 0, y: 20 },
    show: { opacity: 1, y: 0 }
  };

  return (
    <div className="car-list-container">
      <section className="hero">
        <motion.div 
          initial={{ opacity: 0, y: 30 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.8 }}
        >
          <h1>Premium Car Rental</h1>
          <p>Drive your dream car today with our easy and affordable rental service.</p>
        </motion.div>
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

      {loading ? (
        <div className="car-grid">
          {[1, 2, 3, 4, 5, 6].map((n) => (
            <CarSkeleton key={n} />
          ))}
        </div>
      ) : (
        <motion.div 
          className="car-grid"
          variants={container}
          initial="hidden"
          animate="show"
        >
          {cars.map((car) => (
            <motion.div key={car.id} className="car-card" variants={item}>
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
            </motion.div>
          ))}
        </motion.div>
      )}

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
