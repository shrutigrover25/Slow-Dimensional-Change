package scd

import (
  "gorm.io/gorm"
)

type SCDManager[T SCDModel[T]] struct {
	db *gorm.DB
}

func NewManager[T SCDModel[T]](db *gorm.DB) *SCDManager[T] {
  return &SCDManager[T]{db: db}
}

// Get only the latest versions
func (m *SCDManager[T]) GetLatest() *gorm.DB {
  var dummy T
  t := dummy.TableName()
  sub := m.db.Table(t+" as s").
    Select("id, MAX(version) as max_ver").
    Group("id")
  return m.db.Table(t+" as main").
    Joins("JOIN (?) as latest ON main.id = latest.id AND main.version = latest.max_ver", sub)
}

func (m *SCDManager[T]) FindByUID(uid string) (T, error) {
  var entity T
  err := m.db.Where("uid = ?", uid).First(&entity).Error
  return entity, err
}

func (m *SCDManager[T]) Insert(newItem T) error {
  return m.db.Create(&newItem).Error
}

func (m *SCDManager[T]) CreateNewVersion(old T) (T, error) {
  newItem := old.CopyForNewVersion()
  return newItem, m.db.Create(&newItem).Error
}



// TestCases 
// Swagger Docs (optional)
// LOOM VIDEO
// APIS COMPLETION 
// POSTMAN APIS