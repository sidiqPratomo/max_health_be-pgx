package database

const (
	GetAllDoctorSpecializationQuery = `
		SELECT specialization_id, specialization_name
		FROM doctor_specializations
		WHERE deleted_at IS NULL
	`
)