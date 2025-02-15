package database

const (
	GetCouriers = `
		SELECT courier_id, courier_name 
		FROM couriers
		WHERE deleted_at IS NULL
	`
)
