package models

type AISMessage struct {
	MsgType   int    `json:"msg_type"`
	IMO       int    `json:"imo"`
	MMSI      int    `json:"mmsi"`
	CALLSIGN  string `json:"callsign"`
	SHIPNAME  string `json:"shipname"`
	SHIP_TYPE string `json:"ship_type"`
}

// Estructura final enriquecida
type Ship struct {
	MsgType        int     `json:"msg_type"`
	IMO            int     `json:"imo"`
	MMSI           int     `json:"mmsi"`
	Callsign       string  `json:"callsign"`
	Shipname       string  `json:"shipname"`
	ShipType       string  `json:"ship_type"`
	BuiltYear      *string `json:"Built"`
	Shipyard       *string `json:"shipyard"`
	HullNumber     *string `json:"Hull-No."`
	KeelLaying     *string `json:"Keel Laying"`
	LaunchDate     *string `json:"Launch"`
	DeliveryDate   *string `json:"Delivery"`
	GT             *string `json:"gt"`
	NT             *string `json:"nt"`
	CarryingCapTDW *string `json:"Carrying capacity (tdw)"`
	LengthOverall  *string `json:"Length overall (m)"`
	Breadth        *string `json:"Breadth (m)"`
	Depth          *string `json:"Depth (m)"`
	Propulsion     *string `json:"propulsion"`
	Power          *string `json:"power"`
	Screws         *string `json:"screws"`
	Speed          *string `json:"speed"`
}
