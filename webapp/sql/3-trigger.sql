DELIMITER //

CREATE TRIGGER update_statistics
AFTER INSERT ON chair_locations
FOR EACH ROW
BEGIN
  -- dispatching -> dispatched サイクルの開始時刻と終了時刻を取得
  WITH status_cycles AS (
    SELECT chair_id,
           status,
           created_at,
           LAG(created_at) OVER (PARTITION BY chair_id ORDER BY created_at) AS prev_time,
           LAG(status) OVER (PARTITION BY chair_id ORDER BY created_at) AS prev_status
    FROM ride_requests
    WHERE chair_id = NEW.chair_id
  ),

  dispatching_cycles AS (
    SELECT s1.created_at AS dispatching_time,
           s2.created_at AS dispatched_time
    FROM status_cycles s1
    JOIN status_cycles s2 ON s1.chair_id = s2.chair_id
    WHERE s1.status = 'dispatching'
      AND s2.status = 'dispatched'
      AND s1.created_at < s2.created_at
  ),

  carrying_cycles AS (
    SELECT s1.created_at AS carrying_time,
           s2.created_at AS arrived_time
    FROM status_cycles s1
    JOIN status_cycles s2 ON s1.chair_id = s2.chair_id
    WHERE s1.status = 'carrying'
      AND s2.status = 'arrived'
      AND s1.created_at < s2.created_at
  ),

  -- 各 dispatching -> dispatched サイクルごとの距離と時間を計算
  dispatching_stats AS (
    SELECT
      SUM((loc2.latitude - loc1.latitude) - (loc2.longitude - loc1.longitude)) AS dispatching_total_distance,
      SUM(TIMESTAMPDIFF(SECOND, loc1.created_at, loc2.created_at)) AS dispatching_total_times
    FROM chair_locations loc1
    JOIN chair_locations loc2 ON loc1.chair_id = loc2.chair_id
    JOIN dispatching_cycles dc ON loc1.created_at <= dc.dispatching_time
                               AND loc2.created_at BETWEEN dc.dispatching_time AND dc.dispatched_time
    WHERE loc1.chair_id = NEW.chair_id
  ),

  -- 各 carrying -> arrived サイクルごとの距離と時間を計算
  carrying_stats AS (
    SELECT
      SUM((loc2.latitude - loc1.latitude) - (loc2.longitude - loc1.longitude)) AS carrying_total_distance,
      SUM(TIMESTAMPDIFF(SECOND, loc1.created_at, loc2.created_at)) AS carrying_total_times
    FROM chair_locations loc1
    JOIN chair_locations loc2 ON loc1.chair_id = loc2.chair_id
    JOIN carrying_cycles cc ON loc1.created_at <= cc.carrying_time
                             AND loc2.created_at BETWEEN cc.carrying_time AND cc.arrived_time
    WHERE loc1.chair_id = NEW.chair_id
  )

  -- 結果を statistics テーブルに挿入または更新
  INSERT INTO statistics (chair_id, dispatching_total_distance, dispatching_total_times, carrying_total_distance, carrying_total_times)
  SELECT
    NEW.chair_id,
    COALESCE(dispatching_total_distance, 0),
    COALESCE(dispatching_total_times, 0),
    COALESCE(carrying_total_distance, 0),
    COALESCE(carrying_total_times, 0)
  FROM dispatching_stats, carrying_stats
  ON DUPLICATE KEY UPDATE
    dispatching_total_distance = VALUES(dispatching_total_distance),
    dispatching_total_times = VALUES(dispatching_total_times),
    carrying_total_distance = VALUES(carrying_total_distance),
    carrying_total_times = VALUES(carrying_total_times);
END //

DELIMITER ;
