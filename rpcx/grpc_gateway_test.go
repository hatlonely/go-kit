package rpcx

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/pkg/errors"
	. "github.com/smartystreets/goconvey/convey"
	"google.golang.org/grpc/codes"

	"github.com/hatlonely/go-kit/examples/rpcx/api/gen/go/api"
)

type ExampleService struct {
	api.ExampleServiceServer
}

func (s *ExampleService) Echo(ctx context.Context, req *api.EchoReq) (*api.EchoRes, error) {
	return &api.EchoRes{
		Message: req.Message,
	}, nil
}

func (s *ExampleService) Add(ctx context.Context, req *api.AddReq) (*api.AddRes, error) {
	if req.I1 < 0 {
		return nil, NewError(errors.New("i1 should be positive"), codes.InvalidArgument, "InvalidArgument", "i1 should be positive")
	}
	if req.I2 < 0 {
		return nil, errors.Wrap(NewError(errors.New("i2 should be positive"), codes.InvalidArgument, "BadRequest", "i2 should be positive"), "wrap error")
	}
	if req.I1 == 101 {
		return nil, errors.New("hello world")
	}
	if req.I1 == 102 {
		panic("panic")
	}
	if req.I1 == 103 {
		return nil, errors.Wrap(NewError(
			errors.New("hello world"), codes.Unknown, "PermissionDenied", "permissionDeny").
			SetStatus(http.StatusUnauthorized).
			SetBody(`{"errCode": "hello", "message": "world"}`), "wrap error")
	}

	return &api.AddRes{
		Val: req.I1 + req.I2,
	}, nil
}

func waitPortOpen(port int) {
	for {
		timeout := 50 * time.Millisecond
		conn, err := net.DialTimeout("tcp", net.JoinHostPort("127.0.0.1", strconv.Itoa(port)), timeout)
		if err != nil {
			continue
		}
		if conn != nil {
			_ = conn.Close()
			return
		}
	}
}

func waitPortClose(port int) {
	for {
		timeout := 50 * time.Millisecond
		conn, err := net.DialTimeout("tcp", net.JoinHostPort("127.0.0.1", strconv.Itoa(port)), timeout)
		if err != nil {
			return
		}
		if conn != nil {
			_ = conn.Close()
		}
	}
}

func ensureServiceUp() {
	waitPortOpen(6080)
	waitPortOpen(80)
	client := NewHttpClient()
	var res api.AddRes
	resMeta := map[string]string{}
	for i := 0; i < 10; i++ {
		err := client.Get(
			"http://127.0.0.1/v1/echo",
			map[string]string{"message": "hello world"},
			map[string]string{"x-request-id": "test-request-id", "x-user-id": "121231", "origin": "http://abc.example.com"},
			nil,
			&resMeta,
			&res,
		)
		fmt.Println(err, res)
		if err == nil {
			break
		}
		time.Sleep(1000 * time.Millisecond)
	}
}

func TestGrpcGateway_AddHttpPostHandler(t *testing.T) {
	Convey("TestGrpcGateway_AddHttpPostHandler", t, func() {
		server, err := NewGrpcGatewayWithOptions(&GrpcGatewayOptions{
			HttpPort:         80,
			GrpcPort:         6080,
			EnableTrace:      false,
			EnableMetric:     false,
			EnablePprof:      false,
			ExitTimeout:      10 * time.Second,
			Validators:       []string{"Default"},
			RequestIDMetaKey: "x-request-id",
			Headers:          []string{"X-Request-Id", "X-User-Id"},
		})
		So(err, ShouldBeNil)

		server.AddHttpPostHandler(func(w *BufferedHttpResponseWriter, r *http.Request) (*BufferedHttpResponseWriter, error) {
			fmt.Println(w.status)
			fmt.Println(w.body.String())
			fmt.Println(w.header)
			var v map[string]interface{}
			if err := json.Unmarshal(w.Body(), &v); err != nil {
				return w, err
			}

			v["httpStatusCode"] = w.Status()

			nw := NewBufferedHttpResponseWriter()
			nw.SetHeader(w.header)
			nw.WriteHeader(http.StatusOK)
			buf, _ := json.Marshal(v)
			_, _ = nw.Write(buf)
			return nw, nil
		})

		api.RegisterExampleServiceServer(server.GRPCServer(), &ExampleService{})
		So(server.RegisterServiceHandlerFunc(api.RegisterExampleServiceHandlerFromEndpoint), ShouldBeNil)

		go server.Run()
		defer func() {
			server.Stop()
			waitPortClose(80)
			waitPortClose(6080)
		}()

		ensureServiceUp()
		client := NewHttpClient()

		{
			var res map[string]interface{}
			resMeta := map[string]string{}
			err := client.Post(
				"http://127.0.0.1/v1/add",
				nil,
				map[string]string{"x-request-id": "test-request-id", "x-user-id": "121231"},
				&api.AddReq{I1: 12, I2: 34},
				&resMeta,
				&res,
			)
			So(err, ShouldBeNil)
			fmt.Println(res)
			So(res["val"], ShouldEqual, 46)
			So(res["httpStatusCode"], ShouldEqual, 200)
		}
	})
}

