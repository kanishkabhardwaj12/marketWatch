import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import Navbar from './Navbar'; 
import './HomePage.css'; 

function HomePage() {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);
  const navigate = useNavigate();

  const [selectedFiles, setSelectedFiles] = useState(null);
  const [uploadMessage, setUploadMessage] = useState('');

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
  
  const handleUploadSubmit = async (e) => {
    e.preventDefault();
    if (!selectedFiles) {
      setUploadMessage('Please select a file to upload.');
      return;
    }
    setUploadMessage('Uploading...');

    const formData = new FormData();
    for (let i = 0; i < selectedFiles.length; i++) {
        formData.append('csvfiles', selectedFiles[i]);
    }
    
    try {
      const response = await fetch('/api/upload', {
        method: 'POST',
        body: formData,
      });
      if (response.ok) {
        setUploadMessage('File(s) uploaded successfully!');
      } else {
        const errorText = await response.text();
        setUploadMessage(`Upload failed: ${errorText}`);
      }
    } catch (error) {
      setUploadMessage('Upload failed. Please try again.');
    }
  };

  if (loading) {
    return <div>Loading your dashboard...</div>;
  }

  return (
    <div className="home-page-container">
      <Navbar user={user} />
      <div className="content">
        {user && user.Broker !== 'Zerodha' ? (
          <p>Support for {user.Broker} is under development.</p>
        ) : (
          <form onSubmit={handleUploadSubmit}>
          <h2>Upload Zerodha Trades</h2>
          <p>Please upload your tradebook CSV file.</p>
          <input 
            type="file" 
            name="csvfiles" 
            id="csvUpload" 
            accept=".csv" 
            multiple 
            required 
            onChange={(e) => setSelectedFiles(e.target.files)} 
          />
          <button type="submit">Upload</button>
          {uploadMessage && <p className="message">{uploadMessage}</p>}
        </form>
        )}
      </div>
    </div>
  );
}

export default HomePage;