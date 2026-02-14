package sqlstorage

import (
	"calendar/internal/storage"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostgresStorage_Integration(t *testing.T) {
	// Пропускаем если тесты запускаются в коротком режиме
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Конфигурация для тестов
	cfg := Config{
		Host:     "localhost",
		Port:     5432,
		Database: "calendar",
		Username: "calendar_user",
		Password: "calendar_pass",
		SSLMode:  "disable",
	}

	// Создаем хранилище
	store := New(cfg)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	store.Connect(ctx)
	defer store.Close(ctx)

	// Очищаем тестовые данные перед тестом
	cleanupTestData(t, store)

	t.Run("Complete event lifecycle test", func(t *testing.T) {
		// Создаем события
		now := time.Now()
		event1 := &storage.Event{
			ID:          "test-user1-1",
			Title:       "User 1 Event",
			Description: "First event for user 1",
			StartTime:   now.Add(24 * time.Hour),
			EndTime:     now.Add(25 * time.Hour),
			UserID:      "user-1",
		}

		event2 := &storage.Event{
			ID:          "test-user2-1",
			Title:       "User 2 Event",
			Description: "First event for user 2",
			StartTime:   now.Add(24 * time.Hour),
			EndTime:     now.Add(25 * time.Hour),
			UserID:      "user-2",
		}

		event3 := &storage.Event{
			ID:          "test-user1-2",
			Title:       "User 1 Second Event",
			Description: "Second event for user 1",
			StartTime:   now.Add(26 * time.Hour),
			EndTime:     now.Add(27 * time.Hour),
			UserID:      "user-1",
		}

		// 1. Добавляем первое событие
		err := store.CreateEvent(ctx, event1)
		require.NoError(t, err, "Failed to create event1")

		// 2. Проверяем только ID через GetEvent
		gotEvent1, err := store.GetEvent(ctx, event1.ID)
		assert.NoError(t, err, "Failed to get event1")
		assert.NotNil(t, gotEvent1, "Event1 should not be nil")
		assert.Equal(t, event1.ID, gotEvent1.ID, "Event1 ID mismatch")

		// 3. Добавляем событие другого пользователя
		err = store.CreateEvent(ctx, event2)
		require.NoError(t, err, "Failed to create event2")

		// 4. Добавляем второе событие первого пользователя
		err = store.CreateEvent(ctx, event3)
		require.NoError(t, err, "Failed to create event3")

		// 5. Получаем список для первого пользователя
		user1Events, err := store.ListEvents(ctx, "user-1")
		assert.NoError(t, err, "Failed to list events for user1")

		// 6. Проверяем список для первого пользователя
		assert.Len(t, user1Events, 2, "User1 should have 2 events")

		// Проверяем ID событий
		foundIDs := make(map[string]bool)
		for _, e := range user1Events {
			foundIDs[e.ID] = true
		}
		assert.True(t, foundIDs["test-user1-1"], "Event test-user1-1 not found")
		assert.True(t, foundIDs["test-user1-2"], "Event test-user1-2 not found")

		// 7. Изменяем самый первый event: userId на второго пользователя
		event1.UserID = "user-2"
		err = store.UpdateEvent(ctx, event1)
		assert.NoError(t, err, "Failed to update event1")

		// 8. Получаем список для первого пользователя (должен быть 1 event)
		user1Events, err = store.ListEvents(ctx, "user-1")
		assert.NoError(t, err, "Failed to list events for user1 after update")
		assert.Len(t, user1Events, 1, "User1 should have 1 event after update")
		assert.Equal(t, "test-user1-2", user1Events[0].ID, "Remaining event should be test-user1-2")

		// 9. Получаем список для второго пользователя (должен быть 2 event)
		user2Events, err := store.ListEvents(ctx, "user-2")
		assert.NoError(t, err, "Failed to list events for user2 after update")
		assert.Len(t, user2Events, 2, "User2 should have 2 events after update")

		// Проверяем ID событий у второго пользователя
		foundIDs = make(map[string]bool)
		for _, e := range user2Events {
			foundIDs[e.ID] = true
		}
		assert.True(t, foundIDs["test-user2-1"], "Event test-user2-1 not found")
		assert.True(t, foundIDs["test-user1-1"], "Event test-user1-1 not found")

		// 10. Удаляем event первого пользователя (последний оставшийся)
		err = store.DeleteEvent(ctx, "test-user1-2")
		assert.NoError(t, err, "Failed to delete event test-user1-2")

		// 11. Получаем список для первого пользователя (должен быть 0 event)
		user1Events, err = store.ListEvents(ctx, "user-1")
		assert.NoError(t, err, "Failed to list events for user1 after deletion")
		assert.Len(t, user1Events, 0, "User1 should have 0 events after deletion")

		// 12. Проверяем что при попытке получить несуществующий event возвращается ошибка
		_, err = store.GetEvent(ctx, "test-user1-2")
		assert.ErrorIs(t, err, storage.ErrEventNotFound, "Should return ErrEventNotFound for deleted event")

		// 13. Получаем список для второго пользователя (должен быть 2 event)
		user2Events, err = store.ListEvents(ctx, "user-2")
		assert.NoError(t, err, "Failed to list events for user2 after deletion")
		assert.Len(t, user2Events, 2, "User2 should still have 2 events")

		// Проверяем ID событий у второго пользователя
		foundIDs = make(map[string]bool)
		for _, e := range user2Events {
			foundIDs[e.ID] = true
		}
		assert.True(t, foundIDs["test-user2-1"], "Event test-user2-1 not found")
		assert.True(t, foundIDs["test-user1-1"], "Event test-user1-1 not found")

		// 14. Проверяем что Event2 и Event1 (теперь user-2) все еще существуют
		gotEvent2, err := store.GetEvent(ctx, "test-user2-1")
		assert.NoError(t, err, "Failed to get event2")
		assert.NotNil(t, gotEvent2, "Event2 should not be nil")
		assert.Equal(t, "user-2", gotEvent2.UserID, "Event2 user ID should be user-2")

		gotEvent1, err = store.GetEvent(ctx, "test-user1-1")
		assert.NoError(t, err, "Failed to get event1 after user change")
		assert.NotNil(t, gotEvent1, "Event1 should not be nil")
		assert.Equal(t, "user-2", gotEvent1.UserID, "Event1 user ID should now be user-2")
	})
}

// Вспомогательная функция для очистки тестовых данных
func cleanupTestData(t *testing.T, store *SQLStorage) {
	ctx := context.Background()
	_, err := store.db.ExecContext(ctx, "DELETE FROM events WHERE id LIKE 'test-%'")
	require.NoError(t, err, "Failed to clean test data")
}
