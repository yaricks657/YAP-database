package main

import (
	"database/sql"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	// randSource источник псевдо случайных чисел.
	// Для повышения уникальности в качестве seed
	// используется текущее время в unix формате (в виде числа)
	randSource = rand.NewSource(time.Now().UnixNano())
	// randRange использует randSource для генерации случайных чисел
	randRange = rand.New(randSource)
)

// getTestParcel возвращает тестовую посылку
func getTestParcel() Parcel {
	return Parcel{
		Client:    1000,
		Status:    ParcelStatusRegistered,
		Address:   "test",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

// TestAddGetDelete проверяет добавление, получение и удаление посылки
func TestAddGetDelete(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db") // настройте подключение к БД
	if err != nil {
		require.NoError(t, err)
		return
	}
	defer db.Close()

	// подготовка входных данных
	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	res, err := store.Add(parcel)

	require.NoError(t, err)
	require.NotEmpty(t, res)

	// get
	item, err := store.Get(parcel.Number)

	require.NoError(t, err)
	assert.Equal(t, parcel.Number, item.Number)
	assert.Equal(t, parcel.Client, item.Client)
	assert.Equal(t, parcel.Address, item.Address)
	assert.Equal(t, parcel.CreatedAt, item.CreatedAt)
	assert.Equal(t, parcel.Status, item.Status)

	// delete
	// удалите добавленную посылку, убедитесь в отсутствии ошибки
	// проверьте, что посылку больше нельзя получить из БД
	err = store.Delete(parcel.Number)

	require.NoError(t, err)
	element, err := store.Get(parcel.Number)
	require.NoError(t, err)
	assert.NotEmpty(t, element)
}

// TestSetAddress проверяет обновление адреса
func TestSetAddress(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db") // настройте подключение к БД
	if err != nil {
		require.NoError(t, err)
		return
	}
	defer db.Close()

	// подготовка входных данных
	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	res, err := store.Add(parcel)

	require.NoError(t, err)
	require.NotEmpty(t, res)

	// set address
	newAddress := "new test address"
	err = store.SetAddress(parcel.Number, newAddress)

	require.NoError(t, err)

	// check
	item, err := store.Get(parcel.Number)

	require.NoError(t, err)
	require.NotEmpty(t, item)
	assert.Equal(t, newAddress, item.Address)
}

// TestSetStatus проверяет обновление статуса
func TestSetStatus(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db") // настройте подключение к БД
	if err != nil {
		require.NoError(t, err)
		return
	}
	defer db.Close()

	// подготовка входных данных
	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	res, err := store.Add(parcel)

	require.NoError(t, err)
	require.NotEmpty(t, res)

	// set status
	newStatus := "new test status"
	err = store.SetStatus(parcel.Number, newStatus)

	require.NoError(t, err)

	// check
	item, err := store.Get(parcel.Number)

	require.NoError(t, err)
	require.NotEmpty(t, item)
	assert.Equal(t, newStatus, item.Status)
}

// TestGetByClient проверяет получение посылок по идентификатору клиента
func TestGetByClient(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db") // настройте подключение к БД
	if err != nil {
		require.NoError(t, err)
		return
	}
	defer db.Close()

	parcels := []Parcel{
		getTestParcel(),
		getTestParcel(),
		getTestParcel(),
	}
	parcelMap := map[int]Parcel{}

	// задаём всем посылкам один и тот же идентификатор клиента
	client := randRange.Intn(10_000_000)
	parcels[0].Client = client
	parcels[1].Client = client
	parcels[2].Client = client

	store := NewParcelStore(db)

	// add
	for i := 0; i < len(parcels); i++ {
		id, err := store.Add(parcels[i])

		require.NoError(t, err)
		require.NotEmpty(t, id)

		// обновляем идентификатор добавленной у посылки
		parcels[i].Number = id

		// сохраняем добавленную посылку в структуру map, чтобы её можно было легко достать по идентификатору посылки
		parcelMap[id] = parcels[i]
	}

	// get by client
	storedParcels, err := store.GetByClient(client)

	require.NoError(t, err)
	require.Equal(t, len(parcels), len(storedParcels))

	// check
	for _, parcel := range storedParcels {
		item, exists := parcelMap[parcel.Number]

		require.Equal(t, true, exists)
		assert.Equal(t, item.Number, parcel.Number)
		assert.Equal(t, item.Address, parcel.Address)
		assert.Equal(t, item.Client, parcel.Client)
		assert.Equal(t, item.CreatedAt, parcel.CreatedAt)
		assert.Equal(t, item.Status, parcel.Status)
	}
}
