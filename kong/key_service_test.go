package kong

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestKeyService(T *testing.T) {
	T.Skip("Key without set fails until Kong/kong@302f2f7")
	RunWhenKong(T, ">=3.1.0")
	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	assert.NoError(err)
	assert.NotNil(client)

	//nolint:lll // JSON can't split lines so the encoded field is way too long for Go
	key := &Key{
		Name: String("foo"),
		KID:  String("foo-1"),
		JWK: String(`{
			"kty": "RSA",
			"kid": "foo-1",
			"n": "v2KAzzfruqctVHaE9WSCWIg1xAhMwxTIK-i56WNqPtpWBo9AqxcVea8NyVctEjUNq_mix5CklNy3ru7ARh7rBG_LU65fzs4fY_uYalul3QZSnr61Gj-cTUB3Gy4PhA63yXCbYRR3gDy6WR_wfis1MS61j0R_AjgXuVufmmC0F7R9qSWfR8ft0CbQgemEHY3ddKeW7T7fKv1jnRwYAkl5B_xtvxRFIYT-uR9NNftixNpUIW7q8qvOH7D9icXOg4_wIVxTRe5QiRYwEFoUbV1V9bFtu5FLal0vZnLaWwg5tA6enhzBpxJNdrS0v1RcPpyeNP-9r3cUDGmeftwz9v95UQ",
			"e": "AQAB",
			"alg": "A256GCM"
		}`),
	}

	createdKey, err := client.Keys.Create(defaultCtx, key)
	assert.NoError(err)
	require.NotNil(createdKey)

	key, err = client.Keys.Get(defaultCtx, createdKey.ID)
	assert.NoError(err)
	require.NotNil(key)

	key.Name = String("bar")
	key, err = client.Keys.Update(defaultCtx, key)
	assert.NoError(err)
	require.NotNil(key)
	assert.Equal("bar", *key.Name)

	err = client.Keys.Delete(defaultCtx, createdKey.ID)
	assert.NoError(err)

	// ID can be specified
	id := uuid.NewString()
	//nolint:lll // JSON can't split lines so the encoded field is way too long for Go
	key = &Key{
		Name: String("foo"),
		ID:   String(id),
		KID:  String("foo-2"),
		JWK: String(`{
			"kty": "RSA",
			"kid": "foo-2",
			"n": "v2KAzzfruqctVHaE9WSCWIg1xAhMwxTIK-i56WNqPtpWBo9AqxcVea8NyVctEjUNq_mix5CklNy3ru7ARh7rBG_LU65fzs4fY_uYalul3QZSnr61Gj-cTUB3Gy4PhA63yXCbYRR3gDy6WR_wfis1MS61j0R_AjgXuVufmmC0F7R9qSWfR8ft0CbQgemEHY3ddKeW7T7fKv1jnRwYAkl5B_xtvxRFIYT-uR9NNftixNpUIW7q8qvOH7D9icXOg4_wIVxTRe5QiRYwEFoUbV1V9bFtu5FLal0vZnLaWwg5tA6enhzBpxJNdrS0v1RcPpyeNP-9r3cUDGmeftwz9v95UQ",
			"e": "AQAB",
			"alg": "A256GCM"
		}`),
	}

	createdKey, err = client.Keys.Create(defaultCtx, key)
	assert.NoError(err)
	assert.NotNil(createdKey)
	assert.Equal(id, *createdKey.ID)

	err = client.Keys.Delete(defaultCtx, createdKey.ID)
	assert.NoError(err)
}

