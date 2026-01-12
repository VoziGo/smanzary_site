import MainLayout from '@/layout/MainLayout';
import {
    Home, About, Login, Register, Profile, MediaManager,
    UpdateMedia, NotFound, MediaManagerCards, AlbumList,
    AlbumDetail, Videos, UserManagement
} from '@/pages';

// Map string keys to actual components
export const ComponentMap = {
    MainLayout: <MainLayout />,
    Home: <Home />,
    Videos: <Videos />,
    About: <About />,
    Login: <Login />,
    Register: <Register />,
    Profile: <Profile />,
    MediaManager: <MediaManager />,
    UpdateMedia: <UpdateMedia />,
    MediaManagerCards: <MediaManagerCards />,
    AlbumList: <AlbumList />,
    AlbumDetail: <AlbumDetail />,
    UserManagement: <UserManagement />,
    NotFound: <NotFound />
};