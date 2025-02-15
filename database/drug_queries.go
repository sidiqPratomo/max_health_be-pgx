package database

const (
	GetDrugById = `
		select drug_id, drug_name ,generic_name ,content ,manufacture ,description ,classification_id,
		form_id, unit_in_pack ,selling_unit ,weight ,height ,length,width,image,drug_category_id,
		is_prescription_required,is_active
		from drugs where drug_id = $1
	`
	GetOneActiveDrugByIdQuery = `
		SELECT
			d.drug_id,
			d.drug_name,
			d.generic_name,
			d.content,
			d.manufacture,
			d.description,
			d.classification_id,
			dc.classification_name,
			d.form_id,
			df.form_name,
			d.unit_in_pack,
			d.selling_unit,
			d.weight,
			d.height,
			d.length,
			d.width,
			d.image,
			d.drug_category_id,
			dc2.drug_category_name,
			dc2.drug_category_url,
			d.is_prescription_required
		FROM drugs d
		JOIN drug_classifications dc ON dc.classification_id = d.classification_id
		JOIN drug_forms df ON df.form_id = d.form_id
		JOIN drug_categories dc2 ON dc2.drug_category_id = d.drug_category_id
		WHERE d.drug_id = $1 AND d.deleted_at IS NULL AND d.is_active is TRUE
	`

	GetDrugByNameQuery = `
		SELECT 
			drug_id,
			drug_name,
			generic_name,
			content,
			manufacture
		FROM drugs d
		WHERE d.drug_name ILIKE $1
		AND deleted_at IS NULL
	`

	UpdateOneDrugQuery = `
		UPDATE drugs
		SET
			drug_name = $2,
			generic_name = $3,
			content = $4,
			manufacture = $5,
			description = $6,
			classification_id = $7,
			form_id = $8,
			drug_category_id = $9,
			unit_in_pack = $10,
			selling_unit = $11,
			weight = $12,
			height = $13,
			length = $14,
			width = $15,
			image = $16,
			is_active = $17,
			is_prescription_required = $18
		WHERE drug_id = $1
	`

	GetAllDrugsByIdQuery = `
	SELECT
		d.drug_id,
		d.drug_name,
		d.generic_name,
		d.content,
		d.manufacture,
		d.description,
		d.classification_id,
		dc.classification_name,
		d.form_id,
		df.form_name,
		d.unit_in_pack,
		d.selling_unit,
		d.weight,
		d.height,
		d.length,
		d.width,
		d.image,
		d.drug_category_id,
		dc2.drug_category_name,
		dc2.drug_category_url,
		d.is_active,
		d.is_prescription_required,
		count(*) OVER() AS total_item_count
	FROM drugs d
	JOIN drug_classifications dc ON dc.classification_id = d.classification_id
	JOIN drug_forms df ON df.form_id = d.form_id
	JOIN drug_categories dc2 ON dc2.drug_category_id = d.drug_category_id
	WHERE d.deleted_at IS NULL
	AND d.drug_name ILIKE
`

	GetOneDrugByIdQuery = `
		SELECT
			d.drug_id,
			d.drug_name,
			d.generic_name,
			d.content,
			d.manufacture,
			d.description,
			d.classification_id,
			dc.classification_name,
			d.form_id,
			df.form_name,
			d.unit_in_pack,
			d.selling_unit,
			d.weight,
			d.height,
			d.length,
			d.width,
			d.image,
			d.drug_category_id,
			dc2.drug_category_name,
			dc2.drug_category_url,
			d.is_active,
			d.is_prescription_required
		FROM drugs d
		JOIN drug_classifications dc ON dc.classification_id = d.classification_id
		JOIN drug_forms df ON df.form_id = d.form_id
		JOIN drug_categories dc2 ON dc2.drug_category_id = d.drug_category_id
		WHERE d.drug_id = $1 AND d.deleted_at IS NULL
	`

	GetDrugIdByNameQuery = `
		SELECT d.drug_id
		FROM drugs d
		WHERE d.drug_name ILIKE $1
		AND deleted_at IS NULL
	`

	CreateOneDrugQuery = `
		INSERT INTO drugs (drug_name, generic_name, content, manufacture, description, classification_id, form_id, drug_category_id, unit_in_pack, selling_unit, weight, height, length, width, image, is_prescription_required, is_active) VALUES
		($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
	`

	DeleteOneDrugQuery = `
		UPDATE drugs
		SET deleted_at = NOW(), updated_at = NOW()
		WHERE drug_id = $1
	`

	GetDrugsByPharmacyId = `
	SELECT
		pd.pharmacy_drug_id,
		pd.price,
		pd.stock,
		d.drug_id,
		d.drug_name,
		d.generic_name,
		d.content,
		d.manufacture,
		d.description,
		d.classification_id,
		dc.classification_name,
		d.form_id,
		df.form_name,
		d.unit_in_pack,
		d.selling_unit,
		d.weight,
		d.height,
		d.length,
		d.width,
		d.image,
		d.drug_category_id,
		dc2.drug_category_name,
		dc2.drug_category_url,
		d.is_active,
		d.is_prescription_required
		FROM pharmacy_drugs pd
		JOIN drugs d ON pd.drug_id = d.drug_id
		JOIN drug_classifications dc ON dc.classification_id = d.classification_id
		JOIN drug_forms df ON df.form_id = d.form_id
		JOIN drug_categories dc2 ON dc2.drug_category_id = d.drug_category_id
		WHERE pd.pharmacy_id= $1 AND d.deleted_at IS NULL and pd.deleted_at is null and d.drug_name ILIKE '%' || $4 || '%'
		LIMIT $2
		OFFSET $3
	`
)
