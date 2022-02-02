# UserHub
Servicio de authenticacion de usuarios.

---
## Testing
Se usa `go test ./... -coverprofile coverage` para realizar los test.
Luego se usa `go tool cover -html=coverage` para mostrar el coverage de los tests.

---
## Configuracion
El archivo de configuracion es `config.json`, cada usuario tiene que configurarlo localmente siguiendo como ejemplo el archivo `config.json.example`
