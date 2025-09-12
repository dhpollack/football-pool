import { render, screen, waitFor } from "@testing-library/react";
import AdminUsersPage from "./AdminUsersPage";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { setupServer } from "msw/node";
import { HttpResponse, delay, http } from "msw";
import { getAdminListUsersMockHandler } from "../../services/api/user/user.msw";

const server = setupServer();

beforeAll(() => server.listen());
afterEach(() => server.resetHandlers());
afterAll(() => server.close());

describe("Debug Test", () => {
  const queryClient = new QueryClient();
  const wrapper = ({ children }: { children: React.ReactNode }) => (
    <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
  );

  it("should debug the component rendering", async () => {
    // Add debug logging to the MSW handler
    const debugHandler = http.get("*/api/admin/users", async (info) => {
      console.log("MSW intercepted request:", info.request.url);
      console.log("Request headers:", Object.fromEntries(info.request.headers.entries()));
      
      const responseData = {
        users: [{
          id: 1,
          email: "test@example.com",
          name: "Test User",
          role: "user",
          pick_count: 0,
          total_wins: 0,
          created_at: "2024-01-01T00:00:00Z",
          updated_at: "2024-01-01T00:00:00Z"
        }],
        pagination: { total: 1, page: 1, limit: 20 }
      };
      
      console.log("MSW returning response:", responseData);
      
      const response = new HttpResponse(
        JSON.stringify(responseData),
        { status: 200, headers: { "Content-Type": "application/json" } }
      );
      
      console.log("MSW response status:", response.status);
      console.log("MSW response headers:", Object.fromEntries(response.headers.entries()));
      
      return response;
    });
    
    server.use(debugHandler);

    render(<AdminUsersPage />, { wrapper });

    // Wait for the component to render
    await screen.findByText("User Management");
    
    // Wait for React Query to finish loading and component to re-render
    await waitFor(() => {
      expect(queryClient.getQueryState(['/api/admin/users'])?.status).toBe('success');
    });
    
    // Check React Query state
    const queryState = queryClient.getQueryState(["/api/admin/users"]);
    console.log("React Query state:", queryState);
    
    // Check if query has data
    const queryData = queryClient.getQueryData(["/api/admin/users"]);
    console.log("React Query data:", queryData);
    
    // Check all queries
    const allQueries = queryClient.getQueryCache().getAll();
    console.log("All queries:", allQueries.map(q => ({
      queryKey: q.queryKey,
      state: q.state
    })));
    
    // Debug: log the entire document
    console.log("Document body:", document.body.innerHTML);
    
    // Check if the table is rendered
    const table = screen.queryByRole("table");
    console.log("Table found:", !!table);
    
    if (table) {
      console.log("Table content:", table.innerHTML);
    }
    
    // Check if loading state is stuck
    const loading = screen.queryByRole("progressbar");
    console.log("Loading:", !!loading);
    
    // Check for error messages
    const error = screen.queryByRole("alert");
    console.log("Error:", error?.textContent);
  });
});