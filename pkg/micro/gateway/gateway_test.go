package gateway

import (
	"github.com/gin-gonic/gin"
	"testing"
)

func Test_Wrapper(t *testing.T) {
	expectedValue := "test_value"
	engine := gin.New()
	engine.Use(func(g *gin.Context) {
		value := g.GetHeader("test_key")
		if value != expectedValue {
			t.Fatal("test value not equal expected")
		}
	})
	//cors(engine)
}
