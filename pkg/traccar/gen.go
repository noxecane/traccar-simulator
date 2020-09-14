package traccar

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func ReadDump(path string) [][]string {
	var f *os.File
	var err error

	if f, err = os.Open(path); err != nil {
		panic(err)
	}
	defer f.Close()

	scn := bufio.NewScanner(f)

	var csvResult [][]string

	// skip the first line
	if scn.Scan() {
		for scn.Scan() {
			line := scn.Text()
			var cols []string

			for _, col := range strings.Fields(line) {
				if col == "\\N" || col == "{}" {
					cols = append(cols, "")
				} else {
					cols = append(cols, col)
				}
			}

			csvResult = append(csvResult, cols)
		}
	}

	if err = scn.Err(); err != nil {
		panic(err)
	}

	return csvResult
}

func ReadDeviceDump(path string) []Device {
	rows := ReadDump(path)

	var devices []Device

	// we live in Lagos
	zone, err := time.LoadLocation("Africa/Lagos")
	if err != nil {
		panic(err)
	}

	for _, row := range rows {
		id, _ := strconv.Atoi(row[0])
		p, _ := strconv.Atoi(row[5])
		updatedAt, _ := time.Parse(time.RFC3339, fmt.Sprintf("%sT%s+01:00", row[3], row[4]))

		devices = append(devices, Device{
			ID:           uint(id),
			UpdatedAt:    updatedAt.In(zone),
			Name:         row[1],
			ExternalID:   row[2],
			LastPosition: uint(p),
		})
	}

	return devices
}

func ReadEventDump(path string) []Event {
	rows := ReadDump(path)

	var events []Event

	// we live in Lagos
	zone, err := time.LoadLocation("Africa/Lagos")
	if err != nil {
		panic(err)
	}

	for _, row := range rows {
		id, _ := strconv.Atoi(row[0])
		dID, _ := strconv.Atoi(row[4])
		pID, _ := strconv.Atoi(row[5])
		createdAt, _ := time.Parse(time.RFC3339, fmt.Sprintf("%sT%s+01:00", row[2], row[3]))

		events = append(events, Event{
			ID:        uint(id),
			CreatedAt: createdAt.In(zone),
			Type:      row[1],
			Device:    uint(dID),
			Position:  uint(pID),
			Payload:   row[7],
		})
	}

	return events
}

func ReadPositionDump(path string) []Position {
	rows := ReadDump(path)

	var postions []Position

	// we live in Lagos
	zone, err := time.LoadLocation("Africa/Lagos")
	if err != nil {
		panic(err)
	}

	for _, row := range rows {
		id, _ := strconv.Atoi(row[0])
		dId, _ := strconv.Atoi(row[2])
		valid, _ := strconv.ParseBool(row[9])
		lat, _ := strconv.ParseFloat(row[10], 64)
		long, _ := strconv.ParseFloat(row[11], 64)
		speed, _ := strconv.ParseFloat(row[13], 64)
		course, _ := strconv.ParseFloat(row[14], 64)

		createdAt, _ := time.Parse(time.RFC3339, fmt.Sprintf("%sT%s+01:00", row[3], row[4]))
		recordedAt, _ := time.Parse(time.RFC3339, fmt.Sprintf("%sT%s+01:00", row[5], row[6]))
		fixedAt, _ := time.Parse(time.RFC3339, fmt.Sprintf("%sT%s+01:00", row[7], row[8]))

		var payload string
		rem := row[16:]
		if len(rem) == 3 {
			payload = row[16]
		} else {
			payload = strings.Join(rem[:len(rem)-2], " ")
		}

		postions = append(postions, Position{
			ID:         uint(id),
			CreatedAt:  createdAt.In(zone),
			RecordedAt: recordedAt.In(zone),
			Valid:      valid,
			Device:     uint(dId),
			Latitude:   lat,
			Longitude:  long,
			Altitude:   float64(0),
			Speed:      float64(speed),
			Course:     float64(course),
			Payload:    payload,

			Accuracy: "",
			Address:  "",
			Protocol: row[1],
			Network:  "",
			FixedAt:  fixedAt.In(zone),
		})
	}

	return postions
}
