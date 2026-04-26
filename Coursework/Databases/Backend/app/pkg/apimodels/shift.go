package apimodels

import "time"

type DutyEmployee struct {
	ID       string `json:"id"`
	FullName string `json:"full_name"`
	Login    string `json:"login"`
}

type Shift struct {
	ID       string        `json:"id"`
	OpenedAt time.Time     `json:"opened_at"`
	ClosedAt *time.Time    `json:"closed_at"`
	IsActive bool          `json:"is_active"`
	Duty     *DutyEmployee `json:"duty"`
}

type ActiveShiftResponse struct {
	IsActive bool   `json:"is_active"`
	Shift    *Shift `json:"shift"`
}

type OpenShiftRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type OpenShiftResponse struct {
	Shift *Shift `json:"shift"`
}

type CloseShiftResponse struct {
	Shift *Shift `json:"shift"`
}
