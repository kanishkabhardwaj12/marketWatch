import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import Navbar from './Navbar';
import FileList from './FileList'; 
import './ProfilePage.css'; 

function ProfilePage() {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);
  const navigate = useNavigate();

  const [newPhoto, setNewPhoto] = useState(null);
  const [newPassword, setNewPassword] = useState('');
  const [message, setMessage] = useState(''); 

  useEffect(() => {
    const fetchUserData = async () => {
      try {
        const response = await fetch('/api/profile');
        if (!response.ok) {
          navigate('/'); 
          return;
        }
        const data = await response.json();
        setUser(data);
      } catch (error) {
        console.error('Failed to fetch user data:', error);
        navigate('/');
      } finally {
        setLoading(false);
      }
    };
    fetchUserData();
  }, [navigate]);

  const handleFormSubmit = async (e) => {
    e.preventDefault();
    setMessage('');

    if (!newPhoto && !newPassword) {
      setMessage('Please select a new photo or enter a new password to update.');
      return;
    }

    const formData = new FormData();
    if (newPhoto) {
      formData.append('photo', newPhoto);
    }
    if (newPassword) {
      formData.append('password', newPassword);
    }

    try {
      const response = await fetch('/api/profile', {
        method: 'POST',
        body: formData,
      });

      if (response.ok) {
        setMessage('Profile updated successfully! It will be reflected on your next login.');

        setNewPhoto(null);
        setNewPassword('');
        document.getElementById('photoUpload').value = null; // Reset file input
      } else {
        const errorText = await response.text();
        setMessage(`Update failed: ${errorText}`);
      }
    } catch (error) {
      setMessage('An error occurred. Please try again.');
    }
  };

  if (loading) {
    return <div>Loading profile...</div>;
  }

  return (
    <div className="profile-page-background">
      <Navbar user={user} />
      <div className="profile-container">
        <h2>User Profile</h2>
        <div className="profile-photo">
          <img src={user.Photo} alt="Profile" />
        </div>

        <form onSubmit={handleFormSubmit}>
          <div className="form-group">
            <label>Email:</label>
            <input type="email" value={user.Email} readOnly />
          </div>
          <div className="form-group">
            <label>Username:</label>
            <input type="text" value={user.Name} readOnly />
          </div>
          <div className="form-group">
            <label>Broker:</label>
            <input type="text" value={user.Broker} readOnly />
          </div>

          <div className="section">
            <label htmlFor="photoUpload">Change Profile Photo:</label>
            <input
              type="file"
              id="photoUpload"
              name="photo"
              accept="image/*"
              onChange={(e) => setNewPhoto(e.target.files[0])}
            />
          </div>

          <div className="section">
            <label htmlFor="newPassword">New Password:</label>
            <input
              type="password"
              id="newPassword"
              name="password"
              placeholder="Enter new password"
              value={newPassword}
              onChange={(e) => setNewPassword(e.target.value)}
            />
          </div>

          <button type="submit" className="btn">Save Changes</button>
          
          {message && <p className="message">{message}</p>}
        </form>

        <FileList />
      </div>
    </div>
  );
}

export default ProfilePage;