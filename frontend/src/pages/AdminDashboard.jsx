import { useState, useEffect, useCallback } from 'react';
import api from '../api/axios';
import { 
  BarChart3, 
  Car, 
  Users, 
  Plus, 
  Edit3, 
  Trash2, 
  Upload, 
  TrendingUp, 
  DollarSign,
  X,
  Check,
  Loader2
} from 'lucide-react';
import { toast } from 'sonner';

const AdminDashboard = () => {
  const [activeTab, setActiveTab] = useState('overview');
  const [stats, setStats] = useState(null);
  const [cars, setCars] = useState([]);
  const [loading, setLoading] = useState(true);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [editingCar, setEditingCar] = useState(null);
  const [uploading, setUploading] = useState(false);

  const [formData, setFormData] = useState({
    name: '',
    category: 'Sedan',
    rental_costs: '',
    stock_availability: '',
    availability: true,
    description: '',
    image_url: ''
  });

  const fetchStats = useCallback(async () => {
    try {
      const response = await api.get('/admin/dashboard');
      setStats(response.data.data);
    } catch {
      toast.error('Failed to fetch stats');
    }
  }, []);

  const fetchCars = useCallback(async () => {
    try {
      // Get all cars (not just available ones)
      const response = await api.get('/cars');
      setCars(response.data.data.cars);
    } catch {
      toast.error('Failed to fetch cars');
    }
  }, []);

  const initDashboard = useCallback(async () => {
    setLoading(true);
    await Promise.all([fetchStats(), fetchCars()]);
    setLoading(false);
  }, [fetchStats, fetchCars]);

  useEffect(() => {
    initDashboard();
  }, [initDashboard]);

  const handleInputChange = (e) => {
    const { name, value, type, checked } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: type === 'checkbox' ? checked : value
    }));
  };

  const handleImageUpload = async (e) => {
    const file = e.target.files[0];
    if (!file) return;

    const data = new FormData();
    data.append('image', file);

    try {
      setUploading(true);
      const response = await api.post('/cars/upload', data, {
        headers: { 'Content-Type': 'multipart/form-data' }
      });
      setFormData(prev => ({ ...prev, image_url: response.data.data.url }));
      toast.success('Image uploaded successfully');
    } catch {
      toast.error('Failed to upload image');
    } finally {
      setUploading(false);
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    const data = {
      ...formData,
      rental_costs: parseFloat(formData.rental_costs),
      stock_availability: parseInt(formData.stock_availability)
    };

    try {
      if (editingCar) {
        await api.put(`/cars/${editingCar.id}`, data);
        toast.success('Car updated successfully');
      } else {
        await api.post('/cars', data);
        toast.success('Car added successfully');
      }
      setIsModalOpen(false);
      initDashboard();
    } catch (err) {
      toast.error(err.response?.data?.message || 'Action failed');
    }
  };

  const handleDelete = async (id) => {
    if (!window.confirm('Are you sure you want to delete this car?')) return;
    try {
      await api.delete(`/cars/${id}`);
      toast.success('Car deleted');
      initDashboard();
    } catch {
      toast.error('Failed to delete car');
    }
  };

  const openEditModal = (car) => {
    setEditingCar(car);
    setFormData({
      name: car.name,
      category: car.category,
      rental_costs: car.rental_costs,
      stock_availability: car.stock_availability,
      availability: car.availability,
      description: car.description,
      image_url: car.image_url
    });
    setIsModalOpen(true);
  };

  const openAddModal = () => {
    setEditingCar(null);
    setFormData({
      name: '',
      category: 'Sedan',
      rental_costs: '',
      stock_availability: '',
      availability: true,
      description: '',
      image_url: ''
    });
    setIsModalOpen(true);
  };

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-dark">
        <Loader2 className="w-12 h-12 text-primary animate-spin" />
      </div>
    );
  }

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-10">
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-6 mb-10">
        <div>
          <h1 className="text-3xl font-black text-white tracking-tight">Admin Control</h1>
          <p className="text-gray-500 font-medium">Manage your fleet and track performance.</p>
        </div>
        <div className="flex bg-dark-card border border-white/5 p-1 rounded-xl">
          <button 
            onClick={() => setActiveTab('overview')}
            className={`px-6 py-2 rounded-lg font-bold transition-all ${activeTab === 'overview' ? 'bg-primary text-white shadow-lg' : 'text-gray-500 hover:text-gray-300'}`}
          >
            Overview
          </button>
          <button 
            onClick={() => setActiveTab('fleet')}
            className={`px-6 py-2 rounded-lg font-bold transition-all ${activeTab === 'fleet' ? 'bg-primary text-white shadow-lg' : 'text-gray-500 hover:text-gray-300'}`}
          >
            Manage Fleet
          </button>
        </div>
      </div>

      {activeTab === 'overview' ? (
        <div className="space-y-10 animate-in fade-in slide-in-from-bottom-4 duration-500">
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
            <div className="bg-dark-card border border-white/5 p-8 rounded-2xl shadow-xl">
              <div className="w-12 h-12 bg-emerald-500/10 text-emerald-500 rounded-xl flex items-center justify-center mb-4">
                <DollarSign size={24} />
              </div>
              <p className="text-xs font-black text-gray-500 uppercase tracking-widest mb-1">Total Revenue</p>
              <h2 className="text-3xl font-black text-white tracking-tight">IDR {stats?.stats?.total_revenue?.toLocaleString()}</h2>
            </div>
            <div className="bg-dark-card border border-white/5 p-8 rounded-2xl shadow-xl">
              <div className="w-12 h-12 bg-blue-500/10 text-blue-500 rounded-xl flex items-center justify-center mb-4">
                <TrendingUp size={24} />
              </div>
              <p className="text-xs font-black text-gray-500 uppercase tracking-widest mb-1">Total Rentals</p>
              <h2 className="text-3xl font-black text-white tracking-tight">{stats?.stats?.total_rentals}</h2>
            </div>
            <div className="bg-dark-card border border-white/5 p-8 rounded-2xl shadow-xl">
              <div className="w-12 h-12 bg-purple-500/10 text-purple-500 rounded-xl flex items-center justify-center mb-4">
                <Users size={24} />
              </div>
              <p className="text-xs font-black text-gray-500 uppercase tracking-widest mb-1">Active Users</p>
              <h2 className="text-3xl font-black text-white tracking-tight">{stats?.stats?.active_users}</h2>
            </div>
          </div>

          <div className="bg-dark-card border border-white/5 rounded-2xl overflow-hidden shadow-2xl">
            <div className="px-8 py-6 border-b border-white/5 bg-white/[0.02] flex items-center gap-3">
              <div className="text-primary bg-primary/10 p-2 rounded-lg"><TrendingUp size={18} /></div>
              <h2 className="font-bold text-white tracking-tight text-lg">Top Performing Vehicles</h2>
            </div>
            <div className="overflow-x-auto">
              <table className="w-full text-left">
                <thead>
                  <tr className="bg-white/[0.01]">
                    <th className="px-8 py-4 text-[10px] font-black text-gray-500 uppercase tracking-widest">Car Model</th>
                    <th className="px-8 py-4 text-[10px] font-black text-gray-500 uppercase tracking-widest text-center">Rental Count</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-white/5">
                  {stats?.popular_cars?.map((car, idx) => (
                    <tr key={idx} className="hover:bg-white/[0.02] transition-colors">
                      <td className="px-8 py-5">
                        <div className="flex items-center gap-4">
                          <div className="w-8 h-8 bg-white/5 rounded-lg flex items-center justify-center text-xs font-black text-gray-400 border border-white/5">{idx + 1}</div>
                          <span className="font-bold text-white tracking-tight">{car.car_name}</span>
                        </div>
                      </td>
                      <td className="px-8 py-5 text-center">
                        <span className="bg-primary/10 text-primary px-4 py-1 rounded-full text-xs font-black border border-primary/20">{car.rental_count} bookings</span>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>
        </div>
      ) : (
        <div className="animate-in fade-in slide-in-from-bottom-4 duration-500">
          <div className="bg-dark-card border border-white/5 rounded-2xl overflow-hidden shadow-2xl">
            <div className="px-8 py-6 border-b border-white/5 bg-white/[0.02] flex items-center justify-between">
              <div className="flex items-center gap-3">
                <div className="text-primary bg-primary/10 p-2 rounded-lg"><Car size={18} /></div>
                <h2 className="font-bold text-white tracking-tight text-lg">Fleet Management</h2>
              </div>
              <button 
                onClick={openAddModal}
                className="flex items-center gap-2 bg-primary text-white px-5 py-2 rounded-xl font-black text-xs uppercase tracking-widest hover:bg-blue-600 transition-all active:scale-95 shadow-lg shadow-primary/20"
              >
                <Plus size={16} /> <span>Add Vehicle</span>
              </button>
            </div>
            <div className="overflow-x-auto">
              <table className="w-full text-left border-collapse">
                <thead>
                  <tr className="bg-white/[0.01]">
                    <th className="px-8 py-4 text-[10px] font-black text-gray-500 uppercase tracking-widest">Model</th>
                    <th className="px-8 py-4 text-[10px] font-black text-gray-500 uppercase tracking-widest">Category</th>
                    <th className="px-8 py-4 text-[10px] font-black text-gray-500 uppercase tracking-widest">Stock</th>
                    <th className="px-8 py-4 text-[10px] font-black text-gray-500 uppercase tracking-widest">Price</th>
                    <th className="px-8 py-4 text-[10px] font-black text-gray-500 uppercase tracking-widest text-right">Actions</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-white/5">
                  {cars.map(car => (
                    <tr key={car.id} className="hover:bg-white/[0.02] transition-colors">
                      <td className="px-8 py-5">
                        <div className="flex items-center gap-4">
                          <img src={car.image_url} alt="" className="w-12 h-12 rounded-lg object-cover bg-white/5" />
                          <div>
                            <p className="font-black text-white tracking-tight leading-tight">{car.name}</p>
                            {!car.availability && <span className="text-[9px] font-black text-red-500 uppercase tracking-widest">Hidden</span>}
                          </div>
                        </div>
                      </td>
                      <td className="px-8 py-5 text-xs font-bold text-gray-400 uppercase tracking-widest">{car.category}</td>
                      <td className="px-8 py-5 text-sm font-black text-white">{car.stock_availability}</td>
                      <td className="px-8 py-5 text-sm font-black text-primary">IDR {car.rental_costs?.toLocaleString()}</td>
                      <td className="px-8 py-5 text-right space-x-2">
                        <button 
                          onClick={() => openEditModal(car)}
                          className="p-2 text-blue-400 hover:bg-blue-500/10 rounded-lg transition-colors"
                        >
                          <Edit3 size={18} />
                        </button>
                        <button 
                          onClick={() => handleDelete(car.id)}
                          className="p-2 text-red-400 hover:bg-red-500/10 rounded-lg transition-colors"
                        >
                          <Trash2 size={18} />
                        </button>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>
        </div>
      )}

      {isModalOpen && (
        <div className="fixed inset-0 z-[100] flex items-center justify-center p-4">
          <div className="absolute inset-0 bg-black/80 backdrop-blur-md" onClick={() => setIsModalOpen(false)}></div>
          <div className="relative bg-dark-card border border-white/10 p-8 rounded-[2rem] w-full max-w-2xl max-h-[90vh] overflow-y-auto shadow-2xl">
            <div className="flex justify-between items-center mb-10">
              <h3 className="text-2xl font-black text-white tracking-tight">{editingCar ? 'Edit Vehicle' : 'New Vehicle'}</h3>
              <button onClick={() => setIsModalOpen(false)} className="text-gray-500 hover:text-white transition-colors"><X /></button>
            </div>

            <form onSubmit={handleSubmit} className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div className="space-y-2">
                <label className="text-[10px] font-black text-gray-500 uppercase tracking-widest">Model Name</label>
                <input
                  name="name"
                  className="w-full bg-black/40 border border-white/10 rounded-2xl py-3 px-5 text-white focus:outline-none focus:border-primary transition-all"
                  value={formData.name}
                  onChange={handleInputChange}
                  required
                />
              </div>

              <div className="space-y-2">
                <label className="text-[10px] font-black text-gray-500 uppercase tracking-widest">Category</label>
                <select
                  name="category"
                  className="w-full bg-black/40 border border-white/10 rounded-2xl py-3 px-5 text-white focus:outline-none focus:border-primary transition-all"
                  value={formData.category}
                  onChange={handleInputChange}
                >
                  <option value="Sedan">Sedan</option>
                  <option value="SUV">SUV</option>
                  <option value="Luxury">Luxury</option>
                  <option value="Sport">Sport</option>
                </select>
              </div>

              <div className="space-y-2">
                <label className="text-[10px] font-black text-gray-500 uppercase tracking-widest">Daily Cost (IDR)</label>
                <input
                  name="rental_costs"
                  type="number"
                  className="w-full bg-black/40 border border-white/10 rounded-2xl py-3 px-5 text-white focus:outline-none focus:border-primary transition-all"
                  value={formData.rental_costs}
                  onChange={handleInputChange}
                  required
                />
              </div>

              <div className="space-y-2">
                <label className="text-[10px] font-black text-gray-500 uppercase tracking-widest">Stock Level</label>
                <input
                  name="stock_availability"
                  type="number"
                  className="w-full bg-black/40 border border-white/10 rounded-2xl py-3 px-5 text-white focus:outline-none focus:border-primary transition-all"
                  value={formData.stock_availability}
                  onChange={handleInputChange}
                  required
                />
              </div>

              <div className="md:col-span-2 space-y-2">
                <label className="text-[10px] font-black text-gray-500 uppercase tracking-widest">Description</label>
                <textarea
                  name="description"
                  rows="3"
                  className="w-full bg-black/40 border border-white/10 rounded-2xl py-3 px-5 text-white focus:outline-none focus:border-primary transition-all resize-none"
                  value={formData.description}
                  onChange={handleInputChange}
                />
              </div>

              <div className="md:col-span-2 space-y-4">
                <label className="text-[10px] font-black text-gray-500 uppercase tracking-widest">Vehicle Image</label>
                <div className="flex flex-col md:flex-row gap-4">
                  {formData.image_url && (
                    <img src={formData.image_url} className="w-24 h-24 rounded-2xl object-cover border border-white/10" alt="Preview" />
                  )}
                  <div className="flex-1 relative">
                    <input 
                      type="file" 
                      onChange={handleImageUpload}
                      className="hidden" 
                      id="image-upload" 
                      disabled={uploading}
                    />
                    <label 
                      htmlFor="image-upload" 
                      className="w-full h-full min-h-[100px] border-2 border-dashed border-white/10 rounded-2xl flex flex-col items-center justify-center gap-2 hover:border-primary/50 hover:bg-primary/5 transition-all cursor-pointer group"
                    >
                      {uploading ? (
                        <Loader2 className="animate-spin text-primary" />
                      ) : (
                        <>
                          <Upload className="text-gray-500 group-hover:text-primary transition-colors" />
                          <span className="text-xs font-bold text-gray-500 group-hover:text-primary transition-colors uppercase tracking-widest">Click to Upload</span>
                        </>
                      )}
                    </label>
                  </div>
                </div>
              </div>

              <div className="md:col-span-2 flex items-center gap-3 bg-white/5 p-4 rounded-2xl">
                <input 
                  type="checkbox" 
                  name="availability" 
                  id="availability"
                  checked={formData.availability}
                  onChange={handleInputChange}
                  className="w-5 h-5 accent-primary" 
                />
                <label htmlFor="availability" className="text-sm font-bold text-gray-300">Publish this vehicle immediately</label>
              </div>

              <div className="md:col-span-2 flex gap-4 pt-6">
                <button 
                  type="button" 
                  className="flex-1 border border-white/10 text-white font-bold py-4 rounded-2xl hover:bg-white/5 transition-all" 
                  onClick={() => setIsModalOpen(false)}
                >
                  Discard
                </button>
                <button 
                  type="submit" 
                  className="flex-1 bg-primary text-white font-black py-4 rounded-2xl hover:bg-blue-600 shadow-2xl shadow-primary/30 transition-all active:scale-95 uppercase tracking-widest text-sm"
                >
                  {editingCar ? 'Update Vehicle' : 'Add Vehicle'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
};

export default AdminDashboard;