func TestGrpcGateway_SetErrorDetailMarshaler(t *testing.T) {
	Convey("TestGrpcGateway_SetErrorDetailMarshaler", t, func() {
		server, err := NewGrpcGatewayWithOptions(&GrpcGatewayOptions{
			HttpPort:         80,
			GrpcPort:         6080,
			EnableTrace:      false,
			EnableMetric:     false,
			EnablePprof:      false,
			ExitTimeout:      10 * time.Second,
			Validators:       []string{"Default"},
			RequestIDMetaKey: "x-request-id",
			Headers:          []string{"X-Request-Id", "X-User-Id"},
		})
		So(err, ShouldBeNil)

		server.SetErrorDetailMarshaler(func(detail *ErrorDetail) []byte {
			buf, _ := json.Marshal(&struct {
				RequestId      string
				HttpStatusCode int32  `json:"httpStatusCode"`
				DynamicCode    string `json:"dynamicCode"`
				DynamicMessage string `json:"dynamicMessage"`
				Success        bool   `json:"success"`
				Code           string `json:"code"`
				Message        string `json:"message"`
			}{
				RequestId:      detail.RequestID,
				HttpStatusCode: detail.Status,
				DynamicCode:    detail.Code,
				DynamicMessage: detail.Message,
				Code:           detail.Code,
				Message:        detail.Message,
				Success:        false,
			})

			return buf
		})

		api.RegisterExampleServiceServer(server.GRPCServer(), &ExampleService{})
		So(server.RegisterServiceHandlerFunc(api.RegisterExampleServiceHandlerFromEndpoint), ShouldBeNil)

		go server.Run()
		defer func() {
			server.Stop()
			waitPortClose(80)
			waitPortClose(6080)
		}()
		ensureServiceUp()
		client := NewHttpClient()

		{
			var res api.AddRes
			resMeta := map[string]string{}
			err := client.Post(
				"http://127.0.0.1/v1/add",
				nil,
				map[string]string{"x-request-id": "test-request-id", "x-user-id": "121231"},
				&api.AddReq{I1: -12, I2: 34},
				&resMeta,
				&res,
			)
			e := err.(*HttpError)
			fmt.Println(e)
			So(e.RequestID, ShouldEqual, "test-request-id")
		}
	})
}

