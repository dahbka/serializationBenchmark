package main

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"time"
)

const (
	ALL = 1000
)

/**
Struct {} for testing serialization
*/
type TestStruct struct {
	StringData string
	Slice      []int64
	IntData    int64
	FloatData  float64
}

/*
	generate seed for Struct init
*/
var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

/*
	charset for string generation
*/
const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateTestStruct(length int) TestStruct {
	testStruct := TestStruct{
		StringData: generateString(length, charset),
		Slice:      generateIntSlice(length),
		IntData:    rand.Int63(),
		FloatData:  rand.Float64(),
	}
	return testStruct
}

/*
	return generated string using @seededRand, charset and length of string
*/
func generateString(length int, charset string) string {
	buff := make([]byte, length)
	for i := range buff {
		buff[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(buff)
}

/*
	return generated string using @seededRand and length for int Slice
*/
func generateIntSlice(length int) []int64 {
	buff := make([]int64, length)
	for i := range buff {
		buff[i] = rand.Int63()
	}
	return buff
}

func (s TestStruct) String() string {
	v, err := json.Marshal(s)
	if err != nil {
		log.Fatal("Marshal failed")
	}
	return string(v)
}

var testSlice []TestStruct
var testSliceJson []TestStruct
var testSliceGob []TestStruct
var testSliceXml []TestStruct

func init() {
	testSlice = make([]TestStruct, ALL)
	for i := 0; i < ALL; i++ {
		testSlice[i] = GenerateTestStruct(ALL)
	}

}

func toJsonBytes() []byte {
	answer, err := json.Marshal(testSlice)
	if err != nil {
		log.Fatal(err)
	}
	return answer
}

func loadJsonBytes(input []byte) []TestStruct {
	ss := make([]TestStruct, ALL)
	err := json.Unmarshal(input, &ss)
	if err != nil {
		log.Fatal(err)
	}
	return ss
}

func toGobBytes() []byte {
	stream := &bytes.Buffer{}
	en := gob.NewEncoder(stream)
	err := en.Encode(testSlice)
	if err != nil {
		log.Fatal(err)
	}
	return stream.Bytes()
}

func loadGobBytes(input []byte) []TestStruct {
	dec := gob.NewDecoder(bytes.NewBuffer(input))
	ss := make([]TestStruct, ALL)
	err := dec.Decode(&ss)
	if err != nil {
		log.Fatal(err)
	}
	return ss
}

func toXmlBytes() []byte {
	answer, err := xml.Marshal(testSlice)
	if err != nil {
		log.Fatal(err)
	}
	return answer
}
func loadXmlBytes(input []byte) []TestStruct {
	ss := make([]TestStruct, 0, ALL)
	for len(ss) < ALL {
		err := xml.Unmarshal(input, &ss)
		if err != nil {
			log.Fatal(err)
		}

	}
	return ss
}

func benchmark() {
	//create your file with desired read/write permissions

	f, err := os.OpenFile("benchmark.txt", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}
	//defer to close when you're done with it, not because you think it's idiomatic!
	defer f.Close()

	mw := io.MultiWriter(os.Stdout, f)
	//set output of logs to f
	log.SetOutput(mw)

	startSer := time.Now()
	gobbytes := toGobBytes()
	ss1 := loadGobBytes(gobbytes)
	endSer := time.Now()

	fGobWrite, err := os.OpenFile("gob.txt", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal("GobError: OpenFile ", err)
	}

	defer fGobWrite.Close()

	encGob := gob.NewEncoder(fGobWrite)
	if err := encGob.Encode(ss1); err != nil {
		log.Fatal("GobError: Encode ", err)
	}

	startDes := time.Now()
	fGobRead, err := ioutil.ReadFile("./gob.txt")
	if err != nil {
		log.Fatal("GobError: Read ", err)
	}

	dec := gob.NewDecoder(bytes.NewBuffer(fGobRead))
	err = dec.Decode(&testSliceGob)
	if err != nil {
		log.Fatal("GobError: Decode ", err)
	}

	endDes := time.Now()

	log.Println("GOB")
	log.Printf("Serialization time: %d ns", endSer.Sub(startSer)/1000000)
	log.Printf("Deserialization time: %d ns", endDes.Sub(startDes)/1000000)
	log.Printf("Overall time: %d ns", (endSer.Sub(startSer)+endDes.Sub(startDes))/1000000)

	log.Println("serialized size in bytes: ", len(gobbytes))
	if len(testSliceGob) != len(testSlice) {
		fmt.Println("bug")
	}

	startSer = time.Now()
	jsonbytes := toJsonBytes()
	ss2 := loadJsonBytes(jsonbytes)
	endSer = time.Now()

	fJsonWrite, err := os.OpenFile("json.txt", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal("JsonError: Open ", err)
	}

	defer fJsonWrite.Close()

	encJson := json.NewEncoder(fJsonWrite)
	if err := encJson.Encode(ss2); err != nil {
		log.Fatal("JsonError: Encode ", err)
	}

	startDes = time.Now()
	fJsonRead, err := ioutil.ReadFile("./json.txt")
	if err != nil {
		log.Fatal("JsonError: Read ", err)
	}

	err = json.Unmarshal(fJsonRead, &testSliceJson)
	if err != nil {
		log.Fatal("JsonError: Unmarshal ", err)
	}
	endDes = time.Now()

	log.Println("JSON")
	log.Printf("Serialization time: %d ns", endSer.Sub(startSer)/1000000)
	log.Printf("Deserialization time: %d ns", endDes.Sub(startDes)/1000000)
	log.Printf("Overall time: %d ns", (endSer.Sub(startSer)+endDes.Sub(startDes))/1000000)

	log.Println("serialized size in bytes: ", len(jsonbytes))

	if len(testSliceJson) != len(testSlice) {
		fmt.Println("bug")
	}

	startSer = time.Now()
	xmlbytes := toXmlBytes()
	ss3 := loadXmlBytes(xmlbytes)
	endSer = time.Now()

	fXmlWrite, err := os.OpenFile("xml.txt", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal("XmlError: Open ", err)
	}

	defer fXmlWrite.Close()

	encXml := xml.NewEncoder(fXmlWrite)
	if err := encXml.Encode(ss3); err != nil {
		log.Fatal("XmlError: Encode ", err)
	}

	startDes = time.Now()
	fXmlRead, err := ioutil.ReadFile("./xml.txt")
	if err != nil {
		log.Fatal("XmlError: Read ", err)
	}

	decxml := xml.NewDecoder(bytes.NewBuffer(fXmlRead))
	for {
		var buff TestStruct
		err := decxml.Decode(&buff)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal("XmlError: Decode ", err)
		}
		testSliceXml = append(testSliceXml, buff)
	}
	endDes = time.Now()

	log.Println("XML")
	log.Printf("Serialization time: %d ns", endSer.Sub(startSer)/1000000)
	log.Printf("Deserialization time: %d ns", endDes.Sub(startDes)/1000000)
	log.Printf("Overall time: %d ns", (endSer.Sub(startSer)+endDes.Sub(startDes))/1000000)

	log.Println("serialized size in bytes: ", len(xmlbytes))
	if len(testSliceXml) != len(testSlice) {
		fmt.Println("bug")
	}

}

func main() {
	fmt.Println("Array size: ", len(testSlice))
	benchmark()
}
