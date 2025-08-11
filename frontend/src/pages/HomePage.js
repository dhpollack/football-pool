import React from 'react';
import { Button, Typography } from '@mui/material';
import { Link } from 'react-router-dom';

const HomePage = () => {
  return (
    <div>
      <Typography variant="h4" gutterBottom>
        Welcome to the Football Pool!
      </Typography>
      <Typography variant="body1" paragraph>
        Get ready to make your picks for the upcoming games.
      </Typography>
      <Button variant="contained" color="primary" component={Link} to="/picks">
        Make Picks
      </Button>
    </div>
  );
};

export default HomePage;
