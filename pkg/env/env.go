package env

import "os"

var (
	IDFM_API_KEY = os.Getenv("IDFM_API_KEY")
)
