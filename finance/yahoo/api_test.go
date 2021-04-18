package yahoo

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestClient_checkPerSecond(t *testing.T) {
	c := &Client{}

	err := c.checkPerSecond()
	require.NoError(t, err)
	require.Equal(t, 1, c.perSecondCount)

	err = c.checkPerSecond()
	require.NoError(t, err)
	require.Equal(t, 2, c.perSecondCount)

	err = c.checkPerSecond()
	require.NoError(t, err)
	require.Equal(t, 3, c.perSecondCount)

	err = c.checkPerSecond()
	require.NoError(t, err)
	require.Equal(t, 4, c.perSecondCount)

	err = c.checkPerSecond()
	require.NoError(t, err)
	require.Equal(t, 5, c.perSecondCount)

	err = c.checkPerSecond()
	require.EqualError(t, err, "5 requests per second limit exceeded")

	c.perSecondStart = c.perSecondStart.Add(-(perSecondDuration + time.Second))

	err = c.checkPerSecond()
	require.NoError(t, err)
}
