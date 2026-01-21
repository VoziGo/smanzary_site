# Album Management Frontend

## Overview

A complete album management interface has been added to the React SPA frontend for organizing and managing media collections.

## New Files Created

### 1. **AlbumList Page Component**
- **Location**: `src/pages/AlbumList/index.jsx`
- **Features**:
  - View all user albums in a grid
  - Create new albums with title and description
  - Quick access to album management
  - Delete albums (soft delete)

### 2. **AlbumDetail Page Component**
- **Location**: `src/pages/AlbumDetail/index.jsx`
- **Features**:
  - View specific album details
  - Add media files from the global library
  - Remove media files from the album
  - Update album title and description
  - Optimized media fetching (independent request from album metadata)

### 3. **AlbumCard Component**
- **Location**: `src/components/AlbumCard/index.jsx`
- **Features**:
  - Reusable card for displaying album info
  - Dynamic cover image fetching (uses the first media item in the album)
  - Accurate media count badge powered by dedicated queries

### 2. **AlbumManager Styles**
- **Location**: `src/pages/AlbumList/index.module.scss`, `src/pages/AlbumDetail/index.module.scss`
- **Features**:
  - Responsive grid layout
  - Dark/light theme support using CSS variables
  - Mobile-friendly design
  - Smooth transitions and hover effects
  - Card-based UI for albums and media

## Files Modified

### 1. **src/routes/index.jsx**
- Added AlbumManager to imports
- Added new route: `/albums` -> AlbumManager component

### 2. **src/pages/index.jsx**
- Exported AlbumList and AlbumDetail components

### 3. **src/components/Navbar/index.jsx**
- Added "Albums" navigation link (visible when logged in)
- Added to both desktop and mobile navigation
- Route: `/albums` (List), `/albums/:id` (Detail)

## Component Architecture

### AlbumList & AlbumDetail Components

**Decoupled Architecture**:
- Album metadata (Title, Description) and Media Content are fetched in separate requests.
- This allows for faster initial page loads and more efficient cache management.

**React Query Hooks**:
- `useQuery(['albums'])` - Fetch list of user albums
- `useQuery(['albums', id])` - Fetch single album metadata
- `useQuery(['albums', id, 'media'])` - Fetch media files for specific album
- `useMutation` - Create, update, delete, add/remove media operations
- `useQueryClient` - Granular cache invalidation (e.g., only updating media list without refetching album info)

**Features**:
1. **Create Album**
   - Form with title (required) and description (optional)
   - Form validation
   - Loading state during submission
   - Cache invalidation after success

2. **View Albums**
   - List all user albums
   - Display album info (title, description, media count)
   - Expandable details view
   - Empty state message

3. **Album Details**
   - View media files in album
   - Search and add media
   - Remove media from album
   - Loading states for operations

4. **Media Management**
   - Search media files by filename
   - Add media with visual feedback
   - Remove media with confirmation
   - Display media in grid layout

## API Integration

### Endpoints Used

**Albums**:
- `GET /albums` - Get user's albums (metadata only)
- `POST /albums` - Create new album
- `GET /albums/:id` - Get album details (metadata only)
- `PUT /albums/:id` - Update album
- `DELETE /albums/:id` - Delete album
- `POST /albums/:id/media` - Add media to album
- `DELETE /albums/:id/media` - Remove media from album

**Media**:
- `GET /media/album/:id` - List media files for a specific album (Optimized)
- `GET /media?limit=1000` - Get all media files (for adding to albums)

## UI/UX Features

### Album Card
- Title and description display
- Media count indicator
- View/Hide details toggle
- Delete button

### Album Details Section
- Media grid display
- Add media button with search
- Remove media buttons for each file
- Responsive grid layout

### Create Album Form
- Title input field (required)
- Description textarea (optional)
- Submit button with loading state
- Cancel button

### Add Media Section
- Search/filter input
- Scrollable media list
- Add buttons for each media file
- Visual feedback during operation

## Styling

### CSS Variables Used
- `--bg-primary` - Primary background
- `--bg-secondary` - Secondary background
- `--text-primary` - Primary text
- `--text-secondary` - Secondary text
- `--text-tertiary` - Tertiary text
- `--border-color` - Border color
- `--accent-color` - Accent/highlight color
- `--shadow-color` - Shadow color

### Responsive Breakpoints
- Mobile: < 768px
- Tablet: ≥ 768px
- Desktop: ≥ 1024px

### Classes
- `.container` - Main container with max-width
- `.albumsList` - Flex column layout for albums
- `.albumCard` - Individual album card
- `.mediaGrid` - Responsive grid for media items
- `.createForm` - Album creation form styling
- `.mediaList` - Scrollable media list

## Usage

### For Users
1. Navigate to "Albums" in the navbar (when logged in)
2. Click "New Album" to create an album
3. Enter title and optional description
4. Click "Create Album"
5. Click "View" on an album to see details
6. Click "+" to add media to the album
7. Search for media files and click "Add"
8. Remove media by clicking the "X" button
9. Delete album by clicking the trash icon

### For Developers
```jsx
import AlbumManager from '@/pages/AlbumManager';

// Component is already integrated in routes
// Access via /albums route when authenticated
```

## Error Handling

- Network errors displayed via alerts
- Loading states for all operations
- Cache invalidation for consistency
- Form validation for required fields

## Future Enhancements

- Edit album title/description
- Bulk operations (add multiple media at once)
- Album sharing/permissions
- Album cover image selection
- Album sorting and filtering
- Download entire album
- Album statistics (total size, media count)
