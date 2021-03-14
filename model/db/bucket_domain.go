// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/3/4

package db

import "time"

type Bucket struct {
	Bucket string `json:"bucket" gorm:"primary_key:bucket;type:varchar(50);not null;"`
	Domain string `json:"domain" gorm:"column:domain;type:varchar(50);not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (b *Bucket) TableName() string {
	return "bucket_domain"
}
