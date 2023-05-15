CREATE TABLE info (
  key             TEXT PRIMARY KEY,
  value           TEXT NOT NULL,
  created_at      TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at      TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
INSERT INTO info (key, value)
VALUES ('Author', 'Kittipat Poonyakariyakorn');

CREATE TRIGGER info_updated_at_modtime BEFORE UPDATE ON info FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();