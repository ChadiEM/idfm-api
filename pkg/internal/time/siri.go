package time

import "time"

type Siri struct {
	ServiceDelivery ServiceDelivery `json:"ServiceDelivery"`
}

type ServiceDelivery struct {
	ResponseTimestamp         time.Time                `json:"ResponseTimestamp"`
	ProducerRef               string                   `json:"ProducerRef"`
	ResponseMessageIdentifier string                   `json:"ResponseMessageIdentifier"`
	StopMonitoringDelivery    []StopMonitoringDelivery `json:"StopMonitoringDelivery"`
}

type StopMonitoringDelivery struct {
	ResponseTimestamp  time.Time            `json:"ResponseTimestamp"`
	Version            string               `json:"Version"`
	Status             string               `json:"Status"`
	MonitoredStopVisit []MonitoredStopVisit `json:"MonitoredStopVisit"`
	StopLineNotice     []interface{}        `json:"StopLineNotice"`
	ServiceException   []interface{}        `json:"ServiceException"`
}

type MonitoredStopVisit struct {
	RecordedAtTime          time.Time               `json:"RecordedAtTime"`
	ItemIdentifier          string                  `json:"ItemIdentifier"`
	MonitoringRef           ValueWrapper            `json:"MonitoringRef"`
	MonitoredVehicleJourney MonitoredVehicleJourney `json:"MonitoredVehicleJourney"`
}

type MonitoredVehicleJourney struct {
	LineRef                 ValueWrapper            `json:"LineRef"`
	OperatorRef             map[string]interface{}  `json:"OperatorRef"`
	FramedVehicleJourneyRef FramedVehicleJourneyRef `json:"FramedVehicleJourneyRef"`
	DirectionName           []interface{}           `json:"DirectionName"`
	DestinationRef          ValueWrapper            `json:"DestinationRef"`
	DestinationName         []ValueWrapper          `json:"DestinationName"`
	VehicleJourneyName      []ValueWrapper          `json:"VehicleJourneyName"`
	JourneyNote             []ValueWrapper          `json:"JourneyNote"`
	MonitoredCall           MonitoredCall           `json:"MonitoredCall"`
	TrainNumbers            TrainNumbers            `json:"TrainNumbers"`
	VehicleFeatureRef       []string                `json:"VehicleFeatureRef"`
	DirectionRef            ValueWrapper            `json:"DirectionRef"`
}

type FramedVehicleJourneyRef struct {
	DataFrameRef           ValueWrapper `json:"DataFrameRef"`
	DatedVehicleJourneyRef string       `json:"DatedVehicleJourneyRef"`
}

type MonitoredCall struct {
	StopPointName           []ValueWrapper `json:"StopPointName"`
	VehicleAtStop           bool           `json:"VehicleAtStop"`
	DestinationDisplay      []ValueWrapper `json:"DestinationDisplay"`
	ArrivalStopAssignment   StopAssignment `json:"ArrivalStopAssignment"`
	DepartureStopAssignment StopAssignment `json:"DepartureStopAssignment"`
	ExpectedArrivalTime     time.Time      `json:"ExpectedArrivalTime"`
	ExpectedDepartureTime   time.Time      `json:"ExpectedDepartureTime"`
	DepartureStatus         string         `json:"DepartureStatus"`
	Order                   int            `json:"Order"`
	AimedArrivalTime        time.Time      `json:"AimedArrivalTime"`
	ArrivalPlatformName     ValueWrapper   `json:"ArrivalPlatformName"`
	AimedDepartureTime      time.Time      `json:"AimedDepartureTime"`
	ArrivalStatus           string         `json:"ArrivalStatus"`
}

type StopAssignment struct {
	ExpectedQuayRef ValueWrapper `json:"ExpectedQuayRef"`
}

type TrainNumbers struct {
	TrainNumberRef []ValueWrapper `json:"TrainNumberRef"`
}

type ValueWrapper struct {
	Value string `json:"value"`
}
