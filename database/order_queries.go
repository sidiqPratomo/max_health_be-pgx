package database

const (
	FindAllOrders = `
		SELECT DISTINCT o.order_id, o.created_at, count(o.order_id) AS total_item_count
		FROM orders o
		JOIN order_pharmacies op on op.order_id = o.order_id
		WHERE o.deleted_at IS NULL
	`

	FindAllOrdersWithDetails = `
		SELECT DISTINCT o.order_id, o.total_amount, o.payment_proof, o.address, o.updated_at, op.order_pharmacy_id, op.order_status_id, op.subtotal_amount, op.delivery_fee, p.pharmacy_name, a.profile_picture, c.courier_name, COUNT(oi.order_item_id) OVER (PARTITION BY op.order_pharmacy_id) as total_order_items
		FROM orders o
		JOIN order_pharmacies op on op.order_id = o.order_id
		JOIN pharmacy_couriers pc ON pc.pharmacy_courier_id = op.pharmacy_courier_id
		JOIN couriers c ON c.courier_id = pc.courier_id
		JOIN pharmacies p ON p.pharmacy_id = pc.pharmacy_id
		JOIN pharmacy_managers pm ON pm.pharmacy_manager_id = p.pharmacy_manager_id
		JOIN accounts a ON a.account_id = pm.account_id
		JOIN order_items oi ON oi.order_pharmacy_id = op.order_pharmacy_id
		WHERE o.deleted_at IS NULL
	`

	FindAllPendingOrdersByUserId = `
		SELECT DISTINCT o.order_id, o.created_at, count(*) OVER() AS total_item_count
		FROM orders o
		JOIN order_pharmacies op on op.order_id = o.order_id
		WHERE o.user_id = $1
		AND op.order_status_id = 1
		AND o.deleted_at IS NULL
		GROUP BY o.order_id
		ORDER BY o.created_at DESC
		LIMIT $2
		OFFSET $3
	`

	FindAllPendingOrdersWithDetailsByUserId = `
		SELECT DISTINCT o.order_id, o.total_amount, op.order_pharmacy_id, op.order_status_id, op.subtotal_amount, op.delivery_fee, p.pharmacy_name, a.profile_picture, COUNT(oi.order_item_id) OVER (PARTITION BY op.order_pharmacy_id) as total_order_items
		FROM orders o
		JOIN order_pharmacies op on op.order_id = o.order_id
		JOIN pharmacy_couriers pc ON pc.pharmacy_courier_id = op.pharmacy_courier_id
		JOIN pharmacies p ON p.pharmacy_id = pc.pharmacy_id
		JOIN pharmacy_managers pm ON pm.pharmacy_manager_id = p.pharmacy_manager_id
		JOIN accounts a ON a.account_id = pm.account_id
		JOIN order_items oi ON oi.order_pharmacy_id = op.order_pharmacy_id
		WHERE o.user_id = $1
		AND op.order_status_id = 1
		AND o.deleted_at IS NULL
	`

	CreateOneOrder = `
		INSERT INTO orders(user_id, address, total_amount)
		VALUES
		($1, $2, $3)
		RETURNING order_id
	`

	CreateOrderPharmacies = `
		INSERT INTO order_pharmacies(order_id, order_status_id, pharmacy_courier_id, subtotal_amount, delivery_fee)
		VALUES
	`

	CreateOrderItems = `
		INSERT INTO order_items(order_pharmacy_id, drug_id, drug_name, pharmacy_drug_id, drug_price, drug_unit, quantity)
		VALUES
	`

	UpdatePaymentProofOneOrder = `
		UPDATE orders 
		SET payment_proof = $1,
		updated_at = NOW()
		WHERE order_id = $2
	`

	GetOneOrderByOrderId = `
		SELECT order_id, user_id, address, payment_proof, total_amount, expired_at
		FROM orders
		WHERE order_id = $1 AND deleted_at IS NULL
	`
)
