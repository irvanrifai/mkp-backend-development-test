-- Table Users
CREATE TABLE IF NOT EXISTS users (
  id BIGSERIAL PRIMARY KEY,
  name VARCHAR(150),
  username VARCHAR(50) UNIQUE NOT NULL,
  email VARCHAR(150) UNIQUE NOT NULL,
  password VARCHAR(255) NOT NULL,
  phone VARCHAR(20),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP
);

-- Table Movies
CREATE TABLE IF NOT EXISTS movies (
  id BIGSERIAL PRIMARY KEY,
  title VARCHAR(255) NOT NULL,
  description TEXT,
  duration_minutes INT NOT NULL,
  genre VARCHAR(100),
  poster_url VARCHAR(255),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP
);

-- Table Studios
CREATE TABLE IF NOT EXISTS studios (
  id BIGSERIAL PRIMARY KEY,
  branch_name VARCHAR(150) NOT NULL,
  studio_number INT,
  studio_type VARCHAR(50) DEFAULT 'REGULAR',
  total_seats INT NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP
);

-- Table Seats
CREATE TABLE IF NOT EXISTS seats (
  id BIGSERIAL PRIMARY KEY,
  studio_id BIGINT REFERENCES studios(id) ON DELETE CASCADE,
  seat_number VARCHAR(10) NOT NULL,
  status VARCHAR(50) DEFAULT 'ACTIVE', -- ACTIVE, BROKEN/INACTIVE
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Table Schedules
CREATE TABLE IF NOT EXISTS schedules (
  id BIGSERIAL PRIMARY KEY,
  movie_id BIGINT REFERENCES movies(id) ON DELETE RESTRICT,
  studio_id BIGINT REFERENCES studios(id) ON DELETE RESTRICT,
  show_time TIMESTAMP NOT NULL,
  price_per_ticket NUMERIC(10, 2) NOT NULL,
  status VARCHAR(50) DEFAULT 'ACTIVE',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP
);

-- Table Transactions
CREATE TABLE IF NOT EXISTS transactions (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT REFERENCES users(id) ON DELETE SET NULL,
  transactionable_type VARCHAR(50) NOT NULL, -- e.g., 'schedules'
  transactionable_id BIGINT NOT NULL,
  invoice_number VARCHAR(100) UNIQUE NOT NULL,
  subtotal NUMERIC(10, 2) NOT NULL,
  admin_fee NUMERIC(10, 2) DEFAULT 0.00,
  discount NUMERIC(10, 2) DEFAULT 0.00,
  total_price NUMERIC(10, 2) NOT NULL,
  payment_status VARCHAR(50) DEFAULT 'PENDING',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP
);

-- Table Ticket Seats
CREATE TABLE IF NOT EXISTS ticket_seats (
  id BIGSERIAL PRIMARY KEY,
  transaction_id BIGINT REFERENCES transactions(id) ON DELETE CASCADE,
  seat_id BIGINT REFERENCES seats(id) ON DELETE RESTRICT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  UNIQUE(transaction_id, seat_id)
);

-- Indexes
CREATE INDEX idx_schedules_show_time ON schedules(show_time);
CREATE INDEX idx_transactions_user_status ON transactions(user_id, payment_status);
