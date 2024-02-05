CREATE TABLE IF NOT EXISTS national_park (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT UNIQUE,
  link TEXT,
  year INTEGER,
  total_area_in_km INTEGER,
  total_area_in_miles INTEGER,
  water_percentages TEXT,
  region TEXT,
  description TEXT,
  coordinate_latitude REAL,
  coordinate_longitude REAL,
  map_url TEXT,
  location TEXT,
  established INTEGER,
  visitors TEXT,
  management TEXT
);

CREATE TABLE IF NOT EXISTS image (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  link TEXT,
  title TEXT,
  date TEXT,
  source TEXT,
  author TEXT,
  src TEXT
);

CREATE TABLE IF NOT EXISTS license (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  type TEXT,
  name TEXT,
  link TEXT
);

CREATE TABLE IF NOT EXISTS intl_status (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT,
  link TEXT
);

-- This table establishes a many-to-many relationship with license
CREATE TABLE IF NOT EXISTS national_park_image (
  national_park_id INTEGER,
  image_id INTEGER,
  FOREIGN KEY (image_id) REFERENCES image(id),
  FOREIGN KEY (national_park_id) REFERENCES national_park(id)
);

-- This table establishes a many-to-many relationship with license
CREATE TABLE IF NOT EXISTS image_license (
  image_id INTEGER,
  license_id INTEGER,
  FOREIGN KEY (image_id) REFERENCES image(id),
  FOREIGN KEY (license_id) REFERENCES license(id)
);

-- This table establishes a many-to-many relationship with national park
CREATE TABLE IF NOT EXISTS national_park_intl_status (
  national_park_id INTEGER,
  intl_status_id INTEGER,
  FOREIGN KEY (national_park_id) REFERENCES national_park(id),
  FOREIGN KEY (intl_status_id) REFERENCES intl_status(id)
)