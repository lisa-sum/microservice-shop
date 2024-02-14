package data

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io/ioutil"
	"net/http"
	v1 "shop/api/user/v1"
	"shop/internal/biz"
	"testing"
)

var conn *grpc.ClientConn

// var userClient v1.UserClient
var userClient v1.UserClient

func Init() {
	var err error
	conn, err = grpc.Dial("127.0.0.1:9000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	// conn, err = grpc.Dial("192.168.3.220:9000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic("grpc link err" + err.Error())
	}
	userClient = v1.NewUserClient(conn)
}

// Init 初始化 grpc 链接
func HTTPClient(path string) {
	url := "https://127.0.0.1:8001" + path

	// 发起 GET 请求
	response, err := http.Get(url)
	if err != nil {
		fmt.Printf("请求失败: %v", err)
		panic(err)
	}
	defer response.Body.Close()

	// 读取响应内容
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("读取响应失败: %v", err)
		panic(err)
	}

	// 打印响应结果
	fmt.Println(string(body))
}

func TestCreateUser(t *testing.T) {
	mockUser := biz.MockUserRepo{}

	mockUser.On("CreateUser", mock.Anything, &biz.User{
		Mobile:   fmt.Sprintf("1388888888%d", 2),
		Nickname: fmt.Sprintf("YWWW%d", 2),
		Password: "123",
	}).Return(&biz.User{
		Mobile:   fmt.Sprintf("1388888888%d", 2),
		Nickname: fmt.Sprintf("YWWW%d", 2),
	}, nil)

	rsp, err := userClient.CreateUser(context.Background(), &v1.CreateUserInfo{
		Mobile:   fmt.Sprintf("2388888888%d", 1),
		Nickname: fmt.Sprintf("YWWW%d", 2),
		Password: "123",
	})

	// rsp, err := mockUser.CreateUser(context.Background(), &biz.User{
	// 	Mobile:   fmt.Sprintf("1388888888%d", 1),
	// 	Nickname: fmt.Sprintf("YWWW%d", 1),
	// 	Password: "123",
	// })

	// 断言期望的返回值
	assert.NoError(t, err)
	assert.Equal(t, "13888888882", rsp.Mobile)
	assert.Equal(t, "YWWW2", rsp.Nickname)

	// 断言模拟对象的期望行为是否被满足
	mockUser.AssertExpectations(t)
}
func TestCreateUser2(t *testing.T) {
	Init()
	// mockUser := biz.MockUserRepo{}
	// mockUser.On("CreateUser", mock.Anything, &biz.User{
	// 	Mobile:   fmt.Sprintf("1388888888%d", 1),
	// 	Nickname: fmt.Sprintf("YWWW%d", 1),
	// 	Password: "123",
	// }).Return(&biz.User{
	// 	Mobile:   fmt.Sprintf("1388888888%d", 1),
	// 	Nickname: fmt.Sprintf("YWWW%d", 1),
	// }, nil)

	rsp, err := userClient.CreateUser(context.Background(), &v1.CreateUserInfo{
		Mobile:   fmt.Sprintf("1388888888%d", 2),
		Nickname: fmt.Sprintf("YWWW%d", 2),
		Password: "123",
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(rsp)

	// rsp, err := mockUser.CreateUser(context.Background(), &biz.User{
	// 	Mobile:   fmt.Sprintf("1388888888%d", 1),
	// 	Nickname: fmt.Sprintf("YWWW%d", 1),
	// 	Password: "123",
	// })

	// 断言期望的返回值
	// assert.NoError(t, err)
	// assert.Equal(t, "13888888881", rsp.Mobile)
	// assert.Equal(t, "YWWW1", rsp.Nickname)

	// 断言模拟对象的期望行为是否被满足
	// mockUser.AssertExpectations(t)
}

func TestGetUserByMobile(t *testing.T) {
	rsp, err := userClient.GetUserByMobile(context.Background(), &v1.MobileRequest{
		Mobile: fmt.Sprintf("1388888888%d", 1),
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(rsp)
}
