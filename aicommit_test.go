package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"
)

func TestWriteCommitMessage(t *testing.T) {
    // create a temporary file
    f, err := ioutil.TempFile("", "test-file")
    if err != nil {
        t.Fatalf("failed to create temporary file: %v", err)
    }
    defer os.Remove(f.Name())

    // write some content to the file
    originalContent := []byte("Hello, world!\n")
    if _, err := f.Write(originalContent); err != nil {
        t.Fatalf("failed to write to file: %v", err)
    }

    // call the function with a message to append
    newMessage := "Add some new feature"
    if err := WriteCommitMessage(newMessage, f.Name()); err != nil {
        t.Fatalf("writeCommitMessage returned an error: %v", err)
    }

    // read the contents of the file and check if the message was appended
    updatedContent, err := ioutil.ReadFile(f.Name())
    if err != nil {
        t.Fatalf("failed to read file: %v", err)
    }
    expectedContent := append([]byte(newMessage), originalContent...)
    if !bytes.Equal(updatedContent, expectedContent) {
        t.Errorf("file content doesn't match expected value: got %q, want %q", updatedContent, expectedContent)
    }
}