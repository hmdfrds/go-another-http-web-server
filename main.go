package main

type Config struct {
	Host         string `json:"host"`
	Port         int    `json:"port"`
	AdminPort    int    `json:"admin_port"`
	DocumentRoot string `json:"document_root"`
	MaxThreads   int    `json:"max_threads"`
	LogFile      string `json:"log_file"`
}





func main() {
	print("Hello")
}
