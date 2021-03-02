package gateway

//
//import (
//	"encoding/json"
//	"github.com/gin-gonic/gin"
//	"github.com/yametech/yamecloud/common"
//	"github.com/yametech/yamecloud/pkg/micro/gateway"
//	"github.com/yametech/yamecloud/pkg/uri"
//	"net/http"
//	"net/http/httptest"
//	"reflect"
//	"testing"
//	"time"
//)
//
//func performRequest(r http.Handler, req *http.Request) *httptest.ResponseRecorder {
//	w := httptest.NewRecorder()
//	r.ServeHTTP(w, req)
//	return w
//}
//
//func Test_IsNeedSkip(t *testing.T) {
//	expectedStatus := http.StatusOK
//	expectedBody := gin.H{"msg": "success"}
//	engine := gin.New()
//	auth := NewFakeAuthorization()
//	engine.Use(IsNeedSkip(auth))
//	engine.GET("/test", func(c *gin.Context) {
//		isSkip := c.GetBool(IsSkip)
//		if isSkip {
//			c.JSON(http.StatusOK, gin.H{"msg": "success"})
//		} else {
//			c.JSON(http.StatusOK, gin.H{"msg": "failed"})
//		}
//	})
//	req, _ := http.NewRequest("GET", "/test", nil)
//
//	responseRecorder := performRequest(engine, req)
//	if expectedStatus != responseRecorder.Code {
//		t.Fatal("expect not equal")
//	}
//	var getBody = gin.H{}
//	json.Unmarshal([]byte(responseRecorder.Body.String()), &getBody)
//	if !reflect.DeepEqual(getBody, expectedBody) {
//		t.Fatal("expect not equal")
//	}
//}
//
//func Test_ValidateToken(t *testing.T) {
//	expectedStatus := http.StatusOK
//	expectedBody := gin.H{"msg": "success"}
//	engine := gin.New()
//	auth := NewFakeAuthorization()
//	engine.Use(ValidateToken(auth))
//	engine.GET("/test", func(c *gin.Context) {
//		username := c.GetString(AuthorizationUserName)
//		if username == "admin" {
//			c.JSON(http.StatusOK, gin.H{"msg": "success"})
//		} else {
//			c.JSON(http.StatusOK, gin.H{"msg": "failed"})
//		}
//	})
//	req, _ := http.NewRequest("GET", "/test", nil)
//	username := "admin"
//	token := generateToken(username)
//	req.Header.Add(common.AuthorizationHeader, token)
//
//	responseRecorder := performRequest(engine, req)
//	if expectedStatus != responseRecorder.Code {
//		t.Fatal("expect not equal")
//	}
//	var getBody = gin.H{}
//	json.Unmarshal([]byte(responseRecorder.Body.String()), &getBody)
//	if !reflect.DeepEqual(getBody, expectedBody) {
//		t.Fatal("expect not equal")
//	}
//
//}
//
//func Test_IsAdmin(t *testing.T) {
//	expectedStatus := http.StatusOK
//	expectedBody := gin.H{"msg": "success"}
//	engine := gin.New()
//	auth := NewFakeAuthorization()
//	engine.Use(
//		func(c *gin.Context) {
//			c.Set(AuthorizationUserName, "admin")
//		},
//		IsAdmin(auth),
//	)
//	engine.GET("/test", func(c *gin.Context) {
//		userIdentification, _ := c.Get(UserIdentification)
//		if userIdentification == Admin {
//			c.JSON(http.StatusOK, gin.H{"msg": "success"})
//		} else {
//			c.JSON(http.StatusOK, gin.H{"msg": "failed"})
//		}
//	})
//	req, _ := http.NewRequest("GET", "/test", nil)
//	responseRecorder := performRequest(engine, req)
//	if expectedStatus != responseRecorder.Code {
//		t.Fatal("expect not equal")
//	}
//
//	var getBody = gin.H{}
//	json.Unmarshal([]byte(responseRecorder.Body.String()), &getBody)
//	if !reflect.DeepEqual(getBody, expectedBody) {
//		t.Fatal("expect not equal")
//	}
//}
//
//func Test_IsTenantOwner(t *testing.T) {
//	expectedStatus := http.StatusOK
//	expectedBody := gin.H{"msg": "success"}
//	engine := gin.New()
//	auth := NewFakeAuthorization()
//	engine.Use(
//		func(c *gin.Context) {
//			c.Set(AuthorizationUserName, "tenant")
//		},
//		IsTenantOwner(auth),
//	)
//	engine.GET("/test", func(c *gin.Context) {
//		userIdentification, _ := c.Get(UserIdentification)
//		if userIdentification == TenantOwner {
//			c.JSON(http.StatusOK, gin.H{"msg": "success"})
//		} else {
//			c.JSON(http.StatusOK, gin.H{"msg": "failed"})
//		}
//	})
//	req, _ := http.NewRequest("GET", "/test", nil)
//	responseRecorder := performRequest(engine, req)
//	if expectedStatus != responseRecorder.Code {
//		t.Fatal("expect not equal")
//	}
//
//	var getBody = gin.H{}
//	json.Unmarshal([]byte(responseRecorder.Body.String()), &getBody)
//	if !reflect.DeepEqual(getBody, expectedBody) {
//		t.Fatal("expect not equal")
//	}
//}
//
//func Test_IsDepartmentOwner(t *testing.T) {
//	expectedStatus := http.StatusOK
//	expectedBody := gin.H{"msg": "success"}
//	engine := gin.New()
//	auth := NewFakeAuthorization()
//	engine.Use(
//		func(c *gin.Context) {
//			c.Set(AuthorizationUserName, "departmentOwner")
//		},
//		IsDepartmentOwner(auth),
//	)
//	engine.GET("/test", func(c *gin.Context) {
//		userIdentification, _ := c.Get(UserIdentification)
//		if userIdentification == DepartmentOwner {
//			c.JSON(http.StatusOK, gin.H{"msg": "success"})
//		} else {
//			c.JSON(http.StatusOK, gin.H{"msg": "failed"})
//		}
//	})
//	req, _ := http.NewRequest("GET", "/test", nil)
//	responseRecorder := performRequest(engine, req)
//	if expectedStatus != responseRecorder.Code {
//		t.Fatal("expect not equal")
//	}
//
//	var getBody = gin.H{}
//	json.Unmarshal([]byte(responseRecorder.Body.String()), &getBody)
//	if !reflect.DeepEqual(getBody, expectedBody) {
//		t.Fatal("expect not equal")
//	}
//}
//
//func Test_CheckNamespace(t *testing.T) {
//	expectedStatus := http.StatusOK
//	expectedBody := gin.H{"msg": "success"}
//	engine := gin.New()
//	auth := NewFakeAuthorization()
//	engine.Use(
//		func(c *gin.Context) {
//			c.Set(AuthorizationUserName, "ordinaryUser")
//			//c.Set(UserIdentification, Admin)
//		},
//		CheckNamespace(auth),
//	)
//	engine.GET("/test/apis/test.io/v1/namespaces/test/tests/mytest", func(c *gin.Context) {
//		c.JSON(http.StatusOK, gin.H{"msg": "success"})
//	})
//	req, _ := http.NewRequest("GET", "/test/apis/test.io/v1/namespaces/test/tests/mytest", nil)
//	responseRecorder := performRequest(engine, req)
//	if expectedStatus != responseRecorder.Code {
//		t.Fatal("expect not equal")
//	}
//
//	var getBody = gin.H{}
//	json.Unmarshal([]byte(responseRecorder.Body.String()), &getBody)
//	if !reflect.DeepEqual(getBody, expectedBody) {
//		t.Fatal("expect not equal")
//	}
//}
//
//func Test_CheckPermission(t *testing.T) {
//	expectedStatus := http.StatusOK
//	expectedBody := gin.H{"msg": "success"}
//	engine := gin.New()
//	auth := NewFakeAuthorization()
//	engine.Use(
//		func(c *gin.Context) {
//			c.Set(AuthorizationUserName, "ordinaryUser")
//			//c.Set(UserIdentification, Admin)
//		},
//		CheckPermission(auth),
//	)
//	engine.GET("/test/apis/test.io/v1/namespaces/test/tests/mytest", func(c *gin.Context) {
//		c.JSON(http.StatusOK, gin.H{"msg": "success"})
//	})
//	req, _ := http.NewRequest("GET", "/test/apis/test.io/v1/namespaces/test/tests/mytest", nil)
//	responseRecorder := performRequest(engine, req)
//	if expectedStatus != responseRecorder.Code {
//		t.Fatal("expect not equal")
//	}
//
//	var getBody = gin.H{}
//	json.Unmarshal([]byte(responseRecorder.Body.String()), &getBody)
//	if !reflect.DeepEqual(getBody, expectedBody) {
//		t.Fatal("expect not equal")
//	}
//}
//
////func Test_CheckUser(t *testing.T) {
////	expectedStatus := http.StatusOK
////	expectedBody := gin.H{"msg": "success"}
////	engine := gin.New()
////	auth := NewFakeAuthorization()
////	engine.Use(
////		func(c *gin.Context) {
////			c.Set(AuthorizationUserName, "ordinaryUser")
////		},
////		CheckUser(auth),
////	)
////	engine.GET("/test", func(c *gin.Context) {
////		c.JSON(http.StatusOK, gin.H{"msg": "success"})
////	})
////	req, _ := http.NewRequest("GET", "/test", nil)
////	responseRecorder := performRequest(engine, req)
////	if expectedStatus != responseRecorder.Code {
////		t.Fatal("expect not equal")
////	}
////
////	var getBody = gin.H{}
////	json.Unmarshal([]byte(responseRecorder.Body.String()), &getBody)
////	if !reflect.DeepEqual(getBody, expectedBody) {
////		t.Fatal("expect not equal")
////	}
////}
//
//func Test_IsWithGranted(t *testing.T) {
//	expectedStatus := http.StatusOK
//	expectedBody := gin.H{"msg": "success"}
//	engine := gin.New()
//	auth := NewFakeAuthorization()
//	engine.Use(
//		func(c *gin.Context) {
//			c.Set(AuthorizationUserName, "userWithGranted")
//		},
//		IsWithGranted(auth),
//	)
//	engine.GET("/grantPath", func(c *gin.Context) {
//		c.JSON(http.StatusOK, gin.H{"msg": "success"})
//	})
//	req, _ := http.NewRequest("GET", "/grantPath", nil)
//	responseRecorder := performRequest(engine, req)
//	if expectedStatus != responseRecorder.Code {
//		t.Fatal("expect not equal")
//	}
//
//	var getBody = gin.H{}
//	json.Unmarshal([]byte(responseRecorder.Body.String()), &getBody)
//	if !reflect.DeepEqual(getBody, expectedBody) {
//		t.Fatal("expect not equal")
//	}
//}
//
//func generateToken(username string) string {
//	expireTime := time.Now().Add(time.Hour * 24).Unix()
//	token, _ := (&gateway.Token{}).Encode(common.MicroSaltUserHeader, username, expireTime)
//	return token
//}
//
//var _ IAuthorization = (*FakeAuthorization)(nil)
//
//type FakeAuthorization struct {
//}
//
//func NewFakeAuthorization() *FakeAuthorization {
//	return &FakeAuthorization{}
//}
//
//func (auth *FakeAuthorization) IsNeedSkip(method, path string) (bool, error) {
//	excludeMap := map[string]string{
//		"/test": "GET",
//	}
//	value := excludeMap[path]
//	if method == value {
//		return true, nil
//	}
//	return false, nil
//}
//
//func (auth *FakeAuthorization) ValidateToken(token string) (*gateway.CustomClaims, error) {
//	decodeToken, err := (&gateway.Token{}).Decode(token)
//	if err != nil {
//		return nil, err
//	}
//	return decodeToken, nil
//}
//
////check whether a user is an admin
//func (auth *FakeAuthorization) IsAdmin(userName string) (bool, error) {
//	adminList := []string{"admin"}
//	for _, _item := range adminList {
//		if userName == _item {
//			return true, nil
//		}
//	}
//	return false, nil
//}
//
////check whether a user is a tenant owner
//func (auth *FakeAuthorization) IsTenantOwner(userName string) (bool, error) {
//	tenantList := []string{"tenant"}
//	for _, _item := range tenantList {
//		if userName == _item {
//			return true, nil
//		}
//	}
//	return false, nil
//}
//
////check whether a user is a department owner
//func (auth *FakeAuthorization) IsDepartmentOwner(userName string) (bool, error) {
//	departmentOwnerList := []string{"departmentOwner"}
//	for _, _item := range departmentOwnerList {
//		if userName == _item {
//			return true, nil
//		}
//	}
//	return false, nil
//}
//
////check whether a user is with granted
//func (auth *FakeAuthorization) IsWithGranted(userName string) (bool, error) {
//	if userName == "userWithGranted" {
//		return true, nil
//	}
//	return false, nil
//}
//
////check whether a user has specified uri permission
//func (auth *FakeAuthorization) CheckPermission(userName string, uri *uri.URI) (bool, error) {
//	if userName == "ordinaryUser" {
//		return true, nil
//	}
//	return false, nil
//}
//
////check whether a user allow access specified namespace
//func (auth *FakeAuthorization) CheckNamespace(userName, namespace string, isAdmin, isTenantOwner, isDepartmentOwner bool) (bool, error) {
//	if namespace == "test" {
//		return true, nil
//	}
//	return false, nil
//}
//
//func Test_NilArray(t *testing.T) {
//	var testArray []string = nil
//	for _, item := range testArray {
//		t.Log(item)
//	}
//}
