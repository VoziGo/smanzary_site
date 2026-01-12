import React, { useContext, useMemo } from 'react';
import { createBrowserRouter, RouterProvider } from 'react-router-dom';
import { RouteContext } from '@/context/RouteContext';
import { ComponentMap } from './componentMap';
import { NotFound } from '@/pages'; // Fallback import

const AppRouter = () => {
    const { routes, isLoading } = useContext(RouteContext);

    const router = useMemo(() => {
        // While loading or if no routes, return null or a basic loader router
        if (isLoading || !routes.length) return null;

        // Recursive function to map JSON routes to Router objects
        const generateRoutes = (routeList) => {
            return routeList.map((route) => {
                const routeObj = {
                    element: ComponentMap[route.component] || <NotFound />,
                };

                // Handle index vs path
                if (route.index) {
                    routeObj.index = true;
                } else {
                    routeObj.path = route.path;
                }

                // Handle nested children
                if (route.children) {
                    routeObj.children = generateRoutes(route.children);
                }

                return routeObj;
            });
        };

        return createBrowserRouter([
            {
                path: '/',
                element: ComponentMap['MainLayout'],
                errorElement: ComponentMap['NotFound'],
                children: [
                    ...generateRoutes(routes),
                    // Catch-all must be last
                    { path: '*', element: ComponentMap['NotFound'] }
                ],
            },
        ]);
    }, [routes, isLoading]);

    if (isLoading) return <div>Loading Application...</div>;
    if (!router) return <div>Error loading routes.</div>;

    return <RouterProvider router={router} />;
};

export default AppRouter;