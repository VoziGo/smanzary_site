import { useState, useEffect } from "react";
import VideoCard from "@/components/VideoCard";
import Pagination from "@/components/Pagination";
import api from '@/services/api';

import styles from "./index.module.scss";

export default function Home() {
    const [videos, setVideos] = useState([]);
    const [loading, setLoading] = useState(true);
    const [page, setPage] = useState(1);
    const [totalPages, setTotalPages] = useState(0);
    const limit = 6;

    useEffect(() => {
        const controller = new AbortController();
        fetchVideos(page, controller.signal);
        return () => controller.abort();
    }, [page]);

    async function fetchVideos(pageNum, signal) {
        setLoading(true);
        try {
            const response = await api.get(`/videos?page=${pageNum}&limit=${limit}`, { signal });
            const data = response.data;

            // Map backend data to match VideoCard expectations
            // Backend sends: video_id, title, views (int), etc.
            // VideoCard expects: id, title, views (string/int)
            const mappedVideos = data.videos
                ? data.videos.map((v) => ({
                    id: v.video_id,
                    title: v.title,
                    views: v.views
                        ? `${Number(v.views).toLocaleString()} views`
                        : "Views hidden",
                    ...v,
                }))
                : [];

            setVideos(mappedVideos);


            // Calculate total pages based on total count if available, otherwise just use logic
            if (data.total) {
                setTotalPages(Math.ceil(data.total / limit));
            } else {
                setTotalPages(1);
            }
        } catch (error) {
            if (api.isCancel(error) || error.name === 'CanceledError' || error.name === 'AbortError') {
                console.log('Request canceled');
            } else {
                console.error("Error fetching videos:", error);
                // Only clear videos on error if we want to show an error state, 
                // but keeping previous videos might be better UX depending on error.
                // For now, let's keep behavior simple but maybe not clear immediately on temporary errors?
                // existing logic cleared it: setVideos([]); 
                // Let's keep it defined to avoid breaking changes, but maybe only if it's not a cancel
                setVideos([]);
            }
        } finally {
            // Only turn off loading if not aborted (though aborted requests throw, so caught above)
            // If we use signal.aborted check:
            if (!signal?.aborted) {
                setLoading(false);
            }
        }
    }

    const handlePageChange = (newPage) => {
        setPage(newPage);
        window.scrollTo({ top: 0, behavior: "smooth" });
    };

    return (
        <div className={styles.home}>
            <section className={styles.videosSection}>
                <div className={`${styles.videosGrid} ${loading && videos.length > 0 ? styles.gridLoading : ''}`}>
                    {videos.length > 0 ? (
                        videos.map((video) => <VideoCard key={video.id} video={video} />)
                    ) : loading ? (
                        <p className={styles.loading}>Loading videos...</p>
                    ) : (
                        <p className={styles.noVideos}>No videos found.</p>
                    )}
                </div>

                {/* Pagination Controls */}
                {!loading && videos.length > 0 && (
                    <Pagination
                        currentPage={page}
                        totalPages={totalPages}
                        onPageChange={handlePageChange}
                        maxWidth="1300px"
                    />
                )}
            </section>
        </div>
    );
}
