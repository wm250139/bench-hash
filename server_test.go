package main_test

import (
	client_grpc "bench-hash/client-grpc"
	client_inproc "bench-hash/client-inproc"
	client_mqtt "bench-hash/client-mqtt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func BenchmarkHTTPGet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		resp, err := http.Get("http://localhost:3030?q=hello")
		if err != nil {
			b.Fatal("unable to get", err)
		}
		resp.Body.Close()
	}
}

func BenchmarkHTTPClient(b *testing.B) {
	client := http.Client{}

	for i := 0; i < b.N; i++ {
		req, _ := http.NewRequest("GET", "http://localhost:3030?q=hello", nil)
		resp, err := client.Do(req)
		if err != nil {
			b.Fatal("unable to get", err)
		}

		hashBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			b.Fatal("unable to read response body", err)
		}

		hash := strings.TrimSpace(string(hashBytes))

		if hash != "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c" {
			b.Fatalf("got wrong hash response: '%s'", hash)
		}

		resp.Body.Close()
	}
}

func BenchmarkGRPC(b *testing.B) {
	client, err := client_grpc.NewHashClient("localhost:3033")
	if err != nil {
		b.Fatal("Unable to create client", err)
	}
	defer client.Close()

	for i := 0; i < b.N; i++ {
		resp, err := client.Hash("hello")
		if err != nil {
			b.Fatal("unable to get", err)
		}

		if resp != "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c" {
			b.Fatal("got wrong hash response:", resp)
		}
	}
}

func BenchmarkMQTT(b *testing.B) {
	client, err := client_mqtt.NewMQTTClient("tcp", "localhost:1883")
	if err != nil {
		b.Fatal("Unable to create client", err)
	}
	defer client.Close()

	for i := 0; i < b.N; i++ {
		hash, err := client.Hash("hello")
		if err != nil {
			b.Fatal(err)
		}

		if hash != "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c" {
			b.Fatal("got wrong hash response:", hash)
		}
	}
}

func BenchmarkInProc(b *testing.B) {
	client := client_inproc.NewInProcessClient()

	for i := 0; i < b.N; i++ {
		hash, err := client.Hash("hello")
		if err != nil {
			b.Fatal(err)
		}

		if hash != "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c" {
			b.Fatal("got wrong hash response:", hash)
		}
	}
}
