import { Link, Outlet } from "react-router-dom";
import {
  AppBar,
  Toolbar,
  Typography,
  Drawer,
  List,
  ListItem,
  ListItemText,
  Box,
  Button,
} from "@mui/material";
import { useAuth } from "../hooks/useAuth";
import UserMenu from "./UserMenu";

const drawerWidth = 240;

const Layout = () => {
  const { isAuthenticated, loading } = useAuth();

  return (
    <Box sx={{ display: "flex" }}>
      <AppBar
        position="fixed"
        sx={{ zIndex: (theme) => theme.zIndex.drawer + 1 }}
      >
        <Toolbar>
          <Typography variant="h6" noWrap component="div" sx={{ flexGrow: 1 }}>
            Football Pool
          </Typography>
          <Box sx={{ display: "flex", alignItems: "center" }}>
            {loading ? (
              // Show nothing while loading to avoid flash
              <Box sx={{ width: 80 }} />
            ) : isAuthenticated ? (
              <UserMenu />
            ) : (
              <Button
                color="inherit"
                component={Link}
                to="/login"
                data-testid="login-button"
              >
                Login
              </Button>
            )}
          </Box>
        </Toolbar>
      </AppBar>
      <Drawer
        variant="permanent"
        sx={{
          width: drawerWidth,
          flexShrink: 0,
          "& .MuiDrawer-paper": {
            width: drawerWidth,
            boxSizing: "border-box",
          },
        }}
      >
        <Toolbar />
        <Box sx={{ overflow: "auto" }}>
          <List>
            <ListItem component={Link} to="/">
              <ListItemText primary="Home" />
            </ListItem>
            {!loading && isAuthenticated && (
              <>
                <ListItem component={Link} to="/picks">
                  <ListItemText primary="Picks" />
                </ListItem>
                <ListItem component={Link} to="/results">
                  <ListItemText primary="Results" />
                </ListItem>
                <ListItem component={Link} to="/weekly-results">
                  <ListItemText primary="Weekly Results" />
                </ListItem>
                <ListItem component={Link} to="/overall-results">
                  <ListItemText primary="Overall Results" />
                </ListItem>
                <ListItem component={Link} to="/survivor-pool">
                  <ListItemText primary="Survivor Pool" />
                </ListItem>
              </>
            )}
            {!loading && !isAuthenticated && (
              <ListItem component={Link} to="/register">
                <ListItemText primary="Register" />
              </ListItem>
            )}
          </List>
        </Box>
      </Drawer>
      <Box component="main" sx={{ flexGrow: 1, p: 3 }}>
        <Toolbar />
        <Outlet />
      </Box>
    </Box>
  );
};

export default Layout;
