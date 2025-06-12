package database

type BaseConfig interface {
	Connect() error
}
