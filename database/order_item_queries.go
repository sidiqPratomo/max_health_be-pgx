package database

const (
	FindAllOrderItemByOrderPharmacyId = `
		SELECT oi.order_item_id, oi.drug_name, oi.drug_price, oi.drug_unit, oi.quantity, d.image
		FROM order_items oi
		JOIN order_pharmacies op ON op.order_pharmacy_id = oi.order_pharmacy_id
		JOIN drugs d ON oi.drug_id = d.drug_id
		WHERE oi.order_pharmacy_id = $1
		AND oi.deleted_at IS NULL;
	`

	FindPharmacyDrugCategorySalesVolumeRevenue = `
		SELECT dc.drug_category_id, dc.drug_category_name, COUNT (dc.drug_category_id) sales_volume, SUM (oi.quantity * oi.drug_price) revenue
		FROM order_items oi
		JOIN drugs d ON d.drug_id = oi.drug_id
		JOIN drug_categories dc ON dc.drug_category_id = d.drug_category_id
		JOIN order_pharmacies op ON op.order_pharmacy_id = oi.order_pharmacy_id
		JOIN pharmacy_couriers pc ON pc.pharmacy_courier_id = op.pharmacy_courier_id
		JOIN pharmacies p ON p.pharmacy_id = pc.pharmacy_id
		WHERE p.pharmacy_id = $1
		AND oi.created_at <= $2 AND oi.created_at >= $3
		AND op.order_status_id = 5
		GROUP BY dc.drug_category_id
	`

	FindPharmacyDrugSalesVolumeRevenue = `
		SELECT d.drug_id, d.drug_name, COUNT (d.drug_id) sales_volume, SUM (oi.quantity * oi.drug_price) revenue
		FROM order_items oi
		JOIN drugs d ON d.drug_id = oi.drug_id
		JOIN order_pharmacies op ON op.order_pharmacy_id = oi.order_pharmacy_id
		JOIN pharmacy_couriers pc ON pc.pharmacy_courier_id = op.pharmacy_courier_id
		JOIN pharmacies p ON p.pharmacy_id = pc.pharmacy_id
		WHERE p.pharmacy_id = $1
		AND oi.created_at <= $2 AND oi.created_at >= $3
		AND op.order_status_id = 5
		GROUP BY d.drug_id 
	`
)
