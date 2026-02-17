package memorystorage

import (
	"calendar/internal/storage"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStorage_CreateEvent(t *testing.T) {
	t1 := &storage.Event{
		ID: "id",
	}

	t.Run("successful", func(t *testing.T) {
		t0 := New()
		err := t0.CreateEvent(context.Background(), t1)
		assert.NoError(t, err)

		savedEvent, exists := t0.events[t1.ID]
		assert.True(t, exists)
		assert.Equal(t, "id", savedEvent.ID)
	})

	t.Run("fail: duplicate event", func(t *testing.T) {
		t0 := New()
		err := t0.CreateEvent(context.Background(), t1)
		assert.NoError(t, err)

		err = t0.CreateEvent(context.Background(), t1)
		assert.Error(t, err)
		assert.ErrorIs(t, err, storage.ErrEventExists)
	})

	t.Run("fail: context cancellation", func(t *testing.T) {
		t0 := New()

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err := t0.CreateEvent(ctx, t1)
		assert.Error(t, err)
		assert.ErrorIs(t, err, context.Canceled)

		// Check what cancelling happened before adding
		_, exists := t0.events[t1.ID]
		assert.False(t, exists, "Event should NOT exist in storage")
	})
}

func TestStorage_UpdateEvent(t *testing.T) {
	t1 := &storage.Event{
		ID:    "id",
		Title: "title1",
	}

	t2 := &storage.Event{
		ID:    "id",
		Title: "title2",
	}

	t.Run("successful", func(t *testing.T) {
		t0 := New()
		err := t0.CreateEvent(context.Background(), t1)
		assert.NoError(t, err)
		assert.Equal(t, "id", t0.events[t1.ID].ID)
		assert.Equal(t, "title1", t0.events[t1.ID].Title)

		err = t0.UpdateEvent(context.Background(), t2)
		assert.NoError(t, err)
		assert.Equal(t, "id", t0.events[t1.ID].ID)
		assert.Equal(t, "title2", t0.events[t1.ID].Title)
	})

	t.Run("fail: no such event", func(t *testing.T) {
		t0 := New()
		err := t0.UpdateEvent(context.Background(), t1)
		assert.Error(t, err)
		assert.ErrorIs(t, err, storage.ErrEventNotFound)
	})

	t.Run("fail: context cancellation", func(t *testing.T) {
		t0 := New()
		err := t0.CreateEvent(context.Background(), t1)
		assert.NoError(t, err)
		assert.Equal(t, "title1", t0.events[t1.ID].Title)

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err = t0.UpdateEvent(ctx, t2)
		assert.ErrorIs(t, err, context.Canceled)

		savedEvent, exists := t0.events[t1.ID]
		assert.True(t, exists)
		assert.Equal(t, "id", savedEvent.ID)
		// Check what cancelling happened before adding
		assert.Equal(t, "title1", savedEvent.Title)
	})
}

func TestStorage_DeleteEvent(t *testing.T) {
	t1 := &storage.Event{
		ID:    "id",
		Title: "title1",
	}

	t.Run("successful", func(t *testing.T) {
		t0 := New()
		err := t0.CreateEvent(context.Background(), t1)
		assert.NoError(t, err)
		assert.Equal(t, "id", t0.events[t1.ID].ID)
		assert.Equal(t, "title1", t0.events[t1.ID].Title)

		err = t0.DeleteEvent(context.Background(), t1.ID)
		assert.NoError(t, err)

		_, exists := t0.events[t1.ID]
		assert.False(t, exists)
	})

	t.Run("fail: no such event", func(t *testing.T) {
		t0 := New()
		err := t0.DeleteEvent(context.Background(), t1.ID)
		assert.Error(t, err)
		assert.ErrorIs(t, err, storage.ErrEventNotFound)
	})

	t.Run("fail: context cancellation", func(t *testing.T) {
		t0 := New()
		err := t0.CreateEvent(context.Background(), t1)
		assert.NoError(t, err)
		assert.Equal(t, "title1", t0.events[t1.ID].Title)

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err = t0.DeleteEvent(ctx, t1.ID)
		assert.ErrorIs(t, err, context.Canceled)

		savedEvent, exists := t0.events[t1.ID]
		assert.True(t, exists)
		assert.Equal(t, "id", savedEvent.ID)
	})
}

func TestStorage_GetEvent(t *testing.T) {
	t1 := &storage.Event{
		ID: "id",
	}

	t.Run("successful", func(t *testing.T) {
		t0 := New()
		err := t0.CreateEvent(context.Background(), t1)
		assert.NoError(t, err)
		assert.Equal(t, "id", t0.events[t1.ID].ID)

		ret, err := t0.GetEvent(context.Background(), t1.ID)
		assert.NoError(t, err)
		assert.Equal(t, "id", ret.ID)
	})

	t.Run("fail: no such event", func(t *testing.T) {
		t0 := New()
		ret, err := t0.GetEvent(context.Background(), t1.ID)
		assert.Error(t, err)
		assert.ErrorIs(t, err, storage.ErrEventNotFound)
		assert.Nil(t, ret)
	})

	t.Run("fail: context cancellation", func(t *testing.T) {
		t0 := New()
		err := t0.CreateEvent(context.Background(), t1)
		assert.NoError(t, err)

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		ret, err := t0.GetEvent(ctx, t1.ID)
		assert.ErrorIs(t, err, context.Canceled)
		assert.Nil(t, ret)
	})
}

func TestStorage_ListEvents(t *testing.T) {
	t1 := &storage.Event{ID: "t1", Title: "Title_t1", UserID: "user1"}
	t2 := &storage.Event{ID: "t2", Title: "Title_t2", UserID: "user2"}
	t3 := &storage.Event{ID: "t3", Title: "Title_t3", UserID: "user1"}
	t4 := &storage.Event{ID: "t4", Title: "Title_t4", UserID: "user4"}
	t5 := &storage.Event{ID: "t5", Title: "Title_t5", UserID: "user1"}

	initStorage := func(ms *MemoryStorage) {
		ms.CreateEvent(context.Background(), t1)
		ms.CreateEvent(context.Background(), t2)
		ms.CreateEvent(context.Background(), t3)
		ms.CreateEvent(context.Background(), t4)
		ms.CreateEvent(context.Background(), t5)
	}

	t.Run("successful", func(t *testing.T) {
		t0 := New()
		initStorage(t0)

		ret, err := t0.ListEvents(context.Background(), "user1")
		assert.NoError(t, err)
		require.Equal(t, 3, len(ret))

		// check list by converting it to map
		found := make(map[string]bool)
		for _, event := range ret {
			found[event.ID] = true
		}

		assert.True(t, found["t1"], "Event with ID t1 not found")
		assert.True(t, found["t3"], "Event with ID t2 not found")
		assert.True(t, found["t5"], "Event with ID t3 not found")
	})

	t.Run("fail: context cancellation", func(t *testing.T) {
		t0 := New()
		initStorage(t0)

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// cancelling before colllecting list of elemets
		ret, err := t0.ListEvents(ctx, "user1")
		assert.ErrorIs(t, err, context.Canceled)
		assert.Equal(t, 0, len(ret))
	})
}
