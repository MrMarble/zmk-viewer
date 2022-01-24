package keyboard

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test(t *testing.T) {
	// Start a local HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Send response to be tested
		rw.Write([]byte(`{
			"keyboards": {
			  "cradio": {
				"keyboard_name": "Cradio",
				"layouts": {
				  "LAYOUT_split_3x5_2": {
					"layout": [
					  { "x": 0, "y": 1.27, "w": 1, "label": "L01", "matrix": [0, 0] },
					  { "x": 1, "y": 0.31, "w": 1, "label": "L02", "matrix": [0, 1] },
					  { "x": 2, "y": 0, "w": 1, "label": "L03", "matrix": [0, 2] },
					  { "x": 3, "y": 0.28, "w": 1, "label": "L04", "matrix": [0, 3] },
					  { "x": 4, "y": 0.42, "w": 1, "label": "L05", "matrix": [0, 4] },
					  { "x": 8, "y": 0.42, "w": 1, "label": "R01", "matrix": [4, 0] },
					  { "x": 9, "y": 0.28, "w": 1, "label": "R02", "matrix": [4, 1] },
					  { "x": 10, "y": 0, "w": 1, "label": "R03", "matrix": [4, 2] },
					  { "x": 11, "y": 0.31, "w": 1, "label": "R04", "matrix": [4, 3] },
					  { "x": 12, "y": 1.27, "w": 1, "label": "R05", "matrix": [4, 4] },
					  { "x": 0, "y": 2.27, "w": 1, "label": "L06", "matrix": [1, 0] },
					  { "x": 1, "y": 1.31, "w": 1, "label": "L07", "matrix": [1, 1] },
					  { "x": 2, "y": 1, "w": 1, "label": "L08", "matrix": [1, 2] },
					  { "x": 3, "y": 1.28, "w": 1, "label": "L09", "matrix": [1, 3] },
					  { "x": 4, "y": 1.42, "w": 1, "label": "L10", "matrix": [1, 4] },
					  { "x": 8, "y": 1.42, "w": 1, "label": "R06", "matrix": [5, 0] },
					  { "x": 9, "y": 1.28, "w": 1, "label": "R07", "matrix": [5, 1] },
					  { "x": 10, "y": 1, "w": 1, "label": "R08", "matrix": [5, 2] },
					  { "x": 11, "y": 1.31, "w": 1, "label": "R09", "matrix": [5, 3] },
					  { "x": 12, "y": 2.27, "w": 1, "label": "R10", "matrix": [5, 4] },
					  { "x": 0, "y": 3.27, "w": 1, "label": "L11", "matrix": [2, 0] },
					  { "x": 1, "y": 2.31, "w": 1, "label": "L12", "matrix": [2, 1] },
					  { "x": 2, "y": 2, "w": 1, "label": "L13", "matrix": [2, 2] },
					  { "x": 3, "y": 2.28, "w": 1, "label": "L14", "matrix": [2, 3] },
					  { "x": 4, "y": 2.42, "w": 1, "label": "L15", "matrix": [2, 4] },
					  { "x": 8, "y": 2.42, "w": 1, "label": "R11", "matrix": [6, 0] },
					  { "x": 9, "y": 2.28, "w": 1, "label": "R12", "matrix": [6, 1] },
					  { "x": 10, "y": 2, "w": 1, "label": "R13", "matrix": [6, 2] },
					  { "x": 11, "y": 2.31, "w": 1, "label": "R14", "matrix": [6, 3] },
					  { "x": 12, "y": 3.27, "w": 1, "label": "R15", "matrix": [6, 4] },
					  { "x": 4, "y": 3.9, "w": 1, "label": "L16", "matrix": [3, 0] },
					  { "x": 5, "y": 3.7, "w": 1, "label": "L17", "matrix": [3, 1] },
					  { "x": 7, "y": 3.7, "w": 1, "label": "R16", "matrix": [7, 0] },
					  { "x": 8, "y": 3.9, "w": 1, "label": "R17", "matrix": [7, 1] }
					]
				  }
				}
			  }
			}
		  }
		  `))
	}))
	// Close the server when test finishes
	defer server.Close()

	file, err := fetch(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	_, ok := file.Keyboards["cradio"]
	if !ok {
		t.Fatal("Unmarshaling error")
	}
}
