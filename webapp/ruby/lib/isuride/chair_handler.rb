# frozen_string_literal: true

require 'ulid'

require 'isuride/base_handler'

module Isuride
  class ChairHandler < BaseHandler
    CurrentChair = Data.define(
      :id,
      :owner_id,
      :name,
      :model,
      :is_active,
      :access_token,
      :created_at,
      :updated_at,
    )

    before do
      if request.path == '/api/chair/chairs'
        next
      end

      access_token = cookies[:chair_session]
      if access_token.nil?
        raise HttpError.new(401, 'chair_session cookie is required')
      end
      chair = db.xquery('SELECT * FROM chairs WHERE access_token = ?', access_token).first
      if chair.nil?
        raise HttpError.new(401, 'invalid access token')
      end

      @current_chair = CurrentChair.new(**chair)
    end

    ChairPostChairsRequest = Data.define(:name, :model, :chair_register_token)

    # POST /api/chair/chairs
    post '/chairs' do
      req = bind_json(ChairPostChairsRequest)
      if req.name.nil? || req.model.nil? || req.chair_register_token.nil?
        raise HttpError.new(400, 'some of required fields(name, model, chair_register_token) are empty')
      end

      owner = db.xquery('SELECT * FROM owners WHERE chair_register_token = ?', req.chair_register_token).first
      if owner.nil?
        raise HttpError.new(401, 'invalid chair_register_token')
      end

      chair_id = ULID.generate
      access_token = SecureRandom.hex(32)

      db.xquery('INSERT INTO chairs (id, owner_id, name, model, is_active, access_token) VALUES (?, ?, ?, ?, ?, ?)', chair_id, owner.fetch(:id), req.name, req.model, false, access_token)

      cookies.set(:chair_session, value: access_token, path: '/')
      status(201)
      json(id: chair_id, owner_id: owner.fetch(:id))
    end

    PostChairActivityRequest = Data.define(:is_active)

    # POST /api/chair/activity
    post '/activity' do
      req = bind_json(PostChairActivityRequest)

      db.xquery('UPDATE chairs SET is_active = ? WHERE id = ?', req.is_active, @current_chair.id)

      status(204)
    end

    PostChairCoordinateRequest = Data.define(:latitude, :longitude)

    # POST /api/chair/coordinate
    post '/coordinate' do
      req = bind_json(PostChairCoordinateRequest)

      response = db_transaction do |tx|
        chair_location_id = ULID.generate
        tx.xquery('INSERT INTO chair_locations (id, chair_id, latitude, longitude) VALUES (?, ?, ?, ?)', chair_location_id, @current_chair.id, req.latitude, req.longitude)

        location = tx.xquery('SELECT * FROM chair_locations WHERE id = ?', chair_location_id).first

        ride = tx.xquery('SELECT * FROM rides WHERE chair_id = ? ORDER BY updated_at DESC LIMIT 1', @current_chair.id).first
        unless ride.nil?
          status = get_latest_ride_status(tx, ride.fetch(:id))
          if status != 'COMPLETED' && status != 'CANCELED'
            if req.latitude == ride.fetch(:pickup_latitude) && req.longitude == ride.fetch(:pickup_longitude) && status == 'ENROUTE'
              tx.xquery('INSERT INTO ride_statuses (id, ride_id, status) VALUES (?, ?, ?)', ULID.generate, ride.fetch(:id), 'PICKUP')
            end

            if req.latitude == ride.fetch(:destination_latitude) && req.longitude == ride.fetch(:destination_longitude) && status == 'CARRYING'
              tx.xquery('INSERT INTO ride_statuses (id, ride_id, status) VALUES (?, ?, ?)', ULID.generate, ride.fetch(:id), 'ARRIVED')
            end
          end
        end

        { recorded_at: time_msec(location.fetch(:created_at)) }
      end

      json(response)
    end

    # GET /api/chair/notification
    get '/notification' do
      response = db_transaction do |tx|
        ride = tx.xquery('SELECT * FROM rides WHERE chair_id = ? ORDER BY updated_at DESC LIMIT 1', @current_chair.id).first
        unless ride
          halt json(data: nil, retry_after_ms: 30)
        end

        yet_sent_ride_status = tx.xquery('SELECT * FROM ride_statuses WHERE ride_id = ? AND chair_sent_at IS NULL ORDER BY created_at ASC LIMIT 1', ride.fetch(:id)).first
        status =
          if yet_sent_ride_status.nil?
            get_latest_ride_status(tx, ride.fetch(:id))
          else
            yet_sent_ride_status.fetch(:status)
          end

        user = tx.xquery('SELECT * FROM users WHERE id = ? FOR SHARE', ride.fetch(:user_id)).first

        unless yet_sent_ride_status.nil?
          tx.xquery('UPDATE ride_statuses SET chair_sent_at = CURRENT_TIMESTAMP(6) WHERE id = ?', yet_sent_ride_status.fetch(:id))
        end

        {
          data: {
            ride_id: ride.fetch(:id),
            user: {
              id: user.fetch(:id),
              name: "#{user.fetch(:firstname)} #{user.fetch(:lastname)}",
            },
            pickup_coordinate: {
              latitude: ride.fetch(:pickup_latitude),
              longitude: ride.fetch(:pickup_longitude),
            },
            destination_coordinate: {
              latitude: ride.fetch(:destination_latitude),
              longitude: ride.fetch(:destination_longitude),
            },
            status:,
          },
          retry_after_ms: 30,
        }
      end

      json(response)
    end

    PostChairRidesRideIDStatusRequest = Data.define(:status)

    # POST /api/chair/rides/:ride_id/status
    post '/rides/:ride_id/status' do
      ride_id = params[:ride_id]
      req = bind_json(PostChairRidesRideIDStatusRequest)

      db_transaction do |tx|
        ride = tx.xquery('SELECT * FROM rides WHERE id = ? FOR UPDATE', ride_id).first
        if ride.fetch(:chair_id) != @current_chair.id
          raise HttpError.new(400, 'not assigned to this ride')
        end

        case req.status
	# Acknowledge the ride
        when 'ENROUTE'
          tx.xquery('INSERT INTO ride_statuses (id, ride_id, status) VALUES (?, ?, ?)', ULID.generate, ride.fetch(:id), 'ENROUTE')
	# After Picking up user
        when 'CARRYING'
          status = get_latest_ride_status(tx, ride.fetch(:id))
          if status != 'PICKUP'
            raise HttpError.new(400, 'chair has not arrived yet')
          end
          tx.xquery('INSERT INTO ride_statuses (id, ride_id, status) VALUES (?, ?, ?)', ULID.generate, ride.fetch(:id), 'CARRYING')
        else
          raise HttpError.new(400, 'invalid status')
        end
      end

      status(204)
    end
  end
end
