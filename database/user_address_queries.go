package database

const (
	SetAllIsMainFalseQuery = `
		UPDATE user_addresses
		SET is_main = FALSE, updated_at = NOW()
		WHERE user_id = $1
		AND deleted_at IS NULL
	`

	PostOneUserAddressQuery = `
		INSERT INTO user_addresses (user_id, province_id, city_id, district_id, subdistrict_id, latitude, longitude, label, address, is_active, is_main, geom)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, TRUE, $10, ST_SetSRID(ST_MakePoint($11, $12), 4326))
	`

	GetOneUserAddressByAddressIdQuery = `
		SELECT user_id, province_id, city_id, district_id, subdistrict_id, latitude, longitude, label, address, is_active, is_main
		FROM user_addresses
		WHERE user_address_id = $1
		AND deleted_at IS NULL
	`

	UpdateUserAddressQuery = `
		UPDATE user_addresses
		SET
		province_id = $1,
		city_id = $2,
		district_id = $3,
		subdistrict_id = $4,
		latitude = $5,
		longitude = $6,
		label = $7,
		address = $8,
		is_active = $9,
		is_main = $10,
		geom = ST_SetSRID(ST_MakePoint($11, $12), 4326),
		updated_at = NOW()
		WHERE user_address_id = $13
		AND deleted_at IS NULL
	`

	FindAllUserAddressByUserIdQuery = `
		SELECT ua.user_address_id, p.province_id, p.province_code, p.province_name, c.city_id, c.city_code, c.city_name, d.district_id, d.district_code, d.district_name, s.subdistrict_id, s.subdistrict_code, s.subdistrict_name, ua.latitude, ua.longitude, ua.label, ua.address, ua.is_main
		FROM user_addresses ua
		LEFT JOIN provinces p on p.province_id = ua.province_id
		LEFT JOIN cities c on c.city_id = ua.city_id
		LEFT JOIN districts d on d.district_id = ua.district_id
		LEFT JOIN subdistricts s on s.subdistrict_id = ua.subdistrict_id
		WHERE ua.user_id = $1
		AND ua.deleted_at IS NULL
		ORDER BY (ua.is_main IS TRUE) DESC
	`

	GetOneUserAddressById = `
		SELECT user_id
		FROM user_addresses
		WHERE user_address_id = $1 
		AND deleted_at IS NULL
	`

	DeleteOneUserAddressQuery = `
		UPDATE user_addresses
		SET deleted_at = NOW(), updated_at = NOW()
		WHERE user_address_id = $1
	`
)
