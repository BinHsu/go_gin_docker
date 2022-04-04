package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

var local_port_test int = 8080

type Result_test struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Status int    `json:"status"`
}

func TestReadNoTask(t *testing.T) {
	path := fmt.Sprintf("http://127.0.0.1:%d/tasks", local_port_test)
	response, err := http.Get(path)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Testing!")
	//defer req.Body.Close()
	if response == nil {
		t.Fatal("nil response")
	}
	if response.StatusCode != 204 {
		t.Fatal("failed status code")
	}
	t.Log("Testing!")
}

func TestCreatTaskFailed(t *testing.T) {
	// empty request body
	var empty_value map[string]string
	json_data, err := json.Marshal(empty_value)
	if err != nil {
		t.Fatal(err)
	}

	path := fmt.Sprintf("http://127.0.0.1:%d/task", local_port_test)
	resp, err := http.Post(path, "application/json",
		bytes.NewBuffer(json_data))
	if err != nil {
		t.Fatal(err)
	}

	if resp == nil || resp.StatusCode != 400 {
		t.Fatal("unexpected status code")
	}
	// invalid name type
	task := 1
	value := map[string]int{"name": task}
	json_data, err = json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}

	resp, err = http.Post(path, "application/json",
		bytes.NewBuffer(json_data))
	if err != nil {
		t.Fatal(err)
	}

	if resp == nil || resp.StatusCode != 400 {
		t.Fatal("unexpected status code")
	}
}

func TestCreateTask(t *testing.T) {
	task := "買晚餐"
	value := map[string]string{"name": task}
	json_data, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}

	path := fmt.Sprintf("http://127.0.0.1:%d/task", local_port_test)
	resp, err := http.Post(path, "application/json",
		bytes.NewBuffer(json_data))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		t.Fatal("unexpected status code")
	}

	var res map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&res)
	if nil == res["result"] {
		t.Fatal("null result")
	}

	r := res["result"].(map[string]interface{})
	result := Result_test{int(r["id"].(float64)), r["name"].(string), int(r["status"].(float64))}
	verifyResult := Result_test{1, task, 0}
	if result != verifyResult {
		t.Fatalf("unexpected result: %v", result)
	}

	task = "買早餐"
	value = map[string]string{"name": task}
	json_data, err = json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}

	resp, err = http.Post(path, "application/json",
		bytes.NewBuffer(json_data))
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil || resp.Body == nil {
		t.Fatal("null resp")
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		t.Fatal("unexpected status code")
	}

	json.NewDecoder(resp.Body).Decode(&res)
	if nil == res["result"] {
		t.Fatal("null result")
	}

	r = res["result"].(map[string]interface{})
	result = Result_test{int(r["id"].(float64)), r["name"].(string), int(r["status"].(float64))}
	verifyResult = Result_test{2, task, 0}
	if result != verifyResult {
		t.Fatalf("unexpected result: %v", result)
	}
}

func TestUpdateTaskFailed(t *testing.T) {
	task := "買午餐"
	value := map[string]string{"name": task}
	json_data, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}

	path := fmt.Sprintf("http://127.0.0.1:%d/task/999", local_port_test)
	req, err := http.NewRequest(http.MethodPut, path,
		bytes.NewBuffer(json_data))
	if err != nil {
		t.Fatal(err)
	}

	// set the request header Content-Type for json
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	// initialize http client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp == nil || resp.StatusCode != 404 {
		t.Fatal("unexpected status code")
	}

	path = fmt.Sprintf("http://127.0.0.1:%d/task/1", local_port_test)
	invalid_value := map[string]int{"name": -1}
	json_data, err = json.Marshal(invalid_value)
	if err != nil {
		t.Fatal(err)
	}
	req, err = http.NewRequest(http.MethodPut, path,
		bytes.NewBuffer(json_data))
	if err != nil {
		t.Fatal(err)
	}

	// set the request header Content-Type for json
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	// initialize http client
	resp, err = client.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp == nil || resp.StatusCode != 400 {
		t.Fatal("unexpected status code")
	}

	path = fmt.Sprintf("http://127.0.0.1:%d/task/1", local_port_test)
	invalid_value = map[string]int{"status": -1}
	json_data, err = json.Marshal(invalid_value)
	if err != nil {
		t.Fatal(err)
	}
	req, err = http.NewRequest(http.MethodPut, path,
		bytes.NewBuffer(json_data))
	if err != nil {
		t.Fatal(err)
	}

	// set the request header Content-Type for json
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	// initialize http client
	resp, err = client.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp == nil || resp.StatusCode != 400 {
		t.Fatal("unexpected status code")
	}
}

