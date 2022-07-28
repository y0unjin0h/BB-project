package main

// import
import(
	"encoding/json"
	"fmt"
	"time"
	"log"

	"github.com/thoas/go-funk"

	"github.com/golang/protobuf/ptypes"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)
// 체인코드 구조체
type SimpleChaincode struct {
	contractapi.Contract
}
// WS Bubble 구조체
type Bubble struct {
	ObjectType string `json:"docType"`
	BType string `json:"btype"` // "WashingBubble", "FreshBubble", "GatheringBubble"
	BID string `json:"bid"`
	TotalPrice int `json:"totalprice"`
	UniCount int `json:"unicount"`
	Status string `json:"status"` // "Registered", "Executing", "Finalized"
	Unies []string `json:"unies"`
}

type HistoryQueryResult struct {
	Record    *Bubble    `json:"record"`
	TxId     string    `json:"txId"`
	Timestamp time.Time `json:"timestamp"`
	IsDelete  bool      `json:"isDelete"`
}

// RegisterBubble 함수
func (t *SimpleChaincode) RegisterBubble(ctx contractapi.TransactionContextInterface, btype string, bid string, totalprice int, unicount int, uniID string) error {
	fmt.Println("- start register bubble")

	// 기등록 버블 검색
	bubbleAsBytes, err := ctx.GetStub().GetState(bid)
	if err != nil {
		return fmt.Errorf("Failed to get bubble: " + err.Error())
	} else if bubbleAsBytes != nil {
		return fmt.Errorf("This bubble already exists: " + bid)
	}

	// 구조체 생성 -> 마샬링 -> PutState
	bubble := &Bubble{"ServiceBubble", btype, bid, totalprice, unicount, "registered", []string{uniID}}
	bubbleJSONasBytes, err := json.Marshal(bubble)
	if err != nil {
		return err
	}
	err = ctx.GetStub().PutState(bid, bubbleJSONasBytes)
	if err != nil {
		return err
	}
	return nil
}

// ReadBubble 함수
func (t *SimpleChaincode) ReadBubble(ctx contractapi.TransactionContextInterface, bid string) (*Bubble, error) {
	fmt.Println("- start read bubble")
	
	// 기등록 버블 검색
	bubbleAsBytes, err := ctx.GetStub().GetState(bid)
	if err != nil {
		return nil, fmt.Errorf("Failed to get bubble: " + err.Error())
	} else if bubbleAsBytes == nil {
		return nil, fmt.Errorf("This bubble does not exists: " + bid)
	}
	
	var bubble Bubble
	err = json.Unmarshal(bubbleAsBytes, &bubble)
	if err != nil {
		return nil, err
	}

	return &bubble, nil

}

// JoinBubble 함수
func (t *SimpleChaincode) JoinBubble(ctx contractapi.TransactionContextInterface, bid string, uniID string ) error {
	fmt.Println("- start transfer bubble")
	
	// 기등록 버블 검색
	bubbleAsBytes, err := ctx.GetStub().GetState(bid)
	if err != nil {
		return fmt.Errorf("Failed to get bubble: " + err.Error())
	} else if bubbleAsBytes == nil {
		return fmt.Errorf("This bubble does not exist: " + bid)
	}

	// unmarshal 시키는거 먼저
	bubble := Bubble{}
	_ = json.Unmarshal(bubbleAsBytes, &bubble)

	fmt.Println(bubble.UniCount , len(bubble.Unies))
	fmt.Println(bubble.Unies , uniID)

	if bubble.UniCount <= len(bubble.Unies) { // 버블 정원 다 찼으면 에러
		return fmt.Errorf("Error: Bubble is full!")
	} else if funk.Contains(bubble.Unies, uniID) == true { // 이미 참여한 유니버블이면 에러
		return fmt.Errorf("Error: UniBubble is already registered!")
	}


	bubble.Status = "Joining"
	bubble.Unies = append(bubble.Unies, uniID)


	if bubble.UniCount == len(bubble.Unies) { // 버블 정원 다 찼으면 에러
		bubble.Status = "Executing"
		//return fmt.Errorf("Error: Bubble is not ready!")
	} 
	

	bubbleJSONasBytes, err := json.Marshal(bubble)
	if err != nil {
		return err
	}
	err = ctx.GetStub().PutState(bid, bubbleJSONasBytes)
	if err != nil {
		return err
	}
	return nil
}