func TestGrpcGateway_Cors(t *testing.T) {
	Convey("TestGrpcGateway_AddHttpHandler", t, func() {
		Convey("allow all", func() {
			server, err := NewGrpcGatewayWithOptions(&GrpcGatewayOptions{
				HttpPort:         80,
				GrpcPort:         6080,
				EnableTrace:      false,
				EnableMetric:     false,
				EnablePprof:      false,
				ExitTimeout:      10 * time.Second,
				Validators:       []string{"Default"},
				RequestIDMetaKey: "x-request-id",
				Headers:          []string{"X-Request-Id", "X-User-Id"},
				EnableCors:       true,
				Cors: CORSOptions{
					AllowAll:    true,
					AllowMethod: []string{"GET", "POST"},
					AllowHeader: []string{"X-Request-Id", "X-User-Id"},
				},
			})
			So(err, ShouldBeNil)

			api.RegisterExampleServiceServer(server.GRPCServer(), &ExampleService{})
			So(server.RegisterServiceHandlerFunc(api.RegisterExampleServiceHandlerFromEndpoint), ShouldBeNil)

			go server.Run()
			defer func() {
				server.Stop()
				waitPortClose(80)
				waitPortClose(6080)
			}()

			ensureServiceUp()

			client := NewHttpClient()

			{
				var res api.AddRes
				resMeta := map[string]string{}
				So(client.Post(
					"http://127.0.0.1/v1/add",
					nil,
					map[string]string{"x-request-id": "test-request-id", "x-user-id": "121231", "origin": "http://abc.example.com"},
					&api.AddReq{I1: 12, I2: 34},
					&resMeta,
					&res,
				), ShouldBeNil)
				fmt.Println(resMeta)
				So(resMeta["X-Request-Id"], ShouldResemble, "test-request-id")
				So(resMeta["X-User-Id"], ShouldResemble, "121231")
				So(resMeta["Access-Control-Allow-Origin"], ShouldEqual, "*")
				So(res.Val, ShouldEqual, 46)
			}

			{
				resMeta := map[string]string{}
				So(client.Do("OPTIONS", "http://127.0.0.1/", nil,
					map[string]string{"x-request-id": "test-request-id", "x-user-id": "121231", "origin": "http://abc.example.com"},
					nil, &resMeta, nil), ShouldBeNil)
				fmt.Println(resMeta)
				So(resMeta["Access-Control-Allow-Headers"], ShouldEqual, "X-Request-Id,X-User-Id")
				So(resMeta["Access-Control-Allow-Methods"], ShouldEqual, "GET,POST")
				So(resMeta["Access-Control-Allow-Origin"], ShouldEqual, "*")
			}
		})

		Convey("allow origin", func() {
			server, err := NewGrpcGatewayWithOptions(&GrpcGatewayOptions{
				HttpPort:         80,
				GrpcPort:         6080,
				EnableTrace:      false,
				EnableMetric:     false,
				EnablePprof:      false,
				ExitTimeout:      10 * time.Second,
				Validators:       []string{"Default"},
				RequestIDMetaKey: "x-request-id",
				Headers:          []string{"X-Request-Id", "X-User-Id"},
				EnableCors:       true,
				Cors: CORSOptions{
					AllowAll:    false,
					AllowOrigin: []string{"http://abc.example.com", "http://def.example.com"},
				},
			})
			So(err, ShouldBeNil)

			api.RegisterExampleServiceServer(server.GRPCServer(), &ExampleService{})
			So(server.RegisterServiceHandlerFunc(api.RegisterExampleServiceHandlerFromEndpoint), ShouldBeNil)

			go server.Run()
			defer func() {
				server.Stop()
				waitPortClose(80)
				waitPortClose(6080)
			}()

			ensureServiceUp()

			client := NewHttpClient()

			{
				var res api.AddRes
				resMeta := map[string]string{}
				So(client.Post(
					"http://127.0.0.1/v1/add",
					nil,
					map[string]string{"x-request-id": "test-request-id", "x-user-id": "121231", "origin": "http://abc.example.com"},
					&api.AddReq{I1: 12, I2: 34},
					&resMeta,
					&res,
				), ShouldBeNil)
				fmt.Println(resMeta)
				So(resMeta["X-Request-Id"], ShouldResemble, "test-request-id")
				So(resMeta["X-User-Id"], ShouldResemble, "121231")
				So(resMeta["Access-Control-Allow-Origin"], ShouldEqual, "http://abc.example.com")
				So(res.Val, ShouldEqual, 46)
			}

			{
				var res api.AddRes
				resMeta := map[string]string{}
				So(client.Post(
					"http://127.0.0.1/v1/add",
					nil,
					map[string]string{"x-request-id": "test-request-id", "x-user-id": "121231", "origin": "http://def.example.com"},
					&api.AddReq{I1: 12, I2: 34},
					&resMeta,
					&res,
				), ShouldBeNil)
				fmt.Println(resMeta)
				So(resMeta["X-Request-Id"], ShouldResemble, "test-request-id")
				So(resMeta["X-User-Id"], ShouldResemble, "121231")
				So(resMeta["Access-Control-Allow-Origin"], ShouldEqual, "http://def.example.com")
				So(res.Val, ShouldEqual, 46)
			}

			{
				var res api.AddRes
				resMeta := map[string]string{}
				err := client.Post(
					"http://127.0.0.1/v1/add",
					nil,
					map[string]string{"x-request-id": "test-request-id", "x-user-id": "121231", "origin": "http://ghi.example.com"},
					&api.AddReq{I1: 12, I2: 34},
					&resMeta,
					&res,
				)
				fmt.Println(err)
				e := err.(*HttpError)
				So(e.Status, ShouldEqual, 403)
			}
		})

		Convey("allow regex", func() {
			server, err := NewGrpcGatewayWithOptions(&GrpcGatewayOptions{
				HttpPort:         80,
				GrpcPort:         6080,
				EnableTrace:      false,
				EnableMetric:     false,
				EnablePprof:      false,
				ExitTimeout:      10 * time.Second,
				Validators:       []string{"Default"},
				RequestIDMetaKey: "x-request-id",
				Headers:          []string{"X-Request-Id", "X-User-Id"},
				EnableCors:       true,
				Cors: CORSOptions{
					AllowAll:   false,
					AllowRegex: []string{`.*\.example\..*`},
				},
			})
			So(err, ShouldBeNil)

			api.RegisterExampleServiceServer(server.GRPCServer(), &ExampleService{})
			So(server.RegisterServiceHandlerFunc(api.RegisterExampleServiceHandlerFromEndpoint), ShouldBeNil)

			go server.Run()
			defer func() {
				server.Stop()
				waitPortClose(80)
				waitPortClose(6080)
			}()

			ensureServiceUp()

			client := NewHttpClient()

			{
				var res api.AddRes
				resMeta := map[string]string{}
				So(client.Post(
					"http://127.0.0.1/v1/add",
					nil,
					map[string]string{"x-request-id": "test-request-id", "x-user-id": "121231", "origin": "http://abc.example.com"},
					&api.AddReq{I1: 12, I2: 34},
					&resMeta,
					&res,
				), ShouldBeNil)
				fmt.Println(resMeta)
				So(resMeta["X-Request-Id"], ShouldResemble, "test-request-id")
				So(resMeta["X-User-Id"], ShouldResemble, "121231")
				So(resMeta["Access-Control-Allow-Origin"], ShouldEqual, "http://abc.example.com")
				So(res.Val, ShouldEqual, 46)
			}

			{
				var res api.AddRes
				resMeta := map[string]string{}
				err := client.Post(
					"http://127.0.0.1/v1/add",
					nil,
					map[string]string{"x-request-id": "test-request-id", "x-user-id": "121231", "origin": "http://abc.notallow.com"},
					&api.AddReq{I1: 12, I2: 34},
					&resMeta,
					&res,
				)
				fmt.Println(err)
				e := err.(*HttpError)
				So(e.Status, ShouldEqual, 403)
			}
		})
	})
}

