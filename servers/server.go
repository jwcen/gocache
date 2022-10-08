// 用于存放服务器相关方法
package servers

// Server 是服务器结构的接口。
// 抽象出一个接口是因为后续还会提供 TCP 服务
type Server interface {
	// Run 在 address 上启动服务器，并返回错误信息
	Run(address string) error
}
