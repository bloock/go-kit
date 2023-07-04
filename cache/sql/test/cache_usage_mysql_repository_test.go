package test

import (
	"context"
	"github.com/bloock/go-kit/cache/sql"
	"github.com/bloock/go-kit/domain"
	"github.com/bloock/go-kit/errors"
	"github.com/bloock/go-kit/observability"
	"github.com/bloock/go-kit/test_utils"
	"github.com/huandu/go-sqlbuilder"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestCacheUsageMysqlRepository(t *testing.T) {
	mysqlClient := test_utils.GetMysqlClient()
	cr := sql.NewPostgresCacheUsageRepository(mysqlClient.DB(), time.Second*30, observability.Logger{}, "test")
	ct := test_utils.NewPostgresCrudRepository(mysqlClient, sql.SqlCacheUsageTable, sqlbuilder.NewStruct(new(sql.SqlCacheUsage)))

	key := "core:37e1a574-d76e-47ef-8960-dcc970e5a893:limit"
	value := -1
	cacheUsage := domain.NewCacheUsage(key, value)

	t.Run("Given a existent key value should be returned", func(t *testing.T) {
		err := cr.Save(context.Background(), cacheUsage, "test")
		assert.NoError(t, err)

		res, err := cr.GetValueByKey(context.Background(), key, "test")
		assert.NoError(t, err)
		assert.Equal(t, cacheUsage, res)

		err = ct.Truncate()
		require.NoError(t, err)
	})

	t.Run("Given a key that not exist, should return error", func(t *testing.T) {
		_, err := cr.GetValueByKey(context.Background(), "non_existent_key", "test")
		assert.Error(t, err)
		assert.Equal(t, errors.ErrNotFound, err)

		err = ct.Truncate()
		require.NoError(t, err)
	})

	t.Run("Given a key that not exist when finding key, should return no error", func(t *testing.T) {
		_, err := cr.FindValueByKey(context.Background(), "non_existent_key", "test")
		assert.NoError(t, err)

		err = ct.Truncate()
		require.NoError(t, err)
	})

	t.Run("Updating keys, should work", func(t *testing.T) {
		newValue := 0
		updateCacheUsage := domain.NewCacheUsage(key, newValue)

		err := cr.Save(context.Background(), cacheUsage, "test")
		require.NoError(t, err)

		res, err := cr.GetValueByKey(context.Background(), key, "test")
		assert.NoError(t, err)
		assert.Equal(t, cacheUsage, res)

		err = cr.Update(context.Background(), updateCacheUsage, "test")
		assert.NoError(t, err)

		res, err = cr.GetValueByKey(context.Background(), key, "test")
		assert.NoError(t, err)
		assert.Equal(t, updateCacheUsage, res)

		err = ct.Truncate()
		require.NoError(t, err)
	})
}