func TestKeyServiceWithSet(T *testing.T) {
	RunWhenKong(T, ">=3.1.0")
	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	assert.NoError(err)
	assert.NotNil(client)

	createdKeySet, err := client.KeySets.Create(defaultCtx, &KeySet{Name: String("foo")})
	assert.NoError(err)
	require.NotNil(createdKeySet)

	//nolint:lll // JSON can't split lines so the encoded field is way too long for Go
	key := &Key{
		Name: String("foo"),
		KID:  String("foo-1"),
		Set:  &KeySet{ID: createdKeySet.ID},
		JWK: String(`{
			"kty": "RSA",
			"kid": "foo-1",
			"n": "v2KAzzfruqctVHaE9WSCWIg1xAhMwxTIK-i56WNqPtpWBo9AqxcVea8NyVctEjUNq_mix5CklNy3ru7ARh7rBG_LU65fzs4fY_uYalul3QZSnr61Gj-cTUB3Gy4PhA63yXCbYRR3gDy6WR_wfis1MS61j0R_AjgXuVufmmC0F7R9qSWfR8ft0CbQgemEHY3ddKeW7T7fKv1jnRwYAkl5B_xtvxRFIYT-uR9NNftixNpUIW7q8qvOH7D9icXOg4_wIVxTRe5QiRYwEFoUbV1V9bFtu5FLal0vZnLaWwg5tA6enhzBpxJNdrS0v1RcPpyeNP-9r3cUDGmeftwz9v95UQ",
			"e": "AQAB",
			"alg": "A256GCM"
		}`),
	}

	createdKey, err := client.Keys.Create(defaultCtx, key)
	assert.NoError(err)
	require.NotNil(createdKey)

	key, err = client.Keys.Get(defaultCtx, createdKey.ID)
	assert.NoError(err)
	require.NotNil(key)

	key.Name = String("bar")
	key, err = client.Keys.Update(defaultCtx, key)
	assert.NoError(err)
	require.NotNil(key)
	assert.Equal("bar", *key.Name)

	err = client.Keys.Delete(defaultCtx, createdKey.ID)
	assert.NoError(err)

	// ID can be specified
	id := uuid.NewString()
	//nolint:lll // JSON can't split lines so the encoded field is way too long for Go
	key = &Key{
		Name: String("foo"),
		ID:   String(id),
		Set:  &KeySet{ID: createdKeySet.ID},
		KID:  String("foo-2"),
		JWK: String(`{
			"kty": "RSA",
			"kid": "foo-2",
			"n": "v2KAzzfruqctVHaE9WSCWIg1xAhMwxTIK-i56WNqPtpWBo9AqxcVea8NyVctEjUNq_mix5CklNy3ru7ARh7rBG_LU65fzs4fY_uYalul3QZSnr61Gj-cTUB3Gy4PhA63yXCbYRR3gDy6WR_wfis1MS61j0R_AjgXuVufmmC0F7R9qSWfR8ft0CbQgemEHY3ddKeW7T7fKv1jnRwYAkl5B_xtvxRFIYT-uR9NNftixNpUIW7q8qvOH7D9icXOg4_wIVxTRe5QiRYwEFoUbV1V9bFtu5FLal0vZnLaWwg5tA6enhzBpxJNdrS0v1RcPpyeNP-9r3cUDGmeftwz9v95UQ",
			"e": "AQAB",
			"alg": "A256GCM"
		}`),
	}

	createdKey, err = client.Keys.Create(defaultCtx, key)
	assert.NoError(err)
	assert.NotNil(createdKey)
	assert.Equal(id, *createdKey.ID)

	err = client.Keys.Delete(defaultCtx, createdKey.ID)
	assert.NoError(err)

	err = client.KeySets.Delete(defaultCtx, createdKeySet.ID)
	assert.NoError(err)
}

func TestKeyWithTags(T *testing.T) {
	T.Skip("Key without set fails until Kong/kong@302f2f7")
	RunWhenKong(T, ">=3.1.0")
	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	assert.NoError(err)
	assert.NotNil(client)

	//nolint:lll // JSON can't split lines so the encoded field is way too long for Go
	key := &Key{
		Name: String("foo"),
		KID:  String("foo-1"),
		JWK: String(`{
			"kty": "RSA",
			"kid": "foo-1",
			"n": "v2KAzzfruqctVHaE9WSCWIg1xAhMwxTIK-i56WNqPtpWBo9AqxcVea8NyVctEjUNq_mix5CklNy3ru7ARh7rBG_LU65fzs4fY_uYalul3QZSnr61Gj-cTUB3Gy4PhA63yXCbYRR3gDy6WR_wfis1MS61j0R_AjgXuVufmmC0F7R9qSWfR8ft0CbQgemEHY3ddKeW7T7fKv1jnRwYAkl5B_xtvxRFIYT-uR9NNftixNpUIW7q8qvOH7D9icXOg4_wIVxTRe5QiRYwEFoUbV1V9bFtu5FLal0vZnLaWwg5tA6enhzBpxJNdrS0v1RcPpyeNP-9r3cUDGmeftwz9v95UQ",
			"e": "AQAB",
			"alg": "A256GCM"
		}`),
		Tags: StringSlice("tag1", "tag2"),
	}

	createdKey, err := client.Keys.Create(defaultCtx, key)
	assert.NoError(err)
	require.NotNil(createdKey)
	assert.Equal(StringSlice("tag1", "tag2"), createdKey.Tags)

	err = client.Keys.Delete(defaultCtx, createdKey.ID)
	assert.NoError(err)
}

