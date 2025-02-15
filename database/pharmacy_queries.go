package database

const (
	GetAllPharmacyByPharmacyManagerId = `
		select p.pharmacy_id , p.pharmacy_manager_id, p.pharmacy_name, p.pharmacist_name, p.pharmacist_license_number,  
		p.pharmacist_phone_number, p.city, p.address, p.longitude, p.latitude
		from pharmacies p 
		where p.pharmacy_manager_id =$1 and pharmacy_name ILIKE '%' || $4 || '%' and deleted_at is null
		LIMIT $2
		OFFSET $3
	`

	FindOnePharmacyById = `
		SELECT pharmacy_id, pharmacy_name, pharmacy_manager_id
		FROM pharmacies
		WHERE pharmacy_id = $1
		AND deleted_at IS NULL
	`

	CreateOnePharmacy = `
		INSERT 
		INTO pharmacies (pharmacy_manager_id, pharmacy_name, pharmacist_name, pharmacist_license_number, pharmacist_phone_number, address, city, latitude, longitude, geom)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, ST_SetSRID(ST_MakePoint($10, $11), 4326))
		RETURNING pharmacy_id
	`

	UpdateOnePharmacy = `
		UPDATE pharmacies
		SET
		pharmacy_name = $1,
		pharmacist_name = $2,
		pharmacist_license_number = $3,
		pharmacist_phone_number = $4,
		address = $5,
		city = $6,
		latitude = $7,
		longitude = $8,
		geom = ST_SetSRID(ST_MakePoint($9, $10), 4326),
		updated_at = NOW()
		WHERE pharmacy_id = $11
	`

	DeleteOnePharmacyById = `
		UPDATE pharmacies
		SET deleted_at = NOW(), updated_at = NOW()
		WHERE pharmacy_id = $1
	`

	CreatePharmacyOperational = `
		INSERT INTO pharmacy_operationals (pharmacy_id, operational_day)
		VALUES
	`

	UpdateOnePharmacyOperational = `
		UPDATE pharmacy_operationals 
		SET operational_day = $1,
		open_hour = $2,
		close_hour = $3,
		is_open = $4,
		updated_at = NOW()
		WHERE pharmacy_operational_id = $5
	`

	DeleteBulkPharmacyOperationalByPharmacyId = `
		UPDATE pharmacy_operationals
		SET deleted_at = NOW(), updated_at = NOW()
		WHERE pharmacy_id = $1
	`

	CreatePharmacyCourier = `
		INSERT INTO pharmacy_couriers (pharmacy_id, courier_id, is_active)
		VALUES
	`

	UpdateOnePharmacyCourier = `
		UPDATE pharmacy_couriers 
		SET is_active = $1,
		updated_at = NOW()
		WHERE pharmacy_courier_id = $2
	`

	DeleteBulkPharmacyCourierByPharmacyId = `
		UPDATE pharmacy_couriers
		SET deleted_at = NOW(), updated_at = NOW()
		WHERE pharmacy_id = $1
	`
)
