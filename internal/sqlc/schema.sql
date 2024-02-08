CREATE TABLE IF NOT EXISTS national_park (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT UNIQUE NOT NULL,
  link TEXT,
  total_area_in_km INTEGER,
  total_area_in_miles INTEGER,
  water_percentages TEXT,
  region TEXT NOT NULL,
  description TEXT NOT NULL,
  coordinate_latitude REAL NOT NULL,
  coordinate_longitude REAL NOT NULL,
  map_url TEXT,
  location TEXT NOT NULL,
  established_year INTEGER NOT NULL,
  visitors TEXT,
  management TEXT
);

CREATE TABLE IF NOT EXISTS image (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  link TEXT NOT NULL,
  title TEXT NOT NULL,
  date TEXT,
  source TEXT,
  author TEXT,
  src TEXT
);

CREATE TABLE IF NOT EXISTS license (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  type TEXT NOT NULL,
  name TEXT NOT NULL,
  link TEXT
);

CREATE TABLE IF NOT EXISTS intl_status (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL,
  link TEXT
);

-- This table establishes a many-to-many relationship with license
CREATE TABLE IF NOT EXISTS national_park_image (
  national_park_id INTEGER NOT NULL,
  image_id INTEGER NOT NULL,
  FOREIGN KEY (image_id) REFERENCES image(id),
  FOREIGN KEY (national_park_id) REFERENCES national_park(id)
);

-- This table establishes a many-to-many relationship with license
CREATE TABLE IF NOT EXISTS image_license (
  image_id INTEGER NOT NULL,
  license_id INTEGER NOT NULL,
  FOREIGN KEY (image_id) REFERENCES image(id),
  FOREIGN KEY (license_id) REFERENCES license(id)
);

-- This table establishes a many-to-many relationship with national park
CREATE TABLE IF NOT EXISTS national_park_intl_status (
  national_park_id INTEGER NOT NULL,
  intl_status_id INTEGER NOT NULL,
  FOREIGN KEY (national_park_id) REFERENCES national_park(id),
  FOREIGN KEY (intl_status_id) REFERENCES intl_status(id)
)