package types

import (
	"log"
	"ms-common/utils"
)

type ServiceInfo struct {
	Name    string         `json:"name"`
	Address string         `json:"address"`
	Meta    map[string]any `json:"meta"`
}

func (si ServiceInfo) GetGrpcURI() string {
	port, err := utils.ToString(si.Meta["port"])
	if err != nil {
		log.Println(err)
		return ""
	}
	return si.Address + ":" + port
}
