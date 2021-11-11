package models

type AuditAddRequest struct {
	HashArray   string `json:"hash_array"`
	Account     string `json:"account"`
	Description string `json:"description"`
	Algorithm   string `json:"algorithm"`
}

type AuditVerifyRequest struct {
	Hash string `json:"audit_hash"`
}
