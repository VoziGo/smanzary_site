import styles from "./index.module.scss";
import { Edit, Trash2 } from "lucide-react";
import Button from "@/components/Button";
import { getMediaUrl, getThumbnailUrl } from "@/utils/fileUtils";
import { useQuery } from "@tanstack/react-query";
import api from "@/services/api";

export default function AlbumCard({ album, onManage, onDelete, isDeleting }) {
  const { data: mediaRows } = useQuery({
    queryKey: ["albums", album.id, "media"],
    queryFn: () => api.get(`/media/album/${album.id}`).then((res) => res.data),
    enabled: !!album.id,
  });

  const albumMedia = mediaRows?.data || [];
  const coverImage = albumMedia[0];

  return (
    <div className={styles.albumCard}>
      <div className={styles.albumCover} onClick={onManage}>
        {coverImage && coverImage.mime_type.startsWith("image/") ? (
          <img src={getMediaUrl(coverImage)} alt={album.title} />
        ) : (
          <div className={styles.placeholderCover}>
            <Edit size={48} />
          </div>
        )}
        <div className={styles.mediaCountBadge}>{albumMedia?.length || 0}</div>
        {album.user_name && (
          <div className={styles.userNameOverlay}>{album.user_name}</div>
        )}
      </div>

      <div className={styles.albumInfo}>
        <h3 className={styles.albumTitle} onClick={onManage}>
          {album.title}
        </h3>
        {album.description && (
          <p className={styles.albumDescription}>{album.description}</p>
        )}
      </div>

      <div className={styles.albumActions}>
        <Button
          onClick={onManage}
          variant="primary"
          className={styles.manageBtn}
        >
          <Edit size={16} />
          Manage
        </Button>
        <Button
          onClick={(e) => {
            e.stopPropagation();
            if (confirm("Are you sure you want to delete this album?")) {
              onDelete(album.id);
            }
          }}
          disabled={isDeleting}
          variant="danger"
        >
          <Trash2 size={16} />
        </Button>
      </div>
    </div>
  );
}
