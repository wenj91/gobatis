package gobatis

type OptionType int

const (
	OptionTypeFile OptionType = 1
	OptionTypeDS   OptionType = 2
	OptionTypeDB   OptionType = 3
)

type IOption interface {
	Type() OptionType
	ToDBConf() *DBConfig
}
