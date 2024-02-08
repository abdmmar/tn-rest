-- name: GetNationalParks :many
SELECT
  np.id AS id,
  np.name AS name,
  np.description AS description,
  np.region AS region,
  JSON_GROUP_ARRAY(
    JSON_OBJECT(
      'id',
      img.id,
      'link',
      img.link,
      'title',
      img.title,
      'date',
      img.date,
      'source',
      img.source,
      'author',
      img.author,
      'src',
      img.src
    )
  ) AS images,
  JSON_GROUP_ARRAY(
    JSON_OBJECT(
      'id',
      intl.id,
      'name',
      intl.name,
      'link',
      intl.link
    )
  ) as intl_statuses,
  JSON_OBJECT(
    'km',
    np.total_area_in_km,
    'miles',
    np.total_area_in_miles
  ) as total_area,
  JSON_OBJECT(
    'lat',
    np.coordinate_latitude,
    'long',
    np.coordinate_longitude
  ) as coordinate,
  np.water_percentages AS water_percentages,
  np.map_url AS map_url,
  np.location AS location,
  np.established_year AS established_year,
  np.visitors AS visitors,
  np.management AS management
FROM
  national_park np
  LEFT JOIN national_park_image npi ON np.id = npi.national_park_id
  LEFT JOIN image img ON npi.image_id = img.id
  LEFT JOIN image_license il ON img.id = il.image_id
  LEFT JOIN license lic ON il.license_id = lic.id
  LEFT JOIN national_park_intl_status nps ON np.id = nps.national_park_id
  LEFT JOIN intl_status intl ON nps.intl_status_id = intl.id
GROUP BY
  np.id;

-- name: GetNationalPark :one
SELECT
  np.id AS id,
  np.name AS name,
  np.description AS description,
  np.region AS region,
  JSON_GROUP_ARRAY(
    JSON_OBJECT(
      'id',
      img.id,
      'link',
      img.link,
      'title',
      img.title,
      'date',
      img.date,
      'source',
      img.source,
      'author',
      img.author,
      'src',
      img.src
    )
  ) AS images,
  JSON_GROUP_ARRAY(
    JSON_OBJECT(
      'id',
      intl.id,
      'name',
      intl.name,
      'link',
      intl.link
    )
  ) as intl_statuses,
  JSON_OBJECT(
    'km',
    np.total_area_in_km,
    'miles',
    np.total_area_in_miles
  ) as total_area,
  JSON_OBJECT(
    'lat',
    np.coordinate_latitude,
    'long',
    np.coordinate_longitude
  ) as coordinate,
  np.water_percentages AS water_percentages,
  np.map_url AS map_url,
  np.location AS location,
  np.established_year AS established_year,
  np.visitors AS visitors,
  np.management AS management
FROM
  national_park np
  LEFT JOIN national_park_image npi ON np.id = npi.national_park_id
  LEFT JOIN image img ON npi.image_id = img.id
  LEFT JOIN image_license il ON img.id = il.image_id
  LEFT JOIN license lic ON il.license_id = lic.id
  LEFT JOIN national_park_intl_status nps ON np.id = nps.national_park_id
  LEFT JOIN intl_status intl ON nps.intl_status_id = intl.id
WHERE
  np.name = ?
GROUP BY
  np.id;

-- name: CreateNationalPark :one
INSERT INTO
  national_park (
    name,
    description,
    region,
    link,
    total_area_in_km,
    total_area_in_miles,
    water_percentages,
    coordinate_latitude,
    coordinate_longitude,
    map_url,
    location,
    established_year,
    visitors,
    management
  )
VALUES
  (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) RETURNING id;

-- name: GetLicense :one
SELECT * FROM license WHERE id = ? LIMIT 1; 

-- name: GetLicenseByName :one
SELECT * FROM license WHERE name = ? LIMIT 1; 

-- name: GetIntlStatus :one
SELECT * FROM intl_status WHERE id = ? LIMIT 1;

-- name: GetIntlStatusByName :one
SELECT * FROM intl_status WHERE name = ? LIMIT 1;

-- name: CreateImage :one
INSERT INTO
  image (
    link,
    title,
    date,
    source,
    author,
    src
  )
VALUES
  (?, ?, ?, ?, ?, ?) RETURNING id;

-- name: CreateIntlStatus :one
INSERT INTO intl_status (
  name, link
) VALUES (?,?) RETURNING id;

-- name: CreateLicense :one
INSERT INTO license (
  type, name, link
) VALUES (?,?,?) RETURNING id;

-- name: CreateImageLicense :exec
INSERT INTO image_license (
  image_id, license_id
) VALUES (?, ?);

-- name: CreateNationalParkImage :exec
INSERT INTO national_park_image (
  national_park_id, image_id
) VALUES (?, ?);

-- name: CreateNationalParkIntlStatus :exec
INSERT INTO national_park_intl_status (
  national_park_id, intl_status_id
) VALUES (?, ?);

