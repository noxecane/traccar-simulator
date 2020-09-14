package traccar

import (
	"context"
	"time"

	"github.com/go-pg/pg/v9"
)

type Device struct {
	tableName struct{} `pg:"tc_devices"`
	ID        uint
	UpdatedAt time.Time `pg:"lastupdate"`

	Name         string
	ExternalID   string `pg:"uniqueid"`
	LastPosition uint   `pg:"positionid"`
	Payload      string `pg:"attributes"`

	Phone    string
	Model    string
	Contact  string
	Category string
	Disabled bool
	Group    uint `pg:"groupid"`
}

type Event struct {
	tableName struct{} `pg:"tc_events"`
	ID        uint
	CreatedAt time.Time `pg:"servertime"`

	Type     string
	Device   uint   `pg:"deviceid"`
	Position uint   `pg:"position"`
	Payload  string `pg:"attributes"`

	GeoFence    uint `pg:"geofenceid"`
	Maintenance uint `pg:"maintenanceid"`
}

type Position struct {
	tableName  struct{} `pg:"tc_positions"`
	ID         uint
	CreatedAt  time.Time `pg:"servertime"`
	RecordedAt time.Time `pg:"devicetime"`
	Valid      bool      `pg:",use_zero"`
	Device     uint      `pg:"deviceid"`
	Latitude   float64   `pg:",use_zero"`
	Longitude  float64   `pg:",use_zero"`
	Altitude   float64   `pg:",use_zero"`
	Speed      float64   `pg:",use_zero"`
	Course     float64   `pg:",use_zero"`
	Payload    string    `pg:"attributes"`
	Accuracy   string
	Address    string
	Protocol   string
	Network    string
	FixedAt    time.Time `pg:"fixtime"`
}

type Attributes struct {
	FuelConsumption     float32 `json:"fuelConsumption,omitempty"`
	Raw                 string  `json:"raw,omitempty"`
	GSensor             string  `json:"gSensor,omitempty"`
	Result              string  `json:"result,omitempty"`
	Status              uint    `json:"status,omitempty"`
	Motion              bool    `json:"motion,omitempty"`
	ClearedDistance     float32 `json:"clearedDistance,omitempty"`
	TotalDistance       float32 `json:"totalDistance,omitempty"`
	RPM                 uint    `json:"rpm,omitempty"`
	Alarm               string  `json:"alarm,omitempty"`
	Ignition            bool    `json:"ignition,omitempty"`
	DTC                 string  `json:"dtcs,omitempty"`
	OBDSpeed            uint    `json:"obdSpeed,omitempty"`
	EngineLoad          int     `json:"engineLoad,omitempty"`
	CoolantTemperature  int     `json:"coolantTemp,omitempty"`
	Distance            float32 `json:"distance,omitempty"`
	TripOdometer        uint    `json:"tripOdometer,omitempty"`
	IntakeTemperature   int     `json:"intakeTemp,omitempty"`
	Odometer            uint64  `json:"odometer,omitempty"`
	MapIntake           int     `json:"mapIntake,omitempty"`
	Throttle            float32 `json:"throttle,omitempty"`
	MilDistance         float32 `json:"milDistance,omitempty"`
	Satellites          uint    `json:"sat,omitempty"`
	TripFuelConsumption float32 `json:"tripFuelConsumption,omitempty"`
}

func getDevice(db *pg.DB, id uint) (*Device, error) {
	device := &Device{ID: id}
	if err := db.Model(device).WherePK().Select(); err != nil {
		return nil, err
	}

	return device, nil
}

func getPosition(db *pg.DB, id uint) (*Position, error) {
	position := &Position{ID: id}
	if err := db.Model(position).WherePK().Select(); err != nil {
		return nil, err
	}

	return position, nil
}

func getPositions(ctx context.Context, db *pg.DB, id, device uint) ([]Position, error) {
	var positions []Position
	err := db.
		ModelContext(ctx, &positions).
		Where("id > ?", id).
		Where("deviceid = ?", device).
		Select()

	return positions, err
}

func createDevices(db *pg.DB, devices []Device) error {
	return db.Insert(&devices)
}

func createPosition(db *pg.DB, position Position) error {
	return db.Insert(&position)
}
