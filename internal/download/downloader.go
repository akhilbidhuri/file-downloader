package download

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
)

// 1MB chunks
const chunkSize = 1024 * 1024

func getRangeData(u string) (bool, int64, error) {

	request, _ := http.NewRequest("HEAD", u, nil)

	client := http.Client{}

	response, err := client.Do(request)
	if err != nil {
		return false, 0, fmt.Errorf("failed HTTP head, err : %v", err)
	}
	defer response.Body.Close()

	statusCode, headers := response.StatusCode, response.Header

	if statusCode != 200 && statusCode != 206 {
		return false, 0, fmt.Errorf("failed, unsuccessful response: %v", statusCode)
	}

	contentLength := headers.Get("Content-Length")
	contentLengthInt, err := strconv.ParseInt(contentLength, 10, 64)
	if err != nil {
		return false, 0, fmt.Errorf("could not find content length, err : %v", err)
	}
	if headers.Get("Accept-Ranges") == "bytes" {
		return true, contentLengthInt, nil
	}

	return false, contentLengthInt, nil
}

func downloadChunk(url string, start int64, end int64, ch chan<- Chunk) {
	req, _ := http.NewRequest("GET", url, nil)
	rangeHeader := fmt.Sprintf("bytes=%d-%d", start, end)
	req.Header.Add("Range", rangeHeader)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error downloading chunk:", err)
		return
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading chunk:", err)
		return
	}
	fmt.Printf("got chunk %d size %d\n", int(start/chunkSize), len(data))
	ch <- Chunk{
		SeqNum: int(start / chunkSize),
		Data:   data,
	}
}

func Process(url string, outputPath string) error {
	acceptsRange, size, err := getRangeData(url)
	if err != nil {
		panic(err)
	}

	if size == -1 {
		panic("Content-Length header not provided by server")
	}

	numChunks := size / chunkSize
	if size%chunkSize != 0 {
		numChunks++
	}

	if !acceptsRange {
		numChunks = 1
	}

	var wg sync.WaitGroup
	ch := make(chan Chunk, numChunks)

	for i := 0; i < int(numChunks); i++ {
		start := int64(i) * chunkSize
		end := start + chunkSize - 1
		if i == int(numChunks)-1 {
			end = size - 1
		}

		wg.Add(1)
		fmt.Printf("starting thread: %d\n", i)
		go func(start, end int64) {
			defer wg.Done()
			downloadChunk(url, start, end, ch)
		}(start, end)
	}

	wg.Wait()
	close(ch)

	chunks := make([]Chunk, 0, numChunks)
	for chunk := range ch {
		chunks = append(chunks, chunk)
	}

	sort.Slice(chunks, func(i, j int) bool {
		return chunks[i].SeqNum < chunks[j].SeqNum
	})

	if outputPath == "" {
		splitURL := strings.Split(url, "/")
		outputPath = splitURL[len(splitURL)-1]
	}

	outputFile, err := os.Create(outputPath)
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()

	for _, chunk := range chunks {
		fmt.Printf("storing chunk: %d\n", chunk.SeqNum)
		_, err := outputFile.Write(chunk.Data)
		if err != nil {
			panic(err)
		}
	}

	fmt.Printf("File downloaded and saved to: %s\n", outputPath)
	return nil
}
