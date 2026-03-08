const CarSkeleton = () => {
  return (
    <div className="bg-dark-card rounded-2xl overflow-hidden border border-white/5 shadow-xl">
      <div className="h-48 w-full bg-white/5 animate-pulse"></div>
      <div className="p-6 space-y-4">
        <div className="h-6 w-2/3 bg-white/5 animate-pulse rounded-lg"></div>
        <div className="flex gap-3">
          <div className="h-4 w-1/4 bg-white/5 animate-pulse rounded-md"></div>
          <div className="h-4 w-1/4 bg-white/5 animate-pulse rounded-md"></div>
        </div>
        <div className="space-y-2">
          <div className="h-3.5 w-full bg-white/5 animate-pulse rounded-md"></div>
          <div className="h-3.5 w-4/5 bg-white/5 animate-pulse rounded-md"></div>
        </div>
        <div className="pt-4 border-t border-white/5 flex justify-between items-center">
          <div className="h-6 w-1/3 bg-white/5 animate-pulse rounded-lg"></div>
          <div className="h-9 w-24 bg-white/5 animate-pulse rounded-lg"></div>
        </div>
      </div>
    </div>
  );
};

export default CarSkeleton;
