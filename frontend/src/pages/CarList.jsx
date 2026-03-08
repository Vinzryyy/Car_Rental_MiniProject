import React, { useState, useEffect } from 'react';
import api from '../api/axios';
import { Car, Search, Fuel, MapPin } from 'lucide-react';
import { motion } from 'framer-motion';
import CarSkeleton from '../components/CarSkeleton';

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
      transition: { staggerChildren: 0.1 }
    }
  };

  const item = {
    hidden: { opacity: 0, y: 20 },
    show: { opacity: 1, y: 0 }
  };

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
      {/* Hero Section */}
      <section className="relative overflow-hidden bg-dark-card rounded-3xl mb-12 p-12 text-center border border-white/5">
        <div className="absolute inset-0 bg-[url('https://images.unsplash.com/photo-1492144534655-ae79c964c9d7?auto=format&fit=crop&q=80&w=1200')] bg-cover bg-center opacity-20"></div>
        <div className="absolute inset-0 bg-gradient-to-b from-dark-card/50 to-dark-card"></div>
        
        <motion.div 
          initial={{ opacity: 0, y: 30 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.8 }}
          className="relative z-10"
        >
          <h1 className="text-4xl md:text-6xl font-extrabold text-white mb-6 tracking-tight">
            Premium Car <span className="text-primary">Rental</span>
          </h1>
          <p className="text-lg md:text-xl text-gray-400 max-w-2xl mx-auto font-medium">
            Drive your dream car today with our easy and affordable rental service.
          </p>
        </motion.div>
      </section>

      {/* Filter Bar */}
      <div className="flex flex-col md:flex-row justify-between items-center gap-6 mb-12">
        <form onSubmit={handleSearch} className="flex-1 w-full flex items-center bg-dark-card border border-white/10 rounded-2xl px-4 py-2 focus-within:border-primary transition-all shadow-lg">
          <Search size={20} className="text-gray-500" />
          <input
            type="text"
            placeholder="Search cars by name or spec..."
            className="bg-transparent border-none w-full p-3 text-white focus:outline-none placeholder:text-gray-600"
            value={search}
            onChange={(e) => setSearch(e.target.value)}
          />
          <button type="submit" className="bg-primary text-white px-6 py-2 rounded-xl font-bold hover:bg-blue-600 transition-colors hidden sm:block">
            Search
          </button>
        </form>

        <div className="flex gap-2 p-1.5 bg-dark-card border border-white/5 rounded-2xl">
          {['', 'SUV', 'Sedan', 'Luxury'].map((cat) => (
            <button 
              key={cat}
              className={`px-6 py-2 rounded-xl text-sm font-bold transition-all ${
                category === cat 
                ? 'bg-primary text-white shadow-lg shadow-primary/20' 
                : 'text-gray-400 hover:text-white hover:bg-white/5'
              }`} 
              onClick={() => setCategory(cat)}
            >
              {cat || 'All'}
            </button>
          ))}
        </div>
      </div>

      {error && <div className="bg-red-500/10 border border-red-500/20 text-red-400 p-4 rounded-xl mb-8 text-center">{error}</div>}

      {loading ? (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
          {[1, 2, 3, 4, 5, 6].map((n) => (
            <CarSkeleton key={n} />
          ))}
        </div>
      ) : (
        <motion.div 
          className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8"
          variants={container}
          initial="hidden"
          animate="show"
        >
          {cars.map((car) => (
            <motion.div key={car.id} className="group bg-dark-card rounded-2xl overflow-hidden border border-white/5 hover:border-primary/50 shadow-xl transition-all duration-300" variants={item}>
              <div className="relative h-52 overflow-hidden">
                <img 
                  src={car.image_url || 'https://images.unsplash.com/photo-1503376780353-7e6692767b70?auto=format&fit=crop&q=80&w=800'} 
                  alt={car.name} 
                  className="w-full h-full object-cover group-hover:scale-110 transition-transform duration-500"
                />
                <div className="absolute top-4 right-4 bg-primary/90 text-white text-[10px] uppercase tracking-widest font-black px-3 py-1.5 rounded-lg backdrop-blur-md">
                  {car.category}
                </div>
              </div>
              
              <div className="p-6">
                <h3 className="text-xl font-bold text-white mb-2">{car.name}</h3>
                
                <div className="flex gap-4 mb-4 text-gray-500 text-xs font-semibold uppercase tracking-wider">
                  <span className="flex items-center gap-1.5"><MapPin size={14} className="text-primary" /> Jakarta</span>
                  <span className="flex items-center gap-1.5"><Fuel size={14} className="text-primary" /> Hybrid</span>
                </div>
                
                <p className="text-gray-400 text-sm line-clamp-2 mb-6 h-10 font-medium">
                  {car.description || "Premium driving experience with top-tier comfort and modern features."}
                </p>
                
                <div className="flex justify-between items-center pt-5 border-t border-white/5">
                  <div>
                    <span className="text-2xl font-black text-primary">IDR {car.rental_costs?.toLocaleString()}</span>
                    <span className="text-gray-600 text-xs font-bold block">/ DAY</span>
                  </div>
                  <button 
                    onClick={() => window.location.href = `/cars/${car.id}`}
                    className="bg-white/5 border border-white/10 text-white hover:bg-primary hover:border-primary px-5 py-2.5 rounded-xl text-sm font-bold transition-all shadow-lg active:scale-95"
                  >
                    Details
                  </button>
                </div>
              </div>
            </motion.div>
          ))}
        </motion.div>
      )}

      {cars.length === 0 && !loading && (
        <div className="text-center py-24 opacity-50">
          <Car size={64} className="mx-auto mb-4" />
          <h3 className="text-2xl font-bold">No cars found</h3>
          <p className="font-medium mt-2">Try adjusting your search or filters.</p>
        </div>
      )}
    </div>
  );
};

export default CarList;
