package main

import (
	boost "github.com/primev/builder-boost/pkg"
)

const (
	url = "http://localhost:18550" + boost.PathSubmitBlock
)

// func main() {
// 	fmt.Print("submitting test... ")
// 	if err := sendTest(); err != nil {
// 		panic(err)
// 	}
// }

// func sendTest() error {

// 	request := validExampleTestRequest()

// 	b, err := json.Marshal(request)
// 	if err != nil {
// 		return err
// 	}

// 	fmt.Println(hexutil.Bytes(b[:]).String())

// 	resp, err := http.Post(url, "application/json", bytes.NewBuffer(b))
// 	if err != nil {
// 		return err
// 	}

// 	if resp.StatusCode != 200 {
// 		body, err := io.ReadAll(resp.Body)
// 		if err != nil {
// 			return err
// 		}
// 		return errors.WithMessage(fmt.Errorf("invalid return code, expected 200 - received %d", resp.StatusCode), string(body))
// 	}

// 	fmt.Println(resp.Status)
// 	return nil
// }

// func validExampleTestRequest() *types.ExRequest {
// 	return &types.ExRequest{
// 		Field: "123",
// 		Test:  "123",
// 	}
// }