// ExecuteBubble 함수
func (t *SimpleChaincode) ExecuteBubble(ctx contractapi.TransactionContextInterface, bid string) error {
	fmt.Println("- start transfer bubble")
	
	// 기등록 버블 검색
	bubbleAsBytes, err := ctx.GetStub().GetState(bid)
	if err != nil {
		return fmt.Errorf("Failed to get bubble: " + err.Error())
	} else if bubbleAsBytes == nil {
		return fmt.Errorf("This bubble does not exist: " + bid)
	}

	// unmarshal 시키는거 먼저
	bubble := Bubble{}
	_ = json.Unmarshal(bubbleAsBytes, &bubble)


	if bubble.UniCount != len(bubble.Unies) { // 버블 정원 다 찼으면 에러
		return fmt.Errorf("Error: Bubble is not ready!")
	} 

	bubble.Status = "Executing"

	bubbleJSONasBytes, err := json.Marshal(bubble)
	if err != nil {
		return err
	}
	err = ctx.GetStub().PutState(bid, bubbleJSONasBytes)
	if err != nil {
		return err
	}
	return nil
}

// FinalizeBubble 함수
func (t *SimpleChaincode) FinalizeBubble(ctx contractapi.TransactionContextInterface, bid string, uniID string) error {
	fmt.Println("- start transfer bubble")
	
	// 기등록 버블 검색
	bubbleAsBytes, err := ctx.GetStub().GetState(bid)
	if err != nil {
		return fmt.Errorf("Failed to get bubble: " + err.Error())
	} else if bubbleAsBytes == nil {
		return fmt.Errorf("This bubble does not exist: " + bid)
	}

	// unmarshal 시키는거 먼저
	bubble := Bubble{}
	_ = json.Unmarshal(bubbleAsBytes, &bubble)

	if uniID != bubble.Unies[0] { // 버블 정원 다 찼으면 에러
		return fmt.Errorf("Error: UniBubble is not StartUniBubble!")
	} 

	bubble.Status = "Finalized"

	bubbleJSONasBytes, err := json.Marshal(bubble)
	if err != nil {
		return err
	}
	err = ctx.GetStub().PutState(bid, bubbleJSONasBytes)
	if err != nil {
		return err
	}
	return nil
}



// GetHistoryForBubble 함수
func (t *SimpleChaincode) GetBubbleHistory(ctx contractapi.TransactionContextInterface, bid string) ([]HistoryQueryResult, error) {
	log.Printf("GetBubbleHistory: ID %v", bid)

	resultsIterator, err := ctx.GetStub().GetHistoryForKey(bid)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var records []HistoryQueryResult
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var bubble Bubble
		if len(response.Value) > 0 {
			err = json.Unmarshal(response.Value, &bubble)
			if err != nil {
				return nil, err
			}
		} else {
			bubble = Bubble{
				BID: bid,
			}
		}

		timestamp, err := ptypes.Timestamp(response.Timestamp)
		if err != nil {
			return nil, err
		}

		record := HistoryQueryResult{
			TxId:      response.TxId,
			Timestamp: timestamp,
			Record:    &bubble,
			IsDelete:  response.IsDelete,
		}
		records = append(records, record)
	}

	return records, nil
}

// main 함수
func main() {
	chaincode, err := contractapi.NewChaincode(&SimpleChaincode{})
	if err != nil {
		log.Panicf("Error creating Bubble chaincode: %v", err)
	}

	if err := chaincode.Start(); err != nil {
		log.Panicf("Error starting Bubble chaincode: %v", err)
	}
}
