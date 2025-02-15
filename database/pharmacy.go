package database

const (
	GetAllCourierOptionsByPharmacyId = `
		SELECT 
			pc.pharmacy_courier_id,
			c.courier_name,
			CEIL(ST_DistanceSphere(ua.geom, p.geom) / 1000) * c.price,
			c2.raja_ongkir_id,
			c3.raja_ongkir_id 
		FROM pharmacies p
		JOIN user_addresses ua ON ua.user_address_id = $1
		JOIN pharmacy_couriers pc ON pc.pharmacy_id = p.pharmacy_id
		JOIN couriers c ON c.courier_id = pc.courier_id
		JOIN cities c2 ON c2.city_id = ua.city_id
		JOIN cities c3 ON c3.city_name ILIKE CONCAT('%', p.city)
		WHERE p.pharmacy_id = $2
	`

	GetOnePharmacyByPharmacyId = `
		select 
		pharmacy_id, 
		pharmacy_manager_id, 
		pharmacy_name, 
		pharmacist_name, 
		pharmacist_license_number, 
		pharmacist_phone_number, 
		city, 
		address,
		latitude, 
		longitude 
		FROM pharmacies
		WHERE pharmacy_id = $1 AND deleted_at IS NULL
	`
)
