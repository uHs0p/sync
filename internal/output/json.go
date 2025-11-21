package output

import (
    "encoding/json"
    "terraformsync/internal/compare"
)

func ToJSON(r compare.Result) ([]byte, error) {
    return json.MarshalIndent(r, "", "  ")
}