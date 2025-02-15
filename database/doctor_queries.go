package database

const (
	PostOneDoctorQuery = `
		INSERT 
		INTO doctors (account_id, specialization_id, certificate, fee_per_patient, is_online)
		VALUES ($1, $2, $3, 0, FALSE)
	`

	UpdateOneDoctorQuery = `
		UPDATE doctors 
		SET fee_per_patient = $1,
		experience = $2,
		updated_at = NOW()
		WHERE account_id = $3
	`

	FindSpecializationById = `
		SELECT specialization_name 
		FROM doctor_specializations
		WHERE specialization_id = $1
	`

	FindDoctorByAccountIdQuery = `
		SELECT d.doctor_id, d.experience, d.specialization_id, ds.specialization_name, d.fee_per_patient, d.certificate
		FROM doctors d
		JOIN accounts a
		ON d.account_id = a.account_id
		LEFT JOIN doctor_specializations ds
		ON d.specialization_id = ds.specialization_id
		WHERE d.account_id = $1
		AND a.verified_at IS NOT NULL
		AND (d.deleted_at IS NULL AND a.deleted_at IS NULL)
	`

	FindDoctorByDoctorIdQuery = `
		SELECT d.doctor_id, a.email, a.account_name, a.profile_picture, d.experience, d.specialization_id, ds.specialization_name, d.fee_per_patient
		FROM doctors d
		JOIN accounts a
		ON d.account_id = a.account_id
		LEFT JOIN doctor_specializations ds
		ON d.specialization_id = ds.specialization_id
		WHERE d.doctor_id = $1
		AND a.verified_at IS NOT NULL
		AND (d.deleted_at IS NULL AND a.deleted_at IS NULL)
	`

	UpdateDoctorStatusQuery = `
		UPDATE doctors
		SET is_online = $1, updated_at = NOW()
		WHERE account_id = $2
		AND deleted_at IS NULL
	`

	GetDoctorIsOnlineQuery = `
		SELECT d.is_online
		FROM doctors d
		WHERE d.account_id = $1 
		AND deleted_at IS NULL
	`
)
