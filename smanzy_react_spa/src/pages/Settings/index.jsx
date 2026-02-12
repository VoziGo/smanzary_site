import { useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { Panel, Button, Alert } from "@/components";
import styles from "./index.module.scss";
import api from "@/services/api";
import { useTheme } from "@/context/ThemeContext";
import { getMediaUrl } from "@/utils/fileUtils";
import clsx from "clsx";

export default function Settings() {
    const queryClient = useQueryClient();
    const [syncedData, setSyncedData] = useState(null);
    const [alert, setAlert] = useState(null);
    const { backgroundImage, setBackgroundImage } = useTheme();

    // Fetch images for background selection
    const { isPending, error, data: mediaData } = useQuery({
        queryKey: ["media-images"],
        queryFn: () =>
            api.get(`/media?limit=50&type=image`).then((res) => res.data),
    });

    const isMediaLoading = isPending;
    const mediaError = error;
    const images = mediaData?.data?.files || [];

    // Fetch current background setting
    const { data: settingData } = useQuery({
        queryKey: ["site-bg-setting"],
        queryFn: () => api.get("/settings/site-bg-image").then((res) => res.data),
    });
    const currentBgId = settingData?.data?.value;

    // Mutation to update site settings on server
    const updateSettingMutation = useMutation({
        mutationFn: (value) =>
            api.put("/settings/site-bg-image", { value: String(value) }),
        onSuccess: (res) => {
            queryClient.invalidateQueries({ queryKey: ["site-bg-setting"] });
            const val = res.data?.data?.value;
            if (val) {
                // Update local context state with the direct endpoint URL
                const baseUrl = api.defaults.baseURL;
                setBackgroundImage(`${baseUrl}/site-background?v=${Date.now()}`);
            } else {
                setBackgroundImage(null);
            }
        },
        onError: (err) => {
            console.error(err);
            alert("Failed to update site background: " + (err.response?.data?.error || err.message));
        }
    });

    // sync video database
    const syncVideoDatabase = () => {
        api.post("/videos/sync").then((res) => {
            setSyncedData(res.data);
            setAlert({
                variant: 'success',
                message: `${res.data.fetched} - ${res.data.message}`
            });
        });
    };

    const handleSelectBackground = (image) => {
        updateSettingMutation.mutate(image.id);
    };

    const handleClearBackground = () => {
        updateSettingMutation.mutate("");
    };

    return (
        <div className={styles.container}>
            <Panel title="Settings">
                <div>
                    <h2>Site Settings</h2>
                    <p>Manage site settings here</p>
                    <hr className={styles.lineBreak} />

                    <ul className={styles.list}>
                        <li className={styles.listItem}>
                            <h3>Sync video database</h3>
                            <Button variant="primary" onClick={syncVideoDatabase}>Sync</Button>
                        </li>
                        {alert && (
                            <Alert
                                variant={alert.variant}
                                dismissible
                                onClose={() => setAlert(null)}
                            >
                                {alert.message}
                            </Alert>
                        )}
                    </ul>

                    <div className={styles.section}>
                        <h3>Choose Site Background</h3>
                        <p className={styles.textSecondary}>Select an uploaded image to as your site background</p>

                        {isMediaLoading ? (
                            <div className={styles.loading}>Loading images...</div>
                        ) : mediaError ? (
                            <div className={styles.error}>Error loading images: {mediaError.message}</div>
                        ) : images.length === 0 ? (
                            <div className={styles.loading}>No images found. Upload some in Media Manager first.</div>
                        ) : (
                            <div className={styles.bgGrid}>
                                {images.filter(img => img.mime_type?.startsWith('image/')).map((image) => {
                                    const url = getMediaUrl(image);
                                    return (
                                        <div
                                            key={image.id}
                                            className={clsx(styles.bgItem, String(currentBgId) === String(image.id) && styles.selected)}
                                            onClick={() => handleSelectBackground(image)}
                                            title={image.filename}
                                        >
                                            <img src={url} alt={image.filename} />
                                        </div>
                                    );
                                })}
                            </div>
                        )}

                        <Button
                            variant="secondary"
                            className={styles.clearBtn}
                            onClick={handleClearBackground}
                            disabled={!backgroundImage}
                        >
                            Reset to Default Background
                        </Button>
                    </div>
                </div>
            </Panel>
        </div>
    );
}

