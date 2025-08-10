import React, { useState, useEffect } from 'react';
import axios from 'axios';

const ProfilePage = () => {
  const [profile, setProfile] = useState({});
  const [name, setName] = useState('');
  const [address, setAddress] = useState('');

  useEffect(() => {
    const fetchProfile = async () => {
      const { data } = await axios.get('/api/users/me');
      setProfile(data);
      setName(data.Name);
      setAddress(data.Address);
    };
    fetchProfile();
  }, []);

  const handleUpdate = async () => {
    try {
      await axios.put('/api/users/me/update', { name, address });
    } catch (error) {
      console.error('Failed to update profile', error);
    }
  };

  return (
    <div>
      <h1>Profile</h1>
      <input
        type="text"
        placeholder="Name"
        value={name}
        onChange={(e) => setName(e.target.value)}
      />
      <input
        type="text"
        placeholder="Address"
        value={address}
        onChange={(e) => setAddress(e.target.value)}
      />
      <button onClick={handleUpdate}>Update</button>
    </div>
  );
};

export default ProfilePage;
