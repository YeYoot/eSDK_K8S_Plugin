/*
 *  Copyright (c) Huawei Technologies Co., Ltd. 2022-2022. All rights reserved.
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *       http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

package client

import (
	"context"
	"os"
	"path"
	"reflect"
	"testing"

	"bou.ke/monkey"
	. "github.com/smartystreets/goconvey/convey"

	"huawei-csi-driver/utils/log"
)

const (
	logDir  = "/var/log/huawei/"
	logName = "clientNamespaceTest.log"
)

var testClient *Client

func TestMain(m *testing.M) {
	if err := log.InitLogging(logName); err != nil {
		log.Errorf("init logging: %s failed. error: %v", logName, err)
		os.Exit(1)
	}

	logFile := path.Join(logDir, logName)
	defer func() {
		if err := os.RemoveAll(logFile); err != nil {
			log.Errorf("Remove file: %s failed. error: %s", logFile, err)
		}
	}()

	testClient = NewClient("https://192.168.125.*:8088", "dev-account", "dev-password", "50")

	m.Run()
}

func TestAllowNfsShareAccess(t *testing.T) {
	Convey("Normal", t, func() {
		guard := monkey.PatchInstanceMethod(reflect.TypeOf(testClient), "Post",
			func(_ *Client, _ context.Context, _ string, _ map[string]interface{}) (map[string]interface{}, error) {
				return map[string]interface{}{
					"data":   map[string]interface{}{},
					"result": map[string]interface{}{"code": float64(0), "description": ""},
				}, nil
			})
		defer guard.Unpatch()

		err := testClient.AllowNfsShareAccess(context.TODO(), &AllowNfsShareAccessRequest{
			AccessName:  "test",
			ShareId:     "test",
			AccessValue: 0,
			AllSquash:   1,
			RootSquash:  1,
			AccountId:   "0",
		})
		So(err, ShouldBeNil)
	})

	Convey("Result Code Not Exist", t, func() {
		guard := monkey.PatchInstanceMethod(reflect.TypeOf(testClient), "Post",
			func(_ *Client, _ context.Context, _ string, _ map[string]interface{}) (map[string]interface{}, error) {
				return map[string]interface{}{
					"data": map[string]interface{}{},
				}, nil
			})
		defer guard.Unpatch()

		err := testClient.AllowNfsShareAccess(context.TODO(), &AllowNfsShareAccessRequest{
			AccessName:  "test",
			ShareId:     "test",
			AccessValue: 0,
			AllSquash:   1,
			RootSquash:  1,
			AccountId:   "0",
		})
		So(err, ShouldBeError)
	})

	Convey("Client Already Exist", t, func() {
		guard := monkey.PatchInstanceMethod(reflect.TypeOf(testClient), "Post",
			func(_ *Client, _ context.Context, _ string, _ map[string]interface{}) (map[string]interface{}, error) {
				return map[string]interface{}{
					"data":   map[string]interface{}{},
					"result": map[string]interface{}{"code": float64(clientAlreadyExist), "description": ""},
				}, nil
			})
		defer guard.Unpatch()

		err := testClient.AllowNfsShareAccess(context.TODO(), &AllowNfsShareAccessRequest{
			AccessName:  "test",
			ShareId:     "test",
			AccessValue: 0,
			AllSquash:   1,
			RootSquash:  1,
			AccountId:   "0",
		})
		So(err, ShouldBeNil)
	})

	Convey("Error code is not zero", t, func() {
		guard := monkey.PatchInstanceMethod(reflect.TypeOf(testClient), "Post",
			func(_ *Client, _ context.Context, _ string, _ map[string]interface{}) (map[string]interface{}, error) {
				return map[string]interface{}{
					"data":   map[string]interface{}{},
					"result": map[string]interface{}{"code": float64(100), "description": ""},
				}, nil
			})
		defer guard.Unpatch()

		err := testClient.AllowNfsShareAccess(context.TODO(), &AllowNfsShareAccessRequest{
			AccessName:  "test",
			ShareId:     "test",
			AccessValue: 0,
			AllSquash:   1,
			RootSquash:  1,
			AccountId:   "0",
		})
		So(err, ShouldBeError)
	})
}
