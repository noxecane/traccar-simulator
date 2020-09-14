package traccar

import (
	"context"

	"github.com/go-pg/pg/v9"
)

func Listen(ctx context.Context, db *pg.DB, channel string, out chan<- []byte) {
	l := db.Listen(channel)
	defer l.Close()

	ch := l.Channel()
	for {
		select {
		case <-ctx.Done():
			return
		case n := <-ch:
			out <- []byte(n.Payload)
		}
	}
}

func GetPosition(db *pg.DB, id uint) (*Position, error) {
	return getPosition(db, id)
}

func GetPositions(ctx context.Context, db *pg.DB, id, device uint) ([]Position, error) {
	return getPositions(ctx, db, id, device)
}

func GetDevice(db *pg.DB, id uint) (*Device, error) {
	return getDevice(db, id)
}

func CreateDevices(db *pg.DB, devices []Device) error {
	return createDevices(db, devices)
}

func CreatePosition(db *pg.DB, p Position) error {
	return createPosition(db, p)
}
