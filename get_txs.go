package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"math/big"
	"net/http"
	"os"
)

const ApiKey = "xxx"

func main() {
	// 1 get tnx api
	tnxURL := fmt.Sprintf("https://api.etherscan.io/api?module=account&action=txlist&address=0xdAC17F958D2ee523a2206206994597C13D831ec7&startblock=16238340&endblock=16238360&sort=asc&apikey=%s", ApiKey)
	tnxResp, err := http.Get(tnxURL)
	if err != nil {
		fmt.Println("请求失败：", err)
		return
	}
	defer tnxResp.Body.Close()

	var data interface{} // 定义一个变量来保存JSON数据
	err = json.NewDecoder(tnxResp.Body).Decode(&data)
	if err != nil {
		return
	}

	// 格式化JSON数据
	indentData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return
	}
	err = os.WriteFile("tnx-format.json", indentData, 0644)
	if err != nil {
		fmt.Println("write file error:", err)
		panic("write file error")
	}
	//====================
	//tnxBody, err := io.ReadAll(tnxResp.Body)
	//if err != nil {
	//	fmt.Println("读取响应失败：", err)
	//	return
	//}
	//buffer := bytes.NewBuffer(tnxBody)
	//err = os.WriteFile("result.json", buffer.Bytes(), 0644)
	//if err != nil {
	//	fmt.Println("write file error:", err)
	//	panic("write file error")
	//}

	// 2 get balance api
	balanceURL := fmt.Sprintf("https://api.etherscan.io/api?module=account&action=balance&address=0xde0b295669a9fd93d5f28d9ec85e40f4cb697bae&tag=latest&apikey=%s", ApiKey)
	balanceResp, err := http.Get(balanceURL)
	if err != nil {
		fmt.Println("请求失败：", err)
		return
	}
	defer balanceResp.Body.Close()

	tnxBody, err := io.ReadAll(balanceResp.Body)
	if err != nil {
		fmt.Println("读取响应失败：", err)
		return
	}
	fmt.Println("balance=", string(tnxBody))

	// 使用匿名结构体解析响应体
	var data1 struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Result  string `json:"result"`
	}
	err = json.Unmarshal(tnxBody, &data1)
	if err != nil {
		fmt.Println("解析JSON失败：", err)
		return
	}

	num := new(big.Int)
	num, ok := num.SetString(data1.Result, 10)
	if !ok {
		fmt.Println("转换失败")
		return
	}
	divisor := big.NewInt(int64(math.Pow10(18)))
	quotient := new(big.Int)
	remainder := new(big.Int)
	quotient.DivMod(num, divisor, remainder)

	fmt.Println("num:", num)
	fmt.Println("big.Int:", fmt.Sprintf("%s.%s", quotient.String(), remainder.String()))
}
