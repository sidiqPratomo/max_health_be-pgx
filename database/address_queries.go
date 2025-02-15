package database

const (
	FindAllProvinces = `
		SELECT province_id, province_code, province_name
		FROM provinces
		WHERE deleted_at IS NULL
	`

	FindOneProvinceByName = `
		SELECT province_id
		FROM provinces
		WHERE province_name ILIKE $1
		AND deleted_at IS NULL
		LIMIT 1
	`

	FindAllCitiesByProvinceCode = `
		SELECT city_id, city_code, province_code, city_name
		FROM cities
		WHERE province_code = $1
		AND deleted_at IS NULL
	`

	FindOneCityByName = `
		SELECT city_id
		FROM cities
		WHERE city_name ILIKE $1
		AND deleted_at IS NULL
		LIMIT 1
	`

	FindAllDistrictsByCityCode = `
		SELECT district_id, district_code, city_code, district_name
		FROM districts
		WHERE city_code = $1
		AND deleted_at IS NULL
	`

	FindOneDistrictByName = `
		SELECT district_id
		FROM districts
		WHERE district_name ILIKE $1
		AND deleted_at IS NULL
		LIMIT 1
	`

	FindAllSubdistrictsByDistrictCode = `
		SELECT subdistrict_id, subdistrict_code, district_code, subdistrict_name
		FROM subdistricts
		WHERE district_code = $1
		AND deleted_at IS NULL
	`

	FindOneSubdistrictByName = `
		SELECT subdistrict_id
		FROM subdistricts
		WHERE subdistrict_name ILIKE $1
		AND deleted_at IS NULL
		LIMIT 1
	`
)
