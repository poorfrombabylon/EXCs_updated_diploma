package domain

import (
	"time"
)

type newFullModelFromFunc func(createdAt time.Time, updatedAt time.Time, lastCheckedAt time.Time) FullModel

type newFullModelFunc func() FullModel

var (
	NewFullModelFrom newFullModelFromFunc = newFullModelFrom
	NewFullModel     newFullModelFunc     = newFullModel
)

type FullModel struct {
	createdAt     time.Time
	updatedAt     time.Time
	lastCheckedAt time.Time
}

func newFullModel() FullModel {
	data := time.Now().In(time.UTC)

	return FullModel{
		createdAt:     data,
		updatedAt:     data,
		lastCheckedAt: data,
	}
}

func newFullModelFrom(
	createdAt time.Time,
	updatedAt time.Time,
	lastCheckedAt time.Time,
) FullModel {
	return FullModel{
		createdAt:     createdAt.In(time.UTC),
		updatedAt:     updatedAt.In(time.UTC),
		lastCheckedAt: lastCheckedAt.In(time.UTC),
	}
}

func (m *FullModel) Update() {
	m.updatedAt = time.Now().In(time.UTC)
}

func (m *FullModel) Check() {
	m.lastCheckedAt = time.Now().In(time.UTC)
}

func (m FullModel) GetCreatedAt() time.Time {
	return m.createdAt
}

func (m FullModel) GetUpdatedAt() time.Time {
	return m.updatedAt
}

func (m FullModel) GetLastCheckedAt() time.Time {
	return m.lastCheckedAt
}

type newModelFromFunc func(createdAt time.Time, updatedAt time.Time) Model

type newModelFunc func() Model

var (
	NewModelFrom newModelFromFunc = newModelFrom
	NewModel     newModelFunc     = newModel
)

type Model struct {
	createdAt time.Time
	updatedAt time.Time
}

func newModel() Model {
	data := time.Now().In(time.UTC)

	return Model{
		createdAt: data,
		updatedAt: data,
	}
}

func newModelFrom(
	createdAt time.Time,
	updatedAt time.Time,
) Model {
	return Model{
		createdAt: createdAt.In(time.UTC),
		updatedAt: updatedAt.In(time.UTC),
	}
}

func (m *Model) Update() {
	m.updatedAt = time.Now().In(time.UTC)
}

func (m Model) GetCreatedAt() time.Time {
	return m.createdAt
}

func (m Model) GetUpdatedAt() time.Time {
	return m.updatedAt
}
