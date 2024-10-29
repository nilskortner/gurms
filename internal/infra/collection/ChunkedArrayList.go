package collection

import (
	"errors"
	"fmt"
	"gurms/internal/infra/0goutil/pow2"
)

const DEFAULT_CHUNK_CAPACITY = 32
const DEFAULT_INITAL_CHUNKS = 8

type ChunkedArrayList struct {
	elementSizePerChunk     int
	elementSizePerChunkMask int
	chunks                  [][]int
	elementCount            int
}

func NewChunkedArrayListWithInitalChunks(elementSizePerChunk, initialChunks int) *ChunkedArrayList {
	var c = &ChunkedArrayList{}
	c.elementSizePerChunk = pow2.RoundToPowerOfTwo(elementSizePerChunk)
	c.elementSizePerChunkMask = c.elementSizePerChunk - 1
	c.chunks = make([][]int, initialChunks)
	c.elementCount = 0
	return c
}

func NewChunkedArrayListWithElementSize(elementSizePerChunk int) *ChunkedArrayList {
	return NewChunkedArrayListWithInitalChunks(elementSizePerChunk, DEFAULT_INITAL_CHUNKS)
}

func NewChunkedArrayList() *ChunkedArrayList {
	return NewChunkedArrayListWithInitalChunks(DEFAULT_CHUNK_CAPACITY, DEFAULT_INITAL_CHUNKS)
}

func (c *ChunkedArrayList) Size() int {
	return c.elementCount
}

func (c *ChunkedArrayList) Contains(element int) bool {
	for _, chunk := range c.chunks {
		for _, item := range chunk {
			if item == element {
				return true
			}
		}
	}
	return false
}

func (c *ChunkedArrayList) Add(element int) bool {
	var chunkCount = len(c.chunks)
	var chunk = c.chunks[chunkCount-1]
	if chunkCount == 0 || (len(chunk) == c.elementSizePerChunk) {
		chunk = make([]int, c.elementSizePerChunk)
		chunk[0] = element
		c.chunks = append(c.chunks, chunk)
		c.elementCount++
		return true
	}
	chunk = append(chunk, element)
	c.chunks[chunkCount-1] = chunk
	c.elementCount++
	return true
}

func (c *ChunkedArrayList) Remove(element int) bool {
	var chunkIndex = 0
	var chunkCount = len(c.chunks)
	var chunk []int
	var nextChunk []int
	var err error
	var boolean bool
	for i := 0; chunkIndex < chunkCount; i++ {
		chunk, err = getChunk(c.chunks, chunkIndex)
		if err != nil {
			fmt.Println(err)
		}
		chunk, boolean = removeElement(chunk, element)
		if boolean {
			if len(chunk) == 0 {
				c.chunks, err = removeChunk(c.chunks, chunkIndex)
				if err != nil {
					fmt.Println(err)
				}
			} else {
				for chunkIndex+1 < chunkCount {
					chunkIndex += 1
					nextChunk, err = getChunk(c.chunks, chunkIndex)
					if err != nil {
						fmt.Println(err)
					}
					var value int
					nextChunk, value, err = removeValue(nextChunk, 0)
					if err != nil {
						fmt.Println(err)
					}
					chunk = append(chunk, value)
					if len(nextChunk) == 0 {
						c.chunks, err = removeChunk(c.chunks, chunkIndex)
						if err != nil {
							fmt.Println(err)
						}
						break
					}
					c.chunks[chunkIndex] = nextChunk
				}
			}
			c.elementCount--
			return true
		}
	}
	return false
}

func (c *ChunkedArrayList) Clear() {
	for i := range c.chunks {
		c.chunks[i] = nil
	}
	c.chunks = nil
	c.elementCount = 0
}

func (c *ChunkedArrayList) Get(index int) int {
	chunk, err := getChunk(c.chunks, index/c.elementSizePerChunk)
	if err != nil {
		fmt.Println(err)
	}
	value, err := getValue(chunk, index&c.elementSizePerChunkMask)
	if err != nil {
		fmt.Println(err)
	}
	return value
}

func (c *ChunkedArrayList) Set(index int, element int) (int, error) {
	if index < 0 || index >= c.elementCount {
		return *new(int), errors.New("index out of bounds")
	}

	chunk, err := getChunk(c.chunks, index/c.elementSizePerChunk)
	if err != nil {
		fmt.Println(err)
	}

	elementIndex := index & c.elementSizePerChunkMask
	existingElement, err := getValue(chunk, elementIndex)
	if err != nil {
		fmt.Println(err)
	}

	c.elementCount++

	chunk[elementIndex] = element

	return existingElement, nil
}