func TestGrpcGateway_AddHttpHandler(t *testing.T) {
	Convey("TestGrpcGateway_AddHttpHandler", t, func() {
		server, err := NewGrpcGatewayWithOptions(&GrpcGatewayOptions{
			HttpPort:         80,
			GrpcPort:         6080,
			EnableTrace:      false,
			EnableMetric:     false,
			EnablePprof:      false,
			ExitTimeout:      10 * time.Second,
			Validators:       []string{"Default"},
			RequestIDMetaKey: "x-request-id",
			Headers:          []string{"X-Request-Id", "X-User-Id"},
		})
		So(err, ShouldBeNil)

		api.RegisterExampleServiceServer(server.GRPCServer(), &ExampleService{})
		So(server.RegisterServiceHandlerFunc(api.RegisterExampleServiceHandlerFromEndpoint), ShouldBeNil)

		server.AddHttpHandler("/hello", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte(`{"key1": "val1"}`))
		}))

		So(server.HandleHttp("GET", "/hello/{name}", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
			_, _ = w.Write([]byte(fmt.Sprintf(`{"message": "hello %s"}`, pathParams["name"])))
		}), ShouldBeNil)

		go server.Run()
		defer func() {
			server.Stop()
			waitPortClose(80)
			waitPortClose(6080)
		}()

		ensureServiceUp()

		client := NewHttpClient()

		// test addHandler
		{
			var res interface{}
			resMeta := map[string]string{}
			So(client.Get(
				"http://127.0.0.1/hello",
				nil,
				nil,
				nil,
				&resMeta,
				&res,
			), ShouldBeNil)
			fmt.Println(resMeta)
			fmt.Println(res)
			So(res, ShouldResemble, map[string]interface{}{
				"key1": "val1",
			})
		}

		// test handlePath
		{
			var res interface{}
			resMeta := map[string]string{}
			So(client.Get(
				"http://127.0.0.1/hello/world",
				nil,
				nil,
				nil,
				&resMeta,
				&res,
			), ShouldBeNil)
			fmt.Println(resMeta)
			fmt.Println(res)
			So(res, ShouldResemble, map[string]interface{}{
				"message": "hello world",
			})
		}
	})
}

