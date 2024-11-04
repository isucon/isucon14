DELIMITER //

CREATE TRIGGER update_statistics
AFTER INSERT ON chair_locations
FOR EACH ROW
BEGIN
  -- status_changes サブクエリ
  DECLARE dispatching_time DATETIME;
  DECLARE dispatched_time DATETIME;
  DECLARE carrying_time DATETIME;
  DECLARE arrived_time DATETIME;

  -- dispatching_ranges サブクエリ
  SET dispatching_time = (
    SELECT s.prev_created_at
    FROM (SELECT chair_id, status, created_at,
                 LAG(created_at) OVER (PARTITION BY chair_id ORDER BY created_at) AS prev_created_at,
                 LAG(status) OVER (PARTITION BY chair_id ORDER BY created_at) AS prev_status
          FROM ride_requests
          WHERE chair_id = NEW.chair_id) s
    WHERE s.prev_status = 'dispatching' AND s.status = 'dispatched'
  );
  SET dispatched_time = (
    SELECT s.created_at
    FROM (SELECT chair_id, status, created_at,
                 LAG(created_at) OVER (PARTITION BY chair_id ORDER BY created_at) AS prev_created_at,
                 LAG(status) OVER (PARTITION BY chair_id ORDER BY created_at) AS prev_status
          FROM ride_requests
          WHERE chair_id = NEW.chair_id) s
    WHERE s.prev_status = 'dispatching' AND s.status = 'dispatched'
  );

  -- carrying_ranges サブクエリ
  SET carrying_time = (
    SELECT s.prev_created_at
    FROM (SELECT chair_id, status, created_at,
                 LAG(created_at) OVER (PARTITION BY chair_id ORDER BY created_at) AS prev_created_at,
                 LAG(status) OVER (PARTITION BY chair_id ORDER BY created_at) AS prev_status
          FROM ride_requests
          WHERE chair_id = NEW.chair_id) s
    WHERE s.prev_status = 'carrying' AND s.status = 'arrived'
  );
  SET arrived_time = (
    SELECT s.created_at
    FROM (SELECT chair_id, status, created_at,
                 LAG(created_at) OVER (PARTITION BY chair_id ORDER BY created_at) AS prev_created_at,
                 LAG(status) OVER (PARTITION BY chair_id ORDER BY created_at) AS prev_status
          FROM ride_requests
          WHERE chair_id = NEW.chair_id) s
    WHERE s.prev_status = 'carrying' AND s.status = 'arrived'
  );

  -- dispatching_stats 集計
  INSERT INTO statistics (chair_id, dispatch_total_distance, dispatch_total_time, carry_total_distance, carry_total_time)
  SELECT
    NEW.chair_id,
    COALESCE(SUM(ABS(loc.latitude - prev_latitude) + ABS(loc.longitude - prev_longitude)), 0),
    COALESCE(SUM(TIMESTAMPDIFF(MICROSECOND, prev_created_at, loc.created_at)), 0),
    0,
    0
  FROM (
    SELECT loc.chair_id,
           loc.latitude,
           loc.longitude,
           loc.created_at,
           LAG(loc.latitude) OVER (PARTITION BY loc.chair_id ORDER BY loc.created_at) AS prev_latitude,
           LAG(loc.longitude) OVER (PARTITION BY loc.chair_id ORDER BY loc.created_at) AS prev_longitude,
           LAG(loc.created_at) OVER (PARTITION BY loc.chair_id ORDER BY loc.created_at) AS prev_created_at
    FROM chair_locations loc
    WHERE loc.created_at BETWEEN dispatching_time AND dispatched_time
  ) loc
  GROUP BY loc.chair_id
  ON DUPLICATE KEY UPDATE
    dispatch_total_distance = VALUES(dispatch_total_distance),
    dispatch_total_time = VALUES(dispatch_total_time);

  -- carrying_stats 集計
  INSERT INTO statistics (chair_id, carry_total_distance, carry_total_time)
  SELECT
    NEW.chair_id,
    COALESCE(SUM(ABS(loc.latitude - prev_latitude) + ABS(loc.longitude - prev_longitude)), 0),
    COALESCE(SUM(TIMESTAMPDIFF(MICROSECOND, prev_created_at, loc.created_at)), 0)
  FROM (
    SELECT loc.chair_id,
           loc.latitude,
           loc.longitude,
           loc.created_at,
           LAG(loc.latitude) OVER (PARTITION BY loc.chair_id ORDER BY loc.created_at) AS prev_latitude,
           LAG(loc.longitude) OVER (PARTITION BY loc.chair_id ORDER BY loc.created_at) AS prev_longitude,
           LAG(loc.created_at) OVER (PARTITION BY loc.chair_id ORDER BY loc.created_at) AS prev_created_at
    FROM chair_locations loc
    WHERE loc.created_at BETWEEN carrying_time AND arrived_time
  ) loc
  GROUP BY loc.chair_id
  ON DUPLICATE KEY UPDATE
    carry_total_distance = VALUES(carry_total_distance),
    carry_total_time = VALUES(carry_total_time);
END //

DELIMITER ;
