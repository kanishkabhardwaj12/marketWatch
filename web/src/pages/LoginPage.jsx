import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import './LoginPage.css'; 

function LoginPage() {
  // initialising everything we will need
  // the navigate variable if we need to redirect to another page
  const navigate = useNavigate();

  // for login, second argument is a funstion that you call to update the value of the first function
  const [loginData, setLoginData] = useState({ email: '', password: '' });
  const [loginError, setLoginError] = useState('');

  // for signup
  const [signupData, setSignupData] = useState({
    email: '',
    name: '',
    password: '',
    broker: '',
  });

  //sent photo seperately incase a user wants to not use photo at signup
  const [signupPhoto, setSignupPhoto] = useState(null);
  const [signupError, setSignupError] = useState('');
  
  const handleLoginChange = (e) => {
    // e.target is the imput field that the user is typing in
    // it pulls the name="email" and value=the typed email
    const { name, value } = e.target;
    setLoginData(prevData => ({ ...prevData, [name]: value }));
  };

  const handleLoginSubmit = async (e) => {
    e.preventDefault();
    setLoginError('');

    const formData = new FormData();
    formData.append('email', loginData.email);
    formData.append('password', loginData.password);
    
    try {
      const response = await fetch('/api/login', {
        method: 'POST',
        body: formData,
      });

      if (response.ok) {
        navigate('/home'); 
      } else {
        const errorText = await response.text();
        setLoginError(errorText || 'Invalid credentials.');
      }
    } catch (error) {
      setLoginError('Login failed. Please try again.');
    }
  };

  const handleSignupChange = (e) => {
    const { name, value } = e.target;
    setSignupData(prevData => ({ ...prevData, [name]: value }));
  };

  const handlePhotoChange = (e) => {
    if (e.target.files[0]) {
      setSignupPhoto(e.target.files[0]);
    }
  };

  const handleSignupSubmit = async (e) => {
    e.preventDefault();
    setSignupError('');

    const formData = new FormData();
    formData.append('email', signupData.email);
    formData.append('name', signupData.name);
    formData.append('password', signupData.password);
    formData.append('broker', signupData.broker);
    if (signupPhoto) {
      formData.append('photo', signupPhoto);
    }

    try {
      const response = await fetch('/api/signup', {
        method: 'POST',
        body: formData,
      });
      
      if (response.ok) {
        navigate('/home'); 
      } else {
        const errorText = await response.text();
        setSignupError(errorText || 'Signup failed.');
      }
    } catch (error) {
      setSignupError('Signup failed. Please try again.');
    }
  };

  return (
    <div className="login-page-background"> 
      <h1 className="login-title">MarketWatch</h1>
      <div className="container">
        {/* Login Form */}
        <div className="form-container">
          <h2>Login</h2>
          {loginError && <div className="error">{loginError}</div>}
          <form onSubmit={handleLoginSubmit}>
            <input 
              type="email" 
              name="email" 
              placeholder="Email ID" 
              value={loginData.email}
              onChange={handleLoginChange}
              required 
            />
            <input 
              type="password" 
              name="password" 
              placeholder="Password" 
              value={loginData.password}
              onChange={handleLoginChange}
              required 
            />
            <button type="submit">Log In</button>
          </form>
        </div>

        {/* Signup Form */}
        <div className="form-container">
          <h2>Sign Up</h2>
          {signupError && <div className="error">{signupError}</div>}
          <form onSubmit={handleSignupSubmit}>
            <input 
              type="email" 
              name="email" 
              placeholder="Email ID" 
              value={signupData.email}
              onChange={handleSignupChange}
              required
            />
            <input 
              type="text" 
              name="name" 
              placeholder="Name" 
              value={signupData.name}
              onChange={handleSignupChange}
              required 
            />
            <input 
              type="password" 
              name="password" 
              placeholder="Password"
              value={signupData.password}
              onChange={handleSignupChange} 
              required 
            />
            <select 
              name="broker" 
              value={signupData.broker}
              onChange={handleSignupChange}
              required
            >
              <option value="" disabled>Select Broker</option>
              <option value="Zerodha">Zerodha</option>
              <option value="Upstox">Upstox</option>
              <option value="Groww">Groww</option>
            </select>
            <label className="upload">
              <input 
                type="file" 
                name="photo" 
                accept="image/*" 
                onChange={handlePhotoChange}
                hidden 
              />
              {signupPhoto ? signupPhoto.name : 'Upload Photo'}
            </label>
            <button type="submit">Sign Up</button>
          </form>
        </div>
      </div>
    </div>
  );
}

export default LoginPage;