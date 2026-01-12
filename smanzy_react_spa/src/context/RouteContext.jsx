import React, { createContext, useState, useEffect } from 'react';

export const RouteContext = createContext();

export const RouteProvider = ({ children }) => {
    // Initial state: You might want to start with public routes only
    const [routes, setRoutes] = useState([]);
    const [isLoading, setIsLoading] = useState(true);

    useEffect(() => {
        // SIMULATION: In a real app, this is where you fetch from an API
        // or check localStorage for a role to decide which JSON to load.
        const loadRoutes = async () => {
            setIsLoading(true);

            // Simulating API delay
            await new Promise(r => setTimeout(r, 500));

            // The data structure you would get from your DB/API
            const appRoutes = [
                { index: true, component: 'Home', title: 'Home' },
                { path: 'videos', component: 'Videos', title: 'Videos' },
                { path: 'about', component: 'About', title: 'About' },
                { path: 'login', component: 'Login', title: 'Login' },
                { path: 'register', component: 'Register', title: 'Register' },
                { path: 'profile', component: 'Profile', title: 'Profile' },
                { path: 'media', component: 'MediaManager', title: 'Media Manager' },
                { path: 'media/edit/:id', component: 'UpdateMedia', title: 'Update Media' },
                { path: 'mediacards', component: 'MediaManagerCards', title: 'Media Manager Cards' },
                { path: 'albums', component: 'AlbumList', title: 'Album List' },
                { path: 'albums/:id', component: 'AlbumDetail', title: 'Album Detail' },
                { path: 'users', component: 'UserManagement', title: 'User Management' },
                { path: 'settings', component: 'Settings', title: 'Settings' },
                { path: '*', component: 'NotFound', title: 'Not Found' },
            ];

            setRoutes(appRoutes);
            setIsLoading(false);
        };

        loadRoutes();
    }, []);

    return (
        <RouteContext.Provider value={{ routes, isLoading, setRoutes }}>
            {children}
        </RouteContext.Provider>
    );
};