func TestKeyWithTagsWithSet(T *testing.T) {
	RunWhenKong(T, ">=3.1.0")
	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	assert.NoError(err)
	assert.NotNil(client)

	createdKeySet, err := client.KeySets.Create(defaultCtx, &KeySet{Name: String("foo")})
	assert.NoError(err)
	require.NotNil(createdKeySet)

	//nolint:lll // JSON can't split lines so the encoded field is way too long for Go
	key := &Key{
		Name: String("foo"),
		KID:  String("foo-1"),
		Set:  &KeySet{ID: createdKeySet.ID},
		JWK: String(`{
			"kty": "RSA",
			"kid": "foo-1",
			"n": "v2KAzzfruqctVHaE9WSCWIg1xAhMwxTIK-i56WNqPtpWBo9AqxcVea8NyVctEjUNq_mix5CklNy3ru7ARh7rBG_LU65fzs4fY_uYalul3QZSnr61Gj-cTUB3Gy4PhA63yXCbYRR3gDy6WR_wfis1MS61j0R_AjgXuVufmmC0F7R9qSWfR8ft0CbQgemEHY3ddKeW7T7fKv1jnRwYAkl5B_xtvxRFIYT-uR9NNftixNpUIW7q8qvOH7D9icXOg4_wIVxTRe5QiRYwEFoUbV1V9bFtu5FLal0vZnLaWwg5tA6enhzBpxJNdrS0v1RcPpyeNP-9r3cUDGmeftwz9v95UQ",
			"e": "AQAB",
			"alg": "A256GCM"
		}`),
		Tags: StringSlice("tag1", "tag2"),
	}

	createdKey, err := client.Keys.Create(defaultCtx, key)
	assert.NoError(err)
	require.NotNil(createdKey)
	assert.Equal(StringSlice("tag1", "tag2"), createdKey.Tags)

	err = client.Keys.Delete(defaultCtx, createdKey.ID)
	assert.NoError(err)

	err = client.KeySets.Delete(defaultCtx, createdKeySet.ID)
	assert.NoError(err)
}

