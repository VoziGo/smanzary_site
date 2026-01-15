import { Panel } from "@/components";
import styles from "./index.module.scss";

export default function Settings() {
    return (
        <div className={styles.container}>
            <Panel title="Settings">
                <div>
                    <h2>Settings</h2>
                    <ul className={styles.list}>
                        <li>Under Construction!</li>
                    </ul>
                </div>
            </Panel>
        </div>
    );
}
