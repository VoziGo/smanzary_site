import React, { useMemo, useContext } from 'react';
import { createBrowserRouter, RouterProvider } from 'react-router-dom';
import { ComponentMap } from './componentMap';
import { NotFound } from '@/pages'; // Fallback import

// Import your Context (Assuming you have one)
import { RouteContext } from '@/context/RouteContext';

const DynamicRouter = () => {
    // 1. Get route definitions from your Context
    // Structure expected: [{ path: '/', component: 'Home', index: true }, ...]
    const { routeDefinitions, isLoading } = useContext(RouteContext);

    // 2. Memoize the router creation so it doesn't recreate on every render
    const router = useMemo(() => {
        if (!routeDefinitions || routeDefinitions.length === 0) return null;

        // Helper to convert config object to React Router Route object
        const mapRoutes = (routes) => {
            return routes.map(route => {
                // Base route object
                const routerRoute = {
                    path: route.path,
                    // Look up the component in the map, default to NotFound if missing
                    element: ComponentMap[route.component] || <NotFound />,
                };

                // Handle index routes
                if (route.index) {
                    routerRoute.index = true;
                    delete routerRoute.path; // Index routes don't have paths
                }

                // Recursively map children if they exist
                if (route.children) {
                    routerRoute.children = mapRoutes(route.children);
                }

                return routerRoute;
            });
        };

        // 3. Define the structure (Root -> Layout -> Dynamic Children)
        const routes = [
            {
                path: '/',
                element: ComponentMap['MainLayout'],
                errorElement: ComponentMap['NotFound'],
                children: [
                    ...mapRoutes(routeDefinitions),
                    // Always append a 404 catch-all at the end
                    { path: '*', element: ComponentMap['NotFound'] }
                ]
            }
        ];

        return createBrowserRouter(routes);

    }, [routeDefinitions]);

    // 4. Handle Loading State
    if (isLoading || !router) {
        return <div>Loading Routes...</div>;
    }

    return <RouterProvider router={router} />;
};

export default DynamicRouter;