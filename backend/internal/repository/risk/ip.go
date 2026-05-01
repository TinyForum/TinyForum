package risk

import (
	"time"
	"tiny-forum/internal/model/do"
)

func (r *riskRepository) CountActiveRiskEventsByIP(ip string) (int, error) {
	var count int64
	err := r.db.Model(&do.IPRiskRecord{}).
		Where("ip = ? AND expire_at > ?", ip, time.Now()).
		Count(&count).Error
	return int(count), err
}

func (r *riskRepository) AddIPRiskRecord(record *do.IPRiskRecord) error {
	return r.db.Create(record).Error
}

func (r *riskRepository) IsIPBlocked(ip string) (bool, error) {
	var count int64
	err := r.db.Model(&do.BlockedIP{}).
		Where("ip = ? AND (expire_at IS NULL OR expire_at > ?)", ip, time.Now()).
		Count(&count).Error
	return count > 0, err
}
