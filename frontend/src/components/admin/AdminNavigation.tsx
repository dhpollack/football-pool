import { Link, useLocation } from "react-router-dom";
import {
  List,
  ListItem,
  ListItemIcon,
  ListItemText,
  ListItemButton,
  Divider,
  Typography,
  Box,
} from "@mui/material";
import {
  Dashboard,
  SportsFootball,
  People,
  EmojiEvents,
  BarChart,
  Settings,
} from "@mui/icons-material";

const AdminNavigation = () => {
  const location = useLocation();

  const menuItems = [
    {
      text: "Dashboard",
      icon: <Dashboard />,
      path: "/admin",
    },
    {
      text: "Game Management",
      icon: <SportsFootball />,
      path: "/admin/games",
    },
    {
      text: "User Management",
      icon: <People />,
      path: "/admin/users",
    },
    {
      text: "Pick Management",
      icon: <EmojiEvents />,
      path: "/admin/picks",
    },
    {
      text: "Statistics",
      icon: <BarChart />,
      path: "/admin/stats",
    },
  ];

  return (
    <Box>
      <Box sx={{ p: 2 }}>
        <Typography variant="h6" color="primary" gutterBottom>
          Admin Panel
        </Typography>
        <Typography variant="body2" color="text.secondary">
          Football Pool Administration
        </Typography>
      </Box>
      <Divider />
      <List>
        {menuItems.map((item) => (
          <ListItem key={item.path} disablePadding>
            <ListItemButton
              component={Link}
              to={item.path}
              selected={location.pathname === item.path}
              sx={{
                "&.Mui-selected": {
                  backgroundColor: "primary.light",
                  color: "primary.contrastText",
                  "&:hover": {
                    backgroundColor: "primary.main",
                  },
                },
                "&:hover": {
                  backgroundColor: "action.hover",
                },
              }}
            >
              <ListItemIcon
                sx={{
                  color: location.pathname === item.path ? "inherit" : "action.active",
                }}
              >
                {item.icon}
              </ListItemIcon>
              <ListItemText primary={item.text} />
            </ListItemButton>
          </ListItem>
        ))}
      </List>
      <Divider />
      <List>
        <ListItem disablePadding>
          <ListItemButton component={Link} to="/">
            <ListItemIcon>
              <Settings />
            </ListItemIcon>
            <ListItemText primary="Back to Main Site" />
          </ListItemButton>
        </ListItem>
      </List>
    </Box>
  );
};

export default AdminNavigation;