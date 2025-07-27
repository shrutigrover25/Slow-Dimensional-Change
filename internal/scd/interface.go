package scd

type SCDModel[T any] interface {
	TableName() string
	GetID() string
	GetUID() string
	GetVersion() int
	CopyForNewVersion() T
}
