package data_json

type ZoneNotification struct {
	Imei         string   `json:"imei"`
	ZoneName     string   `json:"zone_name"`
	ZoneID       int      `json:"zone_id"`
	EventType    string   `json:"event_type"`
	VehiculoId   int      `json:"vehiculo_id"`
	VehiculoName string   `json:"vehiculo_name"`
	Emails       []string `json:"emails"`
	FCMTokens    []string `json:"fcmtokens"`
}
