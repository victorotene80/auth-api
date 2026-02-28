package entities

import "time"

type Role struct {
    ID          string
    Name        string
    Slug        string
    Description *string
    ParentID    *string
    IsSystem    bool
    CreatedAt   time.Time
    UpdatedAt   time.Time
}