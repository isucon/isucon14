USE isuride;

DELIMITER //

CREATE TRIGGER update_statistics_after_insert
AFTER INSERT ON chair_locations
FOR EACH ROW
BEGIN
    -- NEW.chair_id のみを対象に排他ロックを取得し、total_distance を計算して statistics テーブルを更新
    INSERT INTO statistics (chair_id, total_distance, last_chair_location_id)
    SELECT
        NEW.chair_id,
        COALESCE(SUM(ABS(cl1.latitude - cl2.latitude) + ABS(cl1.longitude - cl2.longitude)), 0) AS total_distance,
        (SELECT id FROM chair_locations cl3 WHERE cl3.chair_id = NEW.chair_id ORDER BY created_at DESC LIMIT 1) AS last_chair_location_id
    FROM
        chair_locations cl1
    LEFT JOIN
        chair_locations cl2 ON cl1.chair_id = cl2.chair_id
        AND cl2.created_at = (
            SELECT MIN(created_at)
            FROM chair_locations
            WHERE chair_id = cl1.chair_id AND created_at > cl1.created_at
        )
    WHERE
        cl1.chair_id = NEW.chair_id
    GROUP BY
        cl1.chair_id
    -- FOR UPDATE
    ON DUPLICATE KEY UPDATE
        total_distance = VALUES(total_distance),
        last_chair_location_id = VALUES(last_chair_location_id);
END //

DELIMITER ;
