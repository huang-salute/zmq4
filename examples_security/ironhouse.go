//  The Stonehouse Pattern
//
//  Where we allow any clients to connect, but we promise clients
//  that we are who we claim to be, and our conversations won't be
//  tampered with or modified, or spied on.

package main

import (
	zmq "github.com/pebbe/zmq4"

	"fmt"
	"log"
	"runtime"
)

func main() {

    //  Start authentication engine
	zmq.AuthSetVerbose(true)
	zmq.AuthStart()
    zmq.AuthAllow("127.0.0.1")

    //  We need two certificates, one for the client and one for
    //  the server. The client must know the server's public key
    //  to make a CURVE connection.
    client_public, client_secret, err := zmq.NewCurveKeypair()
	checkErr(err)
    server_public, server_secret, _ := zmq.NewCurveKeypair()
	checkErr(err)

	//  Tell authenticator to use this public client key
	zmq.AuthConfigureCurve("*", client_public)

    //  Create and bind server socket
	server, _ := zmq.NewSocket(zmq.PUSH)
	server.SetCurveSecretkey(server_secret)
	server.SetCurveServer(1)
    server.Bind("tcp://*:9000")

	//  Create and connect client socket
	client, _ := zmq.NewSocket(zmq.PULL)
	client.SetCurveServerkey(server_public)
	client.SetCurvePublickey(client_public)
	client.SetCurveSecretkey(client_secret)
	client.Connect("tcp://127.0.0.1:9000")

	//  Send a single message from server to client
	_, err = server.Send("Hello", 0)
	checkErr(err)
	message, err := client.Recv(0)
	checkErr(err)
	if message != "Hello" {
		log.Fatalln(message, "!= Hello")
	}

	zmq.AuthStop()

	fmt.Println("Ironhouse test OK")

}

func checkErr(err error) {
	if err != nil {
		log.SetFlags(0)
		_, filename, lineno, ok := runtime.Caller(1)
		if ok {
			log.Fatalf("%v:%v: %v", filename, lineno, err)
		} else {
			log.Fatalln(err)
		}
	}
}


