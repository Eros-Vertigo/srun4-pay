package configs

// Config 加载yaml配置文件
type Config interface {
	Load()
}