func TestGrpcGateway_AddHttpPreHandler(t *testing.T) {
	Convey("TestGrpcGateway_AddHttpPreHandler", t, func() {
		server, err := NewGrpcGatewayWithOptions(&GrpcGatewayOptions{
			HttpPort:         80,
			GrpcPort:         6080,
			EnableTrace:      false,
			EnableMetric:     false,
			EnablePprof:      false,
			ExitTimeout:      10 * time.Second,
			Validators:       []string{"Default"},
			RequestIDMetaKey: "x-request-id",
			Headers:          []string{"X-Request-Id", "X-User-Id"},
		})
		So(err, ShouldBeNil)

		api.RegisterExampleServiceServer(server.GRPCServer(), &ExampleService{})
		So(server.RegisterServiceHandlerFunc(api.RegisterExampleServiceHandlerFromEndpoint), ShouldBeNil)

		server.AddHttpPreHandler(func(w http.ResponseWriter, r *http.Request) error {
			if r.Header.Get("x-user-id") == "" {
				return NewError(errors.New("permission deny"), codes.PermissionDenied, "Forbidden", "permission deny")
			}
			return nil
		})

		go server.Run()
		defer func() {
			server.Stop()
			waitPortClose(80)
			waitPortClose(6080)
		}()

		ensureServiceUp()

		client := NewHttpClient()

		// pass
		{
			var res api.EchoRes
			resMeta := map[string]string{}
			So(client.Get(
				"http://127.0.0.1/v1/echo",
				map[string]string{"message": "hello world"},
				map[string]string{"x-request-id": "test-request-id", "x-user-id": "121231"},
				nil,
				&resMeta,
				&res,
			), ShouldBeNil)
			fmt.Println(resMeta)
			So(resMeta["X-Request-Id"], ShouldResemble, "test-request-id")
			So(resMeta["X-User-Id"], ShouldResemble, "121231")
			So(res.Message, ShouldEqual, "hello world")
		}

		// permission deny
		{
			var res api.AddRes
			resMeta := map[string]string{}
			err := client.Post(
				"http://127.0.0.1/v1/add",
				nil,
				map[string]string{"x-request-id": "test-request-id"},
				&api.AddReq{I1: -12, I2: 34},
				&resMeta,
				&res,
			)
			fmt.Println(err)
			e := err.(*HttpError)
			So(e.RequestID, ShouldEqual, "test-request-id")
			So(e.Status, ShouldEqual, 403)
			So(e.Code, ShouldEqual, "Forbidden")
			So(e.Message, ShouldEqual, "permission deny")
		}
	})
}

