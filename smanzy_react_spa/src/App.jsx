import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { RouterProvider } from 'react-router-dom';
import router from '@/routes/index.jsx';
import { ThemeProvider } from '@/components/ThemeProvider';

import { VersionProvider } from '@/context/VersionContext';
import { UserProvider } from '@/context/UserContext';

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

function App() {
    return (
        <ThemeProvider>
            <QueryClientProvider client={queryClient}>
                <VersionProvider>
                    <UserProvider>
                        <RouterProvider router={router} />
                    </UserProvider>
                </VersionProvider>
            </QueryClientProvider>
        </ThemeProvider>
    );
}

export default App;
