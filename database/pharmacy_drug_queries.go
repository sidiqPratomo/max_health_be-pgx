package database

const (
	GetPharmacyDrugByDrugIdQuery = `
		SELECT
			pd.pharmacy_drug_id,
			pd.pharmacy_id,
			p.pharmacy_name,
			p.pharmacist_name,
			p.pharmacist_license_number,
			p.pharmacist_phone_number,
			ST_DistanceSphere((ST_SetSRID(ST_MakePoint($2, $3), 4326)), p.geom),
			pd.drug_id,
			pd.price,
			pd.stock
		FROM pharmacy_drugs pd
		JOIN pharmacies p ON p.pharmacy_id = pd.pharmacy_id
		JOIN drugs d ON d.drug_id = pd.drug_id
		WHERE d.drug_id = $1 
			AND pd.deleted_at IS NULL
			AND ST_DistanceSphere((ST_SetSRID(ST_MakePoint($2, $3), 4326)), p.geom) <= 25000
		ORDER BY ST_DistanceSphere((ST_SetSRID(ST_MakePoint($2, $3), 4326)), p.geom) ASC
		LIMIT $4
		OFFSET $5
	`

	GetPharmacyInRangeQuery = `
		WITH in_range_pharmacy AS (
			SELECT pharmacy_id, CAST((ST_DistanceSphere((ST_SetSRID(ST_MakePoint($1, $2), 4326)), pharmacies.geom)) AS NUMERIC) AS distance
			FROM pharmacies
			WHERE deleted_at IS NULL AND (ST_DistanceSphere((ST_SetSRID(ST_MakePoint($1, $2), 4326)), pharmacies.geom)) <= 25000
	`

	GetDrugListQuery = `
		drug_list AS(
			SELECT pd.pharmacy_drug_id, pd.drug_id, d.drug_name, pd.price, d.image, ip.distance, d.is_prescription_required
			FROM pharmacy_drugs pd 
			JOIN in_range_pharmacy ip 
			ON pd.pharmacy_id = ip.pharmacy_id
			JOIN drugs d
			ON pd.drug_id = d.drug_id
			WHERE pd.deleted_at IS NULL AND d.deleted_at IS NULL AND pd.stock > 0 AND d.is_active AND d.drug_name ILIKE 
	`

	GetPriceRangeQuery = `
		prices AS(
			SELECT dl.drug_id, MIN(dl.price) AS min_price, MAX(dl.price) AS max_price
			FROM drug_list dl
			GROUP BY dl.drug_id),
		closest_drug AS(
			SELECT dl.pharmacy_drug_id, dl.drug_id, dl.drug_name, dl.price, dl.image, dl.is_prescription_required 
			FROM drug_list dl
			JOIN
				(SELECT drug_id, MIN(distance) AS min_distance
				FROM drug_list GROUP BY drug_id
				) dl2
			ON dl.distance = dl2.min_distance AND dl.drug_id = dl2.drug_id)
	`

	GetProductListingQuery = `
		SELECT cd.pharmacy_drug_id, cd.drug_id, cd.drug_name, p.min_price, p.max_price, cd.image, cd.is_prescription_required 
		FROM closest_drug cd
		JOIN prices p 
		ON cd.drug_id = p.drug_id
	`

	GetProductCountQuery = `
		SELECT COUNT(drug_id) FROM closest_drug;
	`

	GetPharmacyDrugById = `
		select pharmacy_drug_id, pharmacy_id, drug_id, price, stock 
		from pharmacy_drugs where pharmacy_drug_id = $1 and deleted_at is null
	`

	GetPharmacyDrugsByCartForUpdate = `
		SELECT * FROM pharmacy_drugs pd
		JOIN cart_items ci 
		ON ci.pharmacy_drug_id = pd.pharmacy_drug_id
	`

	UpdatePharmacyDrugsByCartId = `
		UPDATE pharmacy_drugs pd 
		SET stock = pd.stock - ci.quantity,
		updated_at = NOW()
		FROM cart_items ci
		WHERE ci.pharmacy_drug_id = pd.pharmacy_drug_id
	`

	UpdatePharmacyDrugsFromStockMutation1 = `
		UPDATE pharmacy_drugs AS pd SET
		stock = c.stock
		FROM (VALUES 
	`

	UpdatePharmacyDrugsFromStockMutation2 = `
		) AS c(pharmacy_drug_id, stock) 
		WHERE c.pharmacy_drug_id = pd.pharmacy_drug_id
	`

	GetPharmacyDrugByPharmacyId = `
		select pharmacy_drug_id, pharmacy_id, drug_id, price, stock 
		from pharmacy_drugs where pharmacy_id = $1 and deleted_at IS NULL
	`

	GetNearestAvailablePharmacyDrugByDrugIdQuery = `
		SELECT
			p.pharmacy_id,
			p.pharmacy_name,
			p.address,
			ST_DistanceSphere(ua.geom, p.geom),
			pd.pharmacy_drug_id,
			d.drug_id,
			d.drug_name,
			d.manufacture,
			d.image,
			d.weight,
			d.selling_unit,
			d.unit_in_pack,
			pd.price
		FROM pharmacy_drugs pd
		JOIN pharmacies p ON p.pharmacy_id = pd.pharmacy_id
		JOIN user_addresses ua ON ua.user_address_id = $2
		JOIN drugs d ON d.drug_id = pd.drug_id
		WHERE d.drug_id = $1 
			AND pd.deleted_at IS NULL 
			AND d.deleted_at IS NULL
			AND ST_DistanceSphere(ua.geom, p.geom) <= 25000
			AND pd.stock > 0
		ORDER BY ST_DistanceSphere(ua.geom, p.geom) ASC
		LIMIT 1
	`

	GetPharmacyDrugsByOrderPharmacyId = `
		WITH order_drugs AS (
			SELECT oi.quantity, pd.pharmacy_drug_id, pd.stock 
			FROM order_pharmacies op
			JOIN order_items oi ON op.order_pharmacy_id = oi.order_pharmacy_id
			JOIN pharmacy_drugs pd ON pd.pharmacy_drug_id = oi.pharmacy_drug_id
			WHERE op.order_pharmacy_id = $1 AND op.deleted_at IS NULL)
	`

	UpdatePharmacyDrugsByOrderPharmacyId = `
		UPDATE pharmacy_drugs pd 
		SET stock = pd.stock + od.quantity,
		updated_at = NOW()
		FROM order_drugs od
		WHERE od.pharmacy_drug_id = pd.pharmacy_drug_id
		RETURNING pd.pharmacy_drug_id, pd.stock, od.quantity
	`

	GetPharmacyDrugsByOrderId = `
		WITH order_drugs AS (
			SELECT oi.quantity, pd.pharmacy_drug_id, pd.stock 
			FROM order_pharmacies op
			JOIN orders o ON o.order_id = op.order_id
			JOIN order_items oi ON op.order_pharmacy_id = oi.order_pharmacy_id
			JOIN pharmacy_drugs pd ON pd.pharmacy_drug_id = oi.pharmacy_drug_id
			WHERE op.order_id = $1 AND op.deleted_at IS NULL)
			`

	UpdatePharmacyDrugStockPrice = `
		update pharmacy_drugs set stock=$2, price= $3, updated_at=now()
		where pharmacy_drug_id = $1;
	`

	DeletePharmacyDrug = `
		update pharmacy_drugs set updated_at=now(), deleted_at= now()
		where pharmacy_drug_id = $1;
	`

	AddPharmacyDrug = `
		insert into pharmacy_drugs (pharmacy_id, drug_id, stock, price)
		VALUES ($1, $2, $3, $4)
	`

	GetPossibleStockMutation = `
		WITH detailed_pharmacy_drugs AS (
			SELECT pd.pharmacy_drug_id, pd.drug_id, pd.pharmacy_id, p.pharmacy_manager_id
			FROM pharmacy_drugs pd
			JOIN pharmacies p 
			ON p.pharmacy_id = pd.pharmacy_id 
			WHERE pd.pharmacy_drug_id = $1
		)
		SELECT pd.pharmacy_drug_id, pd.pharmacy_id, p.pharmacy_name, p.address, pd.stock
		FROM pharmacy_drugs pd
		JOIN pharmacies p
		ON p.pharmacy_id = pd.pharmacy_id
		JOIN detailed_pharmacy_drugs dp
		ON dp.drug_id = pd.drug_id
		WHERE pd.pharmacy_drug_id != dp.pharmacy_drug_id AND pd.deleted_at IS NULL 
		AND p.pharmacy_manager_id = dp.pharmacy_manager_id AND p.deleted_at IS NULL
		AND pd.stock > 0
	`

	GetPharmacyDrugByIdForUpdate = `
		select pharmacy_drug_id, pharmacy_id, drug_id, price, stock 
		from pharmacy_drugs where pharmacy_drug_id = $1 and deleted_at is null for update
	`
)