func TestGrpcGateway_AddGrpcPreHandler(t *testing.T) {
	Convey("TestGrpcGateway_AddGrpcPreHandler", t, func() {
		server, err := NewGrpcGatewayWithOptions(&GrpcGatewayOptions{
			HttpPort:         80,
			GrpcPort:         6080,
			EnableTrace:      false,
			EnableMetric:     false,
			EnablePprof:      false,
			ExitTimeout:      10 * time.Second,
			Validators:       []string{"Default"},
			RequestIDMetaKey: "x-request-id",
			Headers:          []string{"X-Request-Id", "X-User-Id"},
		})
		So(err, ShouldBeNil)

		api.RegisterExampleServiceServer(server.GRPCServer(), &ExampleService{})
		So(server.RegisterServiceHandlerFunc(api.RegisterExampleServiceHandlerFromEndpoint), ShouldBeNil)

		server.AddGrpcPreHandler(func(ctx context.Context, method string, req interface{}) error {
			if MetaDataIncomingGet(ctx, "x-user-id") == "" {
				return NewError(errors.New("permission deny"), codes.PermissionDenied, "Forbidden", "permission deny")
			}
			return nil
		})

		go server.Run()
		defer func() {
			server.Stop()
			waitPortClose(80)
			waitPortClose(6080)
		}()

		ensureServiceUp()

		client := NewHttpClient()

		// pass
		{
			var res api.EchoRes
			resMeta := map[string]string{}
			So(client.Get(
				"http://127.0.0.1/v1/echo",
				map[string]string{"message": "hello world"},
				map[string]string{"x-request-id": "test-request-id", "x-user-id": "121231"},
				nil,
				&resMeta,
				&res,
			), ShouldBeNil)
			So(resMeta["X-Request-Id"], ShouldResemble, "test-request-id")
			So(resMeta["X-User-Id"], ShouldResemble, "121231")
			So(res.Message, ShouldEqual, "hello world")
		}

		// permission deny
		{
			var res api.AddRes
			resMeta := map[string]string{}
			err := client.Post(
				"http://127.0.0.1/v1/add",
				nil,
				map[string]string{"x-request-id": "test-request-id"},
				&api.AddReq{I1: -12, I2: 34},
				&resMeta,
				&res,
			)
			e := err.(*HttpError)
			So(e.RequestID, ShouldEqual, "test-request-id")
			So(e.Status, ShouldEqual, 403)
			So(e.Code, ShouldEqual, "Forbidden")
			So(e.Message, ShouldEqual, "permission deny")
		}
	})
}

