import { useState, useEffect, useCallback, useMemo } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import api from '../api/axios';
import { useAuth } from '../context/AuthContext';
import { Shield, Zap, Info, Loader, ArrowLeft, CheckCircle, Calendar as CalendarIcon } from 'lucide-react';
import { toast } from 'sonner';
import DatePicker from 'react-datepicker';
import { addDays, differenceInDays, format } from 'date-fns';
import "react-datepicker/dist/react-datepicker.css";

const CarDetail = () => {
  const { id } = useParams();
  const navigate = useNavigate();
  const { user } = useAuth();
  
  const [car, setCar] = useState(null);
  const [loading, setLoading] = useState(true);
  const [renting, setRenting] = useState(false);
  
  const [startDate, setStartDate] = useState(new Date());
  const [endDate, setEndDate] = useState(addDays(new Date(), 1));

  const rentalDays = useMemo(() => {
    const days = differenceInDays(endDate, startDate);
    return days > 0 ? days : 1;
  }, [startDate, endDate]);

  const fetchCarDetails = useCallback(async () => {
    try {
      setLoading(true);
      const response = await api.get(`/cars/${id}`);
      setCar(response.data.data);
    } catch {
      toast.error('Car not found');
    } finally {
      setLoading(false);
    }
  }, [id]);

  useEffect(() => {
    fetchCarDetails();
  }, [fetchCarDetails]);

  const handleRent = async () => {
    if (!user) {
      navigate('/login');
      return;
    }

    try {
      setRenting(true);
      const response = await api.post('/rentals', {
        car_id: id,
        start_date: format(startDate, 'yyyy-MM-dd'),
        end_date: format(endDate, 'yyyy-MM-dd')
      });
      
      const { payment_url } = response.data.data;
      if (payment_url) {
        toast.success('Redirecting to payment gateway...');
        window.location.href = payment_url;
      } else {
        toast.success('Car rented successfully!');
        navigate('/dashboard');
      }
    } catch (err) {
      toast.error(err.response?.data?.message || 'Failed to rent car');
      setRenting(false);
    }
  };

  if (loading) return <div className="min-h-screen flex items-center justify-center"><Loader className="animate-spin text-primary" size={48} /></div>;
  if (!car) return (
    <div className="min-h-screen flex flex-col items-center justify-center gap-4 text-center p-6">
      <h2 className="text-3xl font-black text-white tracking-tight">Car Not Found</h2>
      <button onClick={() => navigate('/')} className="bg-white/5 border border-white/10 px-8 py-3 rounded-2xl font-bold hover:bg-primary transition-all shadow-xl">Go Back Home</button>
    </div>
  );

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-10">
      <button className="flex items-center gap-2 text-gray-500 hover:text-primary font-bold transition-colors mb-10 group" onClick={() => navigate('/')}>
        <ArrowLeft size={20} className="group-hover:-translate-x-1 transition-transform" /> 
        <span>Back to Gallery</span>
      </button>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-12 bg-dark-card rounded-[2.5rem] border border-white/5 overflow-hidden shadow-2xl">
        <div className="relative h-[400px] lg:h-full group">
          <img 
            src={car.image_url || 'https://images.unsplash.com/photo-1503376780353-7e6692767b70?auto=format&fit=crop&q=80&w=800'} 
            alt={car.name} 
            className="w-full h-full object-cover transition-transform duration-700 group-hover:scale-105"
          />
          <div className="absolute inset-0 bg-gradient-to-t from-dark-card via-transparent to-transparent lg:hidden"></div>
        </div>

        <div className="p-8 lg:p-16 flex flex-col gap-10">
          <header>
            <span className="bg-primary/10 text-primary px-4 py-1.5 rounded-xl text-[10px] font-black uppercase tracking-widest border border-primary/20">{car.category}</span>
            <h1 className="text-4xl lg:text-5xl font-black text-white mt-6 mb-4 tracking-tight leading-tight">{car.name}</h1>
            <div className="flex items-baseline gap-2">
              <span className="text-3xl font-black text-primary">IDR {car.rental_costs?.toLocaleString()}</span>
              <span className="text-gray-600 font-bold uppercase text-xs tracking-wider">/ PER DAY</span>
            </div>
          </header>

          <div className="grid grid-cols-2 gap-6">
            <div className="flex items-center gap-3 text-gray-400 font-medium">
              <div className="p-2 bg-emerald-500/10 text-emerald-500 rounded-lg"><Zap size={18} /></div>
              <span className="text-sm">Instant Confirm</span>
            </div>
            <div className="flex items-center gap-3 text-gray-400 font-medium">
              <div className="p-2 bg-emerald-500/10 text-emerald-500 rounded-lg"><Shield size={18} /></div>
              <span className="text-sm">Full Insurance</span>
            </div>
          </div>

          <div className="space-y-4">
            <h3 className="text-lg font-bold text-white tracking-tight uppercase text-xs text-gray-500">About this vehicle</h3>
            <p className="text-gray-400 leading-relaxed font-medium">
              {car.description || "Indulge in the ultimate driving experience with top-tier comfort and modern features."}
            </p>
          </div>

          <div className="bg-black/30 p-8 rounded-3xl border border-white/5 space-y-8">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div className="space-y-2">
                <label className="text-[10px] font-black text-gray-500 uppercase tracking-widest flex items-center gap-2">
                  <CalendarIcon size={12} /> Start Date
                </label>
                <DatePicker
                  selected={startDate}
                  onChange={(date) => setStartDate(date)}
                  selectsStart
                  startDate={startDate}
                  endDate={endDate}
                  minDate={new Date()}
                  className="w-full bg-dark-card border border-white/10 rounded-xl py-3 px-4 text-white focus:outline-none focus:border-primary transition-all"
                />
              </div>
              <div className="space-y-2">
                <label className="text-[10px] font-black text-gray-500 uppercase tracking-widest flex items-center gap-2">
                  <CalendarIcon size={12} /> End Date
                </label>
                <DatePicker
                  selected={endDate}
                  onChange={(date) => setEndDate(date)}
                  selectsEnd
                  startDate={startDate}
                  endDate={endDate}
                  minDate={addDays(startDate, 1)}
                  className="w-full bg-dark-card border border-white/10 rounded-xl py-3 px-4 text-white focus:outline-none focus:border-primary transition-all"
                />
              </div>
            </div>

            <div className="space-y-3 pt-6 border-t border-white/5">
              <div className="flex justify-between items-center text-gray-500 text-sm font-medium">
                <span>Duration</span>
                <span className="text-white font-bold">{rentalDays} days</span>
              </div>
              <div className="flex justify-between items-center pt-4">
                <span className="text-xs font-black text-gray-400 uppercase tracking-widest">Total Amount</span>
                <span className="text-3xl font-black text-primary italic">IDR {(car.rental_costs * rentalDays).toLocaleString()}</span>
              </div>
            </div>

            <button 
              className="w-full bg-primary hover:bg-blue-600 text-white font-black py-5 rounded-2xl shadow-2xl shadow-primary/30 transition-all active:scale-[0.98] disabled:opacity-50 disabled:cursor-not-allowed uppercase tracking-[0.2em] text-sm"
              onClick={handleRent}
              disabled={renting || !car.availability}
            >
              {renting ? 'Processing...' : !car.availability ? 'Fully Booked' : 'Confirm Booking'}
            </button>
            
            {!user && <p className="text-center text-xs font-bold text-orange-500/80 uppercase tracking-wider italic">Authentication required to complete booking</p>}
          </div>
        </div>
      </div>
    </div>
  );
};

export default CarDetail;
