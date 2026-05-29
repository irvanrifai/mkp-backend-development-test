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

INSERT INTO
  "studios" (branch_name, studio_number, studio_type, total_seats)
VALUES
  ('MKP Semarang Studio', 1, 'REGULAR', 100),
  ('MKP Ungaran Studio', 2, 'REGULAR', 100);
