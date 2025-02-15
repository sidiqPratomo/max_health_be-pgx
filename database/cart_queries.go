package database

const (
	GetCartById = `
		SELECT cart_item_id, user_id 
		FROM cart_items
		WHERE cart_item_id = $1 AND deleted_at IS NULL
	`

	GetCartsByIds = `
		SELECT cart_item_id, user_id 
		FROM cart_items
		WHERE deleted_at IS NULL
	`

	GetAllDeliveryFee1 = `WITH detailed_cart AS (
		SELECT pc.pharmacy_courier_id, p.geom AS origin_point, c.raja_ongkir_id AS origin_id, (d.weight * ci.quantity) AS total_weight, d.is_active, p.pharmacy_name, 
			co.courier_name, co.price, co.is_official, p.pharmacy_id
		FROM cart_items ci
		JOIN pharmacy_drugs pd 
		ON ci.pharmacy_drug_id = pd.pharmacy_drug_id
		JOIN pharmacies p
		ON pd.pharmacy_id = p.pharmacy_id
		JOIN drugs d
		ON pd.drug_id = d.drug_id
		JOIN cities c
		ON c.city_name ILIKE CONCAT('%', p.city)
		JOIN pharmacy_couriers pc
		ON p.pharmacy_id = pc.pharmacy_id
		JOIN couriers co
		ON pc.courier_id = co.courier_id
		WHERE ci.deleted_at ISNULL AND pd.deleted_at ISNULL AND p.deleted_at ISNULL AND d.deleted_at ISNULL AND pc.deleted_at ISNULL AND (`
	GetAllDeliveryFee2 = `address AS (
		SELECT c.raja_ongkir_id AS destination_id, ua.geom AS destination_point
		FROM user_addresses ua
		JOIN cities c 
		ON c.city_id = ua.city_id
		WHERE ua.user_address_id = $`
	GetAllDeliveryFee3 = `full_data AS (
		SELECT dc.pharmacy_courier_id, dc.origin_id, dc.total_weight, dc.is_active, a.destination_id, dc.pharmacy_name, dc.courier_name, 
			(dc.price * CEIL(ST_DistanceSphere(dc.origin_point, a.destination_point) / 1000)) AS total_price, dc.is_official, 
			CEIL(ST_DistanceSphere(dc.origin_point, a.destination_point) / 1000) AS distance, dc.pharmacy_id
		FROM detailed_cart dc, address a),
	pharmacies AS(
		SELECT pharmacy_id, pharmacy_name, courier_name, CEIL(SUM(total_weight)) AS total_weight
		FROM full_data
		GROUP BY pharmacy_name, courier_name, pharmacy_id),
	grouped_full_data AS (
		SELECT pharmacy_id, pharmacy_courier_id, origin_id, is_active, destination_id, total_price, is_official, distance, pharmacy_name, courier_name
		FROM full_data
		GROUP BY pharmacy_id, pharmacy_courier_id, pharmacy_name, courier_name, origin_id, is_active, destination_id, total_price, is_official, distance)
	SELECT p.pharmacy_id, p.pharmacy_name, fd.distance, fd.pharmacy_courier_id, p.courier_name, fd.origin_id, fd.destination_id, p.total_weight, fd.total_price, fd.is_active, fd.is_official
	FROM pharmacies p 
	JOIN grouped_full_data fd
	ON p.pharmacy_id = fd.pharmacy_id AND p.courier_name = fd.courier_name
	ORDER BY p.pharmacy_id ASC, fd.is_official DESC, p.courier_name`

	CheckUserQuery = `
		select user_id from users where account_id = $1
	`

	PostOneCartQuery = `
		INSERT INTO cart_items (user_id, pharmacy_drug_id,quantity)
		VALUES($1, $2, $3)
		RETURNING cart_item_id
	`

	UpdateOneCartQuery = `
		update cart_items set quantity = $3, updated_at = now() 
		where cart_item_id = $1 and user_id = $2
	`

	DeleteOneCartQuery = `
		update cart_items set updated_at = now(), deleted_at = now()
		where cart_item_id = $1 and user_id = $2
	`

	GetAllCartQuery = `
		Select ci.cart_item_id, ci.user_id, ci.pharmacy_drug_id, ci.quantity, pd.pharmacy_id, pd.drug_id, pd.price, pd.stock, d.drug_name, d.image, p.pharmacy_name
		from cart_items ci
		join pharmacy_drugs pd
		on ci.pharmacy_drug_id = pd.pharmacy_drug_id
		join drugs d
		on pd.drug_id = d.drug_id
		join pharmacies p
		on pd.pharmacy_id = p.pharmacy_id
		where ci.user_id = $1 and ci.deleted_at is null
		ORDER BY ci.cart_item_id ASC
		limit $2
		OFFSET $3
		`

	GetStockPharmacyDrug = `
		select stock 
		from pharmacy_drugs pd 
		join cart_items ci 
		on pd.pharmacy_drug_id = ci.pharmacy_drug_id 
		where ci.cart_item_id = $1
	`

	GetAllDetailedCartItems = `
		SELECT ci.cart_item_id, d.drug_id, d.drug_name, pd.pharmacy_drug_id, pd.price, d.unit_in_pack, ci.quantity
		FROM cart_items ci
		JOIN pharmacy_drugs pd
		ON pd.pharmacy_drug_id = ci.pharmacy_drug_id
		JOIN drugs d
		ON pd.drug_id = d.drug_id
	`

	GetAllCartsForChangesByCartIds = `
		SELECT ci.pharmacy_drug_id, pd.stock, ci.quantity
		FROM cart_items ci
		JOIN pharmacy_drugs pd
		ON ci.pharmacy_drug_id = pd.pharmacy_drug_id 
	`

	DeleteAllCarts = `
		UPDATE cart_items
		SET updated_at = NOW(),
		deleted_at = NOW()
	`
)
