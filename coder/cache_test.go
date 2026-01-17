package coder

import (
	"testing"
)

func TestCacheBasic(t *testing.T) {
	filename1 := "example1.mc"
	filename2 := "example2.mc"
	doc1, _ := ParseDocument([]byte(""), filename1)
	doc2, _ := ParseDocument([]byte(""), filename2)

	cache := NewCache()
	cache.Add(filename1, doc1)
	cache.Add(filename2, doc2)

	retrievedDoc1, ok1 := cache.Get(filename1)
	if !ok1 {
		t.Fatalf("Failed to get document '%s' from cache", filename1)
	}
	if retrievedDoc1 != doc1 {
		t.Fatalf("Retrieved document '%s' does not match the original", filename1)
	}

	retrievedDoc2, ok2 := cache.Get(filename2)
	if !ok2 {
		t.Fatalf("Failed to get document '%s' from cache", filename2)
	}
	if retrievedDoc2 != doc2 {
		t.Fatalf("Retrieved document '%s' does not match the original", filename2)
	}

	filenameNotExist := "nonexistent.mc"
	_, ok3 := cache.Get(filenameNotExist)
	if ok3 {
		t.Fatalf("Expected to not find document '%s' in cache, but it was found", filenameNotExist)
	}
}
