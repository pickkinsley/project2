-- PackSmart database schema
-- MySQL 8.0
-- Run: mysql -u root -p < schema.sql

CREATE DATABASE IF NOT EXISTS packsmart
  CHARACTER SET utf8mb4
  COLLATE utf8mb4_unicode_ci;

USE packsmart;

-- ─────────────────────────────────────────
-- Table 1: trips
-- ─────────────────────────────────────────
CREATE TABLE trips (
  id             CHAR(36)      PRIMARY KEY,
  destination    VARCHAR(255)  NOT NULL,
  dest_lat       DECIMAL(9,6)  NOT NULL,
  dest_lon       DECIMAL(9,6)  NOT NULL,
  departure_date DATE          NOT NULL,
  return_date    DATE          NOT NULL,
  trip_type      VARCHAR(50)   NOT NULL,
  companions     VARCHAR(50)   NOT NULL,
  activities     JSON          NOT NULL,
  created_at     TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- ─────────────────────────────────────────
-- Table 2: packing_items
-- ─────────────────────────────────────────
CREATE TABLE packing_items (
  id           INT          PRIMARY KEY AUTO_INCREMENT,
  trip_id      CHAR(36)     NOT NULL,
  name         VARCHAR(255) NOT NULL,
  category     VARCHAR(100) NOT NULL,
  is_essential BOOLEAN      NOT NULL DEFAULT FALSE,
  reason       VARCHAR(500),
  is_checked   BOOLEAN      NOT NULL DEFAULT FALSE,
  sort_order   INT          NOT NULL DEFAULT 0,
  created_at   TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,

  FOREIGN KEY (trip_id) REFERENCES trips(id) ON DELETE CASCADE,
  INDEX idx_trip_id (trip_id)
);

-- ─────────────────────────────────────────
-- Table 3: weather_snapshots
-- ─────────────────────────────────────────
CREATE TABLE weather_snapshots (
  id             INT          PRIMARY KEY AUTO_INCREMENT,
  trip_id        CHAR(36)     NOT NULL UNIQUE,
  temp_min_f     INT          NOT NULL,
  temp_max_f     INT          NOT NULL,
  rain_days      INT          NOT NULL DEFAULT 0,
  snow_days      INT          NOT NULL DEFAULT 0,
  is_forecast    BOOLEAN      NOT NULL,
  daily_forecast JSON,
  fetched_at     TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,

  FOREIGN KEY (trip_id) REFERENCES trips(id) ON DELETE CASCADE
);
