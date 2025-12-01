package colordispanser

import (
	"sync"

	pkgErrors "github.com/pkg/errors"
)

type Dispanser struct {
	Map sync.Map
}

func NewDispanser() *Dispanser {
	values := []string{
		"primary",
		"secondary",
		"success",
		"danger",
		"warning",
		"dark",
	}

	c := &Dispanser{}
	for _, value := range values {
		c.Map.Store(value, false)
	}

	return c
}

func (c *Dispanser) Get() (string, error) {
	result := ""

	c.Map.Range(func(key, value any) bool {
		v := value.(bool)

		if !v {
			result = key.(string)
			c.Map.Store(result, true)
			return false
		}

		return true
	})

	if len(result) == 0 {
		return "", pkgErrors.New("color limit exceeded")
	}

	return result, nil
}

func (c *Dispanser) Reset(key string) {
	c.Map.Store(key, false)
}
