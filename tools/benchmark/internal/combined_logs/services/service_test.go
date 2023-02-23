package services

import (
	"os"
	"testing"
)

func TestGet(t *testing.T) {
	{
		content := "FN:xyz.log\nHey now\nbrown cow\nFN:lorem.log\nHello world\nHey there\nFN:hello_something.log\nLorem ipsum 123\nawfawfawf"
		err := os.WriteFile("all.clog", []byte(content), 0644)
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove("all.clog")
	}

	combinedLogSvc := CombinedLogsService{}
	err := combinedLogSvc.Load()
	if err != nil {
		t.Fatal(err)
	}

	{
		content_, err := combinedLogSvc.Get("xyz.log")
		if err != nil {
			t.Fatal(err)
		}
		content, exists := content_.Get()
		if !exists {
			t.Fatal("Expected combined log entry but found none")
		}

		expectedEntry := "Hey now\nbrown cow\n"
		if content != expectedEntry {
			t.Fatalf("Expected '%s' but got '%s'", expectedEntry, content)
		}
	}
	{
		content_, err := combinedLogSvc.Get("lorem.log")
		if err != nil {
			t.Fatal(err)
		}
		content, exists := content_.Get()
		if !exists {
			t.Fatal("Expected combined log entry but found none")
		}

		expectedEntry := "Hello world\nHey there\n"
		if content != expectedEntry {
			t.Fatalf("Expected '%s' but got '%s'", expectedEntry, content)
		}
	}
	{
		content_, err := combinedLogSvc.Get("hello_something.log")
		if err != nil {
			t.Fatal(err)
		}
		content, exists := content_.Get()
		if !exists {
			t.Fatal("Expected combined log entry but found none")
		}

		expectedEntry := "Lorem ipsum 123\nawfawfawf"
		if content != expectedEntry {
			t.Fatalf("Expected '%s' but got '%s'", expectedEntry, content)
		}
	}
}
