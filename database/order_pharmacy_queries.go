package database

const (
	FindAllOrderPharmaciesByOrderUserId = `
		SELECT op.order_pharmacy_id, op.order_status_id, op.subtotal_amount, op.delivery_fee, p.pharmacy_name, a.profile_picture, first_order_item.drug_name, first_order_item.drug_price, first_order_item.drug_unit, first_order_item.quantity, first_order_item.image, first_order_item.total_order_items, count(*) OVER() AS total_item_count
		FROM order_pharmacies op
		JOIN pharmacy_couriers pc ON pc.pharmacy_courier_id = op.pharmacy_courier_id
		JOIN pharmacies p ON p.pharmacy_id = pc.pharmacy_id
		JOIN pharmacy_managers pm ON pm.pharmacy_manager_id = p.pharmacy_manager_id
		JOIN accounts a ON a.account_id = pm.account_id
		JOIN orders o ON o.order_id = op.order_id
		JOIN (
			SELECT oi.order_pharmacy_id, oi.drug_name, oi.drug_price, oi.quantity, oi.drug_unit, d.image, oi.total_order_items
			FROM (
				 SELECT *, ROW_NUMBER() OVER (PARTITION BY order_pharmacy_id) AS ROW_NUMBER, COUNT(*) OVER(PARTITION BY order_pharmacy_id) AS total_order_items
				 FROM order_items
			) oi
			JOIN drugs d ON d.drug_id = oi.drug_id
			WHERE oi.row_number = 1
		) AS first_order_item ON first_order_item.order_pharmacy_id = op.order_pharmacy_id
		WHERE o.user_id = $1
		AND op.deleted_at IS NULL
	`

	FindAllOrderPharmaciesByPharmacyManagerId = `
		SELECT DISTINCT op.order_pharmacy_id, op.updated_at, count(*) OVER() AS total_item_count
		FROM order_pharmacies op
		JOIN pharmacy_couriers pc ON pc.pharmacy_courier_id = op.pharmacy_courier_id
		JOIN pharmacies p ON p.pharmacy_id = pc.pharmacy_id
		JOIN pharmacy_managers pm ON pm.pharmacy_manager_id = p.pharmacy_manager_id
		WHERE pm.pharmacy_manager_id = $1
		AND op.deleted_at IS NULL
	`

	FindAllOrderPharmacies = `
		SELECT DISTINCT op.order_pharmacy_id, op.updated_at, count(*) OVER() AS total_item_count
		FROM order_pharmacies op
		JOIN pharmacy_couriers pc ON pc.pharmacy_courier_id = op.pharmacy_courier_id
		JOIN pharmacies p ON p.pharmacy_id = pc.pharmacy_id
		JOIN pharmacy_managers pm ON pm.pharmacy_manager_id = p.pharmacy_manager_id
		WHERE op.deleted_at IS NULL
	`

	FindAllOrderPharmaciesWithDetails = `
		SELECT op.order_pharmacy_id, op.order_status_id, op.subtotal_amount, op.delivery_fee, op.updated_at, p.pharmacy_name, p.pharmacist_phone_number, c.courier_name, a.email, oi.drug_name, oi.drug_price, oi.drug_unit, d.image, oi.quantity
		FROM order_pharmacies op
		JOIN pharmacy_couriers pc ON pc.pharmacy_courier_id = op.pharmacy_courier_id
		JOIN couriers c ON c.courier_id = pc.courier_id
		JOIN pharmacies p ON p.pharmacy_id = pc.pharmacy_id
		JOIN pharmacy_managers pm ON pm.pharmacy_manager_id = p.pharmacy_manager_id
		JOIN orders o ON o.order_id = op.order_id
		JOIN order_items oi ON oi.order_pharmacy_id = op.order_pharmacy_id
		JOIN drugs d ON oi.drug_id = d.drug_id
		JOIN accounts a ON a.account_id = pm.account_id
		WHERE op.deleted_at IS NULL
	`

	FindOneOrderPharmacyById = `
		SELECT op.order_pharmacy_id, op.order_status_id, op.subtotal_amount, op.delivery_fee, op.created_at, op.updated_at, p.pharmacy_name, c.courier_name, a.profile_picture, o.address
		FROM order_pharmacies op
		JOIN pharmacy_couriers pc ON pc.pharmacy_courier_id = op.pharmacy_courier_id
		JOIN pharmacies p ON p.pharmacy_id = pc.pharmacy_id
		JOIN couriers c ON c.courier_id = pc.courier_id
		JOIN pharmacy_managers pm ON pm.pharmacy_manager_id = p.pharmacy_manager_id
		JOIN accounts a ON a.account_id = pm.account_id
		JOIN orders o ON o.order_id = op.order_id
		WHERE op.order_pharmacy_id = $1
		AND op.deleted_at IS NULL
	`

	FindAllOrderPharmaciesByOrderId = `
		SELECT op.order_pharmacy_id, o.user_id, op.order_status_id
		FROM orders o
		JOIN order_pharmacies op ON op.order_id = o.order_id
		WHERE o.order_id = $1 
		AND o.deleted_at IS NULL
	`

	FindAllOngoingOrderPharmacyIdsByPharmacyId = `
		SELECT op.order_pharmacy_id
		FROM order_pharmacies op
		JOIN pharmacy_couriers pc ON pc.pharmacy_courier_id = op.pharmacy_courier_id
		WHERE pc.pharmacy_id = $1
		AND op.order_status_id IN (1,2,3,4)
		AND op.deleted_at IS NULL
	`

	UpdateStatusBulkOrderPharmaciesByOrderId = `
		UPDATE order_pharmacies 
		SET order_status_id = $1,
		updated_at = NOW()
		WHERE order_id = $2
	`

	FindOrderPharmacyByOrderPharmacyId = `
		SELECT op.order_pharmacy_id, o.user_id, op.order_id, op.order_status_id, op.pharmacy_courier_id, op.subtotal_amount, 
			op.delivery_fee
		FROM order_pharmacies op
		JOIN orders o ON op.order_id = o.order_id
		WHERE op.order_pharmacy_id = $1 
		AND op.deleted_at IS NULL
	`

	FindOrderPharmacyCountGroupedByOrderStatusIdByPharmacyManagerId = `
		SELECT op.order_status_id, COUNT (op.order_pharmacy_id)
		FROM order_pharmacies op
		JOIN pharmacy_couriers pc ON pc.pharmacy_courier_id = op.pharmacy_courier_id
		JOIN pharmacies p ON p.pharmacy_id = pc.pharmacy_id
		JOIN pharmacy_managers pm ON pm.pharmacy_manager_id = p.pharmacy_manager_id
		WHERE pm.pharmacy_manager_id = $1
		AND op.deleted_at IS NULL
		GROUP BY op.order_status_id
	`

	UpdateOneStatusById = `
		UPDATE order_pharmacies 
		SET order_status_id = $1,
		updated_at = NOW()
		WHERE order_pharmacy_id = $2
	`
)
