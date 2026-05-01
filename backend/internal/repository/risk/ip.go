package risk

import (
	"time"
	"tiny-forum/internal/model/po"
)

func (r *riskRepository) CountActiveRiskEventsByIP(ip string) (int, error) {
	var count int64
	err := r.db.Model(&po.IPRiskRecord{}).
		Where("ip = ? AND expire_at > ?", ip, time.Now()).
		Count(&count).Error
	return int(count), err
}

func (r *riskRepository) AddIPRiskRecord(record *po.IPRiskRecord) error {
	return r.db.Create(record).Error
}

func (r *riskRepository) IsIPBlocked(ip string) (bool, error) {
	var count int64
	err := r.db.Model(&po.BlockedIP{}).
		Where("ip = ? AND (expire_at IS NULL OR expire_at > ?)", ip, time.Now()).
		Count(&count).Error
	return count > 0, err
}