func TestKeyListWithTags(T *testing.T) {
	T.Skip("Key without set fails until Kong/kong@302f2f7")
	RunWhenKong(T, ">=3.1.0")
	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	assert.NoError(err)
	assert.NotNil(client)

	// fixtures
	//nolint:lll // JSON can't split lines so the encoded field is way too long for Go
	keys := []*Key{
		{
			Name: String("user1"),
			KID:  String("user-key-1"),
			JWK: String(`{
			"kty": "RSA",
			"kid": "user-key-1",
			"n": "v2KAzzfruqctVHaE9WSCWIg1xAhMwxTIK-i56WNqPtpWBo9AqxcVea8NyVctEjUNq_mix5CklNy3ru7ARh7rBG_LU65fzs4fY_uYalul3QZSnr61Gj-cTUB3Gy4PhA63yXCbYRR3gDy6WR_wfis1MS61j0R_AjgXuVufmmC0F7R9qSWfR8ft0CbQgemEHY3ddKeW7T7fKv1jnRwYAkl5B_xtvxRFIYT-uR9NNftixNpUIW7q8qvOH7D9icXOg4_wIVxTRe5QiRYwEFoUbV1V9bFtu5FLal0vZnLaWwg5tA6enhzBpxJNdrS0v1RcPpyeNP-9r3cUDGmeftwz9v95UQ",
			"e": "AQAB",
			"alg": "A256GCM"
		}`),
			Tags: StringSlice("tag1", "tag2"),
		},
		{
			Name: String("user2"),
			KID:  String("user-key-2"),
			JWK: String(`{
			"kty": "RSA",
			"kid": "user-key-2",
			"n": "v2KAzzfruqctVHaE9WSCWIg1xAhMwxTIK-i56WNqPtpWBo9AqxcVea8NyVctEjUNq_mix5CklNy3ru7ARh7rBG_LU65fzs4fY_uYalul3QZSnr61Gj-cTUB3Gy4PhA63yXCbYRR3gDy6WR_wfis1MS61j0R_AjgXuVufmmC0F7R9qSWfR8ft0CbQgemEHY3ddKeW7T7fKv1jnRwYAkl5B_xtvxRFIYT-uR9NNftixNpUIW7q8qvOH7D9icXOg4_wIVxTRe5QiRYwEFoUbV1V9bFtu5FLal0vZnLaWwg5tA6enhzBpxJNdrS0v1RcPpyeNP-9r3cUDGmeftwz9v95UQ",
			"e": "AQAB",
			"alg": "A256GCM"
		}`),
			Tags: StringSlice("tag2", "tag3"),
		},
		{
			Name: String("user3"),
			KID:  String("user-key-3"),
			JWK: String(`{
			"kty": "RSA",
			"kid": "user-key-3",
			"n": "v2KAzzfruqctVHaE9WSCWIg1xAhMwxTIK-i56WNqPtpWBo9AqxcVea8NyVctEjUNq_mix5CklNy3ru7ARh7rBG_LU65fzs4fY_uYalul3QZSnr61Gj-cTUB3Gy4PhA63yXCbYRR3gDy6WR_wfis1MS61j0R_AjgXuVufmmC0F7R9qSWfR8ft0CbQgemEHY3ddKeW7T7fKv1jnRwYAkl5B_xtvxRFIYT-uR9NNftixNpUIW7q8qvOH7D9icXOg4_wIVxTRe5QiRYwEFoUbV1V9bFtu5FLal0vZnLaWwg5tA6enhzBpxJNdrS0v1RcPpyeNP-9r3cUDGmeftwz9v95UQ",
			"e": "AQAB",
			"alg": "A256GCM"
		}`),
			Tags: StringSlice("tag1", "tag3"),
		},
		{
			Name: String("user4"),
			KID:  String("user-key-4"),
			JWK: String(`{
			"kty": "RSA",
			"kid": "user-key-4",
			"n": "v2KAzzfruqctVHaE9WSCWIg1xAhMwxTIK-i56WNqPtpWBo9AqxcVea8NyVctEjUNq_mix5CklNy3ru7ARh7rBG_LU65fzs4fY_uYalul3QZSnr61Gj-cTUB3Gy4PhA63yXCbYRR3gDy6WR_wfis1MS61j0R_AjgXuVufmmC0F7R9qSWfR8ft0CbQgemEHY3ddKeW7T7fKv1jnRwYAkl5B_xtvxRFIYT-uR9NNftixNpUIW7q8qvOH7D9icXOg4_wIVxTRe5QiRYwEFoUbV1V9bFtu5FLal0vZnLaWwg5tA6enhzBpxJNdrS0v1RcPpyeNP-9r3cUDGmeftwz9v95UQ",
			"e": "AQAB",
			"alg": "A256GCM"
		}`),
			Tags: StringSlice("tag1", "tag2"),
		},
		{
			Name: String("user5"),
			KID:  String("user-key-5"),
			JWK: String(`{
			"kty": "RSA",
			"kid": "user-key-5",
			"n": "v2KAzzfruqctVHaE9WSCWIg1xAhMwxTIK-i56WNqPtpWBo9AqxcVea8NyVctEjUNq_mix5CklNy3ru7ARh7rBG_LU65fzs4fY_uYalul3QZSnr61Gj-cTUB3Gy4PhA63yXCbYRR3gDy6WR_wfis1MS61j0R_AjgXuVufmmC0F7R9qSWfR8ft0CbQgemEHY3ddKeW7T7fKv1jnRwYAkl5B_xtvxRFIYT-uR9NNftixNpUIW7q8qvOH7D9icXOg4_wIVxTRe5QiRYwEFoUbV1V9bFtu5FLal0vZnLaWwg5tA6enhzBpxJNdrS0v1RcPpyeNP-9r3cUDGmeftwz9v95UQ",
			"e": "AQAB",
			"alg": "A256GCM"
		}`),
			Tags: StringSlice("tag2", "tag3"),
		},
		{
			Name: String("user6"),
			KID:  String("user-key-6"),
			JWK: String(`{
			"kty": "RSA",
			"kid": "user-key-6",
			"n": "v2KAzzfruqctVHaE9WSCWIg1xAhMwxTIK-i56WNqPtpWBo9AqxcVea8NyVctEjUNq_mix5CklNy3ru7ARh7rBG_LU65fzs4fY_uYalul3QZSnr61Gj-cTUB3Gy4PhA63yXCbYRR3gDy6WR_wfis1MS61j0R_AjgXuVufmmC0F7R9qSWfR8ft0CbQgemEHY3ddKeW7T7fKv1jnRwYAkl5B_xtvxRFIYT-uR9NNftixNpUIW7q8qvOH7D9icXOg4_wIVxTRe5QiRYwEFoUbV1V9bFtu5FLal0vZnLaWwg5tA6enhzBpxJNdrS0v1RcPpyeNP-9r3cUDGmeftwz9v95UQ",
			"e": "AQAB",
			"alg": "A256GCM"
		}`),
			Tags: StringSlice("tag1", "tag3"),
		},
	}

	// create fixtures
	for i := 0; i < len(keys); i++ {
		key, err := client.Keys.Create(defaultCtx, keys[i])
		assert.NoError(err)
		require.NotNil(key)
		keys[i] = key
	}

	keysFromKong, next, err := client.Keys.List(defaultCtx, &ListOpt{
		Tags: StringSlice("tag1"),
	})
	assert.NoError(err)
	assert.Nil(next)
	assert.Equal(4, len(keysFromKong))

	keysFromKong, next, err = client.Keys.List(defaultCtx, &ListOpt{
		Tags: StringSlice("tag2"),
	})
	assert.NoError(err)
	assert.Nil(next)
	assert.Equal(4, len(keysFromKong))

	keysFromKong, next, err = client.Keys.List(defaultCtx, &ListOpt{
		Tags: StringSlice("tag1", "tag2"),
	})
	assert.NoError(err)
	assert.Nil(next)
	assert.Equal(6, len(keysFromKong))

	keysFromKong, next, err = client.Keys.List(defaultCtx, &ListOpt{
		Tags:         StringSlice("tag1", "tag2"),
		MatchAllTags: true,
	})
	assert.NoError(err)
	assert.Nil(next)
	assert.Equal(2, len(keysFromKong))

	keysFromKong, next, err = client.Keys.List(defaultCtx, &ListOpt{
		Tags: StringSlice("tag1", "tag2"),
		Size: 3,
	})
	assert.NoError(err)
	assert.NotNil(next)
	assert.Equal(3, len(keysFromKong))

	keysFromKong, next, err = client.Keys.List(defaultCtx, next)
	assert.NoError(err)
	assert.Nil(next)
	assert.Equal(3, len(keysFromKong))

	keysFromKong, next, err = client.Keys.List(defaultCtx, &ListOpt{
		Tags:         StringSlice("tag1", "tag2"),
		MatchAllTags: true,
		Size:         1,
	})
	assert.NoError(err)
	assert.NotNil(next)
	assert.Equal(1, len(keysFromKong))

	keysFromKong, next, err = client.Keys.List(defaultCtx, next)
	assert.NoError(err)
	assert.Nil(next)
	assert.Equal(1, len(keysFromKong))

	for i := 0; i < len(keys); i++ {
		assert.NoError(client.Keys.Delete(defaultCtx, keys[i].Name))
	}
}

