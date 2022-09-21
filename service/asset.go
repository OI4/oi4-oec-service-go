package service

import oi4 "github.com/mzeiher/oi4/api/pkg/types"

type Oi4Asset struct {
	parent *Oi4Application
	mam    *oi4.MasterAssetModel
	health oi4.Health
	data   oi4.Oi4Data
}

func CreateNewAsset(mam *oi4.MasterAssetModel) *Oi4Asset {
	return &Oi4Asset{
		mam:    mam,
		health: oi4.Health{Health: oi4.Health_Normal, HealthScore: 100},
		data:   oi4.Oi4Data{PrimaryValue: nil},
	}
}
