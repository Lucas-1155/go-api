package schemas

import (
	"time"

	"github.com/google/uuid"
)

type QueryMetadata struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Alias        string    `gorm:"uniqueIndex;not null"` // Nome que o Grafana vai usar
	SQLEncrypted string    `gorm:"not null"`             // SQL guardado com AES
	TTLSeconds   int       `gorm:"default:30"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// TableName define o nome da tabela manualmente se necessário
func (QueryMetadata) TableName() string {
	return "queries"
}
