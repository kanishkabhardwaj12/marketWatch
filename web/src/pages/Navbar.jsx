import React, { useState, useEffect, useRef } from 'react';
import { Link, useNavigate } from 'react-router-dom'; 
import './Navbar.css';

function Navbar({ user }) {
  const [isDropdownOpen, setDropdownOpen] = useState(false);
  const dropdownRef = useRef(null);
  const navigate = useNavigate();

  useEffect(() => {
    const handleClickOutside = (event) => {
      if (dropdownRef.current && !dropdownRef.current.contains(event.target)) {
        setDropdownOpen(false);
      }
    };
    document.addEventListener('mousedown', handleClickOutside);
    return () => {
      document.removeEventListener('mousedown', handleClickOutside);
    };
  }, []);

  const handleLogout = async () => {
    try {
      const response = await fetch('/api/logout');

      if (response.ok) {
        navigate('/');
      } else {
        console.error('Logout failed on the server.');
      }
    } catch (error) {
      console.error('Failed to send logout request:', error);
    }
  };


  if (!user) {
    return null;
  }

  return (
    <div className="navbar">
      <h1>MarketWatch</h1>
      <div 
        className="user-section" 
        onClick={() => setDropdownOpen(!isDropdownOpen)} 
        ref={dropdownRef}
      >
        <span className="username">{user.Name}</span>
        <img src={user.Photo} alt="User" className="user-photo" />
        
        {isDropdownOpen && (
          <div className="dropdown">
            <Link to="/profile" className="logout-button">Profile</Link>
            <button onClick={handleLogout} className="logout-button">Logout</button>
          </div>
        )}
      </div>
    </div>
  );
}

export default Navbar;