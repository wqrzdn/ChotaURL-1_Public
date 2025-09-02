// package storages allows multiple implementation on how to store short URLs.
package storages

type IStorage interface {
	Code() string
	Save(string) string
	SaveWithCustom(url, customName string) (string, error)
	Load(string) (string, error)
	Exists(code string) bool
}