func TestUpdateTask(t *testing.T) {
	task := "買午餐"
	value := map[string]string{"name": task}
	json_data, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}

	path := fmt.Sprintf("http://127.0.0.1:%d/task/1", local_port_test)
	req, err := http.NewRequest(http.MethodPut, path,
		bytes.NewBuffer(json_data))
	if err != nil {
		t.Fatal(err)
	}

	// set the request header Content-Type for json
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	// initialize http client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp == nil || resp.Body == nil {
		t.Fatal("null resp")
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Fatal("unexpected status code")
	}

	var res map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&res)
	if nil == res["result"] {
		t.Fatal("null result")
	}

	r := res["result"].(map[string]interface{})
	if r["id"] != nil || r["status"] != nil {
		t.Fatal("unexpected response")
	}
	result := Result_test{0, r["name"].(string), 0}
	verifyResult := Result_test{0, task, 0}
	if result != verifyResult {
		t.Fatalf("unexpected result: %v", result)
	}

	path = fmt.Sprintf("http://127.0.0.1:%d/task/2", local_port_test)
	value_int := map[string]int{"id": 2, "status": 1}
	json_data, err = json.Marshal(value_int)
	if err != nil {
		t.Fatal(err)
	}
	req, err = http.NewRequest(http.MethodPut, path,
		bytes.NewBuffer(json_data))
	if err != nil {
		t.Fatal(err)
	}

	// set the request header Content-Type for json
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	// initialize http client
	resp, err = client.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp == nil || resp.Body == nil {
		t.Fatal("null resp")
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Fatal("unexpected status code")
	}
	json.NewDecoder(resp.Body).Decode(&res)
	if nil == res["result"] {
		t.Fatal("null result")
	}

	r = res["result"].(map[string]interface{})
	if r["name"] != nil {
		t.Fatal("unexpected response")
	}
	result = Result_test{int(r["id"].(float64)), "", int(r["status"].(float64))}
	verifyResult = Result_test{2, "", 1}
	if result != verifyResult {
		t.Fatalf("unexpected result: %v", result)
	}
}

func TestReadTasks(t *testing.T) {
	path := fmt.Sprintf("http://127.0.0.1:%d/tasks", local_port_test)
	resp, err := http.Get(path)
	if err != nil {
		t.Fatal(err)
	}

	if resp == nil || resp.Body == nil {
		t.Fatal("nil response")
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Fatal("failed status code")
	}

	var res map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&res)
	if nil == res["result"] {
		t.Fatal("null result")
	}
	r, _ := json.Marshal(res["result"])

	var verifyResult []Result_test
	verifyResult = append(verifyResult, Result_test{1, "買午餐", 0})
	verifyResult = append(verifyResult, Result_test{2, "買早餐", 1})
	j, _ := json.Marshal(verifyResult)
	if diff := bytes.Compare(r, j); diff != 0 {
		t.Fatal("unexpected response")
	}
}

func TestDeleteTaskFailed(t *testing.T) {
	path := fmt.Sprintf("http://127.0.0.1:%d/task/AA", local_port_test)
	req, err := http.NewRequest(http.MethodDelete, path, nil)
	if err != nil {
		t.Fatal(err)
	}
	// initialize http client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 400 {
		t.Fatal("unexpected status code")
	}
}

func TestDeleteTask(t *testing.T) {
	path := fmt.Sprintf("http://127.0.0.1:%d/task/3", local_port_test)
	req, err := http.NewRequest(http.MethodDelete, path, nil)
	if err != nil {
		t.Fatal(err)
	}
	// initialize http client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp == nil || resp.StatusCode != 200 {
		t.Fatal("unexpected status code")
	}

	path = fmt.Sprintf("http://127.0.0.1:%d/tasks", local_port_test)
	resp, err = http.Get(path)
	if err != nil {
		t.Fatal(err)
	}

	if resp == nil || resp.Body == nil {
		t.Fatal("nil response")
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Fatal("failed status code")
	}

	var res map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&res)
	if nil == res["result"] {
		t.Fatal("null result")
	}
	r, _ := json.Marshal(res["result"])

	var verifyResult []Result_test
	verifyResult = append(verifyResult, Result_test{1, "買午餐", 0})
	verifyResult = append(verifyResult, Result_test{2, "買早餐", 1})
	j, _ := json.Marshal(verifyResult)
	if diff := bytes.Compare(r, j); diff != 0 {
		t.Fatal("unexpected response")
	}

	path = fmt.Sprintf("http://127.0.0.1:%d/task/1", local_port_test)
	req, err = http.NewRequest(http.MethodDelete, path, nil)
	if err != nil {
		t.Fatal(err)
	}
	// initialize http client
	resp, err = client.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp == nil || resp.StatusCode != 200 {
		t.Fatal("unexpected status code")
	}

	path = fmt.Sprintf("http://127.0.0.1:%d/tasks", local_port_test)
	resp, err = http.Get(path)
	if err != nil {
		t.Fatal(err)
	}

	if resp == nil || resp.Body == nil {
		t.Fatal("nil response")
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Fatal("failed status code")
	}

	json.NewDecoder(resp.Body).Decode(&res)
	if nil == res["result"] {
		t.Fatal("null result")
	}
	r, _ = json.Marshal(res["result"])

	var verifyResult2 []Result_test
	verifyResult2 = append(verifyResult2, Result_test{2, "買早餐", 1})
	j, _ = json.Marshal(verifyResult2)
	if diff := bytes.Compare(r, j); diff != 0 {
		t.Fatal("unexpected response")
	}
}
