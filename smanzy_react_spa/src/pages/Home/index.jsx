import { useState, useEffect } from "react";
import VideoCard from "@/components/VideoCard";
import styles from "./index.module.scss";

const API_BASE_URL = "http://localhost:8080/api";

export default function Home() {
    const [videos, setVideos] = useState([]);
    const [loading, setLoading] = useState(true);
    const [page, setPage] = useState(1);
    const [totalPages, setTotalPages] = useState(0);
    const limit = 9; // 3x3 grid feels nice

    useEffect(() => {
        fetchVideos(page);
    }, [page]);

    async function fetchVideos(pageNum) {
        setLoading(true);
        try {
            const response = await fetch(`${API_BASE_URL}/videos?page=${pageNum}&limit=${limit}`);
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            const data = await response.json();

            // Map backend data to match VideoCard expectations
            // Backend sends: video_id, title, views (int), etc.
            // VideoCard expects: id, title, views (string/int)
            const mappedVideos = data.videos ? data.videos.map((v) => ({
                id: v.video_id,
                title: v.title,
                views: v.views ? `${Number(v.views).toLocaleString()} views` : "Views hidden",
                ...v,
            })) : [];

            setVideos(mappedVideos);

            // Calculate total pages based on total count if available, otherwise just use logic
            if (data.total) {
                setTotalPages(Math.ceil(data.total / limit));
            } else {
                setTotalPages(1);
            }
        } catch (error) {
            console.error("Error fetching videos:", error);
            setVideos([]);
        } finally {
            setLoading(false);
        }
    }

    const handlePrevPage = () => {
        if (page > 1) setPage((p) => p - 1);
    };

    const handleNextPage = () => {
        if (page < totalPages) setPage((p) => p + 1);
    };

    return (
        <div className={styles.home}>
            <header className={styles.header}>
                <div className={styles.logo}>SmAnZaRy</div>
                <h2 className={styles.channelName}>
                    Welcome to SmAnZaRy YouTube Channel
                </h2>
                <p className={styles.description}>
                    Cozy comfort escape environments! <br />
                    Relax with ambient jazz piano, crackling fireplaces, snowy nights, and
                    peaceful Christmas vibes. <br />
                    Perfect for focus, study, relaxation, and calm evenings.
                </p>
                <a
                    href="https://www.youtube.com/@smanzary?sub_confirmation=1"
                    className={styles.subscribeBtn}
                    target="_blank"
                    rel="noopener noreferrer"
                >
                    Subscribe Now
                </a>
            </header>

            <section className={styles.videosSection}>
                {/* <h2>Featured Videos</h2>*/}
                <div className={styles.videosGrid}>
                    {loading ? (
                        <p className={styles.loading}>Loading videos...</p>
                    ) : videos.length > 0 ? (
                        videos.map((video) => <VideoCard key={video.id} video={video} />)
                    ) : (
                        <p className={styles.noVideos}>No videos found.</p>
                    )}
                </div>

                {/* Pagination Controls */}
                {!loading && videos.length > 0 && (
                    <div className={styles.pagination}>
                        <button
                            onClick={handlePrevPage}
                            disabled={page === 1}
                            className={styles.pageBtn}
                        >
                            Previous
                        </button>
                        <span className={styles.pageInfo}>
                            Page {page} of {totalPages}
                        </span>
                        <button
                            onClick={handleNextPage}
                            disabled={page === totalPages}
                            className={styles.pageBtn}
                        >
                            Next
                        </button>
                    </div>
                )}

                <p style={{ marginTop: "40px" }}>
                    <a
                        href="https://www.youtube.com/@smanzary/videos"
                        target="_blank"
                        rel="noopener noreferrer"
                    >
                        View All Videos on YouTube
                    </a>
                </p>
            </section>
        </div>
    );
}
