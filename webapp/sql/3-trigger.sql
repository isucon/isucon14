USE isuride;

DELIMITER //

CREATE TRIGGER update_statistics
AFTER INSERT ON chair_locations
FOR EACH ROW
BEGIN
  -- Dispatching rangesの計算
  INSERT INTO statistics (chair_id, dispatch_total_distance, dispatch_total_time, carry_total_distance, carry_total_time)
  SELECT
    NEW.chair_id,
    COALESCE(SUM(ABS(loc.latitude - prev_loc.latitude) + ABS(loc.longitude - prev_loc.longitude)), 0) AS dispatch_total_distance,
    COALESCE(SUM(TIMESTAMPDIFF(MICROSECOND, prev_loc.created_at, loc.created_at)), 0) AS dispatch_total_time,
    0, 0
  FROM chair_locations loc
  JOIN chair_locations prev_loc ON loc.chair_id = prev_loc.chair_id
    AND loc.created_at > prev_loc.created_at
  JOIN ride_requests r ON r.chair_id = loc.chair_id
  JOIN ride_request_statuses s ON r.id = s.ride_request_id
  WHERE loc.chair_id = NEW.chair_id
    AND s.status = 'DISPATCHED'
    AND prev_loc.created_at = (
      SELECT MAX(p.created_at)
      FROM chair_locations p
      WHERE p.chair_id = loc.chair_id AND p.created_at < loc.created_at
    )
    AND (SELECT prev_s.status
         FROM ride_requests r
         JOIN ride_request_statuses prev_s ON prev_s.ride_request_id = r.id
         WHERE r.chair_id = loc.chair_id
         AND prev_s.created_at < s.created_at
         ORDER BY prev_s.created_at DESC LIMIT 1) = 'DISPATCHING'
  GROUP BY loc.chair_id
  ON DUPLICATE KEY UPDATE
    dispatch_total_distance = VALUES(dispatch_total_distance),
    dispatch_total_time = VALUES(dispatch_total_time);

  -- Carrying rangesの計算
  INSERT INTO statistics (chair_id, dispatch_total_distance, dispatch_total_time, carry_total_distance, carry_total_time)
  SELECT
    NEW.chair_id,
    0, 0,
    COALESCE(SUM(ABS(loc.latitude - prev_loc.latitude) + ABS(loc.longitude - prev_loc.longitude)), 0) AS carry_total_distance,
    COALESCE(SUM(TIMESTAMPDIFF(MICROSECOND, prev_loc.created_at, loc.created_at)), 0) AS carry_total_time
  FROM chair_locations loc
  JOIN chair_locations prev_loc ON loc.chair_id = prev_loc.chair_id
    AND loc.created_at > prev_loc.created_at
  JOIN ride_requests r ON r.chair_id = loc.chair_id
  JOIN ride_request_statuses s ON r.id = s.ride_request_id
  WHERE loc.chair_id = NEW.chair_id
    AND s.status = 'ARRIVED'
    AND prev_loc.created_at = (
      SELECT MAX(p.created_at)
      FROM chair_locations p
      WHERE p.chair_id = loc.chair_id AND p.created_at < loc.created_at
    )
    AND (SELECT prev_s.status
         FROM ride_requests r
         JOIN ride_request_statuses prev_s ON prev_s.ride_request_id = r.id
         WHERE r.chair_id = loc.chair_id
         AND prev_s.created_at < s.created_at
         ORDER BY prev_s.created_at DESC LIMIT 1) = 'CARRYING'
  GROUP BY loc.chair_id
  ON DUPLICATE KEY UPDATE
    carry_total_distance = VALUES(carry_total_distance),
    carry_total_time = VALUES(carry_total_time);
END //

DELIMITER ;
