package database

const (
	CreateStockChanges = `
		INSERT INTO stock_changes (pharmacy_drug_id, final_stock, amount, description)
		VALUES
	`

	CreateStockMutations = `
		INSERT INTO stock_mutation_requests (pharmacy_requester_id, pharmacy_target_id, drug_id, stock, status_id)
		VALUES
	`

	GetTwoClosestAvailableStockBase1 = `
		WITH detailed_cart AS (
			SELECT ci.cart_item_id, ci.pharmacy_drug_id, d.drug_id, pd.stock, p.pharmacy_id, p.geom, p.pharmacy_manager_id, ci.quantity
			FROM cart_items ci
			JOIN pharmacy_drugs pd
			ON ci.pharmacy_drug_id = pd.pharmacy_drug_id
			JOIN drugs d 
			ON d.drug_id = pd.drug_id
			JOIN pharmacies p
			ON p.pharmacy_id = pd.pharmacy_id
	`

	GetTwoClosestAvailableStockBase2 = `
		alternatives AS (
			SELECT dc.cart_item_id, dc.quantity AS cart_quantity, dc.drug_id, dc.pharmacy_drug_id AS original_pharmacy_drug, dc.stock AS origin_stock, 
				pd.stock AS alternative_stock, dc.pharmacy_id AS pharmacy_1, pd.pharmacy_drug_id AS alternative,
				p.pharmacy_id AS pharmacy_2, ST_DistanceSphere(dc.geom, p.geom) AS distance,
				ROW_NUMBER() OVER(PARTITION BY cart_item_id ORDER BY ST_DistanceSphere(dc.geom, p.geom)) AS rank
			FROM detailed_cart dc
			JOIN pharmacy_drugs pd
			ON dc.drug_id = pd.drug_id AND dc.pharmacy_drug_id != pd.pharmacy_drug_id
			JOIN pharmacies p
			ON p.pharmacy_id = pd.pharmacy_id AND p.pharmacy_manager_id = dc.pharmacy_manager_id
			WHERE ST_DistanceSphere(dc.geom, p.geom) <= 25000 AND pd.stock > 0
			ORDER BY dc.cart_item_id, ST_DistanceSphere(dc.geom, p.geom))`

	GetTwoClosestAvailableStockLock = `
		SELECT * FROM pharmacy_drugs pd
		JOIN alternatives a ON pd.pharmacy_drug_id = a.alternative
		FOR UPDATE
	`

	GetTwoClosestAvailableStockList = `
		SELECT cart_item_id, cart_quantity, drug_id, original_pharmacy_drug, pharmacy_1, origin_stock, alternative, pharmacy_2, alternative_stock
		FROM alternatives
		WHERE rank <= 2
		ORDER BY cart_item_id, rank
	`

	GetStockChanges = `
			SELECT p.pharmacy_name, p.address, d.drug_name, d.image, sc.final_stock, sc.amount, sc.description
			FROM stock_changes sc
			JOIN pharmacy_drugs pd 
			ON pd.pharmacy_drug_id = sc.pharmacy_drug_id
			JOIN pharmacies p
			ON pd.pharmacy_id = p.pharmacy_id
			JOIN drugs d
			ON d.drug_id = pd.drug_id
			WHERE p.pharmacy_manager_id = $1 AND sc.deleted_at IS NULL AND pd.deleted_at IS NULL AND p.deleted_at IS NULL AND d.deleted_at IS NULL
	`
)