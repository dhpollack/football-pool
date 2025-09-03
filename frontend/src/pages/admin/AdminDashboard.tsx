import {
  Typography,
  Paper,
  Card,
  CardContent,
  Box,
} from "@mui/material";
import Grid from "@mui/material/Grid";
import {
  SportsFootball,
  People,
  EmojiEvents,
  BarChart,
} from "@mui/icons-material";

const AdminDashboard = () => {
  const stats = [
    {
      title: "Total Games",
      value: "0",
      icon: <SportsFootball />,
      color: "primary.main",
    },
    {
      title: "Total Users",
      value: "0",
      icon: <People />,
      color: "secondary.main",
    },
    {
      title: "Total Picks",
      value: "0",
      icon: <EmojiEvents />,
      color: "success.main",
    },
    {
      title: "Active Weeks",
      value: "0",
      icon: <BarChart />,
      color: "info.main",
    },
  ];

  return (
    <Box>
      <Typography variant="h4" gutterBottom>
        Admin Dashboard
      </Typography>
      <Typography variant="body1" color="text.secondary" sx={{ mb: 4 }}>
        Welcome to the Football Pool Administration Panel
      </Typography>

      <Grid container spacing={3}>
        {stats.map((stat) => (
          <Grid size={{ xs: 12, sm: 6, md: 3 }} key={stat.title}>
            <Card>
              <CardContent>
                <Box display="flex" alignItems="center">
                  <Box
                    sx={{
                      backgroundColor: stat.color,
                      color: "white",
                      borderRadius: "50%",
                      p: 1,
                      mr: 2,
                    }}
                  >
                    {stat.icon}
                  </Box>
                  <Box>
                    <Typography variant="h4" component="div">
                      {stat.value}
                    </Typography>
                    <Typography variant="body2" color="text.secondary">
                      {stat.title}
                    </Typography>
                  </Box>
                </Box>
              </CardContent>
            </Card>
          </Grid>
        ))}
      </Grid>

      <Grid container spacing={3} sx={{ mt: 2 }}>
        <Grid size={{ xs: 12, md: 6 }}>
          <Paper sx={{ p: 3 }}>
            <Typography variant="h6" gutterBottom>
              Recent Activity
            </Typography>
            <Typography variant="body2" color="text.secondary">
              No recent activity to display
            </Typography>
          </Paper>
        </Grid>
        <Grid size={{ xs: 12, md: 6 }}>
          <Paper sx={{ p: 3 }}>
            <Typography variant="h6" gutterBottom>
              Quick Actions
            </Typography>
            <Typography variant="body2" color="text.secondary">
              Manage games, users, and picks from the navigation menu
            </Typography>
          </Paper>
        </Grid>
      </Grid>
    </Box>
  );
};

export default AdminDashboard;