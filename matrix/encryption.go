package matrix

import (
	"database/sql"
	"log"

	"maunium.net/go/mautrix/id"
)

func findDeviceId(db *sql.DB, accountId id.UserID) (deviceId id.DeviceID) {
	err := db.QueryRow("SELECT device_id FROM crypto_account WHERE account_id=$1", string(accountId)).Scan(&deviceId)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("Failed to scan device ID: %v", err)
	}
	return
}