import { useState } from "react";
import { useParams, Link } from "react-router-dom";
import {
  Typography,
  Box,
  Alert,
  Chip,
  Card,
  CardContent,
  Grid,
  Breadcrumbs,
  Tab,
  Tabs,
} from "@mui/material";
import { Home, People } from "@mui/icons-material";
import { useAdminGetUser } from "../../services/api/user/user";
import { useAdminGetPicksByUser } from "../../services/api/picks/picks";
import AdminDataTable from "../../components/admin/AdminDataTable";
import { UserResponse, PickResponse } from "../../services/model";

interface TabPanelProps {
  children?: React.ReactNode;
  index: number;
  value: number;
}

function TabPanel(props: TabPanelProps) {
  const { children, value, index, ...other } = props;

  return (
    <div
      role="tabpanel"
      hidden={value !== index}
      id={`user-details-tabpanel-${index}`}
      aria-labelledby={`user-details-tab-${index}`}
      {...other}
    >
      {value === index && <Box sx={{ py: 3 }}>{children}</Box>}
    </div>
  );
}

const UserDetailsPage = () => {
  const { userId } = useParams<{ userId: string }>();
  const [activeTab, setActiveTab] = useState(0);

  const {
    data: userData,
    error: userError,
    isLoading: userLoading,
  } = useAdminGetUser(Number(userId));

  const {
    data: picksData,
    error: picksError,
    isLoading: picksLoading,
  } = useAdminGetPicksByUser(Number(userId));

  const user = userData || ({} as UserResponse);
  const picks = picksData?.picks || [];

  const handleTabChange = (_event: React.SyntheticEvent, newValue: number) => {
    setActiveTab(newValue);
  };

  const pickColumns = [
    {
      id: "id",
      label: "ID",
      align: "center" as const,
    },
    {
      id: "game",
      label: "Game",
      format: (pick: PickResponse) =>
        pick.game
          ? `${pick.game.favorite_team} vs ${pick.game.underdog_team}`
          : "Unknown",
    },
    {
      id: "week",
      label: "Week",
      align: "center" as const,
      format: (pick: PickResponse) => pick.game?.week || "-",
    },
    {
      id: "season",
      label: "Season",
      align: "center" as const,
      format: (pick: PickResponse) => pick.game?.season || "-",
    },
    {
      id: "picked",
      label: "Pick",
      align: "center" as const,
      format: (picked: string) => (
        <Chip
          label={picked}
          color={picked === "favorite" ? "primary" : "secondary"}
          size="small"
        />
      ),
    },
    {
      id: "created_at",
      label: "Submitted",
      format: (dateString: string) => new Date(dateString).toLocaleDateString(),
    },
  ];

  if (!userId) {
    return (
      <Alert severity="error">User ID is required to view user details.</Alert>
    );
  }

  return (
    <Box>
      {/* Breadcrumbs */}
      <Breadcrumbs sx={{ mb: 3 }}>
        <Link
          to="/admin"
          style={{
            display: "flex",
            alignItems: "center",
            textDecoration: "none",
            color: "inherit",
          }}
        >
          <Home sx={{ mr: 0.5 }} fontSize="inherit" />
          Admin
        </Link>
        <Link
          to="/admin/users"
          style={{
            display: "flex",
            alignItems: "center",
            textDecoration: "none",
            color: "inherit",
          }}
        >
          <People sx={{ mr: 0.5 }} fontSize="inherit" />
          Users
        </Link>
        <Typography color="text.primary">
          {userLoading ? "Loading..." : user.email || "User Details"}
        </Typography>
      </Breadcrumbs>

      {userError && (
        <Alert severity="error" sx={{ mb: 3 }}>
          Error loading user: {userError.message}
        </Alert>
      )}

      {/* User Info Card */}
      <Card sx={{ mb: 3 }}>
        <CardContent>
          <Grid container spacing={3}>
            <Grid size={{ xs: 12, md: 6 }}>
              <Typography variant="h6" gutterBottom>
                User Information
              </Typography>

              <Box sx={{ mb: 2 }}>
                <Typography variant="subtitle2" color="text.secondary">
                  Email
                </Typography>
                <Typography variant="body1">{user.email || "-"}</Typography>
              </Box>

              <Box sx={{ mb: 2 }}>
                <Typography variant="subtitle2" color="text.secondary">
                  Name
                </Typography>
                <Typography variant="body1">{user.name || "-"}</Typography>
              </Box>

              <Box sx={{ mb: 2 }}>
                <Typography variant="subtitle2" color="text.secondary">
                  Role
                </Typography>
                <Chip
                  label={user.role || "user"}
                  color={user.role === "admin" ? "primary" : "default"}
                  size="small"
                />
              </Box>
            </Grid>

            <Grid size={{ xs: 12, md: 6 }}>
              <Typography variant="h6" gutterBottom>
                Player Profile
              </Typography>

              <Box sx={{ mb: 2 }}>
                <Typography variant="subtitle2" color="text.secondary">
                  Player Name
                </Typography>
                <Typography variant="body1">
                  {user.player?.name || "-"}
                </Typography>
              </Box>

              <Box sx={{ mb: 2 }}>
                <Typography variant="subtitle2" color="text.secondary">
                  Address
                </Typography>
                <Typography variant="body1">
                  {user.player?.address || "-"}
                </Typography>
              </Box>

              <Box sx={{ mb: 2 }}>
                <Typography variant="subtitle2" color="text.secondary">
                  Member Since
                </Typography>
                <Typography variant="body1">
                  {user.created_at
                    ? new Date(user.created_at).toLocaleDateString()
                    : "-"}
                </Typography>
              </Box>
            </Grid>
          </Grid>
        </CardContent>
      </Card>

      {/* Tabs */}
      <Box sx={{ borderBottom: 1, borderColor: "divider" }}>
        <Tabs value={activeTab} onChange={handleTabChange}>
          <Tab label="Picks History" />
          <Tab label="Statistics" />
        </Tabs>
      </Box>

      {/* Picks Tab */}
      <TabPanel value={activeTab} index={0}>
        {picksError && (
          <Alert severity="error" sx={{ mb: 3 }}>
            Error loading picks: {picksError.message}
          </Alert>
        )}

        <AdminDataTable
          columns={pickColumns}
          data={picks}
          loading={picksLoading}
          error={picksError?.message || null}
          page={0}
          rowsPerPage={10}
          totalCount={picks.length}
          onPageChange={() => {}}
          onRowsPerPageChange={() => {}}
          emptyMessage={
            picksLoading ? "Loading picks..." : "No picks found for this user"
          }
        />
      </TabPanel>

      {/* Statistics Tab */}
      <TabPanel value={activeTab} index={1}>
        <Typography variant="h6" gutterBottom>
          User Statistics
        </Typography>
        <Typography color="text.secondary">
          Statistics functionality is coming soon. This will include win/loss
          records, weekly performance, and overall standings.
        </Typography>
      </TabPanel>
    </Box>
  );
};

export default UserDetailsPage;