func TestGrpcGateway(t *testing.T) {
	Convey("TestGrpcGateway", t, func() {
		server, err := NewGrpcGatewayWithOptions(&GrpcGatewayOptions{
			HttpPort:         80,
			GrpcPort:         6080,
			EnableTrace:      false,
			EnableMetric:     false,
			EnablePprof:      false,
			ExitTimeout:      10 * time.Second,
			Validators:       []string{"Default"},
			RequestIDMetaKey: "x-request-id",
			Headers:          []string{"X-Request-Id", "X-User-Id"},
		})
		So(err, ShouldBeNil)

		api.RegisterExampleServiceServer(server.GRPCServer(), &ExampleService{})
		So(server.RegisterServiceHandlerFunc(api.RegisterExampleServiceHandlerFromEndpoint), ShouldBeNil)

		go server.Run()
		defer func() {
			server.Stop()
			waitPortClose(80)
			waitPortClose(6080)
		}()

		ensureServiceUp()

		client := NewHttpClient()

		// test echo
		{
			var res api.EchoRes
			resMeta := map[string]string{}
			So(client.Get(
				"http://127.0.0.1/v1/echo",
				map[string]string{"message": "hello world"},
				map[string]string{"x-request-id": "test-request-id", "x-user-id": "121231"},
				nil,
				&resMeta,
				&res,
			), ShouldBeNil)
			So(resMeta["X-Request-Id"], ShouldResemble, "test-request-id")
			So(resMeta["X-User-Id"], ShouldResemble, "121231")
			So(res.Message, ShouldEqual, "hello world")
		}

		// test add
		{
			var res api.AddRes
			resMeta := map[string]string{}
			So(client.Post(
				"http://127.0.0.1/v1/add",
				nil,
				map[string]string{"x-request-id": "test-request-id", "x-user-id": "121231"},
				&api.AddReq{I1: 12, I2: 34},
				&resMeta,
				&res,
			), ShouldBeNil)
			So(resMeta["X-Request-Id"], ShouldResemble, "test-request-id")
			So(resMeta["X-User-Id"], ShouldResemble, "121231")
			So(res.Val, ShouldEqual, 46)
		}

		// test without request id
		{
			var res api.AddRes
			resMeta := map[string]string{}
			So(client.Post(
				"http://127.0.0.1/v1/add",
				nil,
				map[string]string{"x-user-id": "121231"},
				&api.AddReq{I1: 12, I2: 34},
				&resMeta,
				&res,
			), ShouldBeNil)
			fmt.Println(resMeta["X-Request-Id"])
			So(resMeta["X-Request-Id"], ShouldNotBeEmpty)
			So(resMeta["X-User-Id"], ShouldResemble, "121231")
			So(res.Val, ShouldEqual, 46)
		}

		// test error
		{
			var res api.AddRes
			resMeta := map[string]string{}
			err := client.Post(
				"http://127.0.0.1/v1/add",
				nil,
				map[string]string{"x-request-id": "test-request-id", "x-user-id": "121231"},
				&api.AddReq{I1: -12, I2: 34},
				&resMeta,
				&res,
			)
			e := err.(*HttpError)
			fmt.Println(e)
			So(e.RequestID, ShouldEqual, "test-request-id")
			So(e.Status, ShouldEqual, 400)
			So(e.Code, ShouldEqual, "InvalidArgument")
			So(e.Message, ShouldEqual, "i1 should be positive")
		}

		// test wrap error
		{
			var res api.AddRes
			resMeta := map[string]string{}
			err := client.Post(
				"http://127.0.0.1/v1/add",
				nil,
				map[string]string{"x-request-id": "test-request-id", "x-user-id": "121231"},
				&api.AddReq{I1: 12, I2: -34},
				&resMeta,
				&res,
			)
			e := err.(*HttpError)
			fmt.Println(e)
			So(e.RequestID, ShouldEqual, "test-request-id")
			So(e.Status, ShouldEqual, 400)
			So(e.Code, ShouldEqual, "BadRequest")
			So(e.Message, ShouldEqual, "i2 should be positive")
		}

		// test internal error
		{
			var res api.AddRes
			resMeta := map[string]string{}
			err := client.Post(
				"http://127.0.0.1/v1/add",
				nil,
				map[string]string{"x-request-id": "test-request-id", "x-user-id": "121231"},
				&api.AddReq{I1: 101, I2: 34},
				&resMeta,
				&res,
			)
			e := err.(*HttpError)
			fmt.Println(e)
			So(e.RequestID, ShouldEqual, "test-request-id")
			So(e.Status, ShouldEqual, 500)
			So(e.Code, ShouldEqual, "InternalError")
			So(e.Message, ShouldEqual, "unknown error")
		}

		// test panic
		{
			var res api.AddRes
			resMeta := map[string]string{}
			err := client.Post(
				"http://127.0.0.1/v1/add",
				nil,
				map[string]string{"x-request-id": "test-request-id", "x-user-id": "121231"},
				&api.AddReq{I1: 102, I2: 34},
				&resMeta,
				&res,
			)
			e := err.(*HttpError)
			fmt.Println(e)
			So(e.RequestID, ShouldEqual, "test-request-id")
			So(e.Status, ShouldEqual, 500)
			So(e.Code, ShouldEqual, "InternalError")
			So(e.Message, ShouldEqual, "unknown error")
		}

		// test not found
		{
			var res interface{}
			err := client.Get(
				"http://127.0.0.1/v1/found",
				map[string]string{"message": "hello world"},
				map[string]string{"x-request-id": "test-request-id", "x-user-id": "121231"},
				nil,
				nil,
				&res,
			)
			So(reflect.TypeOf(err), ShouldEqual, reflect.TypeOf(&HttpError{}))
			e := err.(*HttpError)
			fmt.Println(e)
			So(e.Status, ShouldEqual, 404)
			So(e.Code, ShouldEqual, http.StatusText(http.StatusNotFound))
			So(e.RequestID, ShouldEqual, "test-request-id")
			So(e.Message, ShouldEqual, http.StatusText(http.StatusNotFound))
		}

		// test not implement
		{
			var res api.EchoRes
			resMeta := map[string]string{}
			err := client.Post(
				"http://127.0.0.1/v1/echo",
				map[string]string{"message": "hello world"},
				map[string]string{"x-request-id": "test-request-id", "x-user-id": "121231"},
				nil,
				&resMeta,
				&res,
			)
			e := err.(*HttpError)
			fmt.Println(e)
			So(e.Status, ShouldEqual, 501)
			So(e.Code, ShouldEqual, http.StatusText(http.StatusNotImplemented))
			So(e.RequestID, ShouldEqual, "test-request-id")
			So(e.Message, ShouldEqual, "Method Not Allowed")
		}

		// test error.SetStatus and error.SetBody
		{
			var res interface{}
			resMeta := map[string]string{}
			err := client.Post(
				"http://127.0.0.1/v1/add",
				nil,
				map[string]string{"x-request-id": "test-request-id", "x-user-id": "121231"},
				&api.AddReq{I1: 103, I2: 34},
				&resMeta,
				&res,
			)

			fmt.Println("@@@", err, res, resMeta)
		}
	})
}
