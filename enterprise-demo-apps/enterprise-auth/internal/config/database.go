package config

import (
	"log"
)

func ConnectDatabase() {
	log.Println("Enterprise-Auth runs in stateless/event-driven mode (No local user database needed).")
}