func TestKeyListWithTagsWithSet(T *testing.T) {
	RunWhenKong(T, ">=3.1.0")
	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	assert.NoError(err)
	assert.NotNil(client)

	// fixtures
	createdKeySet, err := client.KeySets.Create(defaultCtx, &KeySet{Name: String("foo")})
	assert.NoError(err)
	require.NotNil(createdKeySet)

	//nolint:lll // JSON can't split lines so the encoded field is way too long for Go
	keys := []*Key{
		{
			Name: String("user1"),
			KID:  String("user-key-1"),
			Set:  &KeySet{ID: createdKeySet.ID},
			JWK: String(`{
			"kty": "RSA",
			"kid": "user-key-1",
			"n": "v2KAzzfruqctVHaE9WSCWIg1xAhMwxTIK-i56WNqPtpWBo9AqxcVea8NyVctEjUNq_mix5CklNy3ru7ARh7rBG_LU65fzs4fY_uYalul3QZSnr61Gj-cTUB3Gy4PhA63yXCbYRR3gDy6WR_wfis1MS61j0R_AjgXuVufmmC0F7R9qSWfR8ft0CbQgemEHY3ddKeW7T7fKv1jnRwYAkl5B_xtvxRFIYT-uR9NNftixNpUIW7q8qvOH7D9icXOg4_wIVxTRe5QiRYwEFoUbV1V9bFtu5FLal0vZnLaWwg5tA6enhzBpxJNdrS0v1RcPpyeNP-9r3cUDGmeftwz9v95UQ",
			"e": "AQAB",
			"alg": "A256GCM"
		}`),
			Tags: StringSlice("tag1", "tag2"),
		},
		{
			Name: String("user2"),
			KID:  String("user-key-2"),
			Set:  &KeySet{ID: createdKeySet.ID},
			JWK: String(`{
			"kty": "RSA",
			"kid": "user-key-2",
			"n": "v2KAzzfruqctVHaE9WSCWIg1xAhMwxTIK-i56WNqPtpWBo9AqxcVea8NyVctEjUNq_mix5CklNy3ru7ARh7rBG_LU65fzs4fY_uYalul3QZSnr61Gj-cTUB3Gy4PhA63yXCbYRR3gDy6WR_wfis1MS61j0R_AjgXuVufmmC0F7R9qSWfR8ft0CbQgemEHY3ddKeW7T7fKv1jnRwYAkl5B_xtvxRFIYT-uR9NNftixNpUIW7q8qvOH7D9icXOg4_wIVxTRe5QiRYwEFoUbV1V9bFtu5FLal0vZnLaWwg5tA6enhzBpxJNdrS0v1RcPpyeNP-9r3cUDGmeftwz9v95UQ",
			"e": "AQAB",
			"alg": "A256GCM"
		}`),
			Tags: StringSlice("tag2", "tag3"),
		},
		{
			Name: String("user3"),
			KID:  String("user-key-3"),
			Set:  &KeySet{ID: createdKeySet.ID},
			JWK: String(`{
			"kty": "RSA",
			"kid": "user-key-3",
			"n": "v2KAzzfruqctVHaE9WSCWIg1xAhMwxTIK-i56WNqPtpWBo9AqxcVea8NyVctEjUNq_mix5CklNy3ru7ARh7rBG_LU65fzs4fY_uYalul3QZSnr61Gj-cTUB3Gy4PhA63yXCbYRR3gDy6WR_wfis1MS61j0R_AjgXuVufmmC0F7R9qSWfR8ft0CbQgemEHY3ddKeW7T7fKv1jnRwYAkl5B_xtvxRFIYT-uR9NNftixNpUIW7q8qvOH7D9icXOg4_wIVxTRe5QiRYwEFoUbV1V9bFtu5FLal0vZnLaWwg5tA6enhzBpxJNdrS0v1RcPpyeNP-9r3cUDGmeftwz9v95UQ",
			"e": "AQAB",
			"alg": "A256GCM"
		}`),
			Tags: StringSlice("tag1", "tag3"),
		},
		{
			Name: String("user4"),
			KID:  String("user-key-4"),
			Set:  &KeySet{ID: createdKeySet.ID},
			JWK: String(`{
			"kty": "RSA",
			"kid": "user-key-4",
			"n": "v2KAzzfruqctVHaE9WSCWIg1xAhMwxTIK-i56WNqPtpWBo9AqxcVea8NyVctEjUNq_mix5CklNy3ru7ARh7rBG_LU65fzs4fY_uYalul3QZSnr61Gj-cTUB3Gy4PhA63yXCbYRR3gDy6WR_wfis1MS61j0R_AjgXuVufmmC0F7R9qSWfR8ft0CbQgemEHY3ddKeW7T7fKv1jnRwYAkl5B_xtvxRFIYT-uR9NNftixNpUIW7q8qvOH7D9icXOg4_wIVxTRe5QiRYwEFoUbV1V9bFtu5FLal0vZnLaWwg5tA6enhzBpxJNdrS0v1RcPpyeNP-9r3cUDGmeftwz9v95UQ",
			"e": "AQAB",
			"alg": "A256GCM"
		}`),
			Tags: StringSlice("tag1", "tag2"),
		},
		{
			Name: String("user5"),
			KID:  String("user-key-5"),
			Set:  &KeySet{ID: createdKeySet.ID},
			JWK: String(`{
			"kty": "RSA",
			"kid": "user-key-5",
			"n": "v2KAzzfruqctVHaE9WSCWIg1xAhMwxTIK-i56WNqPtpWBo9AqxcVea8NyVctEjUNq_mix5CklNy3ru7ARh7rBG_LU65fzs4fY_uYalul3QZSnr61Gj-cTUB3Gy4PhA63yXCbYRR3gDy6WR_wfis1MS61j0R_AjgXuVufmmC0F7R9qSWfR8ft0CbQgemEHY3ddKeW7T7fKv1jnRwYAkl5B_xtvxRFIYT-uR9NNftixNpUIW7q8qvOH7D9icXOg4_wIVxTRe5QiRYwEFoUbV1V9bFtu5FLal0vZnLaWwg5tA6enhzBpxJNdrS0v1RcPpyeNP-9r3cUDGmeftwz9v95UQ",
			"e": "AQAB",
			"alg": "A256GCM"
		}`),
			Tags: StringSlice("tag2", "tag3"),
		},
		{
			Name: String("user6"),
			KID:  String("user-key-6"),
			Set:  &KeySet{ID: createdKeySet.ID},
			JWK: String(`{
			"kty": "RSA",
			"kid": "user-key-6",
			"n": "v2KAzzfruqctVHaE9WSCWIg1xAhMwxTIK-i56WNqPtpWBo9AqxcVea8NyVctEjUNq_mix5CklNy3ru7ARh7rBG_LU65fzs4fY_uYalul3QZSnr61Gj-cTUB3Gy4PhA63yXCbYRR3gDy6WR_wfis1MS61j0R_AjgXuVufmmC0F7R9qSWfR8ft0CbQgemEHY3ddKeW7T7fKv1jnRwYAkl5B_xtvxRFIYT-uR9NNftixNpUIW7q8qvOH7D9icXOg4_wIVxTRe5QiRYwEFoUbV1V9bFtu5FLal0vZnLaWwg5tA6enhzBpxJNdrS0v1RcPpyeNP-9r3cUDGmeftwz9v95UQ",
			"e": "AQAB",
			"alg": "A256GCM"
		}`),
			Tags: StringSlice("tag1", "tag3"),
		},
	}

	// create fixtures
	for i := 0; i < len(keys); i++ {
		key, err := client.Keys.Create(defaultCtx, keys[i])
		assert.NoError(err)
		require.NotNil(key)
		keys[i] = key
	}

	keysFromKong, next, err := client.Keys.List(defaultCtx, &ListOpt{
		Tags: StringSlice("tag1"),
	})
	assert.NoError(err)
	assert.Nil(next)
	assert.Equal(4, len(keysFromKong))

	keysFromKong, next, err = client.Keys.List(defaultCtx, &ListOpt{
		Tags: StringSlice("tag2"),
	})
	assert.NoError(err)
	assert.Nil(next)
	assert.Equal(4, len(keysFromKong))

	keysFromKong, next, err = client.Keys.List(defaultCtx, &ListOpt{
		Tags: StringSlice("tag1", "tag2"),
	})
	assert.NoError(err)
	assert.Nil(next)
	assert.Equal(6, len(keysFromKong))

	keysFromKong, next, err = client.Keys.List(defaultCtx, &ListOpt{
		Tags:         StringSlice("tag1", "tag2"),
		MatchAllTags: true,
	})
	assert.NoError(err)
	assert.Nil(next)
	assert.Equal(2, len(keysFromKong))

	keysFromKong, next, err = client.Keys.List(defaultCtx, &ListOpt{
		Tags: StringSlice("tag1", "tag2"),
		Size: 3,
	})
	assert.NoError(err)
	assert.NotNil(next)
	assert.Equal(3, len(keysFromKong))

	keysFromKong, next, err = client.Keys.List(defaultCtx, next)
	assert.NoError(err)
	assert.Nil(next)
	assert.Equal(3, len(keysFromKong))

	keysFromKong, next, err = client.Keys.List(defaultCtx, &ListOpt{
		Tags:         StringSlice("tag1", "tag2"),
		MatchAllTags: true,
		Size:         1,
	})
	assert.NoError(err)
	assert.NotNil(next)
	assert.Equal(1, len(keysFromKong))

	keysFromKong, next, err = client.Keys.List(defaultCtx, next)
	assert.NoError(err)
	assert.Nil(next)
	assert.Equal(1, len(keysFromKong))

	for i := 0; i < len(keys); i++ {
		assert.NoError(client.Keys.Delete(defaultCtx, keys[i].Name))
	}

	err = client.KeySets.Delete(defaultCtx, createdKeySet.ID)
	assert.NoError(err)
}
