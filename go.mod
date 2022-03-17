module example.com

go 1.17

require (
	example.com/CloudServerMidterm/service v0.0.0
	github.com/gorilla/mux v1.8.0
)

require (
	github.com/go-sql-driver/mysql v1.6.0 // indirect
	github.com/joho/godotenv v1.4.0 // indirect
)

replace example.com/CloudServerMidterm/service => ../CloudServerMidterm/service
