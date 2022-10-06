package main

import (
	"fmt"
	"net/http"
	"io"
	"bytes"
	"io/ioutil"
	"encoding/json"
	"pongpongi.com/psqldb"
	"pongpongi.com/osdep"
	"pongpongi.com/blockchain/evm"
)

type RequestRecordType struct {
	Table string `json:"table"`
	Field string `json:"field"`
	Value string `json:"value"`
}

func writePostJsonRequest(w http.ResponseWriter, req *http.Request) {
	// fmt.Println("Method", req.Method)
	// fmt.Println("URL", req.URL.String())
	// fmt.Println("Proto", req.Proto)
	// fmt.Println("ProtoMajor", req.ProtoMajor)
	// fmt.Println("ProtoMinor", req.ProtoMinor)
	// fmt.Println("Header")
	var clientIp string;
	for k, v := range req.Header {
		// fmt.Println("\t", k)
		for _, vv := range v {
			// fmt.Println("\t\t", vv)
			if k == "X-Real-Ip" {
				clientIp = vv;
			}
		}
	}
	// fmt.Println("Body", req.Body)

	var content string
	if req.Method == "POST" {
		var p RequestRecordType

		if err := json.NewDecoder(req.Body).Decode(&p); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Do something with the Person struct...
		t, _ := json.Marshal(p)
		fmt.Fprintf(w, "Post DB Write: %+v", string(t))
		content = string(t)

		// body, err := ioutil.ReadAll(req.Body)
		// if err != nil {
		// 	fmt.Printf("Error reading body: %v\n", err)
		// 	http.Error(w, "can't read body", http.StatusBadRequest)
		// 	return
		// }
		// if len(body) == 0 {
		// 	fmt.Println("\tEmpty Body")
		// } else {
		// 	fmt.Println("", string(body))
		// 	if err := json.Unmarshal(body, &p); err != nil {
		// 		fmt.Printf("Error reading body: %v\n", err)
		// 		http.Error(w, "can't read body", http.StatusBadRequest)
		// 		return
		// 	} else {
		// 		fmt.Println("\t", p)
		// 	}
		// }
	} else if req.Method == "GET" {
		// https://pongpongi.com/api/write?p[]=q%20ss&p[]=r&p[]=s&owner=My%20Wallet
		// response: {"owner":["My Wallet"],"p[]":["q ss","r","s"]}
		// https://pongpongi.com/api/write?p[]=q&p[]=r&p[]=s
		// response: {"p[]":["q","r","s"]}

		query := req.URL.Query()
		// fmt.Println("\tHost", req.URL.EscapedPath())
		var objMap []map[string]interface{} = make([]map[string]interface{}, 1)
		objMap[0] = make(map[string]interface{})

		for k, v := range query {
			objMap[0][k] = v
			// fmt.Println("\t", k, ":", v)
			// for idx, val := range v {
				// fmt.Println("\t\t", idx, val)
			// }
		}

		procErr := ProcGet(req.URL.EscapedPath(), objMap[0], w)
		if procErr != nil {
			fmt.Println(procErr)
		}

		btArray, err := json.Marshal(objMap[0])
		if err == nil {
			// fmt.Fprintf(w, "Post DB Write: %+v", string(btArray))
			content = string(btArray)

			// var uncomp []map[string]interface{}
			// if err2 := json.Unmarshal(btArray, &uncomp); err2 != nil {
			// 	osdep.Check(err2)
			// }
			// fmt.Println(uncomp[0]["p[]"]) // to parse out your value
		} else {
			osdep.Check(err)
		}
	}

	// fmt.Println("ContentLength", req.ContentLength)
	// fmt.Println("TransferEncoding")
	// for _, v := range req.TransferEncoding {
	// 	fmt.Println("\t", v)
	// }
	// fmt.Println("Close", req.Close)
	// fmt.Println("Host", req.Host)
	// fmt.Println("Form", req.Form)
	// fmt.Println("PostForm", req.PostForm)
	// fmt.Println("MultipartForm", req.MultipartForm)
	// fmt.Println("Trailer")
	// for k, v := range req.Trailer {
	// 	fmt.Println("\t", k)
	// 	for _, vv := range v {
	// 		fmt.Println("\t\t", vv)
	// 	}
	// }
	// fmt.Println("RemoteAddr", req.RemoteAddr)
	// fmt.Println("RequestURI", req.RequestURI)
	// fmt.Println("TLS", req.TLS)
	// fmt.Println("Response", req.Response)

	writeLoginHistory(clientIp, req.Method, content)
}

type MyResponseWriter struct {
    http.ResponseWriter
    buf *bytes.Buffer
}

func (mrw *MyResponseWriter) Write(p []byte) (int, error) {
    return mrw.buf.Write(p)
}

func writeRequest2(handler http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        body, err := ioutil.ReadAll(r.Body)
        if err != nil {
            fmt.Printf("Error reading body: %v\n", err)
            http.Error(w, "can't read body", http.StatusBadRequest)
            return
        }

        // Work / inspect body. You may even modify it!

        // And now set a new body, which will simulate the same data we read:
        r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

        // Create a response wrapper:
        mrw := &MyResponseWriter{
            ResponseWriter: w,
            buf:            &bytes.Buffer{},
        }

        // Call next handler, passing the response wrapper:
        handler.ServeHTTP(mrw, r)

        // Now inspect response, and finally send it out:
        // (You can also modify it before sending it out!)
        if _, err := io.Copy(w, mrw.buf); err != nil {
            fmt.Printf("Failed to send out response: %v\n", err)
        }
    })
}

func writeLoginHistory(ip string, method string, content string) {
	fmtString := fmt.Sprintf("insert into %s(ip, method, content) values($1, $2, $3)", psqldb.GetHistoryTableNable())
	psqldb.DbQuery(fmtString, ip, method, content)
	// ret := psqldb.DbQuery(fmtString, ip, method, content)
	// fmt.Println(ret)
}

func main() {
	evm.InitEVM()
	psqldb.InitDB()

	mux := http.NewServeMux()
	mux.HandleFunc("/", writePostJsonRequest)

	fmt.Println("started server at port 11331")
	http.ListenAndServe(":11331", mux)
}