func (c *ChunkedArrayList) AddIndex(index int, element int) {
	chunkIndex := index / c.elementSizePerChunk
	elementIndex := index & c.elementSizePerChunkMask
	var chunk []int
	var err error
	chunkCount := len(c.chunks)
	for chunkIndex < chunkCount {
		chunk, err = getChunk(c.chunks, chunkIndex)
		if err != nil {
			fmt.Println(err)
		}
		if len(chunk) < c.elementSizePerChunk {
			chunk[elementIndex] = element
			c.elementCount++
			return
		}
		chunk, removedElement, err := removeValue(chunk, c.elementSizePerChunkMask)
		if err != nil {
			fmt.Println(err)
		}
		chunk[elementIndex] = element
		element = removedElement
		chunkIndex++
		elementIndex = 0
	}
	chunk = make([]int, c.elementSizePerChunk)
	chunk[0] = element
	c.chunks = append(c.chunks, chunk)
	c.elementCount++
}

func (c *ChunkedArrayList) RemoveIndex(index int) int {
	var value int
	var removedElement int
	var err error
	chunkIndex := index / c.elementSizePerChunk
	elementIndex := index & c.elementSizePerChunkMask
	chunk, err := getChunk(c.chunks, chunkIndex)
	if err != nil {
		fmt.Println(err)
	}
	chunk, removedElement, err = removeValue(chunk, elementIndex)
	if err != nil {
		fmt.Println(err)
	}
	c.elementCount--
	var nextChunk []int
	if len(chunk) == 0 {
		c.chunks, err = removeChunk(c.chunks, chunkIndex)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		chunkCount := len(c.chunks)
		for chunkIndex+1 < chunkCount {
			chunkIndex += 1
			nextChunk, err = getChunk(c.chunks, chunkIndex)
			if err != nil {
				fmt.Println(err)
			}
			nextChunk, value, err = removeValue(nextChunk, 0)
			if err != nil {
				fmt.Println(err)
				chunk = append(chunk, value)
				if len(nextChunk) == 0 {
					c.chunks, err = removeChunk(c.chunks, chunkIndex)
					if err != nil {
						fmt.Println(err)
					}
					break
				}
				c.chunks[chunkIndex] = nextChunk
			}
		}
	}
	return removedElement
}

func (c *ChunkedArrayList) IndexOf(element int) int {
	var chunk []int
	var result int
	var err error
	chunkCount := len(c.chunks)
	for i := 0; i < chunkCount; i++ {
		chunk, err = getChunk(c.chunks, i)
		if err != nil {
			fmt.Println(err)
		}
		result = indexOf(chunk, element)
		if result >= 0 {
			return i*c.elementSizePerChunk + result
		}
	}
	return -1
}

func (c *ChunkedArrayList) LastIndexOf() {

}

// region helper functions
func getChunk(slice [][]int, index int) ([]int, error) {
	if index < 0 || index >= len(slice) {
		return slice[0], errors.New("index out of bounds")
	}
	return slice[index], nil
}

func getValue(slice []int, index int) (int, error) {
	if index < 0 || index >= len(slice) {
		return 0, errors.New("index out of bounds")
	}
	return slice[index], nil
}

func removeChunk(slice [][]int, index int) ([][]int, error) {
	if index < 0 || index >= len(slice) {
		return slice, errors.New("index out of bounds")
	}
	slice = append(slice[:index], slice[index+1:]...)
	return slice, nil
}

func removeValue(slice []int, index int) ([]int, int, error) {
	if index < 0 || index >= len(slice) {
		return slice, 0, errors.New("index out of bounds")
	}
	value := slice[index]
	slice = append(slice[:index], slice[index+1:]...)
	return slice, value, nil
}

func removeElement(slice []int, element int) ([]int, bool) {
	for i := 0; i < len(slice); i++ {
		if slice[i] == element {
			slice = append(slice[:i], slice[i+1:]...)
			slice = slice[:len(slice)-1]
			return slice, true
		}
	}
	return slice, false
}

func indexOf(slice []int, element int) int {
	for i := 0; i < len(slice); i++ {
		if slice[i] == element {
			return i
		}
	}
	return -1
}

// func lastIndexOf(slice []int, element int) int {

// }

// endregion
