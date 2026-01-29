import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { ThemeProvider } from "@/components/ThemeProvider";
import { VersionProvider } from "@/context/VersionProvider";
import { UserProvider } from "@/context/UserContext";
import { RouteProvider } from "@/context/RouteContext";
import AppRouter from "@/routes/AppRouter";

// The QueryClient is used to manage and cache data throughout the application.
// Create a QueryClient instance
const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 1000 * 60 * 5, // 5 minutes
      retry: 1,
    },
  },
});

const App = () => {
  return (
    <ThemeProvider>
      <QueryClientProvider client={queryClient}>
        <VersionProvider>
          <UserProvider>
            <RouteProvider>
              <AppRouter />
            </RouteProvider>
          </UserProvider>
        </VersionProvider>
      </QueryClientProvider>
    </ThemeProvider>
  );
};

export default App;
