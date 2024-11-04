DELIMITER //

CREATE TRIGGER update_statistics
AFTER INSERT ON chair_locations
FOR EACH ROW
BEGIN
  WITH status_changes AS (
    SELECT chair_id,
           status,
           created_at,
           LAG(created_at) OVER (PARTITION BY chair_id ORDER BY created_at) AS prev_created_at,
           LAG(status) OVER (PARTITION BY chair_id ORDER BY created_at) AS prev_status
    FROM ride_requests
    WHERE chair_id = NEW.chair_id
  ),

  dispatching_ranges AS (
    SELECT s.prev_created_at AS dispatching_time,
           s.created_at AS dispatched_time
    FROM status_changes s
    WHERE s.prev_status = 'dispatching'
      AND s.status = 'dispatched'
  ),

  carrying_ranges AS (
    SELECT s.prev_created_at AS carrying_time,
           s.created_at AS arrived_time
    FROM status_changes s
    WHERE s.prev_status = 'carrying'
      AND s.status = 'arrived'
  ),

  dispatch_moves AS (
    SELECT
      loc.chair_id,
      loc.latitude,
      loc.longitude,
      loc.created_at,
      LAG(loc.latitude) OVER (PARTITION BY loc.chair_id ORDER BY loc.created_at) AS prev_latitude,
      LAG(loc.longitude) OVER (PARTITION BY loc.chair_id ORDER BY loc.created_at) AS prev_longitude,
      LAG(loc.created_at) OVER (PARTITION BY loc.chair_id ORDER BY loc.created_at) AS prev_created_at
    FROM chair_locations loc
    JOIN dispatching_ranges dr ON loc.prev_created_at > dr.dispatching_time AND loc.created_at <= dr.dispatched_time
  ),

  carry_moves AS (
    SELECT
      loc.chair_id,
      loc.latitude,
      loc.longitude,
      loc.created_at,
      LAG(loc.latitude) OVER (PARTITION BY loc.chair_id ORDER BY loc.created_at) AS prev_latitude,
      LAG(loc.longitude) OVER (PARTITION BY loc.chair_id ORDER BY loc.created_at) AS prev_longitude,
      LAG(loc.created_at) OVER (PARTITION BY loc.chair_id ORDER BY loc.created_at) AS prev_created_at
    FROM chair_locations loc
    JOIN carrying_ranges cr ON loc.prev_created_at > cr.carrying_time AND loc.created_at <= cr.arrived_time
  ),

  dispatching_stats AS (
    SELECT
      dm.chair_id,
      SUM(ABS(dm.latitude - dm.prev_latitude) + ABS(dm.longitude - dm.prev_longitude)) AS dispatch_total_distance,
      SUM(TIMESTAMPDIFF(MICROSECOND, dm.prev_created_at, dm.created_at)) AS dispatch_total_time
    FROM dispatch_moves dm
    GROUP BY dm.chair_id
  ),

  carrying_stats AS (
    SELECT
      cm.chair_id,
      SUM(ABS(cm.latitude - cm.prev_latitude) + ABS(cm.longitude - cm.prev_longitude)) AS carry_total_distance,
      SUM(TIMESTAMPDIFF(MICROSECOND, cm.prev_created_at, cm.created_at)) AS carry_total_time
    FROM carry_moves cm
    GROUP BY cm.chair_id
  )

  -- 結果を statistics テーブルに挿入または更新
  INSERT INTO statistics (chair_id, dispatch_total_distance, dispatch_total_time, carry_total_distance, carry_total_time)
  SELECT
    NEW.chair_id,
    COALESCE(dispatch_total_distance, 0),
    COALESCE(dispatch_total_time, 0),
    COALESCE(carry_total_distance, 0),
    COALESCE(carry_total_time, 0)
  FROM dispatching_stats
  LEFT JOIN carrying_stats ON dispatching_stats.chair_id = carrying_stats.chair_id
  ON DUPLICATE KEY UPDATE
    dispatch_total_distance = VALUES(dispatch_total_distance),
    dispatch_total_time = VALUES(dispatch_total_time),
    carry_total_distance = VALUES(carry_total_distance),
    carry_total_time = VALUES(carry_total_time);
END //

DELIMITER ;
