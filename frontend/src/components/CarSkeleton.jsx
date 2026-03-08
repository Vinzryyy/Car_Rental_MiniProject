import React from 'react';
import './CarSkeleton.css';

const CarSkeleton = () => {
  return (
    <div className="car-card skeleton-card">
      <div className="skeleton-image pulse"></div>
      <div className="skeleton-info">
        <div className="skeleton-title pulse"></div>
        <div className="skeleton-specs">
          <div className="skeleton-spec pulse"></div>
          <div className="skeleton-spec pulse"></div>
        </div>
        <div className="skeleton-desc pulse"></div>
        <div className="skeleton-desc pulse" style={{ width: '80%' }}></div>
        <div className="skeleton-footer">
          <div className="skeleton-price pulse"></div>
          <div className="skeleton-btn pulse"></div>
        </div>
      </div>
    </div>
  );
};

export default CarSkeleton;
