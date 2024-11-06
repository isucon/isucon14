DELIMITER //

CREATE TRIGGER update_statistics
AFTER INSERT ON chair_locations
FOR EACH ROW
BEGIN
  -- Dispatching rangesの計算
  INSERT INTO statistics (chair_id, dispatch_total_distance, dispatch_total_time, carry_total_distance, carry_total_time)
  SELECT
    NEW.chair_id,
    COALESCE(SUM(ABS(loc.latitude - loc.prev_latitude) + ABS(loc.longitude - loc.prev_longitude)), 0) AS dispatch_total_distance,
    COALESCE(SUM(TIMESTAMPDIFF(MICROSECOND, loc.prev_created_at, loc.created_at)), 0) AS dispatch_total_time,
    0, 0
  FROM (
    SELECT
      loc.chair_id,
      loc.latitude,
      loc.longitude,
      loc.created_at,
      LAG(loc.latitude) OVER (PARTITION BY loc.chair_id ORDER BY loc.created_at) AS prev_latitude,
      LAG(loc.longitude) OVER (PARTITION BY loc.chair_id ORDER BY loc.created_at) AS prev_longitude,
      LAG(loc.created_at) OVER (PARTITION BY loc.chair_id ORDER BY loc.created_at) AS prev_created_at
    FROM chair_locations loc
    WHERE loc.chair_id = NEW.chair_id
      AND EXISTS (
        SELECT 1
        FROM ride_requests r
        JOIN ride_request_statuses s ON r.id = s.ride_request_id
        WHERE r.chair_id = loc.chair_id
          AND s.status = 'DISPATCHED'
          AND LAG(s.status) OVER (PARTITION BY s.chair_id ORDER BY s.created_at) = 'DISPATCHING'
      )
  ) AS loc
  GROUP BY loc.chair_id
  ON DUPLICATE KEY UPDATE
    dispatch_total_distance = VALUES(dispatch_total_distance),
    dispatch_total_time = VALUES(dispatch_total_time);

  -- Carrying rangesの計算
  INSERT INTO statistics (chair_id, dispatch_total_distance, dispatch_total_time, carry_total_distance, carry_total_time)
  SELECT
    NEW.chair_id,
    0, 0,
    COALESCE(SUM(ABS(loc.latitude - loc.prev_latitude) + ABS(loc.longitude - loc.prev_longitude)), 0) AS carry_total_distance,
    COALESCE(SUM(TIMESTAMPDIFF(MICROSECOND, loc.prev_created_at, loc.created_at)), 0) AS carry_total_time
  FROM (
    SELECT
      loc.chair_id,
      loc.latitude,
      loc.longitude,
      loc.created_at,
      LAG(loc.latitude) OVER (PARTITION BY loc.chair_id ORDER BY loc.created_at) AS prev_latitude,
      LAG(loc.longitude) OVER (PARTITION BY loc.chair_id ORDER BY loc.created_at) AS prev_longitude,
      LAG(loc.created_at) OVER (PARTITION BY loc.chair_id ORDER BY loc.created_at) AS prev_created_at
    FROM chair_locations loc
    WHERE loc.chair_id = NEW.chair_id
      AND EXISTS (
        SELECT 1
        FROM ride_requests r
        JOIN ride_request_statuses s ON r.id = s.ride_request_id
        WHERE r.chair_id = loc.chair_id
          AND s.status = 'ARRIVED'
          AND LAG(s.status) OVER (PARTITION BY s.chair_id ORDER BY s.created_at) = 'CARRYING'
      )
  ) AS loc
  GROUP BY loc.chair_id
  ON DUPLICATE KEY UPDATE
    carry_total_distance = VALUES(carry_total_distance),
    carry_total_time = VALUES(carry_total_time);
END //

DELIMITER ;
